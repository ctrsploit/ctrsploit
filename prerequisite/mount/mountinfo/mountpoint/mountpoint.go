package mountpoint

import (
	"fmt"
	"os"
	"strings"

	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
	"github.com/moby/sys/mountinfo"
	"github.com/ssst0n3/awesome_libs/awesome_error"
)

type ContainsMountPoint struct {
	prerequisite.BasePrerequisite
	ExpectedContains string
	Type             os.FileMode
	realMountPoint   string
}

func (p *ContainsMountPoint) RealMountPoint() string {
	return p.realMountPoint
}

func (p *ContainsMountPoint) Check() (bool, error) {
	return p.CheckTemplate(func() {
		infos, err := mountinfo.GetMounts(nil)
		if err != nil {
			p.Err = p.WrapErr(fmt.Errorf("checking mount points: %w", err))
			return
		}
		for _, info := range infos {
			if strings.Contains(info.Mountpoint, p.ExpectedContains) {
				if p.Type == 0 {
					p.realMountPoint = info.Mountpoint
					p.Satisfied = true
					break
				}
				fi, err := os.Lstat(info.Mountpoint)
				if err != nil {
					awesome_error.CheckWarning(err)
					continue
				}
				// check file type is expected type
				if fi.Mode()&p.Type != 0 {
					p.realMountPoint = info.Mountpoint
					p.Satisfied = true
					break
				}
			}
		}
		return
	})
}

var DockerSock = ContainsMountPoint{
	ExpectedContains: "docker.sock",
	Type:             os.ModeSocket,
}
