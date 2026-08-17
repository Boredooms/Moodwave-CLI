import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Changelog — Version History, Timeline & Metrics",
  description:
    "Follow the evolution of Moodwave CLI across releases: Bubble Tea rebuild, live terminal visualizers, audio state machine, and scan performance metrics.",
  alternates: {
    canonical: "https://www.moodwave-cli.xyz/changelog",
  },
  openGraph: {
    title: "Moodwave Changelog — Version History & Release Notes",
    description:
      "Full release history, interactive execution latency charts, binary footprint tracking, and detailed feature logs for Moodwave CLI.",
    url: "https://www.moodwave-cli.xyz/changelog",
  },
};

export default function ChangelogLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return <>{children}</>;
}
