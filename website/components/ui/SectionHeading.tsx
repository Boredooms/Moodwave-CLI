"use client";

import { motion, useInView } from "framer-motion";
import { useRef } from "react";

interface SectionHeadingProps {
  eyebrow?: string;
  title: string;
  subtitle?: string;
  align?: "left" | "center";
  className?: string;
}

export default function SectionHeading({
  eyebrow,
  title,
  subtitle,
  align = "left",
  className = "",
}: SectionHeadingProps) {
  const ref = useRef<HTMLDivElement>(null);
  const inView = useInView(ref, { once: true, margin: "-80px 0px" });

  const textAlign = align === "center" ? "text-center" : "text-left";
  const maxWidth = align === "center" ? "max-w-2xl mx-auto" : "";

  return (
    <motion.div
      ref={ref}
      initial={{ opacity: 0, y: 24 }}
      animate={inView ? { opacity: 1, y: 0 } : {}}
      transition={{ duration: 0.6, ease: [0.22, 1, 0.36, 1] }}
      className={`${textAlign} ${maxWidth} ${className}`}
    >
      {eyebrow && (
        <p className="text-xs font-mono text-[#555] uppercase tracking-[0.2em] mb-4">
          {eyebrow}
        </p>
      )}
      <h2 className="text-display-md font-mono font-semibold text-white leading-tight mb-4">
        {title}
      </h2>
      {subtitle && (
        <p className="text-body-lg text-[#666] max-w-xl leading-relaxed">
          {subtitle}
        </p>
      )}
    </motion.div>
  );
}
