package pipeprimitive

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type restartRecordingPrimitive struct {
	recordingPrimitive
	loader []byte
}

func (p *restartRecordingPrimitive) RestartLoader(payload []byte) ([]byte, error) {
	if len(p.loader) == 0 {
		return nil, nil
	}
	return append(append([]byte{}, p.loader...), payload...), nil
}

func TestWriteRestartMaterialWithSelfContainedLoader(t *testing.T) {
	dir := t.TempDir()
	paths := restartMaterialPaths{
		loaders: []string{filepath.Join(dir, "missing-loader"), filepath.Join(dir, "loader")},
	}
	if err := os.WriteFile(paths.loaders[1], make([]byte, 16<<20), 0o644); err != nil {
		t.Fatalf("write loader fixture: %v", err)
	}

	primitive := &restartRecordingPrimitive{loader: []byte("restart-loader:")}
	payload := []byte("restart-payload")
	if err := writeRestartMaterial(primitive, payload, paths); err != nil {
		t.Fatalf("writeRestartMaterial returned error: %v", err)
	}

	seen := map[string]string{}
	for _, write := range primitive.writes {
		seen[write.path] = write.content
	}
	if seen[paths.loaders[1]] != "restart-loader:restart-payload" {
		t.Fatalf("loader write = %q", seen[paths.loaders[1]])
	}
	if _, ok := seen[filepath.Join(dir, "writer")]; ok {
		t.Fatalf("self-contained loader path should not write restart writer: %+v", seen)
	}
	if _, ok := seen[filepath.Join(dir, "payload")]; ok {
		t.Fatalf("self-contained loader path should not write payload file: %+v", seen)
	}
	if _, ok := seen[filepath.Join(dir, "primitive")]; ok {
		t.Fatalf("self-contained loader path should not write primitive config: %+v", seen)
	}
}

func TestWriteRestartMaterialWithEmptySelfContainedLoaderFails(t *testing.T) {
	dir := t.TempDir()
	paths := restartMaterialPaths{
		loaders: []string{filepath.Join(dir, "loader")},
	}
	if err := os.WriteFile(paths.loaders[0], make([]byte, 16<<20), 0o644); err != nil {
		t.Fatalf("write loader fixture: %v", err)
	}

	primitive := &restartRecordingPrimitive{}
	err := writeRestartMaterial(primitive, []byte("restart-payload"), paths)
	if err == nil || !strings.Contains(err.Error(), "restart loader is empty") {
		t.Fatalf("writeRestartMaterial error = %v, want empty loader error", err)
	}
	if len(primitive.writes) != 0 {
		t.Fatalf("empty self-contained loader should not write files: %+v", primitive.writes)
	}
}
