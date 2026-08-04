"use client";

import { useEffect, useState, useCallback } from "react";
import { motion, AnimatePresence } from "framer-motion";
import BlackHole from "@/components/originkit/ui/blackhole";

// ─────────────────────────────────────────────────────────────────────────────
// LAUNCH PAGE — BlackHole visible at the top, content clearly readable below
// ─────────────────────────────────────────────────────────────────────────────

const UNLOCK_UTC = new Date("2026-09-03T00:30:00Z");

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
    <div className="flex flex-col items-center gap-2.5">
      <div className="relative w-[72px] h-[84px] sm:w-[88px] sm:h-[100px] md:w-[100px] md:h-[112px]">
        <AnimatePresence mode="popLayout">
          <motion.div
            key={display}
            initial={{ rotateX: -90, opacity: 0 }}
            animate={{ rotateX: 0, opacity: 1 }}
            exit={{ rotateX: 90, opacity: 0 }}
            transition={{ duration: 0.5, ease: [0.16, 1, 0.3, 1] }}
            className="absolute inset-0 flex items-center justify-center rounded-2xl"
            style={{
              background: "linear-gradient(180deg, #1a1a1a 0%, #111111 100%)",
              border: "1px solid rgba(255,255,255,0.08)",
              boxShadow: "0 12px 40px rgba(0,0,0,0.5), inset 0 1px 0 rgba(255,255,255,0.06)",
              backfaceVisibility: "hidden",
            }}
          >
            <span className="font-mono text-3xl sm:text-4xl md:text-5xl font-bold text-white tracking-tight">
              {display}
            </span>
          </motion.div>
        </AnimatePresence>
      </div>
      <span className="font-mono text-[10px] text-white/30 uppercase tracking-[0.2em]">
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
            background: "rgba(255,255,255,0.06)",
            border: "1px solid rgba(255,255,255,0.12)",
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
            className="flex-1 bg-transparent px-4 py-2.5 font-mono text-sm text-white placeholder-white/30 outline-none disabled:opacity-50 min-w-0"
          />
          <button
            type="submit"
            disabled={status === "loading" || status === "success"}
            className="flex-shrink-0 px-5 py-2.5 bg-white text-black font-mono text-xs font-semibold rounded-full hover:bg-white/90 active:scale-[0.97] transition-all disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
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
              status === "success" ? "text-emerald-400/80" : "text-red-400/80"
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
          <span
            className="inline-flex items-center gap-2 px-3 py-1.5 rounded-full"
            style={{ background: "rgba(255,255,255,0.04)", border: "1px solid rgba(255,255,255,0.07)" }}
          >
            <span className="w-1.5 h-1.5 bg-emerald-400 rounded-full animate-pulse" />
            <span className="font-mono text-[10px] text-white/40">
              <span className="text-white/70 font-semibold">{count.toLocaleString()}</span> joined
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
    <div className="relative w-screen h-screen overflow-hidden flex flex-col" style={{ background: "#000" }}>

      {/* ─── Top: BlackHole hero (takes ~45% of viewport, clearly visible) ─── */}
      <div className="relative flex-shrink-0 w-full" style={{ height: "45vh" }}>
        <BlackHole
          particleCount={600}
          particleSize={3}
          colors={["#ffffff", "#cccccc", "#888888"]}
          outerRadius={80}
          tilt={25}
          tiltSideway={165}
          trail={42}
          orbitSpeed={3}
          pullSpeed={1}
          showCenter={true}
          centre={{ voidRadius: 30, voidX: 50, voidY: 55 }}
        />
        {/* Fade to black at the bottom edge so it blends into content */}
        <div
          className="absolute bottom-0 left-0 right-0 h-24 pointer-events-none"
          style={{ background: "linear-gradient(to bottom, transparent, #000)" }}
        />
      </div>

      {/* ─── Bottom: Content (centered, readable, no overlap with particles) ─── */}
      <div className="relative z-10 flex-1 flex flex-col items-center justify-center px-6 -mt-8">

        {/* Title */}
        <motion.h1
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 1, delay: 0.2, ease: [0.16, 1, 0.3, 1] }}
          className="font-mono font-bold text-white text-3xl sm:text-4xl md:text-5xl tracking-tight mb-2 text-center"
        >
          moodwave
        </motion.h1>

        <motion.p
          initial={{ opacity: 0, y: 15 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 1, delay: 0.35, ease: [0.16, 1, 0.3, 1] }}
          className="text-white/40 text-sm text-center max-w-xs mb-10"
        >
          Your terminal&apos;s new soundtrack.
        </motion.p>

        {/* Countdown */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 1, delay: 0.5, ease: [0.16, 1, 0.3, 1] }}
          className="flex items-start gap-2 sm:gap-3 md:gap-4 mb-3"
        >
          <CountdownUnit value={time.days} label="Days" />
          <span className="text-white/15 text-2xl sm:text-3xl font-light mt-6 sm:mt-8">:</span>
          <CountdownUnit value={time.hours} label="Hours" />
          <span className="text-white/15 text-2xl sm:text-3xl font-light mt-6 sm:mt-8">:</span>
          <CountdownUnit value={time.minutes} label="Min" />
          <span className="text-white/15 text-2xl sm:text-3xl font-light mt-6 sm:mt-8">:</span>
          <CountdownUnit value={time.seconds} label="Sec" />
        </motion.div>

        <motion.p
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.8 }}
          className="font-mono text-[10px] text-white/20 uppercase tracking-[0.2em] mb-10"
        >
          Launching September 3, 2026
        </motion.p>

        {/* Waitlist */}
        <motion.div
          initial={{ opacity: 0, y: 15 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 1, delay: 0.9, ease: [0.16, 1, 0.3, 1] }}
          className="w-full max-w-md"
        >
          <WaitlistForm />
        </motion.div>
      </div>

      {/* Footer */}
      <motion.p
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 1.5 }}
        className="absolute bottom-4 left-1/2 -translate-x-1/2 font-mono text-[9px] text-white/10 z-10"
      >
        © 2026 moodwave
      </motion.p>
    </div>
  );
}
