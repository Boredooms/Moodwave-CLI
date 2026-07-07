"use client";

import { useState } from "react";
import { motion } from "framer-motion";
import CommandBlock from "./ui/CommandBlock";

const tabs = [
  {
    id: "unix", label: "macOS / Linux",
    commands: [
      { label: "Install", prompt: "$", command: "curl -fsSL https://raw.githubusercontent.com/Boredooms/Moodwave-CLI/main/cli/scripts/install.sh | bash" },
      { label: "Start", prompt: "$", command: "moodwave" },
      { label: "Update", prompt: "$", command: "moodwave update" },
    ],
  },
  {
    id: "windows", label: "Windows",
    commands: [
      { label: "Install (PowerShell)", prompt: "PS>", command: "irm https://raw.githubusercontent.com/Boredooms/Moodwave-CLI/main/cli/scripts/install.ps1 | iex" },
      { label: "Start", prompt: "PS>", command: "moodwave" },
      { label: "Update", prompt: "PS>", command: "moodwave update" },
    ],
  },
  {
    id: "go", label: "From source",
    commands: [
      { label: "Clone & build", prompt: "$", command: "git clone https://github.com/Boredooms/Moodwave-CLI && cd Moodwave-CLI/cli && go install ./cmd/moodwave" },
    ],
  },
];

export default function Installation() {
  const [activeTab, setActiveTab] = useState("unix");
  const current = tabs.find((t) => t.id === activeTab)!;

  return (
    <section className="divider section-pad" id="install">
      <div className="container-page">
        <div className="grid lg:grid-cols-[240px_1fr] gap-12 lg:gap-20">
          <motion.div
            initial={{ opacity: 0, y: 16 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true, margin: "-60px" }}
            transition={{ duration: 0.55 }}
          >
            <p className="font-mono text-xs text-[#444] uppercase tracking-[0.2em] mb-5">Installation</p>
            <h2 className="font-mono font-semibold text-white leading-tight" style={{ fontSize: "clamp(1.4rem, 2.5vw, 2rem)", letterSpacing: "-0.025em" }}>
              One command.<br />Any platform.
            </h2>
            <p className="text-sm text-[#555] mt-4 leading-relaxed">
              Pre-compiled binaries for every major OS and CPU. No Go runtime, no build step.
            </p>

            <div className="mt-8 space-y-2.5">
              <p className="font-mono text-xs text-[#333] uppercase tracking-widest mb-3">Supported</p>
              {["Windows amd64 / arm64", "macOS Intel / Apple Silicon", "Linux amd64 / arm64 / arm"].map((p) => (
                <div key={p} className="flex items-center gap-3">
                  <span className="w-1 h-1 bg-[#444] rounded-full flex-shrink-0" />
                  <span className="font-mono text-xs text-[#666]">{p}</span>
                </div>
              ))}
            </div>

            <div className="mt-8">
              <a href="https://github.com/Boredooms/Moodwave-CLI/releases" target="_blank" rel="noopener noreferrer" className="font-mono text-xs text-[#555] hover:text-white transition-colors">
                Browse all releases ↗
              </a>
            </div>
          </motion.div>

          <motion.div
            initial={{ opacity: 0, x: 16 }}
            whileInView={{ opacity: 1, x: 0 }}
            viewport={{ once: true, margin: "-60px" }}
            transition={{ duration: 0.55, delay: 0.12 }}
          >
            {/* Tabs */}
            <div className="flex gap-0 border border-white/[0.07] rounded-lg overflow-hidden mb-6 w-fit">
              {tabs.map((tab) => (
                <button
                  key={tab.id}
                  onClick={() => setActiveTab(tab.id)}
                  className={`font-mono text-xs px-5 py-2.5 cursor-pointer transition-colors duration-200 border-r border-white/[0.07] last:border-r-0 ${
                    activeTab === tab.id ? "bg-white/[0.08] text-white" : "text-[#555] hover:text-[#888]"
                  }`}
                >
                  {tab.label}
                </button>
              ))}
            </div>

            <div className="space-y-3">
              {current.commands.map((cmd) => (
                <CommandBlock key={cmd.label} label={cmd.label} prompt={cmd.prompt} command={cmd.command} />
              ))}
            </div>

            <p className="font-mono text-xs text-[#444] mt-6">
              Run <span className="text-[#666]">moodwave doctor</span> after install to verify everything is working.
            </p>
          </motion.div>
        </div>
      </div>
    </section>
  );
}
