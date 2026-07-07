"use client";

import { useEffect, useState, useRef, useCallback } from "react";

const TERMINAL_SEQUENCE = [
  { text: "$ moodwave", delay: 50, color: "#d4d4d4", pause: 500 },
  { text: "", color: "#d4d4d4", pause: 0 },
  { text: "  scanning repository...", delay: 25, color: "#666", pause: 300 },
  { text: "  ├─ 142 files detected", delay: 18, color: "#555", pause: 80 },
  { text: "  ├─ languages: Go, TypeScript, Shell", delay: 18, color: "#555", pause: 80 },
  { text: "  ├─ TODOs: 23  •  FIXMEs: 7  •  BUGs: 3", delay: 18, color: "#555", pause: 80 },
  { text: "  └─ git: 4 staged, 12 modified", delay: 18, color: "#555", pause: 400 },
  { text: "", color: "#d4d4d4", pause: 0 },
  { text: "  mood detected →  DEBUGGING", delay: 22, color: "#e5e5e5", pause: 150 },
  { text: "  confidence: 87%", delay: 18, color: "#666", pause: 500 },
  { text: "", color: "#d4d4d4", pause: 0 },
  { text: "  fetching recommendations...", delay: 22, color: "#666", pause: 600 },
  { text: "  ▶  Lo-Fi Hip Hop Radio — Chillhop Music", delay: 18, color: "#bbb", pause: 150 },
  { text: "  source: youtube  •  bpm: 72  •  energy: low", delay: 18, color: "#555", pause: 300 },
  { text: "", color: "#d4d4d4", pause: 0 },
  { text: "  ▓▓▓▓▓▓▓▓▓▓▓▓░░░░░░░░  02:14 / 03:41", delay: 18, color: "#888", pause: 150 },
  { text: "  [space] pause  [n] next  [l] loop  [q] quit", delay: 18, color: "#444", pause: 3500 },
] as const;

function AnimatedWaveBars() {
  const [heights, setHeights] = useState([2, 4, 6, 8, 5, 3, 7, 4, 6, 2]);

  useEffect(() => {
    const interval = setInterval(() => {
      setHeights((prev) =>
        prev.map(() => Math.floor(Math.random() * 10) + 1)
      );
    }, 150);
    return () => clearInterval(interval);
  }, []);

  return (
    <div className="flex items-end gap-[2px] h-4 mr-1">
      {heights.map((h, i) => (
        <div
          key={i}
          className="w-[3px] bg-[#444] rounded-sm transition-all duration-150"
          style={{ height: `${h * 1.5}px` }}
        />
      ))}
    </div>
  );
}

export default function TerminalWindow({ className = "" }: { className?: string }) {
  const [lines, setLines] = useState<{ text: string; color: string }[]>([]);
  const [seqIdx, setSeqIdx] = useState(0);
  const [charIdx, setCharIdx] = useState(0);
  const [phase, setPhase] = useState<"typing" | "pause" | "done">("typing");
  const bodyRef = useRef<HTMLDivElement>(null);
  const pauseTimer = useRef<NodeJS.Timeout | null>(null);
  const charTimer = useRef<NodeJS.Timeout | null>(null);

  // Reset everything cleanly
  const reset = useCallback(() => {
    setLines([]);
    setSeqIdx(0);
    setCharIdx(0);
    setPhase("typing");
  }, []);

  // Advance to next sequence item
  const advance = useCallback(() => {
    setSeqIdx((prev) => {
      const next = prev + 1;
      if (next >= TERMINAL_SEQUENCE.length) {
        // Done — wait then restart
        pauseTimer.current = setTimeout(reset, 3000);
        return prev;
      }
      setCharIdx(0);
      setPhase("typing");
      return next;
    });
  }, [reset]);

  useEffect(() => {
    const seq = TERMINAL_SEQUENCE[seqIdx];
    if (!seq) return;

    if (phase === "typing") {
      if (seq.text === "") {
        // Empty line — just commit it immediately
        setLines((prev) => [...prev, { text: "", color: seq.color }]);
        setPhase("pause");
        return;
      }

      if (charIdx < seq.text.length) {
        charTimer.current = setTimeout(() => {
          setCharIdx((c) => c + 1);
        }, seq.delay);
        return () => { if (charTimer.current) clearTimeout(charTimer.current); };
      } else {
        // Finished typing this line
        setLines((prev) => [...prev, { text: seq.text, color: seq.color }]);
        setPhase("pause");
      }
    }

    if (phase === "pause") {
      pauseTimer.current = setTimeout(() => {
        if (seqIdx >= TERMINAL_SEQUENCE.length - 1) {
          pauseTimer.current = setTimeout(reset, 2000);
        } else {
          advance();
        }
      }, seq.pause || 0);
      return () => { if (pauseTimer.current) clearTimeout(pauseTimer.current); };
    }
  }, [seqIdx, charIdx, phase, advance, reset]);

  // Auto-scroll terminal body
  useEffect(() => {
    if (bodyRef.current) {
      bodyRef.current.scrollTop = bodyRef.current.scrollHeight;
    }
  }, [lines]);

  const currentSeq = TERMINAL_SEQUENCE[seqIdx];
  const isTyping = phase === "typing" && currentSeq && currentSeq.text !== "" && charIdx < currentSeq.text.length;
  const typingText = isTyping ? currentSeq.text.slice(0, charIdx) : null;

  return (
    <div className={`terminal-frame ${className}`}>
      {/* Titlebar */}
      <div className="terminal-titlebar">
        <div className="terminal-dot" style={{ background: "#3a3a3a" }} />
        <div className="terminal-dot" style={{ background: "#3a3a3a" }} />
        <div className="terminal-dot" style={{ background: "#3a3a3a" }} />
        <span className="text-xs text-[#3a3a3a] font-mono ml-3 flex-1 text-center" style={{ letterSpacing: "0.08em" }}>
          moodwave
        </span>
        <AnimatedWaveBars />
      </div>

      {/* Body */}
      <div
        ref={bodyRef}
        className="terminal-body"
        style={{ height: "320px", overflowY: "auto", overflowX: "hidden", scrollbarWidth: "none" }}
      >
        {lines.map((line, i) => (
          <div key={i} style={{ color: line.color, minHeight: line.text === "" ? "1em" : undefined }}>
            {line.text}
          </div>
        ))}

        {/* Currently typing */}
        {typingText !== null && (
          <div style={{ color: currentSeq?.color }}>
            {typingText}
            <span className="cursor" />
          </div>
        )}

        {/* Idle cursor */}
        {typingText === null && phase !== "done" && (
          <div>
            <span className="cursor" />
          </div>
        )}
      </div>
    </div>
  );
}
