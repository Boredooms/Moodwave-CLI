// Package visuals implements the terminal rendering engine for Moodwave CLI.
//
// Design philosophy (from cli_design.md + motion-system + visual-hierarchy skills):
//
//	Visual hierarchy: brand band → mood panel → playback panel → wave → footer
//	Motion system:    all animations follow duration tokens and easing curves
//	Fallback chain:   rich ANSI → basic ANSI → monochrome → ASCII → log-only
//
// The renderer never blocks audio playback. It operates independently and
// degrades cleanly when the terminal cannot support advanced rendering.
package visuals

import (
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/moodwave/moodwave/internal/config"
	"github.com/moodwave/moodwave/internal/mood"
	"github.com/moodwave/moodwave/internal/platform"
	"github.com/moodwave/moodwave/internal/recommender"
)

// ──────────────────────────────────────────────────────────────────────────────
// ANSI escape codes
// ──────────────────────────────────────────────────────────────────────────────

const (
	ansiReset   = "\033[0m"
	ansiBold    = "\033[1m"
	ansiDim     = "\033[2m"
	ansiItalic  = "\033[3m"
	ansiReverse = "\033[7m"

	// Cursor control
	ansiHideCursor     = "\033[?25l"
	ansiShowCursor     = "\033[?25h"
	ansiClearLine      = "\033[2K"
	ansiClearDown      = "\033[J"
	ansiHome           = "\033[H"
	ansiAltScreenEnter = "\033[?1049h"
	ansiAltScreenExit  = "\033[?1049l"

	// Colors (foreground)
	ansiFgBlack   = "\033[30m"
	ansiFgRed     = "\033[31m"
	ansiFgGreen   = "\033[32m"
	ansiFgYellow  = "\033[33m"
	ansiFgBlue    = "\033[34m"
	ansiFgMagenta = "\033[35m"
	ansiFgCyan    = "\033[36m"
	ansiFgWhite   = "\033[37m"
	ansiFgBright  = "\033[97m"

	// 256-color gray shades (for monochrome visual hierarchy)
	ansiGray10  = "\033[38;5;234m" // near black
	ansiGray30  = "\033[38;5;238m"
	ansiGray50  = "\033[38;5;242m"
	ansiGray70  = "\033[38;5;246m"
	ansiGray90  = "\033[38;5;250m"
	ansiGray100 = "\033[38;5;255m" // near white
)

// SceneType identifies which visual scene is active.
type SceneType string

const (
	SceneIdle      SceneType = "idle"
	SceneScanning  SceneType = "scanning"
	SceneMoodReady SceneType = "mood"
	ScenePlaying   SceneType = "playing"
	ScenePaused    SceneType = "paused"
	SceneSwitching SceneType = "switching"
	SceneError     SceneType = "error"
	SceneQuiet     SceneType = "quiet"
	SceneWelcome   SceneType = "welcome"
)

// RendererConfig configures the renderer.
type RendererConfig struct {
	// VisualMode is the active visual mode.
	VisualMode config.VisualMode

	// Theme is the active color theme.
	Theme config.ThemeID

	// FPS is the target frame rate.
	FPS int

	// NoAnimation disables animation.
	NoAnimation bool

	// NoColor disables ANSI color codes.
	NoColor bool

	// NoUnicode forces ASCII characters.
	NoUnicode bool

	// Caps describes terminal capabilities.
	Caps platform.Capabilities

	// Output is the writer for terminal output (default: os.Stdout).
	Output io.Writer
}

// RenderState is the current state of what should be rendered.
type RenderState struct {
	Scene      SceneType
	Mood       *mood.Profile
	Candidate  *recommender.Candidate
	Progress   float64 // 0.0–1.0
	Error      string
	ScanMsg    string
	InputMode  bool
	InputBuf   string
	RepeatMode string
}

// Renderer is the main terminal rendering engine.
type Renderer struct {
	mu       sync.Mutex
	cfg      RendererConfig
	state    RenderState
	frame    int
	ticker   *time.Ticker
	stop     chan struct{}
	out      io.Writer
	fireGrid [][]int // persistent fire state matrix
}

// New creates a Renderer.
func New(cfg RendererConfig) *Renderer {
	out := cfg.Output
	if out == nil {
		out = os.Stdout
	}
	if cfg.FPS <= 0 || cfg.FPS > 60 {
		cfg.FPS = 12 // conservative default
	}

	// Respect capabilities.
	if !cfg.Caps.HasColor {
		cfg.NoColor = true
	}
	if !cfg.Caps.HasUnicode {
		cfg.NoUnicode = true
	}
	if !cfg.Caps.HasAnimation || cfg.Caps.ReducedMotion {
		cfg.NoAnimation = true
	}

	return &Renderer{
		cfg:  cfg,
		out:  out,
		stop: make(chan struct{}),
		state: RenderState{
			Scene: SceneIdle,
		},
	}
}

