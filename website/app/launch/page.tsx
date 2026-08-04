"use client";

import { useEffect, useState, useCallback } from "react";
import { motion, AnimatePresence } from "framer-motion";
import BlackHole from "@/components/originkit/ui/blackhole";

// ─────────────────────────────────────────────────────────────────────────────
// LAUNCH PAGE — single premium screen with BlackHole background
// Clean, minimal, no clutter. Desktop-first, responsive.
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

// ─── Countdown Unit ─────────────────────────────────────────────────────────

function CountdownUnit({ value, label }: { value: number; label: string }) {
  const display = String(value).padStart(2, "0");

  return (
    <div className="flex flex-col items-center gap-2">
      <div className="relative w-[64px] h-[76px] sm:w-[80px] sm:h-[92px] md:w-[88px] md:h-[100px]">
        <AnimatePresence mode="popLayout">
          <motion.div
            key={display}
            initial={{ rotateX: -90, opacity: 0 }}
            animate={{ rotateX: 0, opacity: 1 }}
            exit={{ rotateX: 90, opacity: 0 }}
            transition={{ duration: 0.5, ease: [0.16, 1, 0.3, 1] }}
            className="absolute inset-0 flex items-center justify-center rounded-xl backdrop-blur-md"
            style={{
              background: "rgba(255,255,255,0.03)",
              border: "1px solid rgba(255,255,255,0.06)",
              boxShadow: "0 8px 32px rgba(0,0,0,0.4), inset 0 1px 0 rgba(255,255,255,0.04)",
              backfaceVisibility: "hidden",
              transformStyle: "preserve-3d",
            }}
          >
            <span className="font-mono text-3xl sm:text-4xl md:text-5xl font-bold text-white/90 tracking-tight">
              {display}
            </span>
          </motion.div>
        </AnimatePresence>
      </div>
      <span className="font-mono text-[9px] sm:text-[10px] text-white/25 uppercase tracking-[0.2em]">
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
    <div className="w-full max-w-sm mx-auto">
      <form onSubmit={handleSubmit}>
        <div
          className="flex items-center gap-2 rounded-full overflow-hidden px-1.5 py-1.5 transition-all"
          style={{
            background: "rgba(255,255,255,0.04)",
            border: "1px solid rgba(255,255,255,0.08)",
          }}
        >
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
            className="flex-1 bg-transparent px-4 py-2.5 font-mono text-sm text-white placeholder-white/20 outline-none disabled:opacity-50 min-w-0"
          />
          <button
            type="submit"
            disabled={status === "loading" || status === "success"}
            className="flex-shrink-0 px-5 py-2.5 bg-white text-black font-mono text-xs font-semibold rounded-full hover:bg-white/90 active:scale-[0.98] transition-all disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
          >
            {status === "loading" ? (
              <span className="inline-block w-3.5 h-3.5 border-2 border-black/20 border-t-black rounded-full animate-spin" />
            ) : status === "success" ? (
              "✓ Done"
            ) : (
              "Join waitlist"
            )}
          </button>
        </div>
      </form>

      <AnimatePresence>
        {message && (
          <motion.p
            initial={{ opacity: 0, y: -6 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0 }}
            className={`mt-3 text-xs font-mono text-center ${
              status === "success" ? "text-emerald-400/70" : "text-red-400/70"
            }`}
          >
            {message}
          </motion.p>
        )}
      </AnimatePresence>

      {count !== null && count > 0 && (
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.3 }}
          className="mt-4 flex justify-center"
        >
          <span className="inline-flex items-center gap-2 px-3 py-1 rounded-full" style={{ background: "rgba(255,255,255,0.03)", border: "1px solid rgba(255,255,255,0.05)" }}>
            <span className="w-1.5 h-1.5 bg-emerald-400/80 rounded-full animate-pulse" />
            <span className="font-mono text-[10px] text-white/30">
              <span className="text-white/60 font-semibold">{count.toLocaleString()}</span> joined
            </span>
          </span>
        </motion.div>
      )}
    </div>
  );
}

