package tui

import (
	"fmt"
	"math"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/moodwave/moodwave/internal/playback"
	"github.com/moodwave/moodwave/internal/recommender"
	"github.com/moodwave/moodwave/internal/updater"
)

// View represents the current screen.
type View int

const (
	ViewHome View = iota
	ViewScanning
	ViewPlaying
	ViewSearchInput
	ViewSearchResults
	ViewPlaylist
	ViewPersonalPlaylist
	ViewThemeSelect
)

// MenuItem represents a menu option.
type MenuItem struct {
	Label string
	Icon  string
	Desc  string
}

// ScanResult holds the result of an async repository scan.
type ScanResult struct {
	MoodLabel  string
	MoodEmoji  string
	Confidence float64
	Err        error
}

// ScanFunc is a function that performs a repository scan and returns mood info.
type ScanFunc func() ScanResult

// Model is the main Bubble Tea model for the Moodwave TUI.
type Model struct {
	CurrentView View
	Width       int
	Height      int
	Frame       int

	// Menu
	MenuItems  []MenuItem
	MenuCursor int

	// Mood
	MoodLabel      string
	MoodConfidence float64
	MoodEmoji      string

	// Playback
	TrackTitle          string
	TrackArtist         string
	IsPlaying           bool
	IsResolving         bool // true while resolving next track (prevents double-trigger)
	ResolveRetried      bool // true after one resolve retry has already been used for the current track
	ConsecutiveFailures int  // resolve/backend failures in a row; guards against looping forever on a permanently broken track under repeat-one
	Elapsed             time.Duration
	Duration            time.Duration
	PlayStarted         time.Time
	PausedAt            time.Time     // when pause started (zero if not paused)
	PausedTotal         time.Duration // total time spent paused
	Controller          *playback.Controller
	// PQ is the single source of truth for "what's playing and what's
	// queued" for auto-generated content (search results, autonomous play,
	// endless-listening refill). See queue.go for why this replaced two
	// independently-mutated []recommender.Candidate slices.
	PQ PlaybackQueue
	// Personal is a separate, deliberately-curated playlist that is always
	// consumed FIRST, ahead of PQ — see queue.go for the ownership rule.
	// Building this while music from PQ is playing does not touch PQ at
	// all; it only affects what plays once the current track ends.
	Personal PersonalPlaylist
	Backend  string

	// Visual
	VisualMode string
	RepeatMode string // "off", "one", "all"

	// Search input
	SearchInput string

	// Search results
	SearchResults      []string
	SearchResultCursor int
	SearchCandidates   []recommender.Candidate // raw candidates from last search
	// ScanForAutoPlay is scoped narrowly to "this repo scan was triggered by
	// Autonomous Play, so search-and-take-over-queue once it completes" —
	// it is only ever read by the ScanResultMsg handler, never by the
	// search-result handler, so it cannot race with a manual search intent
	// the way the old shared AutoPlayNext flag did.
	ScanForAutoPlay bool

	// Status messages
	StatusMsg string

	// Playlist Manager view state — browses/edits the same PQ above.
	// PlaylistCursor is the browse position in the Playlist Manager list,
	// independent of PQ's playing position (you can scroll to look at
	// track 5 while track 2 is still playing).
	PlaylistCursor int
	PlaylistMode   string // "view", "add"

	// Personal Playlist view state — the Personal Playlist has its own
	// dedicated window (ViewPersonalPlaylist), separate from the Live
	// Queue Playlist Manager. PersonalCursor is the browse position
	// within that window, independent of anything else.
	PersonalCursor int

	// Theme selection
	ThemeCursor int

	// Self-update — checked once automatically on every launch (see
	// Init/UpdateCheckCmd), independently of any menu action. UpdateChecked
	// only flips true once the check has actually completed (success or
	// failure), so the home screen never shows a stale/default state.
	UpdateChecked        bool
	UpdateAvailable      bool
	UpdateCurrentVersion string
	UpdateLatestVersion  string
	UpdateRelease        *updater.Release
	UpdateInProgress     bool
	UpdateErr            string

	// Callbacks — the host app sets these to wire up actual logic
	OnAction func(action string, args string)
	ScanFn   ScanFunc
	Quitting bool
}

// tickMsg is sent on each animation frame.
type tickMsg time.Time

// statusClearMsg clears the status message.
type statusClearMsg struct{}

// NewModel creates the initial TUI model.
func NewModel() Model {
	return Model{
		CurrentView: ViewHome,
		MenuItems: []MenuItem{
			{Label: "Autonomous Play", Icon: "▶", Desc: "Scan & play mood-matched music"},
			{Label: "YouTube Search", Icon: "◉", Desc: "Search and play from YouTube"},
			{Label: "Playlist", Icon: "☰", Desc: "Create & manage session playlist"},
			{Label: "Scan Repository", Icon: "◈", Desc: "Analyze codebase mood signals"},
			{Label: "Customize Theme", Icon: "◆", Desc: "Switch visual theme"},
			{Label: "System Doctor", Icon: "◇", Desc: "Run diagnostics"},
			{Label: "Exit", Icon: "□", Desc: "Quit Moodwave"},
		},
		MenuCursor: 0,
		VisualMode: "spectrum",
		PQ:         NewPlaybackQueue(),
		Personal:   NewPersonalPlaylist(),
	}
}

