package runc

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"regexp"

	"github.com/Masterminds/semver/v3"
)

func LookRunC() (string, error) {
	path, e1 := exec.LookPath("runc")
	if e1 == nil {
		return path, nil
	}
	path, e2 := exec.LookPath("docker-runc")
	if e2 == nil {
		return path, nil
	}
	return "", fmt.Errorf("failed to find runc binary: %w", errors.Join(e1, e2))
}

func GetVersionByCliVersion() (ver *semver.Version, err error) {
	path, err := LookRunC()
	if err != nil {
		return nil, fmt.Errorf("could not find runc binary: %v", err)
	}
	var out bytes.Buffer

	cmd := exec.Command(path, "--version")
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
