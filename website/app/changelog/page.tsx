"use client";

import { useEffect, useState, useRef } from "react";
import { motion, AnimatePresence } from "framer-motion";
import Nav from "../../components/Nav";
import SplitText from "../../components/ui/SplitText";
import FadeIn from "../../components/ui/FadeIn";
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
    version: "v1.0.5",
    date: "July 10, 2026",
    title: "Dynamic Changelog Engine & Unshallow CI Actions",
    summary: "Introducing dynamic GitHub release parsing for the website and unshallow repository depth git-log delta builders in the release workflow.",
    githubUrl: "https://github.com/Boredooms/Moodwave-CLI/releases/tag/v1.0.5",
    features: [
      "Dynamic Release Pipeline: Added GitHub API polling and heuristic text parsing to merge live and static releases on `/changelog`.",
      "Unshallow CI Checkout: Set `fetch-depth: 0` in release actions to enable historical git descriptions and commit log analysis.",
      "Auto-parsed Delta Generator: Automates the generation of a clean, commit-level delta to release notes for every version tag."
    ],
    fixes: [
      "Dynamic Nav versioning: Fixed hardcoded navbar versions by reading the latest live version tag from the API state."
    ],
    performance: [
      "Optimized release notes footprint: Generated minimal changelog notes under release descriptions."
    ],
    metrics: {
      binarySize: "8.0 MB",
      scanLatency: "0.3 ms",
      themes: 9,
      visualizers: 6
    }
  },
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
    version: "v1.0.3",
    date: "July 09, 2026",
    title: "Command Block Polish & Layout Refinements",
    summary: "Improving website code block spacing, fixing layout overflows on mobile viewports, and custom branding assets.",
    githubUrl: "https://github.com/Boredooms/Moodwave-CLI/releases/tag/v1.0.3",
    features: [
      "Custom Brand Assets: Embedded a custom SVG favicon, vector logos, and a high-fidelity root repository README architecture chart.",
      "Viewport-safe layout: Redesigned the terminal visualizer emulation grid to flex and scale nicely on mobile viewports."
    ],
    fixes: [
      "CommandBlock scroll fix: Allowed horizontal code scrolling for terminal install commands without stretching containers.",
      "CommandBlock padding adjustments: Shrunk code block sizes to text-xs to achieve perfect visual proportions."
    ],
    performance: [
      "Removed redundant files: Deleted unused build and deploy guides to keep repository footprint clean."
    ],
    metrics: {
      binarySize: "8.1 MB",
      scanLatency: "0.4 ms",
      themes: 9,
      visualizers: 5
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
    version: "v1.0.1",
    date: "July 07, 2026",
    title: "Dynamic Version Tracking & Platform Targets",
    summary: "Implementing dynamic version tag fetching on the web app and renaming platform build names to resolve install script targets.",
    githubUrl: "https://github.com/Boredooms/Moodwave-CLI/releases/tag/v1.0.1",
    features: [
      "Live CLI Release Tracking: Configured the homepage to query live GitHub tags to state CLI version immediately.",
      "Smooth Scrolling & Transitions: Integrated Lenis smooth scroll and GSAP scroll-triggers for interactive section reveals."
    ],
    fixes: [
      "Platform name mismatch: Renamed the armv7 build target to arm to align with the installation script expectations.",
      "CLI Update command: Fixed the self-upgrade command dispatcher routing path."
    ],
    performance: [
      "Go module resolution: Patched x/sys and x/term imports to allow compatibility back to Go 1.22.0."
    ],
    metrics: {
      binarySize: "8.8 MB",
      scanLatency: "18 ms",
      themes: 3,
      visualizers: 3
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


// Helper to generate a smooth cubic bezier path for an array of points
function getBezierPath(pts: { x: number; y: number }[]) {
  if (pts.length === 0) return "";
  let d = `M ${pts[0].x} ${pts[0].y}`;
  for (let i = 0; i < pts.length - 1; i++) {
    const curr = pts[i];
    const next = pts[i + 1];
    const cp1x = curr.x + (next.x - curr.x) / 3;
    const cp1y = curr.y;
    const cp2x = curr.x + 2 * (next.x - curr.x) / 3;
    const cp2y = next.y;
    d += ` C ${cp1x} ${cp1y}, ${cp2x} ${cp2y}, ${next.x} ${next.y}`;
  }
  return d;
}

export default function Changelog() {
  const [releases, setReleases] = useState<ReleaseItem[]>(staticReleases);
  const [activeMetricTab, setActiveMetricTab] = useState<"size" | "speed" | "themes">("speed");
  const [hoveredIdx, setHoveredIdx] = useState<number | null>(null);
  const [expandedCards, setExpandedCards] = useState<Record<string, boolean>>({
    "v1.0.5": true,
    "v1.0.4": true,
    "v1.0.2": true,
    "v1.0.0": false,
  });

  useEffect(() => {
    async function fetchReleases() {
      try {
        const res = await fetch("https://api.github.com/repos/Boredooms/Moodwave-CLI/releases", {
          headers: {
            "User-Agent": "Moodwave-Website-Builder"
          }
        });
        if (!res.ok) return;
        const data = await res.json();
        if (Array.isArray(data)) {
          // Merge API releases with static template metadata
          const merged: ReleaseItem[] = data.map((gitRelease: any) => {
            const version = gitRelease.tag_name;
            const staticMatch = staticReleases.find(r => r.version === version);
            
            const dateObj = new Date(gitRelease.published_at);
            const formattedDate = dateObj.toLocaleDateString("en-US", {
              month: "long",
              day: "2-digit",
              year: "numeric"
            });

            if (staticMatch) {
              return {
                ...staticMatch,
                date: formattedDate,
                githubUrl: gitRelease.html_url
              };
            }

            // Parse markdown release descriptions into categories
            let bodyText = gitRelease.body || "";
            // Extract text ONLY after ## Changelog header if present, to skip installation blocks
            const changelogIdx = bodyText.indexOf("## Changelog");
            if (changelogIdx !== -1) {
              bodyText = bodyText.slice(changelogIdx + 12);
            }

            const bodyLines = bodyText.split("\n");
            const features: string[] = [];
            const fixes: string[] = [];
            const performance: string[] = [];
            
            let currentCat = features;
            for (let line of bodyLines) {
              line = line.trim();
              if (!line) continue;

              // Filter out installation keywords to keep visual clean
              if (line.includes("irm https") || line.includes("curl -") || line.includes("iex") || line.includes("moodwave doctor") || line.includes("sha256sum") || line.includes("| Platform")) {
                continue;
              }

              // Switch categories based on headings or text clues
              const lowerLine = line.toLowerCase();
              if (lowerLine.includes("### 🚀 features") || lowerLine.includes("features")) {
                currentCat = features;
                continue;
              } else if (lowerLine.includes("### 🐛 bug fixes") || lowerLine.includes("fix") || lowerLine.includes("bug")) {
                currentCat = fixes;
                continue;
              } else if (lowerLine.includes("### ⚙️") || lowerLine.includes("perf") || lowerLine.includes("speed") || lowerLine.includes("chore") || lowerLine.includes("ci")) {
                currentCat = performance;
                continue;
              }
              
              if (line.startsWith("-") || line.startsWith("*")) {
                const bullet = line.slice(1).trim();
                if (bullet && !bullet.startsWith("##") && !bullet.startsWith("###")) currentCat.push(bullet);
              } else if (line.match(/^\d+\./)) {
                const bullet = line.replace(/^\d+\./, "").trim();
                if (bullet && !bullet.startsWith("##") && !bullet.startsWith("###")) currentCat.push(bullet);
              } else if (line.startsWith("###") || line.startsWith("##")) {
                // Skip headers
                continue;
              } else {
                // If it is just a plain commit entry or line
                if (line.length > 5 && !line.includes("|") && !line.includes("---")) {
                  currentCat.push(line);
                }
              }
            }

            // Extract a clean title and summary (avoiding ## Installation as title/summary)
            let title = gitRelease.name || `Release ${version}`;
            if (title.startsWith("v")) {
              title = `Release ${title}`;
            }

            let summary = "Updates and performance enhancements.";
            if (features.length > 0) {
              summary = features[0].length > 150 ? features[0].slice(0, 150) + "..." : features[0];
            } else if (fixes.length > 0) {
              summary = fixes[0].length > 150 ? fixes[0].slice(0, 150) + "..." : fixes[0];
            }

            return {
              version,
              date: formattedDate,
              title,
              summary,
              githubUrl: gitRelease.html_url,
              features: features.length > 0 ? features : ["Refer to GitHub release details."],
              fixes,
              performance,
              metrics: {
                binarySize: "8.0 MB",
                scanLatency: "0.3 ms",
                themes: 9,
                visualizers: 6
              }
            };
          });

          // Sort releases to keep newest on top
          merged.sort((a, b) => {
            return b.version.localeCompare(a.version, undefined, { numeric: true, sensitivity: 'base' });
          });
          setReleases(merged);
        }
      } catch (e) {
        console.error("Failed to fetch live GitHub releases:", e);
      }
    }
    fetchReleases();
  }, []);

  const toggleExpand = (ver: string) => {
    setExpandedCards(prev => ({ ...prev, [ver]: !prev[ver] }));
  };

  const latestVersion = releases[0]?.version || "v1.0.5";

  return (
    <div style={{ background: "#080808", minHeight: "100vh", color: "#ffffff", paddingBottom: "100px", position: "relative" }}>
      {/* Subtle grid */}
      <div
        className="absolute inset-0 pointer-events-none opacity-[0.025]"
        style={{
          backgroundImage:
            "linear-gradient(rgba(255,255,255,0.15) 1px, transparent 1px), linear-gradient(90deg, rgba(255,255,255,0.15) 1px, transparent 1px)",
          backgroundSize: "64px 64px",
        }}
      />

      {/* Top subtle glow */}
      <div
        className="absolute top-0 left-1/2 -translate-x-1/2 w-[800px] h-[300px] pointer-events-none"
        style={{
          background: "radial-gradient(ellipse at center top, rgba(255,255,255,0.02) 0%, transparent 70%)",
        }}
      />

      <Nav version={latestVersion} />

      {/* Hero Header */}
      <section className="relative pt-32 pb-16 overflow-hidden border-b border-white/[0.05]">
        <div className="container-page text-center relative z-10">
          <p className="font-mono text-xs text-[#555] uppercase tracking-[0.2em] mb-4">
            <SplitText text="Version History & Timeline" by="chars" delay={0.15} stagger={0.03} direction="down" />
          </p>
          <h1 className="font-mono font-semibold text-white tracking-tight leading-tight mb-5" style={{ fontSize: "clamp(2rem, 5vw, 3.5rem)" }}>
            <SplitText text="Changelog" by="chars" delay={0.35} stagger={0.05} direction="down" />
          </h1>
          <FadeIn delay={0.7} y={15}>
            <p className="text-[#666] max-w-xl mx-auto text-sm md:text-base leading-relaxed">
              Follow the journey of Moodwave CLI as it evolves from a lightweight mood audio scanner to a highly optimized terminal companion.
            </p>
          </FadeIn>
        </div>
      </section>

      {/* Codebase Evolution Metrics Dashboard */}
      <section className="py-12 border-b border-white/[0.05] bg-white/[0.01]">
        <div className="container-page">
          <div className="border border-white/[0.08] rounded-2xl bg-zinc-950/40 p-6 md:p-8 backdrop-blur-md shadow-2xl relative overflow-hidden">
            {/* Glowing background blob */}
            <div className="absolute -right-32 -top-32 w-96 h-96 rounded-full bg-white/[0.01] blur-[120px] pointer-events-none" />

            <div className="flex flex-col md:flex-row items-start md:items-center justify-between gap-6 mb-8 relative z-10">
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

            {/* Interactive Spline Graph */}
            {(() => {
              const chartData = [
                { version: "v1.0.0", speed: 1500, size: 8.9, themes: 3, labelSpeed: "1500ms", labelSize: "8.9MB", labelThemes: "3 themes", date: "Jul 1" },
                { version: "v1.0.1", speed: 18,   size: 8.8, themes: 3, labelSpeed: "18ms",   labelSize: "8.8MB", labelThemes: "3 themes", date: "Jul 7" },
                { version: "v1.0.2", speed: 0.4,  size: 8.1, themes: 9, labelSpeed: "0.4ms",  labelSize: "8.1MB", labelThemes: "9 themes", date: "Jul 7" },
                { version: "v1.0.3", speed: 0.4,  size: 8.1, themes: 9, labelSpeed: "0.4ms",  labelSize: "8.1MB", labelThemes: "9 themes", date: "Jul 9" },
                { version: "v1.0.4", speed: 0.3,  size: 8.0, themes: 9, labelSpeed: "0.3ms",  labelSize: "8.0MB", labelThemes: "9 themes", date: "Jul 9" },
                { version: "v1.0.5", speed: 0.3,  size: 8.0, themes: 9, labelSpeed: "0.3ms",  labelSize: "8.0MB", labelThemes: "9 themes", date: "Jul 10" },
              ];

              const xCoords = [60, 196, 332, 468, 604, 740];
              const points = chartData.map((d, idx) => {
                let y = 160;
                if (activeMetricTab === "speed") {
                  const ys = [160, 95, 35, 35, 31, 30];
                  y = ys[idx];
                } else if (activeMetricTab === "size") {
                  const ys = [160, 148, 45, 45, 32, 30];
                  y = ys[idx];
                } else {
                  const ys = [160, 160, 30, 30, 30, 30];
                  y = ys[idx];
                }
                return { 
                  x: xCoords[idx], 
                  y, 
                  label: activeMetricTab === "speed" ? d.labelSpeed : activeMetricTab === "size" ? d.labelSize : d.labelThemes, 
                  version: d.version,
                  date: d.date
                };
              });

              const strokeColor = activeMetricTab === "speed" ? "#ef4444" : activeMetricTab === "size" ? "#3b82f6" : "#10b981";
              const gradId = `${activeMetricTab}Grad`;
              const linePath = getBezierPath(points);
              const areaPath = `${linePath} L 740 180 L 60 180 Z`;

              return (
                <div className="h-[260px] w-full relative pb-4 z-10" onMouseLeave={() => setHoveredIdx(null)}>
                  {/* Tooltip Overlay */}
                  {hoveredIdx !== null && (
                    <motion.div 
                      initial={{ opacity: 0, y: 5 }}
                      animate={{ opacity: 1, y: 0 }}
                      className="absolute bg-[#121214]/90 border border-white/[0.08] rounded-xl p-3 shadow-2xl backdrop-blur-md pointer-events-none transition-all duration-150"
                      style={{
                        left: `${(points[hoveredIdx].x / 800) * 100}%`,
                        bottom: `${100 - (points[hoveredIdx].y / 200) * 100 + 8}%`,
                        transform: "translateX(-50%)",
                        zIndex: 30
                      }}
                    >
                      <div className="flex items-center gap-1.5 mb-1">
                        <span className="font-mono text-[9px] font-semibold px-1.5 py-0.5 bg-white/[0.06] border border-white/[0.08] rounded text-white/50">{points[hoveredIdx].version}</span>
                        <span className="font-mono text-[9px] text-[#555]">{points[hoveredIdx].date}</span>
                      </div>
                      <div className="font-mono text-xs font-bold text-white">{points[hoveredIdx].label}</div>
                      <div className="font-mono text-[8px] text-[#666] mt-0.5">
                        {activeMetricTab === "speed" ? "Execution speed latency" : activeMetricTab === "size" ? "CLI static binary footprint" : "Configurable color presets"}
                      </div>
                    </motion.div>
                  )}

                  <svg className="w-full h-full" viewBox="0 0 800 200" preserveAspectRatio="none">
                    <defs>
                      <linearGradient id="speedGrad" x1="0" y1="0" x2="0" y2="1">
                        <stop offset="0%" stopColor="#ef4444" stopOpacity="0.25" />
                        <stop offset="100%" stopColor="#ef4444" stopOpacity="0" />
                      </linearGradient>
                      <linearGradient id="sizeGrad" x1="0" y1="0" x2="0" y2="1">
                        <stop offset="0%" stopColor="#3b82f6" stopOpacity="0.25" />
                        <stop offset="100%" stopColor="#3b82f6" stopOpacity="0" />
                      </linearGradient>
                      <linearGradient id="themesGrad" x1="0" y1="0" x2="0" y2="1">
                        <stop offset="0%" stopColor="#10b981" stopOpacity="0.25" />
                        <stop offset="100%" stopColor="#10b981" stopOpacity="0" />
                      </linearGradient>
                    </defs>

                    {/* Horizontal helper dashed grid lines */}
                    {[30, 80, 130, 180].map((gridY, i) => (
                      <g key={gridY} className="opacity-40">
                        <line x1="60" y1={gridY} x2="740" y2={gridY} stroke="white" strokeWidth="0.5" strokeDasharray="3 6" />
                        <text x="25" y={gridY + 3} fill="#444" fontSize="8" fontFamily="monospace">
                          {i === 0 ? "PEAK" : i === 1 ? "MID" : i === 2 ? "SLOW" : "INIT"}
                        </text>
                      </g>
                    ))}

                    {/* Chart Area Fill */}
                    <motion.path 
                      d={areaPath} 
                      fill={`url(#${gradId})`} 
                      animate={{ d: areaPath }}
                      transition={{ duration: 0.4, ease: "easeInOut" }}
                    />

                    {/* Spline Line */}
                    <motion.path 
                      d={linePath} 
                      fill="none" 
                      stroke={strokeColor} 
                      strokeWidth="2" 
                      animate={{ d: linePath }}
                      transition={{ duration: 0.4, ease: "easeInOut" }}
                    />

                    {/* Version nodes and labels */}
                    {points.map((pt, idx) => {
                      const isHovered = hoveredIdx === idx;
                      return (
                        <g key={pt.version}>
                          {/* Inner glowing circle */}
                          <circle 
                            cx={pt.x} 
                            cy={pt.y} 
                            r={isHovered ? "7" : "4"} 
                            fill={strokeColor} 
                            className="transition-all duration-150"
                            style={{ filter: isHovered ? `drop-shadow(0 0 6px ${strokeColor})` : "none" }}
                          />
                          {/* Outer halo */}
                          <circle 
                            cx={pt.x} 
                            cy={pt.y} 
                            r={isHovered ? "12" : "8"} 
                            fill="none" 
                            stroke={strokeColor} 
                            strokeWidth="1" 
                            strokeOpacity={isHovered ? "0.6" : "0.15"}
                            className="transition-all duration-150"
                          />
                          {/* X-axis version tags */}
                          <text 
                            x={pt.x} 
                            y="196" 
                            textAnchor="middle" 
                            fill={isHovered ? "#fff" : "#444"} 
                            fontSize="8" 
                            fontFamily="monospace"
                            className="transition-colors duration-150"
                          >
                            {pt.version}
                          </text>

                          {/* Interactive invisible hover target */}
                          <circle 
                            cx={pt.x} 
                            cy={pt.y} 
                            r="24" 
                            fill="transparent" 
                            className="cursor-pointer"
                            onMouseEnter={() => setHoveredIdx(idx)}
                          />
                        </g>
                      );
                    })}
                  </svg>
                </div>
              );
            })()}

            {/* Visual indicators */}
            <div className="grid grid-cols-3 gap-4 mt-6 text-center relative z-10 border-t border-white/[0.04] pt-6">
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
            
            {releases.map((release, index) => {
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



