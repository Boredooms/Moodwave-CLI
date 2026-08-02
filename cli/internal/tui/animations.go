package tui

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ─────────────────────────────────────────────────────────────────────────────
// Matrix Rain — falling green characters
// ─────────────────────────────────────────────────────────────────────────────

func RenderMatrix(width, height, frame int) string {
	lines := make([]string, height)
	t := float64(frame) * 0.12
	matrixChars := "ｱｲｳｴｵｶｷｸｹｺ0123456789"

	for row := 0; row < height; row++ {
		var b strings.Builder
		for x := 0; x < width; x++ {
			// Each column has a "drop" position
			speed := 1.0 + math.Sin(float64(x)*0.7)*0.5
			dropPos := math.Mod(t*speed+float64(x*37%100), float64(height+8))
			dist := dropPos - float64(row)

			if dist >= 0 && dist < 1.0 {
				// Head of drop — bright white/green
				ch := string([]rune(matrixChars)[(frame+x*7+row*3)%len([]rune(matrixChars))])
				b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Render(ch))
			} else if dist >= 1.0 && dist < 12.0 {
				// Trail — fading green
				fade := 1.0 - dist/12.0
				green := int(180*fade) + 40
				if green > 255 {
					green = 255
				}
				color := lipgloss.Color(fmt.Sprintf("#00%02x00", green))
				ch := string([]rune(matrixChars)[(frame/2+x*13+row*7)%len([]rune(matrixChars))])
				b.WriteString(lipgloss.NewStyle().Foreground(color).Render(ch))
			} else {
				b.WriteString(" ")
			}
		}
		lines[row] = b.String()
	}
	return strings.Join(lines, "\n")
}

// ─────────────────────────────────────────────────────────────────────────────
// Rain — falling drops with splash
// ─────────────────────────────────────────────────────────────────────────────

func RenderRain(width, height, frame int) string {
	lines := make([]string, height)
	t := float64(frame)

	for row := 0; row < height; row++ {
		var b strings.Builder
		for x := 0; x < width; x++ {
			hash := (x*131 + row*97 + int(t)*3) % 1000
			dropActive := math.Sin(float64(hash)*0.01+t*0.15) > 0.85

			if row == height-1 && dropActive {
				// Splash at bottom
				b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#87ceeb")).Render("·"))
			} else if dropActive && (int(t)+hash)%height == row {
				b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#4fc3f7")).Render("│"))
			} else if dropActive && (int(t)+hash+1)%height == row {
				b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#29b6f6")).Render("╎"))
			} else {
				b.WriteString(" ")
			}
		}
		lines[row] = b.String()
	}
	return strings.Join(lines, "\n")
}

// ─────────────────────────────────────────────────────────────────────────────
// Plasma — overlapping sine waves with smooth color
// ─────────────────────────────────────────────────────────────────────────────

func RenderPlasma(width, height, frame int) string {
	lines := make([]string, height)
	t := float64(frame) * 0.04
	chars := []string{"░", "▒", "▓", "█"}

	for row := 0; row < height; row++ {
		var b strings.Builder
		for x := 0; x < width; x++ {
			fx := float64(x) / float64(width) * 4.0
			fy := float64(row) / float64(height) * 4.0

			v := math.Sin(fx + t)
			v += math.Sin(fy + t*0.7)
			v += math.Sin((fx+fy)*0.7 + t*1.3)
			v += math.Sin(math.Sqrt(fx*fx+fy*fy) + t*0.5)
			v = (v + 4.0) / 8.0 // normalize to [0,1]

			// Map to hue
			r := int(math.Sin(v*math.Pi*2.0)*127 + 128)
			g := int(math.Sin(v*math.Pi*2.0+2.094)*127 + 128)
			bl := int(math.Sin(v*math.Pi*2.0+4.189)*127 + 128)
			color := lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", r, g, bl))
			ch := chars[int(v*float64(len(chars)-1))]
			b.WriteString(lipgloss.NewStyle().Foreground(color).Render(ch))
		}
		lines[row] = b.String()
	}
	return strings.Join(lines, "\n")
}

