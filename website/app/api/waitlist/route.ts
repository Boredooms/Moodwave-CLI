import { NextResponse } from "next/server";
import { promises as fs } from "fs";
import path from "path";
import { headers } from "next/headers";

// ─────────────────────────────────────────────────────────────────────────────
// /api/waitlist
//
// POST — Submit an email to the waitlist (rate-limited, validated, deduped)
// GET  — Returns the current real waitlist count (for the live counter on the
//         launch page)
//
// Anti-fake protections:
// 1. Rate limit: max 3 submissions per IP per hour (in-memory, resets on cold start)
// 2. Disposable email domain blocklist (mailinator, guerrillamail, etc.)
// 3. Email format validation (RFC-ish regex)
// 4. Duplicate email rejection (same email can't register twice)
// 5. IP logged per entry so you can audit/purge bot clusters later
//
// For production persistence, swap the JSON file with a real database (Supabase,
// Turso, Vercel KV, etc.). The Resend integration is ready — just set
// RESEND_API_KEY in your environment to enable confirmation emails.
// ─────────────────────────────────────────────────────────────────────────────

const WAITLIST_FILE =
  process.env.NODE_ENV === "production"
    ? "/tmp/waitlist.json"
    : path.join(process.cwd(), "waitlist.json");

// ── Rate limiter (in-memory, per-instance) ──────────────────────────────────

const RATE_WINDOW_MS = 60 * 60 * 1000; // 1 hour
const RATE_LIMIT = 3; // max signups per IP per window

const ipHits: Map<string, { count: number; windowStart: number }> = new Map();

function isRateLimited(ip: string): boolean {
  const now = Date.now();
  const entry = ipHits.get(ip);

  if (!entry || now - entry.windowStart > RATE_WINDOW_MS) {
    ipHits.set(ip, { count: 1, windowStart: now });
    return false;
  }

  if (entry.count >= RATE_LIMIT) {
    return true;
  }

  entry.count++;
  return false;
}

// ── Disposable email domain blocklist ───────────────────────────────────────

const DISPOSABLE_DOMAINS = new Set([
  "mailinator.com", "guerrillamail.com", "guerrillamail.net", "tempmail.com",
  "throwaway.email", "yopmail.com", "sharklasers.com", "grr.la", "guerrillamail.info",
  "guerrillamail.biz", "guerrillamailblock.com", "pokemail.net", "spam4.me",
  "trashmail.com", "trashmail.me", "trashmail.net", "dispostable.com",
  "maildrop.cc", "fakeinbox.com", "emailondeck.com", "getnada.com",
  "temp-mail.org", "tempail.com", "tempr.email", "10minutemail.com",
  "minutemail.com", "mohmal.com", "burnermail.io", "harakirimail.com",
  "mailnesia.com", "mailsac.com", "mytemp.email", "throwawaymail.com",
  "tmail.ws", "tmails.net", "tmpmail.net", "tmpmail.org", "trash-mail.at",
  "trashmail.ws", "wegwerfmail.de", "wegwerfmail.net",
]);

function isDisposableEmail(email: string): boolean {
  const domain = email.split("@")[1]?.toLowerCase();
  return DISPOSABLE_DOMAINS.has(domain);
}

// ── Waitlist data ───────────────────────────────────────────────────────────

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

// ── Resend confirmation email (optional — requires RESEND_API_KEY env) ──────

async function sendConfirmationEmail(email: string): Promise<void> {
  const apiKey = process.env.RESEND_API_KEY;
  if (!apiKey) return; // Skip if not configured

  const fromAddress = process.env.RESEND_FROM || "Moodwave <noreply@moodwave.dev>";

  try {
    await fetch("https://api.resend.com/emails", {
      method: "POST",
      headers: {
        Authorization: `Bearer ${apiKey}`,
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        from: fromAddress,
        to: [email],
        subject: "You're on the Moodwave waitlist 🌊",
        html: `
          <div style="font-family: 'JetBrains Mono', monospace; background: #080808; color: #ffffff; padding: 40px; max-width: 500px; margin: 0 auto;">
            <h1 style="font-size: 24px; margin-bottom: 16px;">moodwave</h1>
            <p style="color: #888; font-size: 14px; line-height: 1.6;">
              You're on the list. We'll email you the moment Moodwave goes live on September 3, 2026.
            </p>
            <p style="color: #555; font-size: 12px; margin-top: 32px;">
              A terminal-native music companion that reads your code, detects your mood, and plays the perfect soundtrack.
            </p>
            <hr style="border: none; border-top: 1px solid #222; margin: 24px 0;" />
            <p style="color: #444; font-size: 10px;">
              You received this because you signed up at moodwave.dev
            </p>
          </div>
        `,
      }),
    });
  } catch {
    // Non-critical — don't fail the signup if the email doesn't send
  }
}

// ── POST handler ────────────────────────────────────────────────────────────

export async function POST(request: Request) {
  try {
    // Get client IP
    const headersList = await headers();
    const ip =
      headersList.get("x-forwarded-for")?.split(",")[0]?.trim() ||
      headersList.get("x-real-ip") ||
      "unknown";

    // Rate limit check
    if (isRateLimited(ip)) {
      return NextResponse.json(
        { error: "Too many signups from this address. Try again later." },
        { status: 429 }
      );
    }

    const body = await request.json();
    const email = body?.email?.trim()?.toLowerCase();

    // Validate format
    if (!email || !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
      return NextResponse.json(
        { error: "Please provide a valid email address." },
        { status: 400 }
      );
    }

    // Block disposable emails
    if (isDisposableEmail(email)) {
      return NextResponse.json(
        { error: "Disposable email addresses are not allowed. Use a real email." },
        { status: 400 }
      );
    }

    // Check duplicates
    const list = await readWaitlist();
    if (list.some((entry) => entry.email === email)) {
      return NextResponse.json(
        { message: "You're already on the list!", count: list.length, alreadyExists: true },
        { status: 200 }
      );
    }

    // Block multiple emails from same IP (beyond rate limit — cap at 3 total)
    const ipEntries = list.filter((entry) => entry.ip === ip);
    if (ipEntries.length >= 3) {
      return NextResponse.json(
        { error: "Maximum signups reached for this network." },
        { status: 429 }
      );
    }

    // Store
    list.push({ email, ip, timestamp: new Date().toISOString() });
    await writeWaitlist(list);

    // Send confirmation (non-blocking, best-effort)
    sendConfirmationEmail(email);

    return NextResponse.json(
      { message: "You're in! We'll notify you on launch day.", count: list.length },
      { status: 201 }
    );
  } catch {
    return NextResponse.json(
      { error: "Something went wrong. Try again." },
      { status: 500 }
    );
  }
}

// ── GET handler — returns live count ────────────────────────────────────────

export async function GET() {
  try {
    const list = await readWaitlist();
    return NextResponse.json({ count: list.length }, { status: 200 });
  } catch {
    return NextResponse.json({ count: 0 }, { status: 200 });
  }
}
