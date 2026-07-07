"use client";

import { motion } from "framer-motion";
import TerminalWindow from "./TerminalWindow";
import CommandBlock from "./ui/CommandBlock";

import type { Variants } from "framer-motion";

const FADE_UP: Variants = {
  hidden: { opacity: 0, y: 20 },
  show: { opacity: 1, y: 0, transition: { duration: 0.6, ease: "easeOut" } },
};

const STAGGER: Variants = {
  hidden: {},
  show: { transition: { staggerChildren: 0.08 } },
};

export default function Hero() {
  return (
    <section className="relative min-h-screen flex items-center pt-14 overflow-hidden">
      {/* Subtle grid */}
      <div
        className="absolute inset-0 pointer-events-none opacity-[0.025]"
        style={{
          backgroundImage:
            "linear-gradient(rgba(255,255,255,0.15) 1px, transparent 1px), linear-gradient(90deg, rgba(255,255,255,0.15) 1px, transparent 1px)",
          backgroundSize: "64px 64px",
        }}
      />

      <div className="container-page w-full py-20 lg:py-28">
        <div className="grid lg:grid-cols-2 gap-12 lg:gap-16 items-center">
          {/* Left */}
          <motion.div
            initial="hidden"
            animate="show"
            variants={STAGGER}
          >
            <motion.div variants={FADE_UP} className="inline-flex items-center gap-2 mb-8">
              <span className="w-1.5 h-1.5 rounded-full bg-[#444]" />
              <span className="font-mono text-xs text-[#555] tracking-widest uppercase">
                v1.0.1 — open source
              </span>
            </motion.div>

            <motion.h1
              variants={FADE_UP}
              className="font-mono font-semibold text-white mb-6"
              style={{ fontSize: "clamp(2.2rem, 5vw, 4rem)", lineHeight: 1.05, letterSpacing: "-0.03em" }}
            >
              It scans.<br />It infers.<br />It plays.
            </motion.h1>

            <motion.p variants={FADE_UP} className="text-[#666] max-w-md mb-10 leading-relaxed" style={{ fontSize: "1.1rem" }}>
              Moodwave reads your repository, detects your working mood from code signals,
              and streams perfectly matched music right in your terminal.
            </motion.p>

            <motion.div variants={FADE_UP} className="space-y-3 max-w-lg">
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

            <motion.div variants={FADE_UP} className="mt-8 flex items-center gap-6">
              <a
                href="https://github.com/Boredooms/Moodwave-CLI"
                target="_blank"
                rel="noopener noreferrer"
                className="font-mono text-sm text-[#555] hover:text-white transition-colors duration-200"
              >
                GitHub ↗
              </a>
              <span className="text-[#2a2a2a]">·</span>
              <a href="#how-it-works" className="font-mono text-sm text-[#555] hover:text-white transition-colors duration-200">
                How it works
              </a>
            </motion.div>
          </motion.div>

          {/* Right: terminal */}
          <motion.div
            initial={{ opacity: 0, x: 30 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ duration: 0.8, delay: 0.3, ease: [0.22, 1, 0.36, 1] }}
          >
            <TerminalWindow />
          </motion.div>
        </div>
      </div>
    </section>
  );
}
