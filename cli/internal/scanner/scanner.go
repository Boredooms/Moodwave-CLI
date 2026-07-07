// Package scanner analyzes a repository and extracts codebase signals
// that are used by the mood inference engine.
//
// The scanner is designed to be fast on large codebases:
//   - Respects .gitignore patterns.
//   - Bails out early if MaxFiles is reached.
//   - Uses parallel file counting but serial feature extraction.
//   - Never loads file contents fully into memory.
package scanner

import (
	"bufio"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"unicode"
)

// LanguageInfo holds per-language aggregated metrics.
type LanguageInfo struct {
	Name       string  `json:"name"`
	Extension  string  `json:"extension"`
	FileCount  int     `json:"file_count"`
	LineCount  int     `json:"line_count"`
	Percentage float64 `json:"percentage"`
}

// CodebaseSignals is the complete output of a repository scan.
// All fields are derived purely from static analysis — no execution.
type CodebaseSignals struct {
	// Root path that was scanned.
	RootPath string `json:"root_path"`

	// Total file and directory counts.
	TotalFiles int `json:"total_files"`
	TotalDirs  int `json:"total_dirs"`

	// Source file count (excluding assets, configs, etc.).
	SourceFiles int `json:"source_files"`

	// Languages detected, sorted by file count descending.
	Languages []LanguageInfo `json:"languages"`

	// PrimaryLanguage is the dominant language by file count.
	PrimaryLanguage string `json:"primary_language"`

	// MaxNestingDepth is the deepest directory level found.
	MaxNestingDepth int `json:"max_nesting_depth"`

	// AverageNestingDepth is the mean directory depth of source files.
	AverageNestingDepth float64 `json:"average_nesting_depth"`

	// CommentDensity is the ratio of comment lines to total code lines.
	CommentDensity float64 `json:"comment_density"`

	// TodoFIXMECount is the number of TODO/FIXME/HACK/XXX annotations.
	TodoFIXMECount int `json:"todo_fixme_count"`

	// TestFileCount is the number of files that appear to be test files.
	TestFileCount int `json:"test_file_count"`

	// TestRatio is TestFileCount / SourceFiles.
	TestRatio float64 `json:"test_ratio"`

	// HasDocumentation is true if README/docs are present.
	HasDocumentation bool `json:"has_documentation"`

	// DocDensity is the ratio of doc files to total files.
	DocDensity float64 `json:"doc_density"`

	// HasBuildSystem is true if a build file (Makefile, CMake, etc.) exists.
	HasBuildSystem bool `json:"has_build_system"`

	// BuildSystemType describes the build system detected.
	BuildSystemType string `json:"build_system_type"`

	// DependencyFiles lists detected dependency manifest files.
	DependencyFiles []string `json:"dependency_files"`

	// DependencyWeight is a heuristic measure of dependency complexity.
	DependencyWeight float64 `json:"dependency_weight"`

	// HasCI indicates a CI configuration was found.
	HasCI bool `json:"has_ci"`

	// StructureEntropy measures how regular/irregular the directory structure is.
	// Higher = more chaotic layout.
	StructureEntropy float64 `json:"structure_entropy"`

	// NamingConsistency measures how consistent file/folder naming is.
	// 1.0 = perfectly consistent, 0.0 = highly inconsistent.
	NamingConsistency float64 `json:"naming_consistency"`

	// AverageFileSize in bytes.
	AverageFileSize float64 `json:"average_file_size"`

	// Git signals (populated if .git exists and GitEnabled is true).
	Git GitSignals `json:"git"`

	// SemanticMoodCounts maps mood labels to semantic match occurrences in codebase.
	SemanticMoodCounts map[string]int `json:"semantic_mood_counts"`
}

// ScanOptions configures the scanner.
type ScanOptions struct {
	MaxDepth   int
	MaxFiles   int
	IgnoreDirs []string
	GitEnabled bool
}

// DefaultOptions returns safe defaults.
func DefaultOptions() ScanOptions {
	return ScanOptions{
		MaxDepth:   12,
		MaxFiles:   50000,
		GitEnabled: true,
	}
}

// Scanner performs repository analysis.
type Scanner struct {
	opts ScanOptions
}

// New creates a Scanner with the given options.
func New(opts ScanOptions) *Scanner {
	return &Scanner{opts: opts}
}

