# Moodwave CLI — Website Design

## 1. Purpose of the Website

The website is a **promotional and experiential layer** for Moodwave CLI.

It must not become the core product. The core product is the CLI. The website exists to:
- explain the product quickly
- show the terminal experience beautifully
- demonstrate the mood-to-music concept
- provide installation commands
- build excitement and trust
- drive users to try the CLI

The website should feel like an experience, not a brochure.

## 2. Design Direction

### Visual theme
- strict black and white
- monochrome first
- minimal gray hierarchy
- no colorful branding system
- no loud gradients
- no generic SaaS palette

### Tone
- premium
- cinematic
- technical
- calm
- sharp
- confident
- highly polished

### Design feeling
The website should feel like:
- a high-end developer product launch page
- a terminal-inspired motion showcase
- a refined, editorial, modern experience
- something built with taste from scratch

## 3. Design System Philosophy

The design must be built from scratch around a coherent system.

### Core system values
- spacing discipline
- typography-led hierarchy
- strong content rhythm
- precise alignment
- restrained motion
- consistent borders and surfaces
- readable installation commands
- component reusability

### Design objectives
- every section should feel connected
- every visual decision should support the product story
- nothing should feel random or decorative without purpose
- the whole page should flow like a guided experience

## 4. Skills to Explicitly Harness in Antigravity

The website design should deliberately use the skills already installed in Antigravity for each part of the page.

### Key skills to apply
- `frontend-design` for distinctive UI direction and avoiding templated layouts
- `ui-design` for visual hierarchy, layout, color, typography, and screen composition
- `interaction-design` for scroll behavior, transitions, state changes, and motion logic
- `motion-system` for animation timing, easing, and consistency
- `micro-interaction-spec` for hover, click, reveal, and terminal-like feedback
- `layout-grid` for section structure and responsive layout control
- `spacing-system` for exact spacing rhythm across the page
- `typography-scale` for headings, copy, and code block hierarchy
- `visual-hierarchy` for strong focal points and page scanning order
- `responsive-design` for adaptive behavior across screen sizes
- `theme-factory` for generating and maintaining the monochrome theme
- `accessibility-audit` for contrast, focus states, and reduced-motion behavior
- `heuristic-evaluation` for reviewing the page before shipping
- `design-qa-checklist` for final QA before implementation
- `design-rationale` for documenting why each design choice exists

### How these skills should be used
The site should not be built by guessing. The agent should actively route work through the relevant skills:
- use `frontend-design` before choosing the visual direction
- use `ui-design` before finalizing layout and page hierarchy
- use `interaction-design` when defining scroll and entrance motion
- use `motion-system` before writing animations
- use `typography-scale` and `visual-hierarchy` before finalizing headings and body copy
- use `spacing-system` and `layout-grid` before building components
- use `accessibility-audit` and `design-qa-checklist` before release

## 5. Stack Direction

The website should use a modern stack that supports smooth, premium motion.

### Recommended stack
- Next.js
- TypeScript
- Tailwind CSS
- Radix UI
- Lenis for smooth scrolling
- Framer Motion for standard motion
- GSAP for timeline-heavy sections where needed
- Anime.js for small expressive text or icon animation if useful

### Principles for library use
- use only what improves the final experience
- do not overload the project with unnecessary animation libraries
- keep the implementation maintainable
- keep performance strong

## 6. Page Experience Goals

The page should feel like a single coherent journey.

### What users should feel
- immediate curiosity at the hero
- clarity about the product in seconds
- excitement when they see the terminal preview
- confidence when they read installation steps
- trust when they see architecture and technical detail
- motivation to install and try the CLI

## 7. Page Structure

### 7.1 Hero section
The hero must establish the identity instantly.

Should include:
- bold product title
- sharp one-line positioning statement
- short supporting explanation
- visible install command
- animated terminal mockup
- subtle ambient motion

The hero should make the CLI feel alive before the user scrolls.

### 7.2 Product explanation section
Explain the idea in plain terms:
- the CLI scans codebases
- it infers mood
- it recommends music
- it plays in terminal
- it reacts to code changes

### 7.3 Terminal experience section
Show what the CLI looks like.
- command output
- animated status lines
- mood states
- waveform and dot-matrix visuals
- minimal but expressive terminal framing

### 7.4 Feature section
Show the core benefits:
- mood detection
- music matching
- terminal visuals
- lightweight runtime
- source flexibility
- cross-platform support

### 7.5 How it works section
A simple 1–2–3–4 flow:
1. install
2. scan
3. infer mood
4. play and adapt

### 7.6 Technical section
Briefly explain the CLI architecture:
- scanner
- mood engine
- recommendation engine
- source adapters
- playback layer
- visual layer

### 7.7 Installation section
Show copy-paste commands.
- curl install for Unix-like systems
- PowerShell install for Windows
- optional release download button

