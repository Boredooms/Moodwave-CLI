"use client";

import { motion, useInView } from "framer-motion";
import { useRef } from "react";

const features = [
  {
    icon: "◈",
    title: "Codebase mood scan",
    desc: "Analyzes file types, TODO density, git markers, and code patterns to classify your current cognitive state into one of 10 developer moods.",
  },
  {
    icon: "⋯",
    title: "Music matching engine",
    desc: "A BPM-aware, tag-based recommender ranks tracks against your detected mood profile, weighted by energy, genre, and session history.",
  },
  {
    icon: "▶",
    title: "Stream from 3 sources",
    desc: "YouTube (via yt-dlp), Jamendo Creative Commons catalog, and 30,000+ Internet Radio stations via Radio Browser.",
  },
  {
    icon: "▒",
    title: "Live terminal visuals",
    desc: "Six visual themes — waveform, spectrum, pulse, vinyl, minimal, quiet — animated in real time during playback.",
  },
  {
    icon: "↺",
    title: "Self-updating",
    desc: "Checks GitHub Releases on launch and swaps the binary in-place. Always on the latest version, zero manual updates.",
  },
  {
    icon: "⊞",
    title: "Cross-platform binary",
    desc: "Single compiled binary for Windows (amd64/arm64), macOS (Intel/Apple Silicon), and Linux (amd64/arm64/arm). No Go runtime needed.",
  },
];

export default function Features() {
  const ref = useRef<HTMLDivElement>(null);
  const inView = useInView(ref, { once: true, margin: "-80px 0px" });

  return (
    <section className="divider section-pad" id="features">
      <div className="container-page">
        <motion.div
          ref={ref}
          initial={{ opacity: 0, y: 20 }}
          animate={inView ? { opacity: 1, y: 0 } : {}}
          transition={{ duration: 0.6 }}
          className="mb-16"
        >
          <p className="font-mono text-xs text-[#444] uppercase tracking-[0.2em] mb-4">Features</p>
          <h2 className="font-mono text-display-md text-white font-semibold max-w-md">
            Built for developers who don&rsquo;t want to think about music.
          </h2>
        </motion.div>

        <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-px bg-white/[0.06] rounded-lg overflow-hidden">
          {features.map((f, i) => (
            <motion.div
              key={f.title}
              initial={{ opacity: 0, y: 20 }}
              animate={inView ? { opacity: 1, y: 0 } : {}}
              transition={{ duration: 0.55, delay: 0.05 + i * 0.08, ease: [0.22, 1, 0.36, 1] }}
              className="bg-bg p-8 group hover:bg-surface transition-colors duration-300"
            >
              <div className="font-mono text-lg text-[#555] mb-5 group-hover:text-[#888] transition-colors">
                {f.icon}
              </div>
              <h3 className="font-mono text-sm font-semibold text-white mb-3">{f.title}</h3>
              <p className="text-sm text-[#666] leading-relaxed">{f.desc}</p>
            </motion.div>
          ))}
        </div>
      </div>
    </section>
  );
}
