// Package config manages all Moodwave configuration with a clean
// priority chain: CLI flags → env vars → per-project file →
// global user file → built-in defaults.
//
// Config files use TOML format and are stored at:
//   - Global:      ~/.config/moodwave/config.toml
//   - Per-project: .moodwave.toml  (in the scanned repo root)
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// Version is injected at build time via ldflags.
var Version = "dev"

// BuildTime is injected at build time via ldflags.
var BuildTime = "unknown"

// VisualMode describes which terminal visual mode is active.
type VisualMode string

const (
	VisualWave     VisualMode = "wave"
	VisualSpectrum VisualMode = "spectrum"
	VisualPulse    VisualMode = "pulse"
	VisualMinimal  VisualMode = "minimal"
	VisualQuiet    VisualMode = "quiet"
	VisualFire     VisualMode = "fireplace"
)

// ThemeID identifies a visual color theme.
type ThemeID string

const (
	ThemeMonochrome ThemeID = "monochrome"
	ThemeDark       ThemeID = "dark"
	ThemeAsh        ThemeID = "ash"
	ThemeGhost      ThemeID = "ghost"
	ThemeOcean      ThemeID = "ocean"
	ThemeNeon       ThemeID = "neon"
	ThemeSunset     ThemeID = "sunset"
	ThemeMatrix     ThemeID = "matrix"
	ThemeLavender   ThemeID = "lavender"
)

// Config is the fully resolved configuration for a Moodwave session.
type Config struct {
	// Core
	Version   string `json:"version"`
	BuildTime string `json:"build_time"`

	// Scanner settings
	Scanner ScannerConfig `json:"scanner"`

	// Mood engine settings
	Mood MoodConfig `json:"mood"`

	// Music sources
	Sources SourcesConfig `json:"sources"`

	// Playback
	Playback PlaybackConfig `json:"playback"`

	// Visual / rendering
	Visual VisualConfig `json:"visual"`

	// Cache
	Cache CacheConfig `json:"cache"`

	// Paths (resolved at load time)
	Paths PathConfig `json:"paths"`

	// Debug / dev
	Debug bool   `json:"debug"`
	Log   string `json:"log"`
}

// ScannerConfig controls repository analysis behavior.
type ScannerConfig struct {
	// MaxDepth limits directory traversal depth (0 = unlimited).
	MaxDepth int `json:"max_depth"`

	// IgnoreDirs lists additional directories to skip (beyond standard ones).
	IgnoreDirs []string `json:"ignore_dirs"`

	// MaxFiles limits how many files are analyzed (0 = unlimited).
	MaxFiles int `json:"max_files"`

	// GitEnabled enables git metadata reading.
	GitEnabled bool `json:"git_enabled"`

	// CommitWindow is how many recent commits to analyze.
	CommitWindow int `json:"commit_window"`
}

// MoodConfig controls the mood inference engine.
type MoodConfig struct {
	// Sensitivity controls how strongly signals affect the mood score.
	// Range: 0.1 (conservative) to 2.0 (aggressive). Default: 1.0.
	Sensitivity float64 `json:"sensitivity"`

	// UpdateIntervalSecs controls how often auto-scan refreshes the mood.
	UpdateIntervalSecs int `json:"update_interval_secs"`
}

// SourcesConfig lists music source adapter priorities and credentials.
type SourcesConfig struct {
	// Priority is the ordered list of source IDs to try.
	// Defaults: ["radio_browser", "jamendo", "musicbrainz"]
	Priority []string `json:"priority"`

	// RadioBrowser settings
	RadioBrowser RadioBrowserConfig `json:"radio_browser"`

	// MusicBrainz settings
	MusicBrainz MusicBrainzConfig `json:"music_brainz"`

	// Jamendo settings (optional — requires client_id)
	Jamendo JamendoConfig `json:"jamendo"`

	// ListenBrainz settings (optional — requires user token)
	ListenBrainz ListenBrainzConfig `json:"listen_brainz"`

	// LRCLIB settings (lyrics, always enabled)
	LRCLIB LRCLIBConfig `json:"lrclib"`

	// YouTube settings (optional)
	YouTube YouTubeConfig `json:"youtube"`
}

// RadioBrowserConfig configures the Radio Browser adapter.
type RadioBrowserConfig struct {
	Enabled    bool   `json:"enabled"`
	APIBaseURL string `json:"api_base_url"` // empty = auto-discover
}

// MusicBrainzConfig configures the MusicBrainz metadata adapter.
type MusicBrainzConfig struct {
	Enabled   bool   `json:"enabled"`
	UserAgent string `json:"user_agent"`
}

// JamendoConfig configures the Jamendo adapter (optional).
type JamendoConfig struct {
	Enabled  bool   `json:"enabled"`
	ClientID string `json:"client_id"` // optional, from env JAMENDO_CLIENT_ID
}

// ListenBrainzConfig configures the ListenBrainz adapter (optional).
type ListenBrainzConfig struct {
	Enabled   bool   `json:"enabled"`
	UserToken string `json:"user_token"` // optional, from env LISTENBRAINZ_TOKEN
	Username  string `json:"username"`   // optional
}

