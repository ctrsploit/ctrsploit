package root

import (
	"fmt"
	"os"
	"strings"

	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
	"github.com/moby/sys/mountinfo"
	"github.com/ssst0n3/awesome_libs/awesome_error"
)

type ContainsUncertainMountPointWithType struct {
	prerequisite.BasePrerequisite
	ExpectedContains string
	Type             os.FileMode
	realMountPoint   string
}

func (p *ContainsUncertainMountPointWithType) RealMountPoint() string {
	return p.realMountPoint
}

func (p *ContainsUncertainMountPointWithType) Check() (bool, error) {
	return p.CheckTemplate(func() (bool, error) {
		infos, err := mountinfo.GetMounts(func(info *mountinfo.Info) (skip, stop bool) {
			skip = !strings.Contains(info.Root, p.ExpectedContains)
			return skip, false
		})
		if err != nil {
			p.Err = fmt.Errorf("failed to check [%s], caused by getting mountinfo: %w", p.GetName(), err)
			return p.Satisfied, p.Err
		}
		for _, info := range infos {
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
		return p.Satisfied, p.Err
	})
}

var DockerSock = ContainsUncertainMountPointWithType{
	ExpectedContains: "docker.sock",
	Type:             os.ModeSocket,
}
