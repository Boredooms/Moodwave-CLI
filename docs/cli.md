# Moodwave CLI — CLI Design

## 1. What the CLI is

Moodwave CLI is the **core product**. It is a terminal-native, cross-platform music companion for developers that scans a codebase, infers the coding mood, recommends music, and renders a rich ambient terminal experience.

The website is secondary and promotional only. The CLI is what users actually run every day.

## 2. CLI Mission

The CLI should:
- understand the structure and style of a codebase
- infer the current coding mood or session energy
- recommend the right music from open, legal sources
- stream or play music in a lightweight way
- react to code changes and evolve the soundtrack
- provide a premium terminal experience with animations and visual feedback
- stay fast, compact, and easy to install on Windows, macOS, and Linux

## 3. CLI Design Goals

### 3.1 Lightweight
The CLI should use minimal memory and startup time.

### 3.2 Native feeling
The interface should feel like a serious terminal tool, not a web app inside a terminal.

### 3.3 Automatic
It should do useful work after one scan and continue suggesting updates with minimal user effort.

### 3.4 Adaptive
Music and visuals should change with the codebase mood, coding patterns, and session behavior.

### 3.5 Beautiful
The terminal experience should feel rich, atmospheric, and highly polished.

### 3.6 Portable
It should work through shell installers, batch wrappers, and platform-friendly binaries or compact runtimes.

## 4. Core CLI Behavior

The CLI should support two operating styles:

### Interactive mode
Used when the user wants to explore, select tracks, change theme, or review mood output.

### Auto mode
Used when the CLI scans once and then continues running with automatic updates, track suggestions, and ambient visuals.

## 5. Main Responsibilities

The CLI must handle:
- repository scan
- feature extraction
- mood inference
- music recommendation
- source lookup
- stream or preview playback
- queue management
- terminal animation rendering
- source fallbacks
- config and preferences
- status reporting
- error recovery
- cross-platform launch behavior

## 6. Command Philosophy

The command set should be small but powerful.

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
- `moodwave theme`
- `moodwave visual`
- `moodwave source`
- `moodwave doctor`

### Auto-suggested behavior
After a scan, the CLI should suggest the best next action rather than forcing the user to keep typing commands.

## 7. Recommendation Engine

The recommendation engine is the heart of the CLI.

It should:
- map codebase signals to a mood label
- map mood labels to music traits
- rank candidate songs or stations
- avoid repetition
- support smooth transitions
- allow user override at any moment

### Inputs to the recommendation engine
- repo structure
- file counts
- language mix
- naming style
- comment density
- git churn
- TODO / FIXME density
- test coverage signals
- recent activity rhythm
- user preferences
- current playback history

### Music traits used for ranking
- BPM
- energy
- tempo
- instrumentalness
- mood tags
- vocal presence
- ambient intensity
- loopability
- concentration fit
- transition compatibility

## 8. Mood Prediction Layer

The mood layer should not depend on a large model by default.

### Stage 1: Heuristics
Fast, explainable, and always available.
- code density
- folder regularity
- complexity indicators
- diff size
- docs-to-code ratio
- error signal density

### Stage 2: Lightweight models
Optional compact models or embeddings to improve nuance.
- small transformer embeddings
- compact classifier models
- session-level style vectors
- preference blending

### Stage 3: Hybrid inference
Final decision combines:
- heuristic score
- model score
- confidence
- playback history
- user override

### Example mood outputs
- focused
- calm
- intense
- chaotic
- experimental
- late-night
- debugging
- sprint mode
- clean
- minimal

## 9. Music Behavior Logic

The music should adapt gradually.

### Transition rules
- lo-fi should stay near lo-fi when the vibe remains stable
- intense coding should not jump to something jarring unless requested
- if the repo becomes more chaotic, the soundtrack should move toward more energetic or textured options
- if the session becomes calm, the soundtrack should soften
- repeated tracks should be penalized

### Change triggers
- new repository scan
- major file changes
- branch switch
- commit activity spikes
- playback duration thresholds
- manual user override

