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
	return p.CheckTemplate(func() (bool, error) {
		path, err := net.ContainerdShimAbstractUnixSocketPath(p.PrefixSocketName)
		if err != nil {
			p.Err = fmt.Errorf("failed to check [%s], caused by getting abstract unix socket path: %w", p.GetName(), err)
			return p.Satisfied, p.Err
		}
		p.Satisfied = path != ""
		return p.Satisfied, p.Err
	})
}
