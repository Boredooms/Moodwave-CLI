"use client";

import { motion, useInView } from "framer-motion";
import { useRef } from "react";

const layers = [
  {
    name: "Scanner",
    desc: "Walks the file tree, extracts language composition, TODO/FIXME density, git metadata, and dependency manifest signals.",
  },
  {
    name: "Mood Engine",
    desc: "Applies 13 weighted heuristic rules and a naive Bayes semantic classifier to map code signals to one of 10 developer mood profiles.",
  },
  {
    name: "Recommender",
    desc: "Ranks candidate tracks by BPM range, energy, tag overlap, and artist diversity. Boosts tracks matching the currently playing song.",
  },
  {
    name: "Source Adapters",
    desc: "Pluggable adapters for YouTube (yt-dlp), Jamendo API, and Radio Browser. Falls back automatically if a source fails.",
  },
  {
    name: "Playback Layer",
    desc: "Streams audio via mpv, ffplay, or Windows PowerShell Media.SoundPlayer. Auto-retries on network drops with Range resume.",
  },
  {
    name: "Visual Renderer",
    desc: "ANSI escape-based TUI with six animated visual modes. Single goroutine keyboard listener — no input lag.",
  },
];

export default function Architecture() {
  const ref = useRef<HTMLDivElement>(null);
  const inView = useInView(ref, { once: true, margin: "-80px 0px" });

  return (
    <section className="divider section-pad" id="architecture">
      <div className="container-page">
        <div className="grid lg:grid-cols-[280px_1fr] gap-12 lg:gap-20">
          <motion.div
            ref={ref}
            initial={{ opacity: 0, y: 20 }}
            animate={inView ? { opacity: 1, y: 0 } : {}}
            transition={{ duration: 0.6 }}
          >
            <p className="font-mono text-xs text-[#444] uppercase tracking-[0.2em] mb-5">Architecture</p>
            <h2 className="font-mono text-display-md text-white font-semibold leading-tight">
              Six layers. Clean separation.
            </h2>
            <p className="text-body-sm text-[#555] mt-4 leading-relaxed">
              Every component is independently replaceable. The scanner, mood engine, and playback layer are fully decoupled.
            </p>
          </motion.div>

          <div className="space-y-0 border border-white/[0.07] rounded-lg overflow-hidden">
            {layers.map((layer, i) => (
              <motion.div
                key={layer.name}
                initial={{ opacity: 0, x: 20 }}
                animate={inView ? { opacity: 1, x: 0 } : {}}
                transition={{ duration: 0.5, delay: 0.08 + i * 0.1, ease: [0.22, 1, 0.36, 1] }}
                className={`flex items-start gap-6 p-6 group hover:bg-surface transition-colors duration-200 ${
                  i < layers.length - 1 ? "border-b border-white/[0.06]" : ""
                }`}
              >
                {/* Layer number */}
                <div className="font-mono text-xs text-[#333] pt-0.5 w-6 flex-shrink-0 group-hover:text-[#555] transition-colors">
                  {String(i + 1).padStart(2, "0")}
                </div>
                <div className="min-w-[140px] flex-shrink-0">
                  <span className="font-mono text-sm font-semibold text-[#aaa] group-hover:text-white transition-colors">
                    {layer.name}
                  </span>
                </div>
                <p className="text-sm text-[#555] leading-relaxed group-hover:text-[#777] transition-colors">
                  {layer.desc}
                </p>
              </motion.div>
            ))}
          </div>
        </div>
      </div>
    </section>
  );
}
