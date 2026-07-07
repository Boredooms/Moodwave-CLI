# Moodwave CLI — Promotional Website

## 1. Purpose

The website is **not** the core product. It is a promotional layer for Moodwave CLI.

Its job is to:
- introduce the product
- explain the CLI-first concept
- show why the tool is unique
- present screenshots, motion, and demos
- provide installation commands
- document the CLI at a high level
- drive users to download and try the CLI

The website should feel premium, cinematic, minimal, and carefully animated.

## 2. Design Direction

### Overall style
- black and white
- monochrome only
- high-contrast premium aesthetic
- terminal-inspired but polished
- futuristic without becoming flashy
- calm, elegant, and high-end

### Visual identity
- deep black backgrounds
- soft white typography
- subtle gray hierarchy
- restrained accent use if absolutely necessary
- thin lines, borders, and overlays
- ambient motion instead of loud motion

### Emotional goal
The website should feel like:
- a premium developer product page
- a refined terminal showcase
- a serious tool with a fun soul
- a launch page that makes people want to try the CLI immediately

## 3. Motion Philosophy

Motion should be smooth, elegant, and purposeful.

### Recommended motion behavior
- soft fade-ins
- staggered text reveals
- slow parallax depth
- scroll-triggered section transitions
- subtle cursor or waveform motion
- gentle hover lift
- ambient animated noise or grain
- fluid section entrances

### Motion rules
- motion must support content, not distract from it
- no excessive bounce
- no chaotic motion
- keep transitions refined and readable
- animation should feel like a premium product, not a demo reel

## 4. Suggested Stack

The website should be built with a modern front-end stack suitable for premium motion and performance.

### Recommended tools
- Next.js
- TypeScript
- Tailwind CSS
- Framer Motion
- GSAP where needed for advanced timeline control
- Lenis for smooth scrolling
- Anime.js where small expressive sequences are useful
- optional canvas or WebGL helpers only if truly needed

### Principle
Use motion libraries selectively. Do not stack libraries just because they exist. Each one should solve a specific problem.

## 5. Page Structure

### 5.1 Hero section
The first screen should immediately communicate:
- what the product is
- that it is CLI-first
- that it reacts to codebase mood
- that it plays ambient music in the terminal
- that it is lightweight and cross-platform

Hero elements:
- strong headline
- short supporting subtitle
- visible install command
- animated terminal preview
- subtle background motion

### 5.2 Product story
Explain the idea in a compact, compelling narrative:
- codebase analysis
- mood inference
- music recommendation
- terminal ambient experience

### 5.3 Feature highlights
Show the major pillars:
- code scanning
- mood prediction
- music fetching
- terminal visuals
- playback control
- cross-platform support

### 5.4 How it works
Use a step-by-step flow:
1. install
2. scan
3. infer mood
4. fetch music
5. play
6. adapt over time

### 5.5 CLI preview section
Show what the terminal looks like.
- command examples
- animated mock terminal
- command output style
- music mood states
- visual modes

### 5.6 Architecture or engine section
Briefly show the technical thoughtfulness:
- scanner
- mood engine
- recommender
- source adapters
- playback layer
- terminal renderer

### 5.7 Installation section
The website should include copy-paste install commands.
- curl install
- PowerShell install
- release download
- optional package manager info later

### 5.8 FAQ section
Answer common questions:
- is the CLI heavy?
- does it store music locally?
- does it work on Windows?
- what music sources are used?
- is the website the main app?

### 5.9 Call to action
End with a strong CTA:
- download the CLI
- read the docs
- try the first scan
- join the project

## 6. Layout Principles

The layout should be elegant and deliberate.

### Principles
- generous whitespace
- strong grid structure
- clear hierarchy
- one main idea per section
- no cluttered cards everywhere
- no generic SaaS layout
- no overused landing page templates

### Recommended content rhythm
- large hero
- short explanation block
- feature grid
- terminal showcase
- deep technical section
- install and CTA
- minimal footer

## 7. Typography Direction

Typography should be a major visual feature.

### Rules
- large, confident headings
- readable body text
- strong contrast
- careful line length
- typography should feel editorial and precise
- use motion sparingly on text reveal

