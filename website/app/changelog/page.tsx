"use client";

import { useEffect, useState, useRef } from "react";
import { motion, AnimatePresence } from "framer-motion";
import Nav from "../../components/Nav";
import { 
  Flame, 
  Terminal, 
  Sliders, 
  GitBranch, 
  Cpu, 
  Volume2, 
  Activity, 
  TrendingDown, 
  FolderPlus, 
  Palette,
  ExternalLink,
  ChevronDown,
  ChevronUp
} from "lucide-react";

// Release metadata with custom tracking fields
interface ReleaseItem {
  version: string;
  date: string;
  title: string;
  summary: string;
  githubUrl: string;
  features: string[];
  fixes: string[];
  performance: string[];
  metrics: {
    binarySize: string;
    scanLatency: string;
    themes: number;
    visualizers: number;
  };
}

const staticReleases: ReleaseItem[] = [
  {
    version: "v1.0.4",
    date: "July 09, 2026",
    title: "The Fireplace Update & PresentationCore Migration",
    summary: "Introducing a beautiful cozy fireplace visualizer and migrating Windows audio backend to PresentationCore for native modern format streaming.",
    githubUrl: "https://github.com/Boredooms/Moodwave-CLI/releases/tag/v1.0.4",
    features: [
      "Cozy Fireplace Visualizer (fireplace / fire): A cellular automaton fire simulation that propagates flickering flames upwards.",
      "Windows PresentationCore MediaPlayer Backend: Native support for streaming AAC, MP3, and modern audio formats out of the box.",
      "STA Single-Threaded Apartment: Headless PowerShell subprocesses run in STA mode for safe media pipeline initialization."
    ],
    fixes: [
      "Fixed TUI line wrapping: Restricted visualizer line widths and padded lines to prevent terminal layout breaking.",
      "Standard basic 16-color ANSI overrides: Fire colors now render correctly on legacy terminal emulators without falling back to white."
    ],
    performance: [
      "Audio pipeline optimizations: Windows Media Foundation decoder streams directly to the speakers with zero local storage overhead."
    ],
    metrics: {
      binarySize: "8.0 MB",
      scanLatency: "0.3 ms",
      themes: 9,
      visualizers: 6
    }
  },
  {
    version: "v1.0.2",
    date: "July 07, 2026",
    title: "Raw Mode Fixes & Dynamic Visual Themes",
    summary: "Resolving keyboard listener race conditions, duplicate stop panics, and adding full visual theme colorization to the welcome screen TUI.",
    githubUrl: "https://github.com/Boredooms/Moodwave-CLI/releases/tag/v1.0.2",
    features: [
      "Interactive Theme Customizer Menu: Supports switching color presets in welcome screen, waveforms, and visualizer equalizers.",
      "Expanded 9-Theme Palette: Monochrome, Dark, Ash, Ghost, Ocean, Neon, Sunset, Matrix, Lavender presets.",
      "Graceful Panic Recovery: Top-level error recovery catches runtime crashes and prints clean notifications rather than stack dumps."
    ],
    fixes: [
      "Welcome keyboard race fixed: Safe listener goroutine teardown and TTY restore prior to launching subcommands.",
      "Prevented Stop channel panic: Renderer Stop checks if channels are already closed before closing them."
    ],
    performance: [
      "Ignored directory filters: Added case-insensitive ignores for AppData, OneDrive, Windows, and Program Files, speeding up home directory scans by 99%."
    ],
    metrics: {
      binarySize: "8.1 MB",
      scanLatency: "0.4 ms",
      themes: 9,
      visualizers: 5
    }
  },
  {
    version: "v1.0.0",
    date: "July 01, 2026",
    title: "Initial Launch — Inferring Codebase Mood",
    summary: "The initial launch of Moodwave: scanning repository codebases, inferring developer moods, and playing matching music directly in the terminal.",
    githubUrl: "https://github.com/Boredooms/Moodwave-CLI/releases/tag/v1.0.0",
    features: [
      "Weighted Mood Profiling: Heuristics engine mapping TODOs, file count, languages, and git changes to 10 developer mood profiles.",
      "CLI-first Audio Controller: Backend subprocess manager wrapper supporting mpv, ffplay, and afplay.",
      "ASCII Equalizers: Waveform, Spectrum, and Pulse visualizers."
    ],
    fixes: [
      "Initial production release with clean codebase boundaries."
    ],
    performance: [
      "Concurrent directory scanner with lock-free file categorization."
    ],
    metrics: {
      binarySize: "8.9 MB",
      scanLatency: "1500 ms",
      themes: 3,
      visualizers: 3
    }
  }
];

