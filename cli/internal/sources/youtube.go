// youtube.go implements the unofficial YouTube music search adapter.
//
// It searches YouTube by scraping the search results page without requiring
// any API key or authentication, and returns streamable video results.
package sources

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	youtubeUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36"
)

// YouTubeAdapter implements the Adapter interface for YouTube search.
type YouTubeAdapter struct {
	client *http.Client
}

// sharedYouTubeTransport is reused across adapter instances so TCP/TLS
// connections to YouTube can be pooled instead of renegotiated on every
// search. Explicit sub-timeouts make network stalls fail fast instead of
// silently eating the whole request budget on a slow hop.
var sharedYouTubeTransport = &http.Transport{
	DialContext: (&net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext,
	TLSHandshakeTimeout:   5 * time.Second,
	ResponseHeaderTimeout: 8 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
	IdleConnTimeout:       60 * time.Second,
	MaxIdleConnsPerHost:   4,
}

// NewYouTubeAdapter creates a YouTube adapter.
func NewYouTubeAdapter() *YouTubeAdapter {
	return &YouTubeAdapter{
		client: &http.Client{
			Timeout:   15 * time.Second,
			Transport: sharedYouTubeTransport,
		},
	}
}

// Name returns the adapter identifier.
func (a *YouTubeAdapter) Name() string {
	return "youtube"
}

// HealthCheck verifies the YouTube search page is reachable.
func (a *YouTubeAdapter) HealthCheck(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://www.youtube.com", nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", youtubeUserAgent)

	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("youtube: health check failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("youtube: server returned status %d", resp.StatusCode)
	}
	return nil
}

// SearchTracks searches YouTube for matching videos.
func (a *YouTubeAdapter) SearchTracks(ctx context.Context, q SearchQuery) ([]Track, error) {
	if q.Text == "" {
		return nil, nil
	}

	searchURL := "https://www.youtube.com/results?search_query=" + url.QueryEscape(q.Text)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, searchURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", youtubeUserAgent)

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("youtube: search failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("youtube: server returned status %d", resp.StatusCode)
	}

	// Stream the response and stop as soon as we've captured the
	// ytInitialData JSON blob — the page has megabytes of trailing
	// scripts/assets we never need, and reading all of it needlessly
	// exposes every search to the full page-load latency.
	jsonData, err := extractYtInitialData(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("youtube: reading response: %w", err)
	}
	if jsonData == "" {
		return nil, nil // not found
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		return nil, fmt.Errorf("youtube: unmarshalling JSON: %w", err)
	}

	contents, ok := navigateJSON(data, "contents", "twoColumnSearchResultsRenderer", "primaryContents", "sectionListRenderer", "contents")
	if !ok {
		return nil, nil
	}

	contentsSlice, ok := contents.([]interface{})
	if !ok || len(contentsSlice) == 0 {
		return nil, nil
	}

	firstSection, ok := contentsSlice[0].(map[string]interface{})
	if !ok {
		return nil, nil
	}

	itemSection, ok := navigateJSON(firstSection, "itemSectionRenderer", "contents")
	if !ok {
		return nil, nil
	}

	items, ok := itemSection.([]interface{})
	if !ok {
		return nil, nil
	}

	limit := coalesce(q.Limit, 10)
	tracks := make([]Track, 0, limit)
	count := 0
	// YouTube's results page occasionally repeats the same video across
	// multiple shelves/renderers (e.g. a "people also watched" style
	// insert further down the list). Without deduping by videoID, those
	// duplicates flow straight into the queue/playlist UI as visually
	// identical repeated rows with the same track number. Track what
	// we've already emitted so each video appears at most once per search.
	seen := make(map[string]bool, limit)

	for _, item := range items {
		if count >= limit {
			break
		}

		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		video, ok := itemMap["videoRenderer"].(map[string]interface{})
		if !ok {
			continue
		}

		videoID, ok := video["videoId"].(string)
		if !ok || videoID == "" {
			continue
		}
		if seen[videoID] {
			continue
		}
		seen[videoID] = true

		// Get title.
		titleObj, _ := navigateJSON(video, "title", "runs")
		titleSlice, _ := titleObj.([]interface{})
		title := ""
		if len(titleSlice) > 0 {
			firstRun, ok := titleSlice[0].(map[string]interface{})
			if ok {
				title, _ = firstRun["text"].(string)
			}
		}

		// Get channel name (artist).
		channelObj, _ := navigateJSON(video, "longBylineText", "runs")
		channelSlice, _ := channelObj.([]interface{})
		channel := ""
		if len(channelSlice) > 0 {
			firstRun, ok := channelSlice[0].(map[string]interface{})
			if ok {
				channel, _ = firstRun["text"].(string)
			}
		}

		var moodTags []string
		lowerTitle := strings.ToLower(title)
		musicKeywords := map[string][]string{
			"lofi":      {"lofi", "chill", "relax", "calm"},
			"synthwave": {"synthwave", "electronic", "retrowave", "sprint"},
			"ambient":   {"ambient", "minimal", "calm"},
			"focus":     {"focus", "study", "concentration"},
			"chill":     {"chill", "relax", "calm"},
			"workout":   {"intense", "workout", "sprint"},
			"relax":     {"relax", "calm", "minimal"},
			"metal":     {"intense", "metal", "chaotic"},
			"rock":      {"rock", "intense"},
			"classical": {"classical", "calm", "focus"},
			"coding":    {"focus", "study", "coding"},
			"jazz":      {"jazz", "calm", "late-night"},
			"dark":      {"late-night", "dark"},
		}
		for kw, tags := range musicKeywords {
			if strings.Contains(lowerTitle, kw) {
				moodTags = append(moodTags, tags...)
			}
		}

		t := Track{
			ID:        videoID,
			Source:    a.Name(),
			Title:     title,
			Artist:    channel,
			StreamURL: "https://www.youtube.com/watch?v=" + videoID,
			License:   "CC",
			MoodTags:  moodTags,
		}
		tracks = append(tracks, t)
		count++
	}

	return tracks, nil
}

// SearchStations is not applicable for YouTube.
func (a *YouTubeAdapter) SearchStations(_ context.Context, _ SearchQuery) ([]Station, error) {
	return nil, nil
}

// ResolveTrack resolves YouTube video details into a direct audio stream URL.
func (a *YouTubeAdapter) ResolveTrack(ctx context.Context, id string) (*Track, error) {
	videoURL := "https://www.youtube.com/watch?v=" + id

	// Determine yt-dlp executable name and download URL.
	execName := "yt-dlp"
	downloadURL := "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp"
	if runtime.GOOS == "windows" {
		execName = "yt-dlp.exe"
		downloadURL = "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp.exe"
	}

	// Try system PATH first.
	ytDlpPath, err := exec.LookPath(execName)
	if err != nil {
		// If not in PATH, use user cache directory.
		cacheDir, cErr := os.UserCacheDir()
		if cErr == nil {
			mwCache := filepath.Join(cacheDir, "moodwave")
			_ = os.MkdirAll(mwCache, 0755)
			ytDlpPath = filepath.Join(mwCache, execName)

			// Download if it does not exist.
			if _, statErr := os.Stat(ytDlpPath); os.IsNotExist(statErr) {
				fmt.Printf("  Downloading yt-dlp helper (one-time setup)...\n")

				tmpPath := ytDlpPath + ".tmp"
				req, rErr := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL, nil)
				if rErr == nil {
					resp, dErr := a.client.Do(req)
					if dErr == nil {
						defer resp.Body.Close()

						out, fErr := os.OpenFile(tmpPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
						if fErr == nil {
							_, copyErr := io.Copy(out, resp.Body)
							out.Close() // Close immediately to release handle!

							if copyErr == nil {
								// Rename atomically
								_ = os.Rename(tmpPath, ytDlpPath)
							} else {
								_ = os.Remove(tmpPath)
							}
						}
					}
				}
			}
		}
	}

	// If we have a valid yt-dlp executable, try it first!
	if _, statErr := os.Stat(ytDlpPath); statErr == nil {
		// Run yt-dlp to get direct audio URL + metadata.
		// Use --print to get URL, duration, title, and artist in one call.
		args := []string{
			"--print", "urls",
			"--print", "duration",
			"--print", "title",
			"--print", "channel",
			"-f", "ba[ext=m4a]/ba",
			"--no-warnings",
			"--socket-timeout", "10",
			"--retries", "2",
		}

		// Check for available JS runtimes and add them
		jsRuntimes := detectJSRuntimes()
		if jsRuntimes != "" {
			args = append(args, "--js-runtimes", jsRuntimes)
		}

		args = append(args, videoURL)

		// Use a tight timeout for yt-dlp — fall through to Invidious quickly if it fails
		ytCtx, ytCancel := context.WithTimeout(ctx, 15*time.Second)
		defer ytCancel()

		cmd := exec.CommandContext(ytCtx, ytDlpPath, args...)
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		_ = cmd.Run()
		lines := strings.Split(strings.TrimSpace(stdout.String()), "\n")

		// Parse output: line 0 = URL, line 1 = duration, line 2 = title, line 3 = channel
		if len(lines) >= 1 {
			directURL := strings.TrimSpace(lines[0])
			if strings.Contains(directURL, "googlevideo.com") {
				dur := 0
				title := "YouTube Video"
				artist := ""
				if len(lines) >= 2 {
					fmt.Sscanf(strings.TrimSpace(lines[1]), "%d", &dur)
				}
				if len(lines) >= 3 && strings.TrimSpace(lines[2]) != "" {
					title = strings.TrimSpace(lines[2])
				}
				if len(lines) >= 4 && strings.TrimSpace(lines[3]) != "" {
					artist = strings.TrimSpace(lines[3])
				}
				return &Track{
					ID:        id,
					Source:    a.Name(),
					Title:     title,
					Artist:    artist,
					Duration:  dur,
					StreamURL: directURL,
					License:   "CC",
				}, nil
			}
		}
		// If extraction failed, fall through to Invidious silently
	}

	// Fallback to Invidious public instances if local extraction failed or wasn't available!
	// Give this its own short-lived sub-timeout so a slow/unreachable
	// instances.json fetch can't silently consume whatever's left of the
	// caller's overall deadline and surface as a bare "context deadline
	// exceeded" with no useful signal.
	invCtx, invCancel := context.WithTimeout(ctx, 6*time.Second)
	defer invCancel()

	fallbackReq, err := http.NewRequestWithContext(invCtx, http.MethodGet, "https://api.invidious.io/instances.json", nil)
	if err != nil {
		return nil, fmt.Errorf("youtube: failed to resolve audio stream URL for video %s", id)
	}
	fallbackReq.Header.Set("User-Agent", youtubeUserAgent)

	resp, err := a.client.Do(fallbackReq)
	if err != nil {
		// Invidious directory unreachable — no point trying per-instance
		// lookups with no instance list. Fail clean instead of leaking the
		// raw network error.
		return nil, fmt.Errorf("youtube: failed to resolve audio stream URL for video %s", id)
	}
	defer resp.Body.Close()

	var raw [][]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		// Malformed/unexpected instances.json — fail clean instead of
		// leaking the raw decode error, matching the request-error path above.
		return nil, fmt.Errorf("youtube: failed to resolve audio stream URL for video %s", id)
	}

	// Try up to 12 instances.
	count := 0
	for _, item := range raw {
		if len(item) < 2 {
			continue
		}
		name := item[0].(string)

		apiURL := "https://" + name
		count++
		if count > 12 {
			break
		}

		streamURL := fmt.Sprintf("%s/api/v1/videos/%s", apiURL, id)
		sReq, err := http.NewRequestWithContext(ctx, http.MethodGet, streamURL, nil)
		if err != nil {
			continue
		}
		sReq.Header.Set("User-Agent", youtubeUserAgent)
		sReq.Header.Set("Accept", "application/json")

		sResp, err := a.client.Do(sReq)
		if err != nil {
			continue
		}
		defer sResp.Body.Close()

		if sResp.StatusCode != 200 {
			continue
		}

		var invData struct {
			AdaptiveFormats []struct {
				Type string `json:"type"`
				URL  string `json:"url"`
			} `json:"adaptiveFormats"`
		}
		if err := json.NewDecoder(sResp.Body).Decode(&invData); err != nil {
			continue
		}

		for _, format := range invData.AdaptiveFormats {
			if strings.Contains(format.Type, "audio/") && format.URL != "" {
				resolvedURL := format.URL
				if strings.HasPrefix(resolvedURL, "/") {
					resolvedURL = apiURL + resolvedURL
				}
				return &Track{
					ID:        id,
					Source:    a.Name(),
					Title:     "YouTube Video",
					StreamURL: resolvedURL,
					License:   "CC",
				}, nil
			}
		}
	}

	return nil, fmt.Errorf("youtube: failed to resolve audio stream URL for video %s", id)
}

