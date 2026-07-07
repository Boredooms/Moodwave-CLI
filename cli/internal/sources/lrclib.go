// lrclib.go implements the LRCLIB synchronized lyrics adapter.
//
// LRCLIB is a free, open-source, community-driven synchronized lyrics database.
// No API key required.
// Docs: https://lrclib.net/docs
package sources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const lrclibBase = "https://lrclib.net/api"

// LRCLIBAdapter implements LyricsAdapter for LRCLIB.
type LRCLIBAdapter struct {
	client *http.Client
}

// NewLRCLIBAdapter creates a LRCLIB lyrics adapter.
func NewLRCLIBAdapter() *LRCLIBAdapter {
	return &LRCLIBAdapter{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// Name returns the adapter identifier.
func (a *LRCLIBAdapter) Name() string {
	return "lrclib"
}

// HealthCheck pings the LRCLIB API.
func (a *LRCLIBAdapter) HealthCheck(ctx context.Context) error {
	endpoint := lrclibBase + "/search?q=love&limit=1"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", radioBrowserUserAgent)
	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("lrclib: health check failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("lrclib: server returned %d", resp.StatusCode)
	}
	return nil
}

// SearchTracks is not the primary use case for LRCLIB — returns empty.
func (a *LRCLIBAdapter) SearchTracks(_ context.Context, _ SearchQuery) ([]Track, error) {
	return nil, nil
}

// SearchStations is not applicable for LRCLIB.
func (a *LRCLIBAdapter) SearchStations(_ context.Context, _ SearchQuery) ([]Station, error) {
	return nil, nil
}

// ResolveTrack is not applicable for LRCLIB.
func (a *LRCLIBAdapter) ResolveTrack(_ context.Context, _ string) (*Track, error) {
	return nil, fmt.Errorf("lrclib: track resolution not supported; use FetchLyrics")
}

// FetchLyrics retrieves synchronized lyrics for a track by artist and title.
// Implements LyricsAdapter.
func (a *LRCLIBAdapter) FetchLyrics(ctx context.Context, artist, title string) ([]LyricLine, error) {
	// First try the get endpoint (exact match).
	endpoint := fmt.Sprintf(
		"%s/search?q=%s %s&limit=5",
		lrclibBase,
		url.QueryEscape(artist),
		url.QueryEscape(title),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", radioBrowserUserAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("lrclib: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, nil // no lyrics available
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("lrclib: server returned %d", resp.StatusCode)
	}

	var results []lrclibResult
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("lrclib: parse error: %w", err)
	}

	// Find the best match.
	var best *lrclibResult
	for i := range results {
		r := &results[i]
		if r.SyncedLyrics == "" {
			continue
		}
		// Fuzzy match: check if artist/title are contained.
		titleMatch := strings.Contains(strings.ToLower(r.TrackName), strings.ToLower(title))
		artistMatch := strings.Contains(strings.ToLower(r.ArtistName), strings.ToLower(artist))
		if titleMatch && artistMatch {
			best = r
			break
		}
		if best == nil && (titleMatch || artistMatch) {
			best = r
		}
	}

	if best == nil || best.SyncedLyrics == "" {
		return nil, nil
	}

	return parseLRC(best.SyncedLyrics), nil
}

// parseLRC parses LRC format lyrics into a slice of LyricLine.
// LRC format: [mm:ss.xx] lyric text
func parseLRC(lrc string) []LyricLine {
	var lines []LyricLine
	for _, line := range strings.Split(lrc, "\n") {
		line = strings.TrimSpace(line)
		if len(line) < 10 {
			continue
		}
		// Format: [01:23.45] some lyric text
		if line[0] != '[' {
			continue
		}
		end := strings.Index(line, "]")
		if end < 0 {
			continue
		}
		timeStr := line[1:end]
		text := strings.TrimSpace(line[end+1:])
		seconds := parseLRCTime(timeStr)
		if seconds < 0 {
			continue
		}
		lines = append(lines, LyricLine{Time: seconds, Text: text})
	}
	return lines
}

// parseLRCTime parses "mm:ss.xx" into total seconds.
func parseLRCTime(s string) float64 {
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return -1
	}
	var minutes, secs float64
	fmt.Sscanf(parts[0], "%f", &minutes)
	fmt.Sscanf(parts[1], "%f", &secs)
	return minutes*60 + secs
}

type lrclibResult struct {
	ID           int     `json:"id"`
	TrackName    string  `json:"trackName"`
	ArtistName   string  `json:"artistName"`
	AlbumName    string  `json:"albumName"`
	Duration     float64 `json:"duration"`
	PlainLyrics  string  `json:"plainLyrics"`
	SyncedLyrics string  `json:"syncedLyrics"`
}
