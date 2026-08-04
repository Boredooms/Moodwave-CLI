"use client";

import { useEffect, useState, useRef, useCallback } from "react";
import { motion, AnimatePresence } from "framer-motion";

// ─────────────────────────────────────────────────────────────────────────────
// LAUNCH COUNTDOWN PAGE
//
// This page is the ONLY thing visible on the entire website until September 3,
// 2026 at 06:00 AM IST. The server-side middleware (middleware.ts) redirects
// every other route here before any HTML is generated, so there is nothing to
// bypass or inspect. After the date passes, this page is no longer reachable
// (middleware stops redirecting, and nobody links to /launch).
//
// Features:
// - Real-time countdown clock with calendar-page-turn animation on day flip
// - Email waitlist form (submits to /api/waitlist)
// - Fully self-contained — no Nav, no links, no route hints
// - Smooth Lenis-style scroll (page is short enough not to scroll, but the
//   wrapper prevents jank on overscroll)
// - GSAP-style entrance animation via Framer Motion
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

// ─── Calendar Page Turn Digit ───────────────────────────────────────────────

function FlipDigit({ value, label }: { value: number; label: string }) {
  const display = String(value).padStart(2, "0");

  return (
    <div className="flex flex-col items-center gap-2">
      <div className="relative w-[72px] h-[88px] sm:w-[96px] sm:h-[116px] perspective-[600px]">
        <AnimatePresence mode="popLayout">
          <motion.div
            key={display}
            initial={{ rotateX: -90, opacity: 0 }}
            animate={{ rotateX: 0, opacity: 1 }}
            exit={{ rotateX: 90, opacity: 0 }}
            transition={{ duration: 0.5, ease: [0.16, 1, 0.3, 1] }}
            className="absolute inset-0 flex items-center justify-center rounded-xl border border-white/[0.08] bg-[#0c0c0c] shadow-[0_8px_32px_rgba(0,0,0,0.6)]"
            style={{ backfaceVisibility: "hidden", transformStyle: "preserve-3d" }}
          >
            {/* Top fold line */}
            <div className="absolute top-1/2 left-0 right-0 h-[1px] bg-white/[0.04]" />
            <span className="font-mono text-3xl sm:text-5xl font-bold text-white tracking-tight">
              {display}
            </span>
          </motion.div>
        </AnimatePresence>
      </div>
      <span className="font-mono text-[10px] sm:text-xs text-[#555] uppercase tracking-[0.15em]">
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
  const inputRef = useRef<HTMLInputElement>(null);

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
        } else {
          setStatus("error");
          setMessage(data.error || "Something went wrong.");
        }
      } catch {
        setStatus("error");
        setMessage("Network error. Try again.");
      }
    },
    [email, status]
  );

  return (
    <form onSubmit={handleSubmit} className="w-full max-w-md mx-auto">
      <div className="flex gap-2">
        <input
          ref={inputRef}
          type="email"
          value={email}
          onChange={(e) => {
            setEmail(e.target.value);
            if (status === "error") setStatus("idle");
          }}
          placeholder="you@example.com"
          required
          disabled={status === "loading" || status === "success"}
          className="flex-1 bg-white/[0.04] border border-white/[0.1] rounded-lg px-4 py-3 font-mono text-sm text-white placeholder-[#555] outline-none focus:border-white/30 transition-colors disabled:opacity-50"
        />
        <button
          type="submit"
          disabled={status === "loading" || status === "success"}
          className="px-5 py-3 bg-white text-[#080808] font-mono text-sm font-semibold rounded-lg hover:bg-white/90 transition-colors disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer whitespace-nowrap"
        >
          {status === "loading"
            ? "..."
            : status === "success"
            ? "✓"
            : "Join Waitlist"}
        </button>
      </div>
      <AnimatePresence>
        {message && (
          <motion.p
            initial={{ opacity: 0, y: -4 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0 }}
            className={`mt-3 text-xs font-mono text-center ${
              status === "success" ? "text-emerald-400" : "text-red-400"
            }`}
          >
            {message}
          </motion.p>
        )}
      </AnimatePresence>
    </form>
  );
}

// ─── Floating particles ─────────────────────────────────────────────────────

function Particles() {
  const canvasRef = useRef<HTMLCanvasElement>(null);

  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;
    const ctx = canvas.getContext("2d");
    if (!ctx) return;

    let animId: number;
    const particles: { x: number; y: number; vx: number; vy: number; r: number; a: number }[] = [];

    const resize = () => {
      canvas.width = window.innerWidth;
      canvas.height = window.innerHeight;
    };
    resize();
    window.addEventListener("resize", resize);

    for (let i = 0; i < 60; i++) {
      particles.push({
        x: Math.random() * canvas.width,
        y: Math.random() * canvas.height,
        vx: (Math.random() - 0.5) * 0.3,
        vy: (Math.random() - 0.5) * 0.3,
        r: Math.random() * 1.5 + 0.5,
        a: Math.random() * 0.3 + 0.05,
      });
    }

    const draw = () => {
      ctx.clearRect(0, 0, canvas.width, canvas.height);
      for (const p of particles) {
        p.x += p.vx;
        p.y += p.vy;
        if (p.x < 0) p.x = canvas.width;
        if (p.x > canvas.width) p.x = 0;
        if (p.y < 0) p.y = canvas.height;
        if (p.y > canvas.height) p.y = 0;

        ctx.beginPath();
        ctx.arc(p.x, p.y, p.r, 0, Math.PI * 2);
        ctx.fillStyle = `rgba(255,255,255,${p.a})`;
        ctx.fill();
      }
      animId = requestAnimationFrame(draw);
    };
    draw();

    return () => {
      cancelAnimationFrame(animId);
      window.removeEventListener("resize", resize);
    };
  }, []);

  return <canvas ref={canvasRef} className="fixed inset-0 pointer-events-none z-0" />;
}

