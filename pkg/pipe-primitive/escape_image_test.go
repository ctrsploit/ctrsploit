package pipeprimitive

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateEscapeImage(t *testing.T) {
	dir := t.TempDir()
	writer := []byte("package main\nfunc main() {}\n")
	extraFiles := map[string][]byte{
		"go.mod":              []byte("module test\n"),
		"writerdep/helper.go": []byte("package writerdep\n"),
	}
	if err := GenerateEscapeImage(dir, writer, ShellPayload("touch /escaped"), extraFiles); err != nil {
		t.Fatalf("GenerateEscapeImage returned error: %v", err)
	}

	for _, name := range []string{"Dockerfile", "ld.go", "writer.go", "payload", "go.mod", filepath.Join("writerdep", "helper.go")} {
		if _, err := os.Stat(filepath.Join(dir, name)); err != nil {
			t.Fatalf("expected %s to be generated: %v", name, err)
		}
	}

	payload, err := os.ReadFile(filepath.Join(dir, "payload"))
	if err != nil {
		t.Fatalf("read payload: %v", err)
	}
	if string(payload) != "#!/bin/bash\ntouch /escaped\n" {
		t.Fatalf("payload = %q", payload)
	}
}

func TestGenerateEscapeImageExecMode(t *testing.T) {
	dir := t.TempDir()
	writer := []byte("package main\nfunc main() {}\n")
	if err := GenerateEscapeImageWithMode(dir, EscapeImageModeExec, writer, ShellPayload("touch /escaped")); err != nil {
		t.Fatalf("GenerateEscapeImageWithMode returned error: %v", err)
	}

	for _, name := range []string{"Dockerfile", "runc-capture.go", "runcwatch/runcwatch.go", "writer.go", "payload"} {
		if _, err := os.Stat(filepath.Join(dir, name)); err != nil {
			t.Fatalf("expected %s to be generated: %v", name, err)
		}
	}
	if _, err := os.Stat(filepath.Join(dir, "ld.go")); !os.IsNotExist(err) {
		t.Fatalf("ld.go should not be generated in exec mode: %v", err)
	}

	dockerfile, err := os.ReadFile(filepath.Join(dir, "Dockerfile"))
	if err != nil {
		t.Fatalf("read Dockerfile: %v", err)
	}
	content := string(dockerfile)
	for _, want := range []string{"HEALTHCHECK", "CMD [\"/proc/self/exe\"]", "ENTRYPOINT [\"/runc-capture\"]", "runc exec"} {
		if !strings.Contains(content, want) {
			t.Fatalf("Dockerfile does not contain %q:\n%s", want, content)
		}
	}
}

func TestGenerateEscapeImageRejectsUnsupportedMode(t *testing.T) {
	err := GenerateEscapeImageWithMode(t.TempDir(), "unknown", []byte("package main\nfunc main() {}\n"), ShellPayload("true"))
	if err == nil {
		t.Fatal("expected unsupported mode error")
	}
}
