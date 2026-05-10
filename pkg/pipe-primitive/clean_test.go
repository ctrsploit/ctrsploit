package pipeprimitive

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCleanDropsAllCaches(t *testing.T) {
	path := filepath.Join(t.TempDir(), "drop_caches")
	if err := os.WriteFile(path, nil, 0o644); err != nil {
		t.Fatalf("create drop_caches fixture: %v", err)
	}
	originalPath := dropCachesPath
	originalSync := syncFS
	dropCachesPath = path
	var synced bool
	syncFS = func() {
		synced = true
	}
	defer func() {
		dropCachesPath = originalPath
		syncFS = originalSync
	}()

	if err := Clean(false); err != nil {
		t.Fatalf("Clean returned error: %v", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read drop_caches fixture: %v", err)
	}
	if string(content) != dropCachesAll {
		t.Fatalf("drop_caches content = %q, want %q", content, dropCachesAll)
	}
	if synced {
		t.Fatal("sync should not run unless requested")
	}
}

func TestCleanSyncsBeforeDroppingCachesWhenRequested(t *testing.T) {
	path := filepath.Join(t.TempDir(), "drop_caches")
	if err := os.WriteFile(path, nil, 0o644); err != nil {
		t.Fatalf("create drop_caches fixture: %v", err)
	}
	originalPath := dropCachesPath
	originalSync := syncFS
	dropCachesPath = path
	var synced bool
	syncFS = func() {
		synced = true
	}
	defer func() {
		dropCachesPath = originalPath
		syncFS = originalSync
	}()

	if err := Clean(true); err != nil {
		t.Fatalf("Clean returned error: %v", err)
	}
	if !synced {
		t.Fatal("sync was not called")
	}
}

func TestCleanWrapsDropCachesError(t *testing.T) {
	originalPath := dropCachesPath
	originalSync := syncFS
	dropCachesPath = filepath.Join(t.TempDir(), "missing", "drop_caches")
	syncFS = func() {}
	defer func() {
		dropCachesPath = originalPath
		syncFS = originalSync
	}()

	if err := Clean(false); err == nil {
		t.Fatal("expected error")
	}
}
