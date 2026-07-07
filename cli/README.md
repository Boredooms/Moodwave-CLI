# Moodwave CLI

**A terminal-native developer mood music companion.**

[![License: MIT](https://img.shields.io/badge/License-MIT-white.svg)](LICENSE)
[![Platform: Windows | macOS | Linux](https://img.shields.io/badge/Platform-Windows%20%7C%20macOS%20%7C%20Linux-white.svg)]()
[![Runtime: Go](https://img.shields.io/badge/Runtime-Go-00ADD8.svg)](https://go.dev)

---

Moodwave scans your codebase, infers the current coding mood, and plays matching music directly in your terminal — with ambient ASCII animations, waveforms, and reactive visuals.

**The CLI is the product. The website is promotional only.**

---

## Install

### macOS / Linux

```sh
curl -sSL https://raw.githubusercontent.com/moodwave/moodwave/main/scripts/install.sh | sh
```

### Windows (PowerShell)

```powershell
irm https://raw.githubusercontent.com/moodwave/moodwave/main/scripts/install.ps1 | iex
```

### Build from source

```sh
git clone https://github.com/moodwave/moodwave
cd moodwave
make install
```

---

## Quick Start

```sh
# Initialize config
moodwave init

# Scan repo and detect mood
moodwave scan

# Check system health
moodwave doctor

# Play music matched to detected mood
moodwave play

# See all commands
moodwave --help
```

---

## Commands

| Command            | Description                                 |
|--------------------|---------------------------------------------|
| `moodwave init`    | Initialize config and create directories    |
| `moodwave scan`    | Scan repo, extract signals, infer mood      |
| `moodwave mood`    | Show current mood profile and explanation   |
| `moodwave play`    | Start playback matched to detected mood     |
| `moodwave pause`   | Pause playback                              |
| `moodwave stop`    | Stop playback and release audio resources   |
| `moodwave next`    | Skip to next recommended track             |
| `moodwave queue`   | Show current track queue                    |
| `moodwave status`  | Show current CLI and playback status        |
| `moodwave config`  | View or edit configuration                  |
| `moodwave theme`   | Switch visual theme                         |
| `moodwave visual`  | Switch visual mode                          |
| `moodwave source`  | View or switch music source                 |
| `moodwave doctor`  | Run diagnostics on all subsystems           |

---

## Architecture

```
moodwave/
├─ cmd/moodwave/        CLI entry point and command dispatcher
├─ internal/
│  ├─ config/           Configuration loading and management
│  ├─ platform/         OS and terminal capability detection
│  ├─ scanner/          Repository analysis and signal extraction
│  ├─ mood/             Heuristic mood inference engine
│  ├─ recommender/      Mood-to-track matching and ranking
│  ├─ sources/          Music API adapters (MusicBrainz, Radio Browser, etc.)
│  ├─ playback/         Audio playback controller
│  ├─ visuals/          Terminal UI renderer and animation engine
│  └─ cache/            LRU metadata cache
├─ docs/                Project documentation
├─ scripts/             Install and build scripts
├─ assets/              Design tokens and static assets
├─ tests/               Integration and end-to-end tests
└─ web/                 Promotional website (Phase 3)
```

---

## Music Sources

Moodwave uses public, legal, documented APIs only. No API keys required for the default experience.

| Source         | Type                  | Auth Required |
|----------------|-----------------------|---------------|
| Radio Browser  | Internet radio        | No            |
| MusicBrainz    | Metadata              | No            |
| ListenBrainz   | History/recommendations | Optional (user token) |
| LRCLIB         | Synchronized lyrics   | No            |
| Jamendo        | CC music catalog      | Optional (client_id) |

---

## Mood System

Detected moods and their music mapping:

| Mood         | BPM Range | Energy | Genre Hint                  |
|--------------|-----------|--------|------------------------------|
| focused      | 70-90     | low    | Lo-fi, ambient instrumental  |
| calm         | 60-80     | low    | Acoustic, soft electronic    |
| intense      | 120-150   | high   | Electronic, driving          |
| chaotic      | 130-160   | high   | Fast, textured               |
| experimental | 80-120    | medium | Generative, avant-garde      |
| late-night   | 60-80     | low    | Dark ambient, soft            |
| debugging    | 70-90     | low    | Repetitive, stable            |
| sprint       | 120-140   | high   | Rhythmic, motivational        |
| minimal      | 60-70     | low    | Drone, minimal                |
| polished     | 90-110    | medium | Smooth jazz, clean electronic |

---

## Terminal Visual Modes

- `wave` — Animated sine waveform responding to music energy
- `spectrum` — Equalizer bars responding to track energy
- `pulse` — Ambient pulse rings
- `minimal` — Status-only compact mode
- `quiet` — No animation, text only

---

## Design Principles

1. **CLI first** — the website showcases the CLI, not the other way around
2. **Native binaries** — no runtime required beyond the executable
3. **Legal sources** — only public APIs with documented terms
4. **Graceful fallback** — works in any terminal, even without color or animation
5. **Low memory** — stream metadata, never cache full audio
6. **Modular adapters** — swap sources without touching core logic

---

## Development

```sh
# Run locally
go run ./cmd/moodwave

# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
make test

# Cross-compile
make cross
```

---

## Documentation

- [Idea & Vision](docs/idea.md)
- [Technical Design](docs/technical.md)
- [Architecture](docs/architecture.md)
- [CLI Reference](docs/cli.md)
- [CLI Design System](docs/cli_design.md)
- [Music Sources](docs/sources.md)
- [Command Reference](docs/commands.md)
- [Installation Guide](docs/install.md)
- [Project Status](docs/STATUS.md)
- [Website Design](docs/website_design.md)

---

## License

MIT © Moodwave Contributors

---

*The terminal is the UI. The music is the atmosphere. The mood is the signal.*