// SetState updates the render state. Thread-safe.
func (r *Renderer) SetState(s RenderState) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if s.RepeatMode == "" {
		s.RepeatMode = r.state.RepeatMode
	}
	r.state = s
}

// SetScene changes the current scene. Thread-safe.
func (r *Renderer) SetScene(scene SceneType) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.state.Scene = scene
}

// SetVisualMode updates the visual mode dynamically. Thread-safe.
func (r *Renderer) SetVisualMode(mode config.VisualMode) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cfg.VisualMode = mode
}

// Start begins the animation loop. Non-blocking.
func (r *Renderer) Start() {
	if !r.cfg.NoColor {
		fmt.Fprint(r.out, ansiAltScreenEnter)
	}

	if r.cfg.NoAnimation {
		// Render once statically.
		r.RenderFrame()
		return
	}

	interval := time.Second / time.Duration(r.cfg.FPS)
	r.ticker = time.NewTicker(interval)

	go func() {
		for {
			select {
			case <-r.ticker.C:
				r.mu.Lock()
				r.frame++
				r.mu.Unlock()
				r.RenderFrame()
			case <-r.stop:
				return
			}
		}
	}()
}

// Stop halts the animation loop and restores the cursor.
func (r *Renderer) Stop() {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if already closed/stopped to prevent panic
	select {
	case <-r.stop:
		return
	default:
	}

	if r.ticker != nil {
		r.ticker.Stop()
	}
	close(r.stop)
	// Restore cursor and exit alternate screen.
	if !r.cfg.NoColor {
		fmt.Fprint(r.out, ansiAltScreenExit)
	}
	fmt.Fprint(r.out, ansiShowCursor)
}

// RenderFrame renders a single frame based on the current state.
func (r *Renderer) RenderFrame() {
	r.mu.Lock()
	state := r.state
	frame := r.frame
	r.mu.Unlock()

	width := r.cfg.Caps.SafeWidth()
	out := &strings.Builder{}

	// Hide cursor during render.
	if !r.cfg.NoColor {
		out.WriteString(ansiHideCursor)
	}

	// Move to top.
	out.WriteString(ansiHome)
	out.WriteString(ansiClearDown)

	if state.Scene == SceneWelcome {
		r.renderWelcome(out, width, state, frame)
		fmt.Fprint(r.out, out.String())
		return
	}

	// Render each section.
	r.renderHeader(out, width, state)
	r.renderMoodPanel(out, width, state)
	r.renderPlaybackPanel(out, width, state)
	r.renderVisualPanel(out, width, state, frame)
	r.renderFooter(out, width, state)

	// Flush to output.
	fmt.Fprint(r.out, out.String())
}

// ──────────────────────────────────────────────────────────────────────────────
// Section renderers
// ──────────────────────────────────────────────────────────────────────────────

// renderHeader renders the brand band at the top.
func (r *Renderer) renderHeader(out *strings.Builder, width int, _ RenderState) {
	title := "MOODWAVE"

	if r.cfg.Caps.SafeWidth() >= 60 {
		// Dot-matrix banner mode.
		banner := renderDotMatrixBanner(title, r.cfg.NoUnicode)
		for _, line := range banner {
			r.writeLine(out, centerPad(line, width))
		}
	} else {
		// Compact mode.
		r.writeBold(out, centerPad("◆ MOODWAVE ◆", width))
	}

	r.writeDim(out, centerPad("terminal mood music companion", width))
	r.renderSeparator(out, width)
}

// renderMoodPanel renders the detected mood and confidence.
func (r *Renderer) renderMoodPanel(out *strings.Builder, width int, state RenderState) {
	if state.Mood == nil {
		switch state.Scene {
		case SceneIdle:
			r.writeDim(out, centerPad("No mood detected. Run 'moodwave scan' to begin.", width))
		case SceneScanning:
			dots := strings.Repeat(".", (r.frame%4)+1)
			r.writeLine(out, centerPad("Scanning"+dots+" "+state.ScanMsg, width))
		case SceneError:
			r.writeError(out, centerPad("Error: "+state.Error, width))
		}
		out.WriteString("\n")
		return
	}

	m := state.Mood
	label := strings.ToUpper(string(m.Label))

	dynamicConfidence := m.Confidence
	if state.Scene == ScenePlaying {
		energy := 0.5
		if state.Candidate != nil && state.Candidate.Track != nil {
			energy = state.Candidate.Track.Energy
		}
		bounce := math.Sin(float64(r.frame)*0.15) * 0.04 * energy
		dynamicConfidence = math.Max(0.0, math.Min(1.0, m.Confidence+bounce))
	}

	confidence := fmt.Sprintf("%.0f%% confidence", dynamicConfidence*100)

	r.writeBold(out, centerPad(m.Label.Emoji()+" "+label, width))
	r.writeDim(out, centerPad(confidence, width))

	// Mood meter bar.
	meterWidth := width / 2
	if meterWidth < 20 {
		meterWidth = 20
	}
	meterWidth = min2(meterWidth, 50)
	filled := int(math.Round(float64(meterWidth) * dynamicConfidence))
	meter := strings.Repeat("█", filled) + strings.Repeat("░", meterWidth-filled)
	if r.cfg.NoUnicode {
		meter = strings.Repeat("#", filled) + strings.Repeat("-", meterWidth-filled)
	}
	r.writeLine(out, centerPad("["+meter+"]", width))

	out.WriteString("\n")
}