// Init starts the animation tick and kicks off a one-shot background
// check for a newer release. The update check runs independently of
// anything else in the app — it never blocks startup, and a failure or
// slow network just means the banner never appears, nothing else changes.
func (m Model) Init() tea.Cmd {
	return tea.Batch(tickCmd(), tea.WindowSize(), UpdateCheckCmd())
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*42, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Update handles messages.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case tickMsg:
		m.Frame++
		// Update elapsed time in real-time during playback (excluding paused time)
		if m.IsPlaying && !m.PlayStarted.IsZero() {
			m.Elapsed = time.Since(m.PlayStarted) - m.PausedTotal
			if m.Duration > 0 && m.Elapsed > m.Duration {
				m.Elapsed = m.Duration
			}
			// Detect track end — only if not already handling it
			if !m.IsResolving && m.Duration > 0 && m.Elapsed >= m.Duration && m.Controller != nil {
				status := m.Controller.Status()
				if status.State == playback.StatePlaying {
					if cmd := m.resolveOutcome(outcomeEnded, nil); cmd != nil {
						return m, tea.Batch(cmd, tickCmd())
					}
				}
			}
		}
		return m, tickCmd()

	case statusClearMsg:
		m.StatusMsg = ""
		return m, nil

	// ── Async playback messages ──
	case SearchResultMsg:
		if msg.Err != nil {
			m.StatusMsg = "Search failed: " + msg.Err.Error()
			m.CurrentView = ViewHome
			return m, tea.Tick(3*time.Second, func(t time.Time) tea.Msg { return statusClearMsg{} })
		}
		if len(msg.Candidates) == 0 {
			m.StatusMsg = "No results found for: " + msg.Query
			m.CurrentView = ViewHome
			return m, tea.Tick(3*time.Second, func(t time.Time) tea.Msg { return statusClearMsg{} })
		}

		switch msg.Intent {
		case SearchIntentAutoPlayAll, SearchIntentRefillQueue:
			// This search's own intent (not shared/racy model state) says
			// to take over the queue and start playing — used by
			// autonomous play and by the endless-listening refill when the
			// queue runs dry. A manual search the user typed can never
			// land here regardless of timing, because its message was
			// built with SearchIntentBrowse from the start.
			m.PQ.Replace(msg.Candidates, 0)
			m.SearchResults = make([]string, len(msg.Candidates))
			for i, c := range msg.Candidates {
				m.SearchResults[i] = c.DisplayName()
			}
			m.CurrentView = ViewPlaying
			m.TrackTitle = m.PQ.Current().DisplayName()
			m.TrackArtist = "Resolving..."
			m.PlayStarted = time.Time{}
			m.Elapsed = 0
			m.IsResolving = true
			m.ResolveRetried = false
			m.ConsecutiveFailures = 0
			m.StatusMsg = "Playing similar..."
			return m, ResolveAndPlayCmd(m.PQ.Current(), m.Backend)

		default: // SearchIntentBrowse
			// A manual search NEVER touches the live queue — it only ever
			// shows a results list for the user to pick from explicitly.
			// This is what actually fixes "searching one song added 10+
			// songs to my playlist": browsing results and adding/playing
			// a specific one are separate, deliberate actions now, never
			// an implicit side effect of the search itself.
			m.SearchCandidates = msg.Candidates
			m.SearchResults = make([]string, len(msg.Candidates))
			for i, c := range msg.Candidates {
				m.SearchResults[i] = c.DisplayName()
			}
			m.SearchResultCursor = 0
			m.CurrentView = ViewSearchResults
			m.StatusMsg = fmt.Sprintf("Found %d results", len(msg.Candidates))
			return m, nil
		}

	case PlaybackStartedMsg:
		// Stop any previously playing controller to prevent double audio
		if m.Controller != nil && m.Controller != msg.Controller {
			m.Controller.Stop()
		}
		m.Controller = msg.Controller
		m.IsPlaying = true
		m.IsResolving = false // resolution complete, playback started
		m.ResolveRetried = false // re-arm the one-shot retry budget for this new track
		m.ConsecutiveFailures = 0 // a successful start clears the failure streak
		m.PlayStarted = time.Now()
		m.Elapsed = 0
		m.PausedAt = time.Time{}
		m.PausedTotal = 0
		m.StatusMsg = ""
		if m.Controller != nil {
			status := m.Controller.Status()
			m.TrackTitle = status.Title
			m.TrackArtist = status.Artist
		}
		// Set duration from the candidate's track info
		if candidate := m.PQ.Current(); candidate != nil {
			if candidate.Track != nil && candidate.Track.Duration > 0 {
				m.Duration = time.Duration(candidate.Track.Duration) * time.Second
			} else {
				m.Duration = 0 // unknown duration (live stream)
			}
		}
		return m, PollPlaybackCmd(m.Controller)

	case PlaybackErrorMsg:
		return m, m.resolveOutcome(outcomeErrored, msg.Err)

	case PlaybackStatusMsg:
		// Discard stale status from a controller that's already been replaced
		if msg.Controller != m.Controller {
			return m, nil
		}
		if msg.Status.State == playback.StatePlaying {
			m.IsPlaying = true
		} else if msg.Status.State == playback.StatePaused {
			m.IsPlaying = false
		}
		// Keep polling regardless of current view (audio plays in background)
		if m.Controller != nil {
			return m, PollPlaybackCmd(m.Controller)
		}
		return m, nil

	case PlaybackEndedMsg:
		// Discard stale end signal from a controller that's already been replaced
		// (nil Controller in the message means "no controller was ever set" —
		// only honor that case if we truly have none right now)
		if msg.Controller != nil && msg.Controller != m.Controller {
			return m, nil
		}
		if msg.Controller == nil && m.Controller != nil {
			return m, nil
		}
		if m.IsResolving {
			return m, nil
		}
		if !m.IsPlaying && m.Controller == nil {
			return m, nil
		}
		if m.PQ.Current() == nil {
			if m.Controller != nil {
				m.Controller.Stop()
				m.Controller = nil
			}
			m.IsPlaying = false
			return m, nil
		}
		return m, m.resolveOutcome(outcomeEnded, nil)

	case ScanResultMsg:
		if msg.Err != nil {
			m.StatusMsg = "Scan failed: " + msg.Err.Error()
			m.CurrentView = ViewHome
			m.ScanForAutoPlay = false
			return m, tea.Tick(3*time.Second, func(t time.Time) tea.Msg { return statusClearMsg{} })
		}
		m.MoodLabel = msg.MoodLabel
		m.MoodEmoji = msg.MoodEmoji
		m.MoodConfidence = msg.Confidence
		m.StatusMsg = fmt.Sprintf("Mood: %s %s (%.0f%%)", msg.MoodEmoji, msg.MoodLabel, msg.Confidence*100)

		// If autonomous play triggered this scan (as opposed to a plain
		// "Scan Repository" menu selection), now search for mood-matched
		// music and take over the queue. ScanForAutoPlay is scoped only to
		// this scan→search bridge — unlike the old AutoPlayNext, it is
		// never read by the search-result handler, so it can't collide
		// with a manual browse search the way that shared flag did.
		if m.ScanForAutoPlay {
			m.ScanForAutoPlay = false
			query := buildMoodQuery(m.MoodLabel)
			m.StatusMsg = fmt.Sprintf("Finding %s music...", m.MoodLabel)
			return m, SearchCmd(query, SearchIntentAutoPlayAll)
		}
		m.CurrentView = ViewHome
		return m, tea.Tick(3*time.Second, func(t time.Time) tea.Msg { return statusClearMsg{} })

	case UpdateCheckMsg:
		m.UpdateChecked = true
		if msg.Err != nil {
			// Silent failure — never interrupt the user with a background
			// check error. The banner simply doesn't appear.
			m.UpdateAvailable = false
			return m, nil
		}
		m.UpdateAvailable = msg.Available
		m.UpdateCurrentVersion = msg.CurrentVersion
		m.UpdateLatestVersion = msg.LatestVersion
		m.UpdateRelease = msg.Release
		return m, nil

	case UpdateApplyMsg:
		m.UpdateInProgress = false
		if msg.Err != nil {
			m.UpdateErr = msg.Err.Error()
			m.StatusMsg = "Update failed: " + msg.Err.Error()
			return m, tea.Tick(4*time.Second, func(t time.Time) tea.Msg { return statusClearMsg{} })
		}
		m.UpdateAvailable = false
		m.StatusMsg = fmt.Sprintf("Updated to %s — restart moodwave to use it", msg.NewVersion)
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	// Global keys
	switch key {
	case "ctrl+c":
		m.Quitting = true
		return m, tea.Quit
	}

	switch m.CurrentView {
	case ViewHome:
		return m.handleHomeKey(key)
	case ViewSearchInput:
		return m.handleSearchInputKey(msg)
	case ViewSearchResults:
		return m.handleSearchResultsKey(key)
	case ViewPlaying:
		return m.handlePlayingKey(key)
	case ViewPlaylist:
		return m.handlePlaylistKey(key)
	case ViewPersonalPlaylist:
		return m.handlePersonalPlaylistKey(key)
	case ViewThemeSelect:
		return m.handleThemeKey(key)
	case ViewScanning:
		if key == "q" {
			m.CurrentView = ViewHome
		}
		return m, nil
	}
	return m, nil
}

func (m Model) handleHomeKey(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "q":
		m.Quitting = true
		return m, tea.Quit
	case "up", "k":
		m.MenuCursor--
		if m.MenuCursor < 0 {
			m.MenuCursor = len(m.MenuItems) - 1
		}
	case "down", "j":
		m.MenuCursor++
		if m.MenuCursor >= len(m.MenuItems) {
			m.MenuCursor = 0
		}
	case "v":
		m.cycleVisual()
	case "t":
		m.CurrentView = ViewThemeSelect
		return m, nil
	case "u":
		// Install the update that was found by the background check on
		// launch. No-op if nothing's available or one is already running,
		// so mashing 'u' can't fire multiple concurrent downloads.
		if m.UpdateAvailable && !m.UpdateInProgress && m.UpdateRelease != nil {
			m.UpdateInProgress = true
			m.StatusMsg = fmt.Sprintf("Downloading v%s...", m.UpdateLatestVersion)
			return m, UpdateApplyCmd(m.UpdateRelease)
		}
	case "enter":
		return m.handleMenuSelect()
	}
	return m, nil
}

func (m *Model) cycleVisual() {
	modes := []string{
		"spectrum", "fire", "matrix", "rain", "plasma",
		"aurora", "snow", "fireflies", "lava", "dna",
		"campfire", "smoke", "ripple", "lightning", "spiral",
		"petals", "galaxy", "zen", "heartbeat", "pendulum",
		"city", "ocean", "starfield", "pulse",
	}
	current := 0
	for i, mode := range modes {
		if mode == m.VisualMode {
			current = i
			break
		}
	}
	m.VisualMode = modes[(current+1)%len(modes)]
}

func (m Model) handleMenuSelect() (tea.Model, tea.Cmd) {
	switch m.MenuCursor {
	case 0: // Autonomous Play
		// If no mood detected, scan first
		if m.MoodLabel == "" && m.ScanFn != nil {
			m.CurrentView = ViewScanning
			m.StatusMsg = "Scanning codebase for mood..."
			m.ScanForAutoPlay = true // ScanResultMsg will fire the search once this completes
			return m, ScanRepoCmd(m.ScanFn)
		}
		// Build a rich query from mood traits
		query := buildMoodQuery(m.MoodLabel)
		m.CurrentView = ViewScanning
		m.StatusMsg = fmt.Sprintf("Finding %s music...", m.MoodLabel)
		return m, SearchCmd(query, SearchIntentAutoPlayAll)
	case 1: // YouTube Search
		m.CurrentView = ViewSearchInput
		m.SearchInput = ""
	case 2: // Playlist
		m.CurrentView = ViewPlaylist
		m.PlaylistMode = "view"
	case 3: // Scan Repository
		m.CurrentView = ViewScanning
		m.StatusMsg = "Scanning repository..."
		return m, ScanRepoCmd(m.ScanFn)
	case 4: // Customize Theme
		m.CurrentView = ViewThemeSelect
		return m, nil
	case 5: // System Doctor
		if m.OnAction != nil {
			m.OnAction("doctor", "")
		}
		m.Quitting = true
		return m, tea.Quit
	case 6: // Exit
		if m.Controller != nil {
			m.Controller.Stop()
		}
		m.Quitting = true
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) handleSearchInputKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()
	switch key {
	case "esc":
		m.CurrentView = ViewHome
		m.SearchInput = ""
	case "enter":
		if m.SearchInput != "" {
			// Fire async search — stay in TUI. Explicitly a browse search:
			// this must only ever show a results list, never take over
			// the live queue by itself.
			m.CurrentView = ViewScanning
			m.StatusMsg = "Searching: " + m.SearchInput
			return m, SearchCmd(m.SearchInput, SearchIntentBrowse)
		}
	case "backspace":
		if len(m.SearchInput) > 0 {
			m.SearchInput = m.SearchInput[:len(m.SearchInput)-1]
		}
	default:
		if len(key) == 1 && key[0] >= 32 && key[0] <= 126 {
			m.SearchInput += key
		} else if key == "space" {
			m.SearchInput += " "
		}
	}
	return m, nil
}

