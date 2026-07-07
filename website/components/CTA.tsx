"use client";

import { motion, useInView } from "framer-motion";
import { useRef } from "react";
import CommandBlock from "./ui/CommandBlock";

export default function CTA() {
  const ref = useRef<HTMLDivElement>(null);
  const inView = useInView(ref, { once: true, margin: "-80px 0px" });

  return (
    <>
      {/* CTA */}
      <section className="divider section-pad" id="cta">
        <div className="container-page">
          <motion.div
            ref={ref}
            initial={{ opacity: 0, y: 24 }}
            animate={inView ? { opacity: 1, y: 0 } : {}}
            transition={{ duration: 0.7, ease: [0.22, 1, 0.36, 1] }}
            className="border border-white/[0.07] rounded-xl p-12 lg:p-20 text-center relative overflow-hidden"
          >
            {/* ambient background */}
            <div
              className="absolute inset-0 pointer-events-none"
              style={{
                background: "radial-gradient(ellipse at center, rgba(255,255,255,0.02) 0%, transparent 65%)",
              }}
            />

            <p className="font-mono text-xs text-[#444] uppercase tracking-[0.2em] mb-6 relative">
              Get started
            </p>
            <h2 className="font-mono text-display-lg text-white font-semibold mb-4 relative leading-tight">
              Ready to listen to<br />your codebase?
            </h2>
            <p className="text-body-md text-[#666] max-w-md mx-auto mb-10 relative">
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
              <a
                href="https://github.com/Boredooms/Moodwave-CLI"
                target="_blank"
                rel="noopener noreferrer"
                className="font-mono text-sm text-[#555] hover:text-white transition-colors duration-200 flex items-center gap-2"
              >
                View source ↗
              </a>
              <span className="text-[#2a2a2a]">·</span>
              <a
                href="https://github.com/Boredooms/Moodwave-CLI/releases"
                target="_blank"
                rel="noopener noreferrer"
                className="font-mono text-sm text-[#555] hover:text-white transition-colors duration-200"
              >
                Download binary ↗
              </a>
            </div>
          </motion.div>
        </div>
      </section>

      {/* Footer */}
      <footer className="divider">
        <div className="container-page py-10 flex flex-col sm:flex-row items-center justify-between gap-4">
          <div className="flex items-center gap-6">
            <span className="font-mono text-sm text-white">moodwave</span>
            <span className="font-mono text-xs text-[#333]">v1.0.1</span>
          </div>
          <div className="flex items-center gap-6">
            <a
              href="https://github.com/Boredooms/Moodwave-CLI"
              target="_blank"
              rel="noopener noreferrer"
              className="font-mono text-xs text-[#444] hover:text-[#888] transition-colors"
            >
              GitHub
            </a>
            <a
              href="https://github.com/Boredooms/Moodwave-CLI/releases"
              target="_blank"
              rel="noopener noreferrer"
              className="font-mono text-xs text-[#444] hover:text-[#888] transition-colors"
            >
              Releases
            </a>
            <a
              href="https://github.com/Boredooms/Moodwave-CLI/blob/main/cli/docs/commands.md"
              target="_blank"
              rel="noopener noreferrer"
              className="font-mono text-xs text-[#444] hover:text-[#888] transition-colors"
            >
              Docs
            </a>
          </div>
          <span className="font-mono text-xs text-[#333]">MIT License</span>
        </div>
      </footer>
    </>
  );
}