// renderPlaybackPanel renders the current track/station info.
func (r *Renderer) renderPlaybackPanel(out *strings.Builder, width int, state RenderState) {
	if state.Candidate == nil {
		out.WriteString("\n")
		return
	}

	c := state.Candidate
	name := c.DisplayName()
	sourceLabel := ""
	if c.IsStation() {
		sourceLabel = "📻 " + c.Station.Source
	} else if c.Track != nil {
		sourceLabel = "🎵 " + c.Track.Source
	}

	r.writeBold(out, centerPad(name, width))
	if sourceLabel != "" {
		r.writeDim(out, centerPad(sourceLabel, width))
	}

	// Progress bar (if available).
	if state.Progress > 0 {
		barWidth := min2(width-4, 60)
		filled := int(math.Round(float64(barWidth) * state.Progress))
		if r.cfg.NoUnicode {
			bar := "[" + strings.Repeat("=", filled) + strings.Repeat("-", barWidth-filled) + "]"
			r.writeLine(out, centerPad(bar, width))
		} else {
			bar := "▐" + strings.Repeat("▓", filled) + strings.Repeat("░", barWidth-filled) + "▌"
			r.writeLine(out, centerPad(bar, width))
		}
	}

	// State indicator.
	loopStatus := ""
	if state.RepeatMode != "" && state.RepeatMode != "off" {
		loopStatus = fmt.Sprintf("  🔄 LOOP: %s", strings.ToUpper(state.RepeatMode))
		if r.cfg.NoUnicode {
			loopStatus = fmt.Sprintf("  [LOOP: %s]", strings.ToUpper(state.RepeatMode))
		}
	}

	switch state.Scene {
	case ScenePlaying:
		indicator := "▶ PLAYING" + loopStatus
		if r.cfg.NoUnicode {
			indicator = "> PLAYING" + loopStatus
		}
		r.writeLine(out, centerPad(indicator, width))
	case ScenePaused:
		indicator := "⏸ PAUSED" + loopStatus
		if r.cfg.NoUnicode {
			indicator = "|| PAUSED" + loopStatus
		}
		r.writeDim(out, centerPad(indicator, width))
	case SceneSwitching:
		r.writeDim(out, centerPad("~ switching track ~", width))
	}

	out.WriteString("\n")
}

// renderArtwork renders an animated retro vinyl record.
func (r *Renderer) renderArtwork(out *strings.Builder, width int, frame int) {
	if r.cfg.NoAnimation {
		r.writeDim(out, centerPad("    .---.    ", width))
		r.writeDim(out, centerPad("   / (O) \\   ", width))
		r.writeDim(out, centerPad("    '---'    ", width))
		return
	}

	frames := [][]string{
		{
			"    .---.    ",
			"   /  |  \\   ",
			"  |  (O)  |  ",
			"   \\     /   ",
			"    '---'    ",
		},
		{
			"    .---.    ",
			"   /   /  \\   ",
			"  |  (O)  |  ",
			"   \\  /  /   ",
			"    '---'    ",
		},
		{
			"    .---.    ",
			"   /  -  \\   ",
			"  |  (O)  |  ",
			"   \\  -  /   ",
			"    '---'    ",
		},
		{
			"    .---.    ",
			"   /  \\   \\   ",
			"  |  (O)  |  ",
			"   \\   \\ /   ",
			"    '---'    ",
		},
	}

	currentFrame := frames[(frame/2)%len(frames)]
	for _, line := range currentFrame {
		r.writeDim(out, centerPad(line, width))
	}
}

// renderVisualPanel renders the animated wave, spectrum, or pulse.
func (r *Renderer) renderVisualPanel(out *strings.Builder, width int, state RenderState, frame int) {
	if r.cfg.VisualMode == config.VisualQuiet {
		return
	}

	// Render spinning artwork if a candidate is actively playing.
	if state.Scene == ScenePlaying && state.Candidate != nil {
		r.renderArtwork(out, width, frame)
		out.WriteString("\n")
	}

	energy := 0.25
	if state.Scene == ScenePlaying {
		energy = 0.5 // default energy when no track info
		if state.Candidate != nil && state.Candidate.Track != nil {
			energy = state.Candidate.Track.Energy
		}
		if state.Mood != nil {
			// Blend mood energy into visual energy.
			moodEnergy := (state.Mood.Traits.EnergyMin + state.Mood.Traits.EnergyMax) / 2.0
			energy = energy*0.4 + moodEnergy*0.6
		}
	}
	if energy < 0.25 {
		energy = 0.25
	}

	switch r.cfg.VisualMode {
	case config.VisualWave:
		r.renderWave(out, width, energy, frame)
	case config.VisualSpectrum:
		r.renderSpectrum(out, width, energy, frame)
	case config.VisualPulse:
		r.renderPulse(out, width, energy, frame)
	case config.VisualFire:
		r.renderFireplace(out, width, energy, frame)
	case config.VisualMinimal:
		r.renderMinimalVisual(out, width, state)
	}
}

