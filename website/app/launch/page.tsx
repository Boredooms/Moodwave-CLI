"use client";

import { useEffect, useState, useRef, useCallback } from "react";
import { motion, AnimatePresence } from "framer-motion";

// ─────────────────────────────────────────────────────────────────────────────
// LAUNCH COUNTDOWN PAGE — Premium dark aesthetic with flowing silk gradient
// ─────────────────────────────────────────────────────────────────────────────

const UNLOCK_UTC = new Date("2026-09-03T00:30:00Z"); // 06:00 AM IST

interface TimeLeft {
  days: number;
  hours: number;
  minutes: number;
  seconds: number;
  total: number;
}

function getTimeLeft(): TimeLeft {
  const total = Math.max(0, UNLOCK_UTC.getTime() - Date.now());
  return {
    days: Math.floor(total / (1000 * 60 * 60 * 24)),
    hours: Math.floor((total / (1000 * 60 * 60)) % 24),
    minutes: Math.floor((total / (1000 * 60)) % 60),
    seconds: Math.floor((total / 1000) % 60),
    total,
  };
}

// ─── Flip Clock Digit ───────────────────────────────────────────────────────

function FlipUnit({ value, label }: { value: number; label: string }) {
  const display = String(value).padStart(2, "0");
  const top = display;
  const bottom = display;

  return (
    <div className="flex flex-col items-center gap-2 sm:gap-3">
      <div className="relative">
        <AnimatePresence mode="popLayout">
          <motion.div
            key={display}
            initial={{ rotateX: -80, opacity: 0, scale: 0.9 }}
            animate={{ rotateX: 0, opacity: 1, scale: 1 }}
            exit={{ rotateX: 80, opacity: 0, scale: 0.9 }}
            transition={{ duration: 0.6, ease: [0.16, 1, 0.3, 1] }}
            className="relative w-[56px] h-[72px] sm:w-[80px] sm:h-[100px] md:w-[96px] md:h-[120px]"
            style={{ perspective: "800px", transformStyle: "preserve-3d" }}
          >
            {/* Card */}
            <div className="absolute inset-0 rounded-xl overflow-hidden border border-white/[0.06] shadow-[0_16px_48px_rgba(0,0,0,0.5)]">
              {/* Top half */}
              <div className="absolute inset-x-0 top-0 h-1/2 bg-gradient-to-b from-[#1a1a1a] to-[#141414] flex items-end justify-center pb-0">
                <span className="font-mono text-2xl sm:text-4xl md:text-5xl font-bold text-white translate-y-[55%]">
                  {top}
                </span>
              </div>
              {/* Bottom half */}
              <div className="absolute inset-x-0 bottom-0 h-1/2 bg-gradient-to-b from-[#111] to-[#0d0d0d] flex items-start justify-center pt-0">
                <span className="font-mono text-2xl sm:text-4xl md:text-5xl font-bold text-white/80 -translate-y-[55%]">
                  {bottom}
                </span>
              </div>
              {/* Center fold line */}
              <div className="absolute top-1/2 inset-x-0 h-[1px] bg-black/40" />
              <div className="absolute top-1/2 inset-x-0 h-[1px] translate-y-[1px] bg-white/[0.03]" />
            </div>
          </motion.div>
        </AnimatePresence>
      </div>
      <span className="font-mono text-[9px] sm:text-[10px] text-white/30 uppercase tracking-[0.2em]">
        {label}
      </span>
    </div>
  );
}

// ─── Waitlist Form ──────────────────────────────────────────────────────────