### 7.8 FAQ section
Answer common concerns:
- is the CLI heavy?
- does it store music locally?
- what sources does it use?
- does it work on Windows?
- is the website the product?

### 7.9 Final CTA section
End with a strong call to action:
- install the CLI
- read the docs
- try the first scan

## 8. Layout Strategy

The page should avoid generic templates.

### Layout rules
- use a strong grid
- keep section widths disciplined
- allow generous whitespace
- keep one focal point per section
- avoid cluttered card walls
- avoid endless feature tiles
- use deliberate pacing between content blocks

### Recommended rhythm
- hero
- brief narrative
- terminal showcase
- feature grid
- technical explanation
- install commands
- FAQ
- CTA

## 9. Typography Strategy

Typography should carry most of the visual weight.

### Type rules
- large, confident headings
- calm, readable body copy
- code-style blocks for commands
- clear hierarchy between heading levels
- no overcomplicated font mixing
- no decorative type gimmicks

### Motion on text
- letter fade-in
- staggered word or line reveal
- gentle clipping transitions
- subtle tracking adjustment on hero lines

## 10. Terminal Mockup Design

This is a signature section.

### The terminal should show
- scan output
- mood classification
- track selection
- music state
- ambient visual bars
- playback progress
- subtle theme changes

### Terminal style
- matte black surface
- thin borders
- soft white text
- subdued contrast hierarchy
- visible command prompt
- minimal but expressive animation

### Visual logic
The terminal mockup should feel practical, not fake.
It should resemble the actual product experience.

## 11. Motion and Scroll Design

The motion system should feel refined.

### Recommended behaviors
- smooth scrolling with Lenis
- fade-up entrances
- staggered reveal sequences
- scroll-triggered section activation
- subtle parallax depth
- hover polish on buttons and cards
- tiny motion on terminal indicators

### Motion rules
- motion should never overpower reading
- no overly springy motion
- no gimmicky rapid animations
- keep transitions premium and calm

## 12. Interaction Design

The website should feel responsive to user intent.

### Important interactions
- command copy buttons
- install command copy-to-clipboard
- hover reveal on cards
- animated terminal tabs or status chips
- FAQ expand/collapse
- button hover states with subtle luminance changes

### Interaction philosophy
The page should reward exploration without requiring it.

## 13. Component Inventory

Useful reusable components:
- hero terminal panel
- command block
- feature card
- stat line
- section heading block
- timeline / process steps
- architecture callout
- FAQ accordion
- CTA strip
- footer

## 14. Visual Rules

### Background
- black or near-black base
- soft layered surfaces
- faint gradients only if necessary
- subtle grain or texture if tasteful

### Foreground
- white text
- light gray secondary text
- clear divider lines
- minimal accent usage

### Borders and surfaces
- thin, precise borders
- soft contrast surfaces
- no heavy shadows
- no noisy glassmorphism

## 15. Accessibility Rules

The design must remain usable and readable.

### Required
- strong contrast
- visible focus states
- keyboard navigation
- reduced motion support
- semantic headings
- accessible code blocks
- no color-only meaning
- readable line lengths

## 16. Performance Rules

The site must feel fast.

### Performance principles
- keep motion efficient
- avoid huge background videos
- keep assets optimized
- lazy load heavy media
- avoid unnecessary dependencies
- render content immediately above the fold

## 17. Content Style

The copy should be concise and technical.

### Voice
- direct
- confident
- clear
- developer-friendly
- not corporate
- not vague
- not over-marketed

### Messaging focus
- CLI-first
- mood-aware
- music-reactive
- terminal-native
- lightweight
- open-source friendly
- premium experience

## 18. Build Order

### Step 1: Apply design skills
Use the installed design skills to define hierarchy, spacing, motion, and accessibility.

### Step 2: Build structure
Create the homepage shell, layout, and content sequence.

### Step 3: Add motion
Implement Lenis, Framer Motion, and any carefully scoped GSAP or Anime.js sequences.

### Step 4: Create terminal showcase
Build the core visual hero and terminal preview.

### Step 5: Add content sections
Add features, architecture, install commands, and FAQ.

### Step 6: Polish and audit
Run heuristic review, accessibility checks, and design QA.

## 19. What the Website Must Avoid

- generic SaaS aesthetics
- overused gradient-heavy visuals
- excessive card grids
- motion with no purpose
- cluttered copy blocks
- flashy but meaningless effects
- hiding the CLI behind the website

## 20. Definition of a Great Website

A great Moodwave website:
- makes the CLI feel premium
- uses monochrome design with intent
- has smooth, elegant motion
- explains the product quickly
- showcases the terminal beautifully
- feels like a crafted experience from scratch
- uses the installed design skills deliberately
- drives users to install the CLI
