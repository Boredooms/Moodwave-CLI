// musicbrainz.go implements the MusicBrainz metadata adapter.
//
// MusicBrainz is the world's largest open music encyclopedia.
// Rate limit: 1 request per second per IP.
// Docs: https://musicbrainz.org/doc/MusicBrainz_API
package sources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	musicBrainzBase = "https://musicbrainz.org/ws/2"
)

// MusicBrainzAdapter implements Adapter for MusicBrainz metadata.
type MusicBrainzAdapter struct {
	userAgent string
	client    *http.Client
	rateMu    sync.Mutex
	lastReq   time.Time
}

// NewMusicBrainzAdapter creates a MusicBrainz adapter.
// userAgent must be set to a meaningful value per MusicBrainz policy:
// "AppName/version (contact-url-or-email)"
func NewMusicBrainzAdapter(userAgent string) *MusicBrainzAdapter {
	if userAgent == "" {
		userAgent = "MoodwaveCLI/1.0 (https://github.com/moodwave/moodwave)"
	}
	return &MusicBrainzAdapter{
		userAgent: userAgent,
		client: &http.Client{
			Timeout: 20 * time.Second,
		},
	}
}

// Name returns the adapter identifier.
func (a *MusicBrainzAdapter) Name() string {
	return "musicbrainz"
}

// HealthCheck pings the MusicBrainz API.
func (a *MusicBrainzAdapter) HealthCheck(ctx context.Context) error {
	endpoint := musicBrainzBase + "/recording?query=piano&limit=1&fmt=json"
	_, err := a.get(ctx, endpoint)
	return err
}

// SearchTracks searches for recordings on MusicBrainz.
func (a *MusicBrainzAdapter) SearchTracks(ctx context.Context, q SearchQuery) ([]Track, error) {
	if q.Text == "" && len(q.Tags) == 0 {
		return nil, nil
	}

	// Build lucene query for MusicBrainz.
	var queryParts []string

	if q.Text != "" {
		// Search by track title/artist text.
		queryParts = append(queryParts, url.QueryEscape(q.Text))
	}

	if len(q.Tags) > 0 {
		// MusicBrainz tag search.
		tagQuery := "tag:(" + strings.Join(q.Tags, " OR ") + ")"
		queryParts = append(queryParts, tagQuery)
	}

	luceneQuery := strings.Join(queryParts, " AND ")
	limit := coalesce(q.Limit, 10)

	endpoint := fmt.Sprintf(
		"%s/recording?query=%s&limit=%d&fmt=json",
		musicBrainzBase,
		url.QueryEscape(luceneQuery),
		limit,
	)

	body, err := a.get(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	var result mbRecordingSearch
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("musicbrainz: parse error: %w", err)
	}

	tracks := make([]Track, 0, len(result.Recordings))
	for _, rec := range result.Recordings {
		t := Track{
			ID:     rec.ID,
			Source: a.Name(),
			Title:  rec.Title,
		}
		if len(rec.ArtistCredit) > 0 {
			t.Artist = rec.ArtistCredit[0].Name
		}
		if rec.Duration > 0 {
			t.Duration = rec.Duration / 1000 // ms → seconds
		}
		for _, rel := range rec.Releases {
			if t.Album == "" {
				t.Album = rel.Title
			}
		}
		// MusicBrainz does not provide streams — this is metadata only.
		tracks = append(tracks, t)
	}

	return tracks, nil
}

// SearchStations is not applicable for MusicBrainz (metadata-only).
func (a *MusicBrainzAdapter) SearchStations(_ context.Context, _ SearchQuery) ([]Station, error) {
	return nil, nil
}

// ResolveTrack fetches full metadata for a recording by MBID.
func (a *MusicBrainzAdapter) ResolveTrack(ctx context.Context, id string) (*Track, error) {
	endpoint := fmt.Sprintf(
		"%s/recording/%s?inc=artist-credits+releases+tags&fmt=json",
		musicBrainzBase, id,
	)

	body, err := a.get(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	var rec mbRecording
	if err := json.Unmarshal(body, &rec); err != nil {
		return nil, fmt.Errorf("musicbrainz: parse error: %w", err)
	}

	t := &Track{
		ID:     rec.ID,
		Source: a.Name(),
		Title:  rec.Title,
	}
	if len(rec.ArtistCredit) > 0 {
		t.Artist = rec.ArtistCredit[0].Name
	}
	if rec.Duration > 0 {
		t.Duration = rec.Duration / 1000
	}
	for _, tag := range rec.Tags {
		t.MoodTags = append(t.MoodTags, tag.Name)
	}

	return t, nil
}

// get performs a rate-limited HTTP GET request with proper headers.
// MusicBrainz requires: User-Agent header and 1 req/sec rate limit.
func (a *MusicBrainzAdapter) get(ctx context.Context, endpoint string) ([]byte, error) {
	// Rate limiting: enforce at least 1 second between requests.
	a.rateMu.Lock()
	since := time.Since(a.lastReq)
	if since < time.Second {
		time.Sleep(time.Second - since)
	}
	a.lastReq = time.Now()
	a.rateMu.Unlock()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", a.userAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("musicbrainz: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 503 {
		return nil, fmt.Errorf("musicbrainz: rate limited (503)")
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("musicbrainz: server returned %d", resp.StatusCode)
	}

	var body []byte
	buf := make([]byte, 32*1024)
	for {
		n, err := resp.Body.Read(buf)
		body = append(body, buf[:n]...)
		if err != nil {
			break
		}
	}
	return body, nil
}

// ── JSON response types ───────────────────────────────────────────────────────

type mbRecordingSearch struct {
	Recordings []mbRecording `json:"recordings"`
}

type mbRecording struct {
	ID           string      `json:"id"`
	Title        string      `json:"title"`
	Duration     int         `json:"length"`
	ArtistCredit []mbArtist  `json:"artist-credit"`
	Releases     []mbRelease `json:"releases"`
	Tags         []mbTag     `json:"tags"`
}

type mbArtist struct {
	Name   string `json:"name"`
	Artist struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"artist"`
}

type mbRelease struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type mbTag struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}