// ─────────────────────────────────────────────────────────────────────────────
// Aurora Borealis — layered color curtains
// ─────────────────────────────────────────────────────────────────────────────

func RenderAurora(width, height, frame int) string {
	lines := make([]string, height)
	t := float64(frame) * 0.03

	for row := 0; row < height; row++ {
		var b strings.Builder
		rowRatio := float64(row) / float64(height)

		for x := 0; x < width; x++ {
			pos := float64(x) / float64(width)

			// Multiple curtain layers
			curtain1 := math.Sin(pos*math.Pi*3.0+t) * math.Cos(rowRatio*math.Pi+t*0.5)
			curtain2 := math.Sin(pos*math.Pi*5.0+t*1.3+1.0) * math.Cos(rowRatio*math.Pi*0.7+t*0.3)
			intensity := (curtain1 + curtain2 + 2.0) / 4.0

			// Fade toward bottom
			intensity *= (1.0 - rowRatio*0.7)

			if intensity < 0.15 {
				b.WriteString(" ")
			} else {
				// Green/purple/cyan colors
				g := int(intensity * 200)
				bl := int(math.Sin(pos*math.Pi*2.0+t)*80 + 100)
				r := int(math.Sin(pos*math.Pi*4.0+t*1.5)*60 + 40)
				if g > 255 { g = 255 }
				if bl > 255 { bl = 255 }
				if r > 255 { r = 255 }
				color := lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", r, g, bl))
				ch := "░"
				if intensity > 0.5 { ch = "▒" }
				if intensity > 0.7 { ch = "▓" }
				b.WriteString(lipgloss.NewStyle().Foreground(color).Render(ch))
			}
		}
		lines[row] = b.String()
	}
	return strings.Join(lines, "\n")
}

// ─────────────────────────────────────────────────────────────────────────────
// Snow — falling snowflakes with accumulation
// ─────────────────────────────────────────────────────────────────────────────

func RenderSnow(width, height, frame int) string {
	lines := make([]string, height)
	t := float64(frame) * 0.08
	flakes := []string{"·", "❄", "❅", "✦", "*", "⊹"}

	for row := 0; row < height; row++ {
		var b strings.Builder
		for x := 0; x < width; x++ {
			hash := float64((x*73 + row*137) % 997)
			// Wind drift
			drift := math.Sin(t*0.3+hash*0.01) * 2.0
			dropY := math.Mod(t*0.7+hash*0.05, float64(height+4))

			if row == height-1 {
				// Snow accumulation at bottom
				accum := math.Sin(float64(x)*0.2+0.5)*0.3 + 0.5
				if accum > 0.4 {
					b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#e0e0e0")).Render("▁"))
				} else {
					b.WriteString(" ")
				}
			} else if math.Abs(dropY-float64(row)) < 0.8 && math.Abs(drift-float64(x%5)) < 1.5 {
				idx := int(hash) % len(flakes)
				b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Render(flakes[idx]))
			} else {
				b.WriteString(" ")
			}
		}
		lines[row] = b.String()
	}
	return strings.Join(lines, "\n")
}

// ─────────────────────────────────────────────────────────────────────────────
// Fireflies — blinking warm dots in darkness
// ─────────────────────────────────────────────────────────────────────────────

