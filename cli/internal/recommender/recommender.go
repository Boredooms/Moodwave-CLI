// Package recommender maps a MoodProfile to a ranked list of music candidates.
// It queries multiple source adapters in priority order, scores all candidates
// against the mood's preferred track traits, and returns a ranked queue.
package recommender

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/moodwave/moodwave/internal/mood"
	"github.com/moodwave/moodwave/internal/sources"
)

// Candidate is a ranked music candidate (track or station).
type Candidate struct {
	Track   *sources.Track   `json:"track,omitempty"`
	Station *sources.Station `json:"station,omitempty"`
	Score   float64          `json:"score"`
	Reason  string           `json:"reason"`
}

// IsStation returns true if this candidate is a radio station.
func (c *Candidate) IsStation() bool {
	return c.Station != nil
}

// DisplayName returns the display-ready name of this candidate.
func (c *Candidate) DisplayName() string {
	if c.IsStation() {
		return c.Station.Name
	}
	if c.Track != nil {
		if c.Track.Artist != "" {
			return c.Track.Artist + " — " + c.Track.Title
		}
		return c.Track.Title
	}
	return "Unknown"
}

// StreamURL returns the playback URL.
func (c *Candidate) StreamURL() string {
	if c.IsStation() {
		return c.Station.StreamURL
	}
	if c.Track != nil {
		return c.Track.StreamURL
	}
	return ""
}

// Config controls recommender behavior.
type Config struct {
	// MaxCandidates is the max number of candidates to return.
	MaxCandidates int

	// StationWeight controls how much stations are preferred over tracks.
	// 0.0 = prefer tracks, 1.0 = prefer stations.
	StationWeight float64

	// DiversityPenalty penalizes candidates from the same source.
	DiversityPenalty float64

	// RecentlyPlayed holds stream URLs of recently played items.
	// These are penalized to avoid repetition.
	RecentlyPlayed []string

	// Timeout for source queries.
	Timeout time.Duration
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		MaxCandidates:    10,
		StationWeight:    0.7, // prefer stations for ambient listening
		DiversityPenalty: 0.15,
		Timeout:          10 * time.Second,
	}
}

// Recommender queries source adapters and ranks candidates.
type Recommender struct {
	registry *sources.Registry
	cfg      Config
}

// New creates a Recommender.
func New(registry *sources.Registry, cfg Config) *Recommender {
	return &Recommender{registry: registry, cfg: cfg}
}

// Recommend queries all healthy adapters and returns ranked candidates.
func (r *Recommender) Recommend(ctx context.Context, profile *mood.Profile, currentTrack *sources.Track) ([]Candidate, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.Timeout)
	defer cancel()

	traits := profile.Traits
	query := buildQuery(profile)

	var all []Candidate

	var similarTracks []sources.Track
	if currentTrack != nil {
		simQuery := sources.SearchQuery{
			Text:      currentTrack.Artist + " " + currentTrack.Title + " similar",
			BPMMin:    traits.BPMMin,
			BPMMax:    traits.BPMMax,
			EnergyMin: traits.EnergyMin,
			EnergyMax: traits.EnergyMax,
			Limit:     5,
		}
		for _, adapter := range r.registry.All() {
			if err := adapter.HealthCheck(ctx); err == nil {
				tracks, err := adapter.SearchTracks(ctx, simQuery)
				if err == nil {
					similarTracks = append(similarTracks, tracks...)
				}
			}
		}
	}

	// Query each adapter in order.
	for _, adapter := range r.registry.All() {
		if err := adapter.HealthCheck(ctx); err != nil {
			continue // skip unhealthy adapters silently
		}

		// Search stations (primary for ambient listening).
		stations, err := adapter.SearchStations(ctx, query)
		if err == nil {
			for i := range stations {
				score := scoreStation(&stations[i], traits)
				score = r.applyPenalties(score, stations[i].StreamURL, adapter.Name(), all)
				all = append(all, Candidate{
					Station: &stations[i],
					Score:   score * r.cfg.StationWeight,
					Reason:  reasonForStation(&stations[i], traits),
				})
			}
		}

		// Search tracks (for when streams are available).
		tracks, err := adapter.SearchTracks(ctx, query)
		if err == nil {
			for i := range tracks {
				if tracks[i].StreamURL == "" {
					continue // no playable stream
				}
				score := scoreTrack(&tracks[i], traits)
				if currentTrack != nil {
					if strings.EqualFold(tracks[i].Artist, currentTrack.Artist) {
						score += 0.20
					}
					sharedTags := 0
					for _, t1 := range tracks[i].MoodTags {
						for _, t2 := range currentTrack.MoodTags {
							if strings.EqualFold(t1, t2) {
								sharedTags++
							}
						}
					}
					if sharedTags > 0 {
						score += 0.15 * float64(sharedTags)
					}
				}

				score = r.applyPenalties(score, tracks[i].StreamURL, adapter.Name(), all)
				all = append(all, Candidate{
					Track:  &tracks[i],
					Score:  score * (1.0 - r.cfg.StationWeight*0.3),
					Reason: reasonForTrack(&tracks[i], traits),
				})
			}
		}
	}

	for i := range similarTracks {
		if similarTracks[i].StreamURL == "" {
			continue
		}
		score := scoreTrack(&similarTracks[i], traits) + 0.25
		score = r.applyPenalties(score, similarTracks[i].StreamURL, "similar_recommender", all)
		all = append(all, Candidate{
			Track:  &similarTracks[i],
			Score:  score * (1.0 - r.cfg.StationWeight*0.3),
			Reason: fmt.Sprintf("Similar to: %s", currentTrack.Title),
		})
	}

	// Sort by score descending.
	sort.Slice(all, func(i, j int) bool {
		return all[i].Score > all[j].Score
	})

	// Return top N.
	limit := r.cfg.MaxCandidates
	if len(all) < limit {
		limit = len(all)
	}
	return all[:limit], nil
}