func (m Model) handleSearchResultsKey(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "esc":
		// Go back: if playing return to playing, else home
		if m.Controller != nil && (m.IsPlaying || m.IsResolving) {
			m.CurrentView = ViewPlaying
		} else {
			m.CurrentView = ViewHome
		}
		m.SearchResults = nil
		m.SearchCandidates = nil
	case "up", "k":
		if len(m.SearchResults) > 0 {
			m.SearchResultCursor--
			if m.SearchResultCursor < 0 {
				m.SearchResultCursor = len(m.SearchResults) - 1
			}
		}
	case "down", "j":
		if len(m.SearchResults) > 0 {
			m.SearchResultCursor++
			if m.SearchResultCursor >= len(m.SearchResults) {
				m.SearchResultCursor = 0
			}
		}
	case "p":
		// Add selected track to the Personal playlist — a separate list
		// that always plays first, ahead of the live queue, once the
		// current track ends. Does NOT touch what's playing right now.
		if len(m.SearchCandidates) > 0 && m.SearchResultCursor < len(m.SearchCandidates) {
			track := m.SearchCandidates[m.SearchResultCursor]
			m.Personal.Add(track)
			m.StatusMsg = fmt.Sprintf("Added to personal playlist (%d) — plays next", m.Personal.Len())
			return m, tea.Tick(2*time.Second, func(t time.Time) tea.Msg { return statusClearMsg{} })
		}
	case "enter":
		if len(m.SearchCandidates) > 0 && m.SearchResultCursor < len(m.SearchCandidates) {
			if m.PlaylistMode == "add" {
				// Add to the Personal playlist, then go back to the
				// playlist view. Does not touch the currently playing
				// track — this track will play once the personal
				// playlist's turn comes up (see advanceIgnoringRepeatOne).
				track := m.SearchCandidates[m.SearchResultCursor]
				m.Personal.Add(track)
				m.StatusMsg = fmt.Sprintf("Added! (%d in personal playlist)", m.Personal.Len())
				// Always land back on the dedicated Personal Playlist
				// window — that's the one place this track is guaranteed
				// to show up, one row at a time, nothing else added.
				m.CurrentView = ViewPersonalPlaylist
				m.PlaylistMode = "view"
				m.PersonalCursor = 0
				return m, tea.Tick(2*time.Second, func(t time.Time) tea.Msg { return statusClearMsg{} })
			}
			// Normal mode: stop current, switch to new track
			if m.Controller != nil {
				m.Controller.Stop()
				m.Controller = nil
			}
			// Replace queue with search results and play selected
			m.PQ.Replace(m.SearchCandidates, m.SearchResultCursor)
			m.CurrentView = ViewPlaying
			m.TrackTitle = m.PQ.Current().DisplayName()
			m.TrackArtist = "Resolving stream..."
			m.IsPlaying = false
			m.IsResolving = true
			m.StatusMsg = "Resolving audio stream..."
			return m, ResolveAndPlayCmd(m.PQ.Current(), m.Backend)
		}
	}
	return m, nil
}

func (m Model) handlePlayingKey(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "q":
		// Stop playback and quit
		if m.Controller != nil {
			m.Controller.Stop()
			m.Controller = nil
		}
		m.Quitting = true
		return m, tea.Quit
	case "esc":
		// Stop playback and go back to search results if available, else home
		if m.Controller != nil {
			m.Controller.Stop()
			m.Controller = nil
		}
		m.IsPlaying = false
		m.PlayStarted = time.Time{}
		if len(m.SearchResults) > 0 {
			m.CurrentView = ViewSearchResults
		} else {
			m.CurrentView = ViewHome
		}
	case " ":
		if m.Controller != nil {
			if m.IsPlaying {
				m.Controller.Pause()
				m.IsPlaying = false
				m.PausedAt = time.Now()
			} else {
				m.Controller.Resume()
				m.IsPlaying = true
				// Add paused duration to total
				if !m.PausedAt.IsZero() {
					m.PausedTotal += time.Since(m.PausedAt)
					m.PausedAt = time.Time{}
				}
			}
		}
	case "n":
		// Next track — routed through the same single outcome handler as
		// every other "stop playing this track" path.
		cmd := m.resolveOutcome(outcomeManualSkip, nil)
		if cmd == nil {
			m.StatusMsg = "End of queue"
			return m, tea.Tick(2*time.Second, func(t time.Time) tea.Msg { return statusClearMsg{} })
		}
		m.StatusMsg = "Next track..."
		return m, cmd
	case "v":
		m.cycleVisual()
	case "t":
		m.CurrentView = ViewThemeSelect
		return m, nil
	case "l":
		// Cycle repeat mode: off → one → all → off
		switch m.RepeatMode {
		case "", "off":
			m.RepeatMode = "one"
			m.StatusMsg = "🔂 Repeat: Track"
		case "one":
			m.RepeatMode = "all"
			m.StatusMsg = "🔁 Repeat: All"
		case "all":
			m.RepeatMode = "off"
			m.StatusMsg = "Repeat: Off"
		}
		return m, tea.Tick(2*time.Second, func(t time.Time) tea.Msg { return statusClearMsg{} })
	case "p":
		// Open interactive playlist manager (audio keeps playing in background)
		m.CurrentView = ViewPlaylist
		m.PlaylistMode = "view"
	case "o":
		// Open the dedicated Personal Playlist window.
		m.CurrentView = ViewPersonalPlaylist
		m.PersonalCursor = 0
	}
	return m, nil
}

func (m Model) handlePlaylistKey(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "esc":
		// Go back: if playing, return to playing view; else home
		if m.Controller != nil && m.IsPlaying {
			m.CurrentView = ViewPlaying
		} else {
			m.CurrentView = ViewHome
		}
	case "up", "k":
		if m.PQ.Len() > 0 {
			m.PlaylistCursor--
			if m.PlaylistCursor < 0 {
				m.PlaylistCursor = m.PQ.Len() - 1
			}
		}
	case "down", "j":
		if m.PQ.Len() > 0 {
			m.PlaylistCursor++
			if m.PlaylistCursor >= m.PQ.Len() {
				m.PlaylistCursor = 0
			}
		}
	case "enter":
		// Jump playback to the browsed position in the (single) queue.
		if m.PQ.Len() > 0 && m.PQ.SelectIndex(m.PlaylistCursor) {
			if m.Controller != nil {
				m.Controller.Stop()
				m.Controller = nil
			}
			m.CurrentView = ViewPlaying
			m.TrackTitle = m.PQ.Current().DisplayName()
			m.TrackArtist = "Resolving..."
			m.PlayStarted = time.Time{}
			m.Elapsed = 0
			m.IsResolving = true
			m.StatusMsg = "Playing playlist..."
			return m, ResolveAndPlayCmd(m.PQ.Current(), m.Backend)
		}
	case "d", "delete", "backspace":
		// Remove the browsed track from the queue.
		if m.PQ.RemoveAt(m.PlaylistCursor) {
			if m.PlaylistCursor >= m.PQ.Len() && m.PlaylistCursor > 0 {
				m.PlaylistCursor--
			}
			m.StatusMsg = fmt.Sprintf("Removed — %d tracks remain", m.PQ.Len())
			return m, tea.Tick(2*time.Second, func(t time.Time) tea.Msg { return statusClearMsg{} })
		}
	case "a":
		// Open search to add tracks to the Personal playlist.
		m.CurrentView = ViewSearchInput
		m.SearchInput = ""
		m.PlaylistMode = "add"
	case "o":
		// Jump to the dedicated Personal Playlist window.
		m.CurrentView = ViewPersonalPlaylist
		m.PersonalCursor = 0
	case "+":
		// Quick "add" of the currently playing track is a no-op now — it's
		// already in the queue by definition. Kept as a friendly no-op
		// status message instead of silently doing nothing.
		if m.PQ.Current() != nil {
			m.StatusMsg = "Already in the queue"
			return m, tea.Tick(2*time.Second, func(t time.Time) tea.Msg { return statusClearMsg{} })
		}
	}
	return m, nil
}

