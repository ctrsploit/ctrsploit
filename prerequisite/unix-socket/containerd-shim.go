package unix_socket

import (
	"fmt"

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
	return p.CheckTemplate(func() {
		path, err := net.ContainerdShimAbstractUnixSocketPath(p.PrefixSocketName)
		if err != nil {
			p.Err = p.WrapErr(fmt.Errorf("getting abstract unix socket path: %w", err))
			return
		}
		p.Satisfied = path != ""
		return
	})
}