// renderFireplace renders a cozy cellular-automaton fireplace.
func (r *Renderer) renderFireplace(out *strings.Builder, width int, energy float64, frame int) {
	fireRows := 7
	fireCols := width - 2
	if fireCols < 1 {
		fireCols = 1
	}

	// Initialize or resize the fire grid if needed
	if r.fireGrid == nil || len(r.fireGrid) != fireRows || len(r.fireGrid[0]) != fireCols {
		r.fireGrid = make([][]int, fireRows)
		for i := range r.fireGrid {
			r.fireGrid[i] = make([]int, fireCols)
		}
	}

	// Define characters and basic 16-color ANSI codes for maximum terminal compatibility.
	// Index 0 represents empty/cool. Highest index represents hottest flame.
	var firePalette = []struct {
		char  string
		color string
	}{
		{" ", ""},
		{".", "\033[90m"},    // dark gray
		{",", "\033[90m"},    // dark gray
		{"*", "\033[31m"},    // red
		{"x", "\033[31m"},    // red
		{"s", "\033[31m"},    // red
		{"o", "\033[91m"},    // bright red
		{"d", "\033[91m"},    // bright red
		{"m", "\033[33m"},    // yellow/orange
		{"0", "\033[33m"},    // yellow/orange
		{"H", "\033[93m"},    // bright yellow
		{"M", "\033[93m"},    // bright yellow
		{"W", "\033[97m"},    // bright white
		{"█", "\033[1;97m"},  // bold bright white
	}

	// Fallback palette for monochrome/NoUnicode
	if r.cfg.NoUnicode {
		firePalette[len(firePalette)-1].char = "W"
	}

	maxHeat := len(firePalette) - 1

	// Seed the bottom row in three distinct fireplace peaks: left (small), center (large), right (medium)
	for col := 0; col < fireCols; col++ {
		t := float64(col) / float64(fireCols)

		// Overlapping Gaussian bell curves to create flame mounds
		intensity := 0.0
		// Center peak (x = 0.5)
		intensity += math.Exp(-math.Pow((t-0.5)/0.12, 2)) * 1.0
		// Left peak (x = 0.22)
		intensity += math.Exp(-math.Pow((t-0.22)/0.07, 2)) * 0.65
		// Right peak (x = 0.78)
		intensity += math.Exp(-math.Pow((t-0.78)/0.08, 2)) * 0.75

		if intensity > 1.0 {
			intensity = 1.0
		}

		// Taper edges completely
		if t < 0.08 || t > 0.92 {
			intensity = 0
		}

		heatVal := int(intensity * float64(maxHeat))

		// Add flickering ember variation
		if heatVal > 0 {
			flicker := rand.Intn(3)
			heatVal -= flicker
			if heatVal < 0 {
				heatVal = 0
			}
		}

		r.fireGrid[0][col] = heatVal
	}

	// Propagate the flame upwards (cellular automaton style)
	for y := 1; y < fireRows; y++ {
		for x := 0; x < fireCols; x++ {
			// Pick a source cell from the row below with a random horizontal shift/wind
			wind := rand.Intn(3) - 1 // -1, 0, or 1
			srcX := (x + wind + fireCols) % fireCols

			// Apply decay based on the average energy of the song.
			// Higher energy means lower decay (flames rise taller).
			decayProb := 38
			if energy > 0.6 {
				decayProb = 24
			} else if energy < 0.35 {
				decayProb = 52
			}

			decay := rand.Intn(2) // 0 or 1
			if rand.Intn(100) < decayProb {
				decay++
			}

			val := r.fireGrid[y-1][srcX] - decay
			if val < 0 {
				val = 0
			}
			r.fireGrid[y][x] = val
		}
	}

	// Render the fire grid (top to bottom) with safe margins to prevent wrapping
	for y := fireRows - 1; y >= 0; y-- {
		var sb strings.Builder
		sb.WriteString(" ") // left padding
		for x := 0; x < fireCols; x++ {
			heat := r.fireGrid[y][x]
			cell := firePalette[heat]
			if r.cfg.NoColor {
				sb.WriteString(cell.char)
			} else {
				if cell.color != "" {
					sb.WriteString(cell.color + cell.char + ansiReset)
				} else {
					sb.WriteString(cell.char)
				}
			}
		}
		sb.WriteString(" \n") // right padding and newline
		out.WriteString(sb.String())
	}
}

