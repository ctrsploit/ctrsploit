package docker_compose

import (
	"fmt"
	"os/exec"
	"regexp"

	"github.com/Masterminds/semver/v3"
	"github.com/ssst0n3/awesome_libs/awesome_error"
)

func Version() (ver *semver.Version, err error) {
	ver, err = GetVersionByCli()
	if err != nil {
		return
	}
	return
}

func GetVersionByCli() (ver *semver.Version, err error) {
	output, err := exec.Command("docker", "compose", "version", "--short").Output()
	if err != nil {
		awesome_error.CheckWarning(err)
		return
	}
	// Output format: v2.40.1 or 2.40.1
	versionStr := string(output)
	versionStr = regexp.MustCompile(`\s+`).ReplaceAllString(versionStr, "")

	ver, err = semver.NewVersion(versionStr)
	if err != nil {
		// Try parsing without 'v' prefix
		versionStr = regexp.MustCompile(`^v`).ReplaceAllString(versionStr, "")
		ver, err = semver.NewVersion(versionStr)
		if err != nil {
			err = fmt.Errorf("failed to parse docker compose version from output: %s", string(output))
			awesome_error.CheckErr(err)
			return
		}
	}
	return
}
