// Package mood provides heuristic-based mood inference from codebase signals.
//
// The engine uses a weighted scoring system across multiple mood dimensions.
// Each mood gets a score (0.0–1.0) based on how strongly the codebase signals
// match it. The top-scoring mood wins, with a confidence value derived from
// how far ahead it is from the runner-up.
//
// Architecture note: This is Phase 1 (heuristics only). The interface is
// designed to accept optional model scores in Phase 2 without API breakage.
package mood

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/moodwave/moodwave/internal/scanner"
)

// Label is a typed mood identifier.
type Label string

const (
	MoodFocused      Label = "focused"
	MoodCalm         Label = "calm"
	MoodIntense      Label = "intense"
	MoodChaotic      Label = "chaotic"
	MoodExperimental Label = "experimental"
	MoodMinimal      Label = "minimal"
	MoodPolished     Label = "polished"
	MoodLateNight    Label = "late-night"
	MoodSprint       Label = "sprint"
	MoodDebugging    Label = "debugging"
)

// AllMoods is the complete list of supported mood labels.
var AllMoods = []Label{
	MoodFocused, MoodCalm, MoodIntense, MoodChaotic,
	MoodExperimental, MoodMinimal, MoodPolished,
	MoodLateNight, MoodSprint, MoodDebugging,
}

// TrackTraits describes the preferred music characteristics for a mood.
type TrackTraits struct {
	BPMMin           int      `json:"bpm_min"`
	BPMMax           int      `json:"bpm_max"`
	EnergyMin        float64  `json:"energy_min"` // 0.0–1.0
	EnergyMax        float64  `json:"energy_max"`
	Instrumentalness float64  `json:"instrumentalness"` // 0.0–1.0 preference
	AmbientLevel     float64  `json:"ambient_level"`    // 0.0–1.0
	Concentration    float64  `json:"concentration"`    // 0.0–1.0
	Tags             []string `json:"tags"`
	RadioTags        []string `json:"radio_tags"` // tags for radio station search
}

// Profile is the output of the mood inference engine.
type Profile struct {
	// Label is the detected mood.
	Label Label `json:"label"`

	// Confidence is how certain the engine is (0.0–1.0).
	Confidence float64 `json:"confidence"`

	// Scores contains the score for every mood (for transparency).
	Scores map[Label]float64 `json:"scores"`

	// Signals lists the top contributing signals with their effect.
	Signals []SignalContribution `json:"signals"`

	// Traits describes the preferred music characteristics.
	Traits TrackTraits `json:"traits"`

	// Explanation is a human-readable explanation of why this mood was chosen.
	Explanation string `json:"explanation"`

	// TransitionStrategy describes how music should change.
	TransitionStrategy string `json:"transition_strategy"`
}

// SignalContribution records how a specific signal contributed to the mood.
type SignalContribution struct {
	Signal    string  `json:"signal"`
	Effect    string  `json:"effect"`    // "increased", "decreased", "neutral"
	Magnitude float64 `json:"magnitude"` // 0.0–1.0
}

// Engine is the mood inference engine.
type Engine struct {
	sensitivity float64 // multiplier applied to all scores
}

// NewEngine creates a mood inference engine.
// sensitivity is typically 1.0; lower = more conservative, higher = more reactive.
func NewEngine(sensitivity float64) *Engine {
	if sensitivity <= 0 {
		sensitivity = 1.0
	}
	return &Engine{sensitivity: sensitivity}
}

