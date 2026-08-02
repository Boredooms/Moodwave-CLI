// Package tui provides the premium terminal user interface for Moodwave CLI.
// Built with Bubble Tea + Lipgloss for a polished, interactive developer experience.
package tui

import "github.com/charmbracelet/lipgloss"

// Color palette — dark, atmospheric, terminal-native.
var (
	// Base colors
	ColorBg        = lipgloss.Color("#0d0d0d")
	ColorFg        = lipgloss.Color("#c9c9c9")
	ColorDim       = lipgloss.Color("#555555")
	ColorBorder    = lipgloss.Color("#333333")
	ColorHighlight = lipgloss.Color("#00d4aa") // teal/cyan accent
	ColorAccent    = lipgloss.Color("#7c5cff") // purple accent
	ColorWarm      = lipgloss.Color("#ff6b35") // warm orange
	ColorSuccess   = lipgloss.Color("#39e75f") // green
	ColorError     = lipgloss.Color("#ff3b3b") // red
	ColorMuted     = lipgloss.Color("#666666")

	// Mood-specific colors
	ColorMoodFocused      = lipgloss.Color("#00d4aa")
	ColorMoodCalm         = lipgloss.Color("#5b9bd5")
	ColorMoodIntense      = lipgloss.Color("#ff4444")
	ColorMoodChaotic      = lipgloss.Color("#ff00ff")
	ColorMoodExperimental = lipgloss.Color("#ffaa00")
	ColorMoodMinimal      = lipgloss.Color("#888888")
	ColorMoodPolished     = lipgloss.Color("#c9a0dc")
	ColorMoodLateNight    = lipgloss.Color("#4a4aff")
	ColorMoodSprint       = lipgloss.Color("#ff6b35")
	ColorMoodDebugging    = lipgloss.Color("#ffd700")
)

// Styles — these are functions now so they always use the current theme colors.
var (
	// Panel borders
	StyleBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder)

	StyleActiveBorder = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ColorHighlight)

	// Text styles
	StyleTitle = lipgloss.NewStyle().
			Foreground(ColorFg).
			Bold(true)

	StyleSubtitle = lipgloss.NewStyle().
			Foreground(ColorDim)

	StyleHighlight = lipgloss.NewStyle().
			Foreground(ColorHighlight).
			Bold(true)

	StyleAccent = lipgloss.NewStyle().
			Foreground(ColorAccent)

	StyleDim = lipgloss.NewStyle().
			Foreground(ColorDim)

	StyleError = lipgloss.NewStyle().
			Foreground(ColorError)

	StyleSuccess = lipgloss.NewStyle().
			Foreground(ColorSuccess)

	// Status bar
	StyleStatusBar = lipgloss.NewStyle().
			Foreground(ColorDim).
			Background(lipgloss.Color("#1a1a1a")).
			Padding(0, 1)

	StyleStatusKey = lipgloss.NewStyle().
			Foreground(ColorHighlight).
			Bold(true)

	// Menu items
	StyleMenuItem = lipgloss.NewStyle().
			Foreground(ColorFg).
			PaddingLeft(2)

	StyleMenuActive = lipgloss.NewStyle().
			Foreground(ColorHighlight).
			Bold(true).
			PaddingLeft(1)

	// Wave/visualization
	StyleWave = lipgloss.NewStyle().
			Foreground(ColorHighlight)

	StyleWaveDim = lipgloss.NewStyle().
			Foreground(ColorBorder)
)

// MoodColor returns the accent color for a given mood label.
func MoodColor(mood string) lipgloss.Color {
	switch mood {
	case "focused":
		return ColorMoodFocused
	case "calm":
		return ColorMoodCalm
	case "intense":
		return ColorMoodIntense
	case "chaotic":
		return ColorMoodChaotic
	case "experimental":
		return ColorMoodExperimental
	case "minimal":
		return ColorMoodMinimal
	case "polished":
		return ColorMoodPolished
	case "late-night":
		return ColorMoodLateNight
	case "sprint":
		return ColorMoodSprint
	case "debugging":
		return ColorMoodDebugging
	default:
		return ColorHighlight
	}
}
