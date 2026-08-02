package tui

import "github.com/charmbracelet/lipgloss"

// Theme defines a complete color scheme for the TUI.
type Theme struct {
	Name      string
	Bg        lipgloss.Color
	Fg        lipgloss.Color
	Dim       lipgloss.Color
	Border    lipgloss.Color
	Highlight lipgloss.Color
	Accent    lipgloss.Color
	Warm      lipgloss.Color
	Success   lipgloss.Color
	Error     lipgloss.Color
	StatusBg  lipgloss.Color
}

// Built-in themes
var (
	ThemeMidnight = Theme{
		Name:      "Midnight",
		Bg:        "#0d0d0d",
		Fg:        "#c9c9c9",
		Dim:       "#555555",
		Border:    "#333333",
		Highlight: "#00d4aa",
		Accent:    "#7c5cff",
		Warm:      "#ff6b35",
		Success:   "#39e75f",
		Error:     "#ff3b3b",
		StatusBg:  "#111111",
	}

	ThemeDracula = Theme{
		Name:      "Dracula",
		Bg:        "#282a36",
		Fg:        "#f8f8f2",
		Dim:       "#6272a4",
		Border:    "#44475a",
		Highlight: "#bd93f9",
		Accent:    "#ff79c6",
		Warm:      "#ffb86c",
		Success:   "#50fa7b",
		Error:     "#ff5555",
		StatusBg:  "#1e1f29",
	}

	ThemeNord = Theme{
		Name:      "Nord",
		Bg:        "#2e3440",
		Fg:        "#eceff4",
		Dim:       "#4c566a",
		Border:    "#3b4252",
		Highlight: "#88c0d0",
		Accent:    "#81a1c1",
		Warm:      "#d08770",
		Success:   "#a3be8c",
		Error:     "#bf616a",
		StatusBg:  "#242933",
	}

	ThemeTokyoNight = Theme{
		Name:      "Tokyo Night",
		Bg:        "#1a1b26",
		Fg:        "#c0caf5",
		Dim:       "#565f89",
		Border:    "#292e42",
		Highlight: "#7aa2f7",
		Accent:    "#bb9af7",
		Warm:      "#e0af68",
		Success:   "#9ece6a",
		Error:     "#f7768e",
		StatusBg:  "#16161e",
	}

	ThemeGruvbox = Theme{
		Name:      "Gruvbox",
		Bg:        "#1d2021",
		Fg:        "#ebdbb2",
		Dim:       "#665c54",
		Border:    "#3c3836",
		Highlight: "#fabd2f",
		Accent:    "#b8bb26",
		Warm:      "#fe8019",
		Success:   "#b8bb26",
		Error:     "#fb4934",
		StatusBg:  "#1a1a1a",
	}

	ThemeCatppuccin = Theme{
		Name:      "Catppuccin",
		Bg:        "#1e1e2e",
		Fg:        "#cdd6f4",
		Dim:       "#585b70",
		Border:    "#313244",
		Highlight: "#89b4fa",
		Accent:    "#cba6f7",
		Warm:      "#fab387",
		Success:   "#a6e3a1",
		Error:     "#f38ba8",
		StatusBg:  "#181825",
	}

	ThemeSolarized = Theme{
		Name:      "Solarized",
		Bg:        "#002b36",
		Fg:        "#839496",
		Dim:       "#586e75",
		Border:    "#073642",
		Highlight: "#2aa198",
		Accent:    "#6c71c4",
		Warm:      "#cb4b16",
		Success:   "#859900",
		Error:     "#dc322f",
		StatusBg:  "#001f27",
	}

	// AllThemes is the ordered list of available themes
	AllThemes = []Theme{
		ThemeMidnight,
		ThemeDracula,
		ThemeNord,
		ThemeTokyoNight,
		ThemeGruvbox,
		ThemeCatppuccin,
		ThemeSolarized,
	}
)

// ActiveTheme is the currently selected theme
var ActiveTheme = ThemeMidnight

// ApplyTheme sets all global color variables to match the given theme.
func ApplyTheme(t Theme) {
	ActiveTheme = t
	ColorBg = t.Bg
	ColorFg = t.Fg
	ColorDim = t.Dim
	ColorBorder = t.Border
	ColorHighlight = t.Highlight
	ColorAccent = t.Accent
	ColorWarm = t.Warm
	ColorSuccess = t.Success
	ColorError = t.Error
}

// CycleTheme advances to the next theme and applies it.
func CycleTheme() string {
	current := 0
	for i, t := range AllThemes {
		if t.Name == ActiveTheme.Name {
			current = i
			break
		}
	}
	next := (current + 1) % len(AllThemes)
	ApplyTheme(AllThemes[next])
	return AllThemes[next].Name
}