export default function Changelog() {
  const [activeMetricTab, setActiveMetricTab] = useState<"size" | "speed" | "themes">("speed");
  const [expandedCards, setExpandedCards] = useState<Record<string, boolean>>({
    "v1.0.4": true,
    "v1.0.2": true,
    "v1.0.0": false,
  });

  const toggleExpand = (ver: string) => {
    setExpandedCards(prev => ({ ...prev, [ver]: !prev[ver] }));
  };

  return (
    <div style={{ background: "#080808", minHeight: "100vh", color: "#ffffff", paddingBottom: "100px" }}>
      <Nav version="v1.0.4" />

      {/* Hero Header */}
      <section className="relative pt-32 pb-16 overflow-hidden border-b border-white/[0.05]">
        <div className="absolute inset-0 bg-[radial-gradient(ellipse_at_top,_var(--tw-gradient-stops))] from-white/[0.03] via-transparent to-transparent pointer-events-none" />
        <div className="container-page text-center">
          <p className="font-mono text-xs text-[#555] uppercase tracking-[0.2em] mb-4">
            Version History & Timeline
          </p>
          <h1 className="font-mono font-semibold text-white tracking-tight leading-tight mb-5" style={{ fontSize: "clamp(2rem, 5vw, 3.5rem)" }}>
            Changelog
          </h1>
          <p className="text-[#666] max-w-xl mx-auto text-sm md:text-base leading-relaxed">
            Follow the journey of Moodwave CLI as it evolves from a lightweight mood audio scanner to a highly optimized terminal companion.
          </p>
        </div>
      </section>

      {/* Codebase Evolution Metrics Dashboard */}
      <section className="py-12 border-b border-white/[0.05] bg-white/[0.01]">
        <div className="container-page">
          <div className="border border-white/[0.06] rounded-xl bg-white/[0.01] p-6 md:p-8">
            <div className="flex flex-col md:flex-row items-start md:items-center justify-between gap-6 mb-8">
              <div>
                <span className="font-mono text-xs text-[#444] uppercase tracking-widest block mb-1">Interactive Dashboard</span>
                <h3 className="font-mono font-medium text-white text-lg">Codebase Evolution Over Time</h3>
              </div>
              <div className="flex bg-white/[0.03] p-1 rounded-lg border border-white/[0.05] w-full md:w-auto">
                {(["speed", "size", "themes"] as const).map(tab => (
                  <button
                    key={tab}
                    onClick={() => setActiveMetricTab(tab)}
                    className={`font-mono text-xs px-4 py-2 rounded-md transition-colors cursor-pointer w-full md:w-auto text-center ${
                      activeMetricTab === tab ? "bg-white/[0.08] text-white" : "text-[#555] hover:text-[#888]"
                    }`}
                  >
                    {tab === "speed" ? "Scan Speed" : tab === "size" ? "Binary Size" : "Themes count"}
                  </button>
                ))}
              </div>
            </div>

            {/* Interactive SVG Graphs */}
            <div className="h-[240px] w-full flex items-end relative border-b border-white/[0.05] pb-2">
              <AnimatePresence mode="wait">
                {activeMetricTab === "speed" && (
                  <motion.div
                    key="speed-chart"
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    exit={{ opacity: 0 }}
                    transition={{ duration: 0.3 }}
                    className="absolute inset-0 flex items-end"
                  >
                    <svg className="w-full h-full" viewBox="0 0 800 200" preserveAspectRatio="none">
                      <defs>
                        <linearGradient id="speedGrad" x1="0" y1="0" x2="0" y2="1">
                          <stop offset="0%" stopColor="#ef4444" stopOpacity="0.2" />
                          <stop offset="100%" stopColor="#ef4444" stopOpacity="0" />
                        </linearGradient>
                      </defs>
                      {/* Area */}
                      <path 
                        d="M 50 160 L 400 30 L 750 10 L 750 200 L 50 200 Z" 
                        fill="url(#speedGrad)" 
                        className="transition-all duration-700 ease-in-out"
                      />
                      {/* Line */}
                      <path 
                        d="M 50 160 L 400 30 L 750 10" 
                        fill="none" 
                        stroke="#ef4444" 
                        strokeWidth="2" 
                        strokeDasharray="4 4"
                      />
                      {/* Nodes */}
                      <circle cx="50" cy="160" r="5" fill="#ef4444" />
                      <circle cx="400" cy="30" r="5" fill="#ef4444" />
                      <circle cx="750" cy="10" r="5" fill="#ef4444" />
                      {/* Text */}
                      <text x="50" y="140" fill="#666" fontSize="10" fontFamily="monospace">v1.0.0: 1500ms</text>
                      <text x="370" y="50" fill="#666" fontSize="10" fontFamily="monospace">v1.0.2: 0.4ms</text>
                      <text x="690" y="30" fill="#ef4444" fontSize="10" fontFamily="monospace" fontWeight="bold">v1.0.4: 0.3ms</text>
                    </svg>
                  </motion.div>
                )}

                {activeMetricTab === "size" && (
                  <motion.div
                    key="size-chart"
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    exit={{ opacity: 0 }}
                    transition={{ duration: 0.3 }}
                    className="absolute inset-0 flex items-end"
                  >
                    <svg className="w-full h-full" viewBox="0 0 800 200" preserveAspectRatio="none">
                      <defs>
                        <linearGradient id="sizeGrad" x1="0" y1="0" x2="0" y2="1">
                          <stop offset="0%" stopColor="#3b82f6" stopOpacity="0.2" />
                          <stop offset="100%" stopColor="#3b82f6" stopOpacity="0" />
                        </linearGradient>
                      </defs>
                      <path 
                        d="M 50 30 L 400 110 L 750 130 L 750 200 L 50 200 Z" 
                        fill="url(#sizeGrad)"
                      />
                      <path 
                        d="M 50 30 L 400 110 L 750 130" 
                        fill="none" 
                        stroke="#3b82f6" 
                        strokeWidth="2"
                      />
                      <circle cx="50" cy="30" r="5" fill="#3b82f6" />
                      <circle cx="400" cy="110" r="5" fill="#3b82f6" />
                      <circle cx="750" cy="130" r="5" fill="#3b82f6" />
                      <text x="50" y="50" fill="#666" fontSize="10" fontFamily="monospace">v1.0.0: 8.9MB</text>
                      <text x="370" y="90" fill="#666" fontSize="10" fontFamily="monospace">v1.0.2: 8.1MB</text>
                      <text x="690" y="150" fill="#3b82f6" fontSize="10" fontFamily="monospace" fontWeight="bold">v1.0.4: 8.0MB</text>
                    </svg>
                  </motion.div>
                )}

                {activeMetricTab === "themes" && (
                  <motion.div
                    key="themes-chart"
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    exit={{ opacity: 0 }}
                    transition={{ duration: 0.3 }}
                    className="absolute inset-0 flex items-end"
                  >
                    <svg className="w-full h-full" viewBox="0 0 800 200" preserveAspectRatio="none">
                      <defs>
                        <linearGradient id="themesGrad" x1="0" y1="0" x2="0" y2="1">
                          <stop offset="0%" stopColor="#10b981" stopOpacity="0.2" />
                          <stop offset="100%" stopColor="#10b981" stopOpacity="0" />
                        </linearGradient>
                      </defs>
                      <path 
                        d="M 50 150 L 400 40 L 750 40 L 750 200 L 50 200 Z" 
                        fill="url(#themesGrad)"
                      />
                      <path 
                        d="M 50 150 L 400 40 L 750 40" 
                        fill="none" 
                        stroke="#10b981" 
                        strokeWidth="2"
                      />
                      <circle cx="50" cy="150" r="5" fill="#10b981" />
                      <circle cx="400" cy="40" r="5" fill="#10b981" />
                      <circle cx="750" cy="40" r="5" fill="#10b981" />
                      <text x="50" y="130" fill="#666" fontSize="10" fontFamily="monospace">v1.0.0: 3 themes</text>
                      <text x="360" y="60" fill="#666" fontSize="10" fontFamily="monospace">v1.0.2: 9 themes</text>
                      <text x="690" y="60" fill="#10b981" fontSize="10" fontFamily="monospace" fontWeight="bold">v1.0.4: 9 themes</text>
                    </svg>
                  </motion.div>
                )}
              </AnimatePresence>
            </div>
            
            {/* Visual indicators */}
            <div className="grid grid-cols-3 gap-4 mt-6 text-center">
              <div className="border-r border-white/[0.05] last:border-r-0">
                <span className="text-[10px] font-mono text-[#555] uppercase tracking-widest block mb-1">Code Scan Speed</span>
                <span className="font-mono text-sm md:text-base font-semibold text-white">99.98% Optimized</span>
              </div>
              <div className="border-r border-white/[0.05] last:border-r-0">
                <span className="text-[10px] font-mono text-[#555] uppercase tracking-widest block mb-1">Binary Footprint</span>
                <span className="font-mono text-sm md:text-base font-semibold text-white">8.0MB Header-less</span>
              </div>
              <div>
                <span className="text-[10px] font-mono text-[#555] uppercase tracking-widest block mb-1">Total presets</span>
                <span className="font-mono text-sm md:text-base font-semibold text-white">15 combined TUI presets</span>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Main Release Log Timeline */}
      <section className="py-20">
        <div className="container-page max-w-4xl">
          <div className="relative border-l border-white/[0.08] ml-4 md:ml-8 pl-8 md:pl-12 space-y-16">
            
            {staticReleases.map((release, index) => {
              const isOpen = expandedCards[release.version];

              return (
                <div key={release.version} className="relative">
                  {/* Timeline point */}
                  <div className="absolute -left-[45px] md:-left-[61px] top-1.5 flex items-center justify-center">
                    <div className="w-5 h-5 rounded-full bg-[#080808] border-2 border-white/20 flex items-center justify-center group-hover:border-white transition-colors duration-200">
                      <div className="w-2.5 h-2.5 rounded-full bg-white/50" />
                    </div>
                  </div>

                  {/* Release Card */}
                  <motion.div 
                    initial={{ opacity: 0, y: 30 }}
                    whileInView={{ opacity: 1, y: 0 }}
                    viewport={{ once: true, margin: "-100px" }}
                    transition={{ duration: 0.5, delay: index * 0.1 }}
                    className="border border-white/[0.06] rounded-xl overflow-hidden bg-white/[0.01] hover:border-white/[0.12] transition-colors duration-300"
                  >
                    {/* Header */}
                    <div 
                      onClick={() => toggleExpand(release.version)}
                      className="p-6 md:p-8 flex items-start justify-between gap-4 cursor-pointer select-none"
                    >
                      <div className="space-y-2">
                        <div className="flex items-center gap-3">
                          <span className="font-mono text-xs font-semibold px-2.5 py-1 bg-white/[0.06] border border-white/[0.08] rounded-md text-white">
                            {release.version}
                          </span>
                          <span className="font-mono text-xs text-[#555]">{release.date}</span>
                        </div>
                        <h2 className="font-mono font-medium text-white text-xl leading-tight">
                          {release.title}
                        </h2>
                        <p className="text-sm text-[#777] leading-relaxed max-w-2xl">
                          {release.summary}
                        </p>
                      </div>
                      <div className="text-[#555] hover:text-white p-1 transition-colors">
                        {isOpen ? <ChevronUp className="w-5 h-5" /> : <ChevronDown className="w-5 h-5" />}
                      </div>
                    </div>

                    {/* Collapsible Content */}
                    <AnimatePresence initial={false}>
                      {isOpen && (
                        <motion.div
                          initial={{ height: 0, opacity: 0 }}
                          animate={{ height: "auto", opacity: 1 }}
                          exit={{ height: 0, opacity: 0 }}
                          transition={{ duration: 0.3 }}
                          className="overflow-hidden border-t border-white/[0.06] bg-white/[0.005]"
                        >
                          <div className="p-6 md:p-8 space-y-8">
                            {/* Embedded Interactive Visualization */}
                            {release.version === "v1.0.4" && <FireplaceSimulator />}
                            {release.version === "v1.0.2" && <ConsoleMenuSimulator />}
                            {release.version === "v1.0.0" && <AudioEqualizerSimulator />}

                            {/* Detailed Updates */}
                            <div className="grid md:grid-cols-2 gap-8 pt-4">
                              <div className="space-y-4">
                                <div className="flex items-center gap-2 text-white">
                                  <Sliders className="w-4 h-4 text-emerald-400" />
                                  <span className="font-mono text-xs font-semibold tracking-wider uppercase">Features</span>
                                </div>
                                <ul className="space-y-2.5 pl-1.5">
                                  {release.features.map(f => (
                                    <li key={f} className="text-xs text-[#666] leading-relaxed flex items-start gap-2.5">
                                      <span className="w-1 h-1 bg-emerald-400/60 rounded-full flex-shrink-0 mt-2" />
                                      <span>{f}</span>
                                    </li>
                                  ))}
                                </ul>
                              </div>

                              <div className="space-y-6">
                                <div className="space-y-4">
                                  <div className="flex items-center gap-2 text-white">
                                    <GitBranch className="w-4 h-4 text-rose-400" />
                                    <span className="font-mono text-xs font-semibold tracking-wider uppercase">Fixes</span>
                                  </div>
                                  <ul className="space-y-2.5 pl-1.5">
                                    {release.fixes.map(f => (
                                      <li key={f} className="text-xs text-[#666] leading-relaxed flex items-start gap-2.5">
                                        <span className="w-1 h-1 bg-rose-400/60 rounded-full flex-shrink-0 mt-2" />
                                        <span>{f}</span>
                                      </li>
                                    ))}
                                  </ul>
                                </div>

                                <div className="space-y-4">
                                  <div className="flex items-center gap-2 text-white">
                                    <Cpu className="w-4 h-4 text-blue-400" />
                                    <span className="font-mono text-xs font-semibold tracking-wider uppercase">Performance</span>
                                  </div>
                                  <ul className="space-y-2.5 pl-1.5">
                                    {release.performance.map(p => (
                                      <li key={p} className="text-xs text-[#666] leading-relaxed flex items-start gap-2.5">
                                        <span className="w-1 h-1 bg-blue-400/60 rounded-full flex-shrink-0 mt-2" />
                                        <span>{p}</span>
                                      </li>
                                    ))}
                                  </ul>
                                </div>
                              </div>
                            </div>

                            {/* Footer links */}
                            <div className="flex items-center justify-between border-t border-white/[0.05] pt-6 text-xs font-mono">
                              <div className="flex items-center gap-4 text-[#555]">
                                <span>Binary: <strong className="text-[#888]">{release.metrics.binarySize}</strong></span>
                                <span>Scan time: <strong className="text-[#888]">{release.metrics.scanLatency}</strong></span>
                              </div>
                              <a 
                                href={release.githubUrl} 
                                target="_blank" 
                                rel="noopener noreferrer" 
                                className="flex items-center gap-1.5 text-white/50 hover:text-white transition-colors duration-200"
                              >
                                View Release GitHub <ExternalLink className="w-3.5 h-3.5" />
                              </a>
                            </div>
                          </div>
                        </motion.div>
                      )}
                    </AnimatePresence>
                  </motion.div>
                </div>
              );
            })}

          </div>
        </div>
      </section>
    </div>
  );
}

