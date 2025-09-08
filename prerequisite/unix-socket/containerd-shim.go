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

func (p *Available) Check() (bool, error) {
	if p.Checked {
		return p.Satisfied, nil
	}
	path, err := net.ContainerdShimAbstractUnixSocketPath(p.PrefixSocketName)
	if err != nil {
		return false, err
	}
	p.Satisfied = path != ""
	p.Checked = true
	return p.Satisfied, nil
}