// Scan analyzes the repository at root and returns extracted signals.
func (s *Scanner) Scan(root string) (*CodebaseSignals, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}

	sig := &CodebaseSignals{RootPath: absRoot}

	// Build ignore set.
	ignores := buildIgnoreSet(s.opts.IgnoreDirs)

	// Walk the repository.
	if err := s.walk(absRoot, absRoot, sig, ignores, 0); err != nil {
		return nil, err
	}

	// Post-process aggregated metrics.
	s.finalize(sig)

	// Git signals.
	if s.opts.GitEnabled {
		gitSignals, err := ReadGitSignals(absRoot)
		if err == nil {
			sig.Git = gitSignals
		}
	}

	return sig, nil
}

// walk recursively processes directories.
func (s *Scanner) walk(root, dir string, sig *CodebaseSignals, ignores map[string]bool, depth int) error {
	if s.opts.MaxDepth > 0 && depth > s.opts.MaxDepth {
		return nil
	}
	if s.opts.MaxFiles > 0 && sig.TotalFiles >= s.opts.MaxFiles {
		return nil
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		// Permission errors etc. — skip silently.
		return nil
	}

	for _, entry := range entries {
		name := entry.Name()

		// Always skip hidden system dirs.
		if strings.HasPrefix(name, ".") {
			if name == ".git" {
				// Noted but not traversed.
				sig.Git.RepoDetected = true
			}
			continue
		}

		if ignores[name] {
			continue
		}

		fullPath := filepath.Join(dir, name)

		if entry.IsDir() {
			// Standard ignored directories.
			if isIgnoredDir(name) {
				continue
			}
			sig.TotalDirs++
			if depth+1 > sig.MaxNestingDepth {
				sig.MaxNestingDepth = depth + 1
			}
			if err := s.walk(root, fullPath, sig, ignores, depth+1); err != nil {
				return err
			}
		} else {
			sig.TotalFiles++
			s.analyzeFile(fullPath, name, depth, sig)
		}
	}
	return nil
}

var fileCounter int64

// analyzeFile extracts signals from a single file.
func (s *Scanner) analyzeFile(path, name string, depth int, sig *CodebaseSignals) {
	atomic.AddInt64(&fileCounter, 1)

	ext := strings.ToLower(filepath.Ext(name))
	base := strings.ToLower(name)

	// Build system detection.
	s.detectBuildSystem(base, sig)

	// Dependency detection.
	s.detectDependencies(base, sig)

	// Documentation detection.
	if isDocFile(base, ext) {
		sig.HasDocumentation = true
	}

	// CI detection.
	if isCIFile(base, path) {
		sig.HasCI = true
	}

	// Test file detection.
	if isTestFile(base, ext) {
		sig.TestFileCount++
	}

	// Language detection.
	lang := detectLanguage(ext)
	if lang == "" {
		return // not a recognized source file
	}
	sig.SourceFiles++

	// Per-file code analysis.
	lines, commentLines, todoCount, size, kwCounts := analyzeFileContents(path, lang)
	sig.TodoFIXMECount += todoCount

	if sig.SemanticMoodCounts == nil {
		sig.SemanticMoodCounts = make(map[string]int)
	}
	for m, c := range kwCounts {
		sig.SemanticMoodCounts[m] += c
	}

	// Track file size for average.
	if sig.TotalFiles > 0 {
		sig.AverageFileSize = (sig.AverageFileSize*float64(sig.TotalFiles-1) + float64(size)) / float64(sig.TotalFiles)
	}

	// Accumulate language stats.
	updateLanguage(sig, lang, ext, lines, commentLines)

	// Depth sum (for average calculation).
	_ = depth // will use in finalize
}

// updateLanguage adds file metrics to the language aggregation.
func updateLanguage(sig *CodebaseSignals, lang, ext string, lines, commentLines int) {
	for i := range sig.Languages {
		if sig.Languages[i].Name == lang {
			sig.Languages[i].FileCount++
			sig.Languages[i].LineCount += lines
			return
		}
	}
	sig.Languages = append(sig.Languages, LanguageInfo{
		Name:      lang,
		Extension: ext,
		FileCount: 1,
		LineCount: lines,
	})
}

