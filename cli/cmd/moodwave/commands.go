// commands.go implements all moodwave CLI subcommands.
package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"golang.org/x/term"

	"github.com/moodwave/moodwave/internal/cache"
	"github.com/moodwave/moodwave/internal/config"
	"github.com/moodwave/moodwave/internal/mood"
	"github.com/moodwave/moodwave/internal/platform"
	"github.com/moodwave/moodwave/internal/playback"
	"github.com/moodwave/moodwave/internal/recommender"
	"github.com/moodwave/moodwave/internal/scanner"
	"github.com/moodwave/moodwave/internal/sources"
	"github.com/moodwave/moodwave/internal/tui"
	"github.com/moodwave/moodwave/internal/updater"
	"github.com/moodwave/moodwave/internal/visuals"
)

// ──────────────────────────────────────────────────────────────────────────────
// SESSION STATE — persisted between CLI invocations
// ──────────────────────────────────────────────────────────────────────────────

// Session holds state persisted across CLI invocations.
type Session struct {
	MoodProfile      *mood.Profile           `json:"mood_profile"`
	ScannedAt        time.Time               `json:"scanned_at"`
	CurrentCandidate *recommender.Candidate  `json:"current_candidate"`
	Queue            []recommender.Candidate `json:"queue"`
	QueueIndex       int                     `json:"queue_index"`
	PlaybackState    string                  `json:"playback_state"`
	RecentlyPlayed   []string                `json:"recently_played"`
	RepeatMode       string                  `json:"repeat_mode"` // "off", "one", "all"
}

// loadSession reads the session state file.
func (a *App) loadSession() (*Session, error) {
	data, err := os.ReadFile(a.sessionFile())
	if err != nil {
		if os.IsNotExist(err) {
			return &Session{}, nil
		}
		return nil, err
	}
	var s Session
	if err := json.Unmarshal(data, &s); err != nil {
		return &Session{}, nil // treat corrupt session as empty
	}
	return &s, nil
}

// saveSession writes the session state file.
func (a *App) saveSession(s *Session) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(a.sessionFile(), data, 0644)
}

// ──────────────────────────────────────────────────────────────────────────────
// init — Initialize config and directories
// ──────────────────────────────────────────────────────────────────────────────

func (a *App) cmdInit() error {
	fmt.Println("Initializing Moodwave...")

	// Create directories.
	if err := config.EnsureDirectories(a.cfg); err != nil {
		return fmt.Errorf("creating directories: %w", err)
	}

	// Write default config if it doesn't exist.
	if _, err := os.Stat(a.cfg.Paths.ConfigFile); os.IsNotExist(err) {
		if err := config.WriteDefaults(a.cfg); err != nil {
			return fmt.Errorf("writing config: %w", err)
		}
		fmt.Printf("  ✓ Config written to %s\n", a.cfg.Paths.ConfigFile)
	} else {
		fmt.Printf("  ✓ Config already exists at %s\n", a.cfg.Paths.ConfigFile)
	}

	fmt.Printf("  ✓ Cache directory: %s\n", a.cfg.Paths.CacheDir)
	fmt.Println("\nMoodwave is ready. Run 'moodwave scan' to detect your coding mood.")
	return nil
}

// ──────────────────────────────────────────────────────────────────────────────
// scan — Scan the repository
// ──────────────────────────────────────────────────────────────────────────────

func (a *App) cmdScan(args []string) error {
	root := a.projectRoot
	if len(args) > 0 {
		root = args[0]
	}

	absRoot, err := filepath.Abs(root)
	if err != nil {
		return err
	}

	fmt.Printf("Scanning %s...\n\n", absRoot)

	// Run scanner.
	opts := scanner.ScanOptions{
		MaxDepth:   a.cfg.Scanner.MaxDepth,
		MaxFiles:   a.cfg.Scanner.MaxFiles,
		GitEnabled: a.cfg.Scanner.GitEnabled,
		IgnoreDirs: a.cfg.Scanner.IgnoreDirs,
	}
	s := scanner.New(opts)
	signals, err := s.Scan(absRoot)
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	// Print scan summary.
	fmt.Printf("  Files scanned:    %d\n", signals.TotalFiles)
	fmt.Printf("  Source files:     %d\n", signals.SourceFiles)
	fmt.Printf("  Primary language: %s\n", signals.PrimaryLanguage)
	if signals.Git.RepoDetected {
		fmt.Printf("  Branch:           %s\n", signals.Git.CurrentBranch)
		fmt.Printf("  Churn:            %.0f%%\n", signals.Git.ChurnScore*100)
	}
	if signals.HasBuildSystem {
		fmt.Printf("  Build system:     %s\n", signals.BuildSystemType)
	}
	if len(signals.Languages) > 0 {
		langs := make([]string, 0, min3(len(signals.Languages), 5))
		for _, l := range signals.Languages {
			if len(langs) >= 5 {
				break
			}
			langs = append(langs, fmt.Sprintf("%s (%.0f%%)", l.Name, l.Percentage))
		}
		fmt.Printf("  Languages:        %s\n", strings.Join(langs, ", "))
	}
	if len(signals.SemanticMoodCounts) > 0 {
		type match struct {
			label string
			count int
		}
		var list []match
		for m, c := range signals.SemanticMoodCounts {
			if c > 0 {
				list = append(list, match{m, c})
			}
		}
		sort.Slice(list, func(i, j int) bool {
			return list[i].count > list[j].count
		})
		var semStr []string
		for _, item := range list {
			semStr = append(semStr, fmt.Sprintf("%s (%d)", item.label, item.count))
		}
		if len(semStr) > 0 {
			fmt.Printf("  Semantic Vocab:   %s\n", strings.Join(semStr, ", "))
		}
	}

	fmt.Println()

	// Infer mood.
	engine := mood.NewEngine(a.cfg.Mood.Sensitivity)
	profile := engine.Infer(signals)

	fmt.Printf("  Detected mood:    %s %s (%.0f%% confidence)\n",
		profile.Label.Emoji(), profile.Label, profile.Confidence*100)
	fmt.Printf("  Explanation:      %s\n", profile.Explanation)
	fmt.Printf("  BPM range:        %d–%d\n", profile.Traits.BPMMin, profile.Traits.BPMMax)
	fmt.Printf("  Music tags:       %s\n", strings.Join(profile.Traits.Tags[:min3(len(profile.Traits.Tags), 5)], ", "))

	// Save session.
	session := &Session{
		MoodProfile: profile,
		ScannedAt:   time.Now(),
	}
	if err := a.saveSession(session); err != nil {
		a.debugf("save session: %v", err)
	}

	fmt.Println("\nRun 'moodwave play' to start music matched to this mood.")
	return nil
}

// ──────────────────────────────────────────────────────────────────────────────
// mood — Show current mood
// ──────────────────────────────────────────────────────────────────────────────

