package runtime

import (
	"github.com/ctrsploit/ctrsploit/prerequisite/apparmor"
	"github.com/ctrsploit/ctrsploit/prerequisite/cgroups"
	"github.com/ctrsploit/ctrsploit/prerequisite/env"
	"github.com/ctrsploit/ctrsploit/prerequisite/file"
	"github.com/ctrsploit/ctrsploit/prerequisite/hostname"
	"github.com/ctrsploit/ctrsploit/prerequisite/mount/mountinfo/root"
	"github.com/ctrsploit/ctrsploit/prerequisite/mount/mountinfo/source"
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
		&source.RootMountInfoSourceContainsDocker,
		&source.RootMountInfoVFSOptionsContainsDocker,
		&root.HostsRootContainsDocker,
		&cgroups.ContainsDocker,
		apparmor.ProfileDockerDefault,
		&net.UnixContainsDockerSock,
		// TODO: image volume's mountinfo contains 'docker'
	))
}

func Containerd() *Runtime {
	return NewRuntime(prerequisite.Or(
		&source.RootMountInfoVFSOptionsContainsContainerd,
		// ctrsploit assumes that nerdctl only appears with containerd
		&root.HostnameRootContainsContainerd, &root.HostnameRootContainsNerdctl,
		&net.UnixContainsContainerdSock,
		// ctrsploit assumes that nerdctl only appears with containerd
		apparmor.ProfileCriContainerd, apparmor.ProfileNerdctlDefault,
	))
}

func Nerdctl() *Runtime {
	return NewRuntime(prerequisite.Or(
		apparmor.ProfileNerdctlDefault,
		&file.HostsContainsNerdctlMarker,
		&root.HostnameRootContainsNerdctl,
	))
}

func K8s() *Runtime {
	return NewRuntime(prerequisite.Or(
		&file.K8sSecretsExists,
		&root.HostsRootContainsPods,
		&cgroups.ContainsKubepods,
		&hostname.K8sDeploymentHostname,
		&env.KubernetesServiceHostExists,
	))
}
