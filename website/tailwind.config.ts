import type { Config } from "tailwindcss";

const config: Config = {
  content: [
    "./pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./components/**/*.{js,ts,jsx,tsx,mdx}",
    "./app/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      colors: {
        // Monochrome design system
        bg: "#080808",
        surface: "#111111",
        "surface-2": "#181818",
        "surface-3": "#222222",
        border: "rgba(255,255,255,0.07)",
        "border-strong": "rgba(255,255,255,0.12)",
        primary: "#FFFFFF",
        secondary: "#888888",
        tertiary: "#444444",
        accent: "#E8E8E8",
        "code-bg": "#0D0D0D",
        "code-border": "rgba(255,255,255,0.08)",
      },
      fontFamily: {
        mono: ["var(--font-geist-mono)", "JetBrains Mono", "Fira Code", "monospace"],
        sans: ["var(--font-inter)", "Inter", "system-ui", "sans-serif"],
      },
      fontSize: {
        "display-xl": ["clamp(3rem, 7vw, 6rem)", { lineHeight: "1.0", letterSpacing: "-0.04em" }],
        "display-lg": ["clamp(2rem, 4.5vw, 3.75rem)", { lineHeight: "1.05", letterSpacing: "-0.03em" }],
        "display-md": ["clamp(1.5rem, 3vw, 2.25rem)", { lineHeight: "1.1", letterSpacing: "-0.025em" }],
        "body-lg": ["1.125rem", { lineHeight: "1.7", letterSpacing: "-0.01em" }],
        "body-md": ["1rem", { lineHeight: "1.65" }],
        "body-sm": ["0.875rem", { lineHeight: "1.6" }],
        "code-md": ["0.875rem", { lineHeight: "1.7", letterSpacing: "0.01em" }],
        "code-sm": ["0.8125rem", { lineHeight: "1.6" }],
      },
      spacing: {
        section: "clamp(5rem, 10vh, 8rem)",
        "section-sm": "clamp(3rem, 6vh, 5rem)",
      },
      animation: {
        "cursor-blink": "cursor-blink 1s step-end infinite",
        "pulse-slow": "pulse 3s ease-in-out infinite",
        "fade-up": "fade-up 0.6s ease forwards",
      },
      keyframes: {
        "cursor-blink": {
          "0%, 100%": { opacity: "1" },
          "50%": { opacity: "0" },
        },
        "fade-up": {
          from: { opacity: "0", transform: "translateY(20px)" },
          to: { opacity: "1", transform: "translateY(0)" },
        },
      },
      backgroundImage: {
        "noise": "url(\"data:image/svg+xml,%3Csvg viewBox='0 0 200 200' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='n'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.9' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23n)' opacity='0.03'/%3E%3C/svg%3E\")",
      },
    },
  },
  plugins: [],
};

export default config;
