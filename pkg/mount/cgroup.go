package mount

import (
	"fmt"

	v1 "github.com/ctrsploit/ctrsploit/pkg/cgroup/v1"
	"github.com/ctrsploit/sploit-spec/pkg/log"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"golang.org/x/sys/unix"
)

var ErrorNoTopLevelCgroupSubSystems = fmt.Errorf("no top level cgroup sub systems")

func CgroupV1(dest string, options string) (err error) {
	log.Logger.Infof("mount cgroup/%s to %s", options, dest)
	err = unix.Mount("cgroup", dest, "cgroup", 0, options)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	return
}

// https://github.com/torvalds/linux/blob/v5.4/net/core/netprio_cgroup.c#L263
func mountable(name string) {

}

func TopLevelCgroupSubSystem(dest string) error {
	systems, err := v1.ListTopLevelSubSystem()
	if err != nil {
		return fmt.Errorf("mount.ListTopLevelSubSystem: %w", err)
	}
	if len(systems) == 0 {
		return ErrorNoTopLevelCgroupSubSystems
	}
	err = CgroupV1(dest, systems[0])
	if err != nil {
		return err
	}
	return nil
}

func Unmount(dest string) (err error) {
	err = unix.Unmount(dest, 0)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	return
}
