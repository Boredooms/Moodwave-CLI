"use client";

import { motion } from "framer-motion";

const statements = [
  "Most developers listen to music while coding. But the playlist never quite fits — too energetic for deep focus, too mellow when you're debugging a race condition at midnight.",
  "Moodwave reads your repository. It looks at what you're building, how many open TODOs are piling up, what languages are active, how your git tree looks — and it infers the kind of music that actually matches where your head is.",
  "No playlist curation. No manual mode switching. Just run `moodwave` and let it figure out the rest.",
];

export default function WhySection() {
  return (
    <section className="divider section-pad" id="why">
      <div className="container-page">
        <div className="grid lg:grid-cols-[240px_1fr] gap-12 lg:gap-20">
          <motion.div
            initial={{ opacity: 0, y: 16 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true, margin: "-60px" }}
            transition={{ duration: 0.55 }}
          >
            <p className="font-mono text-xs text-[#444] uppercase tracking-[0.2em]">Why Moodwave</p>
          </motion.div>

          <div className="space-y-8 max-w-2xl">
            {statements.map((s, i) => (
              <motion.p
                key={i}
                initial={{ opacity: 0, y: 16 }}
                whileInView={{ opacity: 1, y: 0 }}
                viewport={{ once: true, margin: "-60px" }}
                transition={{ duration: 0.6, delay: i * 0.12, ease: [0.22, 1, 0.36, 1] }}
                className={`leading-relaxed ${i === 2 ? "font-mono text-sm text-[#888]" : "text-[#777]"}`}
                style={{ fontSize: i === 2 ? undefined : "1.05rem" }}
              >
                {s}
              </motion.p>
            ))}
          </div>
        </div>
      </div>
    </section>
  );
}
