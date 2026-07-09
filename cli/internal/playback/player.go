// Package playback provides the audio playback controller.
//
// Design: The playback layer shells out to a platform-native audio backend
// (mpv, ffplay, afplay, or Windows PowerShell). This keeps the binary small,
// avoids C/CGo audio library complexity, and supports virtually any stream
// format the backend handles.
//
// The PlaybackController manages a single background process representing
// the current stream. State transitions (play/pause/stop/next) are explicit
// and safe to call from any goroutine.
package playback

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// State represents the current playback state.
type State int

const (
	StateIdle State = iota
	StatePlaying
	StatePaused
	StateStopped
	StateBuffering
	StateError
)

// String returns a display-friendly state name.
func (s State) String() string {
	switch s {
	case StateIdle:
		return "idle"
	case StatePlaying:
		return "playing"
	case StatePaused:
		return "paused"
	case StateStopped:
		return "stopped"
	case StateBuffering:
		return "buffering"
	case StateError:
		return "error"
	default:
		return "unknown"
	}
}

// Status is the current playback snapshot.
type Status struct {
	State     State
	StreamURL string
	Title     string
	Artist    string
	Backend   string
	StartedAt time.Time
	Error     string
}

// Controller manages audio playback via a subprocess backend.
type Controller struct {
	mu      sync.Mutex
	status  Status
	backend Backend
	cmd     *exec.Cmd
	cancel  context.CancelFunc
}

// NewController creates a playback controller.
// backendOverride may be empty for auto-detection.
func NewController(backendOverride string) (*Controller, error) {
	b, err := detectBackend(backendOverride)
	if err != nil {
		return nil, fmt.Errorf("playback: no audio backend found: %w", err)
	}
	c := &Controller{
		backend: b,
		status:  Status{State: StateIdle, Backend: b.Name()},
	}
	return c, nil
}

// Play starts playing the stream at the given URL.
// If something is already playing, it is stopped first.
func (c *Controller) Play(ctx context.Context, streamURL, title, artist string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Stop any current playback.
	c.stopLocked()

	// Start new playback.
	playCtx, cancel := context.WithCancel(ctx)
	stderrBuf := &tailBuffer{}
	cmd, err := c.backend.Start(playCtx, streamURL, stderrBuf)
	if err != nil {
		cancel()
		c.status = Status{
			State:   StateError,
			Error:   err.Error(),
			Backend: c.backend.Name(),
		}
		return fmt.Errorf("playback: start failed: %w", err)
	}

	c.cmd = cmd
	c.cancel = cancel
	c.status = Status{
		State:     StatePlaying,
		StreamURL: streamURL,
		Title:     title,
		Artist:    artist,
		Backend:   c.backend.Name(),
		StartedAt: time.Now(),
	}

	// Monitor the process in a goroutine.
	go func() {
		err := cmd.Wait()
		c.mu.Lock()
		defer c.mu.Unlock()
		if c.status.State == StatePlaying {
			if err != nil {
				c.status.State = StateError
				c.status.Error = fmt.Sprintf("%v: %s", err, strings.TrimSpace(stderrBuf.String()))
			} else {
				c.status.State = StateStopped
			}
		}
	}()

	return nil
}

// Pause sends a pause signal to the backend (if supported).
// For subprocess-based backends without pause support, this is a no-op.
func (c *Controller) Pause() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.status.State == StatePlaying {
		c.backend.Pause(c.cmd)
		c.status.State = StatePaused
	}
}

// Resume resumes a paused stream.
func (c *Controller) Resume() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.status.State == StatePaused {
		c.backend.Resume(c.cmd)
		c.status.State = StatePlaying
	}
}

// Stop terminates playback and releases all resources.
func (c *Controller) Stop() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.stopLocked()
}

// stopLocked terminates the current subprocess. Must be called with c.mu held.
func (c *Controller) stopLocked() {
	if c.cancel != nil {
		c.cancel()
		c.cancel = nil
	}
	if c.cmd != nil && c.cmd.Process != nil {
		_ = c.cmd.Process.Kill()
		c.cmd = nil
	}
	c.status.State = StateStopped
}

