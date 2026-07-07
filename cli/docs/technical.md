# Moodwave CLI — Technical Design

## 1. Project Summary

Moodwave CLI is a **CLI-first, terminal-native, lightweight developer music companion**. The CLI is the product. The website is promotional only and should be built later.

The tool analyzes a codebase, infers a coding mood, maps that mood to music traits, fetches legal music metadata or streams from open sources, and renders ambient terminal visuals such as dot-matrix banners, waveforms, spectrum bars, and ASCII animations.

The project must be:
- cross-platform
- Windows PowerShell friendly
- shell-installable with curl / PowerShell bootstrap scripts
- low memory
- low storage
- resilient with fallbacks
- fully CLI-native
- built from scratch where practical

## 2. Product Principles

1. **CLI first, website later.**
2. **No heavy runtime dependency for the core CLI.**
3. **Prefer native binaries or compact runtimes over npm-based execution.**
4. **Stream metadata and audio state; do not store media unnecessarily.**
5. **Use open-source, legal, and documented APIs only.**
6. **The terminal experience must feel ambient, alive, and premium.**
7. **Every feature needs a graceful fallback for weak terminals.**

## 3. Core User Flow

1. User installs the CLI.
2. User runs one bootstrap command.
3. The CLI scans the current repository.
4. The mood engine infers the project vibe.
5. The recommender ranks matching tracks or stations.
6. The source adapters fetch legal music options.
7. Playback starts in terminal-controlled mode.
8. Ambient visuals react to music and mood.
9. The CLI continues to suggest updates and track changes automatically.

## 4. Architecture Overview

### 4.1 CLI Layers

**A. Bootstrap Layer**
- one-line installation
- OS detection
- self-update support
- config generation
- shell integration

**B. Core Engine**
- command parser
- repository scanner
- feature extraction
- mood inference
- music selection
- playback control
- cache management

**C. Source Adapter Layer**
- metadata providers
- stream providers
- radio fallback providers
- lyrics providers
- recommendation history providers

**D. Visual Layer**
- ASCII art scenes
- dot-matrix headers
- waveform and spectrum rendering
- theme colors
- 24 fps compatible animation loop where terminal supports it

**E. Platform Layer**
- Windows console support
- PowerShell compatibility
- batch file wrappers
- Unix shell compatibility
- TTY capability detection

## 5. Repository Scanner

The scanner should analyze:
- file types
- folder structure
- code style consistency
- comment density
- naming patterns
- recent git activity
- churn
- test presence
- dependency size
- project age
- language mix
- build tooling
- markdown/docs density

### Output examples
- calm
- focused
- experimental
- intense
- chaotic
- minimal
- polished
- debugging
- late-night
- sprint mode

## 6. Mood Engine

The mood engine should start with **cheap heuristics** and later support lightweight open-source models.

### Phase 1: Heuristics
- keyword density
- diff size
- structure entropy
- folder regularity
- TODO / FIXME presence
- test-to-source ratio
- complexity estimates

### Phase 2: Lightweight ML
- small sentence embedding models
- compact classifiers
- codebase style embeddings
- user preference vector

### Phase 3: Hybrid Decision
Final mood score = heuristic score + model score + preference score + session context.

## 7. Music Matching Logic

The recommendation engine should map mood to track traits.

### Track traits
- BPM
- energy
- tempo
- genre
- acousticness
- instrumentalness
- vocal presence
- loudness
- loopability
- concentration fit
- ambient intensity

### Matching examples
- focused -> ambient / lo-fi / instrumental
- debugging -> repetitive / stable / low-distraction
- sprint mode -> high-energy / rhythmic / faster BPM
- experimental -> textured / synthetic / unusual
- late-night -> dark ambient / soft electronic
- polished repo -> smooth / clean / elegant

### Recommendation sources
- direct music catalog APIs
- metadata APIs
- radio stations
- user history
- local presets

## 8. Open Source / Legal Source Strategy

The CLI should not depend on scraped or unstable sources. Use documented, legal APIs and open music catalogs.

### Primary metadata source candidates
- MusicBrainz
- ListenBrainz
- Jamendo
- Radio Browser
- lyrics providers such as LRCLIB or similar documented APIs

### Why use multiple sources
- metadata source and playback source can differ
- one provider may not have a track or station
- fallback keeps the CLI usable
- modular adapters avoid lock-in

### Source adapter design
Each adapter should provide:
- search
- lookup
- normalization
- rank score
- stream / preview URL when legal and available
- error and rate-limit handling
- cache hints

## 9. Open Source Repositories to Study

### kew
A terminal music player with search, playlist, spectrum, theme, and lyrics capabilities. Good reference for terminal UX, playback flow, visual design, and music-library interaction.

### musikcube
A cross-platform terminal-based audio engine, library, player, and server. Strong reference for architecture, terminal interface behavior, and native cross-platform audio design.

### Spotube
An open-source streaming platform with plugin-based source ideas. Useful as a reference for separating metadata providers, source plugins, and playback layers.

