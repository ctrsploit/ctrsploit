package nvidia_container_runtime

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"

	"github.com/Masterminds/semver/v3"
)

func GetVersion() (ver *semver.Version, err error) {
	var out bytes.Buffer

	cmd := exec.Command("nvidia-container-runtime", "--version")
	cmd.Stdout = &out

	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to run nvidia-container-runtime: %w", err)
	}
	re := regexp.MustCompile(`NVIDIA Container Runtime version ([\w.-]+)`)
	matches := re.FindStringSubmatch(out.String())
	if len(matches) > 1 {
		match := matches[1]
		ver, err = semver.NewVersion(match)
	} else {
		return nil, fmt.Errorf("failed to parse version from output: %s", out.String())
	}
	return
}