// renderFooter renders the keyboard shortcuts and status bar.
func (r *Renderer) renderFooter(out *strings.Builder, width int, state RenderState) {
	r.renderSeparator(out, width)

	if state.InputMode {
		helper := "  Commands: play <song> | playnext <song> | add <song> | visuals [wave|spectrum|pulse|fireplace] | pause | resume"
		if !r.cfg.NoColor {
			helper = "  \033[1;30mCommands: \033[1;32mplay <song>\033[1;30m | \033[1;32mplaynext <song>\033[1;30m | \033[1;32madd <song>\033[1;30m | \033[1;32mvisuals [wave|spectrum|pulse|fireplace]\033[0m"
		}
		out.WriteString(helper + "\n")

		prompt := ": " + state.InputBuf + "▮"
		if !r.cfg.NoColor {
			prompt = "\033[1;33m: \033[1;37m" + state.InputBuf + "\033[5;37m▮\033[0m"
		}
		out.WriteString("  " + prompt + "\n")
		return
	}

	shortcuts := "  [Space] play/pause  [N] next  [L] loop  [S] scan  [V] visuals  [Q] quit  [Enter] command"
	if r.cfg.Caps.IsNarrow() {
		shortcuts = "  [Space] ⏯  [N] ⏭  [L] 🔄  [Q] quit  [Enter] ⌨"
	}

	r.writeDim(out, shortcuts)
}

// ──────────────────────────────────────────────────────────────────────────────
// Visual mode renderers
// ──────────────────────────────────────────────────────────────────────────────

// renderWave renders an animated sine-wave waveform.
func (r *Renderer) renderWave(out *strings.Builder, width int, energy float64, frame int) {
	height := 5
	if r.cfg.Caps.SafeHeight() < 20 {
		height = 3
	}

	amplitude := energy * float64(height) * 0.8
	blockChars := []rune{' ', ' ', '▂', '▃', '▄', '▅', '▆', '▇'}
	if r.cfg.NoUnicode {
		blockChars = []rune{' ', ' ', '.', ',', '-', '~', '=', '#'}
	}

	colorCode := ""
	if r.state.Mood != nil {
		colorCode = r.getWaveColor(r.state.Mood.Label)
	}

	for row := height - 1; row >= 0; row-- {
		line := make([]rune, width-2)
		for col := range line {
			phase := float64(frame)*0.25 + float64(col)*0.12
			val := (float64(height) / 2.0) + math.Sin(phase)*amplitude

			diff := val - float64(row)
			if diff >= 1.0 {
				if r.cfg.NoUnicode {
					line[col] = '#'
				} else {
					line[col] = '█'
				}
			} else if diff > 0.0 {
				charIdx := int(diff * 8)
				if charIdx < 0 {
					charIdx = 0
				}
				if charIdx > 7 {
					charIdx = 7
				}
				line[col] = blockChars[charIdx]
			} else {
				line[col] = ' '
			}
		}

		lineStr := " " + string(line) + " \n"
		if colorCode != "" {
			lineStr = colorCode + " " + string(line) + " " + ansiReset + "\n"
		}
		out.WriteString(lineStr)
	}
}

// renderSpectrum renders animated vertical equalizer bars.
func (r *Renderer) renderSpectrum(out *strings.Builder, width int, energy float64, frame int) {
	height := 5
	if r.cfg.Caps.SafeHeight() < 20 {
		height = 3
	}

	barCount := min2((width-4)/3, 30)
	if barCount < 5 {
		barCount = 5
	}

	// Prepare heights for each bar.
	heights := make([]float64, barCount)
	for i := 0; i < barCount; i++ {
		phase := float64(frame)*0.2 + float64(i)*0.4
		bounce := math.Sin(phase)*0.4 + 0.6
		heights[i] = bounce * energy * float64(height)
	}

	blockChars := []rune{' ', ' ', '▂', '▃', '▄', '▅', '▆', '▇'}
	if r.cfg.NoUnicode {
		blockChars = []rune{' ', ' ', '.', ',', '-', '~', '=', '#'}
	}

	colorCode := ""
	if r.state.Mood != nil {
		colorCode = r.getWaveColor(r.state.Mood.Label)
	}

	// Render row by row from top to bottom.
	for row := height - 1; row >= 0; row-- {
		var sb strings.Builder
		sb.WriteString("  ")
		for i := 0; i < barCount; i++ {
			h := heights[i]
			diff := h - float64(row)
			if diff >= 1.0 {
				if r.cfg.NoUnicode {
					sb.WriteString("## ")
				} else {
					sb.WriteString("██ ")
				}
			} else if diff > 0.0 {
				charIdx := int(diff * 8)
				if charIdx < 0 {
					charIdx = 0
				}
				if charIdx > 7 {
					charIdx = 7
				}
				rChar := blockChars[charIdx]
				sb.WriteRune(rChar)
				sb.WriteRune(rChar)
				sb.WriteRune(' ')
			} else {
				sb.WriteString("   ")
			}
		}

		paddedLine := centerPad(sb.String(), width)
		if colorCode != "" {
			paddedLine = colorCode + paddedLine + ansiReset
		}
		out.WriteString(paddedLine + "\n")
	}

	bottomLine := strings.Repeat("▀", barCount*3)
	if r.cfg.NoUnicode {
		bottomLine = strings.Repeat("=", barCount*3)
	}
	paddedBottom := centerPad(bottomLine, width)
	if colorCode != "" {
		paddedBottom = colorCode + paddedBottom + ansiReset
	}
	out.WriteString(paddedBottom + "\n")
}

