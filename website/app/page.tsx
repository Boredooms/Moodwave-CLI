"use client";

import { useEffect } from "react";
import Lenis from "lenis";
import Nav from "@/components/Nav";
import Hero from "@/components/Hero";
import WhySection from "@/components/WhySection";
import HowItWorks from "@/components/HowItWorks";
import TerminalShowcase from "@/components/TerminalShowcase";
import Features from "@/components/Features";
import Architecture from "@/components/Architecture";
import Installation from "@/components/Installation";
import FAQ from "@/components/FAQ";
import CTA from "@/components/CTA";

export default function Home() {
  useEffect(() => {
    const lenis = new Lenis({
      duration: 1.2,
      easing: (t) => Math.min(1, 1.001 - Math.pow(2, -10 * t)),
      smoothWheel: true,
    });

    function raf(time: number) {
      lenis.raf(time);
      requestAnimationFrame(raf);
    }

    requestAnimationFrame(raf);

    return () => {
      lenis.destroy();
    };
  }, []);

  return (
    <main className="bg-bg min-h-screen overflow-hidden">
      <Nav />
      <Hero />
      <WhySection />
      <HowItWorks />
      <TerminalShowcase />
      <Features />
      <Architecture />
      <Installation />
      <FAQ />
      <CTA />
    </main>
  );
}
