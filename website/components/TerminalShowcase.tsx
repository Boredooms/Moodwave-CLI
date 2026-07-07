"use client";

import TerminalWindow from "./TerminalWindow";
import SplitText from "./ui/SplitText";
import FadeIn from "./ui/FadeIn";

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
          <div>
            <p className="font-mono text-xs text-[#444] uppercase tracking-[0.2em] mb-5">
              <SplitText text="Terminal preview" by="chars" delay={0.1} />
            </p>
            <h2 className="font-mono font-semibold text-white mb-6" style={{ fontSize: "clamp(1.4rem, 2.5vw, 2rem)", letterSpacing: "-0.025em" }}>
              <SplitText text="The whole experience lives in your terminal." by="words" delay={0.25} />
            </h2>
            <FadeIn delay={0.4} y={15}>
              <p className="text-sm text-[#666] leading-relaxed mb-8">
                Moodwave renders entirely in your shell — no browser, no Electron, no GUI.
                The interactive TUI shows the scan result, mood confidence, now playing track,
                playback progress, and a live waveform.
              </p>
            </FadeIn>

            <div className="grid grid-cols-2 gap-3">
              {stats.map((stat, i) => (
                <FadeIn key={stat.label} delay={0.45 + i * 0.08} y={15}>
                  <div className="border border-white/[0.06] rounded-lg p-4">
                    <div className="font-mono text-2xl font-semibold text-white mb-1">{stat.value}</div>
                    <div className="font-mono text-xs text-[#555] uppercase tracking-widest">{stat.label}</div>
                  </div>
                </FadeIn>
              ))}
            </div>
          </div>

          <FadeIn delay={0.2} x={30} y={0} duration={0.8}>
            <TerminalWindow />
          </FadeIn>
        </div>
      </div>
    </section>
  );
}
