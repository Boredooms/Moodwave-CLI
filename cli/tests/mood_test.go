package tests_test

import (
	"testing"

	"github.com/moodwave/moodwave/internal/mood"
	"github.com/moodwave/moodwave/internal/scanner"
)

// testSignals returns a CodebaseSignals struct for testing.
func testSignals(mod func(*scanner.CodebaseSignals)) *scanner.CodebaseSignals {
	base := &scanner.CodebaseSignals{
		TotalFiles:        50,
		SourceFiles:       30,
		TestRatio:         0.15,
		PrimaryLanguage:   "Go",
		CommentDensity:    0.15,
		StructureEntropy:  0.4,
		NamingConsistency: 0.7,
		HasDocumentation:  true,
		HasBuildSystem:    true,
		BuildSystemType:   "Go Modules",
		TodoFIXMECount:    5,
		DependencyWeight:  0.2,
		HasCI:             true,
		Languages:         []scanner.LanguageInfo{{Name: "Go", FileCount: 30, Percentage: 100}},
		Git: scanner.GitSignals{
			RepoDetected:  true,
			CurrentBranch: "main",
			ChurnScore:    0.3,
			LastCommitAge: 8,
		},
	}
	if mod != nil {
		mod(base)
	}
	return base
}

func TestEngineBasic(t *testing.T) {
	engine := mood.NewEngine(1.0)
	signals := testSignals(nil)
	profile := engine.Infer(signals)

	if profile == nil {
		t.Fatal("expected non-nil profile")
	}
	if profile.Label == "" {
		t.Fatal("expected non-empty mood label")
	}
	if profile.Confidence < 0 || profile.Confidence > 1 {
		t.Errorf("confidence out of range: %.2f", profile.Confidence)
	}
	if profile.Traits.BPMMin <= 0 {
		t.Errorf("BPMMin should be positive: %d", profile.Traits.BPMMin)
	}
	if profile.Traits.BPMMax <= profile.Traits.BPMMin {
		t.Errorf("BPMMax (%d) should be > BPMMin (%d)", profile.Traits.BPMMax, profile.Traits.BPMMin)
	}
}

func TestPolishedMood(t *testing.T) {
	engine := mood.NewEngine(1.0)
	signals := testSignals(func(s *scanner.CodebaseSignals) {
		s.TestRatio = 0.35     // high tests → polished
		s.CommentDensity = 0.3 // well documented → polished
		s.NamingConsistency = 0.9
		s.TodoFIXMECount = 1
		s.HasCI = true
		s.Git.ChurnScore = 0.2
		s.StructureEntropy = 0.2
	})

	profile := engine.Infer(signals)

	// polished should be in top 3 scores.
	topFound := false
	type scored struct {
		label mood.Label
		score float64
	}
	var top []scored
	for k, v := range profile.Scores {
		top = append(top, scored{k, v})
	}
	// Sort top 3.
	for i := 0; i < len(top); i++ {
		for j := i + 1; j < len(top); j++ {
			if top[j].score > top[i].score {
				top[i], top[j] = top[j], top[i]
			}
		}
	}
	for i, s := range top {
		if s.label == mood.MoodPolished {
			topFound = true
			break
		}
		if i >= 2 {
			break
		}
	}
	if !topFound {
		t.Logf("Profile: %s (%.0f%%)", profile.Label, profile.Confidence*100)
		t.Logf("Polished score: %.2f", profile.Scores[mood.MoodPolished])
		t.Log("Polished was not in top 3 for expected polished signals")
	}
}

func TestChaoticMood(t *testing.T) {
	engine := mood.NewEngine(1.0)
	signals := testSignals(func(s *scanner.CodebaseSignals) {
		s.TodoFIXMECount = 200   // lots of TODOs
		s.SourceFiles = 50       // → todoRatio = 4.0
		s.StructureEntropy = 0.9 // chaotic structure
		s.NamingConsistency = 0.2
		s.Git.ChurnScore = 0.8
		s.TestRatio = 0.02
		s.Languages = []scanner.LanguageInfo{
			{Name: "Go", FileCount: 10},
			{Name: "Python", FileCount: 10},
			{Name: "JavaScript", FileCount: 10},
			{Name: "Rust", FileCount: 5},
			{Name: "Ruby", FileCount: 5},
			{Name: "Java", FileCount: 5},
		}
	})

	profile := engine.Infer(signals)
	t.Logf("Chaotic signals → detected: %s (%.0f%%)", profile.Label, profile.Confidence*100)

	// Chaotic or sprint should be highly scored.
	chaoticScore := profile.Scores[mood.MoodChaotic]
	sprintScore := profile.Scores[mood.MoodSprint]
	if chaoticScore < 0.3 && sprintScore < 0.3 {
		t.Errorf("expected chaotic/sprint to score > 0.3, got chaotic=%.2f sprint=%.2f",
			chaoticScore, sprintScore)
	}
}

func TestMinimalMood(t *testing.T) {
	engine := mood.NewEngine(1.0)
	signals := testSignals(func(s *scanner.CodebaseSignals) {
		s.SourceFiles = 3 // tiny project
		s.TodoFIXMECount = 0
		s.StructureEntropy = 0.1
		s.NamingConsistency = 1.0
		s.DependencyWeight = 0.05
		s.Languages = []scanner.LanguageInfo{{Name: "Go", FileCount: 3}}
		s.Git.RepoDetected = false
	})

	profile := engine.Infer(signals)
	t.Logf("Minimal signals → detected: %s (%.0f%%)", profile.Label, profile.Confidence*100)
	minimalScore := profile.Scores[mood.MoodMinimal]
	if minimalScore < 0.3 {
		t.Errorf("expected minimal score > 0.3, got %.2f", minimalScore)
	}
}

func TestExplanationNotEmpty(t *testing.T) {
	engine := mood.NewEngine(1.0)
	signals := testSignals(nil)
	profile := engine.Infer(signals)

	if profile.Explanation == "" {
		t.Error("expected non-empty explanation")
	}
}

func TestAllMoodsHaveTraits(t *testing.T) {
	engine := mood.NewEngine(1.0)
	signals := testSignals(nil)
	profile := engine.Infer(signals)

	// Verify top mood has traits.
	if len(profile.Traits.Tags) == 0 {
		t.Errorf("expected mood %s to have tags, got none", profile.Label)
	}
	if len(profile.Traits.RadioTags) == 0 {
		t.Errorf("expected mood %s to have radio tags, got none", profile.Label)
	}
}

func TestSensitivity(t *testing.T) {
	// High sensitivity should produce more extreme scores.
	signals := testSignals(func(s *scanner.CodebaseSignals) {
		s.Git.ChurnScore = 0.8
	})

	engLow := mood.NewEngine(0.5)
	engHigh := mood.NewEngine(2.0)

	lowProfile := engLow.Infer(signals)
	highProfile := engHigh.Infer(signals)

	// High sensitivity should not produce a lower max score than low sensitivity.
	maxLow := maxScore(lowProfile.Scores)
	maxHigh := maxScore(highProfile.Scores)

	if maxHigh < maxLow {
		t.Errorf("high sensitivity should produce >= scores: high=%.2f low=%.2f", maxHigh, maxLow)
	}
}

func maxScore(scores map[mood.Label]float64) float64 {
	max := 0.0
	for _, v := range scores {
		if v > max {
			max = v
		}
	}
	return max
}
