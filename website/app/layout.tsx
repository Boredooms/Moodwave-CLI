import type { Metadata } from "next";
import localFont from "next/font/local";
import { JetBrains_Mono, Inter } from "next/font/google";
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

export const metadata: Metadata = {
  title: "Moodwave — Terminal Mood Music Companion",
  description:
    "Moodwave scans your codebase, infers your working mood, and streams perfectly matched music right in your terminal. CLI-first. Lightweight. Open source.",
  keywords: ["cli", "terminal", "music player", "developer tools", "mood detection", "youtube cli"],
  openGraph: {
    title: "Moodwave — Terminal Mood Music Companion",
    description: "A CLI that scans your codebase and plays music that matches your mood.",
    type: "website",
  },
  twitter: {
    card: "summary_large_image",
    title: "Moodwave — Terminal Mood Music Companion",
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className={`${inter.variable} ${jetbrainsMono.variable}`}>
      <body className="bg-bg text-primary antialiased">{children}</body>
    </html>
  );
}
