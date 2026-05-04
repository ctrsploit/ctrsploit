package pipeprimitive

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type recordingPrimitive struct {
	minOffset int64
	err       error
	writes    []recordedWrite
}

type recordedWrite struct {
	path    string
	offset  int64
	content string
}

func (p *recordingPrimitive) GetExpName() string {
	return "recording"
}

func (p *recordingPrimitive) MinOffset() int64 {
	return p.minOffset
}

func (p *recordingPrimitive) Write(path string, offset int64, content []byte) error {
	p.writes = append(p.writes, recordedWrite{path: path, offset: offset, content: string(content)})
	return p.err
}

func TestWriteImage(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "target")
	if err := os.WriteFile(target, []byte("abc-target"), 0644); err != nil {
		t.Fatalf("write target fixture: %v", err)
	}

	primitive := &recordingPrimitive{minOffset: 3}
	if err := WriteImage(primitive, target, []byte("abc-payload")); err != nil {
		t.Fatalf("WriteImage returned error: %v", err)
	}

	if len(primitive.writes) != 1 {
		t.Fatalf("writes = %d, want 1", len(primitive.writes))
	}
	write := primitive.writes[0]
	if write.path != target || write.offset != 3 || write.content != "-payload" {
		t.Fatalf("write = %+v", write)
	}
}

func TestWriteImagePrefixMismatch(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "target")
	if err := os.WriteFile(target, []byte("abc-target"), 0644); err != nil {
		t.Fatalf("write target fixture: %v", err)
	}

	primitive := &recordingPrimitive{minOffset: 3}
	err := WriteImage(primitive, target, []byte("xyz-payload"))
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "prefix mismatch") {
		t.Fatalf("error = %q, want prefix mismatch context", err)
	}
	if len(primitive.writes) != 0 {
		t.Fatalf("writes = %d, want 0", len(primitive.writes))
	}
}

func TestWriteImageWrapsPrimitiveError(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "target")
	if err := os.WriteFile(target, []byte("abc-target"), 0644); err != nil {
		t.Fatalf("write target fixture: %v", err)
	}

	writeErr := errors.New("write failed")
	primitive := &recordingPrimitive{minOffset: 3, err: writeErr}
	err := WriteImage(primitive, target, []byte("abc-payload"))
	if !errors.Is(err, writeErr) {
		t.Fatalf("error = %v, want wrapped %v", err, writeErr)
	}
	if !strings.Contains(err.Error(), "offset 3") {
		t.Fatalf("error = %q, want offset context", err)
	}
}

func TestImagePollutionUsesSourceAsPayloadAndDestinationAsTarget(t *testing.T) {
	dir := t.TempDir()
	source := filepath.Join(dir, "source")
	destination := filepath.Join(dir, "destination")
	if err := os.WriteFile(source, []byte("abc-payload"), 0644); err != nil {
		t.Fatalf("write source fixture: %v", err)
	}
	if err := os.WriteFile(destination, []byte("abc-target"), 0644); err != nil {
		t.Fatalf("write destination fixture: %v", err)
	}

	primitive := &recordingPrimitive{minOffset: 3}
	if err := ImagePollution(primitive, source, destination); err != nil {
		t.Fatalf("ImagePollution returned error: %v", err)
	}

	if len(primitive.writes) != 1 {
		t.Fatalf("writes = %d, want 1", len(primitive.writes))
	}
	write := primitive.writes[0]
	if write.path != destination || write.offset != 3 || write.content != "-payload" {
		t.Fatalf("write = %+v", write)
	}
}
