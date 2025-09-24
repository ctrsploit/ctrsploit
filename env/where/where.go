package where

import (
	"github.com/ctrsploit/ctrsploit/pkg/runtime"
	"github.com/ctrsploit/ctrsploit/pkg/where"
	"github.com/ctrsploit/sploit-spec/pkg/env/container"
)

const CommandName = "where"

func Docker() (docker container.Type, err error) {
	d := runtime.NewDocker()
	in, _ := d.Is()
	docker = container.Type{
		In: in,
		Rules: map[string]bool{
			"dockerenv": d.DockerEnvFileExists,
			"rootfs":    d.RootfsContainsDocker,
			"cgroups":   d.CgroupContainsDocker,
			"hosts":     d.HostsMountSourceContainsDocker,
			"hostname":  d.HostnameMatchPattern,
		},
	}
	return
}

func K8s() (k8s container.Type, err error) {
	k := where.K8s{}
	in, err := k.IsIn()
	if err != nil {
		return
	}
	k8s = container.Type{
		In: in,
		Rules: map[string]bool{
			"secret":   k.DirSecretsExists,
			"hostname": k.HostnameMatchPattern,
			"hosts":    k.HostsMountSourceContainsPods,
			"cgroups":  k.CgroupContainsKubepods,
		},
	}
	return
}

func Container() (t container.Type, err error) {
	c := where.Container{}
	in, err := c.IsIn()
	if err != nil {
		return
	}
	t = container.Type{
		In:    in,
		Rules: map[string]bool{},
	}
	return
}

func Where() (machine container.Where, err error) {
	docker, err := Docker()
	if err != nil {
		return
	}
	k8s, err := K8s()
	if err != nil {
		return
	}
	c, err := Container()
	if err != nil {
		return
	}
	machine = container.Where{
		Container:  c,
		K8s:        k8s,
		Containerd: container.Type{},
		Docker:     docker,
	}
	return
}
