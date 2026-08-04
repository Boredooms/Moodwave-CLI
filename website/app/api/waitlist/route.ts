import { NextResponse } from "next/server";
import { promises as fs } from "fs";
import path from "path";
import { headers } from "next/headers";

// ─────────────────────────────────────────────────────────────────────────────
// /api/waitlist
//
// POST — Submit an email to the waitlist (rate-limited, validated, deduped)
// GET  — Returns the current real waitlist count (live counter on launch page)
//
// Anti-fake protections:
// 1. Rate limit: max 3 submissions per IP per hour (in-memory)
// 2. Disposable email domain blocklist
// 3. Email format validation
// 4. Duplicate email rejection
// 5. IP hard cap: max 3 registrations per IP ever (blocks bot farms)
// 6. IP logged per entry for audit
//
// No external email service. Emails are stored — export on launch day.
// For production persistence: swap JSON file with Vercel KV / Supabase.
// ─────────────────────────────────────────────────────────────────────────────

const WAITLIST_FILE =
  process.env.NODE_ENV === "production"
    ? "/tmp/waitlist.json"
    : path.join(process.cwd(), "waitlist.json");

// ── Rate limiter ────────────────────────────────────────────────────────────

const RATE_WINDOW_MS = 60 * 60 * 1000;
const RATE_LIMIT = 3;
const ipHits: Map<string, { count: number; windowStart: number }> = new Map();

function isRateLimited(ip: string): boolean {
  const now = Date.now();
  const entry = ipHits.get(ip);
  if (!entry || now - entry.windowStart > RATE_WINDOW_MS) {
    ipHits.set(ip, { count: 1, windowStart: now });
    return false;
  }
  if (entry.count >= RATE_LIMIT) return true;
  entry.count++;
  return false;
}

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

// ── Data ────────────────────────────────────────────────────────────────────

interface WaitlistEntry {
  email: string;
  ip: string;
  timestamp: string;
}

async function readWaitlist(): Promise<WaitlistEntry[]> {
  try {
    const data = await fs.readFile(WAITLIST_FILE, "utf-8");
    return JSON.parse(data);
  } catch {
    return [];
  }
}

async function writeWaitlist(entries: WaitlistEntry[]): Promise<void> {
  await fs.writeFile(WAITLIST_FILE, JSON.stringify(entries, null, 2));
}

// ── POST ────────────────────────────────────────────────────────────────────

export async function POST(request: Request) {
  try {
    const headersList = await headers();
    const ip =
      headersList.get("x-forwarded-for")?.split(",")[0]?.trim() ||
      headersList.get("x-real-ip") ||
      "unknown";

    if (isRateLimited(ip)) {
      return NextResponse.json(
        { error: "Too many attempts. Try again in an hour." },
        { status: 429 }
      );
    }

    const body = await request.json();
    const email = body?.email?.trim()?.toLowerCase();

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

    const list = await readWaitlist();

    if (list.some((entry) => entry.email === email)) {
      return NextResponse.json(
        { message: "You're already on the list!", count: list.length, alreadyExists: true },
        { status: 200 }
      );
    }

    if (list.filter((entry) => entry.ip === ip).length >= 3) {
      return NextResponse.json(
        { error: "Maximum signups reached for this network." },
        { status: 429 }
      );
    }

    list.push({ email, ip, timestamp: new Date().toISOString() });
    await writeWaitlist(list);

    return NextResponse.json(
      { message: "You're in! We'll let you know on launch day.", count: list.length },
      { status: 201 }
    );
  } catch {
    return NextResponse.json(
      { error: "Something went wrong. Try again." },
      { status: 500 }
    );
  }
}

// ── GET — live count ────────────────────────────────────────────────────────

export async function GET() {
  try {
    const list = await readWaitlist();
    return NextResponse.json({ count: list.length }, { status: 200 });
  } catch {
    return NextResponse.json({ count: 0 }, { status: 200 });
  }
}