// LRCLIBConfig configures the LRCLIB lyrics adapter.
type LRCLIBConfig struct {
	Enabled bool `json:"enabled"`
}

// YouTubeConfig configures the unofficial YouTube adapter.
type YouTubeConfig struct {
	Enabled bool `json:"enabled"`
}

// PlaybackConfig controls audio playback behavior.
type PlaybackConfig struct {
	// Backend overrides auto-detected backend. Empty = auto.
	// Valid values: "mpv", "ffplay", "afplay", "powershell", "vlc"
	Backend string `json:"backend"`

	// Volume sets default volume (0-100).
	Volume int `json:"volume"`

	// FadeInSecs is the audio fade-in duration in seconds.
	FadeInSecs float64 `json:"fade_in_secs"`

	// FadeOutSecs is the audio fade-out duration when switching tracks.
	FadeOutSecs float64 `json:"fade_out_secs"`
}

// VisualConfig controls terminal rendering and animation.
type VisualConfig struct {
	// Mode is the default visual mode on startup.
	Mode VisualMode `json:"mode"`

	// Theme is the default color theme.
	Theme ThemeID `json:"theme"`

	// FPS is the target frame rate for animation (0 = auto).
	FPS int `json:"fps"`

	// NoAnimation disables all animation.
	NoAnimation bool `json:"no_animation"`

	// NoColor disables all ANSI color.
	NoColor bool `json:"no_color"`

	// NoUnicode forces ASCII-only output.
	NoUnicode bool `json:"no_unicode"`

	// ShowLyrics enables lyrics display when available.
	ShowLyrics bool `json:"show_lyrics"`

	// CompactMode reduces UI footprint for narrow terminals.
	CompactMode bool `json:"compact_mode"`
}

// CacheConfig controls the metadata cache.
type CacheConfig struct {
	// Enabled enables the disk cache.
	Enabled bool `json:"enabled"`

	// MaxEntries is the maximum number of cached responses.
	MaxEntries int `json:"max_entries"`

	// TTLSecs is the default TTL for cached entries in seconds.
	TTLSecs int `json:"ttl_secs"`
}

// PathConfig holds resolved filesystem paths.
type PathConfig struct {
	// ConfigDir is the global config directory.
	ConfigDir string `json:"config_dir"`

	// ConfigFile is the path to the global config file.
	ConfigFile string `json:"config_file"`

	// CacheDir is the data/cache directory.
	CacheDir string `json:"cache_dir"`

	// LogFile is the path to the log file.
	LogFile string `json:"log_file"`
}

// Defaults returns a Config populated with safe built-in defaults.
func Defaults() Config {
	paths := resolvePaths()
	return Config{
		Version:   Version,
		BuildTime: BuildTime,
		Scanner: ScannerConfig{
			MaxDepth:     12,
			MaxFiles:     50000,
			GitEnabled:   true,
			CommitWindow: 50,
			IgnoreDirs:   []string{},
		},
		Mood: MoodConfig{
			Sensitivity:        1.0,
			UpdateIntervalSecs: 300,
		},
		Sources: SourcesConfig{
			Priority: []string{"youtube", "radio_browser", "jamendo", "musicbrainz"},
			RadioBrowser: RadioBrowserConfig{
				Enabled: true,
			},
			MusicBrainz: MusicBrainzConfig{
				Enabled:   true,
				UserAgent: fmt.Sprintf("MoodwaveCLI/%s (https://github.com/moodwave/moodwave)", Version),
			},
			Jamendo: JamendoConfig{
				Enabled:  true,
				ClientID: "d2e96803",
			},
			ListenBrainz: ListenBrainzConfig{
				Enabled: false, // opt-in; requires user token
			},
			LRCLIB: LRCLIBConfig{
				Enabled: true,
			},
			YouTube: YouTubeConfig{
				Enabled: true,
			},
		},
		Playback: PlaybackConfig{
			Backend:     "",
			Volume:      80,
			FadeInSecs:  1.0,
			FadeOutSecs: 1.5,
		},
		Visual: VisualConfig{
			Mode:        VisualWave,
			Theme:       ThemeMonochrome,
			FPS:         24,
			NoAnimation: false,
			NoColor:     false,
			NoUnicode:   false,
			ShowLyrics:  false,
			CompactMode: false,
		},
		Cache: CacheConfig{
			Enabled:    true,
			MaxEntries: 500,
			TTLSecs:    3600,
		},
		Paths: paths,
		Debug: false,
	}
}

// Load returns a Config by merging defaults, global file, per-project file,
// and environment variables. Does not apply CLI flags — callers do that
// by mutating the returned Config.
func Load(projectRoot string) (*Config, error) {
	cfg := Defaults()

	// Load global config file.
	if err := loadFile(cfg.Paths.ConfigFile, &cfg); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("reading global config: %w", err)
	}

	// Load per-project config file.
	if projectRoot != "" {
		projectCfg := filepath.Join(projectRoot, ".moodwave.toml")
		if err := loadFile(projectCfg, &cfg); err != nil && !os.IsNotExist(err) {
			return nil, fmt.Errorf("reading project config: %w", err)
		}
	}

	// Apply environment variable overrides.
	applyEnvOverrides(&cfg)

	return &cfg, nil
}

