/**
 * analytics.ts — Moodwave GA4 Event Helpers
 *
 * Usage (client components only):
 *   import { trackEvent } from "@/lib/analytics";
 *   trackEvent("install_click", { method: "homebrew" });
 *
 * All events are no-ops when GA_ID is not configured.
 */

export const GA_ID = process.env.NEXT_PUBLIC_GA_ID ?? "";

// ─── Type-safe event catalogue ───────────────────────────────────────────────

type InstallClickEvent = {
  event: "install_click";
  params: { method: "homebrew" | "go_install" | "snap" | "scoop" | "curl" };
};

type DocsClickEvent = {
  event: "docs_click";
  params: { section: string };
};

type ChangelogClickEvent = {
  event: "changelog_click";
  params: { version?: string };
};

type GithubClickEvent = {
  event: "github_click";
  params: { location: "header" | "footer" | "hero" | "cta" };
};

type PageViewEvent = {
  event: "page_view";
  params: { page_path: string; page_title?: string };
};

type GAEvent =
  | InstallClickEvent
  | DocsClickEvent
  | ChangelogClickEvent
  | GithubClickEvent
  | PageViewEvent;

// ─── Send helper ─────────────────────────────────────────────────────────────

declare global {
  // eslint-disable-next-line no-var
  var gtag: (
    command: "event" | "config" | "js",
    targetId: string | Date,
    params?: Record<string, unknown>
  ) => void;
}

/**
 * Track a GA4 custom event.
 * Safe to call on the server — silently skips if window/gtag is unavailable.
 */
export function trackEvent<E extends GAEvent>(
  event: E["event"],
  params?: E extends { params: infer P } ? P : never
): void {
  if (typeof window === "undefined") return;
  if (!GA_ID) return;
  if (typeof window.gtag !== "function") return;
  window.gtag("event", event, params ?? {});
}

/**
 * Manually fire a page_view. Useful for soft navigations if needed.
 */
export function trackPageView(path: string, title?: string): void {
  trackEvent("page_view", { page_path: path, page_title: title });
}
