package where

import (
	"fmt"
	"github.com/ctrsploit/ctrsploit/internal"
	"github.com/ctrsploit/ctrsploit/pkg/where"
	"github.com/ctrsploit/sploit-spec/pkg/result"
	"github.com/ctrsploit/sploit-spec/pkg/result/item"
)

const CommandName = "where"

type Result struct {
	Name  result.Title
	Rules []item.Bool `json:"rules"`
	In    item.Bool   `json:"in"`
}

func (r Result) String() (s string) {
	s += internal.Print(r.Name)
	for _, r := range r.Rules {
		s += internal.Print(r)
	}
	s += internal.Print(r.In)
	return
}

func Container() (err error) {
	c := where.Container{}
	in, err := c.IsIn()
	if err != nil {
		return
	}

	r := Result{
		Name: result.Title{
			Name: "Container",
		},
		In: item.Bool{
			Name:        "Is in Container",
			Description: "",
			Result:      in,
		},
	}
	fmt.Println(r)
	return
}

func Docker() (err error) {
	d := where.Docker{}
	in, err := d.IsIn()
	if err != nil {
		return
	}
	r := Result{
		Name: result.Title{
			Name: "Docker",
		},
		Rules: []item.Bool{
			{
				Name:        "dockerenv",
				Description: ".dockerenv exists",
				Result:      d.DockerEnvFileExists,
			},
			{
				Name:        "rootfs",
				Description: "rootfs contains 'docker'",
				Result:      d.RootfsContainsDocker,
			},
			{
				Name:        "cgroups",
				Description: "cgroups contains 'docker'",
				Result:      d.CgroupContainsDocker,
			},
			{
				Name:        "hosts",
				Description: "the mount source of /etc/hosts contains 'docker'",
				Result:      d.HostsMountSourceContainsDocker,
			},
			{
				Name:        "hostname",
				Description: "hostname match regex ^[0-9a-f]{12}$",
				Result:      d.HostnameMatchPattern,
			},
		},
		In: item.Bool{
			Name:        "Is in docker",
			Description: "",
			Result:      in,
		},
	}
	fmt.Println(r)
	return
}

func K8s() (err error) {
	k := where.K8s{}
	in, err := k.IsIn()
	if err != nil {
		return
	}

	r := Result{
		Name: result.Title{
			Name: "K8S",
		},
		Rules: []item.Bool{
			{
				Name:        "secret",
				Description: fmt.Sprintf("secret path %s exists", where.PathDirSecrets),
				Result:      k.DirSecretsExists,
			},
			{
				Name:        "hostname",
				Description: "hostname match k8s pattern",
				Result:      k.HostnameMatchPattern,
			},
			{
				Name:        "hosts",
				Description: "the mount source of /etc/hosts contains 'pods'",
				Result:      k.HostsMountSourceContainsPods,
			},
			{
				Name:        "cgroups",
				Description: "cgroups contains 'kubepods'",
				Result:      k.CgroupContainsKubepods,
			},
		},
		In: item.Bool{
			Name:        "is in k8s",
			Description: "",
			Result:      in,
		},
	}
	fmt.Println(r)
	return
}
