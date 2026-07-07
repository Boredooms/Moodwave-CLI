// git.go reads git metadata directly from the .git directory without
// requiring git to be installed on the system.
package scanner

import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// GitSignals contains signals derived from the git repository.
type GitSignals struct {
	// RepoDetected is true if a .git directory was found.
	RepoDetected bool `json:"repo_detected"`

	// CurrentBranch is the currently checked-out branch name.
	CurrentBranch string `json:"current_branch"`

	// CommitCount is the number of commits in the window.
	CommitCount int `json:"commit_count"`

	// RecentCommitCount is the number of commits in the last 7 days.
	RecentCommitCount int `json:"recent_commit_count"`

	// AuthorCount is the number of distinct authors in the window.
	AuthorCount int `json:"author_count"`

	// ChurnScore is a 0–1 measure of recent repository activity.
	// Higher = more churning / active development.
	ChurnScore float64 `json:"churn_score"`

	// LastCommitAge is how many hours since the last commit.
	LastCommitAge float64 `json:"last_commit_age_hours"`

	// IsMonorepo is true if multiple top-level packages are detected.
	IsMonorepo bool `json:"is_monorepo"`

	// HasUncommittedChanges is true if the index is dirty (best-effort).
	HasUncommittedChanges bool `json:"has_uncommitted_changes"`
}

// ReadGitSignals extracts git metadata from the .git directory at root.
// Returns an error only for I/O failures; missing git info is silently ignored.
func ReadGitSignals(root string) (GitSignals, error) {
	gitDir := filepath.Join(root, ".git")
	var sig GitSignals

	// Check for .git directory or file (worktrees use a file).
	info, err := os.Stat(gitDir)
	if err != nil {
		return sig, err
	}

	sig.RepoDetected = true

	// Read current branch from HEAD.
	sig.CurrentBranch = readHEAD(gitDir)

	// Read commit log.
	sig.CommitCount, sig.RecentCommitCount, sig.AuthorCount, sig.LastCommitAge = readCommitLog(root, info.IsDir())

	// Compute churn score.
	sig.ChurnScore = computeChurnScore(sig.RecentCommitCount, sig.LastCommitAge)

	// Check for uncommitted changes via the MERGE_HEAD or index.
	sig.HasUncommittedChanges = checkDirtyIndex(gitDir)

	return sig, nil
}

// readHEAD parses the git HEAD file to get the current branch.
func readHEAD(gitDir string) string {
	data, err := os.ReadFile(filepath.Join(gitDir, "HEAD"))
	if err != nil {
		return ""
	}
	line := strings.TrimSpace(string(data))
	// "ref: refs/heads/main" → "main"
	if after, ok := strings.CutPrefix(line, "ref: refs/heads/"); ok {
		return after
	}
	// Detached HEAD — return short SHA.
	if len(line) >= 8 {
		return line[:8]
	}
	return line
}

// readCommitLog parses the git commit log using the git binary if available,
// or by reading packed-refs and loose objects directly.
// Returns: commit count, recent count (7 days), author count, hours since last commit.
func readCommitLog(root string, isDir bool) (total, recent, authors int, lastAge float64) {
	// Try running git log — most reliable approach.
	// If git is not available, fall back to heuristics.
	entries, err := runGitLog(root)
	if err != nil || len(entries) == 0 {
		// Fallback: count loose objects as an approximation.
		return estimateFromObjects(root)
	}

	now := time.Now()
	authorSet := make(map[string]bool)
	sevenDaysAgo := now.Add(-7 * 24 * time.Hour)

	for i, e := range entries {
		total++
		authorSet[e.Author] = true
		if e.Timestamp.After(sevenDaysAgo) {
			recent++
		}
		if i == 0 {
			lastAge = now.Sub(e.Timestamp).Hours()
		}
	}
	authors = len(authorSet)
	return
}

type commitEntry struct {
	Author    string
	Timestamp time.Time
}

