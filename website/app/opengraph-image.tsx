import { ImageResponse } from "next/og";

export const runtime = "edge";
export const alt = "Moodwave — Terminal Mood Music Companion";
export const size = {
  width: 1200,
  height: 630,
};
export const contentType = "image/png";

export default async function Image() {
  return new ImageResponse(
    (
      <div
        style={{
          height: "100%",
          width: "100%",
          display: "flex",
          flexDirection: "column",
          alignItems: "center",
          justifyContent: "center",
          backgroundColor: "#080808",
          backgroundImage:
            "radial-gradient(ellipse at center top, rgba(255,255,255,0.08) 0%, transparent 60%)",
          color: "#ffffff",
          fontFamily: "monospace",
          position: "relative",
          padding: "60px",
        }}
      >
        {/* Subtle border outline */}
        <div
          style={{
            position: "absolute",
            inset: "24px",
            border: "1px solid rgba(255, 255, 255, 0.1)",
            borderRadius: "20px",
          }}
        />

        {/* Top badge */}
        <div
          style={{
            display: "flex",
            alignItems: "center",
            gap: "10px",
            padding: "8px 18px",
            borderRadius: "9999px",
            background: "rgba(255, 255, 255, 0.05)",
            border: "1px solid rgba(255, 255, 255, 0.12)",
            fontSize: 16,
            letterSpacing: "0.2em",
            textTransform: "uppercase",
            color: "#888888",
            marginBottom: "30px",
          }}
        >
          <div
            style={{
              width: "8px",
              height: "8px",
              borderRadius: "50%",
              backgroundColor: "#10b981",
            }}
          />
          v2.0.0 is Live · Open Source CLI
        </div>

        {/* Main Logo & Title */}
        <div
          style={{
            fontSize: 72,
            fontWeight: 800,
            letterSpacing: "-0.04em",
            color: "#ffffff",
            marginBottom: "16px",
          }}
        >
          Moodwave
        </div>

        {/* Subtitle */}
        <div
          style={{
            fontSize: 26,
            color: "#888888",
            maxWidth: "780px",
            textAlign: "center",
            lineHeight: 1.4,
            marginBottom: "40px",
          }}
        >
          Codebase Mood Analysis · Terminal Visualizers · Multi-Source Audio
        </div>

        {/* Terminal Pill Command */}
        <div
          style={{
            display: "flex",
            alignItems: "center",
            background: "rgba(255, 255, 255, 0.03)",
            border: "1px solid rgba(255, 255, 255, 0.15)",
            borderRadius: "12px",
            padding: "14px 28px",
            fontSize: 22,
            color: "#dddddd",
          }}
        >
          <span style={{ color: "#666666", marginRight: "12px" }}>$</span>
          <span>moodwave play</span>
        </div>

        {/* Domain footer */}
        <div
          style={{
            position: "absolute",
            bottom: "40px",
            display: "flex",
            alignItems: "center",
            justifyContent: "space-between",
            width: "calc(100% - 120px)",
            fontSize: 16,
            color: "#555555",
          }}
        >
          <span>https://www.moodwave-cli.xyz</span>
          <span>Windows · macOS · Linux</span>
        </div>
      </div>
    ),
    {
      ...size,
    }
  );
}
