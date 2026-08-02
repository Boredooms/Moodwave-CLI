package tests_test

import (
	"testing"

	"github.com/moodwave/moodwave/internal/updater"
)

// TestIsNewerVersion locks in the ordering rule that decides whether an
// update is offered. This matters because the original implementation used
// a plain string inequality, which meant a locally built binary AHEAD of
// the last published release would be told to "update" — silently
// downgrading itself. Only a genuinely newer release should qualify.
func TestIsNewerVersion(t *testing.T) {
	cases := []struct {
		latest  string
		current string
		want    bool
		reason  string
	}{
		// Genuine upgrades.
		{"2.0.1", "2.0.0", true, "patch bump"},
		{"2.1.0", "2.0.0", true, "minor bump"},
		{"3.0.0", "2.9.9", true, "major bump"},
		{"2.0.10", "2.0.9", true, "multi-digit patch compares numerically, not lexically"},
		{"1.10.0", "1.9.0", true, "multi-digit minor compares numerically, not lexically"},

		// Same version — nothing to do.
		{"2.0.0", "2.0.0", false, "identical"},
		{"2.0.0", "2.0.0", false, "identical repeated"},

		// Local build is newer than the published release: must NOT offer
		// an "update" that would actually be a downgrade.
		{"2.0.0", "2.0.1", false, "released is older patch"},
		{"2.0.0", "2.1.0", false, "released is older minor"},
		{"1.0.5", "2.0.0", false, "released is a whole major behind"},

		// Suffixes are ignored for ordering.
		{"2.0.1-rc1", "2.0.0", true, "prerelease suffix stripped, still newer"},
		{"2.0.0-rc1", "2.0.0", false, "prerelease of same version is not newer"},

		// Unparseable versions must never read as newer.
		{"ci-test", "2.0.0", false, "non-numeric release tag is not newer"},

		// Differing component counts.
		{"2.0", "2.0.0", false, "2.0 equals 2.0.0"},
		{"2.0.1", "2.0", true, "2.0.1 is newer than 2.0"},
	}

	for _, c := range cases {
		got := updater.IsNewerVersion(c.latest, c.current)
		if got != c.want {
			t.Errorf("IsNewerVersion(%q, %q) = %v, want %v (%s)",
				c.latest, c.current, got, c.want, c.reason)
		}
	}
}