// EnsureDirectories creates all required directories if they don't exist.
func EnsureDirectories(cfg *Config) error {
	dirs := []string{cfg.Paths.ConfigDir, cfg.Paths.CacheDir}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return fmt.Errorf("creating directory %s: %w", d, err)
		}
	}
	return nil
}

// WriteDefaults writes a default config file to the global config path.
func WriteDefaults(cfg *Config) error {
	if err := os.MkdirAll(cfg.Paths.ConfigDir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(cfg.Paths.ConfigFile, data, 0644)
}

// resolvePaths computes XDG/platform-appropriate paths.
func resolvePaths() PathConfig {
	var configDir string

	switch runtime.GOOS {
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			home, _ := os.UserHomeDir()
			appData = filepath.Join(home, "AppData", "Roaming")
		}
		configDir = filepath.Join(appData, "moodwave")
	default:
		xdgConfig := os.Getenv("XDG_CONFIG_HOME")
		if xdgConfig == "" {
			home, _ := os.UserHomeDir()
			xdgConfig = filepath.Join(home, ".config")
		}
		configDir = filepath.Join(xdgConfig, "moodwave")
	}

	var cacheDir string
	switch runtime.GOOS {
	case "windows":
		localAppData := os.Getenv("LOCALAPPDATA")
		if localAppData == "" {
			home, _ := os.UserHomeDir()
			localAppData = filepath.Join(home, "AppData", "Local")
		}
		cacheDir = filepath.Join(localAppData, "moodwave", "cache")
	case "darwin":
		home, _ := os.UserHomeDir()
		cacheDir = filepath.Join(home, "Library", "Caches", "moodwave")
	default:
		xdgCache := os.Getenv("XDG_CACHE_HOME")
		if xdgCache == "" {
			home, _ := os.UserHomeDir()
			xdgCache = filepath.Join(home, ".cache")
		}
		cacheDir = filepath.Join(xdgCache, "moodwave")
	}

	return PathConfig{
		ConfigDir:  configDir,
		ConfigFile: filepath.Join(configDir, "config.json"),
		CacheDir:   cacheDir,
		LogFile:    filepath.Join(configDir, "moodwave.log"),
	}
}

// loadFile loads a JSON config file and merges it into cfg.
// Fields in the file override cfg fields; absent fields keep their current value.
func loadFile(path string, cfg *Config) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	// Unmarshal into cfg — JSON merge (existing fields preserved if absent in file).
	return json.Unmarshal(data, cfg)
}

// applyEnvOverrides applies environment variable overrides to cfg.
func applyEnvOverrides(cfg *Config) {
	// Source credentials (never require keys — make them optional env overrides).
	if v := os.Getenv("JAMENDO_CLIENT_ID"); v != "" {
		cfg.Sources.Jamendo.ClientID = v
		cfg.Sources.Jamendo.Enabled = true
	}
	if v := os.Getenv("LISTENBRAINZ_TOKEN"); v != "" {
		cfg.Sources.ListenBrainz.UserToken = v
		cfg.Sources.ListenBrainz.Enabled = true
	}
	if v := os.Getenv("LISTENBRAINZ_USERNAME"); v != "" {
		cfg.Sources.ListenBrainz.Username = v
	}

	// Visual overrides.
	if v := os.Getenv("MOODWAVE_THEME"); v != "" {
		cfg.Visual.Theme = ThemeID(v)
	}
	if v := os.Getenv("MOODWAVE_VISUAL"); v != "" {
		cfg.Visual.Mode = VisualMode(v)
	}
	if v := os.Getenv("MOODWAVE_NO_COLOR"); v == "1" || strings.ToLower(v) == "true" {
		cfg.Visual.NoColor = true
	}
	if v := os.Getenv("NO_COLOR"); v != "" {
		cfg.Visual.NoColor = true
	}
	if v := os.Getenv("MOODWAVE_NO_ANIMATION"); v == "1" || strings.ToLower(v) == "true" {
		cfg.Visual.NoAnimation = true
	}
	if v := os.Getenv("MOODWAVE_NO_UNICODE"); v == "1" || strings.ToLower(v) == "true" {
		cfg.Visual.NoUnicode = true
	}

	// Playback overrides.
	if v := os.Getenv("MOODWAVE_BACKEND"); v != "" {
		cfg.Playback.Backend = v
	}
	if v := os.Getenv("MOODWAVE_VOLUME"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 && n <= 100 {
			cfg.Playback.Volume = n
		}
	}

	// Debug mode.
	if v := os.Getenv("MOODWAVE_DEBUG"); v == "1" || strings.ToLower(v) == "true" {
		cfg.Debug = true
	}
}
