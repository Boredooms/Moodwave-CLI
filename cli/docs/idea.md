# Moodwave CLI

A CLI-first, terminal-native music companion for developers.

## Vision

Moodwave CLI is a lightweight, cross-platform developer tool that analyzes codebases, infers the coding mood and project energy, fetches matching music metadata from legal APIs, and plays ambient soundtracks inside the terminal experience. The website is only promotional; the CLI is the core product.

The project should feel like a living terminal instrument: subtle dot-matrix visuals, ambient color waves, ASCII animation, developer-focused mood modes, and music-reactive UI feedback, all built from scratch with a strong emphasis on performance, low memory use, and Windows PowerShell compatibility.

## Product Principle

The CLI comes first.

The website is secondary and exists only to:
- introduce the project
- explain features
- show demos
- provide install instructions
- document commands and themes
- link to downloads and releases

## What the CLI does

1. Scans a folder or repository.
2. Extracts codebase signals.
3. Infers the current mood of the project or developer session.
4. Selects a matching music profile.
5. Queries music sources through legal APIs.
6. Streams audio through the terminal workflow or a local playback bridge.
7. Displays ambient terminal visuals while the music plays.
8. Adapts the visuals to the track and mood.
9. Stays lightweight and responsive.

## Core Experience

The CLI should feel like:
- a terminal mood engine
- a developer soundscape assistant
- a live coding atmosphere generator
- a clean, modern, minimal but expressive terminal tool

The output should not feel generic or noisy. It should feel carefully composed.

## Mood Detection

The mood engine should analyze the codebase using lightweight heuristics first, then optional models.

### Signals to analyze
- language mix
- folder structure
- file counts
- function and class density
- comment density
- naming style
- formatting consistency
- code complexity
- recent git activity
- churn rate
- TODO / FIXME density
- test coverage presence
- dependency size
- build system type
- project age and scale

### Example mood outputs
- focused
- calm
- intense
- chaotic
- experimental
- minimal
- polished
- late-night
- sprint mode
- debugging mode

## Music Matching Logic

The recommendation engine should map mood to track traits.

### Track traits
- BPM
- energy
- genre
- mood tags
- loudness
- tempo
- instrumentation
- ambient level
- vocal presence
- loopability
- concentration fit

### Matching idea
- focused code -> calm instrumental, low distraction
- debugging -> steady, repetitive, low-energy tracks
- sprint mode -> fast, driving, rhythmic tracks
- experimental -> creative, textured, electronic or lo-fi sounds
- late-night -> dark ambient or soft atmospheric music
- clean codebase -> smooth, polished, minimal audio mood

## Terminal Visuals

The terminal should include optional animated feedback such as:
- dot matrix style banners
- music-reactive waveforms
- ambient pulse lines
- ASCII equalizers
- color shifting bars
- low-key frame animations
- track progress meters
- mood badges
- live status indicators
- subtle scene transitions

### Visual directions
- 24 fps animation where supported
- graceful fallback for low-refresh terminals
- ANSI color support
- no heavy rendering libraries in the core runtime
- Windows PowerShell friendly
- clean fallback mode for plain terminals

## Playback Model

The system should avoid unnecessary memory pressure.

### Preferred behavior
- stream metadata and playback state
- do not cache huge media files in memory
- do not download music permanently for normal playback
- keep local buffers small
- keep background processes minimal
- shut down playback cleanly
- resume or switch tracks without restarting the whole engine

## Cross-Platform Distribution

The CLI should be installable from the internet using simple commands and small bootstrap scripts.

### Target platforms
- Windows PowerShell
- Windows CMD with batch wrappers
- macOS terminal
- Linux shell

### Install philosophy
- one-line bootstrap if possible
- curl-friendly installation
- platform-specific launchers only when needed
- no npm dependency for the core CLI runtime
- minimal setup burden

## Proposed Runtime Direction