func RenderFireflies(width, height, frame int) string {
	lines := make([]string, height)
	t := float64(frame) * 0.05
	numFireflies := (width * height) / 40 // density based on area
	if numFireflies < 8 {
		numFireflies = 8
	}

	// Pre-compute firefly positions
	type firefly struct {
		x, y    float64
		visible bool
	}
	flies := make([]firefly, numFireflies)
	for i := range flies {
		seed := float64(i) * 1.618033988 // golden ratio spacing
		// Distribute across the full area
		baseX := math.Mod(seed*float64(width)*0.73, float64(width))
		baseY := math.Mod(seed*float64(height)*0.91, float64(height))
		// Drift over time
		driftX := math.Sin(t*0.3+seed*2.0) * 3.0
		driftY := math.Cos(t*0.2+seed*1.5) * 2.0
		// Blink cycle
		blink := math.Sin(t*1.2 + seed*3.7)
		flies[i] = firefly{
			x:       baseX + driftX,
			y:       baseY + driftY,
			visible: blink > 0.2,
		}
	}

	for row := 0; row < height; row++ {
		var b strings.Builder
		for x := 0; x < width; x++ {
			drawn := false
			for _, f := range flies {
				if !f.visible {
					continue
				}
				dx := math.Abs(float64(x) - f.x)
				dy := math.Abs(float64(row) - f.y)
				if dx < 1.0 && dy < 0.8 {
					// Core glow
					b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#ffdd44")).Render("✦"))
					drawn = true
					break
				} else if dx < 2.0 && dy < 1.5 {
					b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#554400")).Render("·"))
					drawn = true
					break
				}
			}
			if !drawn {
				b.WriteString(" ")
			}
		}
		lines[row] = b.String()
	}
	return strings.Join(lines, "\n")
}

// ─────────────────────────────────────────────────────────────────────────────
// Lava Lamp — blobs rising and merging
// ─────────────────────────────────────────────────────────────────────────────

func RenderLava(width, height, frame int) string {
	lines := make([]string, height)
	t := float64(frame) * 0.03

	for row := 0; row < height; row++ {
		var b strings.Builder
		fy := float64(row) / float64(height)

		for x := 0; x < width; x++ {
			fx := float64(x) / float64(width)

			// 5 smaller metaballs spread across the area
			d1 := math.Sqrt(math.Pow(fx-0.2-math.Sin(t)*0.08, 2) + math.Pow(fy-0.3-math.Cos(t*0.7)*0.15, 2))
			d2 := math.Sqrt(math.Pow(fx-0.7+math.Cos(t*0.6)*0.1, 2) + math.Pow(fy-0.4+math.Sin(t*0.5)*0.12, 2))
			d3 := math.Sqrt(math.Pow(fx-0.5-math.Sin(t*0.9)*0.12, 2) + math.Pow(fy-0.7-math.Cos(t*0.4)*0.1, 2))
			d4 := math.Sqrt(math.Pow(fx-0.3+math.Cos(t*1.1)*0.06, 2) + math.Pow(fy-0.8+math.Sin(t*0.8)*0.08, 2))
			d5 := math.Sqrt(math.Pow(fx-0.8-math.Sin(t*0.7)*0.07, 2) + math.Pow(fy-0.2-math.Cos(t*1.0)*0.1, 2))

			field := 0.04/d1 + 0.035/d2 + 0.03/d3 + 0.025/d4 + 0.03/d5

			if field > 1.0 {
				b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#ffcc00")).Render("█"))
			} else if field > 0.7 {
				b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#ff4444")).Render("▓"))
			} else if field > 0.5 {
				b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#cc2200")).Render("▒"))
			} else if field > 0.35 {
				b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#551100")).Render("░"))
			} else {
				b.WriteString(" ")
			}
		}
		lines[row] = b.String()
	}
	return strings.Join(lines, "\n")
}

// ─────────────────────────────────────────────────────────────────────────────
// DNA Helix — rotating double helix
// ─────────────────────────────────────────────────────────────────────────────

func RenderDNA(width, height, frame int) string {
	lines := make([]string, height)
	t := float64(frame) * 0.07
	bases := []string{"A", "T", "C", "G"}
	baseColors := []lipgloss.Color{"#ff4444", "#44ff44", "#4444ff", "#ffff44"}

	centerX := float64(width) / 2.0
	amplitude := float64(width) * 0.25
	if amplitude > 30 {
		amplitude = 30
	}

	for row := 0; row < height; row++ {
		var b strings.Builder
		phase := float64(row)*0.4 + t

		// Exact strand positions (integer)
		ix1 := int(centerX + math.Sin(phase)*amplitude)
		ix2 := int(centerX + math.Sin(phase+math.Pi)*amplitude)

		for x := 0; x < width; x++ {
			if x == ix1 {
				idx := (row + frame/5) % len(bases)
				b.WriteString(lipgloss.NewStyle().Foreground(baseColors[idx]).Bold(true).Render(bases[idx]))
			} else if x == ix2 {
				idx := (row + frame/5 + 2) % len(bases)
				b.WriteString(lipgloss.NewStyle().Foreground(baseColors[idx]).Bold(true).Render(bases[idx]))
			} else if x > min2(ix1, ix2) && x < max2(ix1, ix2) && row%2 == 0 {
				// Base pair connection dots between strands
				b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#333344")).Render("·"))
			} else {
				b.WriteString(" ")
			}
		}
		lines[row] = b.String()
	}
	return strings.Join(lines, "\n")
}

