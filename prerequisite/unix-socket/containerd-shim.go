package unix_socket

import (
	"github.com/ctrsploit/ctrsploit/pkg/proc/net"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
)

type Available struct {
	PrefixSocketName string
	prerequisite.BasePrerequisite
}

var ContainerdShimAbstract = Available{
	PrefixSocketName: "@/containerd-shim/",
	BasePrerequisite: prerequisite.BasePrerequisite{
		Name: "containerd-shim socket available",
		Info: "abstract unix socket can be seen",
	},
}

func (p *Available) Check() (err error) {
	err = p.BasePrerequisite.Check()
	if err != nil {
		return
	}
	path, err := net.ContainerdShimAbstractUnixSocketPath(p.PrefixSocketName)
	if err != nil {
		return
	}
	p.Satisfied = path != ""
	return
}