The CLI core should be built in a lightweight runtime suitable for distribution without npm.

Good options:
- Go
- Rust
- Python with a compact packaging strategy

The website can use Next.js later, but the CLI itself should remain native, small, and fast.

## Suggested Architecture

### 1. CLI Core
Responsible for:
- parsing commands
- loading config
- running scans
- generating mood scores
- selecting tracks
- starting playback
- driving terminal visuals

### 2. Mood Engine
Responsible for:
- analyzing codebase signals
- computing a mood profile
- emitting confidence scores
- recommending a sound profile

### 3. Recommendation Engine
Responsible for:
- mapping mood to music tags
- ranking available tracks
- applying user preferences
- caching track metadata only

### 4. Source Connectors
Responsible for:
- querying supported music APIs
- normalizing results
- handling rate limits and errors
- filtering legal/usable sources

### 5. Visual Engine
Responsible for:
- ASCII art
- waves
- animated equalizers
- color transitions
- fallback rendering modes

### 6. Web Promo Site
Responsible for:
- project landing page
- product story
- feature showcase
- install instructions
- screenshots and demos
- release notes

## CLI Commands

Proposed commands:
- `moodwave init`
- `moodwave scan`
- `moodwave mood`
- `moodwave play`
- `moodwave stop`
- `moodwave pause`
- `moodwave next`
- `moodwave queue`
- `moodwave status`
- `moodwave config`
- `moodwave theme`
- `moodwave visuals`
- `moodwave sources`
- `moodwave doctor`

## Example Flow

1. User opens a codebase.
2. Runs `moodwave scan`.
3. The tool analyzes the repo.
4. It prints a mood result such as `focused / calm / 82% confidence`.
5. It fetches suitable music metadata.
6. It starts playback.
7. Terminal waves and ambient ASCII visuals appear.
8. User keeps coding while the soundtrack adapts.

## Design Goals

- fast startup
- tiny memory footprint
- stable terminal rendering
- elegant visuals
- meaningful mood inference
- useful music selection
- cross-platform reliability
- easy installation
- no unnecessary bloat

## Non-Goals for the CLI

- no heavy web app logic inside the CLI
- no large always-running background service unless required
- no massive in-memory audio handling
- no dependency on npm for the core tool
- no feature creep into a full media player

## Website Phase

The website should come after the CLI is solid.

### Website purpose
- explain the idea
- show live screenshots and GIFs
- document installation
- present feature demos
- collect feedback
- route users to downloads

### Website style
- premium promotional landing page
- dark, atmospheric, terminal-inspired look
- subtle motion
- strong typography
- minimal but striking sections

## Repo Structure

```text
moodwave/
├─ apps/
│  ├─ cli/
│  └─ web/
├─ packages/
│  ├─ engine/
│  ├─ mood-model/
│  ├─ recommender/
│  ├─ sources/
│  ├─ visuals/
│  └─ shared/
├─ docs/
│  ├─ idea.md
│  ├─ architecture.md
│  ├─ cli.md
│  ├─ website.md
│  ├─ commands.md
│  └─ release.md
├─ scripts/
├─ assets/
└─ README.md
```

## Important docs to create next

- `architecture.md`
- `cli.md`
- `website.md`
- `commands.md`
- `sources.md`
- `mood-model.md`
- `visual-system.md`
- `install.md`
- `AGENTS.md`
- `README.md`

## MVP Scope

### CLI MVP
- scan repository
- infer mood using rules and lightweight scoring
- choose tracks from one or two music sources
- basic playback state
- one terminal visual mode
- Windows PowerShell support
- simple install command

### Website MVP
- hero section
- product story
- CLI screenshots
- install steps
- feature list
- roadmap

## Long-Term Vision

A developer terminal that feels alive.

A CLI that understands the codebase, sets the mood, and turns coding sessions into ambient experiences without becoming heavy or distracting.

The product should be remembered as a serious terminal tool with a playful heart.

