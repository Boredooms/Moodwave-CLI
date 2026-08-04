import { NextResponse } from "next/server";
import { promises as fs } from "fs";
import path from "path";

// ─────────────────────────────────────────────────────────────────────────────
// POST /api/waitlist — collects emails for the launch waitlist.
//
// In production on Vercel, writes are ephemeral (the filesystem is read-only
// except /tmp). For a real persistent store, swap this with a database call
// (e.g. Supabase, Turso, or a KV store). For the launch gate MVP this writes
// to a local JSON file during development and to /tmp in serverless — the
// important thing is the endpoint exists and validates input, so the frontend
// flow works end to end.
// ─────────────────────────────────────────────────────────────────────────────

const WAITLIST_FILE =
  process.env.NODE_ENV === "production"
    ? "/tmp/waitlist.json"
    : path.join(process.cwd(), "waitlist.json");

async function readWaitlist(): Promise<string[]> {
  try {
    const data = await fs.readFile(WAITLIST_FILE, "utf-8");
    return JSON.parse(data);
  } catch {
    return [];
  }
}

async function writeWaitlist(emails: string[]): Promise<void> {
  await fs.writeFile(WAITLIST_FILE, JSON.stringify(emails, null, 2));
}

export async function POST(request: Request) {
  try {
    const body = await request.json();
    const email = body?.email?.trim()?.toLowerCase();

    if (!email || !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
      return NextResponse.json(
        { error: "Please provide a valid email address." },
        { status: 400 }
      );
    }

    const list = await readWaitlist();

    if (list.includes(email)) {
      return NextResponse.json(
        { message: "You're already on the list!", alreadyExists: true },
        { status: 200 }
      );
    }

    list.push(email);
    await writeWaitlist(list);

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
