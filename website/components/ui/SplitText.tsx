"use client";

import { motion } from "framer-motion";

interface SplitTextProps {
  text: string;
  /** "words" splits by word, "chars" splits by character */
  by?: "words" | "chars";
  className?: string;
  delay?: number;
  stagger?: number;
  duration?: number;
  /** Direction of the fall — "up" rises into place, "down" falls into place */
  direction?: "up" | "down";
  once?: boolean;
  as?: "h1" | "h2" | "h3" | "p" | "span";
}

const ITEM_VARIANTS = (direction: "up" | "down") => ({
  hidden: {
    opacity: 0,
    y: direction === "up" ? 32 : -32,
    rotateX: direction === "up" ? 12 : -12,
  },
  show: {
    opacity: 1,
    y: 0,
    rotateX: 0,
    transition: {
      duration: 0.65,
      ease: [0.16, 1, 0.3, 1] as [number, number, number, number],
    },
  },
});

export default function SplitText({
  text,
  by = "words",
  className = "",
  delay = 0,
  stagger = 0.06,
  duration = 0.65,
  direction = "up",
  once = true,
  as: Tag = "span",
}: SplitTextProps) {
  const tokens = by === "words" ? text.split(" ") : text.split("");
  const itemVariants = ITEM_VARIANTS(direction);

  const containerVariants = {
    hidden: {},
    show: {
      transition: {
        staggerChildren: stagger,
        delayChildren: delay,
      },
    },
  };

  return (
    <Tag style={{ display: "inline", perspective: "800px" }}>
      <motion.span
        initial="hidden"
        whileInView="show"
        viewport={{ once, margin: "-40px 0px" }}
        variants={containerVariants}
        style={{ display: "inline" }}
        aria-label={text}
      >
        {tokens.map((token, i) => (
          <span
            key={i}
            style={{ display: "inline-block", overflow: "hidden", lineHeight: 1.1, paddingBottom: "0.05em" }}
            aria-hidden="true"
          >
            <motion.span
              variants={itemVariants}
              style={{ display: "inline-block" }}
            >
              {token}
              {by === "words" && i < tokens.length - 1 ? "\u00A0" : ""}
            </motion.span>
          </span>
        ))}
      </motion.span>
    </Tag>
  );
}