func min2(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max2(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// ─────────────────────────────────────────────────────────────────────────────
// Smoke — rising wisps
// ─────────────────────────────────────────────────────────────────────────────

func RenderSmoke(width, height, frame int) string {
	lines := make([]string, height)
	t := float64(frame) * 0.04
	smokeChars := []string{".", ":", "∘", "○", "◌", "◯"}

	for row := 0; row < height; row++ {
		var b strings.Builder
		fy := float64(row) / float64(height)

		for x := 0; x < width; x++ {
			fx := float64(x) / float64(width)

			// Perlin-like turbulence
			turb := math.Sin(fx*8.0+t+fy*3.0)*0.3 +
				math.Sin(fx*12.0+t*1.5-fy*5.0)*0.2 +
				math.Sin(fx*4.0+t*0.5+fy*2.0)*0.2

			// Smoke rises from center-bottom
			dist := math.Abs(fx-0.5+turb*0.2) / (0.3 + fy*0.3)
			rise := (1.0 - fy) * 1.5

			intensity := rise * math.Exp(-dist*dist*3.0)

			if intensity > 0.8 {
				b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#aaaaaa")).Render(smokeChars[5]))
			} else if intensity > 0.5 {
				b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#777777")).Render(smokeChars[3]))
			} else if intensity > 0.3 {
				b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#555555")).Render(smokeChars[1]))
			} else if intensity > 0.15 {
				b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#333333")).Render(smokeChars[0]))
			} else {
				b.WriteString(" ")
			}
		}
		lines[row] = b.String()
	}
	return strings.Join(lines, "\n")
}

// ─────────────────────────────────────────────────────────────────────────────
// Ripple — expanding concentric rings
// ─────────────────────────────────────────────────────────────────────────────

func RenderRipple(width, height, frame int) string {
	lines := make([]string, height)
	t := float64(frame) * 0.1
	cx, cy := float64(width)/2.0, float64(height)/2.0

	for row := 0; row < height; row++ {
		var b strings.Builder
		for x := 0; x < width; x++ {
			dx := (float64(x) - cx) / cx
			dy := (float64(row) - cy) / cy * 2.0 // aspect ratio correction
			dist := math.Sqrt(dx*dx + dy*dy)

			// Multiple ripple sources
			v1 := math.Sin(dist*10.0 - t*2.0)
			v2 := math.Sin(math.Sqrt(math.Pow(dx-0.3, 2)+math.Pow(dy-0.3, 2))*8.0 - t*1.5)
			v := (v1 + v2) / 2.0

			if math.Abs(v) > 0.7 {
				hue := math.Mod(dist*2.0+t*0.1, 1.0)
				r, g, bl := hslToRGB(hue, 0.7, 0.5)
				color := lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", r, g, bl))
				b.WriteString(lipgloss.NewStyle().Foreground(color).Render("●"))
			} else if math.Abs(v) > 0.4 {
				b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#333366")).Render("·"))
			} else {
				b.WriteString(" ")
			}
		}
		lines[row] = b.String()
	}
	return strings.Join(lines, "\n")
}

// ─────────────────────────────────────────────────────────────────────────────
// Lightning — branching bolts
// ─────────────────────────────────────────────────────────────────────────────

func RenderLightning(width, height, frame int) string {
	lines := make([]string, height)
	t := frame / 6 // slower update for dramatic effect

	// Generate bolt path
	boltX := width / 2
	boltActive := (t % 12) < 4 // flash every 12 frames

	path := make([]int, height)
	if boltActive {
		path[0] = boltX + (t*7)%10 - 5
		for row := 1; row < height; row++ {
			jitter := ((t*13 + row*7) % 5) - 2
			path[row] = path[row-1] + jitter
			if path[row] < 0 { path[row] = 0 }
			if path[row] >= width { path[row] = width - 1 }
		}
	}

	for row := 0; row < height; row++ {
		var b strings.Builder
		for x := 0; x < width; x++ {
			if boltActive {
				dist := math.Abs(float64(x - path[row]))
				if dist < 1 {
					b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Bold(true).Render("█"))
				} else if dist < 2 {
					b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#aaaaff")).Render("▓"))
				} else if dist < 4 {
					b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#4444aa")).Render("░"))
				} else {
					b.WriteString(" ")
				}
			} else {
				// Dark sky between flashes
				if ((x*31+row*17+t*3)%200) < 2 {
					b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#444444")).Render("·"))
				} else {
					b.WriteString(" ")
				}
			}
		}
		lines[row] = b.String()
	}
	return strings.Join(lines, "\n")
}

// ─────────────────────────────────────────────────────────────────────────────
// Spiral — rotating multi-arm spiral
// ─────────────────────────────────────────────────────────────────────────────

func RenderSpiral(width, height, frame int) string {
	lines := make([]string, height)
	t := float64(frame) * 0.05
	cx, cy := float64(width)/2.0, float64(height)/2.0

	for row := 0; row < height; row++ {
		var b strings.Builder
		for x := 0; x < width; x++ {
			dx := (float64(x) - cx) / cx
			dy := (float64(row) - cy) / cy * 2.0
			dist := math.Sqrt(dx*dx + dy*dy)
			angle := math.Atan2(dy, dx)

			// Spiral arms
			spiral := math.Sin(angle*3.0 - dist*6.0 + t*2.0)

			if dist < 0.05 {
				b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Render("◉"))
			} else if spiral > 0.5 && dist < 1.0 {
				hue := math.Mod(dist+t*0.1, 1.0)
				r, g, bl := hslToRGB(hue, 0.8, 0.5)
				color := lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", r, g, bl))
				b.WriteString(lipgloss.NewStyle().Foreground(color).Render("◆"))
			} else if spiral > 0.2 && dist < 1.0 {
				b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#333355")).Render("·"))
			} else {
				b.WriteString(" ")
			}
		}
		lines[row] = b.String()
	}
	return strings.Join(lines, "\n")
}

