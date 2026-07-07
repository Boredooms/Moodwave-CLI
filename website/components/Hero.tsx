"use client";

import { motion } from "framer-motion";
import TerminalWindow from "./TerminalWindow";
import CommandBlock from "./ui/CommandBlock";

const words = ["scans", "your", "codebase.", "plays", "music", "that", "fits."];

export default function Hero() {
  return (
    <section className="relative min-h-screen flex items-center pt-14 overflow-hidden">
      {/* Ambient grid */}
      <div
        className="absolute inset-0 pointer-events-none opacity-[0.03]"
        style={{
          backgroundImage:
            "linear-gradient(rgba(255,255,255,0.1) 1px, transparent 1px), linear-gradient(90deg, rgba(255,255,255,0.1) 1px, transparent 1px)",
          backgroundSize: "72px 72px",
        }}
      />

      {/* Subtle top radial glow */}
      <div
        className="absolute top-0 left-1/2 -translate-x-1/2 w-[800px] h-[400px] pointer-events-none"
        style={{
          background: "radial-gradient(ellipse at center top, rgba(255,255,255,0.03) 0%, transparent 70%)",
        }}
      />

      <div className="container-page w-full py-20 lg:py-28">
        <div className="grid lg:grid-cols-2 gap-12 lg:gap-16 items-center">
          {/* Left: copy */}
          <div>
            {/* Version badge */}
            <motion.div
              initial={{ opacity: 0, y: 8 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.5, delay: 0.1 }}
              className="inline-flex items-center gap-2 mb-8"
            >
              <span className="inline-block w-1.5 h-1.5 rounded-full bg-[#555]" />
              <span className="font-mono text-xs text-[#555] tracking-widest uppercase">
                v1.0.1 — open source
              </span>
            </motion.div>

            {/* Headline */}
            <h1 className="font-mono text-display-xl font-semibold text-white mb-6 leading-none">
              {["It", "scans.", "It", "infers.", "It", "plays."].map((word, i) => (
                <motion.span
                  key={i}
                  initial={{ opacity: 0, y: 16 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ duration: 0.55, delay: 0.25 + i * 0.08, ease: [0.22, 1, 0.36, 1] }}
                  className={`inline-block mr-[0.25em] ${
                    word.endsWith(".") ? "text-white" : "text-white"
                  }`}
                >
                  {word}
                </motion.span>
              ))}
            </h1>

            {/* Subtitle */}
            <motion.p
              initial={{ opacity: 0, y: 12 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.6, delay: 0.75, ease: [0.22, 1, 0.36, 1] }}
              className="text-body-lg text-[#666] max-w-md mb-10 leading-relaxed"
            >
              Moodwave scans your repository, infers your working mood from code signals,
              and streams perfectly matched music directly in your terminal.
            </motion.p>

            {/* Install command */}
            <motion.div
              initial={{ opacity: 0, y: 12 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.6, delay: 0.9 }}
              className="space-y-3 max-w-lg"
            >
              <CommandBlock
                label="macOS / Linux"
                command="curl -fsSL https://raw.githubusercontent.com/Boredooms/Moodwave-CLI/main/cli/scripts/install.sh | bash"
              />
              <CommandBlock
                label="Windows (PowerShell)"
                prompt="PS>"
                command="irm https://raw.githubusercontent.com/Boredooms/Moodwave-CLI/main/cli/scripts/install.ps1 | iex"
              />
            </motion.div>

            {/* Meta links */}
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              transition={{ duration: 0.6, delay: 1.1 }}
              className="mt-8 flex items-center gap-6"
            >
              <a
                href="https://github.com/Boredooms/Moodwave-CLI"
                target="_blank"
                rel="noopener noreferrer"
                className="font-mono text-sm text-[#555] hover:text-white transition-colors duration-200 flex items-center gap-2"
              >
                <span>View on GitHub</span>
                <span className="text-[#333]">↗</span>
              </a>
              <span className="text-[#2a2a2a]">·</span>
              <a
                href="#how-it-works"
                className="font-mono text-sm text-[#555] hover:text-white transition-colors duration-200"
              >
                How it works
              </a>
            </motion.div>
          </div>

          {/* Right: terminal */}
          <motion.div
            initial={{ opacity: 0, x: 20 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ duration: 0.8, delay: 0.4, ease: [0.22, 1, 0.36, 1] }}
            className="relative"
          >
            {/* Glow behind terminal */}
            <div
              className="absolute -inset-8 pointer-events-none"
              style={{
                background: "radial-gradient(ellipse at center, rgba(255,255,255,0.02) 0%, transparent 70%)",
              }}
            />
            <TerminalWindow />
            {/* Reflection hint */}
            <div className="mt-3 h-8 rounded-b-xl opacity-20 blur-sm scale-y-[-1] origin-top" style={{
              background: "linear-gradient(to bottom, rgba(255,255,255,0.03), transparent)"
            }} />
          </motion.div>
        </div>
      </div>
    </section>
  );
}
