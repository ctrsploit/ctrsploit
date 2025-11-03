package namespace

import (
	"fmt"

	"github.com/ctrsploit/ctrsploit/pkg/namespace"
	"github.com/ctrsploit/sploit-spec/pkg/env/container"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
)

type Namespace struct {
	ExpectedLevel container.NamespaceLevel
	Type          container.NamespaceType
	prerequisite.BasePrerequisite
}

var (
	NetworkNamespaceLevelHost = Namespace{
		ExpectedLevel: container.NamespaceLevelHost,
		Type:          container.NamespaceTypeNetwork,
		BasePrerequisite: prerequisite.BasePrerequisite{
			Name:   "Network_Namespace_Level_Host",
			Info:   "Container with host network namespace can cause network-based attacks even escape",
			ExeEnv: exeenv.InContainer,
		},
	}
	PidNamespaceLevelHost = Namespace{
		ExpectedLevel: container.NamespaceLevelHost,
		Type:          container.NamespaceTypePid,
		BasePrerequisite: prerequisite.BasePrerequisite{
			Name:   "Pid_Namespace_Level_Host",
			Info:   "Container with host pid namespace may cause cross filesystem access even escape",
			ExeEnv: exeenv.InContainer,
		},
	}
)

func (p *Namespace) Check() (bool, error) {
	return p.CheckTemplate(func() {
		arbitrator, err := namespace.NewInoArbitrator()
		if err != nil {
			p.Err = p.WrapErr(fmt.Errorf("creating arbitrator: %w", err))
			return
		}
		namespaceLevels, _, err := namespace.CheckNamespaceLevel(arbitrator)
		if err != nil {
			p.Err = p.WrapErr(fmt.Errorf("checking level: %w", err))
			return
		}
		level, ok := namespaceLevels[container.NamespaceMapType2Name[p.Type]]
		if !ok {
			p.Err = p.WrapErr(fmt.Errorf("unknown namespace type %s", p.Type))
			return
		}
		p.Satisfied = level == container.NamespaceLevelHost
		return
	})
}
