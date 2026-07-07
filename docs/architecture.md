# Moodwave CLI — Architecture

## 1. Architecture Goal

Moodwave CLI is a **CLI-first, terminal-native, cross-platform music companion for developers**. The core product must work without the promotional website. The website is a later, secondary layer for marketing, documentation, and downloads.

The architecture must support:
- fast startup
- low memory usage
- low storage usage
- Windows PowerShell compatibility
- Unix shell compatibility
- legal/open music source integration
- codebase mood analysis
- adaptive music recommendations
- ambient terminal animations
- graceful fallbacks on weak terminals
- future model upgrades without rewriting the core

## 2. High-Level System View

The system is divided into five major layers:

1. **Bootstrap Layer**
   - installation
   - platform detection
   - update handling
   - config generation

2. **CLI Control Layer**
   - command parsing
   - session orchestration
   - interactive prompts
   - status output

3. **Core Intelligence Layer**
   - repository scanning
   - feature extraction
   - mood inference
   - recommendation ranking
   - playback decision-making

4. **Source and Playback Layer**
   - metadata providers
   - music search providers
   - station providers
   - lyrics/preview providers
   - playback adapter

5. **Terminal Experience Layer**
   - dot-matrix visuals
   - spectrum/waveform rendering
   - ASCII animations
   - theme switching
   - reduced-motion and fallback modes

## 3. Design Principles

### 3.1 CLI First
The CLI is the product. All important logic must live in the CLI or shared engine packages.

### 3.2 Minimal Runtime
Avoid heavy dependencies, large runtimes, or always-on background services.

### 3.3 Modular Sources
Music providers must be independent adapters so the system can swap sources without changing the core engine.

### 3.4 Adaptive Intelligence
Mood detection should begin with lightweight heuristics and optionally use compact models later.

### 3.5 Terminal-First UX
Visuals must be designed for terminals first, not copied from web UI patterns.

### 3.6 Fallback Safety
Every advanced feature must degrade cleanly when terminal capability, network access, or provider availability is limited.

## 4. Recommended Repository Layout

```text
moodwave/
├─ apps/
│  ├─ cli/
│  │  ├─ src/
│  │  ├─ commands/
│  │  ├─ renderers/
│  │  ├─ prompts/
│  │  └─ main
│  └─ web/
│     ├─ pages/
│     ├─ components/
│     └─ public/
├─ packages/
│  ├─ core/
│  ├─ scanner/
│  ├─ mood-engine/
│  ├─ recommender/
│  ├─ sources/
│  ├─ playback/
│  ├─ visuals/
│  ├─ config/
│  ├─ cache/
│  └─ shared/
├─ docs/
│  ├─ architecture.md
│  ├─ technical.md
│  ├─ cli.md
│  ├─ install.md
│  ├─ sources.md
│  ├─ visuals.md
│  ├─ model.md
│  └─ website.md
├─ scripts/
├─ assets/
├─ tests/
└─ README.md
```

## 5. Package Responsibilities

### 5.1 `packages/core`
The orchestration backbone.
- session lifecycle
- event bus
- command routing
- configuration loading
- provider selection
- error handling policies

### 5.2 `packages/scanner`
Repository inspection.
- file system traversal
- ignore rules
- language detection
- git metadata reading
- code feature extraction
- codebase size and structure analysis

### 5.3 `packages/mood-engine`
Mood inference logic.
- heuristic scoring
- optional model inference
- mood confidence calculation
- mood state transitions
- developer style estimation

### 5.4 `packages/recommender`
Music matching.
- map mood to track traits
- rank candidate tracks
- apply user preferences
- track diversity and repetition control
- fallback candidate selection

### 5.5 `packages/sources`
Provider adapters.
- metadata APIs
- search APIs
- station APIs
- lyrics/preview APIs
- source health checks
- rate limit handling

### 5.6 `packages/playback`
Audio control.
- stream start/stop/pause/resume
- local backend integration
- buffer control
- playback state sync
- stream fallback behavior

### 5.7 `packages/visuals`
Terminal rendering.
- dot-matrix banners
- waveforms
- spectrum bars
- progress lines
- theme colors
- animation timing

### 5.8 `packages/config`
Configuration.
- defaults
- file-based config
- environment overrides
- per-project config
- provider keys
- visual settings

