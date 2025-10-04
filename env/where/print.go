package where

import (
	"fmt"

	"github.com/ctrsploit/ctrsploit/pkg/where"
	"github.com/ctrsploit/ctrsploit/prerequisite/apparmor"
	"github.com/ctrsploit/ctrsploit/prerequisite/cgroups"
	"github.com/ctrsploit/ctrsploit/prerequisite/file"
	"github.com/ctrsploit/ctrsploit/prerequisite/mount"
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
					Name:        mount.HostsMountInfoRootContainsDocker.Name,
					Description: mount.HostsMountInfoRootContainsDocker.Info,
					Result:      machine.Docker.Rules[mount.HostsMountInfoRootContainsDocker.Name],
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
					Name:        "rootfs",
					Description: "rootfs contains 'containerd'",
					Result:      machine.Containerd.Rules["rootfs"],
				},
				{
					Name:        "hosts",
					Description: "the mount source of /etc/hosts contains 'nerdctl'",
					Result:      machine.Containerd.Rules["hosts"],
				},
				{
					Name:        "hostname",
					Description: "the mount source of /etc/hostname contains 'containerd'/'nerdctl'",
					Result:      machine.Containerd.Rules["hostname"],
				},
				{
					Name:        "socket",
					Description: "/proc/net/unix contains 'containerd.sock', no 'docker.sock'",
					Result:      machine.Containerd.Rules["socket"],
				},
			},
			In: item.Bool{},
		},
		"k8s": {
			Name: result.Title{
				Name: "K8S",
			},
			Rules: []item.Bool{
				{
					Name:        "secret",
					Description: fmt.Sprintf("secret path %s exists", where.PathDirSecrets),
					Result:      machine.K8s.Rules["secret"],
				},
				{
					Name:        "hostname",
					Description: "hostname match k8s pattern",
					Result:      machine.K8s.Rules["hostname"],
				},
				{
					Name:        "hosts",
					Description: "the mount source of /etc/hosts contains 'pods'",
					Result:      machine.K8s.Rules["hosts"],
				},
				{
					Name:        "cgroups",
					Description: "cgroups contains 'kubepods'",
					Result:      machine.K8s.Rules["cgroups"],
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