// Infer produces a MoodProfile from the given codebase signals.
func (e *Engine) Infer(signals *scanner.CodebaseSignals) *Profile {
	scores := make(map[Label]float64, len(AllMoods))
	var contributions []SignalContribution

	// ── Scoring rules ────────────────────────────────────────────────────────
	// Each rule adds or subtracts from one or more mood scores.
	// Rules are designed to be independent and additive.
	// ─────────────────────────────────────────────────────────────────────────

	// 1. Test coverage signals.
	if signals.TestRatio > 0.25 {
		add(scores, MoodPolished, 0.3)
		add(scores, MoodFocused, 0.2)
		contributions = append(contributions, SignalContribution{
			Signal: "high test coverage", Effect: "increased",
			Magnitude: signals.TestRatio,
		})
	} else if signals.TestRatio < 0.05 {
		add(scores, MoodExperimental, 0.2)
		add(scores, MoodSprint, 0.15)
		contributions = append(contributions, SignalContribution{
			Signal: "low test coverage", Effect: "increased experimental",
			Magnitude: 1 - signals.TestRatio,
		})
	}

	// 2. TODO/FIXME density.
	if signals.SourceFiles > 0 {
		todoRatio := float64(signals.TodoFIXMECount) / float64(signals.SourceFiles)
		if todoRatio > 2.0 {
			add(scores, MoodDebugging, 0.35)
			add(scores, MoodChaotic, 0.25)
			contributions = append(contributions, SignalContribution{
				Signal: "high TODO/FIXME density", Effect: "increased",
				Magnitude: math.Min(1.0, todoRatio/5.0),
			})
		} else if todoRatio > 0.5 {
			add(scores, MoodIntense, 0.15)
			add(scores, MoodSprint, 0.2)
		} else if todoRatio < 0.05 {
			add(scores, MoodPolished, 0.25)
			add(scores, MoodMinimal, 0.15)
		}
	}

	// 3. Comment density.
	if signals.CommentDensity > 0.25 {
		add(scores, MoodPolished, 0.2)
		add(scores, MoodFocused, 0.15)
		contributions = append(contributions, SignalContribution{
			Signal: "well-commented codebase", Effect: "increased",
			Magnitude: signals.CommentDensity,
		})
	} else if signals.CommentDensity < 0.05 {
		add(scores, MoodSprint, 0.15)
		add(scores, MoodExperimental, 0.1)
	}

	// 4. Structure entropy (directory chaos).
	if signals.StructureEntropy > 0.7 {
		add(scores, MoodChaotic, 0.35)
		add(scores, MoodExperimental, 0.2)
		contributions = append(contributions, SignalContribution{
			Signal: "irregular directory structure", Effect: "increased",
			Magnitude: signals.StructureEntropy,
		})
	} else if signals.StructureEntropy < 0.3 {
		add(scores, MoodMinimal, 0.25)
		add(scores, MoodPolished, 0.15)
	}

	// 5. Naming consistency.
	if signals.NamingConsistency > 0.8 {
		add(scores, MoodPolished, 0.2)
		add(scores, MoodFocused, 0.1)
	} else if signals.NamingConsistency < 0.4 {
		add(scores, MoodChaotic, 0.2)
	}

	// 6. Documentation presence.
	if signals.HasDocumentation && signals.DocDensity > 0.1 {
		add(scores, MoodPolished, 0.2)
		add(scores, MoodCalm, 0.15)
	}

	// 7. Build system.
	if signals.HasBuildSystem {
		add(scores, MoodPolished, 0.1)
		add(scores, MoodFocused, 0.1)
	}

	// 8. Dependency weight.
	if signals.DependencyWeight > 0.7 {
		add(scores, MoodIntense, 0.15)
		add(scores, MoodChaotic, 0.1)
	} else if signals.DependencyWeight < 0.2 {
		add(scores, MoodMinimal, 0.2)
	}

	// 9. Primary language hints.
	switch signals.PrimaryLanguage {
	case "Go", "Rust", "C":
		add(scores, MoodFocused, 0.15)
		add(scores, MoodMinimal, 0.1)
	case "Python", "Julia", "R":
		add(scores, MoodExperimental, 0.15)
		add(scores, MoodCalm, 0.1)
	case "JavaScript", "TypeScript":
		add(scores, MoodIntense, 0.1)
	case "Shell", "PowerShell":
		add(scores, MoodFocused, 0.1)
		add(scores, MoodDebugging, 0.1)
	}

	// 10. Language diversity.
	if len(signals.Languages) > 5 {
		add(scores, MoodChaotic, 0.2)
		add(scores, MoodExperimental, 0.15)
	} else if len(signals.Languages) == 1 {
		add(scores, MoodMinimal, 0.2)
		add(scores, MoodFocused, 0.15)
	}

	// 11. Git churn signals.
	if signals.Git.RepoDetected {
		if signals.Git.ChurnScore > 0.7 {
			add(scores, MoodSprint, 0.35)
			add(scores, MoodIntense, 0.25)
			contributions = append(contributions, SignalContribution{
				Signal: "high git churn", Effect: "increased",
				Magnitude: signals.Git.ChurnScore,
			})
		} else if signals.Git.ChurnScore > 0.4 {
			add(scores, MoodIntense, 0.2)
			add(scores, MoodFocused, 0.15)
		} else if signals.Git.ChurnScore < 0.1 {
			add(scores, MoodCalm, 0.25)
			add(scores, MoodLateNight, 0.2)
		}

		// Last commit age.
		if signals.Git.LastCommitAge > 72 {
			add(scores, MoodCalm, 0.15)
			add(scores, MoodLateNight, 0.1)
		}

		if signals.Git.HasUncommittedChanges {
			add(scores, MoodDebugging, 0.15)
			add(scores, MoodIntense, 0.1)
		}
	} else {
		// No git repo — could be early project or experimental.
		add(scores, MoodExperimental, 0.15)
	}

	// 12. Project size signals.
	if signals.SourceFiles < 5 {
		add(scores, MoodExperimental, 0.25)
		add(scores, MoodMinimal, 0.2)
	} else if signals.SourceFiles > 500 {
		add(scores, MoodIntense, 0.15)
		add(scores, MoodPolished, 0.1)
	}

	// 13. CI presence.
	if signals.HasCI {
		add(scores, MoodPolished, 0.15)
		add(scores, MoodFocused, 0.1)
	}

	// 14. Semantic Vocabulary (Open-source Naive-Bayes classifier simulation).
	if len(signals.SemanticMoodCounts) > 0 {
		totalSemanticMatches := 0
		for _, count := range signals.SemanticMoodCounts {
			totalSemanticMatches += count
		}

		if totalSemanticMatches > 0 {
			for mLabel, count := range signals.SemanticMoodCounts {
				probability := float64(count) / float64(totalSemanticMatches)
				if probability > 0.0 {
					label := Label(mLabel)
					boost := probability * 0.4
					add(scores, label, boost)
					if boost > 0.1 {
						contributions = append(contributions, SignalContribution{
							Signal:    fmt.Sprintf("semantic keywords for %s", label),
							Effect:    "increased",
							Magnitude: probability,
						})
					}
				}
			}
		}
	}

	// ── Apply sensitivity ─────────────────────────────────────────────────────
	for k, v := range scores {
		scores[k] = math.Min(1.0, v*e.sensitivity)
	}

	// ── Ensure all moods have a baseline score ────────────────────────────────
	for _, m := range AllMoods {
		if _, ok := scores[m]; !ok {
			scores[m] = 0.05 // baseline presence
		}
	}

	// ── Find winner ───────────────────────────────────────────────────────────
	type scored struct {
		mood  Label
		score float64
	}
	var ranked []scored
	for m, s := range scores {
		ranked = append(ranked, scored{m, s})
	}
	sort.Slice(ranked, func(i, j int) bool {
		return ranked[i].score > ranked[j].score
	})

	winner := ranked[0]
	runnerUp := ranked[1]

	// Confidence: how far ahead the winner is from the runner-up.
	confidence := 0.5
	if winner.score > 0 {
		gap := winner.score - runnerUp.score
		confidence = math.Min(1.0, 0.5+gap*2)
	}

	// Build the profile.
	profile := &Profile{
		Label:              winner.mood,
		Confidence:         confidence,
		Scores:             scores,
		Signals:            topContributions(contributions, 5),
		Traits:             traitsFor(winner.mood),
		Explanation:        buildExplanation(winner.mood, signals, contributions),
		TransitionStrategy: transitionStrategyFor(winner.mood),
	}

	return profile
}

