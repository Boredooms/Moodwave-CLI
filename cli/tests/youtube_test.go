package tests_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/moodwave/moodwave/internal/sources"
)

func TestYouTubeSearchDirect(t *testing.T) {
	yt := sources.NewYouTubeAdapter()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := yt.HealthCheck(ctx)
	if err != nil {
		t.Fatalf("YouTube health check failed: %v", err)
	}

	q := sources.SearchQuery{
		Text:  "anirudh ravichander",
		Limit: 3,
	}

	tracks, err := yt.SearchTracks(ctx, q)
	if err != nil {
		t.Fatalf("YouTube SearchTracks failed: %v", err)
	}

	t.Logf("Found %d tracks from YouTube", len(tracks))
	for _, tr := range tracks {
		t.Logf("Track: %s by %s, stream: %s", tr.Title, tr.Artist, tr.StreamURL)
	}

	if len(tracks) == 0 {
		t.Error("Expected to find at least one track from YouTube")
	}
}

func TestYouTubeResolveTrack(t *testing.T) {
	yt := sources.NewYouTubeAdapter()

	// High timeout for dynamic downloading.
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	videoID := "B6ppp6WBCKg" // Hangova Anirudh
	track, err := yt.ResolveTrack(ctx, videoID)
	if err != nil {
		t.Fatalf("YouTube ResolveTrack failed: %v", err)
	}

	t.Logf("Resolved Track: %s, stream URL: %s", track.Title, track.StreamURL)
	if track.StreamURL == "" || !strings.Contains(track.StreamURL, "http") {
		t.Errorf("Expected valid HTTP stream URL, got %q", track.StreamURL)
	}
}
