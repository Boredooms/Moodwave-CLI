"use client";

import { useState } from "react";
import { motion, AnimatePresence } from "framer-motion";

interface CommandBlockProps {
  command: string;
  label?: string;
  prompt?: string;
}

export default function CommandBlock({ command, label, prompt = "$" }: CommandBlockProps) {
  const [copied, setCopied] = useState(false);

  const handleCopy = async () => {
    await navigator.clipboard.writeText(command);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="group relative">
      {label && (
        <p className="text-xs text-[#555] font-mono uppercase tracking-widest mb-2">{label}</p>
      )}
      <div className="code-block group-hover:border-white/[0.12] transition-colors duration-300">
        <div className="flex items-center gap-3 min-w-0 flex-1">
          <span className="text-[#555] font-mono text-sm flex-shrink-0">{prompt}</span>
          <code className="text-[#c9c9c9] font-mono text-sm truncate">{command}</code>
        </div>
        <button
          onClick={handleCopy}
          className="flex-shrink-0 flex items-center gap-1.5 text-xs font-mono text-[#555] hover:text-white transition-colors duration-200 cursor-pointer"
          aria-label="Copy command"
        >
          <AnimatePresence mode="wait">
            {copied ? (
              <motion.span
                key="copied"
                initial={{ opacity: 0, y: 4 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -4 }}
                className="text-[#aaa]"
              >
                copied
              </motion.span>
            ) : (
              <motion.span
                key="copy"
                initial={{ opacity: 0, y: 4 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -4 }}
              >
                copy
              </motion.span>
            )}
          </AnimatePresence>
        </button>
      </div>
    </div>
  );
}
