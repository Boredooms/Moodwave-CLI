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

async function getLatestVersion() {
  try {
    // Fetch repository tags directly from GitHub API
    // This allows immediate tracking of pushed tags (e.g. v1.0.0) without waiting for release assets build
    const res = await fetch("https://api.github.com/repos/Boredooms/Moodwave-CLI/tags", {
      next: { revalidate: 3600 },
      headers: {
        "User-Agent": "Moodwave-Website-Builder"
      }
    });
    if (!res.ok) return "v1.0.0";
    const tags = await res.json();
    if (tags && tags.length > 0) {
      return tags[0].name || "v1.0.0";
    }
    return "v1.0.0";
  } catch (e) {
    return "v1.0.0";
  }
}

export default async function Home() {
  const version = await getLatestVersion();

  return (
    <main style={{ background: "#080808", minHeight: "100vh", overflowX: "hidden" }}>
      <Nav version={version} />
      <Hero version={version} />
      <WhySection />
      <HowItWorks />
      <TerminalShowcase />
      <Features />
      <Architecture />
      <Installation />
      <FAQ />
      <CTA version={version} />
    </main>
  );
}
