package mountinfo

import (
	"github.com/moby/sys/mountinfo"
	"strings"
)

func RootMount() (info *mountinfo.Info, err error) {
	return GetMountByMountpoint("/")
}

func IsDeviceMapper(info *mountinfo.Info) (dm bool) {
	return strings.Contains(info.Source, "/mapper/")
}

func IsOverlay(info *mountinfo.Info) (dm bool) {
	return info.FSType == "overlay" && info.Source == "overlay"
}
