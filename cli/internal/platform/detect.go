// Package platform detects OS capabilities, terminal features, and
// runtime environment. It is the lowest-level package in the stack —
// no other internal package should be imported here.
package platform

import (
	"os"
	"runtime"
	"strconv"
	"strings"

	"golang.org/x/term"
)

// OS constants.
const (
	OSWindows = "windows"
	OSDarwin  = "darwin"
	OSLinux   = "linux"
)

// Capabilities describes what the current terminal can do.
type Capabilities struct {
	// OS is GOOS: "windows", "darwin", "linux"
	OS string

	// Arch is GOARCH: "amd64", "arm64", etc.
	Arch string

	// IsTTY is true when stdin and stdout are a real terminal.
	IsTTY bool

	// HasColor is true when the terminal supports ANSI color codes.
	HasColor bool

	// HasTrueColor is true when 24-bit color is supported.
	HasTrueColor bool

	// HasUnicode is true when the terminal likely supports Unicode.
	HasUnicode bool

	// HasAnimation is true when animation is expected to work smoothly.
	HasAnimation bool

	// ReducedMotion is true when the user has indicated they prefer
	// less motion (NO_MOTION=1 or TERM_REDUCED_MOTION=1).
	ReducedMotion bool

	// Width is the current terminal column count (0 if unknown).
	Width int

	// Height is the current terminal row count (0 if unknown).
	Height int

	// Term is the TERM environment variable value.
	Term string

	// ColorTerm is the COLORTERM environment variable value.
	ColorTerm string
}

// Detect probes the current runtime and returns a Capabilities struct.
// This function is safe to call from any goroutine.
func Detect() Capabilities {
	c := Capabilities{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
		Term: os.Getenv("TERM"),
	}

	// TTY detection — both stdin and stdout must be terminals.
	c.IsTTY = term.IsTerminal(int(os.Stdin.Fd())) &&
		term.IsTerminal(int(os.Stdout.Fd()))

	// Terminal dimensions.
	if c.IsTTY {
		w, h, err := term.GetSize(int(os.Stdout.Fd()))
		if err == nil {
			c.Width = w
			c.Height = h
		}
	}

	// Color detection.
	c.HasColor = detectColor()
	c.HasTrueColor = detectTrueColor()

	// Unicode detection — best-effort heuristic.
	c.HasUnicode = detectUnicode()

	// Animation: on if TTY, has color, and NOT in CI.
	c.HasAnimation = c.IsTTY && c.HasColor && os.Getenv("CI") == ""

	// Reduced motion preference.
	c.ReducedMotion = os.Getenv("NO_MOTION") == "1" ||
		os.Getenv("TERM_REDUCED_MOTION") == "1" ||
		os.Getenv("MOODWAVE_NO_ANIMATION") == "1"

	return c
}

// SafeWidth returns the terminal width, defaulting to 80 if unknown.
func (c Capabilities) SafeWidth() int {
	if c.Width < 20 {
		return 80
	}
	return c.Width
}

// SafeHeight returns the terminal height, defaulting to 24 if unknown.
func (c Capabilities) SafeHeight() int {
	if c.Height < 4 {
		return 24
	}
	return c.Height
}

// IsNarrow returns true for terminals narrower than 60 columns.
func (c Capabilities) IsNarrow() bool {
	return c.SafeWidth() < 60
}

// IsWide returns true for terminals 120+ columns wide.
func (c Capabilities) IsWide() bool {
	return c.SafeWidth() >= 120
}

// detectColor returns true if the terminal appears to support ANSI colors.
func detectColor() bool {
	// Explicit disable.
	if os.Getenv("NO_COLOR") != "" || os.Getenv("MOODWAVE_NO_COLOR") == "1" {
		return false
	}

	// Explicit enable.
	if os.Getenv("FORCE_COLOR") != "" || os.Getenv("MOODWAVE_FORCE_COLOR") == "1" {
		return true
	}

	// Windows: Virtual Terminal Processing available on modern terminals.
	if runtime.GOOS == OSWindows {
		return windowsVTEnabled()
	}

	term := os.Getenv("TERM")
	colorterm := strings.ToLower(os.Getenv("COLORTERM"))

	// Explicit color terminal declarations.
	if colorterm == "truecolor" || colorterm == "24bit" || colorterm == "yes" {
		return true
	}

	// Common color-capable terminal types.
	colorTerms := []string{
		"xterm", "xterm-256color", "xterm-color",
		"screen", "screen-256color",
		"tmux", "tmux-256color",
		"rxvt", "rxvt-unicode",
		"vte", "vte-256color",
		"alacritty", "kitty", "wezterm",
		"ansi", "linux",
	}
	for _, t := range colorTerms {
		if strings.Contains(term, t) {
			return true
		}
	}

	return false
}

// detectTrueColor returns true if 24-bit color is supported.
func detectTrueColor() bool {
	colorterm := strings.ToLower(os.Getenv("COLORTERM"))
	return colorterm == "truecolor" || colorterm == "24bit"
}

// detectUnicode returns true if the terminal likely supports Unicode.
func detectUnicode() bool {
	if runtime.GOOS == OSWindows {
		// Modern Windows Terminal and PowerShell support Unicode.
		// Check WT_SESSION (Windows Terminal session ID).
		return os.Getenv("WT_SESSION") != "" ||
			os.Getenv("TERM_PROGRAM") == "vscode"
	}

	// Check locale settings for UTF-8.
	for _, env := range []string{"LC_ALL", "LC_CTYPE", "LANG"} {
		if strings.Contains(strings.ToUpper(os.Getenv(env)), "UTF") {
			return true
		}
	}

	// Common Unicode-capable terminal emulators.
	term := os.Getenv("TERM_PROGRAM")
	if term == "iTerm.app" || term == "Apple_Terminal" ||
		term == "kitty" || term == "alacritty" ||
		term == "wezterm" || term == "vscode" {
		return true
	}

	return false
}

// CurrentWidth returns the current terminal width or a safe default.
// Safe to call from animation loops.
func CurrentWidth() int {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || w < 20 {
		return 80
	}
	return w
}

// CurrentHeight returns the current terminal height or a safe default.
func CurrentHeight() int {
	_, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || h < 4 {
		return 24
	}
	return h
}

// IsCIEnvironment returns true if running inside a CI/CD system.
func IsCIEnvironment() bool {
	ciVars := []string{"CI", "CONTINUOUS_INTEGRATION", "GITHUB_ACTIONS",
		"GITLAB_CI", "TRAVIS", "CIRCLECI", "BUILDKITE", "DRONE"}
	for _, v := range ciVars {
		if val := os.Getenv(v); val != "" && val != "false" && val != "0" {
			return true
		}
	}
	return false
}

// ColorsEnabled returns true if the MOODWAVE_COLORS env var is explicitly "0" or "false".
func ColorsEnabled() bool {
	v := strings.ToLower(os.Getenv("MOODWAVE_COLORS"))
	if v == "0" || v == "false" || v == "off" {
		return false
	}
	return true
}

// ParseBoolEnv parses a boolean environment variable with a default value.
func ParseBoolEnv(key string, def bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return def
	}
	return b
}