func (a *App) cmdMood() error {
	session, err := a.loadSession()
	if err != nil {
		return err
	}

	if session.MoodProfile == nil {
		fmt.Println("No mood detected yet. Run 'moodwave scan' first.")
		return nil
	}

	m := session.MoodProfile
	age := time.Since(session.ScannedAt).Truncate(time.Minute)

	fmt.Printf("\nMood: %s %s\n", m.Label.Emoji(), strings.ToUpper(string(m.Label)))
	fmt.Printf("Confidence: %.0f%%\n", m.Confidence*100)
	fmt.Printf("Scanned: %v ago\n", age)
	fmt.Printf("\n%s\n", m.Explanation)

	fmt.Printf("\nMusic traits:\n")
	fmt.Printf("  BPM range:    %d–%d\n", m.Traits.BPMMin, m.Traits.BPMMax)
	fmt.Printf("  Energy:       %.0f%%–%.0f%%\n", m.Traits.EnergyMin*100, m.Traits.EnergyMax*100)
	fmt.Printf("  Tags:         %s\n", strings.Join(m.Traits.Tags, ", "))

	if len(m.Signals) > 0 {
		fmt.Printf("\nTop signals:\n")
		for _, sig := range m.Signals {
			fmt.Printf("  · %s (%s)\n", sig.Signal, sig.Effect)
		}
	}

	fmt.Printf("\nAll mood scores:\n")
	for _, label := range mood.AllMoods {
		score := m.Scores[label]
		bar := strings.Repeat("█", int(score*20))
		marker := " "
		if label == m.Label {
			marker = ">"
		}
		fmt.Printf("  %s %-14s %s %.0f%%\n", marker, label, bar, score*100)
	}

	return nil
}

// ──────────────────────────────────────────────────────────────────────────────
// play — Start playback
// ──────────────────────────────────────────────────────────────────────────────

func (a *App) cmdPlay(args []string) error {
	session, err := a.loadSession()
	if err != nil {
		return err
	}

	// Put terminal in raw mode for the entire play command duration
	var oldState *term.State
	if a.caps.IsTTY {
		state, err := term.MakeRaw(int(os.Stdin.Fd()))
		if err == nil {
			oldState = state
			defer term.Restore(int(os.Stdin.Fd()), oldState)
		}
	}

	keyChan := make(chan rune, 100)
	keyDone := make(chan struct{})
	defer close(keyDone)

	if a.caps.IsTTY {
		go func() {
			buf := make([]byte, 1)
			for {
				select {
				case <-keyDone:
					return
				default:
					_, readErr := os.Stdin.Read(buf)
					if readErr != nil {
						return
					}
					select {
					case <-keyDone:
						return
					default:
						keyChan <- rune(buf[0])
					}
				}
			}
		}()
	}

	// Build source registry.
	registry := a.buildSourceRegistry()

	ctx, cancel := context.WithTimeout(a.ctx, 60*time.Second)
	defer cancel()

	var candidates []recommender.Candidate

	if len(args) > 0 {
		queryText := strings.Join(args, " ")
		fmt.Printf("\nSearching for: %q\n\n", queryText)

		query := sources.SearchQuery{
			Text:  queryText,
			Limit: 15,
		}

		for _, adapter := range registry.All() {
			// Skip check if it fails, but try to search if possible.
			_ = adapter.HealthCheck(ctx)

			// Search tracks (primary for direct searches, especially Jamendo).
			tracks, err := adapter.SearchTracks(ctx, query)
			if err == nil {
				for i := range tracks {
					if tracks[i].StreamURL == "" {
						continue
					}
					candidates = append(candidates, recommender.Candidate{
						Track:  &tracks[i],
						Score:  1.0,
						Reason: fmt.Sprintf("Direct match from %s", adapter.Name()),
					})
				}
			}

			// Search stations.
			stations, err := adapter.SearchStations(ctx, query)
			if err == nil {
				for i := range stations {
					if stations[i].StreamURL == "" {
						continue
					}
					candidates = append(candidates, recommender.Candidate{
						Station: &stations[i],
						Score:   0.9,
						Reason:  fmt.Sprintf("Station match from %s", adapter.Name()),
					})
				}
			}
		}

	} else {
		// If no mood profile, scan first.
		if session.MoodProfile == nil {
			fmt.Println("No mood profile found. Scanning first...")
			if err := a.cmdScan(nil); err != nil {
				return err
			}
			session, _ = a.loadSession()
		}

		fmt.Printf("\nFinding music for mood: %s %s\n\n",
			session.MoodProfile.Label.Emoji(), session.MoodProfile.Label)

		// Get recommendations.
		recCfg := recommender.DefaultConfig()
		recCfg.RecentlyPlayed = session.RecentlyPlayed
		rec := recommender.New(registry, recCfg)

		var curTrack *sources.Track
		if session.CurrentCandidate != nil {
			curTrack = session.CurrentCandidate.Track
		}
		candidates, err = rec.Recommend(ctx, session.MoodProfile, curTrack)
		if err != nil {
			return fmt.Errorf("recommendation failed: %w", err)
		}
	}

	if len(candidates) == 0 {
		fmt.Println("No tracks found. Check your internet connection or run 'moodwave doctor'.")
		return nil
	}

	// Save queue to session.
	session.Queue = candidates
	session.QueueIndex = 0

	for {
		if session.QueueIndex < 0 || session.QueueIndex >= len(session.Queue) {
			break
		}

		current := &session.Queue[session.QueueIndex]
		session.CurrentCandidate = current

		streamURL := current.StreamURL()
		if current.Track != nil && current.Track.Source == "youtube" {
			if a.caps.IsTTY {
				fmt.Printf("\r\033[K  Resolving YouTube audio stream for %s...\r", current.DisplayName())
			} else {
				fmt.Println("  Resolving YouTube audio stream...")
			}
			yt := sources.NewYouTubeAdapter()
			resolveCtx, resolveCancel := context.WithTimeout(a.ctx, 30*time.Second)
			resolved, err := yt.ResolveTrack(resolveCtx, current.Track.ID)
			resolveCancel()
			if err == nil {
				streamURL = resolved.StreamURL
				current.Track.StreamURL = resolved.StreamURL
			} else {
				// Resolution failed — skip this track (YouTube watch URLs can't be played directly)
				streamURL = ""
			}
		}

		title, artist := "", ""
		if current.Track != nil {
			title = current.Track.Title
			artist = current.Track.Artist
		} else if current.Station != nil {
			title = current.Station.Name
		}

		// Skip if no stream URL resolved
		if streamURL == "" {
			session.QueueIndex++
			_ = a.saveSession(session)
			continue
		}

		// Start audio playback.
		controller, err := playback.NewController(a.cfg.Playback.Backend)
		if err != nil {
			// Auto-install ffplay if no backend found
			fmt.Printf("  No audio backend — installing ffplay automatically...\n")
			installCtx, installCancel := context.WithTimeout(a.ctx, 5*time.Minute)
			_, installErr := playback.InstallFFplay(installCtx, func(msg string) {
				fmt.Printf("  %s\n", msg)
			})
			installCancel()
			if installErr != nil {
				fmt.Printf("  Auto-install failed: %v\n", installErr)
				fmt.Printf("  Install mpv or ffmpeg manually, or run 'moodwave doctor' to fix.\n")
				fmt.Printf("  Stream URL: %s\n", streamURL)
				session.PlaybackState = "info-only"
				_ = a.saveSession(session)
				return nil
			}
			// Retry with newly installed ffplay
			controller, err = playback.NewController(a.cfg.Playback.Backend)
			if err != nil {
				fmt.Printf("  Still no audio backend after install: %v\n", err)
				return nil
			}
		}

		if err := controller.Play(a.ctx, streamURL, title, artist); err != nil {
			fmt.Printf("\n  Playback failed for track %s: %v\n", current.DisplayName(), err)
			session.QueueIndex = (session.QueueIndex + 1) % len(session.Queue)
			_ = a.saveSession(session)
			continue
		}

		session.PlaybackState = "playing"
		session.RecentlyPlayed = append(session.RecentlyPlayed, streamURL)
		if len(session.RecentlyPlayed) > 20 {
			session.RecentlyPlayed = session.RecentlyPlayed[len(session.RecentlyPlayed)-20:]
		}
		_ = a.saveSession(session)

		var action string
		if a.caps.IsTTY {
			action = a.startPlaybackRenderer(controller, session, keyChan)
		} else {
			fmt.Printf("  Playing: %s (Backend: %s)\n", current.DisplayName(), controller.BackendName())
			fmt.Println("Playing... Press Ctrl-C to stop.")
			<-a.ctx.Done()
			action = "quit"
		}

		controller.Stop()

		if action == "quit" {
			break
		}

		if strings.HasPrefix(action, "play ") {
			query := strings.TrimPrefix(action, "play ")

			var newCandidates []recommender.Candidate
			if strings.HasPrefix(query, "http://") || strings.HasPrefix(query, "https://") {
				newCandidates = []recommender.Candidate{
					{
						Track: &sources.Track{
							ID:        "raw",
							Source:    "custom",
							Title:     query,
							Artist:    "Stream URL",
							StreamURL: query,
						},
						Score:  1,
						Reason: "Direct URL play",
					},
				}
			} else {
				yt := sources.NewYouTubeAdapter()
				searchCtx, searchCancel := context.WithTimeout(a.ctx, 15*time.Second)
				results, err := yt.SearchTracks(searchCtx, sources.SearchQuery{Text: query})
				searchCancel()
				if err == nil && len(results) > 0 {
					for i := range results {
						newCandidates = append(newCandidates, recommender.Candidate{
							Track:  &results[i],
							Score:  1,
							Reason: "YouTube search",
						})
					}
				}
			}

			if len(newCandidates) > 0 {
				session.Queue = newCandidates
				session.QueueIndex = 0
				_ = a.saveSession(session)
				continue
			} else {
				continue
			}
		}

		// Move to the next index in the queue based on repeat mode.
		if session.RepeatMode == "" {
			session.RepeatMode = "off"
		}

		if session.RepeatMode == "one" {
			// Do not change QueueIndex, repeat the same song!
		} else if session.RepeatMode == "all" {
			session.QueueIndex = (session.QueueIndex + 1) % len(session.Queue)
		} else { // "off"
			if session.QueueIndex+1 >= len(session.Queue) {
				// Reached end of queue, stop playback
				break
			}
			session.QueueIndex = session.QueueIndex + 1
		}
		_ = a.saveSession(session)
	}

	return nil
}