// ─── Main ───────────────────────────────────────────────────────────────────

export default function LaunchPage() {
  const [time, setTime] = useState<TimeLeft>(getTimeLeft());

  useEffect(() => {
    const id = setInterval(() => setTime(getTimeLeft()), 1000);
    return () => clearInterval(id);
  }, []);

  useEffect(() => {
    if (time.total <= 0) window.location.href = "/";
  }, [time.total]);

  return (
    <div className="relative w-screen h-screen overflow-hidden" style={{ background: "#000" }}>
      {/* BlackHole — full-screen background */}
      <div className="absolute inset-0 z-0">
        <BlackHole
          particleCount={800}
          particleSize={3}
          colors={["#ffffff", "#aaaaaa", "#666666"]}
          outerRadius={85}
          tilt={25}
          tiltSideway={165}
          trail={45}
          orbitSpeed={3}
          pullSpeed={1}
          showCenter={true}
          centre={{ voidRadius: 35, voidX: 50, voidY: 50 }}
        />
      </div>

      {/* Content overlay */}
      <div className="relative z-10 flex flex-col items-center justify-center h-full px-6">
        {/* Logo */}
        <motion.div
          initial={{ opacity: 0, scale: 0.8 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{ duration: 1.2, delay: 0.2, ease: [0.16, 1, 0.3, 1] }}
          className="mb-8"
        >
          <div
            className="w-14 h-14 rounded-2xl flex items-center justify-center"
            style={{
              background: "rgba(255,255,255,0.05)",
              border: "1px solid rgba(255,255,255,0.1)",
              backdropFilter: "blur(8px)",
            }}
          >
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
          initial={{ opacity: 0, y: 30 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 1, delay: 0.4, ease: [0.16, 1, 0.3, 1] }}
          className="font-mono font-bold text-white text-4xl sm:text-5xl md:text-6xl tracking-tight mb-3 text-center"
        >
          moodwave
        </motion.h1>

        <motion.p
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 1, delay: 0.55, ease: [0.16, 1, 0.3, 1] }}
          className="text-white/35 text-sm sm:text-base text-center max-w-xs leading-relaxed mb-12"
        >
          Your terminal&apos;s new soundtrack.
        </motion.p>

        {/* Countdown */}
        <motion.div
          initial={{ opacity: 0, y: 30 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 1, delay: 0.7, ease: [0.16, 1, 0.3, 1] }}
          className="flex items-start gap-2 sm:gap-3 md:gap-4 mb-4"
        >
          <CountdownUnit value={time.days} label="Days" />
          <span className="text-white/10 text-2xl sm:text-3xl font-light mt-5 sm:mt-7">:</span>
          <CountdownUnit value={time.hours} label="Hours" />
          <span className="text-white/10 text-2xl sm:text-3xl font-light mt-5 sm:mt-7">:</span>
          <CountdownUnit value={time.minutes} label="Min" />
          <span className="text-white/10 text-2xl sm:text-3xl font-light mt-5 sm:mt-7">:</span>
          <CountdownUnit value={time.seconds} label="Sec" />
        </motion.div>

        <motion.p
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 1 }}
          className="font-mono text-[10px] text-white/15 uppercase tracking-[0.25em] mb-12"
        >
          Launching September 3, 2026
        </motion.p>

        {/* Waitlist */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 1, delay: 1.1, ease: [0.16, 1, 0.3, 1] }}
          className="w-full max-w-md"
        >
          <WaitlistForm />
        </motion.div>
      </div>

      {/* Bottom subtle badge */}
      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 2 }}
        className="absolute bottom-5 left-1/2 -translate-x-1/2 z-10"
      >
        <span className="font-mono text-[9px] text-white/10">© 2026 moodwave</span>
      </motion.div>
    </div>
  );
}
