package mountinfo

import (
	"fmt"
	"github.com/moby/sys/mountinfo"
	"github.com/ssst0n3/awesome_libs/awesome_error"
)

func GetMountByMountpoint(mountpoint string) (info *mountinfo.Info, err error) {
	mounts, err := mountinfo.GetMounts(func(info *mountinfo.Info) (skip, stop bool) {
		if info.Mountpoint == mountpoint {
			skip = false
			stop = true
		} else {
			skip = true
		}
		return
	})
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	if len(mounts) != 1 {
		err = fmt.Errorf("there're more or less than one rootfs mount point: %+v", mounts)
		return
	}
	info = mounts[0]
	return
}
