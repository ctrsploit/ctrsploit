package where

import (
	"fmt"

	"github.com/ctrsploit/ctrsploit/prerequisite/apparmor"
	"github.com/ctrsploit/ctrsploit/prerequisite/cgroups"
	"github.com/ctrsploit/ctrsploit/prerequisite/env"
	"github.com/ctrsploit/ctrsploit/prerequisite/file"
	"github.com/ctrsploit/ctrsploit/prerequisite/hostname"
	"github.com/ctrsploit/ctrsploit/prerequisite/mount/mountinfo/root"
	"github.com/ctrsploit/ctrsploit/prerequisite/mount/mountinfo/source"
	"github.com/ctrsploit/ctrsploit/prerequisite/proc/net"
	"github.com/ctrsploit/sploit-spec/pkg/env/container"
	"github.com/ctrsploit/sploit-spec/pkg/printer"
	"github.com/ctrsploit/sploit-spec/pkg/result"
	"github.com/ctrsploit/sploit-spec/pkg/result/item"
)

type Result []struct {
	Name  result.Title `json:"name"`
	Rules []item.Bool  `json:"rules"`
	In    item.Bool    `json:"in"`
}

func Human(machine container.Where) (human Result) {
	human = Result{
		{
			Name: result.Title{
				Name: "Container",
			},
			In: item.Bool{
				Name:        "Is in Container",
				Description: "",
				Result:      machine.Container.In,
			},
		},
		{
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
					Name:        source.RootMountInfoSourceContainsDocker.Name,
					Description: source.RootMountInfoSourceContainsDocker.Info,
					Result:      machine.Docker.Rules[source.RootMountInfoSourceContainsDocker.Name],
				},
				{
					Name:        source.RootMountInfoVFSOptionsContainsDocker.Name,
					Description: source.RootMountInfoVFSOptionsContainsDocker.Info,
					Result:      machine.Docker.Rules[source.RootMountInfoVFSOptionsContainsDocker.Name],
				},
				{
					Name:        root.HostsRootContainsDocker.Name,
					Description: root.HostsRootContainsDocker.Info,
					Result:      machine.Docker.Rules[root.HostsRootContainsDocker.Name],
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
		{
			Name: result.Title{
				Name: "containerd",
			},
			Rules: []item.Bool{
				{
					Name:        source.RootMountInfoVFSOptionsContainsContainerd.Name,
					Description: source.RootMountInfoVFSOptionsContainsContainerd.Info,
					Result:      machine.Containerd.Rules[source.RootMountInfoVFSOptionsContainsContainerd.Name],
				},
				{
					Name:        root.HostnameRootContainsContainerd.Name,
					Description: root.HostnameRootContainsContainerd.Info,
					Result:      machine.Containerd.Rules[root.HostnameRootContainsContainerd.Name],
				},
				{
					Name:        root.HostnameRootContainsNerdctl.Name,
					Description: root.HostnameRootContainsNerdctl.Info,
					Result:      machine.Containerd.Rules[root.HostnameRootContainsNerdctl.Name],
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
		{
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
					Name:        root.HostsRootContainsPods.Name,
					Description: root.HostsRootContainsPods.Info,
					Result:      machine.K8s.Rules[root.HostsRootContainsPods.Name],
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
				{
					Name:        env.KubernetesServiceHostExists.Name,
					Description: env.KubernetesServiceHostExists.Info,
					Result:      machine.K8s.Rules[env.KubernetesServiceHostExists.Name],
				},
			},
			In: item.Bool{
				Name:        "is in k8s",
				Description: "",
				Result:      machine.K8s.In,
			},
		},
		{
			Name: result.Title{
				Name: "Nerdctl",
			},
			Rules: []item.Bool{
				{
					Name:        apparmor.ProfileNerdctlDefault.Name,
					Description: apparmor.ProfileNerdctlDefault.Info,
					Result:      machine.Nerdctl.Rules[apparmor.ProfileNerdctlDefault.Name],
				},
				{
					Name:        file.HostsContainsNerdctlMarker.Name,
					Description: file.HostsContainsNerdctlMarker.Info,
					Result:      machine.Nerdctl.Rules[file.HostsContainsNerdctlMarker.Name],
				},
				{
					Name:        root.HostnameRootContainsNerdctl.Name,
					Description: root.HostnameRootContainsNerdctl.Info,
					Result:      machine.Nerdctl.Rules[root.HostnameRootContainsNerdctl.Name],
				},
			},
			In: item.Bool{
				Name:        "is in nerdctl",
				Description: "",
				Result:      machine.Nerdctl.In,
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