// renderPulse renders ambient pulse rings.
func (r *Renderer) renderPulse(out *strings.Builder, width int, energy float64, frame int) {
	pulse := math.Abs(math.Sin(float64(frame) * 0.15))
	intensity := int(pulse * energy * float64(width-16))
	if intensity < 2 {
		intensity = 2
	}

	leftRing := ""
	rightRing := ""

	ringChars := []string{"·", "•", "o", "O", "◎", "●"}
	if r.cfg.NoUnicode {
		ringChars = []string{".", "-", "o", "O", "@", "#"}
	}

	idx := (frame / 2) % len(ringChars)
	char := ringChars[idx]

	for i := 0; i < intensity/6; i++ {
		leftRing = char + " " + leftRing
		rightRing = rightRing + " " + char
	}

	colorCode := ""
	if r.state.Mood != nil {
		colorCode = r.getWaveColor(r.state.Mood.Label)
	}

	pulseStr := fmt.Sprintf("%s  ◆  %s", leftRing, rightRing)
	paddedPulse := centerPad(pulseStr, width)
	if colorCode != "" {
		paddedPulse = colorCode + paddedPulse + ansiReset
	}
	out.WriteString(paddedPulse + "\n")

	paddedEnergy := centerPad(fmt.Sprintf("~  energy: %.0f%%  ~", energy*100), width)
	if colorCode != "" {
		paddedEnergy = colorCode + paddedEnergy + ansiReset
	}
	out.WriteString(paddedEnergy + "\n")
}