// search — Search YouTube directly and play selection
func (a *App) cmdSearch(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("please provide a search query: moodwave search <query>")
	}

	queryText := strings.Join(args, " ")
	fmt.Printf("\nSearching YouTube for: %q...\n\n", queryText)

	yt := sources.NewYouTubeAdapter()
	ctx, cancel := context.WithTimeout(a.ctx, 15*time.Second)
	results, err := yt.SearchTracks(ctx, sources.SearchQuery{Text: queryText})
	cancel()

	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	if len(results) == 0 {
		fmt.Println("No tracks found on YouTube.")
		return nil
	}

	displayCount := len(results)
	if displayCount > 5 {
		displayCount = 5
	}

	fmt.Println("Search Results:")
	for i := 0; i < displayCount; i++ {
		fmt.Printf("  %d. %s - %s [YouTube]\n", i+1, results[i].Artist, results[i].Title)
	}
	fmt.Println()

	var choice int
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("Select a song number to play (1-%d), or 'q' to quit: ", displayCount)
		input, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		input = strings.TrimSpace(input)
		if input == "q" || input == "Q" {
			return nil
		}

		var val int
		_, scanErr := fmt.Sscanf(input, "%d", &val)
		if scanErr == nil && val >= 1 && val <= displayCount {
			choice = val - 1
			break
		}
		fmt.Println("Invalid selection. Please try again.")
	}

	selectedTrack := &results[choice]
	fmt.Printf("\nResolving audio stream for: %s...\n", selectedTrack.Title)

	resolveCtx, resolveCancel := context.WithTimeout(a.ctx, 30*time.Second)
	resolved, err := yt.ResolveTrack(resolveCtx, selectedTrack.ID)
	resolveCancel()
	if err != nil {
		return fmt.Errorf("failed to resolve audio stream: %w", err)
	}

	selectedTrack.StreamURL = resolved.StreamURL

	session, err := a.loadSession()
	if err != nil {
		session = &Session{}
	}

	candidate := recommender.Candidate{
		Track:  selectedTrack,
		Score:  1.0,
		Reason: "Direct YouTube Search selection",
	}

	session.Queue = []recommender.Candidate{candidate}
	session.QueueIndex = 0
	session.CurrentCandidate = &candidate
	_ = a.saveSession(session)

	controller, err := playback.NewController(a.cfg.Playback.Backend)
	if err != nil {
		// Auto-install ffplay
		fmt.Printf("  No audio backend — installing ffplay...\n")
		installCtx, installCancel := context.WithTimeout(a.ctx, 5*time.Minute)
		_, installErr := playback.InstallFFplay(installCtx, func(msg string) {
			fmt.Printf("  %s\n", msg)
		})
		installCancel()
		if installErr != nil {
			return fmt.Errorf("no audio backend and auto-install failed: %w\nRun 'moodwave doctor' to fix", installErr)
		}
		controller, err = playback.NewController(a.cfg.Playback.Backend)
		if err != nil {
			return fmt.Errorf("no audio backend found after install: %w", err)
		}
	}
	defer controller.Stop()

	if err := controller.Play(a.ctx, selectedTrack.StreamURL, selectedTrack.Title, selectedTrack.Artist); err != nil {
		return fmt.Errorf("playback failed: %w", err)
	}

	session.PlaybackState = "playing"
	_ = a.saveSession(session)

	var action string
	if a.caps.IsTTY {
		var oldState *term.State
		state, err := term.MakeRaw(int(os.Stdin.Fd()))
		if err == nil {
			oldState = state
			defer term.Restore(int(os.Stdin.Fd()), oldState)
		}

		keyChan := make(chan rune, 100)
		keyDone := make(chan struct{})
		defer close(keyDone)

		go func() {
			buf := make([]byte, 1)
			for {
				select {
				case <-keyDone:
					return
				default:
					_, readErr := os.Stdin.Read(buf)
					if readErr != nil {
						return
					}
					select {
					case <-keyDone:
						return
					default:
						keyChan <- rune(buf[0])
					}
				}
			}
		}()

		action = a.startPlaybackRenderer(controller, session, keyChan)
	} else {
		fmt.Printf("  Playing: %s\n", selectedTrack.Title)
		<-a.ctx.Done()
		action = "quit"
	}

	if action == "next" {
		fmt.Println("\nPlayback finished.")
	}

	return nil
}

