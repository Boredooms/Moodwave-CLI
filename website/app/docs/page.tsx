"use client";

import { useState } from "react";
import Nav from "../../components/Nav";
import SplitText from "../../components/ui/SplitText";
import FadeIn from "../../components/ui/FadeIn";
import CommandBlock from "../../components/ui/CommandBlock";
import {
  Terminal,
  Monitor,
  Apple,
  BookOpen,
  Rocket,
  Settings,
  Palette,
  Activity,
  Stethoscope,
} from "lucide-react";
import { BreadcrumbsJsonLd } from "../../components/StructuredData";

// ─────────────────────────────────────────────────────────────────────────────
// Platform-specific setup docs — Windows / macOS / Linux
// ─────────────────────────────────────────────────────────────────────────────

const platforms = [
  {
    id: "windows",
    label: "Windows",
    icon: Monitor,
    install: [
      { label: "Install (PowerShell)", prompt: "PS>", command: "irm https://raw.githubusercontent.com/Boredooms/Moodwave-CLI/main/cli/scripts/install.ps1 | iex" },
    ],
    audio: [
      { label: "mpv (recommended)", prompt: "PS>", command: "winget install mpv" },
      { label: "ffplay (ffmpeg)", prompt: "PS>", command: "winget install ffmpeg" },
    ],
    paths: [
      { label: "Config file", value: "%APPDATA%\\moodwave\\config.json" },
      { label: "Cache directory", value: "%LOCALAPPDATA%\\moodwave\\cache\\" },
    ],
    notes: "Moodwave also ships a native PresentationCore/Media Foundation audio backend, so it can play streams out of the box even without mpv or ffmpeg installed — run moodwave doctor to confirm which backend was picked.",
  },
  {
    id: "macos",
    label: "macOS",
    icon: Apple,
    install: [
      { label: "Install (shell)", prompt: "$", command: "curl -fsSL https://raw.githubusercontent.com/Boredooms/Moodwave-CLI/main/cli/scripts/install.sh | bash" },
    ],
    audio: [
      { label: "mpv (recommended)", prompt: "$", command: "brew install mpv" },
      { label: "ffplay (ffmpeg)", prompt: "$", command: "brew install ffmpeg" },
    ],
    paths: [
      { label: "Config file", value: "~/.config/moodwave/config.json" },
      { label: "Cache directory", value: "~/Library/Caches/moodwave/" },
    ],
    notes: "Playback falls back to afplay through CoreAudio when no other backend is found, with full support for spatial audio on Apple Silicon.",
  },
  {
    id: "linux",
    label: "Linux",
    icon: Terminal,
    install: [
      { label: "Install (shell)", prompt: "$", command: "curl -fsSL https://raw.githubusercontent.com/Boredooms/Moodwave-CLI/main/cli/scripts/install.sh | bash" },
    ],
    audio: [
      { label: "mpv (recommended)", prompt: "$", command: "sudo apt install mpv" },
      { label: "ffplay (ffmpeg)", prompt: "$", command: "sudo apt install ffmpeg" },
    ],
    paths: [
      { label: "Config file", value: "~/.config/moodwave/config.json" },
      { label: "Cache directory", value: "~/.cache/moodwave/" },
    ],
    notes: "Works on amd64, arm64, and armv7 (Raspberry Pi and similar). Moodwave auto-downloads a static ffplay build into its cache directory if none is found on PATH.",
  },
];

// ─────────────────────────────────────────────────────────────────────────────
// Command reference — grouped for the "Installation & User Docs" section
// ─────────────────────────────────────────────────────────────────────────────

const commandGroups = [
  {
    title: "Getting started",
    icon: Rocket,
    commands: [
      { cmd: "moodwave init", desc: "Create config file and cache directories. Run once after install." },
      { cmd: "moodwave doctor", desc: "Diagnose terminal, config, audio backends, and music source connectivity." },
      { cmd: "moodwave", desc: "Launch the interactive TUI — scan, search, play, and manage playlists." },
      { cmd: "moodwave update", desc: "Check for and install the latest release in place. Also runs automatically in the TUI." },
    ],
  },
  {
    title: "Mood & playback",
    icon: Activity,
    commands: [
      { cmd: "moodwave scan [path]", desc: "Analyze the repository and infer a coding mood from language mix, TODOs, git churn, and structure." },
      { cmd: "moodwave mood", desc: "Show the last detected mood profile, confidence, and music traits without re-scanning." },
      { cmd: "moodwave play", desc: "Find and stream music matching the current mood. Runs a scan first if none exists." },
      { cmd: "moodwave search <query>", desc: "Search YouTube directly and pick a track to play." },
      { cmd: "moodwave next", desc: "Advance the recommendation queue to the next candidate." },
      { cmd: "moodwave queue", desc: "Show the current recommendation queue and playing position." },
    ],
  },
  {
    title: "Customization",
    icon: Palette,
    commands: [
      { cmd: "moodwave theme <name>", desc: "Switch the visual color theme (7+ presets — Midnight, Dracula, Nord, Tokyo Night, and more)." },
      { cmd: "moodwave visual <mode>", desc: "Switch the terminal visualizer — spectrum, fire, matrix, plasma, aurora, and 15+ others." },
      { cmd: "moodwave config", desc: "View the resolved configuration as JSON, or edit it directly." },
      { cmd: "moodwave source", desc: "List configured music sources and their live health status." },
    ],
  },
  {
    title: "Diagnostics",
    icon: Stethoscope,
    commands: [
      { cmd: "moodwave status", desc: "Show version, config paths, last mood, and playback state." },
      { cmd: "moodwave doctor", desc: "Full diagnostic pass with a one-key auto-fix for missing audio backends." },
    ],
  },
];