function WaitlistForm() {
  const [email, setEmail] = useState("");
  const [status, setStatus] = useState<"idle" | "loading" | "success" | "error">("idle");
  const [message, setMessage] = useState("");
  const [count, setCount] = useState<number | null>(null);

  const fetchCount = useCallback(async () => {
    try {
      const res = await fetch("/api/waitlist");
      const data = await res.json();
      if (typeof data.count === "number") setCount(data.count);
    } catch { /* silent */ }
  }, []);

  useEffect(() => { fetchCount(); }, [fetchCount]);

  const handleSubmit = useCallback(
    async (e: React.FormEvent) => {
      e.preventDefault();
      if (!email.trim() || status === "loading") return;

      setStatus("loading");
      try {
        const res = await fetch("/api/waitlist", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ email: email.trim() }),
        });
        const data = await res.json();
        if (res.ok) {
          setStatus("success");
          setMessage(data.message || "You're in!");
          setEmail("");
          if (typeof data.count === "number") setCount(data.count);
          else fetchCount();
        } else {
          setStatus("error");
          setMessage(data.error || "Something went wrong.");
        }
      } catch {
        setStatus("error");
        setMessage("Network error. Try again.");
      }
    },
    [email, status, fetchCount]
  );

  return (
    <div className="w-full max-w-md mx-auto">
      <form onSubmit={handleSubmit} className="relative">
        <div className="flex items-center bg-white/[0.04] border border-white/[0.08] rounded-full overflow-hidden transition-all focus-within:border-white/20 focus-within:bg-white/[0.06]">
          {/* Email icon */}
          <div className="pl-4 sm:pl-5 flex-shrink-0">
            <svg className="w-4 h-4 text-white/30" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M21.75 6.75v10.5a2.25 2.25 0 01-2.25 2.25h-15a2.25 2.25 0 01-2.25-2.25V6.75m19.5 0A2.25 2.25 0 0019.5 4.5h-15a2.25 2.25 0 00-2.25 2.25m19.5 0v.243a2.25 2.25 0 01-1.07 1.916l-7.5 4.615a2.25 2.25 0 01-2.36 0L3.32 8.91a2.25 2.25 0 01-1.07-1.916V6.75" />
            </svg>
          </div>
          <input
            type="email"
            value={email}
            onChange={(e) => {
              setEmail(e.target.value);
              if (status === "error") setStatus("idle");
            }}
            placeholder="your@email.com"
            required
            disabled={status === "loading" || status === "success"}
            className="flex-1 bg-transparent px-3 py-3.5 sm:py-4 font-mono text-sm text-white placeholder-white/25 outline-none disabled:opacity-50"
          />
          <button
            type="submit"
            disabled={status === "loading" || status === "success"}
            className="mr-1.5 px-4 sm:px-6 py-2 sm:py-2.5 bg-white text-black font-mono text-xs sm:text-sm font-semibold rounded-full hover:bg-white/90 transition-all disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer whitespace-nowrap"
          >
            {status === "loading" ? (
              <span className="inline-block w-4 h-4 border-2 border-black/20 border-t-black rounded-full animate-spin" />
            ) : status === "success" ? (
              "✓ Joined"
            ) : (
              "Join Waitlist"
            )}
          </button>
        </div>
      </form>

      <AnimatePresence>
        {message && (
          <motion.p
            initial={{ opacity: 0, y: -8 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -8 }}
            transition={{ duration: 0.3, ease: [0.16, 1, 0.3, 1] }}
            className={`mt-4 text-xs font-mono text-center ${
              status === "success" ? "text-emerald-400/80" : "text-red-400/80"
            }`}
          >
            {message}
          </motion.p>
        )}
      </AnimatePresence>

      {count !== null && count > 0 && (
        <motion.p
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.5 }}
          className="mt-4 text-center"
        >
          <span className="inline-flex items-center gap-2 px-3 py-1.5 bg-white/[0.03] border border-white/[0.06] rounded-full">
            <span className="w-1.5 h-1.5 bg-emerald-400 rounded-full animate-pulse" />
            <span className="font-mono text-[11px] text-white/40">
              <span className="text-white/70 font-semibold">{count.toLocaleString()}</span>{" "}
              {count === 1 ? "person" : "people"} waiting
            </span>
          </span>
        </motion.p>
      )}
    </div>
  );
}

// ─── Main Page ──────────────────────────────────────────────────────────────

