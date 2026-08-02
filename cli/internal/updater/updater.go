// Package updater implements self-update checking and in-place binary
// replacement for the moodwave CLI. It is the single source of truth for
// "is there a newer release, and how do we fetch/install it" — both the
// `moodwave update` command and the in-TUI updater on the home screen use
// this package, so there is exactly one implementation of the GitHub
// release lookup and binary-swap logic instead of two copies that could
// drift apart.
package updater

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/moodwave/moodwave/internal/config"
)

// releaseAPIURL is the GitHub API endpoint for the latest release.
const releaseAPIURL = "https://api.github.com/repos/Boredooms/Moodwave-CLI/releases/latest"

// Asset describes a single downloadable file attached to a release.
type Asset struct {
	Name        string `json:"name"`
	DownloadURL string `json:"browser_download_url"`
}

// Release describes a GitHub release relevant to update-checking.
type Release struct {
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
}

// CheckResult is the outcome of an update check.
type CheckResult struct {
	// CurrentVersion is the version of the binary that ran this check.
	CurrentVersion string
	// LatestVersion is the newest version tag found on GitHub (without a
	// leading "v").
	LatestVersion string
	// Available is true when LatestVersion differs from CurrentVersion.
	Available bool
	// Release holds the full release metadata, needed later by Apply.
	Release *Release
}

// httpClient is shared across checks — short timeouts so a flaky network
// never blocks CLI startup or the TUI for more than a couple seconds.
func httpClient() *http.Client {
	return &http.Client{Timeout: 8 * time.Second}
}

// CheckLatest queries GitHub for the latest release and compares it
// against the currently running binary's version. It never returns an
// error for "no update available" — only for actual network/API failures,
// so callers can distinguish "checked, nothing new" from "couldn't check".
func CheckLatest(ctx context.Context) (*CheckResult, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, releaseAPIURL, nil)
	if err != nil {
		return nil, fmt.Errorf("building update request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := httpClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("checking for updates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github API returned status: %s", resp.Status)
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("decoding release info: %w", err)
	}

	currentVersion := strings.TrimPrefix(config.Version, "v")
	latestVersion := strings.TrimPrefix(release.TagName, "v")

	return &CheckResult{
		CurrentVersion: currentVersion,
		LatestVersion:  latestVersion,
		Available:      latestVersion != "" && latestVersion != currentVersion,
		Release:        &release,
	}, nil
}

// findAsset picks the release asset matching the current OS/architecture,
// skipping checksum files.
func findAsset(release *Release) (Asset, error) {
	goos := runtime.GOOS
	goarch := runtime.GOARCH
	for _, asset := range release.Assets {
		nameLower := strings.ToLower(asset.Name)
		if strings.HasSuffix(nameLower, ".sha256") {
			continue
		}
		if strings.Contains(nameLower, goos) && strings.Contains(nameLower, goarch) {
			return asset, nil
		}
	}
	return Asset{}, fmt.Errorf("no release binary found matching your OS/architecture (%s/%s)", goos, goarch)
}

// Apply downloads the release asset matching the current platform and
// replaces the currently running executable with it in place. progress,
// if non-nil, receives human-readable status lines as the update proceeds
// (safe to ignore from callers that don't want to display them, like a
// silent background check).
func Apply(ctx context.Context, release *Release, progress func(string)) (string, error) {
	report := func(msg string) {
		if progress != nil {
			progress(msg)
		}
	}

	asset, err := findAsset(release)
	if err != nil {
		return "", err
	}

	report(fmt.Sprintf("Downloading %s...", asset.Name))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, asset.DownloadURL, nil)
	if err != nil {
		return "", fmt.Errorf("building download request: %w", err)
	}

	resp, err := httpClient().Do(req)
	if err != nil {
		return "", fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download update: status %s", resp.Status)
	}

	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("locating running executable: %w", err)
	}

	oldPath := exePath + ".old"
	_ = os.Remove(oldPath)

	if err := os.Rename(exePath, oldPath); err != nil {
		return "", fmt.Errorf("renaming running binary: %w", err)
	}

	newFile, err := os.OpenFile(exePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		_ = os.Rename(oldPath, exePath)
		return "", fmt.Errorf("creating new binary file: %w", err)
	}

	if _, err := io.Copy(newFile, resp.Body); err != nil {
		newFile.Close()
		_ = os.Rename(oldPath, exePath)
		return "", fmt.Errorf("writing update to file: %w", err)
	}
	newFile.Close()

	_ = os.Remove(oldPath)

	report(fmt.Sprintf("Updated to %s", release.TagName))
	return release.TagName, nil
}
