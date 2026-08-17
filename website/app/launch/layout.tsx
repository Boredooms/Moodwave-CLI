import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Launch — Countdown & Official Release Waitlist",
  description:
    "Join the official Moodwave CLI launch waitlist. Live countdown to September 3, 2026, interactive particle animation, and early release access.",
  alternates: {
    canonical: "https://www.moodwave-cli.xyz/launch",
  },
  openGraph: {
    title: "Moodwave Launch — Your Terminal's New Soundtrack",
    description:
      "Live launch countdown to September 3, 2026. Join the developer waitlist for Moodwave CLI.",
    url: "https://www.moodwave-cli.xyz/launch",
  },
};

export default function LaunchLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return <>{children}</>;
}
