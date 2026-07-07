"use client";

import { motion, useInView } from "framer-motion";
import { useRef } from "react";

const steps = [
  {
    num: "01",
    title: "Install",
    desc: "One curl command or PowerShell snippet. No Go runtime needed. The binary lands in your PATH.",
    code: "curl ... | bash",
  },
  {
    num: "02",
    title: "Scan",
    desc: "Moodwave walks your repository — files, languages, git state, TODOs, code patterns — and builds a semantic profile.",
    code: "moodwave scan",
  },
  {
    num: "03",
    title: "Infer Mood",
    desc: "A weighted heuristic engine maps the code signals to one of 10 developer moods — Focus, Debug, Sprint, Flow, and more.",
    code: "mood: DEBUGGING (87%)",
  },
  {
    num: "04",
    title: "Play",
    desc: "The recommendation engine fetches music matching your mood from YouTube, Jamendo, or Internet Radio and streams it in the terminal.",
    code: "▶  Lo-Fi Hip Hop Radio",
  },
];

export default function HowItWorks() {
  const ref = useRef<HTMLDivElement>(null);
  const inView = useInView(ref, { once: true, margin: "-80px 0px" });

  return (
    <section className="divider section-pad" id="how-it-works">
      <div className="container-page">
        <motion.div
          ref={ref}
          initial={{ opacity: 0, y: 20 }}
          animate={inView ? { opacity: 1, y: 0 } : {}}
          transition={{ duration: 0.6 }}
          className="mb-16"
        >
          <p className="font-mono text-xs text-[#444] uppercase tracking-[0.2em] mb-4">How it works</p>
          <h2 className="font-mono text-display-md text-white font-semibold">
            Four steps. One command.
          </h2>
        </motion.div>

        <div className="grid sm:grid-cols-2 lg:grid-cols-4 gap-px bg-white/[0.06] rounded-lg overflow-hidden">
          {steps.map((step, i) => (
            <motion.div
              key={step.num}
              initial={{ opacity: 0, y: 24 }}
              animate={inView ? { opacity: 1, y: 0 } : {}}
              transition={{ duration: 0.6, delay: 0.1 + i * 0.12, ease: [0.22, 1, 0.36, 1] }}
              className="bg-bg p-8 lg:p-10 group hover:bg-surface transition-colors duration-300"
            >
              <div className="font-mono text-5xl font-semibold text-[#1a1a1a] mb-6 leading-none group-hover:text-[#252525] transition-colors">
                {step.num}
              </div>
              <h3 className="font-mono text-base font-semibold text-white mb-3">{step.title}</h3>
              <p className="text-sm text-[#666] leading-relaxed mb-6">{step.desc}</p>
              <div className="font-mono text-xs text-[#444] bg-[#0d0d0d] border border-white/[0.06] rounded px-3 py-2 group-hover:border-white/[0.1] transition-colors">
                {step.code}
              </div>
            </motion.div>
          ))}
        </div>
      </div>
    </section>
  );
}
