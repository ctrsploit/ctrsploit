package pipeprimitive

import (
	"bytes"
	"debug/elf"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type restartRecordingPrimitive struct {
	recordingPrimitive
	writer []byte
}

func (p *restartRecordingPrimitive) RestartWriter() []byte {
	return p.writer
}

func TestWriteRestartMaterial(t *testing.T) {
	dir := t.TempDir()
	paths := restartMaterialPaths{
		writer:    filepath.Join(dir, "writer"),
		payload:   filepath.Join(dir, "payload"),
		primitive: filepath.Join(dir, "primitive"),
		loaders:   []string{filepath.Join(dir, "missing-loader"), filepath.Join(dir, "loader")},
	}
	for _, path := range []string{paths.writer, paths.payload, paths.primitive, paths.loaders[1]} {
		if err := os.WriteFile(path, make([]byte, 16<<20), 0o644); err != nil {
			t.Fatalf("write %s fixture: %v", path, err)
		}
	}

	primitive := &restartRecordingPrimitive{writer: []byte("restart-writer")}
	payload := []byte("restart-payload")
	if err := writeRestartMaterial(primitive, payload, paths); err != nil {
		t.Fatalf("writeRestartMaterial returned error: %v", err)
	}

	seen := map[string]string{}
	for _, write := range primitive.writes {
		seen[write.path] = write.content
	}
	if seen[paths.loaders[1]] != string(restartLoaderPayload()) {
		t.Fatalf("loader write length = %d, want %d", len(seen[paths.loaders[1]]), len(restartLoaderPayload()))
	}
	if _, err := os.Stat(paths.loaders[0]); !os.IsNotExist(err) {
		t.Fatalf("missing loader should not be created: %v", err)
	}
	if seen[paths.payload] != string(payload) {
		t.Fatalf("payload write = %q, want %q", seen[paths.payload], payload)
	}
	if seen[paths.primitive] != "recording\n15\n" {
		t.Fatalf("primitive config = %q", seen[paths.primitive])
	}
	if seen[paths.writer] != "restart-writer" {
		t.Fatalf("restart writer = %q", seen[paths.writer])
	}
}

func TestRestartLoaderPayloadIsSmallELF(t *testing.T) {
	payload := restartLoaderPayload()
	if len(payload) == 0 {
		t.Fatal("restart loader payload is empty")
	}
	if len(payload) > 16<<10 {
		t.Fatalf("restart loader payload length = %d, want <= 16 KiB", len(payload))
	}
	if !strings.HasPrefix(string(payload[:4]), "\x7fELF") {
		t.Fatalf("restart loader payload does not start with ELF magic: % x", payload[:4])
	}
	if !strings.Contains(string(payload), "/writer") || !strings.Contains(string(payload), "3") {
		t.Fatal("restart loader payload does not contain expected writer arguments")
	}
}

func TestRestartLoaderPayloadIsSharedObject(t *testing.T) {
	payload := restartLoaderPayload()
	f, err := elf.NewFile(bytes.NewReader(payload))
	if err != nil {
		t.Fatalf("parse restart loader elf: %v", err)
	}
	if f.Type != elf.ET_DYN {
		t.Fatalf("restart loader elf type = %v, want %v", f.Type, elf.ET_DYN)
	}
}
