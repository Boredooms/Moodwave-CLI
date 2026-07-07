"use client";

import { motion } from "framer-motion";
import type { ReactNode } from "react";

interface FadeInProps {
  children: ReactNode;
  delay?: number;
  duration?: number;
  y?: number;
  x?: number;
  className?: string;
  once?: boolean;
}

/**
 * FadeIn — wraps any content with a scroll-triggered fade + translate reveal.
 * Uses Framer Motion's whileInView (no useInView hook = no re-render bounce).
 */
export default function FadeIn({
  children,
  delay = 0,
  duration = 0.65,
  y = 24,
  x = 0,
  className = "",
  once = true,
}: FadeInProps) {
  return (
    <motion.div
      initial={{ opacity: 0, y, x }}
      whileInView={{ opacity: 1, y: 0, x: 0 }}
      viewport={{ once, margin: "-50px 0px" }}
      transition={{
        duration,
        delay,
        ease: [0.16, 1, 0.3, 1],
      }}
      className={className}
    >
      {children}
    </motion.div>
  );
}
