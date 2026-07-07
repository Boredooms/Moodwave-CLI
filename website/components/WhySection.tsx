"use client";

import { motion, useInView } from "framer-motion";
import { useRef } from "react";

const statements = [
  {
    text: "Most developers listen to music while coding. But the playlist never quite fits — too energetic for deep focus, too mellow when you're debugging a race condition at midnight.",
  },
  {
    text: "Moodwave reads your repository. It looks at what you're building, how many open TODOs are piling up, what languages are active, how your git tree looks — and it infers the kind of music that actually matches where your head is.",
  },
  {
    text: "No playlist curation. No manual mode switching. Just run `moodwave` and let it figure out the rest.",
  },
];

export default function WhySection() {
  const ref = useRef<HTMLDivElement>(null);
  const inView = useInView(ref, { once: true, margin: "-100px 0px" });

  return (
    <section className="divider section-pad" id="why">
      <div className="container-page">
        <div className="grid lg:grid-cols-[280px_1fr] gap-12 lg:gap-20">
          {/* Label column */}
          <motion.div
            ref={ref}
            initial={{ opacity: 0, y: 16 }}
            animate={inView ? { opacity: 1, y: 0 } : {}}
            transition={{ duration: 0.6 }}
          >
            <p className="font-mono text-xs text-[#444] uppercase tracking-[0.2em]">
              Why Moodwave
            </p>
          </motion.div>

          {/* Content column */}
          <div className="space-y-8 max-w-2xl">
            {statements.map((s, i) => (
              <motion.p
                key={i}
                initial={{ opacity: 0, y: 20 }}
                animate={inView ? { opacity: 1, y: 0 } : {}}
                transition={{ duration: 0.65, delay: 0.15 + i * 0.15, ease: [0.22, 1, 0.36, 1] }}
                className={`text-body-lg leading-relaxed ${
                  i === 0 ? "text-[#999]" : i === 1 ? "text-[#777]" : "text-[#999] font-mono text-base"
                }`}
              >
                {s.text}
              </motion.p>
            ))}
          </div>
        </div>
      </div>
    </section>
  );
}
