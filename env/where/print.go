package where

import (
	"fmt"

	"github.com/ctrsploit/ctrsploit/prerequisite/apparmor"
	"github.com/ctrsploit/ctrsploit/prerequisite/cgroups"
	"github.com/ctrsploit/ctrsploit/prerequisite/file"
	"github.com/ctrsploit/ctrsploit/prerequisite/hostname"
	"github.com/ctrsploit/ctrsploit/prerequisite/mount"
	"github.com/ctrsploit/ctrsploit/prerequisite/mount/mountinfo"
	"github.com/ctrsploit/ctrsploit/prerequisite/proc/net"
	"github.com/ctrsploit/sploit-spec/pkg/env/container"
	"github.com/ctrsploit/sploit-spec/pkg/printer"
	"github.com/ctrsploit/sploit-spec/pkg/result"
	"github.com/ctrsploit/sploit-spec/pkg/result/item"
)

type Result map[string]struct {
	Name  result.Title `json:"name"`
	Rules []item.Bool  `json:"rules"`
	In    item.Bool    `json:"in"`
}

func Human(machine container.Where) (human Result) {
	human = Result{
		"container": {
			Name: result.Title{
				Name: "Container",
			},
			In: item.Bool{
				Name:        "Is in Container",
				Description: "",
				Result:      machine.Container.In,
			},
		},
		"docker": {
			Name: result.Title{
				Name: "Docker",
			},
			Rules: []item.Bool{
				{
					Name:        file.DockerEnvFileExists.Name,
					Description: file.DockerEnvFileExists.Info,
					Result:      machine.Docker.Rules[file.DockerEnvFileExists.Name],
				},
				{
					Name:        mount.RootMountInfoSourceContainsDocker.Name,
					Description: mount.RootMountInfoSourceContainsDocker.Info,
					Result:      machine.Docker.Rules[mount.RootMountInfoSourceContainsDocker.Name],
				},
				{
					Name:        mount.RootMountInfoVFSOptionsContainsDocker.Name,
					Description: mount.RootMountInfoVFSOptionsContainsDocker.Info,
					Result:      machine.Docker.Rules[mount.RootMountInfoVFSOptionsContainsDocker.Name],
				},
				{
					Name:        mountinfo.HostsRootContainsDocker.Name,
					Description: mountinfo.HostsRootContainsDocker.Info,
					Result:      machine.Docker.Rules[mountinfo.HostsRootContainsDocker.Name],
				},
				{
					Name:        cgroups.ContainsDocker.Name,
					Description: cgroups.ContainsDocker.Info,
					Result:      machine.Docker.Rules[cgroups.ContainsDocker.Name],
				},
				{
					Name:        apparmor.ProfileDockerDefault.Name,
					Description: apparmor.ProfileDockerDefault.Info,
					Result:      machine.Docker.Rules[apparmor.ProfileDockerDefault.Name],
				},
				{
					Name:        net.UnixContainsDockerSock.Name,
					Description: net.UnixContainsDockerSock.Info,
					Result:      machine.Docker.Rules[net.UnixContainsDockerSock.Name],
				},
			},
			In: item.Bool{
				Name:        "Is in docker",
				Description: "",
				Result:      machine.Docker.In,
			},
		},
		"containerd": {
			Name: result.Title{
				Name: "containerd",
			},
			Rules: []item.Bool{
				{
					Name:        mount.RootMountInfoVFSOptionsContainsContainerd.Name,
					Description: mount.RootMountInfoVFSOptionsContainsContainerd.Info,
					Result:      machine.Containerd.Rules[mount.RootMountInfoVFSOptionsContainsContainerd.Name],
				},
				{
					Name:        mountinfo.HostnameRootContainsContainerd.Name,
					Description: mountinfo.HostnameRootContainsContainerd.Info,
					Result:      machine.Containerd.Rules[mountinfo.HostnameRootContainsContainerd.Name],
				},
				{
					Name:        mountinfo.HostnameRootContainsNerdctl.Name,
					Description: mountinfo.HostnameRootContainsNerdctl.Info,
					Result:      machine.Containerd.Rules[mountinfo.HostnameRootContainsNerdctl.Name],
				},
				{
					Name:        net.UnixContainsContainerdSock.Name,
					Description: net.UnixContainsContainerdSock.Info,
					Result:      machine.Containerd.Rules[net.UnixContainsContainerdSock.Name],
				},
				{
					Name:        apparmor.ProfileCriContainerd.Name,
					Description: apparmor.ProfileCriContainerd.Info,
					Result:      machine.Containerd.Rules[apparmor.ProfileCriContainerd.Name],
				},
				{
					Name:        apparmor.ProfileNerdctlDefault.Name,
					Description: apparmor.ProfileNerdctlDefault.Info,
					Result:      machine.Containerd.Rules[apparmor.ProfileNerdctlDefault.Name],
				},
			},
			In: item.Bool{
				Name:        "Is in containerd",
				Description: "",
				Result:      machine.Containerd.In,
			},
		},
		"k8s": {
			Name: result.Title{
				Name: "K8S",
			},
			Rules: []item.Bool{
				{
					Name:        file.K8sSecretsExists.Name,
					Description: file.K8sSecretsExists.Info,
					Result:      machine.K8s.Rules[file.K8sSecretsExists.Name],
				},
				{
					Name:        mountinfo.HostsRootContainsPods.Name,
					Description: mountinfo.HostsRootContainsPods.Info,
					Result:      machine.K8s.Rules[mountinfo.HostsRootContainsPods.Name],
				},
				{
					Name:        cgroups.ContainsDocker.Name,
					Description: cgroups.ContainsDocker.Info,
					Result:      machine.K8s.Rules[cgroups.ContainsDocker.Name],
				},
				{
					Name:        hostname.K8sDeploymentHostname.Name,
					Description: hostname.K8sDeploymentHostname.Info,
					Result:      machine.K8s.Rules[hostname.K8sDeploymentHostname.Name],
				},
			},
			In: item.Bool{
				Name:        "is in k8s",
				Description: "",
				Result:      machine.K8s.In,
			},
		},
	}
	return
}

func Print() (err error) {
	machine, err := Where()
	if err != nil {
		return
	}
	u := result.Union{
		Machine: machine,
		Human:   Human(machine),
	}
	fmt.Println(printer.Printer.Print(u))
	return
}
