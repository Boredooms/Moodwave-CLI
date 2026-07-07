"use client";

import { useEffect } from "react";
import Lenis from "lenis";

/**
 * LenisProvider — initializes Lenis smooth scroll once at the app level.
 *
 * Key difference from the broken version:
 * - No React state is mutated inside the RAF loop
 * - Lenis is created once, destroyed on unmount
 * - No interference with Framer Motion's IntersectionObserver (whileInView)
 *   because Lenis doesn't touch element positions, only window.scrollY
 */
export default function LenisProvider({ children }: { children: React.ReactNode }) {
  useEffect(() => {
    const lenis = new Lenis({
      duration: 1.3,
      easing: (t: number) => 1 - Math.pow(1 - t, 4), // ease-out-quart
      smoothWheel: true,
      wheelMultiplier: 0.9,
      touchMultiplier: 1.5,
      infinite: false,
    });

    // Use the native RAF loop — no requestAnimationFrame wrapper that leaks
    let rafId: number;
    const raf = (time: number) => {
      lenis.raf(time);
      rafId = requestAnimationFrame(raf);
    };
    rafId = requestAnimationFrame(raf);

    // Handle anchor links
    const handleAnchor = (e: MouseEvent) => {
      const target = (e.target as Element).closest("a[href^='#']");
      if (!target) return;
      const href = target.getAttribute("href");
      if (!href) return;
      const el = document.querySelector(href);
      if (!el) return;
      e.preventDefault();
      lenis.scrollTo(el as HTMLElement, { offset: -56, duration: 1.4 });
    };
    document.addEventListener("click", handleAnchor);

    return () => {
      cancelAnimationFrame(rafId);
      lenis.destroy();
      document.removeEventListener("click", handleAnchor);
    };
  }, []);

  return <>{children}</>;
}
