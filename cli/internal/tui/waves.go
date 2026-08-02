package tui

import (
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ─────────────────────────────────────────────────────────────────────────────
// Color gradient — smooth rainbow spectrum from left to right.
// ─────────────────────────────────────────────────────────────────────────────

// Rainbow colors — full spectrum left to right like the reference image.
var RainbowColors = []string{
	"#ff0066", "#ff0055", "#ff0044", "#ff0033", "#ff0022",
	"#ff1100", "#ff2200", "#ff3300", "#ff4400", "#ff5500",
	"#ff6600", "#ff7700", "#ff8800", "#ff9900", "#ffaa00",
	"#ffbb00", "#ffcc00", "#ffdd00", "#ffee00", "#ffff00",
	"#ddff00", "#bbff00", "#88ff00", "#55ff00", "#22ff00",
	"#00ff22", "#00ff55", "#00ff88", "#00ffaa", "#00ffcc",
	"#00ffee", "#00eeff", "#00ccff", "#00aaff", "#0088ff",
	"#0066ff", "#0044ff", "#0033ff", "#0022ff", "#0011ff",
	"#1100ff", "#2200ff", "#3300ff", "#4400ff", "#5500ff",
	"#6600ff", "#7700ff", "#8800ff", "#9900ff", "#aa00ff",
	"#bb00ff", "#cc00ff", "#dd00ff", "#ee00ff", "#ff00ff",
	"#ff00dd", "#ff00bb", "#ff0099", "#ff0077", "#ff0066",
}

// AuroraPalette — pink/magenta to cyan.
var AuroraPalette = []string{
	"#ff6b9d", "#ff4ecb", "#e040fb", "#d946ef",
	"#c026d3", "#a855f7", "#8b5cf6", "#7c3aed",
	"#6366f1", "#4f46e5", "#3b82f6", "#2563eb",
	"#0ea5e9", "#06b6d4", "#14b8a6", "#22d3ee",
	"#67e8f9", "#a5f3fc",
}

// OceanPalette — teal to purple.
var OceanPalette = []string{
	"#5eead4", "#2dd4bf", "#14b8a6", "#0d9488",
	"#0891b2", "#0284c7", "#0369a1", "#1d4ed8",
	"#2563eb", "#3b82f6", "#4f46e5", "#6366f1",
	"#7c3aed", "#8b5cf6", "#a855f7", "#c084fc",
}

// FirePalette — yellow to magenta.
var FirePalette = []string{
	"#fef08a", "#fde047", "#facc15", "#eab308",
	"#f59e0b", "#d97706", "#ea580c", "#f97316",
	"#ef4444", "#dc2626", "#be123c", "#e11d48",
	"#db2777", "#c026d3", "#9333ea", "#7c3aed",
}

// ─────────────────────────────────────────────────────────────────────────────
// Doom-style fire — proper heat diffusion algorithm
// Based on the classic Doom fire effect (heat source at bottom, diffuses up)
// ─────────────────────────────────────────────────────────────────────────────

// fireState holds the persistent fire buffer between frames.
var fireState struct {
	buf    []int
	width  int
	height int
}

// Fire color ramp — from black (cool) to white (hottest)
var fireColors = []string{
	"#070707", "#1f0707", "#2f0f07", "#470f07",
	"#571707", "#671f07", "#772707", "#8f2f07",
	"#9f2f07", "#af3f07", "#bf4707", "#c74707",
	"#df4f07", "#df5707", "#df5707", "#d75f07",
	"#d7670f", "#cf6f0f", "#cf770f", "#cf7f0f",
	"#cf8717", "#c78717", "#c78f17", "#c7971f",
	"#bf9f1f", "#bf9f1f", "#bfa727", "#bfa727",
	"#bfaf2f", "#b7af2f", "#b7b72f", "#b7b737",
	"#cfcf6f", "#dfdf9f", "#efefc7", "#ffffff",
}

// Fire ASCII intensity ramp
var fireChars = []byte(" .:^*xsS#$")

func RenderDoomFire(width, height, frame int) string {
	if width < 10 {
		width = 10
	}
	if height < 4 {
		height = 4
	}

	size := width * height

	// Initialize or resize buffer
	if fireState.width != width || fireState.height != height || len(fireState.buf) != size+width+1 {
		fireState.width = width
		fireState.height = height
		fireState.buf = make([]int, size+width+1)
	}

	buf := fireState.buf

	// Seed hot pixels at the bottom row (random intensity bursts)
	for i := 0; i < width/6; i++ {
		pos := int(math.Abs(math.Sin(float64(frame)*0.1+float64(i)*1.618))) * (width - 1)
		pos = pos % width
		idx := pos + width*(height-1)
		if idx < len(buf) {
			buf[idx] = 35 + int(math.Abs(math.Sin(float64(frame+i)*0.3))*30)
			if buf[idx] > 65 {
				buf[idx] = 65
			}
		}
	}

	// Additional random seeds for variety
	for i := 0; i < width/4; i++ {
		// Use deterministic pseudo-random based on frame
		seed := (frame*7 + i*13) % width
		idx := seed + width*(height-1)
		if idx < len(buf) {
			buf[idx] = 50 + (frame*3+i*17)%15
		}
	}

	// Heat diffusion pass — each pixel averages with neighbors and cools
	for i := 0; i < size; i++ {
		right := i + 1
		below := i + width
		belowRight := i + width + 1
		if right >= len(buf) {
			right = i
		}
		if below >= len(buf) {
			below = i
		}
		if belowRight >= len(buf) {
			belowRight = i
		}

		v := (buf[i] + buf[right] + buf[below] + buf[belowRight]) >> 2
		if v > 0 {
			v-- // cooling
		}
		buf[i] = v
	}

	// Render to string
	lines := make([]string, height)
	for row := 0; row < height; row++ {
		var b strings.Builder
		for x := 0; x < width; x++ {
			idx := row*width + x
			v := buf[idx]
			if v < 0 {
				v = 0
			}

			// Map value to character
			charIdx := v * len(fireChars) / 66
			if charIdx >= len(fireChars) {
				charIdx = len(fireChars) - 1
			}
			if charIdx < 0 {
				charIdx = 0
			}
			ch := string(fireChars[charIdx])

			// Map value to color
			colorIdx := v * len(fireColors) / 66
			if colorIdx >= len(fireColors) {
				colorIdx = len(fireColors) - 1
			}
			if colorIdx < 0 {
				colorIdx = 0
			}

			if v < 2 {
				b.WriteString(" ")
			} else {
				color := lipgloss.Color(fireColors[colorIdx])
				style := lipgloss.NewStyle().Foreground(color)
				b.WriteString(style.Render(ch))
			}
		}
		lines[row] = b.String()
	}

	return strings.Join(lines, "\n")
}

// ─────────────────────────────────────────────────────────────────────────────
// Main visualizer — full-width mirrored vertical bars, end to end.
// This is the exact look from the reference: thin bars, variable height,
// mirrored around center, rainbow gradient across full width.
// ─────────────────────────────────────────────────────────────────────────────

// RenderFullVisualizer renders the end-to-end mirrored waveform.
// Each column is one bar. Height varies per frame. Colors gradient left to right.
// totalHeight is the number of character rows (should be 8-12 for best look).
func RenderFullVisualizer(width, totalHeight, frame int, energy float64, palette []string) string {
	if width < 20 {
		width = 20
	}
	if totalHeight < 4 {
		totalHeight = 4
	}
	if len(palette) == 0 {
		palette = RainbowColors
	}

	halfHeight := totalHeight / 2
	t := float64(frame) * 0.06

	// Pre-compute amplitudes for each column
	amplitudes := make([]float64, width)
	for x := 0; x < width; x++ {
		pos := float64(x) / float64(width)

		// Multiple layered sine waves at different frequencies for organic movement
		a := math.Sin(pos*math.Pi*6.0+t*1.0) * 0.25
		b := math.Sin(pos*math.Pi*12.0+t*1.7) * 0.18
		c := math.Sin(pos*math.Pi*18.0+t*2.3) * 0.12
		d := math.Sin(pos*math.Pi*3.0+t*0.5) * 0.20
		e := math.Sin(pos*math.Pi*24.0+t*3.1) * 0.08
		f := math.Sin(pos*math.Pi*9.0+t*1.3+float64(x)*0.1) * 0.10

		// Combine — produces values roughly [-0.93, +0.93]
		combined := a + b + c + d + e + f

		// Normalize to [0, 1]
		val := (combined + 1.0) / 2.0

		// Apply energy multiplier
		val *= energy

		// Add slight random-feel variation using position-dependent phase
		val += math.Sin(float64(x)*1.618+t*2.7) * 0.05

		// Clamp
		if val < 0.02 {
			val = 0.02
		}
		if val > 1.0 {
			val = 1.0
		}

		amplitudes[x] = val
	}

	// Build each row
	lines := make([]string, totalHeight)
	for row := 0; row < totalHeight; row++ {
		var b strings.Builder

		// Distance from center line (0 at center, halfHeight at top/bottom)
		distFromCenter := row - halfHeight
		if distFromCenter < 0 {
			distFromCenter = -distFromCenter
		}
		// Threshold: what amplitude is needed to fill this row
		threshold := float64(distFromCenter) / float64(halfHeight)

		for x := 0; x < width; x++ {
			amp := amplitudes[x]

			if amp >= threshold {
				// Get color for this x position
				colorIdx := int(float64(x) / float64(width) * float64(len(palette)-1))
				if colorIdx >= len(palette) {
					colorIdx = len(palette) - 1
				}

				// Slightly dim bars that are further from center for depth
				color := lipgloss.Color(palette[colorIdx])
				style := lipgloss.NewStyle().Foreground(color)

				// Use full block for solid bars
				b.WriteString(style.Render("│"))
			} else {
				b.WriteString(" ")
			}
		}

		lines[row] = b.String()
	}

	return strings.Join(lines, "\n")
}

// RenderFullSpectrumBars renders end-to-end bars rising from bottom only.
// Classic equalizer look — bars go up, no mirror.
func RenderFullSpectrumBars(width, height, frame int, energy float64, palette []string) string {
	if width < 20 {
		width = 20
	}
	if height < 3 {
		height = 3
	}
	if len(palette) == 0 {
		palette = RainbowColors
	}

	t := float64(frame) * 0.06

	// Compute amplitudes
	amplitudes := make([]float64, width)
	for x := 0; x < width; x++ {
		pos := float64(x) / float64(width)

		a := math.Sin(pos*math.Pi*6.0+t*1.0) * 0.25
		b := math.Sin(pos*math.Pi*12.0+t*1.7) * 0.18
		c := math.Sin(pos*math.Pi*18.0+t*2.3) * 0.12
		d := math.Sin(pos*math.Pi*3.0+t*0.5) * 0.20
		e := math.Sin(pos*math.Pi*24.0+t*3.1) * 0.08

		combined := a + b + c + d + e
		val := (combined + 1.0) / 2.0
		val *= energy

		if val < 0.02 {
			val = 0.02
		}
		if val > 1.0 {
			val = 1.0
		}
		amplitudes[x] = val
	}

	lines := make([]string, height)
	for row := 0; row < height; row++ {
		var b strings.Builder
		// Row 0 = top, row height-1 = bottom
		threshold := float64(height-1-row) / float64(height)

		for x := 0; x < width; x++ {
			if amplitudes[x] >= threshold {
				colorIdx := int(float64(x) / float64(width) * float64(len(palette)-1))
				if colorIdx >= len(palette) {
					colorIdx = len(palette) - 1
				}
				color := lipgloss.Color(palette[colorIdx])
				style := lipgloss.NewStyle().Foreground(color)
				b.WriteString(style.Render("█"))
			} else {
				b.WriteString(" ")
			}
		}
		lines[row] = b.String()
	}

	return strings.Join(lines, "\n")
}

// ─────────────────────────────────────────────────────────────────────────────
// Inline single-row wave for the home screen.
// ─────────────────────────────────────────────────────────────────────────────

// RenderInlineWave renders a single-row colored wave spanning full width.
func RenderInlineWave(width int, frame int, palette []string) string {
	if width < 10 {
		width = 10
	}
	if len(palette) == 0 {
		palette = RainbowColors
	}

	chars := []string{"▁", "▂", "▃", "▄", "▅", "▆", "▇", "█"}
	t := float64(frame) * 0.08
	var b strings.Builder

	for x := 0; x < width; x++ {
		pos := float64(x) / float64(width)

		val := math.Sin(pos*math.Pi*4.0+t) * 0.30
		val += math.Sin(pos*math.Pi*8.0+t*1.5) * 0.20
		val += math.Sin(pos*math.Pi*12.0+t*0.7) * 0.15
		val += math.Sin(pos*math.Pi*2.0+t*0.4) * 0.15
		val = (val + 1.0) / 2.0

		idx := int(val * float64(len(chars)-1))
		if idx < 0 {
			idx = 0
		}
		if idx >= len(chars) {
			idx = len(chars) - 1
		}

		colorIdx := int(pos * float64(len(palette)-1))
		if colorIdx >= len(palette) {
			colorIdx = len(palette) - 1
		}
		color := lipgloss.Color(palette[colorIdx])
		style := lipgloss.NewStyle().Foreground(color)
		b.WriteString(style.Render(chars[idx]))
	}

	return b.String()
}

// ─────────────────────────────────────────────────────────────────────────────
// Compatibility wrappers used by app.go
// ─────────────────────────────────────────────────────────────────────────────

// GetPaletteForMood returns the appropriate color palette for a mood.
func GetPaletteForMood(mood string) []string {
	switch mood {
	case "focused", "calm", "minimal":
		return OceanPalette
	case "intense", "sprint":
		return FirePalette
	case "experimental", "debugging", "late-night":
		return AuroraPalette
	default:
		return RainbowColors
	}
}

// RenderWave wraps InlineWave for backward compat.
func RenderWave(width int, frame int, amplitude float64, color lipgloss.Color) string {
	_ = amplitude
	_ = color
	return RenderInlineWave(width, frame, RainbowColors)
}

// RenderSpectrum wraps the full spectrum bars.
func RenderSpectrum(width int, frame int, energy float64, color lipgloss.Color) string {
	_ = color
	return RenderFullSpectrumBars(width, 6, frame, energy, RainbowColors)
}

// RenderPulse wraps the mirrored visualizer.
func RenderPulse(width int, frame int, intensity float64, color lipgloss.Color) string {
	_ = color
	return RenderFullVisualizer(width, 8, frame, intensity, RainbowColors)
}

// RenderMiniWave is a small inline wave for the status area.
func RenderMiniWave(width int, frame int, color lipgloss.Color) string {
	_ = color
	return RenderInlineWave(width, frame, OceanPalette)
}

// ─────────────────────────────────────────────────────────────────────────────
// City Night — buildings silhouette with twinkling windows and stars
// ─────────────────────────────────────────────────────────────────────────────

func RenderCityNight(width, height, frame int) string {
	if width < 20 {
		width = 20
	}
	if height < 5 {
		height = 5
	}

	t := float64(frame) * 0.04
	lines := make([]string, height)

	// Generate building heights
	buildings := make([]int, width)
	for x := 0; x < width; x++ {
		pos := float64(x) / float64(width)
		h := int((math.Sin(pos*math.Pi*3.0+0.5)*0.3 +
			math.Sin(pos*math.Pi*7.0+1.2)*0.2 +
			0.5) * float64(height))
		if h < 1 {
			h = 1
		}
		if h >= height {
			h = height - 1
		}
		buildings[x] = h
	}

	skyColor := lipgloss.Color("#1a1a3e")
	buildColor := lipgloss.Color("#2d2d4e")
	windowOn := lipgloss.Color("#ffdd44")
	windowOff := lipgloss.Color("#3a3a5e")
	starColor := lipgloss.Color("#ffffff")

	for row := 0; row < height; row++ {
		var b strings.Builder
		for x := 0; x < width; x++ {
			buildingTop := height - buildings[x]
			if row < buildingTop {
				// Sky
				starPhase := math.Sin(float64(x*17+row*31)*0.1 + t*2.0)
				if starPhase > 0.92 {
					b.WriteString(lipgloss.NewStyle().Foreground(starColor).Render("·"))
				} else {
					b.WriteString(lipgloss.NewStyle().Foreground(skyColor).Render(" "))
				}
			} else {
				// Building
				windowPhase := math.Sin(float64(x*7+row*13)*0.3 + t*0.5)
				if x%3 == 1 && row%2 == 0 && windowPhase > 0.2 {
					b.WriteString(lipgloss.NewStyle().Foreground(windowOn).Render("▪"))
				} else if x%3 == 1 && row%2 == 0 {
					b.WriteString(lipgloss.NewStyle().Foreground(windowOff).Render("▪"))
				} else {
					b.WriteString(lipgloss.NewStyle().Foreground(buildColor).Render("█"))
				}
			}
		}
		lines[row] = b.String()
	}

	return strings.Join(lines, "\n")
}

// ─────────────────────────────────────────────────────────────────────────────
// Ocean Waves — calm beach with layered water waves
// ─────────────────────────────────────────────────────────────────────────────

func RenderOceanWaves(width, height, frame int) string {
	if width < 20 {
		width = 20
	}
	if height < 5 {
		height = 5
	}

	t := float64(frame) * 0.05
	lines := make([]string, height)

	// Color layers from sky to deep water
	skyColors := []string{"#1a1a3e", "#1e2a4e", "#223a5e"}
	waterColors := []string{"#0077b6", "#0096c7", "#00b4d8", "#48cae4", "#90e0ef", "#ade8f4"}

	for row := 0; row < height; row++ {
		var b strings.Builder
		rowRatio := float64(row) / float64(height)

		for x := 0; x < width; x++ {
			pos := float64(x) / float64(width)

			if rowRatio < 0.3 {
				// Sky with moon reflection
				colorIdx := int(rowRatio / 0.3 * float64(len(skyColors)-1))
				if colorIdx >= len(skyColors) {
					colorIdx = len(skyColors) - 1
				}
				moonDist := math.Abs(pos - 0.7)
				if moonDist < 0.02 && rowRatio < 0.15 {
					b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#fffacd")).Render("◯"))
				} else {
					b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(skyColors[colorIdx])).Render("░"))
				}
			} else {
				// Water with wave movement
				waveOffset := math.Sin(pos*math.Pi*4.0+t+rowRatio*3.0) * 0.3
				waveOffset += math.Sin(pos*math.Pi*8.0+t*1.5) * 0.15

				colorIdx := int((rowRatio + waveOffset) * float64(len(waterColors)-1))
				if colorIdx < 0 {
					colorIdx = 0
				}
				if colorIdx >= len(waterColors) {
					colorIdx = len(waterColors) - 1
				}

				// Wave crests
				crest := math.Sin(pos*math.Pi*6.0+t*1.2+rowRatio*2.0)
				char := "~"
				if crest > 0.7 {
					char = "≈"
				} else if crest < -0.5 {
					char = "∽"
				}

				b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(waterColors[colorIdx])).Render(char))
			}
		}
		lines[row] = b.String()
	}

	return strings.Join(lines, "\n")
}

