package where

import (
	"fmt"

	"github.com/ctrsploit/ctrsploit/pkg/where"
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
					Name:        "dockerenv",
					Description: ".dockerenv exists",
					Result:      machine.Docker.Rules["dockerenv"],
				},
				{
					Name:        "rootfs",
					Description: "rootfs contains 'docker'",
					Result:      machine.Docker.Rules["rootfs"],
				},
				{
					Name:        "cgroups",
					Description: "cgroups contains 'docker'",
					Result:      machine.Docker.Rules["cgroups"],
				},
				{
					Name:        "hosts",
					Description: "the mount source of /etc/hosts contains 'docker'",
					Result:      machine.Docker.Rules["hosts"],
				},
				//{
				//	Name:        "hostname",
				//	Description: "hostname match regex ^[0-9a-f]{12}$",
				//	Result:      machine.Docker.Rules["hostname"],
				//},
				{
					Name:        "lsm",
					Description: "/proc/1/attr/current contains 'docker'",
					Result:      machine.Docker.Rules["lsm"],
				},
				{
					Name:        "socket",
					Description: "/proc/net/unix contains 'docker.sock'",
					Result:      machine.Docker.Rules["socket"],
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
