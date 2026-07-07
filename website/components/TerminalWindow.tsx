"use client";

import { useEffect, useState, useRef } from "react";
import { motion } from "framer-motion";

const TERMINAL_SEQUENCE = [
  { text: "$ moodwave", delay: 0, color: "#d4d4d4", pause: 600 },
  { text: "", delay: 0, color: "#d4d4d4", pause: 0 },
  { text: "  scanning repository...", delay: 30, color: "#666", pause: 400 },
  { text: "  ├─ 142 files detected", delay: 20, color: "#555", pause: 100 },
  { text: "  ├─ languages: Go, TypeScript, Shell", delay: 20, color: "#555", pause: 100 },
  { text: "  ├─ TODOs: 23  •  FIXMEs: 7  •  BUGs: 3", delay: 20, color: "#555", pause: 100 },
  { text: "  └─ git: 4 staged, 12 modified", delay: 20, color: "#555", pause: 500 },
  { text: "", delay: 0, color: "#d4d4d4", pause: 0 },
  { text: "  mood detected →  DEBUGGING", delay: 25, color: "#e5e5e5", pause: 200 },
  { text: "  confidence: 87%", delay: 20, color: "#666", pause: 600 },
  { text: "", delay: 0, color: "#d4d4d4", pause: 0 },
  { text: "  fetching recommendations...", delay: 25, color: "#666", pause: 800 },
  { text: "  ▶  Lo-Fi Hip Hop Radio — Chillhop Music", delay: 20, color: "#bbb", pause: 200 },
  { text: "  source: youtube  •  bpm: 72  •  energy: low", delay: 20, color: "#555", pause: 400 },
  { text: "", delay: 0, color: "#d4d4d4", pause: 0 },
  { text: "  ▓▓▓▓▓▓▓▓▓▓▓▓░░░░░░░░  02:14 / 03:41", delay: 20, color: "#888", pause: 200 },
  { text: "  [space] pause  [n] next  [l] loop  [q] quit", delay: 20, color: "#444", pause: 200 },
];

function WaveBar({ delay }: { delay: number }) {
  const heights = ["h-1", "h-2", "h-3", "h-4", "h-3", "h-2", "h-1", "h-2", "h-3", "h-4", "h-5", "h-4", "h-3"];
  const [idx, setIdx] = useState(0);

  useEffect(() => {
    const t = setTimeout(() => {
      const interval = setInterval(() => {
        setIdx((i) => (i + 1) % heights.length);
      }, 120 + delay);
      return () => clearInterval(interval);
    }, delay);
    return () => clearTimeout(t);
  }, [delay]);

  return (
    <div className={`w-[3px] bg-[#666] rounded-sm transition-all duration-100 ease-linear ${heights[idx]}`} />
  );
}

export default function TerminalWindow({ className = "" }: { className?: string }) {
  const [lines, setLines] = useState<{ text: string; color: string; id: number }[]>([]);
  const [currentLineIdx, setCurrentLineIdx] = useState(0);
  const [currentCharIdx, setCurrentCharIdx] = useState(0);
  const [showCursor, setShowCursor] = useState(true);
  const [isComplete, setIsComplete] = useState(false);
  const endRef = useRef<HTMLDivElement>(null);
  const cycleRef = useRef<NodeJS.Timeout | null>(null);

  const runSequence = () => {
    setLines([]);
    setCurrentLineIdx(0);
    setCurrentCharIdx(0);
    setIsComplete(false);
    setShowCursor(true);
  };

  useEffect(() => {
    runSequence();
  }, []);

  useEffect(() => {
    if (isComplete) {
      cycleRef.current = setTimeout(() => {
        runSequence();
      }, 4000);
      return () => { if (cycleRef.current) clearTimeout(cycleRef.current); };
    }
  }, [isComplete]);

  useEffect(() => {
    if (currentLineIdx >= TERMINAL_SEQUENCE.length) {
      setIsComplete(true);
      return;
    }

    const seq = TERMINAL_SEQUENCE[currentLineIdx];

    if (seq.text === "") {
      setLines((prev) => [...prev, { text: "", color: seq.color, id: Date.now() }]);
      const t = setTimeout(() => {
        setCurrentLineIdx((i) => i + 1);
      }, seq.pause);
      return () => clearTimeout(t);
    }

    if (currentCharIdx < seq.text.length) {
      const t = setTimeout(() => {
        setCurrentCharIdx((c) => c + 1);
      }, seq.delay);
      return () => clearTimeout(t);
    } else {
      const t = setTimeout(() => {
        setLines((prev) => [...prev, { text: seq.text, color: seq.color, id: Date.now() }]);
        setCurrentLineIdx((i) => i + 1);
        setCurrentCharIdx(0);
      }, seq.pause);
      return () => clearTimeout(t);
    }
  }, [currentLineIdx, currentCharIdx]);

  useEffect(() => {
    endRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [lines]);

  const currentSeq = TERMINAL_SEQUENCE[currentLineIdx];
  const typingLine = currentSeq && currentSeq.text !== "" && currentCharIdx < currentSeq.text.length
    ? { text: currentSeq.text.slice(0, currentCharIdx), color: currentSeq.color }
    : null;

  return (
    <div className={`terminal-frame ${className}`}>
      {/* Titlebar */}
      <div className="terminal-titlebar">
        <div className="terminal-dot" style={{ background: "#444" }} />
        <div className="terminal-dot" style={{ background: "#444" }} />
        <div className="terminal-dot" style={{ background: "#444" }} />
        <span
          className="text-xs text-[#444] font-mono ml-3 flex-1 text-center"
          style={{ letterSpacing: "0.08em" }}
        >
          moodwave
        </span>
        {/* mini waveform */}
        <div className="flex items-end gap-[2px] h-4 mr-1">
          {[0, 80, 160, 40, 120, 200, 60, 140, 20, 100].map((d, i) => (
            <WaveBar key={i} delay={d} />
          ))}
        </div>
      </div>

      {/* Body */}
      <div
        className="terminal-body h-[320px] overflow-y-auto overflow-x-hidden"
        style={{ scrollbarWidth: "none" }}
      >
        {lines.map((line) => (
          <div key={line.id} style={{ color: line.color, minHeight: line.text === "" ? "0.85em" : undefined }}>
            {line.text}
          </div>
        ))}

        {/* Currently typing line */}
        {typingLine && (
          <div style={{ color: typingLine.color }}>
            {typingLine.text}
            <span className="cursor" />
          </div>
        )}

        {/* Idle cursor */}
        {!typingLine && !isComplete && (
          <div>
            <span className="cursor" />
          </div>
        )}

        <div ref={endRef} />
      </div>
    </div>
  );
}
