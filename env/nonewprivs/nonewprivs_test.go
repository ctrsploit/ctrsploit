package nonewprivs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFromStatusFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "status")
	if err := os.WriteFile(path, []byte("Name:\ttest\nNoNewPrivs:\t1\n"), 0o644); err != nil {
		t.Fatalf("write status fixture: %v", err)
	}

	got, err := FromStatusFile(path)
	if err != nil {
		t.Fatalf("FromStatusFile returned error: %v", err)
	}
	if got.StatusPath != path {
		t.Fatalf("StatusPath = %q, want %q", got.StatusPath, path)
	}
	if !got.Enabled {
		t.Fatalf("Enabled = false, want true")
	}
}

func TestHumanMarksSUIDBlockedWhenEnabled(t *testing.T) {
	human := Human(Info{StatusPath: "/proc/self/status", Enabled: true})
	if !human.Enabled.Result {
		t.Fatalf("Enabled.Result = false, want true")
	}
	if !human.SUIDBlocked.Result {
		t.Fatalf("SUIDBlocked.Result = false, want true")
	}
}