// handlePersonalPlaylistKey drives the dedicated Personal Playlist window.
// This window is deliberately separate from the Live Queue Playlist
// Manager: it only ever shows tracks that are actually loaded into
// m.Personal right now, one row per track, added and removed one at a
// time — never a bulk dump of search results. Nothing here ever touches
// m.PQ or the currently playing track.
func (m Model) handlePersonalPlaylistKey(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "esc":
		// Go back: if playing, return to playing view; else home
		if m.Controller != nil && m.IsPlaying {
			m.CurrentView = ViewPlaying
		} else {
			m.CurrentView = ViewHome
		}
	case "up", "k":
		if m.Personal.Len() > 0 {
			m.PersonalCursor--
			if m.PersonalCursor < 0 {
				m.PersonalCursor = m.Personal.Len() - 1
			}
		}
	case "down", "j":
		if m.Personal.Len() > 0 {
			m.PersonalCursor++
			if m.PersonalCursor >= m.Personal.Len() {
				m.PersonalCursor = 0
			}
		}
	case "a":
		// Search and add exactly one track to the personal playlist.
		m.CurrentView = ViewSearchInput
		m.SearchInput = ""
		m.PlaylistMode = "add"
	case "d", "delete", "backspace":
		// Remove exactly the one track under the cursor — the list
		// re-renders from the live m.Personal state, so a removed track
		// disappears from view immediately and nothing else shifts in
		// unexpectedly.
		if m.Personal.RemoveAt(m.PersonalCursor) {
			if m.PersonalCursor >= m.Personal.Len() && m.PersonalCursor > 0 {
				m.PersonalCursor--
			}
			m.StatusMsg = fmt.Sprintf("Removed — %d left in personal playlist", m.Personal.Len())
			return m, tea.Tick(2*time.Second, func(t time.Time) tea.Msg { return statusClearMsg{} })
		}
	case "c":
		// Clear the Personal playlist only — never touches the live queue
		// (m.PQ) or whatever is currently playing.
		if m.Personal.Len() > 0 {
			n := m.Personal.Len()
			m.Personal.Clear()
			m.PersonalCursor = 0
			m.StatusMsg = fmt.Sprintf("Cleared personal playlist (%d tracks removed)", n)
			return m, tea.Tick(2*time.Second, func(t time.Time) tea.Msg { return statusClearMsg{} })
		}
	}
	return m, nil
}

func (m Model) handleThemeKey(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "esc":
		if m.Controller != nil && m.IsPlaying {
			m.CurrentView = ViewPlaying
		} else {
			m.CurrentView = ViewHome
		}
	case "up", "k":
		m.ThemeCursor--
		if m.ThemeCursor < 0 {
			m.ThemeCursor = len(AllThemes) - 1
		}
		ApplyTheme(AllThemes[m.ThemeCursor])
	case "down", "j":
		m.ThemeCursor++
		if m.ThemeCursor >= len(AllThemes) {
			m.ThemeCursor = 0
		}
		ApplyTheme(AllThemes[m.ThemeCursor])
	case "enter":
		ApplyTheme(AllThemes[m.ThemeCursor])
		m.StatusMsg = "Theme: " + AllThemes[m.ThemeCursor].Name
		if m.Controller != nil && m.IsPlaying {
			m.CurrentView = ViewPlaying
		} else {
			m.CurrentView = ViewHome
		}
		return m, tea.Tick(2*time.Second, func(t time.Time) tea.Msg { return statusClearMsg{} })
	}
	return m, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// View rendering
// ─────────────────────────────────────────────────────────────────────────────

func (m Model) View() string {
	if m.Quitting {
		return ""
	}
	if m.Width == 0 {
		return "Loading..."
	}

	switch m.CurrentView {
	case ViewHome:
		return m.viewHome()
	case ViewPlaying:
		return m.viewPlaying()
	case ViewSearchInput:
		return m.viewSearchInput()
	case ViewSearchResults:
		return m.viewSearchResults()
	case ViewPlaylist:
		return m.viewPlaylist()
	case ViewPersonalPlaylist:
		return m.viewPersonalPlaylist()
	case ViewThemeSelect:
		return m.viewThemeSelect()
	case ViewScanning:
		return m.viewScanning()
	default:
		return m.viewHome()
	}
}

func (m Model) viewHome() string {
	var sections []string

	moodColor := ColorHighlight
	if m.MoodLabel != "" {
		moodColor = MoodColor(m.MoodLabel)
	}

	// Banner — always uses theme highlight color
	sections = append(sections, RenderBanner(m.Width, ColorHighlight))
	sections = append(sections, "")

	// Update banner — only shown once the background check (fired from
	// Init on every launch) has actually found something newer. Silent
	// otherwise: no banner, no interruption, nothing to dismiss.
	if m.UpdateAvailable {
		updateMsg := fmt.Sprintf("⬆ Update available: v%s → v%s   [U] Install now", m.UpdateCurrentVersion, m.UpdateLatestVersion)
		if m.UpdateInProgress {
			updateMsg = fmt.Sprintf("⬆ Installing v%s...", m.UpdateLatestVersion)
		}
		updateBanner := lipgloss.NewStyle().
			Foreground(ColorSuccess).Bold(true).
			Render(updateMsg)
		sections = append(sections, lipgloss.PlaceHorizontal(m.Width, lipgloss.Center, updateBanner))
		sections = append(sections, "")
	}

	// Menu panel
	menuPanel := m.renderMenu(moodColor)
	sections = append(sections, lipgloss.PlaceHorizontal(m.Width, lipgloss.Center, menuPanel))
	sections = append(sections, "")

	// Pixel pet animation
	sections = append(sections, "")
	sections = append(sections, RenderPixelPet(m.Width, m.Frame))

	// Status msg
	if m.StatusMsg != "" {
		sections = append(sections, "")
		msg := lipgloss.NewStyle().Foreground(ColorHighlight).Render("  " + m.StatusMsg)
		sections = append(sections, msg)
	}

	// Hints
	sections = append(sections, "")
	sections = append(sections, lipgloss.PlaceHorizontal(m.Width, lipgloss.Center, m.renderHints()))

	// Fill + status bar
	content := strings.Join(sections, "\n")
	content = m.fillAndStatusBar(content)
	return content
}

func (m Model) viewPlaying() string {
	var sections []string
	palette := GetPaletteForMood(m.MoodLabel)

	// Header
	moodStr := ""
	if m.MoodLabel != "" {
		moodStr = fmt.Sprintf(" · %s %s", m.MoodEmoji, m.MoodLabel)
	}
	state := "▶ Playing"
	if !m.IsPlaying {
		state = "⏸ Paused"
	}
	header := lipgloss.NewStyle().Foreground(ColorHighlight).Bold(true).
		Render(fmt.Sprintf("  %s%s", state, moodStr))
	sections = append(sections, header)
	sections = append(sections, "")

	// Track info
	title := m.TrackTitle
	if title == "" {
		title = "Waiting for track..."
	}
	artist := m.TrackArtist
	titleStyle := lipgloss.NewStyle().Foreground(ColorFg).Bold(true)
	artistStyle := lipgloss.NewStyle().Foreground(ColorDim)
	trackPanel := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorder).
		Padding(0, 2).Width(m.Width - 6).
		Render(titleStyle.Render(title) + "\n" + artistStyle.Render(artist))
	sections = append(sections, trackPanel)
	sections = append(sections, "")

	// Visualizer — FULL WIDTH, dynamically sized based on terminal
	visWidth := m.Width - 2
	if visWidth < 20 {
		visWidth = 20
	}
	// Reserve space for: header(2) + track(4) + progress(1) + status(2) + hints(2) + statusbar(1) = ~12 lines
	visHeight := m.Height - 12
	if visHeight > 22 {
		visHeight = 22
	}
	if visHeight < 6 {
		visHeight = 6
	}

	var vis string
	switch m.VisualMode {
	case "fire":
		vis = RenderDoomFire(visWidth, visHeight, m.Frame)
	case "matrix":
		vis = RenderMatrix(visWidth, visHeight, m.Frame)
	case "rain":
		vis = RenderRain(visWidth, visHeight, m.Frame)
	case "plasma":
		vis = RenderPlasma(visWidth, visHeight, m.Frame)
	case "aurora":
		vis = RenderAurora(visWidth, visHeight, m.Frame)
	case "snow":
		vis = RenderSnow(visWidth, visHeight, m.Frame)
	case "fireflies":
		vis = RenderFireflies(visWidth, visHeight, m.Frame)
	case "lava":
		vis = RenderLava(visWidth, visHeight, m.Frame)
	case "dna":
		vis = RenderDNA(visWidth, visHeight, m.Frame)
	case "campfire":
		vis = RenderDoomFire(visWidth, visHeight, m.Frame) // reuse doom fire
	case "smoke":
		vis = RenderSmoke(visWidth, visHeight, m.Frame)
	case "ripple":
		vis = RenderRipple(visWidth, visHeight, m.Frame)
	case "lightning":
		vis = RenderLightning(visWidth, visHeight, m.Frame)
	case "spiral":
		vis = RenderSpiral(visWidth, visHeight, m.Frame)
	case "petals":
		vis = RenderPetals(visWidth, visHeight, m.Frame)
	case "galaxy":
		vis = RenderStarfield(visWidth, visHeight, m.Frame) // reuse starfield
	case "zen":
		vis = RenderOceanWaves(visWidth, visHeight, m.Frame) // reuse ocean
	case "heartbeat":
		vis = RenderHeartbeat(visWidth, visHeight, m.Frame)
	case "pendulum":
		vis = RenderPendulum(visWidth, visHeight, m.Frame)
	case "city":
		vis = RenderCityNight(visWidth, visHeight, m.Frame)
	case "ocean":
		vis = RenderOceanWaves(visWidth, visHeight, m.Frame)
	case "starfield":
		vis = RenderStarfield(visWidth, visHeight, m.Frame)
	case "pulse":
		vis = RenderFullVisualizer(visWidth, visHeight, m.Frame, 0.9, AuroraPalette)
	default: // "spectrum"
		vis = RenderFullSpectrumBars(visWidth, visHeight, m.Frame, 0.85, palette)
	}
	sections = append(sections, vis)
	sections = append(sections, "")

	// Progress bar
	sections = append(sections, lipgloss.PlaceHorizontal(m.Width, lipgloss.Center, m.renderProgress()))

	// Status
	if m.StatusMsg != "" {
		sections = append(sections, "")
		msg := lipgloss.NewStyle().Foreground(ColorHighlight).Render("  " + m.StatusMsg)
		sections = append(sections, msg)
	}

	// Controls
	sections = append(sections, "")
	sections = append(sections, lipgloss.PlaceHorizontal(m.Width, lipgloss.Center, m.renderPlayingHints()))

	content := strings.Join(sections, "\n")
	content = m.fillAndStatusBar(content)
	return content
}

