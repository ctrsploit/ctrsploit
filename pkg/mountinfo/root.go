package mountinfo

import (
	"strings"

	"github.com/moby/sys/mountinfo"
)

func RootMount() (info *mountinfo.Info, err error) {
	return GetMountByMountpoint("/")
}

func IsDeviceMapper(info *mountinfo.Info) bool {
	return strings.Contains(info.Source, "/mapper/")
}

func IsOverlay(info *mountinfo.Info) bool {
	return info.FSType == "overlay" && info.Source == "overlay"
}
