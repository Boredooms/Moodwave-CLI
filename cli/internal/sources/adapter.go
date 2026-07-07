// Package sources defines the music source adapter interface and
// the shared Track and Station data models used across all adapters.
//
// Each adapter implements the Adapter interface and can be registered
// with the source registry. The CLI selects adapters in priority order.
package sources

import (
	"context"
	"time"
)

// Track is the canonical music track representation used across all adapters.
type Track struct {
	// ID is the adapter-specific identifier.
	ID string `json:"id"`

	// Source identifies which adapter provided this track.
	Source string `json:"source"`

	// Title is the track name.
	Title string `json:"title"`

	// Artist is the artist name.
	Artist string `json:"artist"`

	// Album is the album name (may be empty).
	Album string `json:"album"`

	// Duration in seconds (0 if unknown).
	Duration int `json:"duration"`

	// BPM is the tempo in beats per minute (0 if unknown).
	BPM int `json:"bpm"`

	// Energy is a 0.0–1.0 measure of track energy (0 if unknown).
	Energy float64 `json:"energy"`

	// MoodTags are descriptive mood/genre tags from the source.
	MoodTags []string `json:"mood_tags"`

	// StreamURL is a direct audio stream URL (empty if not streamable).
	StreamURL string `json:"stream_url"`

	// PreviewURL is a short preview URL (30s clips, etc.).
	PreviewURL string `json:"preview_url"`

	// LyricsAvailable is true if synchronized lyrics can be fetched.
	LyricsAvailable bool `json:"lyrics_available"`

	// License describes the content license (e.g. "CC-BY").
	License string `json:"license"`

	// CacheTTL is when this entry expires.
	CacheTTL time.Time `json:"cache_ttl"`

	// Score is the recommendation rank score (set by recommender).
	Score float64 `json:"score,omitempty"`
}

// Station represents a radio station (used for ambient/fallback streams).
type Station struct {
	// ID is the adapter-specific station identifier.
	ID string `json:"id"`

	// Source identifies the adapter.
	Source string `json:"source"`

	// Name is the station name.
	Name string `json:"name"`

	// Country is the station's country code.
	Country string `json:"country"`

	// Language is the broadcast language.
	Language string `json:"language"`

	// Tags are genre/mood tags.
	Tags []string `json:"tags"`

	// StreamURL is the station stream URL.
	StreamURL string `json:"stream_url"`

	// Bitrate in kbps (0 if unknown).
	Bitrate int `json:"bitrate"`

	// Codec describes the audio codec (mp3, aac, ogg, etc.).
	Codec string `json:"codec"`

	// IsReliable is the adapter's assessment of stream reliability.
	IsReliable bool `json:"is_reliable"`

	// License describes the content license (e.g. "CC-BY").
	License string `json:"license,omitempty"`

	// Score is the recommendation rank score.
	Score float64 `json:"score,omitempty"`
}

// LyricLine is a single line of synchronized lyrics.
type LyricLine struct {
	// Time in seconds from the start of the track.
	Time float64 `json:"time"`
	// Text of the lyric line.
	Text string `json:"text"`
}

// SearchQuery encapsulates a search request.
type SearchQuery struct {
	// Text is a free-text search string.
	Text string

	// Tags are mood/genre tags to filter by.
	Tags []string

	// BPMMin and BPMMax constrain tempo.
	BPMMin, BPMMax int

	// EnergyMin and EnergyMax constrain energy level.
	EnergyMin, EnergyMax float64

	// Limit is the maximum number of results.
	Limit int

	// StationsOnly requests station results rather than tracks.
	StationsOnly bool
}

// Adapter is the interface all music source adapters must implement.
type Adapter interface {
	// Name returns the unique adapter identifier (e.g. "radio_browser").
	Name() string

	// HealthCheck verifies the adapter can reach its upstream service.
	// Returns nil if healthy, an error otherwise.
	HealthCheck(ctx context.Context) error

	// SearchTracks finds tracks matching the query.
	// Returns an empty slice (not an error) if nothing is found.
	SearchTracks(ctx context.Context, q SearchQuery) ([]Track, error)

	// SearchStations finds radio stations matching the query.
	// Returns an empty slice (not an error) if nothing is found.
	SearchStations(ctx context.Context, q SearchQuery) ([]Station, error)

	// Resolve fetches full metadata for a track or station by ID.
	// This is called after SearchTracks returns partial results.
	ResolveTrack(ctx context.Context, id string) (*Track, error)
}

// LyricsAdapter is an optional extension for adapters that can provide lyrics.
type LyricsAdapter interface {
	Adapter
	// FetchLyrics retrieves synchronized lyrics for a track.
	FetchLyrics(ctx context.Context, artist, title string) ([]LyricLine, error)
}

// Registry holds the registered adapters in priority order.
type Registry struct {
	adapters []Adapter
}

// NewRegistry creates an empty adapter registry.
func NewRegistry() *Registry {
	return &Registry{}
}

// Register adds an adapter to the registry.
// Adapters are tried in registration order.
func (r *Registry) Register(a Adapter) {
	r.adapters = append(r.adapters, a)
}

// All returns the registered adapters in order.
func (r *Registry) All() []Adapter {
	return r.adapters
}

// Get returns the adapter with the given name, or nil.
func (r *Registry) Get(name string) Adapter {
	for _, a := range r.adapters {
		if a.Name() == name {
			return a
		}
	}
	return nil
}

// HealthyAdapters returns adapters that pass their health check.
func (r *Registry) HealthyAdapters(ctx context.Context) []Adapter {
	var healthy []Adapter
	for _, a := range r.adapters {
		if err := a.HealthCheck(ctx); err == nil {
			healthy = append(healthy, a)
		}
	}
	return healthy
}