### 5.9 `packages/cache`
Caching.
- recent track metadata
- provider responses
- mood inference hints
- session history
- LRU eviction

### 5.10 `packages/shared`
Common utilities.
- types
- logging helpers
- validation
- file path utilities
- terminal capability detection

## 6. CLI Flow

### 6.1 Startup Flow
1. Bootstrap script starts the CLI.
2. CLI reads config and environment.
3. Terminal capability is detected.
4. Provider availability is checked.
5. Session state is loaded.
6. Default mode is restored.
7. Command execution begins.

### 6.2 Scan Flow
1. User runs a scan or auto-scan starts.
2. Scanner reads repository structure.
3. Feature extractor derives codebase signals.
4. Mood engine builds a mood profile.
5. Recommender converts mood into music traits.
6. Source adapters search for legal options.
7. Playback starts or suggestions are presented.
8. Visual engine renders the ambient terminal state.

### 6.3 Ongoing Auto Mode
The CLI should not require constant manual commands.
- One run can start a session.
- The system can watch for code changes.
- Mood can update when meaningful repo changes occur.
- Track changes can be suggested automatically.
- The user should still be able to override at any time.

## 7. Mood Inference Architecture

### 7.1 Input Signals
The mood engine should consider:
- language distribution
- file count and nesting
- code density
- naming consistency
- comment density
- test coverage
- build scripts
- TODO / FIXME patterns
- recent diff size
- commit rhythm
- churn intensity
- dependency weight
- docs-to-code ratio
- error / debug signal density

### 7.2 Inference Stages

#### Stage A: Rules
Lightweight scoring from direct code signals.
- fast
- explainable
- no heavy memory cost
- ideal for startup and offline mode

#### Stage B: Compact Models
Optional embeddings or classifier models.
- improve nuance
- better semantic interpretation
- still lightweight enough for local use

#### Stage C: Hybrid Decision
Final mood is produced from:
- heuristic score
- model confidence
- current session activity
- user preferences
- recent playback history

### 7.3 Mood Outputs
Examples:
- focused
- calm
- intense
- chaotic
- experimental
- late-night
- sprint mode
- debugging
- clean
- minimal

### 7.4 Explainability
The CLI should show why a mood was chosen.
- top contributing signals
- confidence score
- detected patterns
- recommended music style

## 8. Music Recommendation Architecture

### 8.1 Matching Dimensions
The recommender should rank songs or stations using:
- BPM
- energy
- tempo
- ambient level
- instrumentalness
- vocals
- genre
- loopability
- concentration fit
- user taste history

### 8.2 Ranking Strategy
1. Match current mood.
2. Apply user preferences.
3. Penalize repetition.
4. Prefer legally streamable candidates.
5. Favor lightweight playback options.
6. Keep fallbacks ready.

### 8.3 Music Selection Modes
- **same mood**: keep current vibe stable
- **adjacent mood**: gradual transition
- **contrast mode**: shift energy based on session state
- **user-picked mode**: user overrides the algorithm

### 8.4 Music Transition Rules
The system should support slow transitions:
- lo-fi stays lo-fi
- ambient stays ambient
- energetic tracks fade gradually into similar energetic tracks
- abrupt switches should be avoided unless the user requests them

## 9. Source Adapter Architecture

### 9.1 Adapter Contract
Each source adapter should implement:
- `search(query)`
- `resolve(id)`
- `normalize(raw)`
- `healthCheck()`
- `rankCandidate(candidate, context)`
- `streamOrPreview(candidate)`
- `fallback(candidate)`

### 9.2 Provider Types
- metadata providers
- preview providers
- station providers
- lyrics providers
- history/recommendation providers

### 9.3 Source Priorities
Recommended source strategy:
1. metadata provider
2. preview/stream provider
3. station fallback
4. offline cache fallback
5. local silence / ambience mode

### 9.4 Legal and Stability Rules
- use documented APIs only
- respect rate limits
- cache only allowed metadata
- avoid scraping as a primary strategy
- use open-source or open APIs where possible

## 10. Playback Architecture

### 10.1 Playback Controller
The playback controller should manage:
- play
- pause
- resume
- stop
- next
- previous
- volume
- mute
- queue
- fade transitions

