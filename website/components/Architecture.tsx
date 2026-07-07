"use client";

import SplitText from "./ui/SplitText";
import FadeIn from "./ui/FadeIn";

const layers = [
  { name: "Scanner", desc: "Walks the file tree, extracts language composition, TODO/FIXME density, git metadata, and dependency manifest signals." },
  { name: "Mood Engine", desc: "Applies 13 weighted heuristic rules and a naive Bayes semantic classifier to map code signals to one of 10 developer mood profiles." },
  { name: "Recommender", desc: "Ranks candidate tracks by BPM range, energy, tag overlap, and artist diversity. Boosts tracks matching the currently playing song." },
  { name: "Source Adapters", desc: "Pluggable adapters for YouTube (yt-dlp), Jamendo API, and Radio Browser. Falls back automatically if a source fails." },
  { name: "Playback Layer", desc: "Streams audio via mpv, ffplay, or Windows PowerShell. Auto-retries on network drops with HTTP Range resume." },
  { name: "Visual Renderer", desc: "ANSI escape-based TUI with six animated visual modes. Single goroutine keyboard listener — no input lag." },
];

export default function Architecture() {
  return (
    <section className="divider section-pad" id="architecture">
      <div className="container-page">
        <div className="grid lg:grid-cols-[240px_1fr] gap-12 lg:gap-20">
          <div>
            <p className="font-mono text-xs text-[#444] uppercase tracking-[0.2em] mb-5">
              <SplitText text="Architecture" by="chars" delay={0.1} />
            </p>
            <h2 className="font-mono font-semibold text-white leading-tight" style={{ fontSize: "clamp(1.4rem, 2.5vw, 2rem)", letterSpacing: "-0.025em" }}>
              <SplitText text="Six layers." by="words" delay={0.25} />
              <br />
              <SplitText text="Clean separation." by="words" delay={0.4} />
            </h2>
            <FadeIn delay={0.55} y={15}>
              <p className="text-sm text-[#555] mt-4 leading-relaxed">
                Every component is independently replaceable.
              </p>
            </FadeIn>
          </div>

          <div className="border border-white/[0.07] rounded-lg overflow-hidden">
            {layers.map((layer, i) => (
              <FadeIn
                key={layer.name}
                delay={0.1 + i * 0.08}
                y={15}
                className="w-full"
              >
                <div
                  className={`flex items-start gap-6 p-6 group hover:bg-[#111] transition-colors duration-200 ${i < layers.length - 1 ? "border-b border-white/[0.06]" : ""}`}
                >
                  <div className="font-mono text-xs text-[#2a2a2a] pt-0.5 w-6 flex-shrink-0 group-hover:text-[#555] transition-colors">
                    {String(i + 1).padStart(2, "0")}
                  </div>
                  <div className="min-w-[130px] flex-shrink-0">
                    <span className="font-mono text-sm font-semibold text-[#999] group-hover:text-white transition-colors">{layer.name}</span>
                  </div>
                  <p className="text-sm text-[#555] leading-relaxed group-hover:text-[#777] transition-colors">{layer.desc}</p>
                </div>
              </FadeIn>
            ))}
          </div>
        </div>
      </div>
    </section>
  );
}