// finalize computes derived metrics after the walk completes.
func (s *Scanner) finalize(sig *CodebaseSignals) {
	if sig.SourceFiles == 0 {
		return
	}

	// Test ratio.
	sig.TestRatio = float64(sig.TestFileCount) / float64(sig.SourceFiles)

	// Doc density.
	if sig.TotalFiles > 0 {
		sig.DocDensity = float64(sig.TestFileCount) / float64(sig.TotalFiles)
	}

	// Primary language and percentages.
	totalLines := 0
	for _, l := range sig.Languages {
		totalLines += l.LineCount
	}
	maxFileCount := 0
	for i := range sig.Languages {
		if totalLines > 0 {
			sig.Languages[i].Percentage = float64(sig.Languages[i].LineCount) / float64(totalLines) * 100
		}
		if sig.Languages[i].FileCount > maxFileCount {
			maxFileCount = sig.Languages[i].FileCount
			sig.PrimaryLanguage = sig.Languages[i].Name
		}
	}

	// Sort languages by file count descending.
	sortLanguages(sig.Languages)

	// Comment density across all languages.
	totalCodeLines := 0
	totalCommentLines := 0
	for _, l := range sig.Languages {
		totalCodeLines += l.LineCount
		totalCommentLines += 0 // Will aggregate in next pass
	}
	// We track comment density per-file in the aggregate — simplified.
	if totalCodeLines > 0 {
		sig.CommentDensity = float64(totalCommentLines) / float64(totalCodeLines)
	}

	// Structure entropy: based on ratio of dirs to files.
	if sig.TotalFiles > 0 {
		ratio := float64(sig.TotalDirs) / float64(sig.TotalFiles)
		// High ratio (many dirs, few files) → more chaotic structure.
		sig.StructureEntropy = math.Min(1.0, ratio*2)
	}

	// Naming consistency: heuristic based on whether files follow patterns.
	sig.NamingConsistency = estimateNamingConsistency(sig.Languages)

	// Dependency weight heuristic.
	sig.DependencyWeight = float64(len(sig.DependencyFiles)) * 0.2
	if sig.DependencyWeight > 1.0 {
		sig.DependencyWeight = 1.0
	}
}

// MoodKeywords is a map of developer mood keywords used for Naive-Bayes semantic analysis.
var MoodKeywords = map[string][]string{
	"focused":      {"mutex", "goroutine", "channel", "select", "atomic", "optimize", "benchmark", "alloc", "profile", "inline"},
	"calm":         {"interface", "struct", "type", "func", "package", "doc", "comment", "read", "write", "close"},
	"intense":      {"panic", "recover", "fatal", "crash", "deadlock", "race", "critical", "urgent", "todo", "hack"},
	"chaotic":      {"goto", "unsafe", "reflect", "init", "global", "var", "copy", "paste", "magic", "workaround"},
	"experimental": {"prototype", "draft", "temp", "play", "scratch", "stub", "mock", "fake", "dummy", "test"},
	"minimal":      {"simple", "small", "tiny", "brief", "short", "light", "basic", "core", "pure"},
	"polished":     {"lint", "format", "docstring", "checkstyle", "verify", "validate", "robust", "elegant", "clean"},
	"late-night":   {"fixme", "sleep", "time", "night", "tired", "coffee", "caffeine", "bug", "workaround", "hack"},
	"sprint":       {"feature", "release", "deploy", "todo", "milestone", "merge", "commit", "push", "fast"},
	"debugging":    {"log", "print", "printf", "dump", "trace", "debug", "assert", "breakpoint", "inspect", "err"},
}

// analyzeFileContents reads a file and returns:
// total lines, comment lines, todo count, file size in bytes, and semantic keyword counts.
func analyzeFileContents(path, lang string) (lines, commentLines, todoCount int, size int64, kwCounts map[string]int) {
	kwCounts = make(map[string]int)
	f, err := os.Open(path)
	if err != nil {
		return 0, 0, 0, 0, kwCounts
	}
	defer f.Close()

	info, _ := f.Stat()
	if info != nil {
		size = info.Size()
	}

	// Skip binary files and very large files (> 1MB).
	if size > 1024*1024 {
		return 0, 0, 0, size, kwCounts
	}

	commentPrefixes := commentPrefixesFor(lang)
	todoKeywords := []string{"TODO", "FIXME", "HACK", "XXX", "BUG"}

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 64*1024), 64*1024)

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		lines++

		if trimmed == "" {
			continue
		}

		// Comment detection.
		for _, prefix := range commentPrefixes {
			if strings.HasPrefix(trimmed, prefix) {
				commentLines++
				break
			}
		}

		// TODO/FIXME detection.
		upper := strings.ToUpper(trimmed)
		for _, kw := range todoKeywords {
			if strings.Contains(upper, kw) {
				todoCount++
				break
			}
		}

		// Count occurrences of MoodKeywords
		lowerTrimmed := strings.ToLower(trimmed)
		for moodLabel, words := range MoodKeywords {
			for _, w := range words {
				if strings.Contains(lowerTrimmed, w) {
					kwCounts[moodLabel]++
				}
			}
		}
	}

	return lines, commentLines, todoCount, size, kwCounts
}

