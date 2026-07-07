"use client";

import { motion, useInView } from "framer-motion";
import { useRef } from "react";
import TerminalWindow from "./TerminalWindow";

export default function TerminalShowcase() {
  const ref = useRef<HTMLDivElement>(null);
  const inView = useInView(ref, { once: true, margin: "-80px 0px" });

  return (
    <section className="divider section-pad" id="preview">
      <div className="container-page">
        <div className="grid lg:grid-cols-[1fr_1.4fr] gap-12 lg:gap-20 items-center">
          {/* Left: copy */}
          <motion.div
            ref={ref}
            initial={{ opacity: 0, x: -24 }}
            animate={inView ? { opacity: 1, x: 0 } : {}}
            transition={{ duration: 0.7, ease: [0.22, 1, 0.36, 1] }}
          >
            <p className="font-mono text-xs text-[#444] uppercase tracking-[0.2em] mb-5">
              Terminal preview
            </p>
            <h2 className="font-mono text-display-md text-white font-semibold mb-6">
              The whole experience lives in your terminal.
            </h2>
            <p className="text-body-md text-[#666] leading-relaxed mb-8">
              Moodwave renders entirely in your shell — no browser, no Electron, no GUI.
              The interactive TUI shows the scan result, mood confidence, now playing track,
              playback progress, and a live waveform.
            </p>

            {/* Stats */}
            <div className="grid grid-cols-2 gap-4">
              {[
                { label: "music sources", value: "3" },
                { label: "mood profiles", value: "10" },
                { label: "visual themes", value: "6" },
                { label: "binary size", value: "~8 MB" },
              ].map((stat) => (
                <div key={stat.label} className="border border-white/[0.06] rounded-lg p-4">
                  <div className="font-mono text-2xl font-semibold text-white mb-1">
                    {stat.value}
                  </div>
                  <div className="font-mono text-xs text-[#555] uppercase tracking-widest">
                    {stat.label}
                  </div>
                </div>
              ))}
            </div>
          </motion.div>

          {/* Right: terminal */}
          <motion.div
            initial={{ opacity: 0, x: 24 }}
            animate={inView ? { opacity: 1, x: 0 } : {}}
            transition={{ duration: 0.7, delay: 0.15, ease: [0.22, 1, 0.36, 1] }}
          >
            <TerminalWindow />
          </motion.div>
        </div>
      </div>
    </section>
  );
}
