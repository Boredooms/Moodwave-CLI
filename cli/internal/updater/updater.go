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
	"net"
	"net/http"
	"os"
	"runtime"
	"strconv"
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

// apiClient is used for the release metadata lookup only. http.Client's
// Timeout covers the whole request *including* reading the body, so this
// short budget is safe here (the payload is a small JSON document) and
// keeps a flaky network from stalling CLI startup or the TUI.
func apiClient() *http.Client {
	return &http.Client{Timeout: 10 * time.Second}
}

// downloadClient is used for fetching the release binary. It deliberately
// sets NO client-level Timeout: the asset is several megabytes, and a
// whole-request deadline here would abort mid-body on any connection that
// can't finish the transfer inside it — which is exactly how self-update
// used to fail. Cancellation instead comes from the caller's context
// (see Apply), while the per-phase timeouts below still fail fast on a
// genuinely dead connection rather than hanging forever.
func downloadClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 30 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
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

	resp, err := apiClient().Do(req)
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
		// An update is only offered when the release is genuinely NEWER.
		// A plain string inequality here would be wrong: a locally built
		// dev binary ahead of the last published release (e.g. 2.1.0 local
		// vs 2.0.0 released) would be told to "update", silently
		// downgrading itself.
		Available: latestVersion != "" && IsNewerVersion(latestVersion, currentVersion),
		Release:   &release,
	}, nil
}

// IsNewerVersion reports whether release version `latest` is strictly newer
// than `current`, and therefore whether an update should be offered at all.
func IsNewerVersion(latest, current string) bool {
	return compareVersions(latest, current) > 0
}

// compareVersions compares two dotted numeric version strings, returning
// >0 if a is newer than b, <0 if older, and 0 if equivalent. Any
// pre-release/build suffix (e.g. "2.0.0-rc1", "ci-test") is ignored for
// ordering; non-numeric components compare as 0, so an unparseable version
// never falsely reads as newer.
func compareVersions(a, b string) int {
	aParts := splitVersion(a)
	bParts := splitVersion(b)

	n := len(aParts)
	if len(bParts) > n {
		n = len(bParts)
	}

	for i := 0; i < n; i++ {
		var av, bv int
		if i < len(aParts) {
			av = aParts[i]
		}
		if i < len(bParts) {
			bv = bParts[i]
		}
		if av != bv {
			if av > bv {
				return 1
			}
			return -1
		}
	}
	return 0
}

// splitVersion turns "2.0.1-rc2" into []int{2, 0, 1}.
func splitVersion(v string) []int {
	// Drop any pre-release / build metadata suffix.
	if i := strings.IndexAny(v, "-+"); i >= 0 {
		v = v[:i]
	}
	var out []int
	for _, part := range strings.Split(v, ".") {
		n, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			n = 0
		}
		out = append(out, n)
	}
	return out
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

	resp, err := downloadClient().Do(req)
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
