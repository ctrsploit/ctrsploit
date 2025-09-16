package mount

import (
	"os"
	"strings"

	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
	"github.com/moby/sys/mountinfo"
	"github.com/ssst0n3/awesome_libs/awesome_error"
)

type Contains struct {
	prerequisite.BasePrerequisite
	ExpectedMountPoint string
	Type               os.FileMode
	realMountPoint     string
}

func (p *Contains) RealMountPoint() string {
	return p.realMountPoint
}

func (p *Contains) Check() (bool, error) {
	if p.Checked {
		return p.Satisfied, nil
	}
	infos, err := mountinfo.GetMounts(nil)
	if err != nil {
		awesome_error.CheckErr(err)
		return false, err
	}
	for _, info := range infos {
		if strings.Contains(info.Mountpoint, p.ExpectedMountPoint) {
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
			if fi.Mode()&p.Type != 0 {
				p.realMountPoint = info.Mountpoint
				p.Satisfied = true
				break
			}
		}
	}
	p.Checked = true
	return p.Satisfied, nil
}

var DockerSock = Contains{
	ExpectedMountPoint: "docker.sock",
	Type:               os.ModeSocket,
}
