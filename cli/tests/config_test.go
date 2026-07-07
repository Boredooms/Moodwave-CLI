package tests_test

import (
	"os"
	"testing"

	"github.com/moodwave/moodwave/internal/config"
)

func TestDefaults(t *testing.T) {
	cfg := config.Defaults()

	if cfg.Scanner.MaxDepth <= 0 {
		t.Errorf("expected MaxDepth > 0, got %d", cfg.Scanner.MaxDepth)
	}
	if cfg.Mood.Sensitivity <= 0 {
		t.Errorf("expected Sensitivity > 0, got %.2f", cfg.Mood.Sensitivity)
	}
	if cfg.Visual.Theme == "" {
		t.Error("expected non-empty default theme")
	}
	if cfg.Visual.Mode == "" {
		t.Error("expected non-empty default visual mode")
	}
	if cfg.Paths.ConfigDir == "" {
		t.Error("expected non-empty config dir")
	}
	if cfg.Paths.CacheDir == "" {
		t.Error("expected non-empty cache dir")
	}
	if len(cfg.Sources.Priority) == 0 {
		t.Error("expected non-empty source priority")
	}
}

func TestEnvOverrides(t *testing.T) {
	// Set env vars.
	os.Setenv("MOODWAVE_NO_COLOR", "1")
	os.Setenv("MOODWAVE_THEME", "dark")
	os.Setenv("MOODWAVE_VOLUME", "50")
	defer func() {
		os.Unsetenv("MOODWAVE_NO_COLOR")
		os.Unsetenv("MOODWAVE_THEME")
		os.Unsetenv("MOODWAVE_VOLUME")
	}()

	cfg, err := config.Load("")
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if !cfg.Visual.NoColor {
		t.Error("expected NoColor=true from env")
	}
	if cfg.Visual.Theme != "dark" {
		t.Errorf("expected theme 'dark', got %q", cfg.Visual.Theme)
	}
	if cfg.Playback.Volume != 50 {
		t.Errorf("expected volume 50, got %d", cfg.Playback.Volume)
	}
}

func TestLoad_NoConfigFile(t *testing.T) {
	// Load with no config file — should return defaults without error.
	cfg, err := config.Load("/nonexistent/path")
	if err != nil {
		t.Fatalf("Load with no file should not fail: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
}

func TestEnsureDirectories(t *testing.T) {
	dir := t.TempDir()
	cfg := config.Defaults()
	cfg.Paths.ConfigDir = dir + "/config"
	cfg.Paths.CacheDir = dir + "/cache"

	if err := config.EnsureDirectories(&cfg); err != nil {
		t.Fatalf("EnsureDirectories failed: %v", err)
	}

	if _, err := os.Stat(cfg.Paths.ConfigDir); err != nil {
		t.Errorf("config dir not created: %v", err)
	}
	if _, err := os.Stat(cfg.Paths.CacheDir); err != nil {
		t.Errorf("cache dir not created: %v", err)
	}
}

func TestVisualModes(t *testing.T) {
	modes := []config.VisualMode{
		config.VisualWave,
		config.VisualSpectrum,
		config.VisualPulse,
		config.VisualMinimal,
		config.VisualQuiet,
	}
	for _, m := range modes {
		if m == "" {
			t.Error("VisualMode constant should not be empty")
		}
	}
}
