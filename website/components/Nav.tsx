"use client";

import { useState } from "react";
import { motion, AnimatePresence } from "framer-motion";

export default function Nav({ version = "v1.0.1" }: { version?: string }) {
  const [menuOpen, setMenuOpen] = useState(false);

  const links = [
    { href: "#how-it-works", label: "How it works" },
    { href: "#install", label: "Install" },
    { href: "https://github.com/Boredooms/Moodwave-CLI", label: "GitHub", external: true },
  ];

  return (
    <header className="fixed top-0 left-0 right-0 z-50 border-b border-white/[0.05]" style={{ background: "rgba(8,8,8,0.85)", backdropFilter: "blur(12px)" }}>
      <div className="container-page flex items-center justify-between" style={{ height: "56px" }}>
        <a href="#" className="flex items-center gap-2.5 group">
          <svg className="w-5 h-5 text-white transition-transform group-hover:rotate-12 duration-300" viewBox="0 0 32 32" fill="none">
            <circle cx="16" cy="16" r="14" fill="#080808" stroke="currentColor" strokeWidth="1.5"/>
            <path d="M9 16h1.5M13 11v10M17 7v18M21 13v6M25 16h-1.5" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
          </svg>
          <span className="font-mono text-sm font-semibold text-white">moodwave</span>
          <span className="font-mono text-xs text-[#444] group-hover:text-[#666] transition-colors">{version}</span>
        </a>

        {/* Desktop */}
        <nav className="hidden md:flex items-center gap-8">
          {links.map((link) => (
            <a
              key={link.href}
              href={link.href}
              target={link.external ? "_blank" : undefined}
              rel={link.external ? "noopener noreferrer" : undefined}
              className="font-mono text-sm text-[#666] hover:text-white transition-colors duration-200"
            >
              {link.label}
            </a>
          ))}
        </nav>

        {/* Mobile toggle */}
        <button
          className="md:hidden font-mono text-xs text-[#666] hover:text-white transition-colors cursor-pointer"
          onClick={() => setMenuOpen((v) => !v)}
        >
          {menuOpen ? "close" : "menu"}
        </button>
      </div>

      <AnimatePresence>
        {menuOpen && (
          <motion.div
            initial={{ height: 0, opacity: 0 }}
            animate={{ height: "auto", opacity: 1 }}
            exit={{ height: 0, opacity: 0 }}
            transition={{ duration: 0.22 }}
            className="md:hidden border-t border-white/[0.05] overflow-hidden"
          >
            <div className="container-page py-4 flex flex-col gap-4">
              {links.map((link) => (
                <a
                  key={link.href}
                  href={link.href}
                  onClick={() => setMenuOpen(false)}
                  className="font-mono text-sm text-[#666] hover:text-white transition-colors"
                >
                  {link.label}
                </a>
              ))}
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </header>
  );
}
