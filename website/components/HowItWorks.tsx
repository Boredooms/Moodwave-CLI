"use client";

import SplitText from "./ui/SplitText";
import FadeIn from "./ui/FadeIn";

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
        <div className="mb-14">
          <p className="font-mono text-xs text-[#444] uppercase tracking-[0.2em] mb-4">
            <SplitText text="How it works" by="chars" delay={0.1} />
          </p>
          <h2 className="font-mono font-semibold text-white" style={{ fontSize: "clamp(1.4rem, 2.5vw, 2rem)", letterSpacing: "-0.025em" }}>
            <SplitText text="Four steps. One command." by="words" delay={0.25} />
          </h2>
        </div>

        <div className="grid sm:grid-cols-2 lg:grid-cols-4 gap-px rounded-lg overflow-hidden" style={{ background: "rgba(255,255,255,0.06)" }}>
          {steps.map((step, i) => (
            <FadeIn
              key={step.num}
              delay={0.1 + i * 0.1}
              y={20}
              className="w-full h-full"
            >
              <div
                className="p-8 lg:p-10 group hover:bg-[#111] transition-colors duration-300 h-full flex flex-col justify-between"
                style={{ background: "#080808" }}
              >
                <div>
                  <div className="font-mono font-semibold text-[#1a1a1a] mb-6 leading-none group-hover:text-[#262626] transition-colors" style={{ fontSize: "3.5rem" }}>
                    {step.num}
                  </div>
                  <h3 className="font-mono text-sm font-semibold text-white mb-3">{step.title}</h3>
                  <p className="text-sm text-[#666] leading-relaxed mb-6">{step.desc}</p>
                </div>
                <div className="font-mono text-xs text-[#444] rounded px-3 py-2 border border-white/[0.06] group-hover:border-white/[0.1] transition-colors" style={{ background: "#0d0d0d" }}>
                  {step.code}
                </div>
              </div>
            </FadeIn>
          ))}
        </div>
      </div>
    </section>
  );
}