// renderMinimalVisual renders a compact status visual.
func (r *Renderer) renderMinimalVisual(out *strings.Builder, width int, state RenderState) {
	if state.Mood != nil {
		conf := fmt.Sprintf("mood: %s (%.0f%%)", state.Mood.Label, state.Mood.Confidence*100)
		r.writeDim(out, centerPad(conf, width))
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// Rendering helpers
// ──────────────────────────────────────────────────────────────────────────────

func (r *Renderer) getMoodColor(moodLabel mood.Label) string {
	if r.cfg.NoColor {
		return ""
	}
	switch moodLabel {
	case "calm", "minimal":
		return "\033[1;34m" // Bright Blue
	case "focused", "debugging":
		return "\033[1;33m" // Bright Yellow
	case "intense", "sprint":
		return "\033[1;31m" // Bright Red
	case "chaotic", "experimental":
		return "\033[1;36m" // Bright Cyan
	case "late-night":
		return "\033[1;35m" // Bright Magenta
	default:
		return "\033[1;32m" // Bright Green
	}
}

// ThemeColors represents terminal colors for visual theme rendering.
type ThemeColors struct {
	Primary   string
	Secondary string
	Text      string
	Selected  string
}

// getThemeColors returns the active theme color escape sequences.
func (r *Renderer) getThemeColors() ThemeColors {
	if r.cfg.NoColor {
		return ThemeColors{}
	}

	theme := r.cfg.Theme
	switch theme {
	case "ocean":
		return ThemeColors{
			Primary:   "\033[1;36m", // Cyan
			Secondary: "\033[1;34m", // Blue
			Text:      "\033[37m",   // White
			Selected:  "\033[1;36m", // Cyan
		}
	case "neon":
		return ThemeColors{
			Primary:   "\033[1;35m", // Magenta/Purple
			Secondary: "\033[1;36m", // Cyan
			Text:      "\033[37m",   // White
			Selected:  "\033[1;35m", // Magenta
		}
	case "sunset":
		return ThemeColors{
			Primary:   "\033[1;31m", // Red
			Secondary: "\033[1;33m", // Yellow/Orange
			Text:      "\033[37m",   // White
			Selected:  "\033[1;31m", // Red
		}
	case "matrix":
		return ThemeColors{
			Primary:   "\033[1;32m", // Green
			Secondary: "\033[32m",   // Dim Green
			Text:      "\033[37m",   // White
			Selected:  "\033[1;32m", // Green
		}
	case "lavender":
		return ThemeColors{
			Primary:   "\033[38;5;141m", // Lavender
			Secondary: "\033[1;35m",     // Purple
			Text:      "\033[37m",       // White
			Selected:  "\033[38;5;141m",
		}
	case "dark", "ash":
		return ThemeColors{
			Primary:   "\033[38;5;252m", // Bright ash
			Secondary: "\033[38;5;240m", // Slate gray
			Text:      "\033[38;5;248m", // Medium ash
			Selected:  "\033[1;37m",     // Bold white
		}
	case "ghost":
		return ThemeColors{
			Primary:   "\033[38;5;244m", // Dim gray
			Secondary: "\033[38;5;236m", // Very dark gray
			Text:      "\033[38;5;240m", // Muted gray
			Selected:  "\033[1;30m",     // Bold dark gray
		}
	case "monochrome":
		fallthrough
	default:
		return ThemeColors{
			Primary:   "\033[1;37m", // Bold White
			Secondary: "\033[1;30m", // Dark Gray
			Text:      "\033[37m",   // White
			Selected:  "\033[1;37m", // Bold White
		}
	}
}

// getWaveColor returns wave colors influenced by active theme settings.
func (r *Renderer) getWaveColor(moodLabel mood.Label) string {
	if r.cfg.NoColor {
		return ""
	}
	// Fall back to mood-based colors for monochrome themes, otherwise use theme primary
	if r.cfg.Theme == config.ThemeMonochrome || r.cfg.Theme == "" {
		return r.getMoodColor(moodLabel)
	}
	return r.getThemeColors().Primary
}

func (r *Renderer) writeLine(out *strings.Builder, s string) {
	out.WriteString(s + "\n")
}

func (r *Renderer) writeBold(out *strings.Builder, s string) {
	if r.cfg.NoColor {
		out.WriteString(s + "\n")
		return
	}
	out.WriteString(ansiBold + ansiGray100 + s + ansiReset + "\n")
}

func (r *Renderer) writeDim(out *strings.Builder, s string) {
	if r.cfg.NoColor {
		out.WriteString(s + "\n")
		return
	}
	out.WriteString(ansiDim + ansiGray70 + s + ansiReset + "\n")
}

func (r *Renderer) writeError(out *strings.Builder, s string) {
	if r.cfg.NoColor {
		out.WriteString("ERROR: " + s + "\n")
		return
	}
	out.WriteString(ansiBold + "\033[31m" + s + ansiReset + "\n")
}

func (r *Renderer) renderSeparator(out *strings.Builder, width int) {
	sep := ""
	if r.cfg.NoUnicode {
		sep = strings.Repeat("-", width)
	} else {
		sep = strings.Repeat("─", width)
	}
	r.writeDim(out, sep)
}

// centerPad pads s to width with spaces for centering.
func centerPad(s string, width int) string {
	// Strip ANSI codes for length measurement.
	displayLen := ansiLen(s)
	if displayLen >= width {
		return s
	}
	padding := (width - displayLen) / 2
	return strings.Repeat(" ", padding) + s
}

// ansiLen returns the display length of a string, ignoring ANSI escape codes.
func ansiLen(s string) int {
	l := 0
	inEscape := false
	for _, c := range s {
		if inEscape {
			if c == 'm' {
				inEscape = false
			}
			continue
		}
		if c == '\033' {
			inEscape = true
			continue
		}
		l++
	}
	return l
}

func min2(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ──────────────────────────────────────────────────────────────────────────────
// PrintOnce renders a single non-interactive frame to stdout.
// Used when TTY is not available or NoAnimation is set.
// ──────────────────────────────────────────────────────────────────────────────

// PrintStatus prints a compact status line. Safe for non-TTY output.
func PrintStatus(w io.Writer, state RenderState, noColor bool) {
	if state.Mood != nil {
		moodStr := fmt.Sprintf("mood: %s (%.0f%%)", state.Mood.Label, state.Mood.Confidence*100)
		if noColor {
			fmt.Fprintln(w, moodStr)
		} else {
			fmt.Fprintf(w, "%s%s%s\n", ansiBold, moodStr, ansiReset)
		}
	}
	if state.Candidate != nil {
		fmt.Fprintf(w, "now: %s\n", state.Candidate.DisplayName())
	}
}

// SetVisualTheme updates the theme dynamically. Thread-safe.
func (r *Renderer) SetVisualTheme(theme string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cfg.Theme = config.ThemeID(theme)
}

func (r *Renderer) renderWelcome(out *strings.Builder, width int, state RenderState, frame int) {
	items := strings.Split(state.Error, ";")
	selectedIdx := int(state.Progress)

	height := r.cfg.Caps.SafeHeight()
	if height < 15 {
		height = 15
	}

	bannerHeight := 6
	if r.cfg.Caps.SafeWidth() < 60 {
		bannerHeight = 1
	}
	contentHeight := bannerHeight + 2 + len(items) + 4 + 4
	topPadding := (height - contentHeight) / 2
	if topPadding < 1 {
		topPadding = 1
	}

	for i := 0; i < topPadding; i++ {
		out.WriteString("\n")
	}

	colors := r.getThemeColors()
	colorCode := ""
	if !r.cfg.NoColor {
		colorCode = colors.Primary
	}

	if r.cfg.Caps.SafeWidth() >= 60 {
		banner := renderDotMatrixBanner("MOODWAVE", r.cfg.NoUnicode)
		for _, line := range banner {
			centered := centerPad(line, width)
			if colorCode != "" {
				centered = colorCode + centered + ansiReset
			}
			out.WriteString(centered + "\n")
		}
	} else {
		centered := centerPad("◆ MOODWAVE ◆", width)
		if colorCode != "" {
			centered = ansiBold + colorCode + centered + ansiReset
		}
		out.WriteString(centered + "\n")
	}

	// Dynamic subtitle coloring
	subtitleText := "terminal mood music companion"
	if !r.cfg.NoColor {
		out.WriteString(centerPad(colors.Secondary+subtitleText+ansiReset, width) + "\n\n")
	} else {
		out.WriteString(centerPad(subtitleText, width) + "\n\n")
	}

	boxWidth := 36
	for _, item := range items {
		if len(item)+8 > boxWidth {
			boxWidth = len(item) + 8
		}
	}
	if boxWidth%2 != 0 {
		boxWidth++
	}

	boxTitle := fmt.Sprintf(" %s ", state.ScanMsg)
	headerLine := "┌" + strings.Repeat("─", (boxWidth-len(boxTitle))/2) + boxTitle + strings.Repeat("─", (boxWidth-len(boxTitle))/2) + "┐"
	if len(headerLine) > boxWidth+2 {
		headerLine = headerLine[:boxWidth+2]
	}
	
	// Dynamic header box border and title coloring
	if !r.cfg.NoColor {
		// Color the title of the box in primary, and borders in secondary
		lhs := "┌" + strings.Repeat("─", (boxWidth-len(boxTitle))/2)
		rhs := strings.Repeat("─", (boxWidth-len(boxTitle))/2) + "┐"
		coloredHeader := colors.Secondary + lhs + ansiReset + colors.Primary + ansiBold + boxTitle + ansiReset + colors.Secondary + rhs + ansiReset
		out.WriteString(centerPad(coloredHeader, width) + "\n")
	} else {
		out.WriteString(centerPad(headerLine, width) + "\n")
	}

	for idx, item := range items {
		content := item
		isSelected := idx == selectedIdx

		padding := boxWidth - len(content) - 4
		if padding < 0 {
			padding = 0
		}
		leftSpace := strings.Repeat(" ", padding/2)
		rightSpace := strings.Repeat(" ", padding-padding/2)

		var line string
		if isSelected {
			pointer := "▶ "
			if r.cfg.NoUnicode {
				pointer = "> "
			}
			rawLine := fmt.Sprintf("│  %s%s%s  │", pointer, content, leftSpace+rightSpace)
			if !r.cfg.NoColor {
				// Borders in secondary, pointer and selection in active selected theme color
				line = centerPad(colors.Secondary+"│  "+ansiReset+colors.Selected+ansiBold+pointer+content+leftSpace+rightSpace+ansiReset+colors.Secondary+"  │"+ansiReset, width)
			} else {
				line = centerPad(rawLine, width)
			}
		} else {
			rawLine := fmt.Sprintf("│     %s%s  │", content, leftSpace+rightSpace)
			if !r.cfg.NoColor {
				// Borders in secondary, standard items in theme text color
				line = centerPad(colors.Secondary+"│     "+ansiReset+colors.Text+content+ansiReset+colors.Secondary+leftSpace+rightSpace+"  │"+ansiReset, width)
			} else {
				line = centerPad(rawLine, width)
			}
		}
		out.WriteString(line + "\n")
	}

	bottomLine := "└" + strings.Repeat("─", boxWidth) + "┘"
	if !r.cfg.NoColor {
		out.WriteString(centerPad(colors.Secondary+bottomLine+ansiReset, width) + "\n\n")
	} else {
		out.WriteString(centerPad(bottomLine, width) + "\n\n")
	}

	helper := "[W/S / Arrow Keys] Navigate  •  [Enter/Space] Select  •  [Q/Esc] Back/Exit"
	if r.cfg.NoUnicode {
		helper = "[W/S] Navigate  *  [Enter] Select  *  [Q] Exit"
	}
	if !r.cfg.NoColor {
		out.WriteString(centerPad(colors.Secondary+helper+ansiReset, width) + "\n")
	} else {
		out.WriteString(centerPad(helper, width) + "\n")
	}

	remainingRows := height - (topPadding + contentHeight)
	for i := 0; i < remainingRows-1; i++ {
		out.WriteString("\n")
	}

	version := state.RepeatMode
	if version == "" {
		version = "dev"
	}
	metaLeft := "  🌊 moodwave dev companion"
	metaRight := fmt.Sprintf("v%s  ", version)
	metaWidth := width - len(metaLeft) - len(metaRight)
	if metaWidth < 2 {
		metaWidth = 2
	}
	metaLine := metaLeft + strings.Repeat(" ", metaWidth) + metaRight
	if !r.cfg.NoColor {
		out.WriteString(colors.Secondary + metaLine + ansiReset)
	} else {
		out.WriteString(metaLine)
	}
}
