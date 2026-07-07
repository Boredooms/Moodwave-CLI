// jamendo.go implements the Jamendo Creative Commons music catalog adapter.
//
// Jamendo provides a REST API for discovering CC-licensed music.
// Personal use is free; a client_id is required but free to register.
// Docs: https://developer.jamendo.com/v3.0
//
// This adapter is optional — it will gracefully degrade to unavailable
// if no client_id is configured.
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

const jamendoBase = "https://api.jamendo.com/v3.0"

// JamendoAdapter implements Adapter for the Jamendo music catalog.
// If ClientID is empty, the adapter reports as unavailable via HealthCheck.
type JamendoAdapter struct {
	clientID string
	client   *http.Client
}

// NewJamendoAdapter creates a Jamendo adapter.
// clientID may be empty — the adapter will be disabled but won't panic.
func NewJamendoAdapter(clientID string) *JamendoAdapter {
	return &JamendoAdapter{
		clientID: clientID,
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// Name returns the adapter identifier.
func (a *JamendoAdapter) Name() string {
	return "jamendo"
}

// HealthCheck verifies the Jamendo API is reachable and the client_id works.
func (a *JamendoAdapter) HealthCheck(ctx context.Context) error {
	if a.clientID == "" {
		return fmt.Errorf("jamendo: no client_id configured (set JAMENDO_CLIENT_ID)")
	}

	endpoint := fmt.Sprintf("%s/tracks?client_id=%s&format=json&limit=1", jamendoBase, a.clientID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", radioBrowserUserAgent)

	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("jamendo: health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 || resp.StatusCode == 403 {
		return fmt.Errorf("jamendo: invalid client_id")
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("jamendo: server returned %d", resp.StatusCode)
	}
	return nil
}

// SearchTracks finds Creative Commons tracks on Jamendo matching the query.
func (a *JamendoAdapter) SearchTracks(ctx context.Context, q SearchQuery) ([]Track, error) {
	if a.clientID == "" {
		return nil, nil // disabled
	}

	params := url.Values{}
	params.Set("client_id", a.clientID)
	params.Set("format", "json")
	params.Set("limit", fmt.Sprintf("%d", coalesce(q.Limit, 15)))
	params.Set("audioformat", "mp32")
	params.Set("include", "musicinfo")

	// Mood/genre filter.
	if len(q.Tags) > 0 {
		params.Set("tags", strings.Join(q.Tags[:min(3, len(q.Tags))], " "))
	}
	if q.Text != "" {
		params.Set("search", q.Text)
	}

	// Sort by popularity.
	params.Set("order", "popularity_total")

	endpoint := fmt.Sprintf("%s/tracks?%s", jamendoBase, params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", radioBrowserUserAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("jamendo: search failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("jamendo: server returned %d", resp.StatusCode)
	}

	var result jamendoTracksResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("jamendo: parse error: %w", err)
	}

	tracks := make([]Track, 0, len(result.Results))
	for _, jt := range result.Results {
		t := Track{
			ID:        jt.ID,
			Source:    a.Name(),
			Title:     jt.Name,
			Artist:    jt.ArtistName,
			Album:     jt.AlbumName,
			Duration:  jt.Duration,
			StreamURL: jt.AudioDownload,
			License:   "CC",
		}
		// Combine genres, instruments, and vartags.
		for _, tag := range jt.MusicInfo.Tags.Genres {
			t.MoodTags = append(t.MoodTags, tag)
		}
		for _, tag := range jt.MusicInfo.Tags.Instruments {
			t.MoodTags = append(t.MoodTags, tag)
		}
		for _, tag := range jt.MusicInfo.Tags.Vartags {
			t.MoodTags = append(t.MoodTags, tag)
		}
		tracks = append(tracks, t)
	}

	return tracks, nil
}

// SearchStations uses Jamendo radios as fallback ambient streams.
func (a *JamendoAdapter) SearchStations(ctx context.Context, q SearchQuery) ([]Station, error) {
	if a.clientID == "" {
		return nil, nil
	}

	params := url.Values{}
	params.Set("client_id", a.clientID)
	params.Set("format", "json")
	params.Set("limit", "10")

	if len(q.Tags) > 0 {
		params.Set("tags", q.Tags[0])
	}

	endpoint := fmt.Sprintf("%s/radios?%s", jamendoBase, params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", radioBrowserUserAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("jamendo: radio search failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, nil // soft fail for radios
	}

	var result jamendoRadiosResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, nil
	}

	stations := make([]Station, 0, len(result.Results))
	for _, jr := range result.Results {
		stations = append(stations, Station{
			ID:        jr.ID,
			Source:    a.Name(),
			Name:      jr.Name,
			Tags:      splitTags(jr.Dispname),
			StreamURL: jr.Stream,
			License:   "CC",
		})
	}

	return stations, nil
}

// ResolveTrack fetches full metadata for a Jamendo track by ID.
func (a *JamendoAdapter) ResolveTrack(ctx context.Context, id string) (*Track, error) {
	if a.clientID == "" {
		return nil, fmt.Errorf("jamendo: not configured")
	}

	params := url.Values{}
	params.Set("client_id", a.clientID)
	params.Set("format", "json")
	params.Set("id", id)
	params.Set("audioformat", "mp3")

	endpoint := fmt.Sprintf("%s/tracks?%s", jamendoBase, params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", radioBrowserUserAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("jamendo: resolve failed: %w", err)
	}
	defer resp.Body.Close()

	var result jamendoTracksResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("jamendo: parse error: %w", err)
	}
	if len(result.Results) == 0 {
		return nil, fmt.Errorf("jamendo: track %s not found", id)
	}

	jt := result.Results[0]
	t := &Track{
		ID:        jt.ID,
		Source:    a.Name(),
		Title:     jt.Name,
		Artist:    jt.ArtistName,
		Album:     jt.AlbumName,
		Duration:  jt.Duration,
		StreamURL: jt.AudioDownload,
		License:   "CC",
	}
	return t, nil
}

// min returns the smaller of two ints.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ── JSON types ────────────────────────────────────────────────────────────────

type jamendoTracksResponse struct {
	Headers struct {
		Status string `json:"status"`
		Code   int    `json:"code"`
	} `json:"headers"`
	Results []jamendoTrack `json:"results"`
}

type jamendoTrack struct {
	ID            string           `json:"id"`
	Name          string           `json:"name"`
	Duration      int              `json:"duration"`
	ArtistName    string           `json:"artist_name"`
	AlbumName     string           `json:"album_name"`
	AudioDownload string           `json:"audiodownload"`
	MusicInfo     jamendoMusicInfo `json:"musicinfo"`
}

type jamendoMusicInfo struct {
	Tags struct {
		Genres      []string `json:"genres"`
		Instruments []string `json:"instruments"`
		Vartags     []string `json:"vartags"`
	} `json:"tags"`
}

type jamendoRadiosResponse struct {
	Results []jamendoRadio `json:"results"`
}

type jamendoRadio struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Dispname string `json:"dispname"`
	Stream   string `json:"stream"`
}
