package buildkitd

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"

	"github.com/Masterminds/semver/v3"
	"github.com/ctrsploit/sploit-spec/pkg/log"
	"github.com/moby/buildkit/client"
	"github.com/ssst0n3/awesome_libs/awesome_error"
)

func Version(addr string) (ver *semver.Version, err error) {
	ver, err = GetVersionByAPI(addr)
	if err != nil {
		log.Logger.Warn(fmt.Sprintf("Failed to get version for %s, trying execute buildkitd --version", addr))
		err = nil
	}
	ver, err = GetVersionByCli()
	if err != nil {
		return
	}
	log.Logger.Debugf("Buildkitd's version: %s", ver.String())
	return
}

func GetVersionByAPI(addr string) (ver *semver.Version, err error) {
	ctx := context.Background()
	c, err := client.New(ctx, addr)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	info, err := c.Info(ctx)
	if err != nil {
		awesome_error.CheckWarning(err)
		return
	}
	ver, err = semver.NewVersion(info.BuildkitVersion.Version)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	return
}

func GetVersionByCli() (ver *semver.Version, err error) {
	output, err := exec.Command("buildkitd", "--version").Output()
	if err != nil {
		awesome_error.CheckWarning(err)
		return
	}
	re := regexp.MustCompile(`buildkitd github.com/moby/buildkit ([\w.-]+) `)
	matches := re.FindStringSubmatch(string(output))
	if len(matches) > 1 {
		ver, err = semver.NewVersion(matches[1])
		if err != nil {
			awesome_error.CheckErr(err)
			return
		}
	} else {
		err = fmt.Errorf("failed to parse buildkitd version from output: %s", string(output))
		awesome_error.CheckErr(err)
		return
	}
	return
}
