package runtime

import (
	"github.com/ctrsploit/ctrsploit/prerequisite/apparmor"
	"github.com/ctrsploit/ctrsploit/prerequisite/cgroups"
	"github.com/ctrsploit/ctrsploit/prerequisite/file"
	"github.com/ctrsploit/ctrsploit/prerequisite/mount"
	"github.com/ctrsploit/ctrsploit/prerequisite/mount/mountinfo"
	"github.com/ctrsploit/ctrsploit/prerequisite/proc/net"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
	"github.com/ssst0n3/awesome_libs/awesome_error"
)

type Runtime struct {
	Prerequisites prerequisite.Set
}

func NewRuntime(prerequisites prerequisite.Set) *Runtime {
	r := &Runtime{
		Prerequisites: prerequisites,
	}
	_, err := r.Prerequisites.Check()
	if err != nil {
		awesome_error.CheckWarning(err)
	}
	return r
}

func (r *Runtime) Is() (bool, error) {
	return r.Prerequisites.Check()
}

func Docker() *Runtime {
	return NewRuntime(prerequisite.Or(
		&file.DockerEnvFileExists,
		&mount.RootMountInfoSourceContainsDocker,
		&mount.RootMountInfoVFSOptionsContainsDocker,
		&mountinfo.HostsRootContainsDocker,
		&cgroups.ContainsDocker,
		apparmor.ProfileDockerDefault,
		&net.UnixContainsDockerSock,
		// TODO: image volume's mountinfo contains 'docker'
	))
}

func Containerd() *Runtime {
	return NewRuntime(prerequisite.Or(
		&mount.RootMountInfoVFSOptionsContainsContainerd,
		&mountinfo.HostnameRootContainsContainerd,
		apparmor.ProfileCriContainerd,
		// ctrsploit assumes that nerdctl only appears with containerd
		apparmor.ProfileNerdctlDefault,
		&net.UnixContainsContainerdSock,
	))
}

func Nerdctl() *Runtime {
	return NewRuntime(prerequisite.Or(
		apparmor.ProfileNerdctlDefault),
	)
}
