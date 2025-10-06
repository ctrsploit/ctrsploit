package rootfs

import (
	"fmt"
	"regexp"

	"github.com/ctrsploit/ctrsploit/pkg/mountinfo"
	"github.com/ctrsploit/ctrsploit/pkg/runtime"
	mountinfo2 "github.com/moby/sys/mountinfo"
)

type DockerOverlay struct{}

func (d DockerOverlay) Is() (bool, error) {
	// 1. is docker
	if is, _ := runtime.Docker().Is(); !is {
		return false, nil
	}
	// 2. is overlay
	info, err := mountinfo.RootMount()
	if err != nil {
		return false, fmt.Errorf("error getting root's mount info: %w", err)
	}
	if !mountinfo.IsOverlay(info) {
		return false, nil
	}
	return true, nil
}

func (d DockerOverlay) RootPath() (string, error) {
	info, err := mountinfo.RootMount()
	if err != nil {
		return "", fmt.Errorf("error getting root's mount info: %w", err)
	}
	return d.parseHostPathFromMountInfo(info)
}

func (d DockerOverlay) parseHostPathFromMountInfo(info *mountinfo2.Info) (string, error) {
	pattern := regexp.MustCompile(",upperdir=(.*)/diff,")
	matches := pattern.FindStringSubmatch(info.VFSOptions)
	if len(matches) != 2 {
		return "", fmt.Errorf("unkown VFSOptions: %+v, please report an issue to ctrsploit project", info)
	}
	return matches[1] + "/merged", nil
}
