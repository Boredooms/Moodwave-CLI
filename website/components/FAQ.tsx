"use client";

import { useState } from "react";
import { motion, AnimatePresence } from "framer-motion";
import SplitText from "./ui/SplitText";
import FadeIn from "./ui/FadeIn";

const faqs = [
  { q: "Is the CLI heavy? Will it slow down my machine?", a: "No. Moodwave is a single compiled Go binary of ~8 MB. It has zero background processes and only runs when you invoke it. It uses your system's native audio player (mpv, ffplay) and does not transcode or cache audio locally." },
  { q: "Does it store music files locally?", a: "No. Moodwave streams audio from the internet. Nothing is cached or stored beyond your configuration file and a small track metadata index used for recommendations." },
  { q: "Does it work on Windows?", a: "Yes. Windows is a first-class target. The installer (install.ps1) places the binary in a local AppData path and adds it to PATH. Playback uses Windows PowerShell's built-in Media.SoundPlayer or ffplay if installed." },
  { q: "What music sources does it use?", a: "Three primary sources: YouTube (via yt-dlp for stream extraction), Jamendo (Creative Commons licensed music with a public API), and Radio Browser (30,000+ community-curated internet radio stations)." },
  { q: "Is a YouTube or music API key required?", a: "No API keys are needed. Moodwave uses yt-dlp (auto-downloaded on first run), Jamendo's public client ID, and Radio Browser's open REST API." },
  { q: "Is the website the product?", a: "No. This page is purely promotional. The product is entirely in your terminal. Installing the CLI is the only thing this page is here to help you do." },
];

function FAQItem({ q, a, index }: { q: string; a: string; index: number }) {
  const [open, setOpen] = useState(false);

  return (
    <div className="border-b border-white/[0.06] last:border-b-0">
      <button
        onClick={() => setOpen((v) => !v)}
        className="w-full flex items-start justify-between gap-6 py-5 text-left group cursor-pointer"
        aria-expanded={open}
      >
        <span className="font-mono text-sm text-[#999] group-hover:text-white transition-colors duration-200 leading-relaxed">
          {q}
        </span>
        <span className={`font-mono text-xl text-[#444] group-hover:text-[#666] flex-shrink-0 mt-0.5 transition-transform duration-300 ${open ? "rotate-45" : "rotate-0"}`}>
          +
        </span>
      </button>

      <AnimatePresence initial={false}>
        {open && (
          <motion.div
            initial={{ height: 0, opacity: 0 }}
            animate={{ height: "auto", opacity: 1 }}
            exit={{ height: 0, opacity: 0 }}
            transition={{ duration: 0.25, ease: [0.16, 1, 0.3, 1] }}
            className="overflow-hidden"
          >
            <p className="pb-5 text-sm text-[#666] leading-relaxed max-w-2xl">{a}</p>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}

export default function FAQ() {
  return (
    <section className="divider section-pad" id="faq">
      <div className="container-page">
        <div className="grid lg:grid-cols-[240px_1fr] gap-12 lg:gap-20">
          <div>
            <p className="font-mono text-xs text-[#444] uppercase tracking-[0.2em] mb-5">
              <SplitText text="FAQ" by="chars" delay={0.1} />
            </p>
            <h2 className="font-mono font-semibold text-white" style={{ fontSize: "clamp(1.4rem, 2.5vw, 2rem)", letterSpacing: "-0.025em" }}>
              <SplitText text="Common questions." by="words" delay={0.25} />
            </h2>
          </div>

          <div className="border border-white/[0.07] rounded-lg px-6 lg:px-8">
            {faqs.map((item, i) => (
              <FadeIn key={item.q} delay={0.1 + i * 0.06} y={15} className="w-full">
                <FAQItem q={item.q} a={item.a} index={i} />
              </FadeIn>
            ))}
          </div>
        </div>
      </div>
    </section>
  );
}
