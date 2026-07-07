import Nav from "../components/Nav";
import Hero from "../components/Hero";
import WhySection from "../components/WhySection";
import HowItWorks from "../components/HowItWorks";
import TerminalShowcase from "../components/TerminalShowcase";
import Features from "../components/Features";
import Architecture from "../components/Architecture";
import Installation from "../components/Installation";
import FAQ from "../components/FAQ";
import CTA from "../components/CTA";

export default function Home() {
  return (
    <main style={{ background: "#080808", minHeight: "100vh", overflowX: "hidden" }}>
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
