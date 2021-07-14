package mountinfo

import (
	"github.com/moby/sys/mountinfo"
)

func HostsMount() (info *mountinfo.Info, err error) {
	return GetMountByMountpoint("/etc/hosts")
}