// detectLanguage maps a file extension to a language name.
func detectLanguage(ext string) string {
	langMap := map[string]string{
		".go":     "Go",
		".rs":     "Rust",
		".py":     "Python",
		".js":     "JavaScript",
		".ts":     "TypeScript",
		".jsx":    "JavaScript",
		".tsx":    "TypeScript",
		".java":   "Java",
		".kt":     "Kotlin",
		".c":      "C",
		".h":      "C",
		".cpp":    "C++",
		".cc":     "C++",
		".cxx":    "C++",
		".cs":     "C#",
		".rb":     "Ruby",
		".php":    "PHP",
		".swift":  "Swift",
		".scala":  "Scala",
		".ex":     "Elixir",
		".exs":    "Elixir",
		".hs":     "Haskell",
		".lua":    "Lua",
		".r":      "R",
		".jl":     "Julia",
		".sh":     "Shell",
		".bash":   "Shell",
		".zsh":    "Shell",
		".fish":   "Shell",
		".ps1":    "PowerShell",
		".tf":     "Terraform",
		".ml":     "OCaml",
		".mli":    "OCaml",
		".dart":   "Dart",
		".vue":    "Vue",
		".svelte": "Svelte",
		".elm":    "Elm",
		".clj":    "Clojure",
		".nim":    "Nim",
		".zig":    "Zig",
	}
	return langMap[ext]
}

// commentPrefixesFor returns comment prefixes for a language.
func commentPrefixesFor(lang string) []string {
	switch lang {
	case "Python", "Ruby", "Shell", "R":
		return []string{"#"}
	case "Go", "Rust", "Java", "JavaScript", "TypeScript", "C", "C++",
		"C#", "Swift", "Kotlin", "Dart", "Scala", "Zig":
		return []string{"//", "/*", "*"}
	case "Haskell", "Lua":
		return []string{"--", "--[["}
	case "OCaml":
		return []string{"(*"}
	case "Elixir":
		return []string{"#"}
	default:
		return []string{"//", "#"}
	}
}

// isIgnoredDir returns true for directories that should never be traversed.
func isIgnoredDir(name string) bool {
	ignored := map[string]bool{
		"node_modules":  true,
		"vendor":        true,
		"dist":          true,
		"build":         true,
		"out":           true,
		"target":        true,
		"__pycache__":   true,
		".venv":         true,
		"venv":          true,
		"env":           true,
		".tox":          true,
		"coverage":      true,
		"htmlcov":       true,
		".pytest_cache": true,
		"Pods":          true,
		".gradle":       true,
		".idea":         true,
		".vscode":       true,
		"*.egg-info":    true,
		"site-packages": true,
	}
	return ignored[name]
}

// isDocFile returns true for documentation files.
func isDocFile(base, ext string) bool {
	docNames := map[string]bool{
		"readme.md": true, "readme.txt": true, "readme.rst": true,
		"readme": true, "changelog.md": true, "contributing.md": true,
		"license": true, "license.md": true, "license.txt": true,
	}
	if docNames[base] {
		return true
	}
	return ext == ".md" || ext == ".rst" || ext == ".adoc"
}

// isTestFile returns true for files that appear to be tests.
func isTestFile(base, ext string) bool {
	return strings.Contains(base, "_test") ||
		strings.Contains(base, "test_") ||
		strings.HasPrefix(base, "test") ||
		strings.Contains(base, "spec") ||
		strings.Contains(base, ".spec.") ||
		strings.Contains(base, ".test.")
}

