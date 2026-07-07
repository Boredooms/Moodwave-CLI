// radiobrowser.go implements the Radio Browser source adapter.
//
// Radio Browser is a fully public, free, open-source community database
// of internet radio stations. No API key or authentication is required.
// API docs: https://www.radio-browser.info/
package sources

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	radioBrowserFallbackHost = "all.api.radio-browser.info"
	radioBrowserUserAgent    = "MoodwaveCLI/1.0 (https://github.com/moodwave/moodwave)"
)

// RadioBrowserAdapter implements the Adapter interface for Radio Browser.
type RadioBrowserAdapter struct {
	host   string
	client *http.Client
}

// NewRadioBrowserAdapter creates a Radio Browser adapter.
// host may be empty, in which case the fallback host is used.
// The Radio Browser API uses DNS round-robin for load balancing —
// see https://api.radio-browser.info for current server list.
func NewRadioBrowserAdapter(host string) *RadioBrowserAdapter {
	if host == "" {
		host = radioBrowserFallbackHost
	}
	return &RadioBrowserAdapter{
		host: host,
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// Name returns the adapter identifier.
func (a *RadioBrowserAdapter) Name() string {
	return "radio_browser"
}

// HealthCheck pings the Radio Browser API.
func (a *RadioBrowserAdapter) HealthCheck(ctx context.Context) error {
	endpoint := fmt.Sprintf("https://%s/json/stats", a.host)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", radioBrowserUserAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("radio_browser: health check failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("radio_browser: server returned %d", resp.StatusCode)
	}
	return nil
}

// SearchTracks is not applicable for radio stations — returns empty.
func (a *RadioBrowserAdapter) SearchTracks(_ context.Context, _ SearchQuery) ([]Track, error) {
	return nil, nil
}

// SearchStations finds stations matching the mood query.
func (a *RadioBrowserAdapter) SearchStations(ctx context.Context, q SearchQuery) ([]Station, error) {
	if len(q.Tags) == 0 && q.Text == "" {
		return nil, nil
	}

	// Build the API URL.
	// We use the /json/stations/search endpoint with tag filtering.
	params := url.Values{}
	params.Set("limit", fmt.Sprintf("%d", coalesce(q.Limit, 20)))
	params.Set("order", "votes")
	params.Set("reverse", "true")
	params.Set("hidebroken", "true")
	params.Set("has_geo_info", "false")

	if q.Text != "" {
		params.Set("name", q.Text)
	}
	if len(q.Tags) > 0 {
		// Radio Browser accepts comma-separated tags.
		params.Set("tagList", strings.Join(q.Tags, ","))
	}

	endpoint := fmt.Sprintf("https://%s/json/stations/search?%s", a.host, params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", radioBrowserUserAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("radio_browser: search failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("radio_browser: server returned %d", resp.StatusCode)
	}

	var raw []rbStation
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("radio_browser: decoding response: %w", err)
	}

	stations := make([]Station, 0, len(raw))
	for _, r := range raw {
		// Skip stations without a stream URL or with very low votes.
		if r.URLResolved == "" && r.URL == "" {
			continue
		}
		streamURL := r.URLResolved
		if streamURL == "" {
			streamURL = r.URL
		}
		stations = append(stations, Station{
			ID:         r.StationUUID,
			Source:     a.Name(),
			Name:       r.Name,
			Country:    r.Country,
			Language:   r.Language,
			Tags:       splitTags(r.Tags),
			StreamURL:  streamURL,
			Bitrate:    r.Bitrate,
			Codec:      strings.ToLower(r.Codec),
			IsReliable: r.Votes > 50,
		})
	}

	return stations, nil
}

// ResolveTrack is not applicable for this adapter.
func (a *RadioBrowserAdapter) ResolveTrack(_ context.Context, _ string) (*Track, error) {
	return nil, fmt.Errorf("radio_browser: track resolution not supported")
}

// DiscoverServer picks a working server from the DNS round-robin pool.
// Call this once at startup to find the best host for the current session.
func (a *RadioBrowserAdapter) DiscoverServer(ctx context.Context) {
	// The fallback host all.api.radio-browser.info resolves to multiple IPs.
	// If a custom host was set, keep it.
	if a.host != radioBrowserFallbackHost {
		return
	}

	servers := []string{
		"de1.api.radio-browser.info",
		"de2.api.radio-browser.info",
		"nl1.api.radio-browser.info",
	}

	// Pick a random server to distribute load.
	a.host = servers[rand.Intn(len(servers))]
}

// rbStation is the raw JSON response from Radio Browser.
type rbStation struct {
	StationUUID string `json:"stationuuid"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	URLResolved string `json:"url_resolved"`
	Country     string `json:"country"`
	Language    string `json:"language"`
	Tags        string `json:"tags"`
	Codec       string `json:"codec"`
	Bitrate     int    `json:"bitrate"`
	Votes       int    `json:"votes"`
}

// splitTags splits a comma-separated tag string into a slice.
func splitTags(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(strings.ToLower(p))
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

// coalesce returns val if > 0, otherwise def.
func coalesce(val, def int) int {
	if val > 0 {
		return val
	}
	return def
}