// Status returns the current playback status (safe for concurrent reads).
func (c *Controller) Status() Status {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.status
}

// IsPlaying returns true if audio is currently playing.
func (c *Controller) IsPlaying() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.status.State == StatePlaying
}

// BackendName returns the name of the detected audio backend.
func (c *Controller) BackendName() string {
	return c.backend.Name()
}

// ──────────────────────────────────────────────────────────────────────────────
// Backend interface
// ──────────────────────────────────────────────────────────────────────────────

// Backend is the interface for platform audio backends.
type Backend interface {
	Name() string
	// Start launches the backend with the given stream URL.
	Start(ctx context.Context, streamURL string, stderr io.Writer) (*exec.Cmd, error)
	// Pause sends a pause signal (best-effort).
	Pause(cmd *exec.Cmd)
	// Resume sends a resume signal (best-effort).
	Resume(cmd *exec.Cmd)
}

// detectBackend finds the best available audio backend.
func detectBackend(override string) (Backend, error) {
	candidates := backendsFor(runtime.GOOS, override)
	for _, b := range candidates {
		if b.Available() {
			return b, nil
		}
	}
	return nil, fmt.Errorf("none of %v found in PATH", backendNames(candidates))
}

// backendNames returns a slice of backend names for error messages.
func backendNames(bs []availableBackend) []string {
	names := make([]string, len(bs))
	for i, b := range bs {
		names[i] = b.Name()
	}
	return names
}

// availableBackend extends Backend with an availability check.
type availableBackend interface {
	Backend
	Available() bool
}

// backendsFor returns the ordered list of backends to try for this OS.
func backendsFor(goos, override string) []availableBackend {
	all := []availableBackend{
		&mpvBackend{},
		&ffplayBackend{},
	}

	switch goos {
	case "darwin":
		all = append([]availableBackend{&afplayBackend{}}, all...)
	case "windows":
		all = append(all, &windowsMediaBackend{})
	}

	if override != "" {
		// If an override is set, try it first.
		for i, b := range all {
			if b.Name() == override {
				return append([]availableBackend{all[i]}, append(all[:i], all[i+1:]...)...)
			}
		}
	}

	return all
}