const envVars = [
  { name: "MOODWAVE_THEME", desc: "Visual theme ID" },
  { name: "MOODWAVE_VISUAL", desc: "Visual mode" },
  { name: "MOODWAVE_NO_COLOR=1", desc: "Disable ANSI color" },
  { name: "MOODWAVE_NO_ANIMATION=1", desc: "Disable animation" },
  { name: "MOODWAVE_BACKEND", desc: "Force a specific audio backend" },
  { name: "MOODWAVE_VOLUME", desc: "Volume, 0–100" },
  { name: "MOODWAVE_DEBUG=1", desc: "Enable debug logging" },
  { name: "JAMENDO_CLIENT_ID", desc: "Enable the Jamendo source" },
];

export default function Docs() {
  const [activePlatform, setActivePlatform] = useState("windows");
  const current = platforms.find((p) => p.id === activePlatform)!;

  return (
    <div style={{ background: "#080808", minHeight: "100vh", color: "#ffffff", paddingBottom: "100px", position: "relative" }}>
      <BreadcrumbsJsonLd
        items={[
          { name: "Home", url: "https://www.moodwave-cli.xyz" },
          { name: "Docs", url: "https://www.moodwave-cli.xyz/docs" },
        ]}
      />
      <div
        className="absolute inset-0 pointer-events-none opacity-[0.025]"
        style={{
          backgroundImage:
            "linear-gradient(rgba(255,255,255,0.15) 1px, transparent 1px), linear-gradient(90deg, rgba(255,255,255,0.15) 1px, transparent 1px)",
          backgroundSize: "64px 64px",
        }}
      />
      <div
        className="absolute top-0 left-1/2 -translate-x-1/2 w-[800px] h-[300px] pointer-events-none"
        style={{ background: "radial-gradient(ellipse at center top, rgba(255,255,255,0.02) 0%, transparent 70%)" }}
      />

      <Nav version="v2.0.0" />

      {/* Hero */}
      <section className="relative pt-28 md:pt-32 pb-12 md:pb-16 overflow-hidden border-b border-white/[0.05]">
        <div className="container-page text-center relative z-10">
          <p className="font-mono text-xs text-[#555] uppercase tracking-[0.2em] mb-4">
            <SplitText text="Reference & Guides" by="chars" delay={0.15} stagger={0.03} direction="down" />
          </p>
          <h1 className="font-mono font-semibold text-white tracking-tight leading-tight mb-5" style={{ fontSize: "clamp(2rem, 5vw, 3.5rem)" }}>
            <SplitText text="Docs" by="chars" delay={0.35} stagger={0.05} direction="down" />
          </h1>
          <FadeIn delay={0.6} y={15}>
            <p className="text-[#666] max-w-xl mx-auto text-sm md:text-base leading-relaxed">
              Platform setup guides, the full command reference, and everything needed to go from install to a working terminal music companion.
            </p>
          </FadeIn>
        </div>
      </section>

      {/* ── Platform docs ── */}
      <section className="py-16 md:py-20 border-b border-white/[0.05]">
        <div className="container-page">
          <FadeIn>
            <p className="font-mono text-xs text-[#444] uppercase tracking-[0.2em] mb-2">Platform setup</p>
            <h2 className="font-mono font-semibold text-white text-xl md:text-2xl mb-8">Windows, macOS &amp; Linux</h2>
          </FadeIn>

          {/* Platform tabs — wraps on mobile instead of overflowing */}
          <FadeIn delay={0.1}>
            <div className="flex flex-wrap gap-2 mb-8">
              {platforms.map((p) => {
                const Icon = p.icon;
                return (
                  <button
                    key={p.id}
                    onClick={() => setActivePlatform(p.id)}
                    className={`flex items-center gap-2 font-mono text-xs px-4 py-2.5 rounded-lg border transition-colors duration-200 cursor-pointer ${
                      activePlatform === p.id
                        ? "bg-white/[0.08] border-white/20 text-white"
                        : "bg-white/[0.01] border-white/[0.07] text-[#666] hover:text-[#999]"
                    }`}
                  >
                    <Icon className="w-3.5 h-3.5" />
                    {p.label}
                  </button>
                );
              })}
            </div>
          </FadeIn>

          <div className="grid lg:grid-cols-2 gap-8 lg:gap-10 w-full min-w-0">
            <FadeIn delay={0.15} className="w-full min-w-0">
              <div className="space-y-6 w-full min-w-0">
                <div>
                  <p className="font-mono text-xs text-[#444] uppercase tracking-widest mb-3">Install</p>
                  <div className="space-y-3 w-full min-w-0">
                    {current.install.map((c) => (
                      <CommandBlock key={c.command} label={c.label} prompt={c.prompt} command={c.command} />
                    ))}
                  </div>
                </div>

                <div>
                  <p className="font-mono text-xs text-[#444] uppercase tracking-widest mb-3">Audio backend (optional)</p>
                  <div className="space-y-3 w-full min-w-0">
                    {current.audio.map((c) => (
                      <CommandBlock key={c.command} label={c.label} prompt={c.prompt} command={c.command} />
                    ))}
                  </div>
                </div>
              </div>
            </FadeIn>

            <FadeIn delay={0.2} className="w-full min-w-0">
              <div className="border border-white/[0.07] rounded-xl bg-white/[0.01] p-5 md:p-6 space-y-5 h-full">
                <div>
                  <p className="font-mono text-xs text-[#444] uppercase tracking-widest mb-3">File locations</p>
                  <div className="space-y-2.5">
                    {current.paths.map((p) => (
                      <div key={p.label} className="flex flex-col sm:flex-row sm:items-baseline sm:gap-2 min-w-0">
                        <span className="text-xs text-[#666] flex-shrink-0">{p.label}:</span>
                        <code className="font-mono text-xs text-[#aaa] break-all">{p.value}</code>
                      </div>
                    ))}
                  </div>
                </div>
                <div className="border-t border-white/[0.06] pt-4">
                  <p className="text-xs text-[#666] leading-relaxed">{current.notes}</p>
                </div>
                <div className="border-t border-white/[0.06] pt-4">
                  <p className="text-xs text-[#555] leading-relaxed">
                    Verify everything is working after install:
                  </p>
                  <div className="mt-2">
                    <CommandBlock prompt={activePlatform === "windows" ? "PS>" : "$"} command="moodwave doctor" />
                  </div>
                </div>
              </div>
            </FadeIn>
          </div>
        </div>
      </section>

      {/* ── Installation & user docs ── */}
      <section className="py-16 md:py-20">
        <div className="container-page">
          <FadeIn>
            <p className="font-mono text-xs text-[#444] uppercase tracking-[0.2em] mb-2">Installation &amp; user guide</p>
            <h2 className="font-mono font-semibold text-white text-xl md:text-2xl mb-3">Command reference</h2>
            <p className="text-sm text-[#666] max-w-2xl mb-10 leading-relaxed">
              Every moodwave subcommand, grouped by what you&apos;re trying to do. Run <code className="text-[#999]">moodwave --help</code> anytime for the same list in your terminal.
            </p>
          </FadeIn>

          <div className="grid md:grid-cols-2 gap-6 md:gap-8">
            {commandGroups.map((group, gi) => {
              const Icon = group.icon;
              return (
                <FadeIn key={group.title} delay={0.1 + gi * 0.06}>
                  <div className="border border-white/[0.07] rounded-xl bg-white/[0.01] p-5 md:p-6 h-full">
                    <div className="flex items-center gap-2 mb-4 text-white">
                      <Icon className="w-4 h-4 text-[#888]" />
                      <span className="font-mono text-xs font-semibold uppercase tracking-wider">{group.title}</span>
                    </div>
                    <ul className="space-y-3.5">
                      {group.commands.map((c) => (
                        <li key={c.cmd} className="min-w-0">
                          <code className="font-mono text-xs text-white bg-white/[0.05] px-2 py-1 rounded inline-block break-all">{c.cmd}</code>
                          <p className="text-xs text-[#666] mt-1.5 leading-relaxed">{c.desc}</p>
                        </li>
                      ))}
                    </ul>
                  </div>
                </FadeIn>
              );
            })}
          </div>

          {/* Env vars */}
          <FadeIn delay={0.3}>
            <div className="mt-10 border border-white/[0.07] rounded-xl bg-white/[0.01] p-5 md:p-6">
              <div className="flex items-center gap-2 mb-4 text-white">
                <Settings className="w-4 h-4 text-[#888]" />
                <span className="font-mono text-xs font-semibold uppercase tracking-wider">Environment variables</span>
              </div>
              <div className="overflow-x-auto -mx-1">
                <table className="w-full min-w-[420px] text-left">
                  <tbody>
                    {envVars.map((v) => (
                      <tr key={v.name} className="border-b border-white/[0.05] last:border-b-0">
                        <td className="py-2.5 px-1 font-mono text-xs text-[#aaa] whitespace-nowrap">{v.name}</td>
                        <td className="py-2.5 px-1 text-xs text-[#666]">{v.desc}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          </FadeIn>

          <FadeIn delay={0.35}>
            <div className="mt-8 flex items-center gap-2 text-xs font-mono text-[#555]">
              <BookOpen className="w-3.5 h-3.5" />
              <a
                href="https://github.com/Boredooms/Moodwave-CLI/tree/main/docs"
                target="_blank"
                rel="noopener noreferrer"
                className="hover:text-white transition-colors"
              >
                Full docs on GitHub ↗
              </a>
            </div>
          </FadeIn>
        </div>
      </section>
    </div>
  );
}
