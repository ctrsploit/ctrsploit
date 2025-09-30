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
			"lsm":       d.ProcAttrCurrentContainsDocker,
			"socket":    d.ProcNetUnixContainsDockerSock,
		},
	}
	return
}

func Containerd() (containerd container.Type, err error) {
	r := runtime.NewContainerd()
	in, _ := r.Is()
	containerd = container.Type{
		In: in,
		Rules: map[string]bool{
			"rootfs":   r.RootfsContainsContainerd,
			"hosts":    r.HostsMountSourceContainsNerdctl,
			"hostname": r.HostnameMountSourceContainsContainerd,
			"socket":   r.ProcNetUnixContainsContainerdSock,
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
	containerd, err := Containerd()
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
		Containerd: containerd,
		Docker:     docker,
	}
	return
}