// startPlaybackRenderer starts the interactive terminal UI during playback.
func (a *App) startPlaybackRenderer(ctrl *playback.Controller, session *Session, keyChan chan rune) string {
	rendCfg := visuals.RendererConfig{
		VisualMode:  a.cfg.Visual.Mode,
		Theme:       a.cfg.Visual.Theme,
		FPS:         a.cfg.Visual.FPS,
		NoAnimation: a.cfg.Visual.NoAnimation,
		NoColor:     a.cfg.Visual.NoColor,
		NoUnicode:   a.cfg.Visual.NoUnicode,
		Caps:        a.caps,
	}

	renderer := visuals.New(rendCfg)
	renderer.SetState(visuals.RenderState{
		Scene:      visuals.ScenePlaying,
		Mood:       session.MoodProfile,
		Candidate:  session.CurrentCandidate,
		RepeatMode: session.RepeatMode,
	})
	renderer.Start()
	defer renderer.Stop()

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	inInputMode := false
	inputBuf := ""
	var statusMsgClearTime time.Time

	for {
		select {
		case <-a.ctx.Done():
			return "quit"
		case key := <-keyChan:
			if inInputMode {
				switch key {
				case 27: // Esc
					inInputMode = false
					inputBuf = ""
					status := ctrl.Status()
					scene := visuals.ScenePlaying
					if status.State == playback.StatePaused {
						scene = visuals.ScenePaused
					}
					renderer.SetState(visuals.RenderState{
						Scene:     scene,
						Mood:      session.MoodProfile,
						Candidate: session.CurrentCandidate,
						InputMode: false,
						InputBuf:  "",
					})
				case 127, 8: // Backspace (DEL or Backspace)
					if len(inputBuf) > 0 {
						inputBuf = inputBuf[:len(inputBuf)-1]
					}
					status := ctrl.Status()
					scene := visuals.ScenePlaying
					if status.State == playback.StatePaused {
						scene = visuals.ScenePaused
					}
					renderer.SetState(visuals.RenderState{
						Scene:     scene,
						Mood:      session.MoodProfile,
						Candidate: session.CurrentCandidate,
						InputMode: true,
						InputBuf:  inputBuf,
					})
				case 13, 10: // Enter
					inInputMode = false
					cmd := strings.TrimSpace(inputBuf)
					inputBuf = ""
					if cmd != "" {
						if cmd == "q" || cmd == "quit" || cmd == "exit" {
							return "quit"
						} else if cmd == "n" || cmd == "next" {
							renderer.SetState(visuals.RenderState{
								Scene:     visuals.SceneSwitching,
								Mood:      session.MoodProfile,
								Candidate: session.CurrentCandidate,
							})
							return "next"
						} else if cmd == "pause" || cmd == "stop" {
							ctrl.Pause()
							renderer.SetState(visuals.RenderState{
								Scene:     visuals.ScenePaused,
								Mood:      session.MoodProfile,
								Candidate: session.CurrentCandidate,
							})
						} else if cmd == "resume" || cmd == "play" {
							ctrl.Resume()
							renderer.SetState(visuals.RenderState{
								Scene:     visuals.ScenePlaying,
								Mood:      session.MoodProfile,
								Candidate: session.CurrentCandidate,
							})
						} else if strings.HasPrefix(cmd, "visuals ") {
							modeStr := strings.TrimPrefix(cmd, "visuals ")
							var nextMode config.VisualMode
							switch modeStr {
							case "wave":
								nextMode = config.VisualWave
							case "spectrum":
								nextMode = config.VisualSpectrum
							case "pulse":
								nextMode = config.VisualPulse
							case "fireplace", "fire":
								nextMode = config.VisualFire
							case "minimal":
								nextMode = config.VisualMinimal
							case "quiet":
								nextMode = config.VisualQuiet
							default:
								nextMode = a.cfg.Visual.Mode
							}
							a.cfg.Visual.Mode = nextMode
							renderer.SetVisualMode(nextMode)
						} else if cmd == "loop" || cmd == "repeat" || strings.HasPrefix(cmd, "loop ") || strings.HasPrefix(cmd, "repeat ") {
							mode := ""
							if strings.HasPrefix(cmd, "loop ") {
								mode = strings.TrimPrefix(cmd, "loop ")
							} else if strings.HasPrefix(cmd, "repeat ") {
								mode = strings.TrimPrefix(cmd, "repeat ")
							}

							nextRepeat := ""
							switch strings.ToLower(strings.TrimSpace(mode)) {
							case "one", "track", "1":
								nextRepeat = "one"
							case "all", "queue":
								nextRepeat = "all"
							case "off", "none":
								nextRepeat = "off"
							default:
								switch session.RepeatMode {
								case "off", "":
									nextRepeat = "one"
								case "one":
									nextRepeat = "all"
								case "all":
									nextRepeat = "off"
								}
							}
							session.RepeatMode = nextRepeat
							_ = a.saveSession(session)

							status := ctrl.Status()
							scene := visuals.ScenePlaying
							if status.State == playback.StatePaused {
								scene = visuals.ScenePaused
							}
							renderer.SetState(visuals.RenderState{
								Scene:      scene,
								Mood:       session.MoodProfile,
								Candidate:  session.CurrentCandidate,
								RepeatMode: session.RepeatMode,
								ScanMsg:    fmt.Sprintf("Repeat mode set to %s", strings.ToUpper(nextRepeat)),
							})
							statusMsgClearTime = time.Now().Add(2 * time.Second)
						} else if strings.HasPrefix(cmd, "add ") || strings.HasPrefix(cmd, "queue ") || strings.HasPrefix(cmd, "playnext ") {
							isPlayNext := strings.HasPrefix(cmd, "playnext ")
							query := ""
							if isPlayNext {
								query = strings.TrimPrefix(cmd, "playnext ")
							} else if strings.HasPrefix(cmd, "add ") {
								query = strings.TrimPrefix(cmd, "add ")
							} else {
								query = strings.TrimPrefix(cmd, "queue ")
							}
							query = strings.TrimSpace(query)

							renderer.SetState(visuals.RenderState{
								Scene:     visuals.SceneScanning,
								Mood:      session.MoodProfile,
								Candidate: session.CurrentCandidate,
								ScanMsg:   "Searching & queuing track...",
							})

							yt := sources.NewYouTubeAdapter()
							searchCtx, searchCancel := context.WithTimeout(a.ctx, 15*time.Second)
							results, err := yt.SearchTracks(searchCtx, sources.SearchQuery{Text: query})
							searchCancel()

							if err == nil && len(results) > 0 {
								newCandidate := recommender.Candidate{
									Track:  &results[0],
									Score:  1.0,
									Reason: "User queued track",
								}

								if isPlayNext {
									insertIdx := session.QueueIndex + 1
									if insertIdx >= len(session.Queue) {
										session.Queue = append(session.Queue, newCandidate)
									} else {
										session.Queue = append(session.Queue[:insertIdx+1], session.Queue[insertIdx:]...)
										session.Queue[insertIdx] = newCandidate
									}
									renderer.SetState(visuals.RenderState{
										Scene:     visuals.SceneScanning,
										Mood:      session.MoodProfile,
										Candidate: session.CurrentCandidate,
										ScanMsg:   fmt.Sprintf("Playing next: %s", results[0].Title),
									})
								} else {
									session.Queue = append(session.Queue, newCandidate)
									renderer.SetState(visuals.RenderState{
										Scene:     visuals.SceneScanning,
										Mood:      session.MoodProfile,
										Candidate: session.CurrentCandidate,
										ScanMsg:   fmt.Sprintf("Queued to end: %s", results[0].Title),
									})
								}
								_ = a.saveSession(session)
								statusMsgClearTime = time.Now().Add(2 * time.Second)
							} else {
								renderer.SetState(visuals.RenderState{
									Scene:     visuals.SceneScanning,
									Mood:      session.MoodProfile,
									Candidate: session.CurrentCandidate,
									ScanMsg:   "No tracks found to queue.",
								})
								statusMsgClearTime = time.Now().Add(2 * time.Second)
							}
						} else {
							// Treat anything else as a play query!
							query := cmd
							if strings.HasPrefix(cmd, "play ") {
								query = strings.TrimPrefix(cmd, "play ")
							}
							renderer.SetState(visuals.RenderState{
								Scene:     visuals.SceneSwitching,
								Mood:      session.MoodProfile,
								Candidate: session.CurrentCandidate,
							})
							return "play " + query
						}
					}
					status := ctrl.Status()
					scene := visuals.ScenePlaying
					if status.State == playback.StatePaused {
						scene = visuals.ScenePaused
					}
					renderer.SetState(visuals.RenderState{
						Scene:     scene,
						Mood:      session.MoodProfile,
						Candidate: session.CurrentCandidate,
						InputMode: false,
						InputBuf:  "",
					})
				default:
					if key >= 32 && key <= 126 {
						inputBuf += string(key)
						status := ctrl.Status()
						scene := visuals.ScenePlaying
						if status.State == playback.StatePaused {
							scene = visuals.ScenePaused
						}
						renderer.SetState(visuals.RenderState{
							Scene:     scene,
							Mood:      session.MoodProfile,
							Candidate: session.CurrentCandidate,
							InputMode: true,
							InputBuf:  inputBuf,
						})
					}
				}
			} else {
				switch key {
				case 'q', 'Q', 3: // 'q' or Ctrl-C
					return "quit"
				case ' ': // Space to play/pause
					status := ctrl.Status()
					if status.State == playback.StatePlaying {
						ctrl.Pause()
						renderer.SetState(visuals.RenderState{
							Scene:     visuals.ScenePaused,
							Mood:      session.MoodProfile,
							Candidate: session.CurrentCandidate,
						})
					} else if status.State == playback.StatePaused {
						ctrl.Resume()
						renderer.SetState(visuals.RenderState{
							Scene:     visuals.ScenePlaying,
							Mood:      session.MoodProfile,
							Candidate: session.CurrentCandidate,
						})
					}
				case 'n', 'N': // Next
					renderer.SetState(visuals.RenderState{
						Scene:     visuals.SceneSwitching,
						Mood:      session.MoodProfile,
						Candidate: session.CurrentCandidate,
					})
					return "next"
				case 'v', 'V': // Cycle visuals
					modes := []config.VisualMode{config.VisualWave, config.VisualSpectrum, config.VisualPulse, config.VisualFire, config.VisualMinimal, config.VisualQuiet}
					nextMode := config.VisualWave
					for idx, m := range modes {
						if m == a.cfg.Visual.Mode {
							nextMode = modes[(idx+1)%len(modes)]
							break
						}
					}
					a.cfg.Visual.Mode = nextMode
					renderer.SetVisualMode(nextMode)
				case 's', 'S': // Show temporary scan scanning UI feedback
					renderer.SetState(visuals.RenderState{
						Scene:     visuals.SceneScanning,
						Mood:      session.MoodProfile,
						Candidate: session.CurrentCandidate,
						ScanMsg:   "Initiating scan...",
					})
					time.Sleep(500 * time.Millisecond)
					status := ctrl.Status()
					scene := visuals.ScenePlaying
					if status.State == playback.StatePaused {
						scene = visuals.ScenePaused
					}
					renderer.SetState(visuals.RenderState{
						Scene:     scene,
						Mood:      session.MoodProfile,
						Candidate: session.CurrentCandidate,
					})
				case 'l', 'L': // Toggle repeat/loop mode
					nextRepeat := "off"
					switch session.RepeatMode {
					case "off", "":
						nextRepeat = "one"
					case "one":
						nextRepeat = "all"
					case "all":
						nextRepeat = "off"
					}
					session.RepeatMode = nextRepeat
					_ = a.saveSession(session)

					status := ctrl.Status()
					scene := visuals.ScenePlaying
					if status.State == playback.StatePaused {
						scene = visuals.ScenePaused
					}
					renderer.SetState(visuals.RenderState{
						Scene:      scene,
						Mood:       session.MoodProfile,
						Candidate:  session.CurrentCandidate,
						RepeatMode: session.RepeatMode,
						ScanMsg:    fmt.Sprintf("Repeat mode: %s", strings.ToUpper(nextRepeat)),
					})
					statusMsgClearTime = time.Now().Add(2 * time.Second)
				case 13, 10: // Enter key in normal mode starts input mode!
					inInputMode = true
					status := ctrl.Status()
					scene := visuals.ScenePlaying
					if status.State == playback.StatePaused {
						scene = visuals.ScenePaused
					}
					renderer.SetState(visuals.RenderState{
						Scene:     scene,
						Mood:      session.MoodProfile,
						Candidate: session.CurrentCandidate,
						InputMode: true,
						InputBuf:  "",
					})
				}
			}
		case <-ticker.C:
			status := ctrl.Status()
			if status.State == playback.StateError || status.State == playback.StateStopped {
				return "done"
			}
			if !statusMsgClearTime.IsZero() && time.Now().After(statusMsgClearTime) {
				statusMsgClearTime = time.Time{}
				scene := visuals.ScenePlaying
				if status.State == playback.StatePaused {
					scene = visuals.ScenePaused
				}
				renderer.SetState(visuals.RenderState{
					Scene:     scene,
					Mood:      session.MoodProfile,
					Candidate: session.CurrentCandidate,
					InputMode: false,
					InputBuf:  "",
					ScanMsg:   "",
				})
			}
		}
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// pause / stop / next
// ──────────────────────────────────────────────────────────────────────────────

func (a *App) cmdPause() error {
	session, err := a.loadSession()
	if err != nil {
		return err
	}
	if session.PlaybackState != "playing" {
		fmt.Println("Nothing is currently playing.")
		return nil
	}
	fmt.Println("Pause: stop the current 'moodwave play' process (Ctrl-C) to pause.")
	fmt.Println("(Interactive pause requires the play session to be running in this terminal.)")
	return nil
}

func (a *App) cmdStop() error {
	session, err := a.loadSession()
	if err != nil {
		return err
	}
	session.PlaybackState = "stopped"
	_ = a.saveSession(session)
	fmt.Println("Stopped. Run 'moodwave play' to resume.")
	return nil
}

func (a *App) cmdNext() error {
	session, err := a.loadSession()
	if err != nil {
		return err
	}

	if len(session.Queue) == 0 {
		fmt.Println("No queue. Run 'moodwave play' first.")
		return nil
	}

	session.QueueIndex = (session.QueueIndex + 1) % len(session.Queue)
	session.CurrentCandidate = &session.Queue[session.QueueIndex]
	_ = a.saveSession(session)

	fmt.Printf("Next: %s\n", session.CurrentCandidate.DisplayName())
	fmt.Println("Run 'moodwave play' to start this track.")
	return nil
}

// ──────────────────────────────────────────────────────────────────────────────
// queue — Show the music queue
// ──────────────────────────────────────────────────────────────────────────────

func (a *App) cmdQueue() error {
	session, err := a.loadSession()
	if err != nil {
		return err
	}

	if len(session.Queue) == 0 {
		fmt.Println("Queue is empty. Run 'moodwave play' to build a queue.")
		return nil
	}

	fmt.Printf("\nMood: %s %s\n\n", session.MoodProfile.Label.Emoji(), session.MoodProfile.Label)
	fmt.Printf("Queue (%d tracks):\n\n", len(session.Queue))

	for i, c := range session.Queue {
		marker := "  "
		if i == session.QueueIndex {
			marker = "▶ "
		}
		score := fmt.Sprintf("%.0f%%", c.Score*100)
		reason := c.Reason
		fmt.Printf("%s%d. %-40s %s  %s\n", marker, i+1, c.DisplayName(), score, reason)
	}
	return nil
}

// ──────────────────────────────────────────────────────────────────────────────
// status — Show current status
// ──────────────────────────────────────────────────────────────────────────────

func (a *App) cmdStatus() error {
	session, err := a.loadSession()
	if err != nil {
		return err
	}

	fmt.Printf("\nMoodwave Status\n%s\n\n", strings.Repeat("─", 40))
	fmt.Printf("Version:      %s\n", config.Version)
	fmt.Printf("Config:       %s\n", a.cfg.Paths.ConfigFile)
	fmt.Printf("Cache:        %s\n", a.cfg.Paths.CacheDir)
	fmt.Printf("Project:      %s\n", a.projectRoot)

	if session.MoodProfile != nil {
		age := time.Since(session.ScannedAt).Truncate(time.Minute)
		fmt.Printf("\nMood:         %s %s (%.0f%%) — %v ago\n",
			session.MoodProfile.Label.Emoji(),
			session.MoodProfile.Label,
			session.MoodProfile.Confidence*100,
			age)
	} else {
		fmt.Printf("\nMood:         (not scanned — run 'moodwave scan')\n")
	}

	fmt.Printf("Playback:     %s\n", session.PlaybackState)

	if session.CurrentCandidate != nil {
		fmt.Printf("Now playing:  %s\n", session.CurrentCandidate.DisplayName())
	}

	fmt.Printf("\nTerminal:     %dx%d  color=%v  unicode=%v  tty=%v\n",
		a.caps.Width, a.caps.Height,
		a.caps.HasColor, a.caps.HasUnicode, a.caps.IsTTY)

	return nil
}

// ──────────────────────────────────────────────────────────────────────────────
// config — View or edit configuration
// ──────────────────────────────────────────────────────────────────────────────

func (a *App) cmdConfig(args []string) error {
	if len(args) == 0 {
		// Print current config as JSON.
		data, err := json.MarshalIndent(a.cfg, "", "  ")
		if err != nil {
			return err
		}
		fmt.Printf("Config file: %s\n\n", a.cfg.Paths.ConfigFile)
		fmt.Println(string(data))
		return nil
	}

	switch args[0] {
	case "edit":
		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "notepad"
			if os.Getenv("TERM") != "" {
				editor = "vi"
			}
		}
		fmt.Printf("Opening config in %s...\n", editor)
		// We don't exec here to keep it simple.
		fmt.Printf("File: %s\n", a.cfg.Paths.ConfigFile)
	case "path":
		fmt.Println(a.cfg.Paths.ConfigFile)
	default:
		fmt.Printf("Unknown config subcommand: %s\nUsage: moodwave config [edit|path]\n", args[0])
	}
	return nil
}

// ──────────────────────────────────────────────────────────────────────────────
// theme — Switch visual theme
// ──────────────────────────────────────────────────────────────────────────────

func (a *App) cmdTheme(args []string) error {
	themes := []config.ThemeID{
		config.ThemeMonochrome,
		config.ThemeDark,
		config.ThemeAsh,
		config.ThemeGhost,
	}

	if len(args) == 0 {
		fmt.Println("Available themes:")
		for _, t := range themes {
			marker := "  "
			if t == a.cfg.Visual.Theme {
				marker = "▶ "
			}
			fmt.Printf("%s%s\n", marker, t)
		}
		fmt.Printf("\nCurrent: %s\nUse 'moodwave theme <name>' to switch.\n", a.cfg.Visual.Theme)
		return nil
	}

	newTheme := config.ThemeID(args[0])
	for _, t := range themes {
		if t == newTheme {
			fmt.Printf("Switching theme to: %s\n", newTheme)
			// Theme is runtime-only for now; config persistence can be added.
			return nil
		}
	}
	return fmt.Errorf("unknown theme: %s", args[0])
}

// ──────────────────────────────────────────────────────────────────────────────
// visual — Switch visual mode
// ──────────────────────────────────────────────────────────────────────────────

func (a *App) cmdVisual(args []string) error {
	modes := []config.VisualMode{
		config.VisualWave,
		config.VisualSpectrum,
		config.VisualPulse,
		config.VisualMinimal,
		config.VisualQuiet,
	}

	if len(args) == 0 {
		fmt.Println("Available visual modes:")
		for _, m := range modes {
			marker := "  "
			if m == a.cfg.Visual.Mode {
				marker = "▶ "
			}
			fmt.Printf("%s%s\n", marker, m)
		}
		fmt.Printf("\nCurrent: %s\nUse 'moodwave visual <mode>' to switch.\n", a.cfg.Visual.Mode)
		return nil
	}

	newMode := config.VisualMode(args[0])
	for _, m := range modes {
		if m == newMode {
			fmt.Printf("Switching visual mode to: %s\n", newMode)
			return nil
		}
	}
	return fmt.Errorf("unknown visual mode: %s", args[0])
}

// ──────────────────────────────────────────────────────────────────────────────
// source — View or switch music sources
// ──────────────────────────────────────────────────────────────────────────────

func (a *App) cmdSource(args []string) error {
	registry := a.buildSourceRegistry()

	if len(args) == 0 {
		fmt.Println("Music sources (in priority order):")
		fmt.Println()
		for i, adapter := range registry.All() {
			ctx, cancel := context.WithTimeout(a.ctx, 5*time.Second)
			status := "✓ healthy"
			if err := adapter.HealthCheck(ctx); err != nil {
				status = "✗ " + err.Error()
			}
			cancel()
			fmt.Printf("  %d. %-20s %s\n", i+1, adapter.Name(), status)
		}
		fmt.Println("\nSet JAMENDO_CLIENT_ID or LISTENBRAINZ_TOKEN to enable optional sources.")
		return nil
	}

	switch args[0] {
	case "list":
		return a.cmdSource(nil)
	default:
		fmt.Printf("Unknown source subcommand: %s\nUsage: moodwave source [list]\n", args[0])
	}
	return nil
}

// ──────────────────────────────────────────────────────────────────────────────
// doctor — Run diagnostics
// ──────────────────────────────────────────────────────────────────────────────

func (a *App) cmdDoctor() error {
	runChecks := func() []tui.DoctorCheck {
		var checks []tui.DoctorCheck

		// ── System ──
		checks = append(checks, tui.DoctorCheck{
			Category: "System",
			Label:    "Platform",
			Status:   tui.CheckPass,
			Detail:   fmt.Sprintf("%s/%s", a.caps.OS, a.caps.Arch),
		})

		// ── Terminal ──
		if a.caps.IsTTY {
			checks = append(checks, tui.DoctorCheck{
				Category: "Terminal",
				Label:    "TTY",
				Status:   tui.CheckPass,
				Detail:   fmt.Sprintf("%dx%d", a.caps.Width, a.caps.Height),
			})
		} else {
			checks = append(checks, tui.DoctorCheck{
				Category: "Terminal",
				Label:    "TTY",
				Status:   tui.CheckWarn,
				Detail:   "Not a terminal — some features limited",
			})
		}
		if a.caps.HasColor {
			checks = append(checks, tui.DoctorCheck{
				Category: "Terminal",
				Label:    "Color support",
				Status:   tui.CheckPass,
				Detail:   "ANSI colors available",
			})
		} else {
			checks = append(checks, tui.DoctorCheck{
				Category: "Terminal",
				Label:    "Color support",
				Status:   tui.CheckWarn,
				Detail:   "No color support detected",
			})
		}

		// ── Configuration ──
		if _, err := os.Stat(a.cfg.Paths.ConfigFile); err == nil {
			checks = append(checks, tui.DoctorCheck{
				Category: "Configuration",
				Label:    "Config file",
				Status:   tui.CheckPass,
				Detail:   a.cfg.Paths.ConfigFile,
			})
		} else {
			checks = append(checks, tui.DoctorCheck{
				Category: "Configuration",
				Label:    "Config file",
				Status:   tui.CheckWarn,
				Detail:   "Missing — run 'moodwave init'",
			})
		}

		// ── Audio Backends ──
		ffStatus := playback.CheckFFplay()
		if ffStatus.Available {
			checks = append(checks, tui.DoctorCheck{
				Category: "Audio",
				Label:    "ffplay",
				Status:   tui.CheckPass,
				Detail:   fmt.Sprintf("Found (%s: %s)", ffStatus.Source, ffStatus.Path),
			})
		} else {
			checks = append(checks, tui.DoctorCheck{
				Category: "Audio",
				Label:    "ffplay",
				Status:   tui.CheckFail,
				Detail:   "Not installed — required for audio playback",
				CanFix:   true,
				FixLabel: "Auto-download ffplay",
			})
		}

		// Check mpv (optional)
		if _, err := findExecutable("mpv"); err == nil {
			checks = append(checks, tui.DoctorCheck{
				Category: "Audio",
				Label:    "mpv",
				Status:   tui.CheckPass,
				Detail:   "Found in PATH",
			})
		}

		// Check yt-dlp
		if _, err := findExecutable("yt-dlp"); err == nil {
			checks = append(checks, tui.DoctorCheck{
				Category: "Audio",
				Label:    "yt-dlp",
				Status:   tui.CheckPass,
				Detail:   "YouTube stream resolver available",
			})
		} else {
			checks = append(checks, tui.DoctorCheck{
				Category: "Audio",
				Label:    "yt-dlp",
				Status:   tui.CheckFail,
				Detail:   "Not found — required for YouTube playback",
				CanFix:   true,
				FixLabel: "Auto-download yt-dlp",
			})
		}

		// ── Music Sources ──
		registry := a.buildSourceRegistry()
		ctx, cancel := context.WithTimeout(a.ctx, 20*time.Second)
		defer cancel()

		for _, adapter := range registry.All() {
			err := adapter.HealthCheck(ctx)
			if err == nil {
				checks = append(checks, tui.DoctorCheck{
					Category: "Sources",
					Label:    adapter.Name(),
					Status:   tui.CheckPass,
					Detail:   "Connected",
				})
			} else {
				checks = append(checks, tui.DoctorCheck{
					Category: "Sources",
					Label:    adapter.Name(),
					Status:   tui.CheckFail,
					Detail:   err.Error(),
				})
			}
		}

		// ── Cache ──
		c, err := cache.New(
			a.cfg.Paths.CacheDir,
			a.cfg.Cache.MaxEntries,
			time.Duration(a.cfg.Cache.TTLSecs)*time.Second,
		)
		if err == nil {
			count, max := c.Stats()
			checks = append(checks, tui.DoctorCheck{
				Category: "Storage",
				Label:    "Cache",
				Status:   tui.CheckPass,
				Detail:   fmt.Sprintf("%d/%d entries", count, max),
			})
		}

		// ── Session ──
		session, _ := a.loadSession()
		if session.MoodProfile != nil {
			age := time.Since(session.ScannedAt).Truncate(time.Minute)
			checks = append(checks, tui.DoctorCheck{
				Category: "Storage",
				Label:    "Session",
				Status:   tui.CheckPass,
				Detail:   fmt.Sprintf("Last scan %v ago (mood: %s)", age, session.MoodProfile.Label),
			})
		} else {
			checks = append(checks, tui.DoctorCheck{
				Category: "Storage",
				Label:    "Session",
				Status:   tui.CheckWarn,
				Detail:   "No scan yet — run 'moodwave scan'",
			})
		}

		return checks
	}

	fixCheck := func(idx int) error {
		checks := runChecks()
		if idx >= len(checks) {
			return fmt.Errorf("invalid check index")
		}

		check := checks[idx]
		switch check.Label {
		case "ffplay":
			ctx, cancel := context.WithTimeout(a.ctx, 5*time.Minute)
			defer cancel()
			_, err := playback.InstallFFplay(ctx, nil)
			return err
		case "yt-dlp":
			return installYtDlp(a.ctx)
		default:
			return fmt.Errorf("no auto-fix available for: %s", check.Label)
		}
	}

	// If not a TTY, fall back to plain text output
	if !a.caps.IsTTY {
		checks := runChecks()
		fmt.Printf("\nMoodwave Doctor\n%s\n\n", strings.Repeat("═", 50))
		currentCat := ""
		for _, c := range checks {
			if c.Category != currentCat {
				if currentCat != "" {
					fmt.Println()
				}
				currentCat = c.Category
				fmt.Printf("%s:\n", c.Category)
			}
			icon := "✓"
			switch c.Status {
			case tui.CheckFail:
				icon = "✗"
			case tui.CheckWarn:
				icon = "⚠"
			}
			detail := ""
			if c.Detail != "" {
				detail = " — " + c.Detail
			}
			fmt.Printf("  %s %s%s\n", icon, c.Label, detail)
		}
		fmt.Printf("\n%s\nDiagnostics complete.\n", strings.Repeat("═", 50))
		return nil
	}

	return tui.RunDoctor(runChecks, fixCheck)
}

// ──────────────────────────────────────────────────────────────────────────────
// Helpers
// ──────────────────────────────────────────────────────────────────────────────

// buildSourceRegistry creates and registers all configured source adapters.
func (a *App) buildSourceRegistry() *sources.Registry {
	registry := sources.NewRegistry()

	// YouTube search (no auth needed).
	if a.cfg.Sources.YouTube.Enabled {
		yt := sources.NewYouTubeAdapter()
		registry.Register(yt)
	}

	// Always try Radio Browser first (no auth needed).
	if a.cfg.Sources.RadioBrowser.Enabled {
		rb := sources.NewRadioBrowserAdapter(a.cfg.Sources.RadioBrowser.APIBaseURL)
		registry.Register(rb)
	}

	// Jamendo (optional — needs client_id).
	if a.cfg.Sources.Jamendo.Enabled || a.cfg.Sources.Jamendo.ClientID != "" {
		j := sources.NewJamendoAdapter(a.cfg.Sources.Jamendo.ClientID)
		registry.Register(j)
	}

	// MusicBrainz metadata.
	if a.cfg.Sources.MusicBrainz.Enabled {
		mb := sources.NewMusicBrainzAdapter(a.cfg.Sources.MusicBrainz.UserAgent)
		registry.Register(mb)
	}

	// LRCLIB lyrics.
	if a.cfg.Sources.LRCLIB.Enabled {
		lrc := sources.NewLRCLIBAdapter()
		registry.Register(lrc)
	}

	return registry
}

// findExecutable searches PATH and moodwave cache for a binary.
func findExecutable(name string) (string, error) {
	// Check system PATH first
	paths := strings.Split(os.Getenv("PATH"), string(os.PathListSeparator))
	for _, dir := range paths {
		full := filepath.Join(dir, name)
		if _, err := os.Stat(full); err == nil {
			return full, nil
		}
		// Windows .exe
		fullExe := full + ".exe"
		if _, err := os.Stat(fullExe); err == nil {
			return fullExe, nil
		}
	}
	// Also check moodwave cache directory (auto-downloaded binaries)
	if cacheDir, err := os.UserCacheDir(); err == nil {
		mwCache := filepath.Join(cacheDir, "moodwave")
		full := filepath.Join(mwCache, name)
		if _, err := os.Stat(full); err == nil {
			return full, nil
		}
		fullExe := full + ".exe"
		if _, err := os.Stat(fullExe); err == nil {
			return fullExe, nil
		}
	}
	return "", fmt.Errorf("%s not found", name)
}

// installYtDlp downloads yt-dlp to the moodwave cache directory.
func installYtDlp(ctx context.Context) error {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return fmt.Errorf("cannot determine cache dir: %w", err)
	}
	mwCache := filepath.Join(cacheDir, "moodwave")
	_ = os.MkdirAll(mwCache, 0755)

	execName := "yt-dlp"
	var dlURL string
	switch runtime.GOOS {
	case "windows":
		execName = "yt-dlp.exe"
		dlURL = "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp.exe"
	case "darwin":
		dlURL = "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_macos"
	default:
		dlURL = "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_linux"
	}

	targetPath := filepath.Join(mwCache, execName)
	client := &http.Client{Timeout: 2 * time.Minute}
	req, err := http.NewRequestWithContext(ctx, "GET", dlURL, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("download returned status %d", resp.StatusCode)
	}

	tmpPath := targetPath + ".tmp"
	out, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}

	_, err = io.Copy(out, resp.Body)
	out.Close()
	if err != nil {
		os.Remove(tmpPath)
		return err
	}

	_ = os.Remove(targetPath)
	return os.Rename(tmpPath, targetPath)
}