// ─── Main Page ──────────────────────────────────────────────────────────────

export default function LaunchPage() {
  const [time, setTime] = useState<TimeLeft>(getTimeLeft());

  useEffect(() => {
    const id = setInterval(() => setTime(getTimeLeft()), 1000);
    return () => clearInterval(id);
  }, []);

  // If the countdown is over and somehow this page is still rendered
  // (shouldn't happen due to middleware, but defensive), redirect to home.
  useEffect(() => {
    if (time.total <= 0) {
      window.location.href = "/";
    }
  }, [time.total]);

  return (
    <div className="relative min-h-screen flex flex-col items-center justify-center overflow-hidden px-6" style={{ background: "#050505" }}>
      <Particles />

      {/* Subtle radial glow */}
      <div
        className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[600px] h-[600px] pointer-events-none z-0"
        style={{ background: "radial-gradient(circle, rgba(255,255,255,0.015) 0%, transparent 70%)" }}
      />

      <motion.div
        initial={{ opacity: 0, y: 40 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 1.2, ease: [0.16, 1, 0.3, 1] }}
        className="relative z-10 flex flex-col items-center text-center max-w-2xl"
      >
        {/* Logo mark */}
        <motion.div
          initial={{ scale: 0.8, opacity: 0 }}
          animate={{ scale: 1, opacity: 1 }}
          transition={{ duration: 0.8, delay: 0.2, ease: [0.16, 1, 0.3, 1] }}
          className="mb-8"
        >
          <svg className="w-12 h-12 text-white" viewBox="0 0 32 32" fill="none">
            <circle cx="16" cy="16" r="14" fill="#080808" stroke="currentColor" strokeWidth="1.5" />
            <path
              d="M9 16h1.5M13 11v10M17 7v18M21 13v6M25 16h-1.5"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
            />
          </svg>
        </motion.div>

        {/* Title */}
        <motion.h1
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8, delay: 0.4, ease: [0.16, 1, 0.3, 1] }}
          className="font-mono font-bold text-white text-3xl sm:text-5xl tracking-tight mb-3"
        >
          moodwave
        </motion.h1>

        <motion.p
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8, delay: 0.55, ease: [0.16, 1, 0.3, 1] }}
          className="font-mono text-[#666] text-sm sm:text-base mb-12 max-w-md leading-relaxed"
        >
          A terminal-native music companion that reads your code, detects your mood, and plays the perfect soundtrack.
        </motion.p>

        {/* Countdown */}
        <motion.div
          initial={{ opacity: 0, scale: 0.95 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{ duration: 0.8, delay: 0.7, ease: [0.16, 1, 0.3, 1] }}
          className="flex items-center gap-3 sm:gap-5 mb-4"
        >
          <FlipDigit value={time.days} label="Days" />
          <span className="text-[#333] text-2xl sm:text-4xl font-mono font-light mt-[-20px]">:</span>
          <FlipDigit value={time.hours} label="Hours" />
          <span className="text-[#333] text-2xl sm:text-4xl font-mono font-light mt-[-20px]">:</span>
          <FlipDigit value={time.minutes} label="Min" />
          <span className="text-[#333] text-2xl sm:text-4xl font-mono font-light mt-[-20px]">:</span>
          <FlipDigit value={time.seconds} label="Sec" />
        </motion.div>

        <motion.p
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 1.0 }}
          className="font-mono text-[10px] text-[#444] uppercase tracking-[0.2em] mb-12"
        >
          Launching September 3, 2026
        </motion.p>

        {/* Waitlist */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8, delay: 1.1, ease: [0.16, 1, 0.3, 1] }}
          className="w-full"
        >
          <p className="font-mono text-xs text-[#555] mb-4">
            Get notified on launch day
          </p>
          <WaitlistForm />
        </motion.div>

        {/* Features teaser */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8, delay: 1.4, ease: [0.16, 1, 0.3, 1] }}
          className="mt-16 grid grid-cols-2 sm:grid-cols-4 gap-6 text-center"
        >
          {[
            { icon: "◈", label: "Mood Detection" },
            { icon: "▶", label: "Live Playback" },
            { icon: "◆", label: "24 Visualizers" },
            { icon: "☰", label: "Smart Playlists" },
          ].map((f) => (
            <div key={f.label} className="flex flex-col items-center gap-2">
              <span className="text-white text-lg">{f.icon}</span>
              <span className="font-mono text-[10px] text-[#555] uppercase tracking-wider">
                {f.label}
              </span>
            </div>
          ))}
        </motion.div>
      </motion.div>

      {/* Footer */}
      <motion.p
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 1.8 }}
        className="absolute bottom-6 font-mono text-[10px] text-[#333] z-10"
      >
        © 2026 moodwave
      </motion.p>
    </div>
  );
}