// ─────────────────────────────────────────────────────────────────────────────
// Petals — cherry blossom drifting
// ─────────────────────────────────────────────────────────────────────────────

func RenderPetals(width, height, frame int) string {
	lines := make([]string, height)
	t := float64(frame) * 0.06
	petalChars := []string{"❀", "✿", "❁", "✾", "⚘"}
	petalColors := []string{"#ffb7c5", "#ff69b4", "#ffc0cb", "#ff1493", "#ffb6c1"}

	for row := 0; row < height; row++ {
		var b strings.Builder
		for x := 0; x < width; x++ {
			hash := float64((x*89 + row*151) % 1013)
			// Wind drift + gravity
			wind := math.Sin(t*0.3+hash*0.02) * 3.0
			fall := math.Mod(t*0.5+hash*0.03, float64(height+6))

			if hash/1013.0 < 0.04 && math.Abs(fall-float64(row)) < 1.0 {
				idx := int(hash) % len(petalChars)
				cIdx := int(hash*0.7) % len(petalColors)
				_ = wind
				b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(petalColors[cIdx])).Render(petalChars[idx]))
			} else {
				b.WriteString(" ")
			}
		}
		lines[row] = b.String()
	}
	return strings.Join(lines, "\n")
}

// ─────────────────────────────────────────────────────────────────────────────
// Heartbeat — pulsing ECG-style line
// ─────────────────────────────────────────────────────────────────────────────

