import { NextResponse } from "next/server";
import { headers } from "next/headers";
import { createClient } from "redis";

// ─────────────────────────────────────────────────────────────────────────────
// /api/waitlist — Persistent waitlist backed by Vercel Redis (redis-pink-dog)
//
// Uses the `redis` package with REDIS_URL env var. Globally consistent,
// survives deployments, every device sees the same real-time count.
//
// POST — Submit an email (rate-limited, validated, deduped)
// GET  — Returns the live count
// ─────────────────────────────────────────────────────────────────────────────

// ── Redis client (singleton, lazy-connected) ────────────────────────────────

let redisClient: ReturnType<typeof createClient> | null = null;

async function getRedis() {
  if (!redisClient) {
    const url = process.env.REDIS_URL;
    if (!url) {
      throw new Error("REDIS_URL environment variable is not set");
    }
    redisClient = createClient({ url });
    redisClient.on("error", (err) => console.error("Redis error:", err));
    await redisClient.connect();
  }
  // Reconnect if disconnected (serverless cold start edge case)
  if (!redisClient.isOpen) {
    await redisClient.connect();
  }
  return redisClient;
}

// ── Config ──────────────────────────────────────────────────────────────────

const RATE_LIMIT = 3;
const IP_MAX_TOTAL = 3;
const RATE_WINDOW_SECS = 3600;

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
    const redis = await getRedis();
    const headersList = await headers();
    const ip =
      headersList.get("x-forwarded-for")?.split(",")[0]?.trim() ||
      headersList.get("x-real-ip") ||
      "unknown";

    // Rate limit
    const rateLimitKey = `waitlist:ip:rate:${ip}`;
    const currentRate = await redis.get(rateLimitKey);
    if (currentRate && parseInt(currentRate) >= RATE_LIMIT) {
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
    const alreadyExists = await redis.sIsMember("waitlist:emails", email);
    if (alreadyExists) {
      const count = await redis.sCard("waitlist:emails");
      return NextResponse.json(
        { message: "You're already on the list!", count, alreadyExists: true },
        { status: 200 }
      );
    }

    // IP hard cap
    const ipTotalKey = `waitlist:ip:total:${ip}`;
    const ipTotal = await redis.get(ipTotalKey);
    if (ipTotal && parseInt(ipTotal) >= IP_MAX_TOTAL) {
      return NextResponse.json(
        { error: "Maximum signups reached for this network." },
        { status: 429 }
      );
    }

    // Store
    await redis.sAdd("waitlist:emails", email);
    await redis.rPush("waitlist:entries", JSON.stringify({
      email,
      ip,
      timestamp: new Date().toISOString(),
    }));

    // Increment rate limit + total
    await redis.incr(rateLimitKey);
    await redis.expire(rateLimitKey, RATE_WINDOW_SECS);
    await redis.incr(ipTotalKey);

    const count = await redis.sCard("waitlist:emails");

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

// ── GET — live count ────────────────────────────────────────────────────────

export async function GET() {
  try {
    const redis = await getRedis();
    const count = await redis.sCard("waitlist:emails");
    return NextResponse.json({ count: count || 0 }, { status: 200 });
  } catch {
    return NextResponse.json({ count: 0 }, { status: 200 });
  }
}
