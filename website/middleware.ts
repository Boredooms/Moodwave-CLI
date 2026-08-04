import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

// ─────────────────────────────────────────────────────────────────────────────
// LAUNCH GATE — server-side middleware that blocks every route on the website
// until the unlock timestamp. This runs on the edge before any page renders,
// so no HTML, no JSON, no static asset for protected pages can be returned
// before the launch date regardless of what URL, slug, or query parameter
// someone tries. Not bypassable from the client: the decision is made before
// the response exists.
//
// Unlock: September 3, 2026 at 06:00 AM IST (UTC+05:30) → 00:30 UTC
// ─────────────────────────────────────────────────────────────────────────────

const UNLOCK_UTC = new Date("2026-09-03T00:30:00Z"); // 06:00 AM IST

// Paths that are ALWAYS accessible (the launch page itself, its assets, and
// the waitlist API route). Everything else is gated.
const ALLOWED = [
  "/launch",        // the countdown page
  "/api/waitlist",  // the email submission endpoint
  "/_next",         // Next.js internal assets (JS chunks, images, etc.)
  "/favicon.ico",
  "/icon.svg",
  "/logo.svg",
];

export function middleware(request: NextRequest) {
  const now = new Date();

  // After the unlock timestamp, the gate disappears forever — no redirect,
  // no launch page, nothing. The middleware becomes a transparent pass-through.
  if (now >= UNLOCK_UTC) {
    return NextResponse.next();
  }

  const { pathname } = request.nextUrl;

  // Allow the launch page itself and its dependencies through.
  if (ALLOWED.some((prefix) => pathname.startsWith(prefix))) {
    return NextResponse.next();
  }

  // Everything else → redirect to the launch/countdown page.
  const url = request.nextUrl.clone();
  url.pathname = "/launch";
  return NextResponse.redirect(url);
}

// Match all routes. The function above decides which ones pass through.
export const config = {
  matcher: ["/((?!_next/static|_next/image).*)"],
};