func (m Model) viewSearchInput() string {
	var sections []string

	// Header
	header := lipgloss.NewStyle().Foreground(ColorHighlight).Bold(true).
		Render("  ◉ YouTube Search")
	sections = append(sections, "")
	sections = append(sections, header)
	sections = append(sections, "")

	// Input box
	prompt := lipgloss.NewStyle().Foreground(ColorDim).Render("  Search: ")
	cursor := lipgloss.NewStyle().Foreground(ColorHighlight).Bold(true).Render("█")
	inputText := lipgloss.NewStyle().Foreground(ColorFg).Render(m.SearchInput)

	inputLine := prompt + inputText + cursor

	inputPanel := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorHighlight).
		Padding(1, 2).Width(m.Width - 10).
		Render(inputLine)
	sections = append(sections, lipgloss.PlaceHorizontal(m.Width, lipgloss.Center, inputPanel))
	sections = append(sections, "")

	// Pixel pet
	sections = append(sections, RenderPixelPet(m.Width, m.Frame))
	sections = append(sections, "")

	// Hints
	keyStyle := lipgloss.NewStyle().Foreground(ColorHighlight).Bold(true)
	descStyle := lipgloss.NewStyle().Foreground(ColorDim)
	hints := keyStyle.Render("[Enter]") + descStyle.Render(" Search & Play") +
		"  ·  " + keyStyle.Render("[Esc]") + descStyle.Render(" Back")
	sections = append(sections, lipgloss.PlaceHorizontal(m.Width, lipgloss.Center, hints))

	content := strings.Join(sections, "\n")
	content = m.fillAndStatusBar(content)
	return content
}

func (m Model) viewSearchResults() string {
	var sections []string

	header := lipgloss.NewStyle().Foreground(ColorHighlight).Bold(true).
		Render("  ◉ Search Results — Select a track")
	sections = append(sections, "")
	sections = append(sections, header)
	sections = append(sections, "")

	// Results list
	var items []string
	maxShow := 8
	if len(m.SearchResults) < maxShow {
		maxShow = len(m.SearchResults)
	}

	activeStyle := lipgloss.NewStyle().Foreground(ColorHighlight).Bold(true)
	normalStyle := lipgloss.NewStyle().Foreground(ColorFg)
	dimStyle := lipgloss.NewStyle().Foreground(ColorDim)
	cursorStyle := lipgloss.NewStyle().Foreground(ColorHighlight).Bold(true)

	for i := 0; i < maxShow; i++ {
		num := dimStyle.Render(fmt.Sprintf("%d.", i+1))
		if i == m.SearchResultCursor {
			cursor := cursorStyle.Render("▶")
			label := activeStyle.Render(m.SearchResults[i])
			items = append(items, fmt.Sprintf("  %s %s %s", cursor, num, label))
		} else {
			label := normalStyle.Render(m.SearchResults[i])
			items = append(items, fmt.Sprintf("    %s %s", num, label))
		}
	}

	resultPanel := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorder).
		Padding(1, 2).Width(m.Width - 8).
		Render(strings.Join(items, "\n"))
	sections = append(sections, lipgloss.PlaceHorizontal(m.Width, lipgloss.Center, resultPanel))
	sections = append(sections, "")

	// Hints
	k := lipgloss.NewStyle().Foreground(ColorHighlight).Bold(true)
	d := lipgloss.NewStyle().Foreground(ColorDim)
	hints := k.Render("[↑/↓]") + d.Render(" Select") + "  ·  " +
		k.Render("[Enter]") + d.Render(" Play") + "  ·  " +
		k.Render("[P]") + d.Render(" Add to Playlist") + "  ·  " +
		k.Render("[Esc]") + d.Render(" Back")
	sections = append(sections, lipgloss.PlaceHorizontal(m.Width, lipgloss.Center, hints))

	content := strings.Join(sections, "\n")
	content = m.fillAndStatusBar(content)
	return content
}

func (m Model) viewScanning() string {
	var sections []string

	spinChars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	spin := spinChars[m.Frame%len(spinChars)]

	header := lipgloss.NewStyle().Foreground(ColorHighlight).Bold(true).
		Render(fmt.Sprintf("  %s  Scanning repository...", spin))
	sections = append(sections, "")
	sections = append(sections, header)
	sections = append(sections, "")

	waveWidth := m.Width - 2
	palette := GetPaletteForMood(m.MoodLabel)
	wave := RenderInlineWave(waveWidth, m.Frame, palette)
	sections = append(sections, wave)

	content := strings.Join(sections, "\n")
	content = m.fillAndStatusBar(content)
	return content
}

func (m Model) viewPlaylist() string {
	var sections []string

	// If music is playing, show mini now-playing bar at top
	if m.Controller != nil && m.IsPlaying {
		nowPlaying := lipgloss.NewStyle().Foreground(ColorHighlight).
			Render(fmt.Sprintf("  ▶ %s", m.TrackTitle))
		sections = append(sections, nowPlaying)
		sections = append(sections, lipgloss.NewStyle().Foreground(ColorBorder).
			Render("  "+strings.Repeat("─", m.Width-4)))
	}

	header := lipgloss.NewStyle().Foreground(ColorHighlight).Bold(true).
		Render("  ☰ Live Queue")
	subtitle := lipgloss.NewStyle().Foreground(ColorDim).
		Render(fmt.Sprintf("  %d tracks · %d in your personal playlist (see [O])", m.PQ.Len(), m.Personal.Len()))
	sections = append(sections, "")
	sections = append(sections, header)
	sections = append(sections, subtitle)
	sections = append(sections, "")

	if m.PQ.Len() == 0 {
		emptyBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			Padding(2, 4).Width(m.Width - 10).
			Foreground(ColorDim).
			Render("No tracks yet\n\n[A] Search & add to personal playlist\n[O] Open personal playlist\n[Esc] Go back")
		sections = append(sections, lipgloss.PlaceHorizontal(m.Width, lipgloss.Center, emptyBox))
	} else {
		var items []string
		maxShow := m.Height - 12
		if maxShow > 15 {
			maxShow = 15
		}
		if maxShow < 5 {
			maxShow = 5
		}
		if maxShow > m.PQ.Len() {
			maxShow = m.PQ.Len()
		}

		// Scroll offset
		scrollOffset := 0
		if m.PlaylistCursor >= maxShow {
			scrollOffset = m.PlaylistCursor - maxShow + 1
		}

		activeStyle := lipgloss.NewStyle().Foreground(ColorHighlight).Bold(true)
		normalStyle := lipgloss.NewStyle().Foreground(ColorFg)
		dimStyle := lipgloss.NewStyle().Foreground(ColorDim)
		cursorStyle := lipgloss.NewStyle().Foreground(ColorHighlight).Bold(true)
		nowPlayingStyle := lipgloss.NewStyle().Foreground(ColorSuccess)

		for i := scrollOffset; i < scrollOffset+maxShow && i < m.PQ.Len(); i++ {
			num := dimStyle.Render(fmt.Sprintf("%2d.", i+1))
			name := m.PQ.At(i).DisplayName()
			// Truncate
			maxLen := m.Width - 16
			if maxLen > 60 {
				maxLen = 60
			}
			if len(name) > maxLen {
				name = name[:maxLen-3] + "..."
			}

			playingMark := ""
			if i == m.PQ.CurrentIndex() {
				playingMark = nowPlayingStyle.Render(" ♪ playing")
			}

			if i == m.PlaylistCursor {
				cursor := cursorStyle.Render("▶")
				label := activeStyle.Render(name)
				items = append(items, fmt.Sprintf("  %s %s %s%s", cursor, num, label, playingMark))
			} else {
				label := normalStyle.Render(name)
				items = append(items, fmt.Sprintf("    %s %s%s", num, label, playingMark))
			}
		}

		if m.PQ.Len() > scrollOffset+maxShow {
			items = append(items, dimStyle.Render(fmt.Sprintf("    ↓ %d more below", m.PQ.Len()-scrollOffset-maxShow)))
		}

		panel := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			Padding(0, 1).Width(m.Width - 6).
			Render(strings.Join(items, "\n"))
		sections = append(sections, lipgloss.PlaceHorizontal(m.Width, lipgloss.Center, panel))
	}

	// Status
	if m.StatusMsg != "" {
		sections = append(sections, "")
		msg := lipgloss.NewStyle().Foreground(ColorHighlight).Render("  " + m.StatusMsg)
		sections = append(sections, msg)
	}

	sections = append(sections, "")

	// Hints
	k := lipgloss.NewStyle().Foreground(ColorHighlight).Bold(true)
	d := lipgloss.NewStyle().Foreground(ColorDim)
	sep := d.Render("  ·  ")
	hints := k.Render("[Enter]") + d.Render(" Play from here") + sep +
		k.Render("[A]") + d.Render(" Search & Add") + sep +
		k.Render("[D]") + d.Render(" Remove") + sep +
		k.Render("[O]") + d.Render(" Personal Playlist") + sep +
		k.Render("[Esc]") + d.Render(" Back")
	sections = append(sections, lipgloss.PlaceHorizontal(m.Width, lipgloss.Center, hints))

	content := strings.Join(sections, "\n")
	content = m.fillAndStatusBar(content)
	return content
}

