package where

import (
	"fmt"
	"github.com/ctrsploit/ctrsploit/pkg/where"
	"github.com/ctrsploit/sploit-spec/pkg/printer"
	"github.com/ctrsploit/sploit-spec/pkg/result"
	"github.com/ctrsploit/sploit-spec/pkg/result/item"
)

const CommandName = "where"

type Result struct {
	Name  result.Title `json:"name"`
	Rules []item.Bool  `json:"rules"`
	In    item.Bool    `json:"in"`
}

func Where() (err error) {
	r := map[string]Result{
		"container": Container(),
		"docker":    Docker(),
		"k8s":       K8s(),
	}
	fmt.Println(printer.Printer.Print(r))
	return
}

func Container() (r Result) {
	c := where.Container{}
	in, err := c.IsIn()
	if err != nil {
		return
	}

	r = Result{
		Name: result.Title{
			Name: "Container",
		},
		In: item.Bool{
			Name:        "Is in Container",
			Description: "",
			Result:      in,
		},
	}
	return
}

func Docker() (r Result) {
	d := where.Docker{}
	in, err := d.IsIn()
	if err != nil {
		return
	}
	r = Result{
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
	return
}

func K8s() (r Result) {
	k := where.K8s{}
	in, err := k.IsIn()
	if err != nil {
		return
	}

	r = Result{
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
	return
}