// ──────────────────────────────────────────────────────────────────────────────
// v1.0.4 Interactive Fireplace Simulator Component
// ──────────────────────────────────────────────────────────────────────────────
function FireplaceSimulator() {
  const [grid, setGrid] = useState<string[]>([]);
  const width = 64;
  const height = 7;
  const timer = useRef<NodeJS.Timeout | null>(null);

  useEffect(() => {
    // Standard basic 16-color ANSI characters palette
    const chars = [" ", ".", ",", "*", "x", "s", "o", "d", "m", "0", "H", "M", "W", "█"];
    const colors = [
      "",               // empty
      "text-gray-800",  // .
      "text-gray-700",  // ,
      "text-red-900",   // *
      "text-red-800",   // x
      "text-red-700",   // s
      "text-red-500",   // o
      "text-red-400",   // d
      "text-orange-500",// m
      "text-orange-400",// 0
      "text-yellow-500",// H
      "text-yellow-400",// M
      "text-white/80",  // W
      "text-white font-bold" // █
    ];

    let fireGrid = Array(height).fill(0).map(() => Array(width).fill(0));

    const step = () => {
      // Seed bottom row with Gaussian curves
      for (let col = 0; col < width; col++) {
        const t = col / width;
        let intensity = 0;
        intensity += Math.exp(-Math.pow((t - 0.5) / 0.12, 2)) * 1.0;     // Center peak
        intensity += Math.exp(-Math.pow((t - 0.22) / 0.07, 2)) * 0.65;   // Left peak
        intensity += Math.exp(-Math.pow((t - 0.78) / 0.08, 2)) * 0.75;   // Right peak

        if (intensity > 1.0) intensity = 1.0;
        if (t < 0.08 || t > 0.92) intensity = 0;

        let val = Math.floor(intensity * (chars.length - 1));
        if (val > 0) {
          val -= Math.floor(Math.random() * 3);
          if (val < 0) val = 0;
        }
        fireGrid[0][col] = val;
      }

      // Propagate flame upwards
      for (let y = 1; y < height; y++) {
        for (let x = 0; x < width; x++) {
          const wind = Math.floor(Math.random() * 3) - 1;
          const srcX = (x + wind + width) % width;
          let decay = Math.floor(Math.random() * 2);
          if (Math.random() * 100 < 38) {
            decay++;
          }
          let val = fireGrid[y - 1][srcX] - decay;
          if (val < 0) val = 0;
          fireGrid[y][x] = val;
        }
      }

      // Render cells to styled HTML elements
      const lines: string[] = [];
      for (let y = height - 1; y >= 0; y--) {
        let lineHtml = "";
        for (let x = 0; x < width; x++) {
          const heat = fireGrid[y][x];
          const char = chars[heat];
          const colorClass = colors[heat];
          if (heat === 0) {
            lineHtml += " ";
          } else {
            lineHtml += `<span class="${colorClass}">${char}</span>`;
          }
        }
        lines.push(lineHtml);
      }
      setGrid(lines);
    };

    timer.current = setInterval(step, 120);
    return () => {
      if (timer.current) clearInterval(timer.current);
    };
  }, []);

  return (
    <div className="space-y-3">
      <div className="flex items-center gap-2 text-xs font-mono text-[#555]">
        <Flame className="w-3.5 h-3.5 text-red-500 animate-pulse" />
        <span>Cozy TUI Fireplace Visualizer Simulation (v1.0.4)</span>
      </div>
      <div className="terminal-frame">
        <div className="terminal-titlebar">
          <div className="terminal-dot bg-rose-500/80" />
          <div className="terminal-dot bg-amber-500/80" />
          <div className="terminal-dot bg-emerald-500/80" />
          <span className="font-mono text-xs text-[#555] ml-2">moodwave - visuals fireplace</span>
        </div>
        <div className="terminal-body font-mono select-none" style={{ fontSize: "11px", lineHeight: "1.35", letterSpacing: "1px", background: "#060606" }}>
          {grid.map((line, idx) => (
            <div key={idx} dangerouslySetInnerHTML={{ __html: line || "&nbsp;" }} />
          ))}
        </div>
      </div>
    </div>
  );
}