### 10.2 Backend Layer
The backend should remain replaceable.
- one backend on Windows
- one backend on macOS
- one backend on Linux
- if unavailable, use a fallback mode

### 10.3 Memory Rules
- never load unnecessary media into memory
- stream when possible
- use small buffers
- cache only metadata and small previews
- free playback resources when idle

## 11. Visual Rendering Architecture

### 11.1 Visual Modes
- dot-matrix title mode
- waveform mode
- spectrum mode
- pulse mode
- quiet mode
- minimal status mode
- ambient scene mode

### 11.2 Animation System
The renderer should be frame-driven.
- target 24 fps when supported
- degrade to slower frame rates automatically
- animation logic should stay cheap
- rendering should never block audio

### 11.3 Visual Elements
- ASCII art banners
- terminal waves
- equalizer bars
- color changes tied to mood
- progress meter
- song title display
- source badge
- mood badge

### 11.4 Styling Strategy
Visuals can use:
- hand-written terminal art
- ANSI color sequences
- optional external rendering helpers if lightweight and cross-platform

## 12. Platform Architecture

### 12.1 Windows Support
Must support:
- PowerShell
- CMD batch launchers
- Windows console behavior
- path handling
- native executable distribution

### 12.2 Unix Support
Must support:
- bash
- zsh
- fish-friendly patterns where possible
- curl bootstrap install
- executable permissions

### 12.3 Terminal Capability Detection
The CLI should detect:
- color support
- unicode support
- animation support
- TTY availability
- terminal width/height
- reduced motion preference

## 13. Configuration Architecture

### 13.1 Config Sources
Priority order:
1. command-line flags
2. environment variables
3. per-project config
4. global config
5. built-in defaults

### 13.2 Config Domains
- sources
- visuals
- playback
- mood sensitivity
- auto scan behavior
- cache settings
- network settings
- key bindings
- user preferences

## 14. Auto-Scan and Watch Mode

The CLI should support periodic or event-based updates.

### Watch inputs
- file changes
- git changes
- branch switches
- commit activity
- user-triggered rescan

### Watch outputs
- mood recalculation
- track transition suggestion
- visual scene update
- confidence adjustment

## 15. Error Handling Architecture

Every module must have:
- structured errors
- user-facing friendly messages
- debug mode details
- fallback path where possible

### Examples
- source unavailable -> switch provider
- internet missing -> use cache or offline mode
- terminal unsupported -> plain text output
- playback failure -> suggest alternative source
- model unavailable -> heuristic-only inference

## 16. Testing Architecture

### Unit tests
- scanner
- mood engine
- recommender
- source normalization
- config parsing

### Integration tests
- CLI startup
- source adapter chain
- playback state transitions
- terminal renderer

### End-to-end tests
- install
- scan
- infer mood
- fetch candidate track
- start playback
- render visual output

## 17. Performance Targets

- fast startup
- minimal resident memory
- low idle CPU
- small disk footprint
- quick scan path for large repos
- non-blocking rendering loop
- efficient source requests

## 18. Security and Safety

- secrets must not be logged
- API keys must be stored safely
- config should support safe credential loading
- avoid arbitrary remote code execution from source data
- sanitize filenames and paths
- validate provider responses

## 19. Phase Plan

### Phase 1 — CLI Foundation
- repo scanner
- config
- command parser
- mood heuristics
- basic source adapters
- simple playback
- minimal visuals

### Phase 2 — Recommendation System
- better ranking
- music trait matching
- user preference learning
- repeat avoidance

### Phase 3 — Terminal Experience
- richer animations
- dot-matrix art
- waveform visualization
- theme engine

### Phase 4 — Model Upgrade
- small models for better inference
- semantic mood understanding
- better track transitions

### Phase 5 — Website
- promotional homepage
- docs
- screenshots
- install guide

## 20. Future Website Role

The website should not replace the CLI.
It should only explain it, show it, and help users install it.

## 21. Definition of a Strong Architecture

The architecture is successful when:
- the CLI works independently
- users can install it easily
- mood detection feels intelligent
- music suggestions feel relevant
- playback remains lightweight
- visuals feel unique and polished
- every platform has a fallback path
- future changes do not require a rewrite

