package runc

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/opencontainers/runc/libcontainer"
	"github.com/vishvananda/netlink/nl"
	"golang.org/x/sys/unix"
)

func GetVersionByCliVersion() (ver *semver.Version, err error) {
	var out bytes.Buffer

	cmd := exec.Command("runc", "--version")
	cmd.Stdout = &out

	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to run runc --version: %w", err)
	}
	re := regexp.MustCompile(`runc version ([\w.-]+)`)
	matches := re.FindStringSubmatch(out.String())
	if len(matches) > 1 {
		match := matches[1]
		ver, err = semver.NewVersion(match)
	} else {
		return nil, fmt.Errorf("failed to parse version from output: %s", out.String())
	}
	return ver, nil
}

func GetVersionByCliSyscall() (ver *semver.Version, err error) {
	path, err := exec.LookPath("runc")
	if err != nil {
		return nil, fmt.Errorf("failed to find runc: %w", err)
	}
	parent, child, err := newPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create parent pipe: %w", err)
	}
	namespaces := []string{
		// join pid ns of the current process
		fmt.Sprintf("pid:/proc/%d/ns/pid", -1),
	}
	cmd := exec.Cmd{
		Path:       path,
		Args:       []string{"nsenter-exec"},
		ExtraFiles: []*os.File{child},
		Env:        []string{"_LIBCONTAINER_INITPIPE=3"},
		Stdout:     os.Stdout,
		Stderr:     os.Stderr,
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start runc: %w", err)
	}
	// write cloneFlags
	r := nl.NewNetlinkRequest(int(libcontainer.InitMsg), 0)
	r.AddData(&libcontainer.Int32msg{
		Type:  libcontainer.CloneFlagsAttr,
		Value: uint32(unix.CLONE_NEWNET),
	})
	r.AddData(&libcontainer.Bytemsg{
		Type:  libcontainer.NsPathsAttr,
		Value: []byte(strings.Join(namespaces, ",")),
	})
	if _, err := io.Copy(parent, bytes.NewReader(r.Serialize())); err != nil {
		return nil, fmt.Errorf("failed to copy runc init msg: %w", err)
	}

	if err := cmd.Wait(); err == nil {
		return nil, fmt.Errorf("runc init succeeded")
	}
	return
}
