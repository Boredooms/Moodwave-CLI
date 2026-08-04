import { NextResponse } from "next/server";
import { headers } from "next/headers";
import { kv } from "@vercel/kv";

// ─────────────────────────────────────────────────────────────────────────────
// /api/waitlist — Persistent waitlist backed by Vercel KV (Redis)
//
// POST — Submit an email (rate-limited, validated, deduped)
// GET  — Returns the live count (same across all instances/devices)
//
// Data model in KV:
//   "waitlist:emails"     → Redis Set of all registered emails
//   "waitlist:entries"    → Redis List of JSON-stringified {email, ip, timestamp}
//   "waitlist:ip:<ip>"   → Counter of signups from this IP (expires in 1h)
//
// Every device sees the same real-time count because KV is shared global
// state, not per-instance /tmp files. Data survives deployments, cold starts,
// and multi-instance scaling.
// ─────────────────────────────────────────────────────────────────────────────

// ── Config ──────────────────────────────────────────────────────────────────

const RATE_LIMIT = 3;           // max signups per IP per hour
const IP_MAX_TOTAL = 3;         // hard cap per IP ever
const RATE_WINDOW_SECS = 3600;  // 1 hour TTL on rate limit key

// ── Disposable email blocklist ──────────────────────────────────────────────

const DISPOSABLE_DOMAINS = new Set([
  "mailinator.com", "guerrillamail.com", "guerrillamail.net", "tempmail.com",
  "throwaway.email", "yopmail.com", "sharklasers.com", "grr.la",
  "guerrillamail.info", "guerrillamail.biz", "guerrillamailblock.com",
  "pokemail.net", "spam4.me", "trashmail.com", "trashmail.me", "trashmail.net",
  "dispostable.com", "maildrop.cc", "fakeinbox.com", "emailondeck.com",
  "getnada.com", "temp-mail.org", "tempail.com", "tempr.email",
  "10minutemail.com", "minutemail.com", "mohmal.com", "burnermail.io",
  "harakirimail.com", "mailnesia.com", "mailsac.com", "mytemp.email",
  "throwawaymail.com", "tmail.ws", "tmails.net", "tmpmail.net", "tmpmail.org",
  "trash-mail.at", "trashmail.ws", "wegwerfmail.de", "wegwerfmail.net",
]);

function isDisposableEmail(email: string): boolean {
  const domain = email.split("@")[1]?.toLowerCase();
  return DISPOSABLE_DOMAINS.has(domain);
}

// ── POST ────────────────────────────────────────────────────────────────────

export async function POST(request: Request) {
  try {
    const headersList = await headers();
    const ip =
      headersList.get("x-forwarded-for")?.split(",")[0]?.trim() ||
      headersList.get("x-real-ip") ||
      "unknown";

    // Rate limit: check IP counter (auto-expires after 1 hour)
    const rateLimitKey = `waitlist:ip:rate:${ip}`;
    const currentRate = (await kv.get<number>(rateLimitKey)) || 0;
    if (currentRate >= RATE_LIMIT) {
      return NextResponse.json(
        { error: "Too many attempts. Try again in an hour." },
        { status: 429 }
      );
    }

    const body = await request.json();
    const email = body?.email?.trim()?.toLowerCase();

    // Validate
    if (!email || !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
      return NextResponse.json(
        { error: "Please enter a valid email address." },
        { status: 400 }
      );
    }

    if (isDisposableEmail(email)) {
      return NextResponse.json(
        { error: "Temporary email addresses are not allowed." },
        { status: 400 }
      );
    }

    // Check duplicate
    const alreadyExists = await kv.sismember("waitlist:emails", email);
    if (alreadyExists) {
      const count = await kv.scard("waitlist:emails");
      return NextResponse.json(
        { message: "You're already on the list!", count, alreadyExists: true },
        { status: 200 }
      );
    }

    // IP hard cap (total ever)
    const ipTotalKey = `waitlist:ip:total:${ip}`;
    const ipTotal = (await kv.get<number>(ipTotalKey)) || 0;
    if (ipTotal >= IP_MAX_TOTAL) {
      return NextResponse.json(
        { error: "Maximum signups reached for this network." },
        { status: 429 }
      );
    }

    // Store the email
    await kv.sadd("waitlist:emails", email);
    await kv.rpush("waitlist:entries", JSON.stringify({
      email,
      ip,
      timestamp: new Date().toISOString(),
    }));

    // Increment rate limit (with TTL) and total IP counter
    await kv.incr(rateLimitKey);
    await kv.expire(rateLimitKey, RATE_WINDOW_SECS);
    await kv.incr(ipTotalKey);

    const count = await kv.scard("waitlist:emails");

    return NextResponse.json(
      { message: "You're in! We'll let you know on launch day.", count },
      { status: 201 }
    );
  } catch (err) {
    console.error("Waitlist POST error:", err);
    return NextResponse.json(
      { error: "Something went wrong. Try again." },
      { status: 500 }
    );
  }
}

// ── GET — live count (globally consistent) ──────────────────────────────────

export async function GET() {
  try {
    const count = await kv.scard("waitlist:emails");
    return NextResponse.json({ count: count || 0 }, { status: 200 });
  } catch {
    return NextResponse.json({ count: 0 }, { status: 200 });
  }
}