// buildQuery converts a mood profile into a source search query.
func buildQuery(profile *mood.Profile) sources.SearchQuery {
	qText := ""
	if len(profile.Traits.Tags) > 0 {
		tagsToUse := profile.Traits.Tags
		if len(tagsToUse) > 2 {
			tagsToUse = tagsToUse[:2]
		}
		qText = strings.Join(tagsToUse, " ") + " music"
	} else {
		qText = string(profile.Label) + " music"
	}

	return sources.SearchQuery{
		Text:      qText,
		Tags:      profile.Traits.RadioTags,
		BPMMin:    profile.Traits.BPMMin,
		BPMMax:    profile.Traits.BPMMax,
		EnergyMin: profile.Traits.EnergyMin,
		EnergyMax: profile.Traits.EnergyMax,
		Limit:     20,
	}
}

// scoreStation rates a station's fit for the given mood traits.
// Returns a score in [0, 1].
func scoreStation(s *sources.Station, traits mood.TrackTraits) float64 {
	if s.StreamURL == "" {
		return 0
	}

	score := 0.3 // baseline for having a stream URL

	// Tag overlap.
	tagOverlap := countTagOverlap(s.Tags, traits.RadioTags)
	if len(traits.RadioTags) > 0 {
		score += float64(tagOverlap) / float64(len(traits.RadioTags)) * 0.5
	}

	// Reliability bonus.
	if s.IsReliable {
		score += 0.15
	}

	// Codec preference (mp3 is most widely supported).
	if s.Codec == "mp3" || s.Codec == "aac" {
		score += 0.05
	}

	return math.Min(1.0, score)
}

// scoreTrack rates a track's fit for the given mood traits.
func scoreTrack(t *sources.Track, traits mood.TrackTraits) float64 {
	if t.StreamURL == "" {
		return 0
	}

	score := 0.4 // baseline

	// BPM fit.
	if t.BPM > 0 {
		if t.BPM >= traits.BPMMin && t.BPM <= traits.BPMMax {
			score += 0.25
		} else {
			gap := 0
			if t.BPM < traits.BPMMin {
				gap = traits.BPMMin - t.BPM
			} else {
				gap = t.BPM - traits.BPMMax
			}
			penalty := math.Min(0.3, float64(gap)/50.0*0.3)
			score -= penalty
		}
	}

	// Energy fit.
	if t.Energy > 0 {
		if t.Energy >= traits.EnergyMin && t.Energy <= traits.EnergyMax {
			score += 0.2
		}
	}

	// Mood tag overlap.
	tagOverlap := countTagOverlap(t.MoodTags, traits.Tags)
	if len(traits.Tags) > 0 {
		score += float64(tagOverlap) / float64(len(traits.Tags)) * 0.15
	}

	return math.Min(1.0, score)
}

// applyPenalties reduces a score for repetition and source diversity.
func (r *Recommender) applyPenalties(score float64, streamURL, source string, existing []Candidate) float64 {
	// Repetition penalty: recently played.
	for _, recent := range r.cfg.RecentlyPlayed {
		if recent == streamURL {
			score *= 0.2
			break
		}
	}

	// Diversity penalty: penalize over-representation of a single source.
	sourceCount := 0
	for _, c := range existing {
		src := ""
		if c.Track != nil {
			src = c.Track.Source
		} else if c.Station != nil {
			src = c.Station.Source
		}
		if src == source {
			sourceCount++
		}
	}
	if sourceCount >= 5 {
		score *= (1.0 - r.cfg.DiversityPenalty)
	}

	return score
}

// countTagOverlap returns the number of tags from needle that appear in haystack.
func countTagOverlap(haystack, needle []string) int {
	set := make(map[string]bool, len(haystack))
	for _, t := range haystack {
		set[strings.ToLower(t)] = true
	}
	count := 0
	for _, t := range needle {
		if set[strings.ToLower(t)] {
			count++
		}
	}
	return count
}

// reasonForStation generates a human-readable reason for a station recommendation.
func reasonForStation(s *sources.Station, traits mood.TrackTraits) string {
	overlap := countTagOverlap(s.Tags, traits.RadioTags)
	if overlap > 0 {
		return "Matches mood tags"
	}
	if s.IsReliable {
		return "Reliable ambient stream"
	}
	return "Available stream"
}

// reasonForTrack generates a human-readable reason for a track recommendation.
func reasonForTrack(t *sources.Track, traits mood.TrackTraits) string {
	if t.BPM >= traits.BPMMin && t.BPM <= traits.BPMMax {
		return "BPM matches mood"
	}
	overlap := countTagOverlap(t.MoodTags, traits.Tags)
	if overlap > 1 {
		return "Strong tag match"
	}
	return "Available track"
}
