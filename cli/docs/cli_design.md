# Moodwave CLI — CLI Design

## 1. Purpose of the CLI Design

This document defines the **interactive terminal design system** for Moodwave CLI.

The CLI is the primary product. The design must feel like a handcrafted terminal instrument: highly interactive, visually rich, fast, accessible, and polished enough for everyday use by developers on Windows, macOS, and Linux.

The website is secondary. The CLI must be designed first.

## 2. Core Design Goal

The CLI should feel like a premium developer tool with:
- elegant terminal visuals
- responsive controls
- smooth animated waveforms
- reactive music visualization
- clean ASCII art scenes
- strong keyboard interaction
- low memory usage
- no unnecessary clutter
- graceful fallback when terminal features are limited

## 3. Design Direction

### Visual style
- dark terminal-first aesthetic
- monochrome foundation
- subtle grayscale hierarchy
- white-on-black clarity
- precise borders and spacing
- ambient motion rather than noisy motion
- minimal but expressive UI surfaces

### Experience style
- calm when idle
- alive when music is playing
- responsive to code changes
- easy to understand immediately
- powerful without feeling complex

## 4. Skills to Explicitly Use in Antigravity

The CLI design should deliberately use the installed design skills during planning and implementation.

### Required skills for this CLI design
- `frontend-design` for overall aesthetic direction and avoiding generic templates
- `ui-design` for screen composition, hierarchy, and visual clarity
- `interaction-design` for play/pause flows, keyboard states, and transitions
- `micro-interaction-spec` for button feedback, hover states, and subtle motion
- `motion-system` for timing, easing, and animation consistency
- `layout-grid` for structured terminal panel layout
- `spacing-system` for clean visual rhythm
- `typography-scale` for title, status, and command hierarchy
- `visual-hierarchy` for strong focal points and readable scanning order
- `responsive-design` for different terminal sizes and window changes
- `accessibility-audit` for contrast, focus states, and reduced-motion handling
- `heuristic-evaluation` for reviewing the terminal UI before release
- `design-qa-checklist` for final UI quality checks
- `design-rationale` for documenting why each terminal interaction exists

### How these skills should guide the build
- use `frontend-design` before choosing the visual language
- use `ui-design` and `layout-grid` before building panels or dashboards
- use `interaction-design` for playback controls, scans, and navigation
- use `motion-system` before writing animation timing
- use `micro-interaction-spec` for button states and track transitions
- use `accessibility-audit` to ensure the terminal remains usable in reduced-motion or low-color environments

## 5. CLI Experience Model

The CLI is not just a music player.
It is an **interactive ambient developer surface**.

It should support three kinds of engagement:

### 5.1 Passive mode
The CLI shows mood, track, and visual state without demanding action.

### 5.2 Interactive mode
The user can control music, choose sources, switch moods, and change themes.

### 5.3 Autonomous mode
The CLI keeps scanning the codebase and adjusting the soundtrack automatically.

## 6. Main Design Principles

### 6.1 Terminal-first, not web-first
The interface must be built for terminals, not copied from a browser UI.

### 6.2 Fast and lightweight
Rendering must be cheap. The UI should never make the system feel heavy.

### 6.3 Beautiful but usable
The design should be visually rich, but all controls must remain clear and practical.

### 6.4 Highly interactive
The CLI should support meaningful keyboard actions and responsive feedback.

### 6.5 Graceful fallback
If a terminal cannot render advanced graphics, the interface must still work cleanly.

## 7. CLI Layout Structure

The CLI should use a clear layered layout.

### Suggested layout areas
1. **Header / Brand band**
2. **Mood status panel**
3. **Playback panel**
4. **Wave / visual panel**
5. **Track queue / source panel**
6. **Action footer**

### Typical visual hierarchy
- top: product name and session state
- center: current mood and playback state
- lower center: wave animation or ASCII art
- bottom: controls, shortcuts, hints, and status

## 8. Core Interactive Components

### 8.1 Play / Pause / Stop controls
These should behave like real interactive UI actions, even in a terminal.

Required states:
- idle
- playing
- paused
- buffering
- switching track
- error

### 8.2 Next / Previous / Skip
The user should be able to move between tracks without breaking the visual flow.

### 8.3 Mood mode switch
Allow manual override or automatic selection.
- focused
- calm
- intense
- experimental
- late-night
- debugging
- sprint mode

### 8.4 Source selector
Let users pick which source or fallback mode they want.

### 8.5 Visual mode selector
Switch between:
- wave mode
- spectrum mode
- pulse mode
- minimal mode
- quiet mode
- text-only mode

### 8.6 Scan button / action
The CLI should support one explicit scan action that triggers repository analysis and mood detection.

## 9. Terminal Visual Design

The visuals should feel handcrafted.

### 9.1 ASCII art system
Use ASCII art for:
- product banners
- ambient frames
- idle states
- mode changes
- loading moments

### 9.2 Reactive music waves
The wave animation should respond to:
- playback intensity
- mood score
- track energy
- volume level
- track changes

### 9.3 Spectrum rendering
Use a simple and efficient bar system.
- bars rise and fall with amplitude
- spacing should remain tidy
- color should change slowly, not abruptly

### 9.4 Dot-matrix styling
A dot-matrix look can be used for:
- headers
- mood labels
- track names
- status chips
- empty-state messages

