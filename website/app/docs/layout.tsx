import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Docs — Setup, Commands & Audio Backends",
  description:
    "Complete setup guides for Windows, macOS, and Linux, command reference, audio backend configuration, and environment variables for Moodwave CLI.",
  alternates: {
    canonical: "https://www.moodwave-cli.xyz/docs",
  },
  openGraph: {
    title: "Moodwave Docs — Terminal Mood Music Companion",
    description:
      "Platform setup guides, full command reference, and everything needed to go from install to a working terminal music companion.",
    url: "https://www.moodwave-cli.xyz/docs",
  },
};

export default function DocsLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return <>{children}</>;
}
