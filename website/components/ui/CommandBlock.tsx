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
    <div className="group relative w-full min-w-0">
      {label && (
        <p className="text-xs text-[#555] font-mono uppercase tracking-widest mb-2">{label}</p>
      )}
      <div className="code-block w-full min-w-0 group-hover:border-white/[0.12] transition-colors duration-300 flex items-center justify-between">
        {/* Scrollable command container */}
        <div 
          className="flex items-center gap-3 min-w-0 flex-1 overflow-x-auto pr-3 scrollbar-none"
          style={{ scrollbarWidth: "none" }} // Firefox
        >
          <span className="text-[#555] font-mono text-xs flex-shrink-0 select-none">{prompt}</span>
          <code className="text-[#c9c9c9] font-mono text-xs whitespace-nowrap tracking-tight">{command}</code>
        </div>
        
        {/* Copy button - stays pinned to the right */}
        <button
          onClick={handleCopy}
          className="flex-shrink-0 flex items-center gap-1.5 text-xs font-mono text-[#555] hover:text-white transition-colors duration-200 cursor-pointer select-none bg-[#0d0d0d] pl-2"
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