export default function LaunchPage() {
  const [time, setTime] = useState<TimeLeft>(getTimeLeft());
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const id = setInterval(() => setTime(getTimeLeft()), 1000);
    return () => clearInterval(id);
  }, []);

  useEffect(() => {
    if (time.total <= 0) {
      window.location.href = "/";
    }
  }, [time.total]);

  return (
    <div
      ref={containerRef}
      className="relative min-h-screen flex flex-col items-center justify-center overflow-hidden"
      style={{ background: "#080808" }}
    >
      {/* Silk/fabric gradient — the signature visual like Resend's dark page */}
      <div className="absolute inset-0 overflow-hidden pointer-events-none">
        {/* Large flowing gradient blob — top right */}
        <div
          className="absolute -top-[30%] -right-[20%] w-[80%] h-[100%] opacity-[0.12]"
          style={{
            background: "conic-gradient(from 180deg at 50% 50%, #ffffff 0deg, #888888 60deg, #333333 120deg, #666666 180deg, #ffffff 240deg, #aaaaaa 300deg, #444444 360deg)",
            filter: "blur(100px)",
            borderRadius: "50%",
            transform: "rotate(-20deg)",
          }}
        />
        {/* Secondary gradient blob — bottom left */}
        <div
          className="absolute -bottom-[40%] -left-[20%] w-[70%] h-[90%] opacity-[0.06]"
          style={{
            background: "conic-gradient(from 0deg at 50% 50%, #888888 0deg, #ffffff 90deg, #444444 180deg, #aaaaaa 270deg, #888888 360deg)",
            filter: "blur(120px)",
            borderRadius: "50%",
            transform: "rotate(30deg)",
          }}
        />
        {/* Subtle noise overlay for texture */}
        <div
          className="absolute inset-0 opacity-[0.03]"
          style={{
            backgroundImage: `url("data:image/svg+xml,%3Csvg viewBox='0 0 256 256' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='noise'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.65' numOctaves='3' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23noise)' opacity='1'/%3E%3C/svg%3E")`,
          }}
        />
      </div>

      <div className="relative z-10 flex flex-col items-center text-center px-6 max-w-2xl w-full">
        {/* Logo */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 1, delay: 0.1, ease: [0.16, 1, 0.3, 1] }}
          className="mb-10 sm:mb-12"
        >
          <div className="inline-flex items-center justify-center w-14 h-14 rounded-2xl border border-white/[0.08] bg-white/[0.03] backdrop-blur-sm">
            <svg className="w-7 h-7 text-white" viewBox="0 0 32 32" fill="none">
              <path
                d="M9 16h1.5M13 11v10M17 7v18M21 13v6M25 16h-1.5"
                stroke="currentColor"
                strokeWidth="2.5"
                strokeLinecap="round"
                strokeLinejoin="round"
              />
            </svg>
          </div>
        </motion.div>

        {/* Title */}
        <motion.h1
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 1, delay: 0.25, ease: [0.16, 1, 0.3, 1] }}
          className="font-mono font-bold text-white text-4xl sm:text-5xl md:text-6xl tracking-tight mb-4"
        >
          moodwave
        </motion.h1>

        <motion.p
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 1, delay: 0.4, ease: [0.16, 1, 0.3, 1] }}
          className="text-white/40 text-sm sm:text-base max-w-sm leading-relaxed mb-14 sm:mb-16"
        >
          A terminal-native music companion that reads your code and plays the perfect soundtrack.
        </motion.p>

        {/* Countdown */}
        <motion.div
          initial={{ opacity: 0, y: 30 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 1, delay: 0.55, ease: [0.16, 1, 0.3, 1] }}
          className="flex items-start gap-2 sm:gap-4 md:gap-5 mb-6"
        >
          <FlipUnit value={time.days} label="Days" />
          <span className="text-white/10 text-xl sm:text-3xl font-light mt-5 sm:mt-8 md:mt-10">:</span>
          <FlipUnit value={time.hours} label="Hours" />
          <span className="text-white/10 text-xl sm:text-3xl font-light mt-5 sm:mt-8 md:mt-10">:</span>
          <FlipUnit value={time.minutes} label="Min" />
          <span className="text-white/10 text-xl sm:text-3xl font-light mt-5 sm:mt-8 md:mt-10">:</span>
          <FlipUnit value={time.seconds} label="Sec" />
        </motion.div>

        <motion.p
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.8 }}
          className="font-mono text-[10px] sm:text-[11px] text-white/20 uppercase tracking-[0.25em] mb-14 sm:mb-16"
        >
          September 3, 2026
        </motion.p>

        {/* Waitlist */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 1, delay: 0.9, ease: [0.16, 1, 0.3, 1] }}
          className="w-full mb-16 sm:mb-20"
        >
          <p className="font-mono text-[11px] text-white/30 mb-5 uppercase tracking-[0.15em]">
            Get notified on launch day
          </p>
          <WaitlistForm />
        </motion.div>

        {/* Feature pills */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 1, delay: 1.1, ease: [0.16, 1, 0.3, 1] }}
          className="flex flex-wrap items-center justify-center gap-2 sm:gap-3"
        >
          {[
            { icon: "◈", label: "Mood Detection" },
            { icon: "▶", label: "Terminal Playback" },
            { icon: "◆", label: "24 Visualizers" },
            { icon: "☰", label: "Smart Playlists" },
            { icon: "⬆", label: "Auto Updates" },
          ].map((f, i) => (
            <motion.span
              key={f.label}
              initial={{ opacity: 0, scale: 0.9 }}
              animate={{ opacity: 1, scale: 1 }}
              transition={{ delay: 1.2 + i * 0.08, duration: 0.5, ease: [0.16, 1, 0.3, 1] }}
              className="inline-flex items-center gap-2 px-3 py-1.5 bg-white/[0.03] border border-white/[0.06] rounded-full"
            >
              <span className="text-white/50 text-xs">{f.icon}</span>
              <span className="font-mono text-[10px] text-white/30">{f.label}</span>
            </motion.span>
          ))}
        </motion.div>
      </div>

      {/* Footer */}
      <motion.p
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 1.8 }}
        className="absolute bottom-5 font-mono text-[10px] text-white/15 z-10"
      >
        © 2026 moodwave
      </motion.p>
    </div>
  );
}
