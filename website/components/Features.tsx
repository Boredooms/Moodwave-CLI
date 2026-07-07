"use client";

import SplitText from "./ui/SplitText";
import FadeIn from "./ui/FadeIn";

const features = [
  { icon: "◈", title: "Codebase mood scan", desc: "Analyzes file types, TODO density, git markers, and code patterns to classify your cognitive state into one of 10 developer moods." },
  { icon: "⋯", title: "Music matching engine", desc: "BPM-aware, tag-based recommender ranks tracks against your mood profile, weighted by energy, genre, and session history." },
  { icon: "▶", title: "Stream from 3 sources", desc: "YouTube (via yt-dlp), Jamendo Creative Commons catalog, and 30,000+ Internet Radio stations via Radio Browser." },
  { icon: "▒", title: "Live terminal visuals", desc: "Six visual themes — waveform, spectrum, pulse, vinyl, minimal, quiet — animated in real time during playback." },
  { icon: "↺", title: "Self-updating", desc: "Checks GitHub Releases on launch and swaps the binary in-place. Always on the latest version, zero manual updates." },
  { icon: "⊞", title: "Cross-platform binary", desc: "Single binary for Windows (amd64/arm64), macOS (Intel/M-series), and Linux (amd64/arm64/arm). No Go runtime needed." },
];

export default function Features() {
  return (
    <section className="divider section-pad" id="features">
      <div className="container-page">
        <div className="mb-14">
          <p className="font-mono text-xs text-[#444] uppercase tracking-[0.2em] mb-4">
            <SplitText text="Features" by="chars" delay={0.1} />
          </p>
          <h2 className="font-mono font-semibold text-white max-w-sm" style={{ fontSize: "clamp(1.4rem, 2.5vw, 2rem)", letterSpacing: "-0.025em" }}>
            <SplitText text="Built for developers who don't want to think about music." by="words" delay={0.2} stagger={0.04} />
          </h2>
        </div>

        <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-px rounded-lg overflow-hidden" style={{ background: "rgba(255,255,255,0.06)" }}>
          {features.map((f, i) => (
            <FadeIn
              key={f.title}
              delay={0.1 + i * 0.08}
              y={20}
              className="w-full h-full"
            >
              <div
                className="p-8 group hover:bg-[#111] transition-colors duration-300 h-full"
                style={{ background: "#080808" }}
              >
                <div className="font-mono text-lg text-[#444] mb-5 group-hover:text-[#777] transition-colors">{f.icon}</div>
                <h3 className="font-mono text-sm font-semibold text-white mb-3">{f.title}</h3>
                <p className="text-sm text-[#666] leading-relaxed">{f.desc}</p>
              </div>
            </FadeIn>
          ))}
        </div>
      </div>
    </section>
  );
}