// ─────────────────────────────────────────────────────────────────────────────
// Starfield — deep space with moving stars and nebula colors
// ─────────────────────────────────────────────────────────────────────────────

func RenderStarfield(width, height, frame int) string {
	if width < 20 {
		width = 20
	}
	if height < 5 {
		height = 5
	}

	t := float64(frame) * 0.03
	lines := make([]string, height)

	nebulaColors := []string{"#2d1b69", "#4a1c7a", "#6b2fa0", "#8b3fa0", "#5c2d91", "#3d1e6d"}
	starChars := []string{"·", "✦", "⋆", "*", "✧", "˚"}

	for row := 0; row < height; row++ {
		var b strings.Builder
		for x := 0; x < width; x++ {
			// Unique hash per position for consistent star placement
			hash := float64((x*127 + row*311 + frame/8) % 1000)
			pos := float64(x) / float64(width)
			rowPos := float64(row) / float64(height)

			// Nebula background
			nebulaPhase := math.Sin(pos*math.Pi*2.0+t*0.3) * math.Cos(rowPos*math.Pi+t*0.2)
			colorIdx := int((nebulaPhase + 1.0) / 2.0 * float64(len(nebulaColors)-1))
			if colorIdx < 0 {
				colorIdx = 0
			}
			if colorIdx >= len(nebulaColors) {
				colorIdx = len(nebulaColors) - 1
			}

			// Stars appear based on hash
			starPhase := math.Sin(hash*0.01 + t*2.0)
			if hash/1000.0 < 0.06 && starPhase > 0.3 {
				// Bright star
				charIdx := int(hash) % len(starChars)
				brightness := []string{"#ffffff", "#ffffee", "#eeeeff", "#ffddff"}
				bIdx := int(hash*0.7) % len(brightness)
				b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(brightness[bIdx])).Render(starChars[charIdx]))
			} else if hash/1000.0 < 0.12 {
				// Dim star
				b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#555577")).Render("·"))
			} else {
				// Nebula background
				b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(nebulaColors[colorIdx])).Render("░"))
			}
		}
		lines[row] = b.String()
	}

	return strings.Join(lines, "\n")
}
