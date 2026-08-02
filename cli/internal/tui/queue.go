package tui

import "github.com/moodwave/moodwave/internal/recommender"

// PlaybackQueue is the single owner of "what's playing and what's queued".
//
// Design note: this file exists because the previous design used two
// independent []recommender.Candidate slices (Model.Queue and
// Model.Playlist) with ad-hoc syncing logic scattered across several
// handlers (copy-on-takeover, splice-on-add, mirror-on-remove). Slices in
// Go share backing arrays under append/copy in ways that are easy to get
// subtly wrong, and every one of those sync points was a place a future
// edit could silently desync the two lists — which is exactly what
// happened. There is now exactly one slice, tracks, and exactly one
// currently-playing index, current. The playlist IS the queue, always —
// there is no separate "session playlist" that sometimes gets copied into
// the live queue and sometimes doesn't.
type PlaybackQueue struct {
	tracks  []recommender.Candidate
	current int // index into tracks of the currently playing/selected item; -1 if empty
}

// NewPlaybackQueue creates an empty queue.
func NewPlaybackQueue() PlaybackQueue {
	return PlaybackQueue{current: -1}
}

// Len returns the number of tracks in the queue.
func (q *PlaybackQueue) Len() int {
	return len(q.tracks)
}

// IsEmpty reports whether the queue has no tracks.
func (q *PlaybackQueue) IsEmpty() bool {
	return len(q.tracks) == 0
}

// Current returns the currently playing/selected track, or nil if the
// queue is empty or the index is out of range.
func (q *PlaybackQueue) Current() *recommender.Candidate {
	if q.current < 0 || q.current >= len(q.tracks) {
		return nil
	}
	return &q.tracks[q.current]
}

// CurrentIndex returns the index of the currently playing/selected track.
func (q *PlaybackQueue) CurrentIndex() int {
	return q.current
}

// At returns the track at index i, or nil if out of range.
func (q *PlaybackQueue) At(i int) *recommender.Candidate {
	if i < 0 || i >= len(q.tracks) {
		return nil
	}
	return &q.tracks[i]
}

// All returns the full track list (read-only use expected by callers).
func (q *PlaybackQueue) All() []recommender.Candidate {
	return q.tracks
}

// Replace discards the current queue entirely and starts fresh from the
// given tracks, with playback positioned at startAt (clamped to range).
// Used when starting a brand new search/autonomous-play/playlist session.
func (q *PlaybackQueue) Replace(tracks []recommender.Candidate, startAt int) {
	q.tracks = tracks
	if len(q.tracks) == 0 {
		q.current = -1
		return
	}
	if startAt < 0 {
		startAt = 0
	}
	if startAt >= len(q.tracks) {
		startAt = len(q.tracks) - 1
	}
	q.current = startAt
}

// Append adds a track to the end of the queue. If the queue was empty,
// the new track becomes current immediately.
func (q *PlaybackQueue) Append(track recommender.Candidate) {
	q.tracks = append(q.tracks, track)
	if q.current < 0 {
		q.current = 0
	}
}

// InsertNext splices a track in immediately after the current position,
// so it plays next. If the queue is empty, it becomes the current track.
// This is the ONLY way a track gets "played next" — there is no separate
// playlist splice path, because there is no separate playlist slice.
func (q *PlaybackQueue) InsertNext(track recommender.Candidate) {
	if len(q.tracks) == 0 || q.current < 0 {
		q.tracks = []recommender.Candidate{track}
		q.current = 0
		return
	}
	insertAt := q.current + 1
	if insertAt >= len(q.tracks) {
		q.tracks = append(q.tracks, track)
		return
	}
	q.tracks = append(q.tracks, recommender.Candidate{})
	copy(q.tracks[insertAt+1:], q.tracks[insertAt:])
	q.tracks[insertAt] = track
}