// ──────────────────────────────────────────────────────────────────────────────
// v1.0.2 Interactive Console Menu Simulator Component
// ──────────────────────────────────────────────────────────────────────────────
const themesList = [
  { id: "monochrome", primary: "text-white border-white", text: "text-white/80", accent: "bg-white/[0.08] text-white" },
  { id: "dark", primary: "text-gray-300 border-gray-500", text: "text-gray-400", accent: "bg-gray-800 text-gray-200" },
  { id: "ash", primary: "text-neutral-400 border-neutral-600", text: "text-neutral-500", accent: "bg-neutral-800 text-neutral-300" },
  { id: "ghost", primary: "text-zinc-600 border-zinc-800", text: "text-zinc-500", accent: "bg-zinc-900 text-zinc-400" },
  { id: "ocean", primary: "text-cyan-400 border-cyan-700", text: "text-cyan-200/80", accent: "bg-cyan-950/50 text-cyan-300 border-cyan-800" },
  { id: "neon", primary: "text-pink-500 border-pink-700", text: "text-pink-300", accent: "bg-pink-950/40 text-pink-300 border-pink-900" },
  { id: "sunset", primary: "text-orange-400 border-orange-700", text: "text-orange-200", accent: "bg-orange-950/40 text-orange-300 border-orange-900" },
  { id: "matrix", primary: "text-green-500 border-green-700", text: "text-green-300", accent: "bg-green-950/50 text-green-400 border-green-800" },
  { id: "lavender", primary: "text-purple-400 border-purple-700", text: "text-purple-200", accent: "bg-purple-950/40 text-purple-300 border-purple-900" }
];