// min3 returns the minimum of two ints.
func min3(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Ensure platform import is used.
var _ = platform.Detect

func (a *App) cmdWelcome() error {
	// Load session to get current mood if available.
	session, _ := a.loadSession()

	var moodLabel, moodEmoji string
	var confidence float64
	if session != nil && session.MoodProfile != nil {
		moodLabel = string(session.MoodProfile.Label)
		moodEmoji = session.MoodProfile.Label.Emoji()
		confidence = session.MoodProfile.Confidence
	}

	var postAction string

	// Create scan function that runs the full scan pipeline
	scanFn := func() tui.ScanResult {
		root, err := filepath.Abs(a.projectRoot)
		if err != nil {
			return tui.ScanResult{Err: err}
		}

		opts := scanner.ScanOptions{
			MaxDepth:   a.cfg.Scanner.MaxDepth,
			MaxFiles:   a.cfg.Scanner.MaxFiles,
			GitEnabled: a.cfg.Scanner.GitEnabled,
			IgnoreDirs: a.cfg.Scanner.IgnoreDirs,
		}
		s := scanner.New(opts)
		signals, err := s.Scan(root)
		if err != nil {
			return tui.ScanResult{Err: fmt.Errorf("scan failed: %w", err)}
		}

		engine := mood.NewEngine(a.cfg.Mood.Sensitivity)
		profile := engine.Infer(signals)

		// Save to session
		newSession := &Session{
			MoodProfile: profile,
			ScannedAt:   time.Now(),
		}
		_ = a.saveSession(newSession)

		return tui.ScanResult{
			MoodLabel:  string(profile.Label),
			MoodEmoji:  profile.Label.Emoji(),
			Confidence: profile.Confidence,
		}
	}

	err := tui.Run(moodLabel, moodEmoji, confidence, func(action, args string) {
		switch action {
		case "doctor":
			postAction = "doctor"
		}
	}, scanFn)
	if err != nil {
		return err
	}

	// Handle post-TUI actions that need raw terminal
	switch postAction {
	case "doctor":
		return a.cmdDoctor()
	}

	return nil
}

func (a *App) cmdSelfUpdate() error {
	fmt.Println("\nChecking for updates...")

	result, err := updater.CheckLatest(a.ctx)
	if err != nil {
		return fmt.Errorf("checking for updates: %w", err)
	}

	if !result.Available {
		fmt.Printf("You are already on the latest version (v%s).\n", result.CurrentVersion)
		return nil
	}

	fmt.Printf("New version found: v%s (current: v%s)\n", result.LatestVersion, result.CurrentVersion)

	newTag, err := updater.Apply(a.ctx, result.Release, func(msg string) {
		fmt.Printf("%s\n", msg)
	})
	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	fmt.Printf("\nSuccessfully updated to %s! Please run moodwave again.\n", newTag)
	return nil
}