func RenderHeartbeat(width, height, frame int) string {
	lines := make([]string, height)
	t := float64(frame) * 0.15
	centerY := float64(height) / 2.0

	// Pre-compute the waveform for each column
	waveY := make([]float64, width)
	for x := 0; x < width; x++ {
		// Scrolling ECG waveform
		pos := math.Mod(float64(x)/float64(width)*3.0-t+100.0, 3.0)

		var y float64
		switch {
		case pos < 0.4:
			y = 0 // flat
		case pos < 0.5:
			y = (pos - 0.4) * 3.0 // small P wave up
		case pos < 0.6:
			y = (0.6 - pos) * 3.0 // P wave down
		case pos < 0.8:
			y = 0 // flat
		case pos < 0.9:
			y = -(pos - 0.8) * 8.0 // Q dip
		case pos < 1.1:
			y = -0.8 + (pos-0.9)*14.0 // R spike up
		case pos < 1.3:
			y = 2.0 - (pos-1.1)*15.0 // R spike down to S
		case pos < 1.5:
			y = -1.0 + (pos-1.3)*5.0 // S recovery
		case pos < 1.8:
			y = 0 // flat ST segment
		case pos < 2.1:
			y = math.Sin((pos-1.8)/0.3*math.Pi) * 0.5 // T wave
		default:
			y = 0 // flat
		}
		waveY[x] = y
	}

	// Render
	for row := 0; row < height; row++ {
		var b strings.Builder
		for x := 0; x < width; x++ {
			// Map wave value to row position
			targetRow := centerY - waveY[x]*float64(height)*0.35
			dist := math.Abs(float64(row) - targetRow)

			if dist < 0.6 {
				// Main line — bright red
				b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#ff2222")).Bold(true).Render("━"))
			} else if dist < 1.3 {
				// Glow
				b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#881111")).Render("─"))
			} else if row == int(centerY) {
				// Baseline grid
				b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#1a1a1a")).Render("·"))
			} else {
				b.WriteString(" ")
			}
		}
		lines[row] = b.String()
	}
	return strings.Join(lines, "\n")
}

// ─────────────────────────────────────────────────────────────────────────────
// Pendulum Wave — phase-shifted pendulums
// ─────────────────────────────────────────────────────────────────────────────

func RenderPendulum(width, height, frame int) string {
	lines := make([]string, height)
	t := float64(frame) * 0.06
	numPendulums := width / 3
	if numPendulums < 5 { numPendulums = 5 }

	// Clear all lines
	for row := 0; row < height; row++ {
		var b strings.Builder
		for x := 0; x < width; x++ {
			b.WriteString(" ")
		}
		lines[row] = b.String()
	}

	// Draw each pendulum
	grid := make([][]rune, height)
	for i := range grid {
		grid[i] = make([]rune, width)
		for j := range grid[i] {
			grid[i][j] = ' '
		}
	}

	for p := 0; p < numPendulums; p++ {
		freq := 1.0 + float64(p)*0.05
		px := int(float64(p) / float64(numPendulums) * float64(width))
		angle := math.Sin(t * freq)
		py := int((angle + 1.0) / 2.0 * float64(height-1))

		if px >= 0 && px < width && py >= 0 && py < height {
			grid[py][px] = '●'
		}
	}

	// Render grid
	for row := 0; row < height; row++ {
		var b strings.Builder
		for x := 0; x < width; x++ {
			if grid[row][x] == '●' {
				hue := float64(x) / float64(width)
				r, g, bl := hslToRGB(hue, 0.8, 0.6)
				color := lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", r, g, bl))
				b.WriteString(lipgloss.NewStyle().Foreground(color).Render("●"))
			} else {
				b.WriteString(" ")
			}
		}
		lines[row] = b.String()
	}
	return strings.Join(lines, "\n")
}

// ─────────────────────────────────────────────────────────────────────────────
// Helper: HSL to RGB conversion
// ─────────────────────────────────────────────────────────────────────────────

