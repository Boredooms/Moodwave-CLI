"use client";

import { motion } from "framer-motion";
import TerminalWindow from "./TerminalWindow";

const stats = [
  { label: "music sources", value: "3" },
  { label: "mood profiles", value: "10" },
  { label: "visual themes", value: "6" },
  { label: "binary size", value: "~8 MB" },
];

export default function TerminalShowcase() {
  return (
    <section className="divider section-pad" id="preview">
      <div className="container-page">
        <div className="grid lg:grid-cols-[1fr_1.3fr] gap-12 lg:gap-20 items-center">
          <motion.div
            initial={{ opacity: 0, x: -20 }}
            whileInView={{ opacity: 1, x: 0 }}
            viewport={{ once: true, margin: "-60px" }}
            transition={{ duration: 0.65, ease: [0.22, 1, 0.36, 1] }}
          >
            <p className="font-mono text-xs text-[#444] uppercase tracking-[0.2em] mb-5">Terminal preview</p>
            <h2 className="font-mono font-semibold text-white mb-6" style={{ fontSize: "clamp(1.4rem, 2.5vw, 2rem)", letterSpacing: "-0.025em" }}>
              The whole experience lives in your terminal.
            </h2>
            <p className="text-sm text-[#666] leading-relaxed mb-8">
              Moodwave renders entirely in your shell — no browser, no Electron, no GUI.
              The interactive TUI shows the scan result, mood confidence, now playing track,
              playback progress, and a live waveform.
            </p>

            <div className="grid grid-cols-2 gap-3">
              {stats.map((stat) => (
                <div key={stat.label} className="border border-white/[0.06] rounded-lg p-4">
                  <div className="font-mono text-2xl font-semibold text-white mb-1">{stat.value}</div>
                  <div className="font-mono text-xs text-[#555] uppercase tracking-widest">{stat.label}</div>
                </div>
              ))}
            </div>
          </motion.div>

          <motion.div
            initial={{ opacity: 0, x: 20 }}
            whileInView={{ opacity: 1, x: 0 }}
            viewport={{ once: true, margin: "-60px" }}
            transition={{ duration: 0.65, delay: 0.12, ease: [0.22, 1, 0.36, 1] }}
          >
            <TerminalWindow />
          </motion.div>
        </div>
      </div>
    </section>
  );
}
