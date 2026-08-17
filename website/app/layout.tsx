import type { Metadata } from "next";
import { JetBrains_Mono, Inter } from "next/font/google";
import { GoogleAnalytics } from "@next/third-parties/google";
import LenisProvider from "../components/LenisProvider";
import { GlobalJsonLd } from "../components/StructuredData";
import { GA_ID } from "../lib/analytics";
import "./globals.css";

const inter = Inter({
  subsets: ["latin"],
  variable: "--font-inter",
  display: "swap",
});

const jetbrainsMono = JetBrains_Mono({
  subsets: ["latin"],
  variable: "--font-geist-mono",
  display: "swap",
});

const baseUrl = "https://www.moodwave-cli.xyz";

export const metadata: Metadata = {
  metadataBase: new URL(baseUrl),
  title: {
    default: "Moodwave — Terminal Mood Music Companion",
    template: "%s | Moodwave CLI",
  },
  description:
    "Moodwave scans your codebase, infers your working mood across 10 developer states, and streams perfectly matched music right in your terminal. CLI-first, lightweight, open source.",
  keywords: [
    "moodwave",
    "moodwave-cli",
    "cli music player",
    "terminal music companion",
    "developer tools",
    "codebase mood scan",
    "youtube cli",
    "bubble tea tui",
    "terminal visualizer",
    "golang cli",
    "radio browser cli",
    "developer productivity",
  ],
  authors: [{ name: "Boredooms", url: "https://github.com/Boredooms" }],
  creator: "Boredooms",
  publisher: "Moodwave",
  applicationName: "Moodwave",
  alternates: {
    canonical: baseUrl,
  },
  manifest: "/site.webmanifest",
  icons: {
    icon: [
      { url: "/icon.svg", type: "image/svg+xml" },
      { url: "/logo.svg", type: "image/svg+xml" },
    ],
    apple: "/logo.svg",
    shortcut: "/icon.svg",
  },
  robots: {
    index: true,
    follow: true,
    googleBot: {
      index: true,
      follow: true,
      "max-video-preview": -1,
      "max-image-preview": "large",
      "max-snippet": -1,
    },
  },
  openGraph: {
    type: "website",
    locale: "en_US",
    url: baseUrl,
    siteName: "Moodwave",
    title: "Moodwave — Terminal Mood Music Companion",
    description:
      "A CLI that scans your codebase, infers your cognitive mood, and streams music that matches your state directly into your shell.",
    images: [
      {
        url: "/opengraph-image",
        width: 1200,
        height: 630,
        alt: "Moodwave — Terminal Mood Music Companion",
      },
    ],
  },
  twitter: {
    card: "summary_large_image",
    title: "Moodwave — Terminal Mood Music Companion",
    description:
      "A CLI that scans your codebase, infers your cognitive mood, and streams music that matches your state directly into your shell.",
    images: ["/opengraph-image"],
    creator: "@moodwave_cli",
  },
  category: "technology",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className={`${inter.variable} ${jetbrainsMono.variable}`}>
      <head>
        <GlobalJsonLd />
      </head>
      <body style={{ background: "#080808", color: "#fff", overflowX: "hidden" }}>
        <LenisProvider>{children}</LenisProvider>
      </body>
      {/* Google Analytics 4 — loads after page interactive, zero LCP impact */}
      {GA_ID && <GoogleAnalytics gaId={GA_ID} />}
    </html>
  );
}