// extractYtInitialData scans the search-results page incrementally and
// returns just the ytInitialData JSON blob, without buffering the rest of
// the (often multi-megabyte) HTML document. It stops reading from the
// network the moment the closing token is found, which keeps a single
// search fast even when YouTube is slow to finish sending the full page.
func extractYtInitialData(body io.Reader) (string, error) {
	const startToken = "ytInitialData = "
	const maxScan = 4 << 20 // 4MB safety cap in case markers are ever missing

	r := bufio.NewReaderSize(body, 64*1024)
	var buf bytes.Buffer
	found := false
	startIdx := -1

	for buf.Len() < maxScan {
		chunk := make([]byte, 32*1024)
		n, err := r.Read(chunk)
		if n > 0 {
			buf.Write(chunk[:n])
			if !found {
				if idx := strings.Index(buf.String(), startToken); idx >= 0 {
					found = true
					startIdx = idx + len(startToken)
				}
			}
			if found {
				tail := buf.String()[startIdx:]
				if endIdx := strings.Index(tail, ";</script>"); endIdx >= 0 {
					return tail[:endIdx], nil
				}
				if endIdx := strings.Index(tail, ";\n"); endIdx >= 0 {
					return tail[:endIdx], nil
				}
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
	}

	if found {
		// Reached EOF/cap without a clean terminator — fall back to a
		// best-effort trim at the last semicolon.
		tail := buf.String()[startIdx:]
		if endIdx := strings.LastIndex(tail, ";"); endIdx >= 0 {
			return tail[:endIdx], nil
		}
	}
	return "", nil
}

// navigateJSON traverses nested map[string]interface{} maps.
func navigateJSON(obj interface{}, keys ...string) (interface{}, bool) {
	curr := obj
	for _, key := range keys {
		m, ok := curr.(map[string]interface{})
		if !ok {
			return nil, false
		}
		curr, ok = m[key]
		if !ok {
			return nil, false
		}
	}
	return curr, true
}

// detectJSRuntimes checks for available JavaScript runtimes that yt-dlp can use.
// Returns a comma-separated list suitable for --js-runtimes flag, or empty string.
func detectJSRuntimes() string {
	var runtimes []string

	// Check for deno (preferred by yt-dlp)
	if _, err := exec.LookPath("deno"); err == nil {
		runtimes = append(runtimes, "deno")
	}

	// Check for Node.js (yt-dlp uses "node" as the runtime name)
	if _, err := exec.LookPath("node"); err == nil {
		runtimes = append(runtimes, "node")
	}

	// Check for bun
	if _, err := exec.LookPath("bun"); err == nil {
		runtimes = append(runtimes, "bun")
	}

	return strings.Join(runtimes, ",")
}
