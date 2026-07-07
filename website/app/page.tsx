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
    // Fetch latest release tag name dynamically from GitHub API at build/render time
    // Next.js cache keeps this optimized for 1 hour to prevent hitting API rate limits
    const res = await fetch("https://api.github.com/repos/Boredooms/Moodwave-CLI/releases/latest", {
      next: { revalidate: 3600 },
      headers: {
        "User-Agent": "Moodwave-Website-Builder"
      }
    });
    if (!res.ok) return "v1.0.1";
    const data = await res.json();
    return data.tag_name || "v1.0.1";
  } catch (e) {
    return "v1.0.1";
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
