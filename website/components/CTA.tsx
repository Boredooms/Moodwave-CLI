"use client";

import { motion } from "framer-motion";
import CommandBlock from "./ui/CommandBlock";

export default function CTA() {
  return (
    <>
      <section className="divider section-pad" id="cta">
        <div className="container-page">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true, margin: "-60px" }}
            transition={{ duration: 0.65, ease: [0.22, 1, 0.36, 1] }}
            className="border border-white/[0.07] rounded-xl p-12 lg:p-20 text-center relative overflow-hidden"
          >
            <div
              className="absolute inset-0 pointer-events-none"
              style={{ background: "radial-gradient(ellipse at center, rgba(255,255,255,0.02) 0%, transparent 65%)" }}
            />
            <p className="font-mono text-xs text-[#444] uppercase tracking-[0.2em] mb-6 relative">Get started</p>
            <h2 className="font-mono font-semibold text-white mb-4 relative leading-tight" style={{ fontSize: "clamp(1.8rem, 3.5vw, 2.8rem)", letterSpacing: "-0.03em" }}>
              Ready to listen to<br />your codebase?
            </h2>
            <p className="text-sm text-[#666] max-w-md mx-auto mb-10 relative leading-relaxed">
              Install Moodwave in one command. Start scanning. Let it play.
            </p>

            <div className="max-w-lg mx-auto space-y-3 relative">
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

            <div className="mt-10 flex items-center justify-center gap-8 relative">
              <a href="https://github.com/Boredooms/Moodwave-CLI" target="_blank" rel="noopener noreferrer" className="font-mono text-sm text-[#555] hover:text-white transition-colors duration-200">
                View source ↗
              </a>
              <span className="text-[#2a2a2a]">·</span>
              <a href="https://github.com/Boredooms/Moodwave-CLI/releases" target="_blank" rel="noopener noreferrer" className="font-mono text-sm text-[#555] hover:text-white transition-colors duration-200">
                Download binary ↗
              </a>
            </div>
          </motion.div>
        </div>
      </section>

      <footer className="divider">
        <div className="container-page py-10 flex flex-col sm:flex-row items-center justify-between gap-4">
          <div className="flex items-center gap-4">
            <span className="font-mono text-sm text-white">moodwave</span>
            <span className="font-mono text-xs text-[#333]">v1.0.1</span>
          </div>
          <div className="flex items-center gap-6">
            {[
              { href: "https://github.com/Boredooms/Moodwave-CLI", label: "GitHub" },
              { href: "https://github.com/Boredooms/Moodwave-CLI/releases", label: "Releases" },
              { href: "https://github.com/Boredooms/Moodwave-CLI/blob/main/docs/commands.md", label: "Docs" },
            ].map((l) => (
              <a key={l.href} href={l.href} target="_blank" rel="noopener noreferrer" className="font-mono text-xs text-[#444] hover:text-[#888] transition-colors">
                {l.label}
              </a>
            ))}
          </div>
          <span className="font-mono text-xs text-[#333]">MIT License</span>
        </div>
      </footer>
    </>
  );
}