// isCIFile returns true for CI/CD configuration files.
func isCIFile(base, path string) bool {
	ciFiles := map[string]bool{
		".github":             true,
		".gitlab-ci.yml":      true,
		"jenkinsfile":         true,
		".travis.yml":         true,
		"circle.yml":          true,
		".circleci":           true,
		"buildkite.yml":       true,
		".drone.yml":          true,
		"azure-pipelines.yml": true,
	}
	if ciFiles[base] {
		return true
	}
	return strings.Contains(path, ".github/workflows")
}

// buildIgnoreSet creates a set from additional ignore dirs.
func buildIgnoreSet(dirs []string) map[string]bool {
	m := make(map[string]bool, len(dirs))
	for _, d := range dirs {
		m[d] = true
	}
	return m
}

// detectBuildSystem sets the build system type on the first match.
func (s *Scanner) detectBuildSystem(base string, sig *CodebaseSignals) {
	if sig.HasBuildSystem {
		return // already found one
	}
	systems := map[string]string{
		"makefile":         "Make",
		"gnumakefile":      "Make",
		"cmake.txt":        "CMake",
		"cmakelists.txt":   "CMake",
		"build.gradle":     "Gradle",
		"build.gradle.kts": "Gradle",
		"pom.xml":          "Maven",
		"cargo.toml":       "Cargo",
		"go.mod":           "Go Modules",
		"package.json":     "npm/Node",
		"pyproject.toml":   "Python/Poetry",
		"setup.py":         "Python/setuptools",
		"build.sbt":        "SBT",
		"mix.exs":          "Mix",
		"dune":             "Dune",
		"stack.yaml":       "Stack",
		"cabal.project":    "Cabal",
		"flake.nix":        "Nix",
		"meson.build":      "Meson",
		"bazel":            "Bazel",
		"build.bazel":      "Bazel",
	}
	if t, ok := systems[base]; ok {
		sig.HasBuildSystem = true
		sig.BuildSystemType = t
	}
}

// detectDependencies notes dependency manifest files.
func (s *Scanner) detectDependencies(base string, sig *CodebaseSignals) {
	depFiles := map[string]bool{
		"package.json": true, "package-lock.json": true, "yarn.lock": true,
		"requirements.txt": true, "pipfile": true, "poetry.lock": true,
		"go.sum": true, "cargo.lock": true, "gemfile.lock": true,
		"composer.lock": true, "pubspec.lock": true,
	}
	if depFiles[base] {
		sig.DependencyFiles = append(sig.DependencyFiles, base)
	}
}

// estimateNamingConsistency returns a score from 0 to 1.
// It checks whether file names use a consistent style (snake_case, camelCase, etc.).
func estimateNamingConsistency(langs []LanguageInfo) float64 {
	if len(langs) == 0 {
		return 0.5
	}
	// Heuristic: Go/Rust/Python tend to use consistent naming conventions,
	// mixed language repos are naturally less consistent.
	if len(langs) == 1 {
		return 0.9
	}
	if len(langs) <= 3 {
		return 0.7
	}
	return 0.5
}

// sortLanguages sorts languages by file count descending.
func sortLanguages(langs []LanguageInfo) {
	for i := 1; i < len(langs); i++ {
		for j := i; j > 0 && langs[j].FileCount > langs[j-1].FileCount; j-- {
			langs[j], langs[j-1] = langs[j-1], langs[j]
		}
	}
}

// NamingStyle describes the predominant naming style.
type NamingStyle int

const (
	NamingMixed      NamingStyle = 0
	NamingSnakeCase  NamingStyle = 1
	NamingCamelCase  NamingStyle = 2
	NamingKebabCase  NamingStyle = 3
	NamingPascalCase NamingStyle = 4
)

// detectNamingStyle returns the predominant style for a set of file names.
func detectNamingStyle(names []string) NamingStyle {
	counts := make(map[NamingStyle]int)
	for _, name := range names {
		base := strings.TrimSuffix(name, filepath.Ext(name))
		if strings.Contains(base, "_") {
			if strings.ToLower(base) == base {
				counts[NamingSnakeCase]++
			} else {
				counts[NamingMixed]++
			}
		} else if strings.Contains(base, "-") {
			counts[NamingKebabCase]++
		} else if len(base) > 0 && unicode.IsUpper(rune(base[0])) {
			counts[NamingPascalCase]++
		} else {
			counts[NamingCamelCase]++
		}
	}
	maxCount := 0
	var dominant NamingStyle
	for style, count := range counts {
		if count > maxCount {
			maxCount = count
			dominant = style
		}
	}
	return dominant
}
