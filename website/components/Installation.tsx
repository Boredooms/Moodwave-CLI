"use client";

import { useState } from "react";
import { motion, useInView } from "framer-motion";
import { useRef } from "react";
import CommandBlock from "./ui/CommandBlock";

const tabs = [
  {
    id: "unix",
    label: "macOS / Linux",
    commands: [
      {
        label: "Install",
        prompt: "$",
        command: "curl -fsSL https://raw.githubusercontent.com/Boredooms/Moodwave-CLI/main/cli/scripts/install.sh | bash",
      },
      {
        label: "Start",
        prompt: "$",
        command: "moodwave",
      },
      {
        label: "Update",
        prompt: "$",
        command: "moodwave update",
      },
    ],
  },
  {
    id: "windows",
    label: "Windows",
    commands: [
      {
        label: "Install (PowerShell)",
        prompt: "PS>",
        command: "irm https://raw.githubusercontent.com/Boredooms/Moodwave-CLI/main/cli/scripts/install.ps1 | iex",
      },
      {
        label: "Start",
        prompt: "PS>",
        command: "moodwave",
      },
      {
        label: "Update",
        prompt: "PS>",
        command: "moodwave update",
      },
    ],
  },
  {
    id: "go",
    label: "From source",
    commands: [
      {
        label: "Clone & build",
        prompt: "$",
        command: "git clone https://github.com/Boredooms/Moodwave-CLI && cd Moodwave-CLI/cli && go install ./cmd/moodwave",
      },
    ],
  },
];

export default function Installation() {
  const [activeTab, setActiveTab] = useState("unix");
  const ref = useRef<HTMLDivElement>(null);
  const inView = useInView(ref, { once: true, margin: "-80px 0px" });

  const current = tabs.find((t) => t.id === activeTab)!;

  return (
    <section className="divider section-pad" id="install">
      <div className="container-page">
        <div className="grid lg:grid-cols-[280px_1fr] gap-12 lg:gap-20">
          <motion.div
            ref={ref}
            initial={{ opacity: 0, y: 20 }}
            animate={inView ? { opacity: 1, y: 0 } : {}}
            transition={{ duration: 0.6 }}
          >
            <p className="font-mono text-xs text-[#444] uppercase tracking-[0.2em] mb-5">Installation</p>
            <h2 className="font-mono text-display-md text-white font-semibold leading-tight">
              One command. Any platform.
            </h2>
            <p className="text-body-sm text-[#555] mt-4 leading-relaxed">
              Pre-compiled binaries for every major OS and CPU architecture.
              No Go runtime, no build step, no dependency hell.
            </p>

            <div className="mt-8 space-y-3">
              <p className="font-mono text-xs text-[#444] uppercase tracking-widest">Supported platforms</p>
              {["Windows amd64 / arm64", "macOS Intel / Apple Silicon", "Linux amd64 / arm64 / arm"].map((p) => (
                <div key={p} className="flex items-center gap-3">
                  <span className="w-1 h-1 bg-[#444] rounded-full flex-shrink-0" />
                  <span className="font-mono text-xs text-[#666]">{p}</span>
                </div>
              ))}
            </div>

            <div className="mt-8">
              <a
                href="https://github.com/Boredooms/Moodwave-CLI/releases"
                target="_blank"
                rel="noopener noreferrer"
                className="font-mono text-xs text-[#555] hover:text-white transition-colors flex items-center gap-2"
              >
                Browse all releases ↗
              </a>
            </div>
          </motion.div>

          <motion.div
            initial={{ opacity: 0, x: 20 }}
            animate={inView ? { opacity: 1, x: 0 } : {}}
            transition={{ duration: 0.6, delay: 0.15 }}
          >
            {/* Tab row */}
            <div className="flex gap-0 border border-white/[0.07] rounded-lg overflow-hidden mb-6 w-fit">
              {tabs.map((tab) => (
                <button
                  key={tab.id}
                  onClick={() => setActiveTab(tab.id)}
                  className={`font-mono text-xs px-5 py-2.5 transition-colors duration-200 cursor-pointer ${
                    activeTab === tab.id
                      ? "bg-white/[0.08] text-white"
                      : "text-[#555] hover:text-[#888] hover:bg-white/[0.03]"
                  } border-r border-white/[0.07] last:border-r-0`}
                >
                  {tab.label}
                </button>
              ))}
            </div>

            {/* Commands */}
            <div className="space-y-3">
              {current.commands.map((cmd) => (
                <CommandBlock
                  key={cmd.label}
                  label={cmd.label}
                  prompt={cmd.prompt}
                  command={cmd.command}
                />
              ))}
            </div>

            {/* Note */}
            <p className="font-mono text-xs text-[#444] mt-6">
              The installer detects your OS and architecture automatically.
              Run <span className="text-[#666]">moodwave doctor</span> after install to verify everything is working.
            </p>
          </motion.div>
        </div>
      </div>
    </section>
  );
}