// add adds delta to the score for mood m, clamping to [0, 1].
func add(scores map[Label]float64, m Label, delta float64) {
	scores[m] = math.Min(1.0, scores[m]+delta)
}

// topContributions returns the top n contributions by magnitude.
func topContributions(all []SignalContribution, n int) []SignalContribution {
	sort.Slice(all, func(i, j int) bool {
		return all[i].Magnitude > all[j].Magnitude
	})
	if len(all) > n {
		return all[:n]
	}
	return all
}

// traitsFor returns the preferred music traits for a mood.
func traitsFor(m Label) TrackTraits {
	traits := map[Label]TrackTraits{
		MoodFocused: {
			BPMMin: 70, BPMMax: 95, EnergyMin: 0.2, EnergyMax: 0.5,
			Instrumentalness: 0.9, AmbientLevel: 0.7, Concentration: 0.9,
			Tags:      []string{"lo-fi", "ambient", "instrumental", "focus", "study"},
			RadioTags: []string{"ambient", "study", "focus", "lofi", "chillout"},
		},
		MoodCalm: {
			BPMMin: 60, BPMMax: 85, EnergyMin: 0.1, EnergyMax: 0.4,
			Instrumentalness: 0.8, AmbientLevel: 0.8, Concentration: 0.7,
			Tags:      []string{"acoustic", "soft", "calm", "meditation", "peaceful"},
			RadioTags: []string{"calm", "acoustic", "soft", "meditation", "classical"},
		},
		MoodIntense: {
			BPMMin: 120, BPMMax: 155, EnergyMin: 0.7, EnergyMax: 1.0,
			Instrumentalness: 0.5, AmbientLevel: 0.2, Concentration: 0.6,
			Tags:      []string{"electronic", "energetic", "driving", "power", "industrial"},
			RadioTags: []string{"electronic", "techno", "industrial", "metal", "energetic"},
		},
		MoodChaotic: {
			BPMMin: 130, BPMMax: 165, EnergyMin: 0.8, EnergyMax: 1.0,
			Instrumentalness: 0.6, AmbientLevel: 0.1, Concentration: 0.3,
			Tags:      []string{"punk", "noise", "experimental", "glitch", "fast"},
			RadioTags: []string{"punk", "noise", "experimental", "glitch", "hardcore"},
		},
		MoodExperimental: {
			BPMMin: 80, BPMMax: 125, EnergyMin: 0.4, EnergyMax: 0.7,
			Instrumentalness: 0.75, AmbientLevel: 0.5, Concentration: 0.5,
			Tags:      []string{"generative", "avant-garde", "electronic", "weird", "creative"},
			RadioTags: []string{"experimental", "avant-garde", "generative", "electronica", "eclectic"},
		},
		MoodMinimal: {
			BPMMin: 60, BPMMax: 80, EnergyMin: 0.05, EnergyMax: 0.3,
			Instrumentalness: 0.95, AmbientLevel: 0.9, Concentration: 0.8,
			Tags:      []string{"drone", "minimal", "ambient", "silence", "space"},
			RadioTags: []string{"ambient", "drone", "minimal", "space", "silence"},
		},
		MoodPolished: {
			BPMMin: 90, BPMMax: 115, EnergyMin: 0.4, EnergyMax: 0.65,
			Instrumentalness: 0.7, AmbientLevel: 0.5, Concentration: 0.75,
			Tags:      []string{"smooth", "jazz", "clean", "sophisticated", "bossa"},
			RadioTags: []string{"jazz", "smooth", "bossa", "sophisticated", "lounge"},
		},
		MoodLateNight: {
			BPMMin: 60, BPMMax: 85, EnergyMin: 0.1, EnergyMax: 0.4,
			Instrumentalness: 0.85, AmbientLevel: 0.85, Concentration: 0.6,
			Tags:      []string{"dark ambient", "nocturnal", "soft electronic", "night"},
			RadioTags: []string{"ambient", "dark", "nocturnal", "chill", "night"},
		},
		MoodSprint: {
			BPMMin: 120, BPMMax: 145, EnergyMin: 0.65, EnergyMax: 0.9,
			Instrumentalness: 0.6, AmbientLevel: 0.2, Concentration: 0.7,
			Tags:      []string{"workout", "motivational", "driving", "beat", "hype"},
			RadioTags: []string{"dance", "electronic", "workout", "motivation", "hype"},
		},
		MoodDebugging: {
			BPMMin: 70, BPMMax: 95, EnergyMin: 0.25, EnergyMax: 0.5,
			Instrumentalness: 0.85, AmbientLevel: 0.6, Concentration: 0.9,
			Tags:      []string{"repetitive", "stable", "focus", "minimal", "deep"},
			RadioTags: []string{"focus", "study", "deep", "concentration", "lofi"},
		},
	}

	if t, ok := traits[m]; ok {
		return t
	}
	return traits[MoodFocused]
}