function ConsoleMenuSimulator() {
  const [activeTheme, setActiveTheme] = useState(4); // ocean default
  const theme = themesList[activeTheme];

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between text-xs font-mono">
        <div className="flex items-center gap-2 text-[#555]">
          <Palette className="w-3.5 h-3.5 text-cyan-400" />
          <span>Interactive TUI Theme Customizer (v1.0.2)</span>
        </div>
        <span className="text-[#444] text-[10px]">Click a theme to override visualizer palette</span>
      </div>

      <div className="grid grid-cols-3 md:grid-cols-5 gap-2">
        {themesList.map((t, idx) => (
          <button
            key={t.id}
            onClick={() => setActiveTheme(idx)}
            className={`font-mono text-xs px-3 py-2 rounded-lg border transition-all text-center cursor-pointer ${
              activeTheme === idx 
                ? "bg-white/[0.06] border-white/30 text-white shadow-lg" 
                : "bg-white/[0.01] border-white/[0.03] text-[#555] hover:text-[#888] hover:border-white/[0.08]"
            }`}
          >
            {t.id}
          </button>
        ))}
      </div>

      {/* TUI Welcome Screen Border Frame */}
      <div className="terminal-frame" style={{ background: "#050505" }}>
        <div className="terminal-titlebar">
          <div className="terminal-dot bg-rose-500/80" />
          <div className="terminal-dot bg-amber-500/80" />
          <div className="terminal-dot bg-emerald-500/80" />
          <span className="font-mono text-xs text-[#444] ml-2">moodwave - theme selector preset</span>
        </div>
        <div className="terminal-body font-mono py-10 flex flex-col items-center justify-center select-none">
          {/* Logo */}
          <div className={`text-center font-bold text-xs md:text-sm tracking-widest mb-6 ${theme.primary}`}>
            ◆ MOODWAVE ◆
          </div>

          {/* Menu box */}
          <div className={`border rounded-lg p-6 max-w-xs w-full text-center space-y-4 ${theme.primary}`}>
            <div className={`text-xs font-bold uppercase tracking-wider ${theme.primary}`}>
              Active Theme: {theme.id}
            </div>
            
            <div className="space-y-2">
              <div className={`text-xs p-2 rounded border border-transparent ${theme.accent} font-semibold flex items-center justify-center gap-1.5`}>
                ▶ {theme.id} (selected)
              </div>
              <div className={`text-xs p-2 text-[#444]`}>
                ⬅ Back to Main Menu
              </div>
            </div>
          </div>

          <div className="mt-8 text-center text-[10px] text-[#444]">
            [W/S / Arrow Keys] Navigate  •  [Enter] Select theme
          </div>
        </div>
      </div>
    </div>
  );
}