### Effects
- letter fade-ins
- staggered line reveals
- clipped text transitions
- subtle tracking changes on key headings

## 8. Terminal Showcase Section

This is a critical section.

The website should visually simulate the CLI with:
- dark terminal frame
- animated text output
- mood labels
- music status lines
- waveform strips
- command examples
- theme and track changes

The terminal preview should feel real, not decorative.

## 9. Visual Language

### Background
- pure black or near-black surfaces
- subtle gradients if needed
- grain/noise for depth
- thin border lines

### Foreground
- white typography
- soft grays for secondary text
- occasional mono highlights

### Components
- minimal buttons
- outlined command cards
- thin separators
- subtle glass or matte surfaces only if tasteful

## 10. Motion System Ideas

### Recommended animations
- hero text fade and slide in
- command line typing effect
- terminal cursor blink
- waveform pulse
- scroll-based reveal sequencing
- section opacity and translate transitions
- CTA hover glow
- logo and icon micro-animations

### Recommended libraries
- Framer Motion for component transitions
- GSAP for complex timeline moments
- Lenis for premium smooth scrolling
- Anime.js for lightweight text or icon motion

## 11. Content Strategy

The website should speak to developers.

### Tone
- confident
- concise
- technically credible
- a little playful
- not corporate
- not vague
- not overly marketing-heavy

### Messaging pillars
- CLI-first
- mood-aware
- music-reactive
- lightweight
- open-source friendly
- terminal-native
- beautiful and practical

## 12. Recommended Sections in Order

1. Hero
2. Why Moodwave exists
3. How it works
4. Terminal preview
5. Feature breakdown
6. Architecture summary
7. Installation commands
8. FAQ
9. Final CTA

## 13. Performance Rules

The website must feel fast.

### Rules
- optimize motion carefully
- avoid huge background videos unless truly necessary
- avoid bloated libraries that are not needed
- keep assets compressed
- lazy load heavy visuals
- render above-the-fold content immediately
- keep the experience smooth on mid-range machines

## 14. Accessibility Rules

Even in a monochrome design, accessibility matters.

### Must support
- readable contrast
- keyboard navigation
- reduced-motion mode
- clear focus states
- semantic headings
- legible installation commands
- no text hidden only in motion

## 15. Technical Content on the Website

The website should briefly explain the CLI architecture without overwhelming the visitor.

Include:
- repository scanner
- mood engine
- source adapters
- playback layer
- visual layer
- install flow

Do not dump the whole technical doc onto the homepage. Keep the main page concise and elegant.

## 16. Installation Presentation

Use copy-paste blocks that are easy to use.

Examples of presentation goals:
- one-line Unix install
- PowerShell install
- direct release link
- update command

The commands should be visible, clean, and easy to copy.

## 17. Media and Demonstration

The website should include:
- screenshots of the terminal UI
- short GIF-like motion sections
- command output snippets
- mood states
- audio visualizer preview

Keep media stylized and restrained.

## 18. Component Ideas

Useful website components:
- hero terminal mockup
- command strip
- animated stats panel
- feature cards
- architecture diagram block
- timeline section
- FAQ accordion
- install code block
- final CTA banner

## 19. Build Order

### Step 1: Content and structure
Define sections and copy.

### Step 2: Base layout
Create the homepage shell, spacing, typography, and grid.

### Step 3: Motion system
Add smooth scroll, fades, and timeline motion.

### Step 4: Terminal mockup
Build the central terminal preview.

### Step 5: Feature sections
Add architecture, install, and FAQ content.

### Step 6: Polish
Refine spacing, interaction, contrast, and responsiveness.

## 20. What the Website Should Avoid

- generic SaaS hero layouts
- bright rainbow gradients
- unnecessary card overload
- too many motion effects at once
- cluttered marketing speak
- complicated navigation for a simple promo site
- hiding the CLI behind the website

## 21. Definition of a Great Promotional Website

A great Moodwave website:
- feels premium and cinematic
- communicates the CLI-first product clearly
- showcases the terminal beautifully
- uses motion with restraint and taste
- makes installation feel easy
- is simple enough to understand quickly
- makes the CLI look exciting, trustworthy, and worth trying
