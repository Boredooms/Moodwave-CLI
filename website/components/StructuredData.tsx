export function GlobalJsonLd() {
  const baseUrl = "https://www.moodwave-cli.xyz";

  // 1. WebSite Schema (Enables Google Brand Search & Sitelinks Searchbox)
  const websiteSchema = {
    "@context": "https://schema.org",
    "@type": "WebSite",
    "@id": `${baseUrl}/#website`,
    url: baseUrl,
    name: "Moodwave",
    alternateName: ["Moodwave CLI", "Moodwave-CLI", "Moodwave Terminal Music"],
    description:
      "Moodwave scans your codebase, infers your working mood, and streams perfectly matched music right in your terminal. CLI-first, lightweight, open source.",
    publisher: {
      "@id": `${baseUrl}/#organization`,
    },
    inLanguage: "en-US",
  };

  // 2. Organization Schema (Enables Google Knowledge Graph & Social Verification)
  const organizationSchema = {
    "@context": "https://schema.org",
    "@type": "Organization",
    "@id": `${baseUrl}/#organization`,
    name: "Moodwave",
    url: baseUrl,
    logo: `${baseUrl}/logo.svg`,
    sameAs: [
      "https://github.com/Boredooms/Moodwave-CLI",
      "https://www.linkedin.com/company/moodwave-cli/",
    ],
    description: "Open source terminal-native developer mood music companion.",
  };

  // 3. SoftwareApplication Schema (Enables Google Software Rich Cards)
  const softwareAppSchema = {
    "@context": "https://schema.org",
    "@type": "SoftwareApplication",
    "@id": `${baseUrl}/#software`,
    name: "Moodwave CLI",
    operatingSystem: "Windows, macOS, Linux",
    applicationCategory: "DeveloperApplication",
    applicationSubCategory: "Terminal Music Player & Workspace Productivity",
    softwareVersion: "v2.0.0",
    offers: {
      "@type": "Offer",
      price: "0",
      priceCurrency: "USD",
      availability: "https://schema.org/InStock",
    },
    downloadUrl: "https://github.com/Boredooms/Moodwave-CLI/releases",
    softwareRequirements: "Zero external dependencies (native Go binary)",
    fileSize: "8.2MB",
    author: {
      "@id": `${baseUrl}/#organization`,
    },
    featureList: [
      "Heuristic codebase mood analysis across 10 developer states",
      "Multi-source streaming from YouTube, Jamendo, and 30,000+ Internet Radio stations",
      "20+ real-time terminal visualizers including Fire, Matrix, Plasma, and Aurora",
      "7 switchable ANSI color themes with live TUI preview",
      "Personal curated playlist management with FIFO priority",
      "In-app binary self-updater directly from terminal home screen",
      "Zero-dependency cross-platform support for Windows, macOS, and Linux",
    ],
  };

  // 4. SiteNavigationElement Schema (Enables Google Multi-Tier Sitelinks)
  const navigationSchema = {
    "@context": "https://schema.org",
    "@type": "ItemList",
    itemListElement: [
      {
        "@type": "SiteNavigationElement",
        position: 1,
        name: "Docs",
        description: "Platform setup guides, command reference, and audio backends.",
        url: `${baseUrl}/docs`,
      },
      {
        "@type": "SiteNavigationElement",
        position: 2,
        name: "Changelog",
        description: "Release history, interactive performance metrics, and version timeline.",
        url: `${baseUrl}/changelog`,
      },
      {
        "@type": "SiteNavigationElement",
        position: 3,
        name: "Features",
        description: "Codebase mood scan, recommendation engine, and terminal visualizers.",
        url: `${baseUrl}/#features`,
      },
      {
        "@type": "SiteNavigationElement",
        position: 4,
        name: "Installation",
        description: "Zero-dependency single-line installation for Windows, macOS, and Linux.",
        url: `${baseUrl}/#install`,
      },
      {
        "@type": "SiteNavigationElement",
        position: 5,
        name: "Architecture",
        description: "Decoupled 5-layer architecture from scanner to Bubble Tea TUI.",
        url: `${baseUrl}/#architecture`,
      },
      {
        "@type": "SiteNavigationElement",
        position: 6,
        name: "Launch Waitlist",
        description: "Official launch countdown and early access waitlist.",
        url: `${baseUrl}/launch`,
      },
    ],
  };

  return (
    <>
      <script
        type="application/ld+json"
        dangerouslySetInnerHTML={{ __html: JSON.stringify(websiteSchema) }}
      />
      <script
        type="application/ld+json"
        dangerouslySetInnerHTML={{ __html: JSON.stringify(organizationSchema) }}
      />
      <script
        type="application/ld+json"
        dangerouslySetInnerHTML={{ __html: JSON.stringify(softwareAppSchema) }}
      />
      <script
        type="application/ld+json"
        dangerouslySetInnerHTML={{ __html: JSON.stringify(navigationSchema) }}
      />
    </>
  );
}

export function FAQJsonLd() {
  const faqSchema = {
    "@context": "https://schema.org",
    "@type": "FAQPage",
    mainEntity: [
      {
        "@type": "Question",
        name: "Is the CLI heavy? Will it slow down my machine?",
        acceptedAnswer: {
          "@type": "Answer",
          text: "No. Moodwave is a single compiled Go binary of ~8 MB. It has zero background processes and only runs when you invoke it. It uses your system's native audio player (mpv, ffplay) and does not transcode or cache audio locally.",
        },
      },
      {
        "@type": "Question",
        name: "Does it store music files locally?",
        acceptedAnswer: {
          "@type": "Answer",
          text: "No. Moodwave streams audio from the internet. Nothing is cached or stored beyond your configuration file and a small track metadata index used for recommendations.",
        },
      },
      {
        "@type": "Question",
        name: "Does it work on Windows?",
        acceptedAnswer: {
          "@type": "Answer",
          text: "Yes. Windows is a first-class target. The installer (install.ps1) places the binary in a local AppData path and adds it to PATH. Playback uses Windows PowerShell's built-in Media.SoundPlayer or ffplay if installed.",
        },
      },
      {
        "@type": "Question",
        name: "What music sources does it use?",
        acceptedAnswer: {
          "@type": "Answer",
          text: "Three primary sources: YouTube (via yt-dlp for stream extraction), Jamendo (Creative Commons licensed music with a public API), and Radio Browser (30,000+ community-curated internet radio stations).",
        },
      },
      {
        "@type": "Question",
        name: "Is a YouTube or music API key required?",
        acceptedAnswer: {
          "@type": "Answer",
          text: "No API keys are needed. Moodwave uses yt-dlp (auto-downloaded on first run), Jamendo's public client ID, and Radio Browser's open REST API.",
        },
      },
      {
        "@type": "Question",
        name: "Is the website the product?",
        acceptedAnswer: {
          "@type": "Answer",
          text: "No. This page is purely promotional. The product is entirely in your terminal. Installing the CLI is the only thing this page is here to help you do.",
        },
      },
    ],
  };

  return (
    <script
      type="application/ld+json"
      dangerouslySetInnerHTML={{ __html: JSON.stringify(faqSchema) }}
    />
  );
}

export function BreadcrumbsJsonLd({
  items,
}: {
  items: { name: string; url: string }[];
}) {
  const schema = {
    "@context": "https://schema.org",
    "@type": "BreadcrumbList",
    itemListElement: items.map((item, index) => ({
      "@type": "ListItem",
      position: index + 1,
      name: item.name,
      item: item.url,
    })),
  };

  return (
    <script
      type="application/ld+json"
      dangerouslySetInnerHTML={{ __html: JSON.stringify(schema) }}
    />
  );
}