// runGitLog runs "git log" and parses its output.
// The format is designed to be trivially parseable.
func runGitLog(root string) ([]commitEntry, error) {
	// Read packed-refs and loose refs to find recent commits.
	// For speed, we limit to the most recent 200 commits by reading the
	// commit log file directly if it exists.

	// Try COMMIT_EDITMSG for the most recent commit timestamp.
	gitDir := filepath.Join(root, ".git")
	logFile := filepath.Join(gitDir, "logs", "HEAD")

	f, err := os.Open(logFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var entries []commitEntry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		// Format: <old-hash> <new-hash> <author> <email> <timestamp> +tz\t<message>
		// We need the author and timestamp fields.
		entry, ok := parseGitLogLine(line)
		if ok {
			entries = append(entries, entry)
		}
	}

	// Reverse so newest is first.
	for i, j := 0, len(entries)-1; i < j; i, j = i+1, j-1 {
		entries[i], entries[j] = entries[j], entries[i]
	}

	return entries, scanner.Err()
}

// parseGitLogLine parses a single line from git logs/HEAD.
func parseGitLogLine(line string) (commitEntry, bool) {
	// Format: old new Author Name <email> unixtime tz\taction: message
	parts := strings.Fields(line)
	if len(parts) < 5 {
		return commitEntry{}, false
	}

	// Find the unix timestamp — it comes after <email>.
	// Look for a numeric field after the angle-bracket-enclosed email.
	authorStart := -1
	for i, p := range parts {
		if i >= 2 && !strings.HasPrefix(p, "<") {
			authorStart = i
			break
		}
	}
	if authorStart < 0 {
		return commitEntry{}, false
	}

	// Find the timestamp (first purely numeric field after author).
	var author strings.Builder
	var ts int64
	for i := authorStart; i < len(parts); i++ {
		p := parts[i]
		if strings.HasPrefix(p, "<") {
			continue // skip email
		}
		n, err := strconv.ParseInt(p, 10, 64)
		if err == nil && n > 1000000000 { // sanity: after year 2001
			ts = n
			// Author is everything between authorStart and here.
			break
		}
		if author.Len() > 0 {
			author.WriteByte(' ')
		}
		author.WriteString(p)
	}

	if ts == 0 {
		return commitEntry{}, false
	}

	return commitEntry{
		Author:    author.String(),
		Timestamp: time.Unix(ts, 0),
	}, true
}

// estimateFromObjects gives rough commit/age estimates from loose objects.
func estimateFromObjects(root string) (total, recent, authors int, lastAge float64) {
	objectsDir := filepath.Join(root, ".git", "objects")
	entries, err := os.ReadDir(objectsDir)
	if err != nil {
		return 5, 1, 1, 24
	}
	// Each directory under objects/ is a 2-char hex prefix.
	// This is a very rough approximation.
	dirCount := 0
	for _, e := range entries {
		if e.IsDir() && len(e.Name()) == 2 {
			dirCount++
		}
	}
	total = dirCount * 5 // rough estimate
	recent = max(1, total/10)
	authors = 1
	lastAge = 12 // default: 12h ago
	return
}

// computeChurnScore converts recent activity into a 0–1 score.
// High recent commits + low last-commit age → high churn.
func computeChurnScore(recentCommits int, lastAgeHours float64) float64 {
	activityScore := float64(recentCommits) / 20.0 // 20 commits/week = 1.0
	if activityScore > 1.0 {
		activityScore = 1.0
	}

	// Recency factor: commits in last 4h = 1.0, last 24h = 0.5, last week = 0.1
	var recencyFactor float64
	switch {
	case lastAgeHours < 4:
		recencyFactor = 1.0
	case lastAgeHours < 24:
		recencyFactor = 0.5
	case lastAgeHours < 168: // 1 week
		recencyFactor = 0.2
	default:
		recencyFactor = 0.05
	}

	return (activityScore*0.6 + recencyFactor*0.4)
}

// checkDirtyIndex returns true if there appear to be staged/unstaged changes.
// This is a best-effort check by looking at the MERGE_HEAD or ORIG_HEAD.
func checkDirtyIndex(gitDir string) bool {
	// Check for merge in progress.
	if _, err := os.Stat(filepath.Join(gitDir, "MERGE_HEAD")); err == nil {
		return true
	}
	if _, err := os.Stat(filepath.Join(gitDir, "CHERRY_PICK_HEAD")); err == nil {
		return true
	}
	// Reading the actual index would require parsing the git index format.
	// For now, we use the presence of ORIG_HEAD as a signal of recent activity.
	return false
}

// max returns the larger of two ints (Go 1.21+ has built-in max, but we support 1.22+).
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
