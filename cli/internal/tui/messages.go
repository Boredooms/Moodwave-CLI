package tui

import (
	"github.com/moodwave/moodwave/internal/playback"
	"github.com/moodwave/moodwave/internal/recommender"
	"github.com/moodwave/moodwave/internal/updater"
)

// ─────────────────────────────────────────────────────────────────────────────
// Async messages for Bubble Tea — all background work communicates via these.
// ─────────────────────────────────────────────────────────────────────────────

// SearchIntent describes *why* a search was fired, carried explicitly on
// the result message instead of being inferred from shared mutable model
// state at the time the message happens to arrive. This replaces a design
// where a single Model.AutoPlayNext bool was set before dispatching a
// search and read-and-cleared when the result landed — which broke the
// moment two searches could be in flight with different intents (a manual
// search typed by the user, and an "endless listening" refill search
// firing because the queue ran out) and the wrong one's result arrived
// while the flag was still set from the other.
type SearchIntent int

const (
	// SearchIntentBrowse: user-initiated search — show results for manual
	// selection, don't touch the live queue at all.
	SearchIntentBrowse SearchIntent = iota
	// SearchIntentAutoPlayAll: replace the entire queue with these results
	// and start playing immediately (autonomous play, mood-based search).
	SearchIntentAutoPlayAll
	// SearchIntentRefillQueue: the queue ran out (endless listening) —
	// replace the queue with these results and continue playing seamlessly.
	SearchIntentRefillQueue
)

// SearchResultMsg is sent when a YouTube/source search completes.
type SearchResultMsg struct {
	Candidates []recommender.Candidate
	Query      string
	Err        error
	Intent     SearchIntent
}

// ResolveResultMsg is sent when a YouTube stream URL is resolved.
type ResolveResultMsg struct {
	StreamURL string
	Title     string
	Artist    string
	Err       error
}

// PlaybackStartedMsg is sent when audio playback begins.
type PlaybackStartedMsg struct {
	Controller *playback.Controller
	Backend    string
}

// PlaybackErrorMsg is sent when playback fails.
type PlaybackErrorMsg struct {
	Err error
}

// PlaybackStatusMsg is sent periodically with the current playback state.
// Controller identifies which controller generation this status is about —
// used to discard stale messages from a previous (already-replaced) controller.
type PlaybackStatusMsg struct {
	Controller *playback.Controller
	Status     playback.Status
}

// PlaybackEndedMsg is sent when a track finishes playing.
// Controller identifies which controller generation ended — used to discard
// stale end signals from a controller that has already been replaced.
type PlaybackEndedMsg struct {
	Controller *playback.Controller
}

// ScanResultMsg is sent when a repo scan completes.
type ScanResultMsg struct {
	MoodLabel  string
	MoodEmoji  string
	Confidence float64
	Err        error
}

// UpdateCheckMsg is sent when a background "is there a new version"
// check completes. It never carries a hard error for "no update found" —
// Err is only set for genuine network/API failures, and those are treated
// as silent no-ops on the home screen (a failed background check should
// never interrupt the user with a visible error banner).
type UpdateCheckMsg struct {
	Available      bool
	CurrentVersion string
	LatestVersion  string
	Release        *updater.Release
	Err            error
}

// UpdateApplyMsg is sent when the in-TUI "install this update now" action
// completes (success or failure).
type UpdateApplyMsg struct {
	NewVersion string
	Err        error
}
