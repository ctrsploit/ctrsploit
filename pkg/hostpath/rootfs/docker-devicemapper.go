package rootfs

import (
	"fmt"
	"regexp"

	"github.com/ctrsploit/ctrsploit/pkg/mountinfo"
	"github.com/ctrsploit/ctrsploit/pkg/runtime"
	mountinfo2 "github.com/moby/sys/mountinfo"
)

type DockerDeviceMapper struct {
}

func (d DockerDeviceMapper) Is() (bool, error) {
	// 1. is docker
	if is, _ := runtime.Docker().Is(); !is {
		return false, nil
	}
	// 2. is device mapper
	info, err := mountinfo.RootMount()
	if err != nil {
		return false, fmt.Errorf("error getting root's mount info: %v", err)
	}
	if !mountinfo.IsDeviceMapper(info) {
		return false, nil
	}
	return true, nil
}

func (d DockerDeviceMapper) RootPath() (string, error) {
	info, err := mountinfo.RootMount()
	if err != nil {
		return "", fmt.Errorf("error getting root's mount info: %v", err)
	}
	return d.parseHostPathFromMountInfo(info)
}

func (d DockerDeviceMapper) parseHostPathFromMountInfo(info *mountinfo2.Info) (string, error) {
	pattern := regexp.MustCompile("-\\d+:\\d+-[0-9a-f]+-([0-9a-f]+)")
	matches := pattern.FindStringSubmatch(info.Source)
	if len(matches) != 2 {
		return "", fmt.Errorf("unknown mount source as a device mapper: %+v, please report an issue to ctrsploit project", info)
	}
	dm := matches[1]
	return fmt.Sprintf("/var/lib/docker/devicemapper/mnt/%s%s", dm, info.Root), nil
}