func findYtDlp() string {
	execName := "yt-dlp"
	if runtime.GOOS == "windows" {
		execName = "yt-dlp.exe"
	}
	if p, err := exec.LookPath(execName); err == nil {
		return p
	}
	if cacheDir, err := os.UserCacheDir(); err == nil {
		p := filepath.Join(cacheDir, "moodwave", execName)
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return execName
}

func cacheDirBin() string {
	if cacheDir, err := os.UserCacheDir(); err == nil {
		return filepath.Join(cacheDir, "moodwave")
	}
	return ""
}

// ──────────────────────────────────────────────────────────────────────────────
// mpv backend
// ──────────────────────────────────────────────────────────────────────────────

type mpvBackend struct{}

func (b *mpvBackend) Name() string { return "mpv" }

func (b *mpvBackend) Available() bool {
	_, err := exec.LookPath("mpv")
	return err == nil
}

func (b *mpvBackend) Start(ctx context.Context, streamURL string, stderr io.Writer) (*exec.Cmd, error) {
	cmd := exec.CommandContext(ctx, "mpv",
		"--no-video",
		"--no-terminal",
		"--quiet",
		"--cache=yes",
		"--cache-secs=30",
		streamURL,
	)
	if binDir := cacheDirBin(); binDir != "" {
		pathKey := "PATH"
		if runtime.GOOS == "windows" {
			pathKey = "Path"
		}
		pathVal := os.Getenv(pathKey)
		cmd.Env = append(os.Environ(), fmt.Sprintf("%s=%s%c%s", pathKey, binDir, os.PathListSeparator, pathVal))
	}
	cmd.Stdout = io.Discard
	cmd.Stderr = stderr
	return cmd, cmd.Start()
}

func (b *mpvBackend) Pause(cmd *exec.Cmd)  { pauseProcess(cmd) }
func (b *mpvBackend) Resume(cmd *exec.Cmd) { resumeProcess(cmd) }

// ──────────────────────────────────────────────────────────────────────────────
// ffplay backend
// ──────────────────────────────────────────────────────────────────────────────

type ffplayBackend struct{}

func (b *ffplayBackend) Name() string { return "ffplay" }

func (b *ffplayBackend) Available() bool {
	_, err := exec.LookPath("ffplay")
	return err == nil
}

func streamURLToWriter(ctx context.Context, streamURL string, w io.WriteCloser, isStatic bool) {
	defer w.Close()

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	downloaded := int64(0)
	retries := 15

	for i := 0; i < retries; i++ {
		select {
		case <-ctx.Done():
			return
		default:
		}

		req, err := http.NewRequestWithContext(ctx, "GET", streamURL, nil)
		if err != nil {
			return
		}

		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36")
		if isStatic && downloaded > 0 {
			req.Header.Set("Range", fmt.Sprintf("bytes=%d-", downloaded))
		}

		resp, err := client.Do(req)
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		if resp.StatusCode != 200 && resp.StatusCode != 206 {
			resp.Body.Close()
			if resp.StatusCode == 403 {
				// Expired signature / Forbidden, stop retry
				return
			}
			time.Sleep(1 * time.Second)
			continue
		}

		buf := make([]byte, 32*1024)
		for {
			select {
			case <-ctx.Done():
				resp.Body.Close()
				return
			default:
			}

			n, readErr := resp.Body.Read(buf)
			if n > 0 {
				writeN, writeErr := w.Write(buf[:n])
				if writeN > 0 {
					downloaded += int64(writeN)
				}
				if writeErr != nil {
					resp.Body.Close()
					return
				}
			}
			if readErr != nil {
				resp.Body.Close()
				if readErr == io.EOF {
					return
				}
				break // network/read error, break to reconnect
			}
		}

		time.Sleep(500 * time.Millisecond)
	}
}

func (b *ffplayBackend) Start(ctx context.Context, streamURL string, stderr io.Writer) (*exec.Cmd, error) {
	isStatic := strings.Contains(streamURL, "googlevideo.com") || strings.Contains(streamURL, "jamendo.com") || strings.Contains(streamURL, "storage.jamendo.com")

	if isStatic {
		ffplayCmd := exec.CommandContext(ctx, "ffplay",
			"-nodisp",
			"-autoexit",
			"-infbuf",
			"-loglevel", "quiet",
			"-",
		)

		pipe, err := ffplayCmd.StdinPipe()
		if err != nil {
			return nil, err
		}

		ffplayCmd.Stdout = io.Discard
		ffplayCmd.Stderr = stderr

		if err := ffplayCmd.Start(); err != nil {
			return nil, fmt.Errorf("failed to start ffplay: %w", err)
		}

		go streamURLToWriter(ctx, streamURL, pipe, isStatic)

		return ffplayCmd, nil
	}

	cmd := exec.CommandContext(ctx, "ffplay",
		"-nodisp",
		"-autoexit",
		"-infbuf",
		"-loglevel", "quiet",
		"-reconnect", "1",
		"-reconnect_at_eof", "1",
		"-reconnect_streamed", "1",
		"-reconnect_delay_max", "5",
		"-user_agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36",
		streamURL,
	)
	cmd.Stdout = io.Discard
	cmd.Stderr = stderr
	return cmd, cmd.Start()
}

func (b *ffplayBackend) Pause(cmd *exec.Cmd)  { pauseProcess(cmd) }
func (b *ffplayBackend) Resume(cmd *exec.Cmd) { resumeProcess(cmd) }

// ──────────────────────────────────────────────────────────────────────────────
// afplay backend (macOS)
// ──────────────────────────────────────────────────────────────────────────────

type afplayBackend struct{}

func (b *afplayBackend) Name() string { return "afplay" }

func (b *afplayBackend) Available() bool {
	_, err := exec.LookPath("afplay")
	return err == nil
}

func (b *afplayBackend) Start(ctx context.Context, streamURL string, stderr io.Writer) (*exec.Cmd, error) {
	// afplay can play local files; for streams, use curl | afplay pipe.
	// For radio streams, fallback to mpv.
	cmd := exec.CommandContext(ctx, "afplay", streamURL)
	cmd.Stdout = io.Discard
	cmd.Stderr = stderr
	return cmd, cmd.Start()
}

func (b *afplayBackend) Pause(cmd *exec.Cmd)  { pauseProcess(cmd) }
func (b *afplayBackend) Resume(cmd *exec.Cmd) { resumeProcess(cmd) }

// ──────────────────────────────────────────────────────────────────────────────
// Windows PowerShell Media.SoundPlayer backend
// ──────────────────────────────────────────────────────────────────────────────

type windowsMediaBackend struct{}

func (b *windowsMediaBackend) Name() string { return "powershell" }

func (b *windowsMediaBackend) Available() bool {
	if runtime.GOOS != "windows" {
		return false
	}
	_, err := exec.LookPath("powershell")
	return err == nil
}

func (b *windowsMediaBackend) Start(ctx context.Context, streamURL string, stderr io.Writer) (*exec.Cmd, error) {
	// Escape single quotes for PowerShell single-quoted string literal
	escapedURL := strings.ReplaceAll(streamURL, "'", "''")
	psScript := fmt.Sprintf(`
Add-Type -AssemblyName PresentationCore
$player = New-Object System.Windows.Media.MediaPlayer
$player.Volume = 0.8
$player.Open([Uri]'%s')
$player.Play()
while ($true) { Start-Sleep -Seconds 1 }
`, escapedURL)

	cmd := exec.CommandContext(ctx, "powershell",
		"-NonInteractive",
		"-NoProfile",
		"-STA",
		"-Command", psScript,
	)
	cmd.Stdout = io.Discard
	cmd.Stderr = stderr
	return cmd, cmd.Start()
}

func (b *windowsMediaBackend) Pause(cmd *exec.Cmd)  { pauseProcess(cmd) }
func (b *windowsMediaBackend) Resume(cmd *exec.Cmd) { resumeProcess(cmd) }

// NoopController is a playback controller that does nothing.
// Used when no audio backend is available — allows the CLI to work in
// recommendations-only mode.
type NoopController struct {
	mu     sync.Mutex
	status Status
}

// NewNoopController creates a no-op playback controller.
func NewNoopController() *NoopController {
	return &NoopController{
		status: Status{State: StateIdle, Backend: "none"},
	}
}

func (c *NoopController) Play(_ context.Context, streamURL, title, artist string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.status = Status{
		State:     StatePlaying, // show as "playing" for UX, even though no audio
		StreamURL: streamURL,
		Title:     title,
		Artist:    artist,
		Backend:   "none (no audio backend)",
		StartedAt: time.Now(),
	}
	return nil
}

func (c *NoopController) Pause() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.status.State = StatePaused
}

func (c *NoopController) Resume() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.status.State = StatePlaying
}

func (c *NoopController) Stop() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.status.State = StateStopped
}

func (c *NoopController) Status() Status {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.status
}

func (c *NoopController) IsPlaying() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.status.State == StatePlaying
}

func (c *NoopController) BackendName() string {
	return "none"
}

type tailBuffer struct {
	mu   sync.Mutex
	data []byte
}

func (tb *tailBuffer) Write(p []byte) (n int, err error) {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	tb.data = append(tb.data, p...)
	if len(tb.data) > 4096 {
		tb.data = tb.data[len(tb.data)-4096:]
	}
	return len(p), nil
}

func (tb *tailBuffer) String() string {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	return string(tb.data)
}
