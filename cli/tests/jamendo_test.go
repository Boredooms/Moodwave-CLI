package tests_test

import (
	"context"
	"testing"
	"time"

	"github.com/moodwave/moodwave/internal/sources"
)

func TestJamendoSearchDirect(t *testing.T) {
	clientID := "d2e96803"
	j := sources.NewJamendoAdapter(clientID)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := j.HealthCheck(ctx)
	if err != nil {
		t.Fatalf("Jamendo health check failed: %v", err)
	}

	q := sources.SearchQuery{
		Text:  "ambient",
		Limit: 5,
	}

	tracks, err := j.SearchTracks(ctx, q)
	if err != nil {
		t.Fatalf("Jamendo SearchTracks failed: %v", err)
	}

	t.Logf("Found %d tracks from Jamendo", len(tracks))
	for _, tr := range tracks {
		t.Logf("Track: %s by %s, stream: %s", tr.Title, tr.Artist, tr.StreamURL)
	}

	if len(tracks) == 0 {
		t.Error("Expected to find at least one track from Jamendo")
	}
}