### 9.5 Animated scene changes
The CLI should be able to transition between scenes:
- scan scene
- mood scene
- playback scene
- paused scene
- quiet scene
- error scene

## 10. Wave Animation Design

### Goals
- smooth motion
- cheap to render
- readable at terminal size
- good even when the terminal is resized

### Suggested wave logic
- use a small frame buffer
- calculate amplitude from music energy or simulated peaks
- render bars or curves with terminal-safe characters
- update at a target frame rate only if the terminal can handle it

### Fallback behavior
- if animation is not supported, show static waveform snapshots
- if color is unavailable, render monochrome bars
- if width is limited, switch to compact bars

## 11. Music-Responsive UI

The CLI should react to music in a visible way.

### Reactions
- bar height changes with energy
- ambient tint changes with mood
- status labels pulse subtly on track change
- progress indicators advance smoothly
- track transitions fade rather than jump

### Music-driven states
- calm music -> soft motion and slow wave changes
- lo-fi -> low amplitude and relaxed rendering
- intense music -> stronger bars and faster updates
- experimental music -> more dynamic color shifts and sharper motion

## 12. Keyboard Interaction Design

The terminal should feel interactive through keyboard input.

### Suggested key actions
- `Space` → play / pause
- `N` → next track
- `P` → previous track
- `S` → scan repository
- `M` → switch mood mode
- `V` → switch visual mode
- `Q` → queue view
- `Ctrl+C` → clean exit
- `?` → help overlay

### Design rules
- shortcuts must be visible in the footer or help panel
- controls should be discoverable
- shortcuts should not be overloaded
- every major action should have a keyboard path

## 13. UI States

The CLI should explicitly design the following states.

### Idle
No music is playing. Show a calm visual and a prompt to scan or play.

### Scanning
Show progress, repo analysis, and a subtle loading animation.

### Mood detected
Show the inferred mood and confidence.

### Playing
Show track info, progress, and reactive visuals.

### Paused
Keep the layout stable and indicate paused state clearly.

### Switching track
Show a transition moment with fade or sweep motion.

### Error
Show a concise, human-readable error and fallback option.

### Quiet mode
Show analysis and recommendations without playback.

## 14. Layout and Composition

### Structure rules
- maintain strong alignment
- avoid clutter
- keep spacing regular
- keep controls near the bottom or within a stable footer band
- ensure the main visual center is easy to read

### Terminal proportions
The layout should adapt to:
- narrow terminals
- standard terminals
- wide terminal monitors
- resized windows

### Composition strategy
- central focus area for mood and playback
- side or lower panel for queue and source details
- footer for shortcuts and current status

## 15. Typography and Labels

### Text hierarchy
- large title for product / session
- medium labels for mood and track status
- small helper text for shortcuts and technical status
- code-style text for commands and source names

### Typography behavior
- avoid dense text blocks
- keep labels short
- use spacing to separate groups
- keep track names readable

## 16. Color and Contrast

Even if the product is monochrome-first, the CLI may still use color sparingly when supported.

### Rules
- default to black, white, and gray
- use color only for meaning or mood tinting
- do not rely on color alone
- maintain high contrast for all critical text
- provide monochrome fallback for every color treatment

## 17. Animation System

### Animation principles
- animations should be subtle, meaningful, and light
- no flashy or distracting motion
- prefer slow transitions and readable transformations
- keep frame generation cheap

### Recommended animation types
- fade in
- fade out
- slide in
- pulse
- wave motion
- equalizer oscillation
- cursor blink
- soft progress fill
- scene crossfade

## 18. External Library Strategy

The CLI may use external libraries where they improve the final experience.

### Possible needs
- terminal rendering helpers
- color and ANSI helpers
- keyboard input handling
- animation timing
- audio backend integration
- file watching
- HTTP clients
- config parsing

### Rule
External libraries are allowed, but the CLI must still feel like a coherent from-scratch product. Use dependencies as tools, not as the identity of the product.

## 19. Fallback Design

The CLI must degrade gracefully.

### Fallback paths
- no animation support -> static layout
- no color -> monochrome only
- no unicode -> ASCII only
- no TTY -> text logs
- no audio backend -> recommendations only
- no internet -> local or cached state

## 20. Accessibility Rules

The interface must remain usable for everyday users.

### Required
- visible focus states
- keyboard navigation
- readable contrast
- reduced motion support
- clear shortcut discovery
- no hidden meaning in animations

## 21. Build Order

### Step 1: Terminal layout
Build the basic structure and visual hierarchy.

### Step 2: Controls
Add play, pause, stop, next, previous, scan, and mode selection.

### Step 3: Visual engine
Add waveforms, ASCII art, dot-matrix styling, and state animations.

### Step 4: Playback linkage
Connect visuals to music state and track changes.

### Step 5: Behavior polish
Refine motion, spacing, and interaction feedback.

### Step 6: QA and accessibility
Audit the interface for readability, fallback behavior, and keyboard usability.

## 22. Definition of a Great CLI Design

A great Moodwave CLI design:
- feels handcrafted and premium
- uses ASCII and motion with taste
- reacts to music naturally
- gives users real control
- remains lightweight and fast
- works well on Windows, macOS, and Linux
- supports both power users and everyday users
- looks and feels like an artistically designed terminal product