// ─────────────────────────────────────────────────────────────────────────────
// Personal Playlist view — a dedicated window separate from the Live
// Queue Playlist Manager above. It shows ONLY the tracks currently loaded
// into m.Personal, one row per track, rendered fresh from that slice on
// every frame. Adding goes through search (one track at a time, never a
// bulk dump); removing a single track makes it disappear immediately
// without any other row shifting unexpectedly. Nothing here ever mutates
// m.PQ or touches the currently playing track — this list only decides
// what plays once its turn comes up (see advanceIgnoringRepeatOne).
// ─────────────────────────────────────────────────────────────────────────────

func (m Model) viewPersonalPlaylist() string {
	var sections []string

	// If music is playing, show mini now-playing bar at top
	if m.Controller != nil && m.IsPlaying {
		nowPlaying := lipgloss.NewStyle().Foreground(ColorHighlight).
			Render(fmt.Sprintf("  ▶ %s", m.TrackTitle))
		sections = append(sections, nowPlaying)
		sections = append(sections, lipgloss.NewStyle().Foreground(ColorBorder).
			Render("  "+strings.Repeat("─", m.Width-4)))
	}

	header := lipgloss.NewStyle().Foreground(ColorSuccess).Bold(true).
		Render("  📋 Personal Playlist")
	subtitle := lipgloss.NewStyle().Foreground(ColorDim).
		Render(fmt.Sprintf("  %d tracks · plays first, in order, before the live queue", m.Personal.Len()))
	sections = append(sections, "")
	sections = append(sections, header)
	sections = append(sections, subtitle)
	sections = append(sections, "")

	if m.Personal.Len() == 0 {
		emptyBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			Padding(2, 4).Width(m.Width - 10).
			Foreground(ColorDim).
			Render("Nothing here yet\n\nSearch and add tracks one at a time —\nthey'll show up here, and only here,\nonce they're actually loaded in.\n\n[A] Search & add a track\n[Esc] Go back")
		sections = append(sections, lipgloss.PlaceHorizontal(m.Width, lipgloss.Center, emptyBox))
	} else {
		var items []string
		maxShow := m.Height - 12
		if maxShow > 15 {
			maxShow = 15
		}
		if maxShow < 5 {
			maxShow = 5
		}
		if maxShow > m.Personal.Len() {
			maxShow = m.Personal.Len()
		}

		scrollOffset := 0
		if m.PersonalCursor >= maxShow {
			scrollOffset = m.PersonalCursor - maxShow + 1
		}

		activeStyle := lipgloss.NewStyle().Foreground(ColorHighlight).Bold(true)
		normalStyle := lipgloss.NewStyle().Foreground(ColorFg)
		dimStyle := lipgloss.NewStyle().Foreground(ColorDim)
		cursorStyle := lipgloss.NewStyle().Foreground(ColorHighlight).Bold(true)

		maxLen := m.Width - 16
		if maxLen > 60 {
			maxLen = 60
		}

		for i := scrollOffset; i < scrollOffset+maxShow && i < m.Personal.Len(); i++ {
			num := dimStyle.Render(fmt.Sprintf("%2d.", i+1))
			name := m.Personal.At(i).DisplayName()
			if len(name) > maxLen {
				name = name[:maxLen-3] + "..."
			}

			nextMark := ""
			if i == 0 {
				nextMark = lipgloss.NewStyle().Foreground(ColorSuccess).Render(" ♪ plays next")
			}

			if i == m.PersonalCursor {
				cursor := cursorStyle.Render("▶")
				label := activeStyle.Render(name)
				items = append(items, fmt.Sprintf("  %s %s %s%s", cursor, num, label, nextMark))
			} else {
				label := normalStyle.Render(name)
				items = append(items, fmt.Sprintf("    %s %s%s", num, label, nextMark))
			}
		}

		if m.Personal.Len() > scrollOffset+maxShow {
			items = append(items, dimStyle.Render(fmt.Sprintf("    ↓ %d more below", m.Personal.Len()-scrollOffset-maxShow)))
		}

		panel := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorSuccess).
			Padding(0, 1).Width(m.Width - 6).
			Render(strings.Join(items, "\n"))
		sections = append(sections, lipgloss.PlaceHorizontal(m.Width, lipgloss.Center, panel))
	}

	// Status
	if m.StatusMsg != "" {
		sections = append(sections, "")
		msg := lipgloss.NewStyle().Foreground(ColorHighlight).Render("  " + m.StatusMsg)
		sections = append(sections, msg)
	}

	sections = append(sections, "")

	k := lipgloss.NewStyle().Foreground(ColorHighlight).Bold(true)
	d := lipgloss.NewStyle().Foreground(ColorDim)
	sep := d.Render("  ·  ")
	hints := k.Render("[A]") + d.Render(" Search & Add") + sep +
		k.Render("[D]") + d.Render(" Remove") + sep +
		k.Render("[C]") + d.Render(" Clear All") + sep +
		k.Render("[Esc]") + d.Render(" Back")
	sections = append(sections, lipgloss.PlaceHorizontal(m.Width, lipgloss.Center, hints))

	content := strings.Join(sections, "\n")
	content = m.fillAndStatusBar(content)
	return content
}

// ─────────────────────────────────────────────────────────────────────────────
// Theme selection view
// ─────────────────────────────────────────────────────────────────────────────