### Ponytail
Not a music engine, but a useful reference for ultra-light developer tooling and native agent workflows.

## 10. Audio Playback Strategy

The CLI should be able to:
- play tracks through a local backend
- handle station streams
- play previews when full streaming is not possible
- keep audio buffers small
- avoid storing full media in memory
- stop cleanly
- resume quickly
- switch tracks without heavy restarts

### Playback backends
The playback layer should support platform-specific options, such as:
- Windows audio path
- macOS audio path
- Linux audio path

The CLI should prefer a small and reliable backend rather than embedding a large media stack.

## 11. Terminal Visual System

The terminal should support several display modes:
- dot-matrix title mode
- ambient waveform mode
- spectrum mode
- pulse mode
- quiet mode
- minimal status mode

### Visual features
- animated bars
- wave pulses
- color shifting themes
- ASCII art frames
- song progress
- mood intensity meter
- current source badge
- lightweight visual scenes

### Rendering rules
- if terminal supports richer rendering, use it
- if not, degrade to plain ANSI text
- never block playback for visuals
- never require high GPU or huge memory
- keep frame generation cheap

## 12. Command Surface

### Core commands
- `moodwave init`
- `moodwave scan`
- `moodwave mood`
- `moodwave play`
- `moodwave pause`
- `moodwave stop`
- `moodwave next`
- `moodwave queue`
- `moodwave status`
- `moodwave config`
- `moodwave source`
- `moodwave theme`
- `moodwave visual`
- `moodwave doctor`

### Auto mode
A single scan should be enough for the CLI to infer the mood and begin suggesting tracks automatically.

## 13. Installation Strategy

### Goal
A user should be able to install from the internet with a simple command.

### Supported patterns
- `curl` bootstrap on Unix-like systems
- PowerShell bootstrap on Windows
- batch launcher for Windows CMD users
- optional package manager distribution later

### Install output
- binary or compact runtime
- config file
- cache directory
- logs directory
- update metadata

## 14. Memory / Storage Rules

The CLI must stay lean.

### Do not
- cache full songs permanently by default
- require huge models
- run a heavy always-on background process
- keep large temp files
- hold all track metadata in memory at once

### Do
- stream metadata
- cache only small indexes
- keep recent results in an LRU cache
- use on-demand downloads only when allowed
- clear transient playback state after use

## 15. Fallbacks

The tool must keep working even when some features fail.

### Fallback matrix
- no rich terminal -> plain ANSI mode
- no audio source -> radio fallback or local suggestions
- no internet -> local cache or offline mood analysis
- no model -> heuristic-only inference
- no colors -> monochrome visuals
- no waveform support -> status-only mode

## 16. Development Phases

### Phase 1: CLI MVP
- scanner
- mood inference heuristics
- source adapters for a small set of APIs
- terminal playback control
- minimal ASCII visuals
- Windows support

### Phase 2: Better Recommendations
- mood-to-track ranking
- user preferences
- history-aware selection
- better metadata normalization

### Phase 3: Visual Polish
- dot-matrix branding
- waveform and pulse animations
- theme engine
- richer terminal art

### Phase 4: Model Assist
- lightweight embeddings
- better style inference
- smart suggestions
- adaptive music changes

### Phase 5: Promotional Website
- launch page
- screenshots
- demo clips
- install docs
- roadmap

## 17. Suggested Repo Structure

```text
moodwave/
├─ apps/
│  ├─ cli/
│  └─ web/
├─ packages/
│  ├─ engine/
│  ├─ scanner/
│  ├─ mood-model/
│  ├─ recommender/
│  ├─ sources/
│  ├─ audio/
│  ├─ visuals/
│  └─ shared/
├─ docs/
│  ├─ technical.md
│  ├─ architecture.md
│  ├─ cli.md
│  ├─ visuals.md
│  ├─ sources.md
│  ├─ install.md
│  └─ website.md
├─ scripts/
├─ assets/
└─ README.md
```

## 18. Documentation Set

Create these next:
- `architecture.md`
- `cli.md`
- `sources.md`
- `visuals.md`
- `install.md`
- `modeling.md`
- `website.md`
- `README.md`
- `AGENTS.md`

## 19. Recommended Tech Direction

### CLI
Use a compact native runtime suitable for cross-platform distribution and low memory usage.

### Website
Use Next.js only for the promo site, not for the CLI core.

### Data / API layer
Keep adapters modular, testable, and source-agnostic.

## 20. Key Success Criteria

- one-command install
- low memory footprint
- stable terminal visuals
- mood inference that feels correct
- music selection that matches coding context
- Windows compatibility
- graceful fallback behavior
- clean source adapter architecture
- room to expand later without rewriting the core

## 21. Definition of Done

The CLI is complete when:
- it scans a project
- detects a mood
- selects music from a legal source
- plays or suggests music reliably
- renders polished terminal visuals
- works on Windows, Linux, and macOS
- stays lightweight
- supports future model upgrades
- the website can explain and promote it clearly

