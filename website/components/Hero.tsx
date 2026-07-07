"use client";

import TerminalWindow from "./TerminalWindow";
import CommandBlock from "./ui/CommandBlock";
import SplitText from "./ui/SplitText";
import FadeIn from "./ui/FadeIn";

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

      {/* Top subtle glow */}
      <div
        className="absolute top-0 left-1/2 -translate-x-1/2 w-[800px] h-[300px] pointer-events-none"
        style={{
          background: "radial-gradient(ellipse at center top, rgba(255,255,255,0.02) 0%, transparent 70%)",
        }}
      />

      <div className="container-page w-full py-20 lg:py-28">
        <div className="grid lg:grid-cols-2 gap-12 lg:gap-16 items-center">
          {/* Left */}
          <div>
            <FadeIn delay={0.1}>
              <div className="inline-flex items-center gap-2 mb-8">
                <span className="w-1.5 h-1.5 rounded-full bg-[#444]" />
                <span className="font-mono text-xs text-[#555] tracking-widest uppercase">
                  v1.0.1 — open source
                </span>
              </div>
            </FadeIn>

            <h1
              className="font-mono font-semibold text-white mb-6"
              style={{ fontSize: "clamp(2.2rem, 5vw, 4rem)", lineHeight: 1.05, letterSpacing: "-0.03em" }}
            >
              <SplitText text="It scans." by="chars" delay={0.2} stagger={0.05} direction="down" />
              <br />
              <SplitText text="It infers." by="chars" delay={0.65} stagger={0.05} direction="down" />
              <br />
              <SplitText text="It plays." by="chars" delay={1.1} stagger={0.05} direction="down" />
            </h1>

            <FadeIn delay={1.5}>
              <p className="text-[#666] max-w-md mb-10 leading-relaxed" style={{ fontSize: "1.1rem" }}>
                Moodwave reads your repository, detects your working mood from code signals,
                and streams perfectly matched music right in your terminal.
              </p>
            </FadeIn>

            <FadeIn delay={1.6} y={15}>
              <div className="space-y-3 max-w-lg">
                <CommandBlock
                  label="macOS / Linux"
                  command="curl -fsSL https://raw.githubusercontent.com/Boredooms/Moodwave-CLI/main/cli/scripts/install.sh | bash"
                />
                <CommandBlock
                  label="Windows (PowerShell)"
                  prompt="PS>"
                  command="irm https://raw.githubusercontent.com/Boredooms/Moodwave-CLI/main/cli/scripts/install.ps1 | iex"
                />
              </div>
            </FadeIn>

            <FadeIn delay={1.7}>
              <div className="mt-8 flex items-center gap-6">
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
              </div>
            </FadeIn>
          </div>

          {/* Right: terminal */}
          <FadeIn delay={0.4} x={30} y={0} duration={0.8}>
            <TerminalWindow />
          </FadeIn>
        </div>
      </div>
    </section>
  );
}
