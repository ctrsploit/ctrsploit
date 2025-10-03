package runtime

import (
	"github.com/ctrsploit/ctrsploit/prerequisite/apparmor"
	"github.com/ctrsploit/ctrsploit/prerequisite/cgroups"
	"github.com/ctrsploit/ctrsploit/prerequisite/file"
	"github.com/ctrsploit/ctrsploit/prerequisite/mount"
	"github.com/ctrsploit/ctrsploit/prerequisite/proc/net"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
	"github.com/ssst0n3/awesome_libs/awesome_error"
)

type Docker struct {
	Prerequisites prerequisite.Set
}

func NewDocker() *Docker {
	d := &Docker{
		Prerequisites: prerequisite.Or(
			&file.DockerEnvFileExists,
			&mount.RootMountInfoSourceContainsDocker,
			&mount.RootMountInfoVFSOptionsContainsDocker,
			&mount.HostsMountInfoRootContainsDocker,
			&cgroups.ContainsDocker,
			apparmor.ProfileDockerDefault,
			&net.UnixContainsDockerSock,
		),
	}
	_, err := d.Prerequisites.Check()
	if err != nil {
		awesome_error.CheckWarning(err)
	}
	return d
}

func (r *Docker) Is() (bool, error) {
	return r.Prerequisites.Check()
}