func (m Model) viewThemeSelect() string {
	var sections []string

	// Now playing indicator
	if m.Controller != nil && m.IsPlaying {
		np := lipgloss.NewStyle().Foreground(ColorHighlight).
			Render(fmt.Sprintf("  ▶ %s", m.TrackTitle))
		sections = append(sections, np)
		sections = append(sections, lipgloss.NewStyle().Foreground(ColorBorder).
			Render("  "+strings.Repeat("─", m.Width-4)))
	}

	header := lipgloss.NewStyle().Foreground(ColorHighlight).Bold(true).
		Render("  ◆ Select Theme")
	subtitle := lipgloss.NewStyle().Foreground(ColorDim).
		Render("  Navigate ↑/↓ to preview — Enter to apply")
	sections = append(sections, "")
	sections = append(sections, header)
	sections = append(sections, subtitle)
	sections = append(sections, "")

	// Theme list with color previews
	var items []string
	for i, t := range AllThemes {
		// Color preview swatches
		swatch := lipgloss.NewStyle().Foreground(t.Highlight).Render("██") +
			lipgloss.NewStyle().Foreground(t.Accent).Render("██") +
			lipgloss.NewStyle().Foreground(t.Warm).Render("██") +
			lipgloss.NewStyle().Foreground(t.Success).Render("██") +
			lipgloss.NewStyle().Foreground(t.Error).Render("██")

		name := t.Name
		if i == m.ThemeCursor {
			cursor := lipgloss.NewStyle().Foreground(ColorHighlight).Bold(true).Render("▶")
			label := lipgloss.NewStyle().Foreground(ColorHighlight).Bold(true).Render(name)
			items = append(items, fmt.Sprintf("  %s  %s  %s", cursor, swatch, label))
		} else {
			label := lipgloss.NewStyle().Foreground(ColorFg).Render(name)
			items = append(items, fmt.Sprintf("      %s  %s", swatch, label))
		}
	}

	panel := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorder).
		Padding(1, 2).Width(m.Width - 10).
		Render(strings.Join(items, "\n"))
	sections = append(sections, lipgloss.PlaceHorizontal(m.Width, lipgloss.Center, panel))
	sections = append(sections, "")

	// Preview area — show current theme applied to sample text
	preview := lipgloss.NewStyle().Foreground(ColorHighlight).Bold(true).Render("  Highlight") + "  " +
		lipgloss.NewStyle().Foreground(ColorAccent).Render("Accent") + "  " +
		lipgloss.NewStyle().Foreground(ColorFg).Render("Text") + "  " +
		lipgloss.NewStyle().Foreground(ColorDim).Render("Dimmed") + "  " +
		lipgloss.NewStyle().Foreground(ColorSuccess).Render("●") + " " +
		lipgloss.NewStyle().Foreground(ColorError).Render("●") + " " +
		lipgloss.NewStyle().Foreground(ColorWarm).Render("●")
	sections = append(sections, lipgloss.PlaceHorizontal(m.Width, lipgloss.Center, preview))
	sections = append(sections, "")

	// Hints
	k := lipgloss.NewStyle().Foreground(ColorHighlight).Bold(true)
	d := lipgloss.NewStyle().Foreground(ColorDim)
	hints := k.Render("[↑/↓]") + d.Render(" Preview") + "  ·  " +
		k.Render("[Enter]") + d.Render(" Apply") + "  ·  " +
		k.Render("[Esc]") + d.Render(" Cancel")
	sections = append(sections, lipgloss.PlaceHorizontal(m.Width, lipgloss.Center, hints))

	content := strings.Join(sections, "\n")
	content = m.fillAndStatusBar(content)
	return content
}

// ─────────────────────────────────────────────────────────────────────────────
// Shared rendering helpers
// ─────────────────────────────────────────────────────────────────────────────

func (m Model) renderMenu(moodColor lipgloss.Color) string {
	var items []string

	activeStyle := lipgloss.NewStyle().Foreground(moodColor).Bold(true)
	normalStyle := lipgloss.NewStyle().Foreground(ColorFg)
	dimStyle := lipgloss.NewStyle().Foreground(ColorDim)
	cursorStyle := lipgloss.NewStyle().Foreground(moodColor).Bold(true)

	for i, item := range m.MenuItems {
		var line string
		if i == m.MenuCursor {
			line = fmt.Sprintf("  %s %s %s %s",
				cursorStyle.Render("❯"),
				item.Icon,
				activeStyle.Render(item.Label),
				dimStyle.Render("· "+item.Desc))
		} else {
			line = fmt.Sprintf("    %s %s", item.Icon, normalStyle.Render(item.Label))
		}
		items = append(items, line)
	}

	panel := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorder).
		Padding(1, 2).Width(60).
		Render(strings.Join(items, "\n"))

	return panel
}

func (m Model) renderHints() string {
	k := lipgloss.NewStyle().Foreground(ColorHighlight).Bold(true)
	d := lipgloss.NewStyle().Foreground(ColorDim)
	sep := d.Render("  ·  ")
	hints := k.Render("[↑/↓]") + d.Render(" Navigate") + sep +
		k.Render("[Enter]") + d.Render(" Select") + sep +
		k.Render("[T]") + d.Render(" Theme") + sep +
		k.Render("[V]") + d.Render(" Visual")
	if m.UpdateAvailable && !m.UpdateInProgress {
		hints += sep + k.Render("[U]") + d.Render(" Update")
	}
	hints += sep + k.Render("[Q]") + d.Render(" Quit")
	return hints
}

func (m Model) renderPlayingHints() string {
	k := lipgloss.NewStyle().Foreground(ColorHighlight).Bold(true)
	d := lipgloss.NewStyle().Foreground(ColorDim)
	sep := d.Render("  ·  ")

	repeatLabel := "Off"
	switch m.RepeatMode {
	case "one":
		repeatLabel = "🔂"
	case "all":
		repeatLabel = "🔁"
	}

	return k.Render("[Space]") + d.Render(" Pause") + sep +
		k.Render("[N]") + d.Render(" Next") + sep +
		k.Render("[L]") + d.Render(" Loop "+repeatLabel) + sep +
		k.Render("[P]") + d.Render(" Queue") + sep +
		k.Render("[O]") + d.Render(" My Playlist") + sep +
		k.Render("[V]") + d.Render(" Visual") + sep +
		k.Render("[Esc]") + d.Render(" Back") + sep +
		k.Render("[Q]") + d.Render(" Quit")
}

func (m Model) renderProgress() string {
	width := m.Width - 16
	if width > 80 {
		width = 80
	}
	if width < 20 {
		width = 20
	}

	moodColor := ColorHighlight
	if m.MoodLabel != "" {
		moodColor = MoodColor(m.MoodLabel)
	}

	// Calculate real progress
	var progress float64
	if m.Duration > 0 {
		progress = float64(m.Elapsed) / float64(m.Duration)
		if progress > 1.0 {
			progress = 1.0
		}
		if progress < 0 {
			progress = 0
		}
	} else {
		// Unknown duration (live stream) — show a pulsing indicator
		progress = (math.Sin(float64(m.Frame)*0.03) + 1) / 2 * 0.3
	}

	filled := int(progress * float64(width))
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}

	barFull := lipgloss.NewStyle().Foreground(moodColor).Render(strings.Repeat("━", filled))
	barEmpty := lipgloss.NewStyle().Foreground(ColorBorder).Render(strings.Repeat("─", width-filled))
	dot := lipgloss.NewStyle().Foreground(moodColor).Bold(true).Render("●")

	elapsed := formatDuration(m.Elapsed)
	total := formatDuration(m.Duration)
	if m.Duration == 0 {
		total = "LIVE"
	}
	ts := lipgloss.NewStyle().Foreground(ColorDim)

	return ts.Render(elapsed) + " " + barFull + dot + barEmpty + " " + ts.Render(total)
}

func (m Model) renderStatusBar() string {
	moodColor := ColorHighlight
	if m.MoodLabel != "" {
		moodColor = MoodColor(m.MoodLabel)
	}

	left := lipgloss.NewStyle().Foreground(ColorDim).Render("  ◈ moodwave")
	right := ""
	if m.MoodLabel != "" {
		ms := lipgloss.NewStyle().Foreground(moodColor).Bold(true)
		right = fmt.Sprintf("%s %s · %.0f%%", m.MoodEmoji, ms.Render(m.MoodLabel), m.MoodConfidence*100)
	}
	version := lipgloss.NewStyle().Foreground(ColorDim).Render("v1.0.5")

	gap := m.Width - lipgloss.Width(left) - lipgloss.Width(right) - lipgloss.Width(version) - 4
	if gap < 0 {
		gap = 0
	}

	barStyle := lipgloss.NewStyle().Background(lipgloss.Color("#111111")).Width(m.Width)
	bar := left + strings.Repeat(" ", gap) + right + "  " + version + "  "
	return barStyle.Render(bar)
}

func (m Model) fillAndStatusBar(content string) string {
	contentLines := strings.Count(content, "\n") + 1
	if m.Height > contentLines+1 {
		content += strings.Repeat("\n", m.Height-contentLines-2)
	}
	content += "\n" + m.renderStatusBar()
	return content
}

func formatDuration(d time.Duration) string {
	mins := int(d.Minutes())
	secs := int(d.Seconds()) % 60
	return fmt.Sprintf("%d:%02d", mins, secs)
}

// playbackOutcome describes *why* resolveOutcome is being invoked. Every
// entry point that can end current playback funnels through resolveOutcome
// with one of these, instead of each call site independently deciding what
// happens next. This is the single source of truth for IsResolving,
// ResolveRetried, ConsecutiveFailures, Controller lifecycle, and the
// repeat/playlist/similar-music decision tree.
type playbackOutcome int

const (
	// outcomeEnded: the track finished playing normally (duration reached).
	outcomeEnded playbackOutcome = iota
	// outcomeErrored: resolve or backend playback failed with an error.
	outcomeErrored
	// outcomeManualSkip: the user pressed Next — always advances regardless
	// of repeat-one, so the user is never stuck unable to skip a track.
	outcomeManualSkip
)

// maxConsecutiveFailures caps how many times in a row a track (or the
// resolve/playback pipeline in general) can fail before resolveOutcome
// gives up on retrying that specific track and force-advances past it,
// even under repeat-one. Without this cap a permanently broken video ID
// under repeat-one would retry forever.
const maxConsecutiveFailures = 3