func hslToRGB(h, s, l float64) (int, int, int) {
	if s == 0 {
		v := int(l * 255)
		return v, v, v
	}
	var q float64
	if l < 0.5 {
		q = l * (1 + s)
	} else {
		q = l + s - l*s
	}
	p := 2*l - q
	r := hueToRGB(p, q, h+1.0/3.0)
	g := hueToRGB(p, q, h)
	b := hueToRGB(p, q, h-1.0/3.0)
	return int(r * 255), int(g * 255), int(b * 255)
}

func hueToRGB(p, q, t float64) float64 {
	if t < 0 { t += 1 }
	if t > 1 { t -= 1 }
	if t < 1.0/6.0 { return p + (q-p)*6*t }
	if t < 1.0/2.0 { return q }
	if t < 2.0/3.0 { return p + (q-p)*(2.0/3.0-t)*6 }
	return p
}

// ─────────────────────────────────────────────────────────────────────────────
// Pixel Pet — animated companion that hangs out in the home/search views
// ─────────────────────────────────────────────────────────────────────────────

// RenderPixelPet renders an animated pixel shiba that walks and vibes to music.
func RenderPixelPet(width, frame int) string {
	// Shiba walk frames — properly designed multi-line pixel dog
	walk1 := []string{
		`    ╱▔▔╲▂▂`,
		`   ▕ ◕ ◕ ▏`,
		`   ▕  ▽  ▏╱`,
		`   ▕▔▔▔▔▔╱`,
		`    ╱╱  ╲╲`,
	}
	walk2 := []string{
		`    ╱▔▔╲▂▂`,
		`   ▕ ◕ ◕ ▏`,
		`   ▕  ▽  ▏╱`,
		`   ▕▔▔▔▔▔╱`,
		`   ╱╱    ╲╲`,
	}
	walk3 := []string{
		`    ╱▔▔╲▂▂`,
		`   ▕ ◕ ◕ ▏`,
		`   ▕  ▽  ▏~`,
		`   ▕▔▔▔▔▔╱`,
		`    ╲╲  ╱╱`,
	}

	// Idle/vibing frames — ears perk, tail wags, notes float
	idle1 := []string{
		`    ╱▔▔╲▂▂  ♪`,
		`   ▕ ◕‿◕ ▏`,
		`   ▕  ω  ▏~`,
		`   ▕▔▔▔▔▔╱`,
		`    ╱▏  ▕╲`,
	}
	idle2 := []string{
		`    ╱▔▔╲▂▂ ♫`,
		`   ▕ ◕‿◕ ▏`,
		`   ▕  ω  ▏~~`,
		`   ▕▔▔▔▔▔╱`,
		`    ╱▏  ▕╲`,
	}
	idle3 := []string{
		`    ╱▔▔╲▂▂♪`,
		`   ▕ ◕‿◕ ▏`,
		`   ▕  ω  ▏~`,
		`   ▕▔▔▔▔▔╱`,
		`    ╱▏  ▕╲`,
	}

	cycle := (frame / 90) % 2

	var petLines []string
	if cycle == 0 {
		frames := [][]string{walk1, walk2, walk3, walk2}
		idx := (frame / 7) % len(frames)
		petLines = frames[idx]
	} else {
		frames := [][]string{idle1, idle2, idle3, idle2}
		idx := (frame / 14) % len(frames)
		petLines = frames[idx]
	}

	// Position
	var pos int
	if cycle == 0 {
		maxPos := width - 18
		if maxPos < 1 {
			maxPos = 1
		}
		pos = (frame / 3) % maxPos
	} else {
		pos = width/2 - 8
	}
	if pos < 0 {
		pos = 0
	}

	padding := strings.Repeat(" ", pos)
	dogColor := lipgloss.NewStyle().Foreground(lipgloss.Color("#f0a030"))
	faceColor := lipgloss.NewStyle().Foreground(ColorHighlight)

	var result strings.Builder
	for i, line := range petLines {
		switch i {
		case 1, 2:
			result.WriteString(padding + faceColor.Render(line) + "\n")
		default:
			result.WriteString(padding + dogColor.Render(line) + "\n")
		}
	}

	return strings.TrimRight(result.String(), "\n")
}
