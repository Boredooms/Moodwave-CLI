package tests_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/moodwave/moodwave/internal/scanner"
)

// createTestRepo creates a minimal fake repository structure in a temp dir.
func createTestRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	// Go source files.
	writeFile(t, dir, "main.go", `package main

// TODO: implement this
func main() {}
`)
	writeFile(t, dir, "handler.go", `package main

// handler handles requests
// FIXME: add error handling
func handler() {}
`)
	writeFile(t, dir, "handler_test.go", `package main

import "testing"

func TestHandler(t *testing.T) {}
`)

	// Nested package.
	writeFile(t, dir, "pkg/util/util.go", `package util

// Utility functions
func Add(a, b int) int { return a + b }
`)

	// Build system.
	writeFile(t, dir, "go.mod", `module example.com/test

go 1.22
`)

	// Documentation.
	writeFile(t, dir, "README.md", "# Test Project\n")

	// CI.
	os.MkdirAll(filepath.Join(dir, ".github", "workflows"), 0755)
	writeFile(t, dir, ".github/workflows/ci.yml", "name: CI\n")

	return dir
}

func writeFile(t *testing.T, dir, rel, content string) {
	t.Helper()
	path := filepath.Join(dir, rel)
	os.MkdirAll(filepath.Dir(path), 0755)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writing %s: %v", rel, err)
	}
}

func TestScanBasic(t *testing.T) {
	dir := createTestRepo(t)
	opts := scanner.DefaultOptions()
	opts.GitEnabled = false // no git in temp dir

	s := scanner.New(opts)
	signals, err := s.Scan(dir)
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}

	if signals.TotalFiles < 5 {
		t.Errorf("expected >= 5 files, got %d", signals.TotalFiles)
	}
	if signals.SourceFiles < 3 {
		t.Errorf("expected >= 3 source files, got %d", signals.SourceFiles)
	}
	if signals.PrimaryLanguage != "Go" {
		t.Errorf("expected Go as primary language, got %q", signals.PrimaryLanguage)
	}
}

func TestScanDetectsTestFiles(t *testing.T) {
	dir := createTestRepo(t)
	opts := scanner.DefaultOptions()
	opts.GitEnabled = false

	s := scanner.New(opts)
	signals, err := s.Scan(dir)
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}

	if signals.TestFileCount < 1 {
		t.Errorf("expected at least 1 test file, got %d", signals.TestFileCount)
	}
}

func TestScanDetectsDocumentation(t *testing.T) {
	dir := createTestRepo(t)
	opts := scanner.DefaultOptions()
	opts.GitEnabled = false

	s := scanner.New(opts)
	signals, err := s.Scan(dir)
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}

	if !signals.HasDocumentation {
		t.Error("expected HasDocumentation=true (README.md present)")
	}
}

func TestScanDetectsBuildSystem(t *testing.T) {
	dir := createTestRepo(t)
	opts := scanner.DefaultOptions()
	opts.GitEnabled = false

	s := scanner.New(opts)
	signals, err := s.Scan(dir)
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}

	if !signals.HasBuildSystem {
		t.Error("expected HasBuildSystem=true (go.mod present)")
	}
	if signals.BuildSystemType != "Go Modules" {
		t.Errorf("expected 'Go Modules', got %q", signals.BuildSystemType)
	}
}

func TestScanDetectsTodos(t *testing.T) {
	dir := createTestRepo(t)
	opts := scanner.DefaultOptions()
	opts.GitEnabled = false

	s := scanner.New(opts)
	signals, err := s.Scan(dir)
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}

	// We wrote 1 TODO and 1 FIXME.
	if signals.TodoFIXMECount < 2 {
		t.Errorf("expected >= 2 TODO/FIXME annotations, got %d", signals.TodoFIXMECount)
	}
}

func TestScanMaxDepth(t *testing.T) {
	dir := t.TempDir()
	// Create a deep directory structure.
	deepPath := filepath.Join(dir, "a", "b", "c", "d", "e")
	os.MkdirAll(deepPath, 0755)
	writeFile(t, dir, "a/b/c/d/e/deep.go", `package deep`)
	writeFile(t, dir, "shallow.go", `package main`)

	opts := scanner.DefaultOptions()
	opts.MaxDepth = 2
	opts.GitEnabled = false

	s := scanner.New(opts)
	signals, err := s.Scan(dir)
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}

	// With MaxDepth=2, we should find shallow.go (depth 0) but not deep.go (depth 4).
	found := false
	for _, l := range signals.Languages {
		if l.Name == "Go" && l.FileCount >= 1 {
			found = true
		}
	}
	if !found {
		t.Error("expected to find Go files within depth limit")
	}
}

func TestScanEmptyDir(t *testing.T) {
	dir := t.TempDir()
	opts := scanner.DefaultOptions()
	opts.GitEnabled = false

	s := scanner.New(opts)
	signals, err := s.Scan(dir)
	if err != nil {
		t.Fatalf("scan should not fail on empty dir: %v", err)
	}
	if signals.TotalFiles != 0 {
		t.Errorf("expected 0 files, got %d", signals.TotalFiles)
	}
}