## 10. Open Source Integration Model

The CLI should use open and documented source integrations only.

### Primary integration layers
- metadata fetchers
- station search
- stream or preview resolvers
- lyrics fetchers
- local playback backends
- terminal visual renderers

### Important rule
The CLI should be source-agnostic. Music providers are adapters, not hardcoded assumptions.

## 11. Terminal Experience

The terminal should be a first-class design surface.

### Visual features
- dot matrix header
- animated ASCII waves
- spectrum bars
- soft pulse motion
- track progress line
- mood meter
- source badge
- session badge
- theme tinting
- slow ambient transitions

### Quality targets
- 24 fps where supported
- graceful fallback to lower refresh rates
- ANSI color support
- monochrome fallback
- Unicode fallback
- no blocking the playback loop

### Look and feel
The CLI should feel closer to a premium developer terminal experience than a basic player.

## 12. External Libraries

The CLI may use lightweight external libraries where they help with:
- terminal rendering
- colors and ANSI styling
- keyboard input
- audio backend integration
- HTTP requests
- config parsing
- file watching
- rate limiting
- process control
- animation timing

### Rule
Use external libraries only when they improve reliability or performance. Do not use heavy dependencies without a clear reason.

## 13. Storage and Memory Strategy

The CLI must remain lean.

### Do not
- download songs permanently by default
- load large media into memory
- keep huge caches
- keep a heavy model resident all the time
- depend on large always-running services

### Do
- stream metadata
- cache only small state and ranking data
- keep recent history in a small local cache
- unload playback resources when idle
- use compact models only when needed

## 14. CLI Runtime Modes

### Mode A: Scan only
Analyze the repo and show the mood without playback.

### Mode B: Scan and play
Analyze, recommend, and begin playback.

### Mode C: Ambient watch
Monitor code changes and adapt the soundtrack over time.

### Mode D: Quiet mode
Show mood analysis and recommendations but do not play audio.

## 15. Install and Launch Model

The CLI should be installable from the internet with a simple command.

### Desired launch styles
- `curl | sh` on Unix-like systems
- PowerShell bootstrap on Windows
- batch-friendly launcher for CMD users
- single-binary or compact-runtime distribution when possible

### First-run behavior
- create config directory
- detect terminal capabilities
- detect OS
- test source connectivity
- print a concise welcome screen
- offer scan/play/quiet setup

## 16. Error Handling

The CLI should fail gracefully.

### Examples
- no internet -> use cache or offline mode
- source unavailable -> switch provider
- terminal unsupported -> plain text mode
- model unavailable -> heuristics only
- playback failure -> suggest fallback track or station
- bad config -> safe default values

## 17. Accessibility and Usability

The CLI should be usable by everyday developers, not only power users.

### Usability rules
- command names should be predictable
- help output should be clear
- defaults should be safe
- visuals should be optional
- color should never be the only signal
- keyboard shortcuts should be discoverable

## 18. Website Relationship

The website is not the product core.

It should:
- explain the concept
- show screenshots or demos
- document install steps
- describe the CLI commands
- present releases and updates
- attract users to install the CLI

The website should not contain the actual mood engine or music logic.

## 19. Build Order

### Step 1: CLI engine
Build the scanner, mood engine, recommender, and adapter system first.

### Step 2: Playback and visuals
Add source integration, playback control, and terminal rendering.

### Step 3: Auto mode
Make the CLI scan once and continue adapting to code changes.

### Step 4: Refine model layer
Add lightweight model support and better inference.

### Step 5: Promotional website
Build the website only after the CLI feels polished.

## 20. Definition of a Great CLI

A great Moodwave CLI:
- works immediately after install
- feels alive in the terminal
- understands the codebase mood
- picks matching music well
- handles changes smoothly
- stays lightweight and reliable
- supports Windows, macOS, and Linux
- looks premium without becoming heavy
- makes the website feel like a showcase, not the main product

