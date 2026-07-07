"use client";

import { motion } from "framer-motion";

const steps = [
  { num: "01", title: "Install", desc: "One curl command or PowerShell snippet. No Go runtime needed. The binary lands in your PATH.", code: "curl ... | bash" },
  { num: "02", title: "Scan", desc: "Moodwave walks your repository — files, languages, git state, TODOs, code patterns — and builds a semantic profile.", code: "moodwave scan" },
  { num: "03", title: "Infer Mood", desc: "A weighted heuristic engine maps code signals to one of 10 developer moods — Focus, Debug, Sprint, Flow, and more.", code: "mood: DEBUGGING (87%)" },
  { num: "04", title: "Play", desc: "The recommender fetches music matching your mood from YouTube, Jamendo, or Radio Browser and streams it in terminal.", code: "▶  Lo-Fi Hip Hop Radio" },
];

export default function HowItWorks() {
  return (
    <section className="divider section-pad" id="how-it-works">
      <div className="container-page">
        <motion.div
          initial={{ opacity: 0, y: 16 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true, margin: "-60px" }}
          transition={{ duration: 0.55 }}
          className="mb-14"
        >
          <p className="font-mono text-xs text-[#444] uppercase tracking-[0.2em] mb-4">How it works</p>
          <h2 className="font-mono font-semibold text-white" style={{ fontSize: "clamp(1.4rem, 2.5vw, 2rem)", letterSpacing: "-0.025em" }}>
            Four steps. One command.
          </h2>
        </motion.div>

        <div className="grid sm:grid-cols-2 lg:grid-cols-4 gap-px rounded-lg overflow-hidden" style={{ background: "rgba(255,255,255,0.06)" }}>
          {steps.map((step, i) => (
            <motion.div
              key={step.num}
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true, margin: "-40px" }}
              transition={{ duration: 0.55, delay: i * 0.1, ease: [0.22, 1, 0.36, 1] }}
              className="p-8 lg:p-10 group hover:bg-[#111] transition-colors duration-300"
              style={{ background: "#080808" }}
            >
              <div className="font-mono font-semibold text-[#1a1a1a] mb-6 leading-none group-hover:text-[#262626] transition-colors" style={{ fontSize: "3.5rem" }}>
                {step.num}
              </div>
              <h3 className="font-mono text-sm font-semibold text-white mb-3">{step.title}</h3>
              <p className="text-sm text-[#666] leading-relaxed mb-6">{step.desc}</p>
              <div className="font-mono text-xs text-[#444] rounded px-3 py-2 border border-white/[0.06] group-hover:border-white/[0.1] transition-colors" style={{ background: "#0d0d0d" }}>
                {step.code}
              </div>
            </motion.div>
          ))}
        </div>
      </div>
    </section>
  );
}