// RemoveAt removes the track at index i. If it was the current track and
// tracks remain, current stays pointing at the same index (which now holds
// what used to be the next track) — matching normal "remove from a list
// you're looking at" expectations. Returns false if i was out of range.
func (q *PlaybackQueue) RemoveAt(i int) bool {
	if i < 0 || i >= len(q.tracks) {
		return false
	}
	q.tracks = append(q.tracks[:i], q.tracks[i+1:]...)
	if len(q.tracks) == 0 {
		q.current = -1
		return true
	}
	if q.current >= len(q.tracks) {
		q.current = len(q.tracks) - 1
	}
	return true
}

// SelectIndex moves the current-playing pointer to i (e.g. the user picked
// a specific track in the Playlist Manager or search results to play now).
// Returns false if i is out of range.
func (q *PlaybackQueue) SelectIndex(i int) bool {
	if i < 0 || i >= len(q.tracks) {
		return false
	}
	q.current = i
	return true
}

// Advance moves current forward by one, without wrapping. Returns false
// if there is no next track (queue exhausted).
func (q *PlaybackQueue) Advance() bool {
	if q.current+1 >= len(q.tracks) {
		return false
	}
	q.current++
	return true
}

// AdvanceWrapping moves current forward by one, wrapping to 0 past the end.
// Returns false only if the queue is empty.
func (q *PlaybackQueue) AdvanceWrapping() bool {
	if len(q.tracks) == 0 {
		return false
	}
	q.current = (q.current + 1) % len(q.tracks)
	return true
}

// ─────────────────────────────────────────────────────────────────────────────
// PersonalPlaylist — a separate, deliberately-curated list that plays out
// BEFORE the auto-generated queue (autonomous play / endless-listening
// refill). It is genuinely a second list, not a variant of PlaybackQueue,
// because it has different semantics: strict FIFO consumption (no
// "current index" concept, no wrap-around, no repeat-mode interaction),
// and it is meant to be built once and drained, not browsed-while-playing
// the way PlaybackQueue is.
//
// Ownership rule (enforced in app.go's resolveOutcome/advance* functions,
// not here): whenever "what plays next" is being decided, PersonalPlaylist
// is checked FIRST. Only once it's empty does normal PlaybackQueue
// advancement / similar-music refill take over. This is the one and only
// place that priority is decided — see advanceIgnoringRepeatOne.
type PersonalPlaylist struct {
	tracks []recommender.Candidate
}

// NewPersonalPlaylist creates an empty personal playlist.
func NewPersonalPlaylist() PersonalPlaylist {
	return PersonalPlaylist{}
}

// Len returns the number of tracks remaining in the personal playlist.
func (p *PersonalPlaylist) Len() int {
	return len(p.tracks)
}

// IsEmpty reports whether the personal playlist has no tracks left.
func (p *PersonalPlaylist) IsEmpty() bool {
	return len(p.tracks) == 0
}

// At returns the track at index i (0 = next to play), or nil if out of range.
func (p *PersonalPlaylist) At(i int) *recommender.Candidate {
	if i < 0 || i >= len(p.tracks) {
		return nil
	}
	return &p.tracks[i]
}

// All returns the full remaining track list (read-only use expected).
func (p *PersonalPlaylist) All() []recommender.Candidate {
	return p.tracks
}

// Add appends a track to the end of the personal playlist.
func (p *PersonalPlaylist) Add(track recommender.Candidate) {
	p.tracks = append(p.tracks, track)
}

// RemoveAt removes the track at index i. Returns false if out of range.
func (p *PersonalPlaylist) RemoveAt(i int) bool {
	if i < 0 || i >= len(p.tracks) {
		return false
	}
	p.tracks = append(p.tracks[:i], p.tracks[i+1:]...)
	return true
}

// PopNext removes and returns the first (next-to-play) track, or nil if
// the personal playlist is empty. This is how playback actually consumes
// it — always from the front, always in the order tracks were added.
func (p *PersonalPlaylist) PopNext() *recommender.Candidate {
	if len(p.tracks) == 0 {
		return nil
	}
	track := p.tracks[0]
	p.tracks = p.tracks[1:]
	return &track
}

// Clear empties the personal playlist entirely.
func (p *PersonalPlaylist) Clear() {
	p.tracks = nil
}
