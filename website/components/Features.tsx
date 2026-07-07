"use client";

import { motion } from "framer-motion";

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
        <motion.div
          initial={{ opacity: 0, y: 16 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true, margin: "-60px" }}
          transition={{ duration: 0.55 }}
          className="mb-14"
        >
          <p className="font-mono text-xs text-[#444] uppercase tracking-[0.2em] mb-4">Features</p>
          <h2 className="font-mono font-semibold text-white max-w-sm" style={{ fontSize: "clamp(1.4rem, 2.5vw, 2rem)", letterSpacing: "-0.025em" }}>
            Built for developers who don&rsquo;t want to think about music.
          </h2>
        </motion.div>

        <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-px rounded-lg overflow-hidden" style={{ background: "rgba(255,255,255,0.06)" }}>
          {features.map((f, i) => (
            <motion.div
              key={f.title}
              initial={{ opacity: 0, y: 16 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true, margin: "-40px" }}
              transition={{ duration: 0.5, delay: i * 0.07, ease: [0.22, 1, 0.36, 1] }}
              className="p-8 group hover:bg-[#111] transition-colors duration-300"
              style={{ background: "#080808" }}
            >
              <div className="font-mono text-lg text-[#444] mb-5 group-hover:text-[#777] transition-colors">{f.icon}</div>
              <h3 className="font-mono text-sm font-semibold text-white mb-3">{f.title}</h3>
              <p className="text-sm text-[#666] leading-relaxed">{f.desc}</p>
            </motion.div>
          ))}
        </div>
      </div>
    </section>
  );
}