// resolveOutcome is the single place that decides what happens after
// playback stops, for any reason. It owns:
//   - IsResolving / IsPlaying / Controller lifecycle (stopping the old one)
//   - ResolveRetried (one-shot retry-same-track-once on error)
//   - ConsecutiveFailures (hard cap so a broken track can't loop forever)
//   - The repeat-one / repeat-all / repeat-off / playlist / similar-music
//     priority order for picking the next track
//
// No other function in this file should mutate those fields directly —
// every code path that can end current playback (tick-detected end,
// poll-detected end, playback error, manual Next) should call this instead
// of reimplementing any part of this decision on its own. That's the
// specific class of bug this replaces: three/four call sites each holding
// a slightly different copy of "what plays next" logic.
func (m *Model) resolveOutcome(outcome playbackOutcome, resolveErr error) tea.Cmd {
	// Stop whatever was playing — every outcome starts from a clean slate.
	if m.Controller != nil {
		m.Controller.Stop()
		m.Controller = nil
	}
	m.IsPlaying = false
	m.PlayStarted = time.Time{}
	m.Elapsed = 0
	m.PausedTotal = 0
	m.PausedAt = time.Time{}

	switch outcome {
	case outcomeErrored:
		m.ConsecutiveFailures++
		// A resolve failure is frequently transient (network hiccup, rate
		// limiting) — retry the exact same track once before giving up on
		// it, unless we've already hit the hard failure cap for this
		// track/session, in which case retrying further would just hang.
		if !m.ResolveRetried && m.ConsecutiveFailures < maxConsecutiveFailures && m.PQ.Current() != nil {
			m.ResolveRetried = true
			m.IsResolving = true
			m.StatusMsg = "Retrying..."
			return ResolveAndPlayCmd(m.PQ.Current(), m.Backend)
		}
		if resolveErr != nil {
			m.StatusMsg = "Playback error: " + resolveErr.Error()
		}
		m.ResolveRetried = false
		if m.ConsecutiveFailures >= maxConsecutiveFailures {
			// This track is stuck failing — force a real advance even under
			// repeat-one so we don't retry the same broken video forever.
			m.ConsecutiveFailures = 0
			return m.advanceIgnoringRepeatOne()
		}
		return m.advanceRespectingRepeat()

	case outcomeManualSkip:
		// Next always advances, even under repeat-one, so the user is never
		// stuck unable to skip a track they don't want to hear again.
		m.ConsecutiveFailures = 0
		m.ResolveRetried = false
		return m.advanceIgnoringRepeatOne()

	default: // outcomeEnded
		m.ConsecutiveFailures = 0
		m.ResolveRetried = false
		return m.advanceRespectingRepeat()
	}
}

// advanceRespectingRepeat applies the normal repeat-mode priority order:
//  1. Repeat-one  → replay current track
//  2. Repeat-all  → wrap to next track in queue
//  3. Repeat-off  → advance within queue if there's more
//  4. Queue exhausted → search for similar music (endless listening)
func (m *Model) advanceRespectingRepeat() tea.Cmd {
	if m.RepeatMode == "one" {
		if m.PQ.Current() != nil {
			m.IsResolving = true
			return ResolveAndPlayCmd(m.PQ.Current(), m.Backend)
		}
		m.IsResolving = false
		return nil
	}
	return m.advanceIgnoringRepeatOne()
}

// advanceIgnoringRepeatOne moves to a genuinely different track regardless
// of repeat-one (used for manual Next, the failure-cap force-advance, and
// as the shared tail of advanceRespectingRepeat for repeat-all/off).
//
// PRIORITY ORDER (this is the one and only place this is decided):
//  1. Personal playlist (m.Personal) — deliberately curated by the user,
//     always consumed first, one track at a time, front to back, before
//     anything auto-generated is allowed to play.
//  2. Live queue (m.PQ) — search results / autonomous play / prior refill,
//     advanced normally (with repeat-all wrapping if applicable).
//  3. Queue exhausted AND personal playlist exhausted → similar-music
//     refill (endless listening).
func (m *Model) advanceIgnoringRepeatOne() tea.Cmd {
	// Personal playlist takes priority over everything else. It is a
	// strict FIFO — PopNext always returns tracks in the order they were
	// added, and once popped they're gone from the personal list (they
	// live on only as the currently-playing track via m.PQ, below).
	if track := m.Personal.PopNext(); track != nil {
		m.PQ.Replace([]recommender.Candidate{*track}, 0)
		m.TrackTitle = m.PQ.Current().DisplayName()
		m.TrackArtist = "Resolving..."
		m.IsResolving = true
		m.StatusMsg = fmt.Sprintf("📋 Personal playlist (%d left)...", m.Personal.Len())
		return ResolveAndPlayCmd(m.PQ.Current(), m.Backend)
	}

	if m.RepeatMode == "all" && m.PQ.AdvanceWrapping() {
		m.TrackTitle = m.PQ.Current().DisplayName()
		m.TrackArtist = "Resolving..."
		m.IsResolving = true
		return ResolveAndPlayCmd(m.PQ.Current(), m.Backend)
	}
	if m.PQ.Advance() {
		m.TrackTitle = m.PQ.Current().DisplayName()
		m.TrackArtist = "Resolving..."
		m.IsResolving = true
		return ResolveAndPlayCmd(m.PQ.Current(), m.Backend)
	}
	// Queue exhausted AND personal playlist exhausted — endless listening
	// via similar-music search. Explicitly a refill: this search's own
	// intent says to replace the queue once results land, so it can never
	// be confused with a concurrent manual browse search regardless of
	// timing.
	query := buildSimilarQuery(m.PQ.All())
	if query != "" {
		m.StatusMsg = "♾ Finding similar music..."
		m.IsResolving = false // SearchCmd re-arms IsResolving once results land
		return SearchCmd(query, SearchIntentRefillQueue)
	}
	m.IsResolving = false
	m.StatusMsg = "Queue finished"
	return tea.Tick(3*time.Second, func(t time.Time) tea.Msg { return statusClearMsg{} })
}

// extractMoodKeyword pulls a mood-related word from a title for similar search.
func extractMoodKeyword(title string) string {
	keywords := []string{"lofi", "chill", "relax", "calm", "focus", "study",
		"ambient", "synthwave", "jazz", "rock", "classical", "hip hop",
		"electronic", "beats", "vibes", "acoustic", "indie", "pop"}
	lower := strings.ToLower(title)
	for _, kw := range keywords {
		if strings.Contains(lower, kw) {
			return kw
		}
	}
	return "music"
}

// buildSimilarQuery creates a search query from recently played tracks for endless listening.
// It picks artists and mood keywords from the queue's tracks to find similar music.
func buildSimilarQuery(sources []recommender.Candidate) string {
	// Collect artists and keywords from played tracks
	artists := make(map[string]bool)
	keywords := make(map[string]bool)

	for _, c := range sources {
		if c.Track != nil {
			if c.Track.Artist != "" {
				artists[c.Track.Artist] = true
			}
			kw := extractMoodKeyword(c.Track.Title)
			if kw != "music" {
				keywords[kw] = true
			}
		}
	}

	// Build query: pick one artist + one keyword for variety
	var parts []string
	for a := range artists {
		parts = append(parts, a)
		if len(parts) >= 1 {
			break
		}
	}
	for k := range keywords {
		parts = append(parts, k)
		if len(parts) >= 2 {
			break
		}
	}

	if len(parts) == 0 {
		// Fallback: use last track title words
		if len(sources) > 0 && sources[len(sources)-1].Track != nil {
			return sources[len(sources)-1].Track.Artist + " mix"
		}
		return "chill music mix"
	}

	return strings.Join(parts, " ")
}

// buildMoodQuery creates a YouTube search query from the detected mood label.
// Uses semantic understanding of each mood to find appropriate music.
func buildMoodQuery(mood string) string {
	queries := map[string]string{
		"focused":      "deep focus coding music instrumental",
		"calm":         "calm ambient study music lo-fi",
		"intense":      "intense electronic programming music",
		"chaotic":      "experimental electronic glitch music",
		"experimental": "experimental ambient electronic music",
		"minimal":      "minimal ambient piano focus",
		"polished":     "smooth jazz coding playlist",
		"late-night":   "late night coding lofi dark ambient",
		"sprint":       "high energy electronic workout coding",
		"debugging":    "chill beats problem solving music",
	}
	if q, ok := queries[mood]; ok {
		return q
	}
	return "coding music focus playlist"
}

// ─────────────────────────────────────────────────────────────────────────────
// Public API — Run launches the TUI and blocks until the user quits.
// The onAction callback is called when the user triggers an action.
// ─────────────────────────────────────────────────────────────────────────────

// Run launches the persistent TUI. Returns the action to perform after exit.
func Run(moodLabel, moodEmoji string, confidence float64, onAction func(action, args string), scanFn ScanFunc) error {
	m := NewModel()
	m.MoodLabel = moodLabel
	m.MoodEmoji = moodEmoji
	m.MoodConfidence = confidence
	m.OnAction = onAction
	m.ScanFn = scanFn

	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
