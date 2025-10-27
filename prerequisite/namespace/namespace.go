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
	return p.CheckTemplate(func() (bool, error) {
		arbitrator, err := namespace.NewInoArbitrator()
		if err != nil {
			p.Err = fmt.Errorf("failed to check [%s], caused by creating arbitrator: %w", p.GetName(), err)
			return p.Satisfied, p.Err
		}
		namespaceLevels, _, err := namespace.CheckNamespaceLevel(arbitrator)
		if err != nil {
			p.Err = fmt.Errorf("failed to check [%s], caused by checking level: %w", p.GetName(), err)
			return p.Satisfied, p.Err
		}
		level, ok := namespaceLevels[container.NamespaceMapType2Name[p.Type]]
		if !ok {
			p.Err = fmt.Errorf("failed to check [%s], caused by unknown namespace type %s", p.GetName(), p.Type)
			return p.Satisfied, p.Err
		}
		p.Satisfied = level == container.NamespaceLevelHost
		return p.Satisfied, p.Err
	})
}
