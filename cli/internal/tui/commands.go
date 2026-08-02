package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/moodwave/moodwave/internal/playback"
	"github.com/moodwave/moodwave/internal/recommender"
	"github.com/moodwave/moodwave/internal/sources"
	"github.com/moodwave/moodwave/internal/updater"
)

// ─────────────────────────────────────────────────────────────────────────────
// Tea commands — these run in goroutines and return messages to the TUI.
// ─────────────────────────────────────────────────────────────────────────────

// SearchCmd searches YouTube for tracks matching the query. intent is
// carried through unchanged onto the resulting SearchResultMsg so the
// Update loop knows exactly what this particular search was for, without
// needing to consult (or race against) any shared model state.
// Network/timeout errors are retried once before surfacing a clean message
// to the UI, since transient DNS/connection hiccups are common and shouldn't
// dump raw Go error strings (e.g. "context deadline exceeded") on the user.
func SearchCmd(query string, intent SearchIntent) tea.Cmd {
	return func() tea.Msg {
		yt := sources.NewYouTubeAdapter()

		var tracks []sources.Track
		var err error

		for attempt := 0; attempt < 2; attempt++ {
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			tracks, err = yt.SearchTracks(ctx, sources.SearchQuery{
				Text:  query,
				Limit: 10,
			})
			cancel()
			if err == nil {
				break
			}
			if attempt == 0 {
				time.Sleep(500 * time.Millisecond)
			}
		}

		if err != nil {
			return SearchResultMsg{
				Err:    fmt.Errorf("couldn't reach YouTube — check your connection and try again"),
				Query:  query,
				Intent: intent,
			}
		}

		var candidates []recommender.Candidate
		for i := range tracks {
			candidates = append(candidates, recommender.Candidate{
				Track:  &tracks[i],
				Score:  1.0,
				Reason: "YouTube search",
			})
		}

		return SearchResultMsg{
			Candidates: candidates,
			Query:      query,
			Intent:     intent,
		}
	}
}

// ResolveAndPlayCmd resolves a YouTube track and starts playback.
func ResolveAndPlayCmd(candidate *recommender.Candidate, backend string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
		defer cancel()

		streamURL := candidate.StreamURL()
		title := candidate.DisplayName()
		artist := ""

		// Resolve YouTube stream URL if needed
		// Skip re-resolution if we already have a valid googlevideo URL (for loop/repeat)
		if candidate.Track != nil && candidate.Track.Source == "youtube" {
			existingURL := candidate.Track.StreamURL
			if strings.Contains(existingURL, "googlevideo.com") {
				// Already resolved — reuse URL directly (instant replay for loop mode)
				streamURL = existingURL
				title = candidate.Track.Title
				artist = candidate.Track.Artist
			} else {
				// Need to resolve
				yt := sources.NewYouTubeAdapter()
				resolved, err := yt.ResolveTrack(ctx, candidate.Track.ID)
				if err != nil {
					return PlaybackErrorMsg{Err: fmt.Errorf("resolve failed: %w", err)}
				}
				streamURL = resolved.StreamURL
				// Update candidate with resolved metadata
				if resolved.Title != "" && resolved.Title != "YouTube Video" {
					candidate.Track.Title = resolved.Title
					title = resolved.Title
				} else {
					title = candidate.Track.Title
				}
				if resolved.Artist != "" {
					candidate.Track.Artist = resolved.Artist
					artist = resolved.Artist
				} else {
					artist = candidate.Track.Artist
				}
				if resolved.Duration > 0 {
					candidate.Track.Duration = resolved.Duration
				}
				// Store resolved URL for future replays
				candidate.Track.StreamURL = streamURL
			}
		} else if candidate.Track != nil {
			title = candidate.Track.Title
			artist = candidate.Track.Artist
		} else if candidate.Station != nil {
			title = candidate.Station.Name
		}

		if streamURL == "" {
			return PlaybackErrorMsg{Err: fmt.Errorf("no stream URL available")}
		}

		// Ensure ffplay is available (auto-download if needed)
		ensureFFplay()

		// Start playback
		controller, err := playback.NewController(backend)
		if err != nil {
			return PlaybackErrorMsg{Err: fmt.Errorf("no audio backend: %w\nInstall ffmpeg: https://ffmpeg.org/download.html", err)}
		}

		if err := controller.Play(context.Background(), streamURL, title, artist); err != nil {
			return PlaybackErrorMsg{Err: fmt.Errorf("playback failed: %w", err)}
		}

		return PlaybackStartedMsg{
			Controller: controller,
			Backend:    controller.BackendName(),
		}
	}
}

// ensureFFplay checks if ffplay is available — auto-installs if missing.
// This runs inside a tea.Cmd goroutine so it won't block the TUI.
func ensureFFplay() {
	status := playback.CheckFFplay()
	if status.Available {
		return
	}
	// Auto-install ffplay — this is running in a background goroutine
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	_, _ = playback.InstallFFplay(ctx, nil)
}

// PollPlaybackCmd polls the playback controller status every 500ms.
// Every message is tagged with the controller it came from, so the Update
// loop can discard stale messages from a controller that was already
// replaced (e.g. by a fast loop/repeat cycle).
func PollPlaybackCmd(ctrl *playback.Controller) tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
		if ctrl == nil {
			return PlaybackEndedMsg{Controller: nil}
		}
		status := ctrl.Status()
		if status.State == playback.StateStopped || status.State == playback.StateError {
			return PlaybackEndedMsg{Controller: ctrl}
		}
		return PlaybackStatusMsg{Controller: ctrl, Status: status}
	})
}

// ScanRepoCmd runs a repository scan asynchronously using the provided scan function.
func ScanRepoCmd(scanFn ScanFunc) tea.Cmd {
	return func() tea.Msg {
		if scanFn == nil {
			return ScanResultMsg{Err: fmt.Errorf("scan not configured")}
		}
		result := scanFn()
		return ScanResultMsg{
			MoodLabel:  result.MoodLabel,
			MoodEmoji:  result.MoodEmoji,
			Confidence: result.Confidence,
			Err:        result.Err,
		}
	}
}

// buildSearchQuery joins args into a search string.
func BuildSearchQuery(args []string) string {
	return strings.Join(args, " ")
}

// UpdateCheckCmd checks GitHub for a newer release in the background.
// This is what powers the "Update available" banner on the home screen —
// it runs independently of any specific menu action, on every launch,
// end to end: check on start, show the result, let the user apply it
// in-place from inside the TUI without ever leaving to a shell.
func UpdateCheckCmd() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		defer cancel()

		result, err := updater.CheckLatest(ctx)
		if err != nil {
			// Network hiccups / rate limits during a silent background
			// check should never surface as a visible error — just report
			// "no update info available" and move on.
			return UpdateCheckMsg{Err: err}
		}

		return UpdateCheckMsg{
			Available:      result.Available,
			CurrentVersion: result.CurrentVersion,
			LatestVersion:  result.LatestVersion,
			Release:        result.Release,
		}
	}
}

// UpdateApplyCmd downloads and installs the given release in place,
// replacing the running executable. Used by the home screen's "install
// update now" action.
func UpdateApplyCmd(release *updater.Release) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		newVersion, err := updater.Apply(ctx, release, nil)
		if err != nil {
			return UpdateApplyMsg{Err: err}
		}
		return UpdateApplyMsg{NewVersion: newVersion}
	}
}