// transitionStrategyFor describes how music should evolve for a mood.
func transitionStrategyFor(m Label) string {
	strategies := map[Label]string{
		MoodFocused:      "maintain",
		MoodCalm:         "maintain",
		MoodIntense:      "adjacent",
		MoodChaotic:      "adjacent",
		MoodExperimental: "adjacent",
		MoodMinimal:      "maintain",
		MoodPolished:     "maintain",
		MoodLateNight:    "maintain",
		MoodSprint:       "adjacent",
		MoodDebugging:    "maintain",
	}
	if s, ok := strategies[m]; ok {
		return s
	}
	return "maintain"
}

// buildExplanation creates a human-readable explanation for the detected mood.
func buildExplanation(m Label, signals *scanner.CodebaseSignals, contributions []SignalContribution) string {
	var sb strings.Builder

	descriptions := map[Label]string{
		MoodFocused:      "This codebase shows a focused, disciplined style",
		MoodCalm:         "This codebase has a calm, measured character",
		MoodIntense:      "This codebase shows signs of intense, high-output development",
		MoodChaotic:      "This codebase shows chaotic, high-entropy patterns",
		MoodExperimental: "This codebase has an experimental, exploratory character",
		MoodMinimal:      "This codebase is lean and minimal",
		MoodPolished:     "This codebase shows a polished, mature style",
		MoodLateNight:    "This codebase suggests late-night, quiet development",
		MoodSprint:       "This codebase shows intense sprint-mode activity",
		MoodDebugging:    "This codebase is currently in debugging mode",
	}

	sb.WriteString(descriptions[m])

	if len(contributions) > 0 {
		sb.WriteString(fmt.Sprintf(". Top signals: %s", contributions[0].Signal))
		for _, c := range contributions[1:] {
			sb.WriteString(fmt.Sprintf(", %s", c.Signal))
		}
	}

	if signals.PrimaryLanguage != "" {
		sb.WriteString(fmt.Sprintf(". Primary language: %s", signals.PrimaryLanguage))
	}

	sb.WriteString(".")
	return sb.String()
}

// String returns a display-friendly string for a mood label.
func (l Label) String() string {
	return string(l)
}

// Emoji returns an emoji representative of the mood.
func (l Label) Emoji() string {
	emojis := map[Label]string{
		MoodFocused:      "🎯",
		MoodCalm:         "🌊",
		MoodIntense:      "⚡",
		MoodChaotic:      "🌪️",
		MoodExperimental: "🔬",
		MoodMinimal:      "◽",
		MoodPolished:     "✨",
		MoodLateNight:    "🌙",
		MoodSprint:       "🚀",
		MoodDebugging:    "🐛",
	}
	if e, ok := emojis[l]; ok {
		return e
	}
	return "🎵"
}
