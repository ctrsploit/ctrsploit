package pipeprimitive

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ctrsploit/sploit-spec/pkg/log"
)

var (
	//go:embed material/Dockerfile
	escapeImageDockerfile string
	//go:embed material/Dockerfile.exec
	escapeImageExecDockerfile string
	//go:embed material/ld.go
	escapeImageLd string
	//go:embed material/runc-capture.go
	escapeImageRuncCapture string
	//go:embed runcwatch/runcwatch.go
	escapeImageRuncWatch string
)

const (
	EscapeImageModeStart = "start"
	EscapeImageModeExec  = "exec"
)

func ValidateEscapeImageMode(mode string) error {
	switch mode {
	case EscapeImageModeStart, EscapeImageModeExec:
		return nil
	default:
		return fmt.Errorf("invalid escape image mode %q: expected %q or %q", mode, EscapeImageModeStart, EscapeImageModeExec)
	}
}

func GenerateEscapeImage(dir string, writer, payload []byte, extraFileSets ...map[string][]byte) error {
	return GenerateEscapeImageWithMode(dir, EscapeImageModeStart, writer, payload, extraFileSets...)
}

func GenerateEscapeImageWithMode(dir, mode string, writer, payload []byte, extraFileSets ...map[string][]byte) error {
	if len(writer) == 0 {
		return fmt.Errorf("escape image writer is empty")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create escape image directory %q: %w", dir, err)
	}
	baseFiles := map[string][]byte{
		"writer.go": writer,
		"payload":   payload,
	}
	switch mode {
	case EscapeImageModeStart:
		baseFiles["Dockerfile"] = []byte(escapeImageDockerfile)
		baseFiles["ld.go"] = []byte(escapeImageLd)
	case EscapeImageModeExec:
		baseFiles["Dockerfile"] = []byte(escapeImageExecDockerfile)
		baseFiles["runc-capture.go"] = []byte(escapeImageRuncCapture)
		baseFiles["runcwatch/runcwatch.go"] = []byte(escapeImageRuncWatch)
	default:
		return ValidateEscapeImageMode(mode)
	}
	for name, content := range baseFiles {
		if err := writeEscapeImageFile(dir, name, content); err != nil {
			return fmt.Errorf("write escape image file %q: %w", name, err)
		}
	}
	for _, extraFiles := range extraFileSets {
		for name, content := range extraFiles {
			cleanName, err := cleanEscapeImageFileName(name)
			if err != nil {
				return err
			}
			if _, ok := baseFiles[cleanName]; ok {
				return fmt.Errorf("extra escape image file %q would overwrite a generated file", name)
			}
			if err := writeEscapeImageFile(dir, cleanName, content); err != nil {
				return fmt.Errorf("write extra escape image file %q: %w", name, err)
			}
		}
	}
	log.Logger.Infof(
		"escape image generated in %s directory, build it with:\ndocker build -t poc %s\n",
		dir, dir,
	)
	return nil
}

func writeEscapeImageFile(dir, name string, content []byte) error {
	cleanName, err := cleanEscapeImageFileName(name)
	if err != nil {
		return err
	}
	path := filepath.Join(dir, cleanName)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, content, 0o755)
}

func cleanEscapeImageFileName(name string) (string, error) {
	cleanName := filepath.Clean(name)
	if cleanName == "." || filepath.IsAbs(cleanName) || cleanName == ".." || strings.HasPrefix(cleanName, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("invalid escape image file path %q", name)
	}
	return cleanName, nil
}
