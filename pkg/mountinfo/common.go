package mountinfo

import (
	"fmt"

	"github.com/moby/sys/mountinfo"
	"github.com/ssst0n3/awesome_libs/awesome_error"
)

const ErrMountPointNotFound = "mount point not found"

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
	if len(mounts) == 0 {
		err = fmt.Errorf("%s: %s", ErrMountPointNotFound, mountpoint)
		awesome_error.CheckDebug(err)
		return
	}
	if len(mounts) > 1 {
		err = fmt.Errorf("there're more than one mount point %s: %+v", mountpoint, mounts)
		awesome_error.CheckWarning(err)
		return
	}
	info = mounts[0]
	return
}

func MountInfo() (info []*mountinfo.Info, err error) {
	info, err = mountinfo.GetMounts(nil)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	return
}