// ──────────────────────────────────────────────────────────────────────────────
// v1.0.0 Equalizer Visualizer Simulator Component
// ──────────────────────────────────────────────────────────────────────────────
function AudioEqualizerSimulator() {
  const [bars, setBars] = useState<number[]>(Array(18).fill(20));

  useEffect(() => {
    const interval = setInterval(() => {
      setBars(Array(18).fill(0).map(() => Math.floor(Math.random() * 85) + 15));
    }, 150);
    return () => clearInterval(interval);
  }, []);

  return (
    <div className="space-y-3">
      <div className="flex items-center gap-2 text-xs font-mono text-[#555]">
        <Activity className="w-3.5 h-3.5 text-emerald-400" />
        <span>TUI Equalizer Spectrum Simulation (v1.0.0)</span>
      </div>
      <div className="terminal-frame">
        <div className="terminal-titlebar">
          <div className="terminal-dot bg-rose-500/80" />
          <div className="terminal-dot bg-amber-500/80" />
          <div className="terminal-dot bg-emerald-500/80" />
          <span className="font-mono text-xs text-[#555] ml-2">moodwave - visuals spectrum</span>
        </div>
        <div className="terminal-body font-mono flex items-end justify-center gap-1.5 h-[120px] pb-6 bg-[#060606] overflow-hidden">
          {bars.map((h, i) => (
            <div 
              key={i} 
              className="w-2 bg-emerald-400/80 hover:bg-emerald-400 transition-all duration-150 rounded-t-sm"
              style={{ height: `${h}%` }}
            />
          ))}
        </div>
      </div>
    </div>
  );
}
