"use client";

import { useState } from "react";
import { motion, AnimatePresence } from "framer-motion";
import Link from "next/link";

export default function Nav() {
  const [menuOpen, setMenuOpen] = useState(false);

  const links = [
    { href: "#how-it-works", label: "How it works" },
    { href: "#install", label: "Install" },
    { href: "https://github.com/Boredooms/Moodwave-CLI", label: "GitHub", external: true },
  ];

  return (
    <header className="fixed top-0 left-0 right-0 z-50 border-b border-white/[0.05] bg-bg/80 backdrop-blur-md">
      <div className="container-page flex items-center justify-between h-14">
        <a href="#" className="flex items-center gap-2 group">
          <span className="font-mono text-sm font-semibold tracking-tight text-white">moodwave</span>
          <span className="font-mono text-xs text-[#444] group-hover:text-[#666] transition-colors">cli</span>
        </a>

        {/* Desktop nav */}
        <nav className="hidden md:flex items-center gap-8">
          {links.map((link) => (
            <a
              key={link.href}
              href={link.href}
              target={link.external ? "_blank" : undefined}
              rel={link.external ? "noopener noreferrer" : undefined}
              className="text-sm text-[#666] hover:text-white transition-colors duration-200 font-mono"
            >
              {link.label}
            </a>
          ))}
        </nav>

        {/* Mobile menu toggle */}
        <button
          className="md:hidden text-[#666] hover:text-white transition-colors font-mono text-xs"
          onClick={() => setMenuOpen(!menuOpen)}
          aria-label="Toggle menu"
        >
          {menuOpen ? "close" : "menu"}
        </button>
      </div>

      {/* Mobile nav */}
      <AnimatePresence>
        {menuOpen && (
          <motion.div
            initial={{ height: 0, opacity: 0 }}
            animate={{ height: "auto", opacity: 1 }}
            exit={{ height: 0, opacity: 0 }}
            transition={{ duration: 0.25 }}
            className="md:hidden border-t border-white/[0.05] overflow-hidden"
          >
            <div className="container-page py-4 flex flex-col gap-4">
              {links.map((link) => (
                <a
                  key={link.href}
                  href={link.href}
                  onClick={() => setMenuOpen(false)}
                  className="text-sm text-[#666] hover:text-white transition-colors font-mono"
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
