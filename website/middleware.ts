import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

// ─────────────────────────────────────────────────────────────────────────────
// LAUNCH GATE — server-side middleware that gates human browser traffic to
// /launch until the launch unlock timestamp.
//
// Search engine crawlers (Googlebot, Bingbot, Google-InspectionTool, GPTBot,
// ClaudeBot, PerplexityBot, etc.) and SEO metadata assets (robots.txt,
// sitemap.xml, llms.txt, etc.) are ALWAYS allowed through so Google Search
// Console and AI engines can crawl, verify, index, and build rich snippets.
//
// Unlock: September 3, 2026 at 06:00 AM IST (UTC+05:30) → 00:30 UTC
// ─────────────────────────────────────────────────────────────────────────────

const UNLOCK_UTC = new Date("2026-09-03T00:30:00Z"); // 06:00 AM IST

// Known search crawler and inspection bot user-agent tokens
const BOT_USER_AGENTS = [
  "googlebot",
  "google-inspectiontool",
  "google-site-verification",
  "google-structured-data-testing-tool",
  "mediapartners-google",
  "adsbot-google",
  "feedfetcher-google",
  "bingbot",
  "bingpreview",
  "msnbot",
  "yandex",
  "duckduckbot",
  "baiduspider",
  "slurp",
  "twitterbot",
  "facebookexternalhit",
  "linkedinbot",
  "whatsapp",
  "telegrambot",
  "discordbot",
  "applebot",
  "gptbot",
  "claudebot",
  "perplexitybot",
  "google-extended",
  "ccbot",
  "bytespider",
  "screaming frog",
  "lighthouse",
];

// Paths that are ALWAYS accessible for everyone (launch page, assets, APIs, SEO files)
const ALLOWED_EXACT_OR_PREFIX = [
  "/launch",
  "/api/waitlist",
  "/_next",
  "/favicon.ico",
  "/icon.svg",
  "/logo.svg",
  "/robots.txt",
  "/sitemap.xml",
  "/sitemap",
  "/llms.txt",
  "/llms-full.txt",
  "/site.webmanifest",
  "/opengraph-image",
  "/twitter-image",
];

export function middleware(request: NextRequest) {
  const now = new Date();

  // 1. After unlock, pass through everything
  if (now >= UNLOCK_UTC) {
    return NextResponse.next();
  }

  const { pathname } = request.nextUrl;

  // 2. Allow explicitly whitelisted routes (launch page, static assets, sitemaps, robots, llms)
  if (ALLOWED_EXACT_OR_PREFIX.some((p) => pathname === p || pathname.startsWith(p))) {
    return NextResponse.next();
  }

  // 3. Allow search engine crawlers and inspection bots to crawl pages for SEO indexing
  const userAgent = (request.headers.get("user-agent") || "").toLowerCase();
  const isSearchCrawler = BOT_USER_AGENTS.some((bot) => userAgent.includes(bot));
  if (isSearchCrawler) {
    return NextResponse.next();
  }

  // 4. Redirect human browser traffic to the launch countdown page
  const url = request.nextUrl.clone();
  url.pathname = "/launch";
  return NextResponse.redirect(url);
}

// Match all routes except Next.js internal static assets
export const config = {
  matcher: [
    "/((?!_next/static|_next/image|favicon.ico|icon.svg|logo.svg|robots.txt|sitemap.xml|llms.txt|llms-full.txt|site.webmanifest).*)",
  ],
};
