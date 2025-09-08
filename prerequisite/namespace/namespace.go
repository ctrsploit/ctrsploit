package namespace

import (
	"fmt"

	"github.com/ctrsploit/ctrsploit/pkg/namespace"
	"github.com/ctrsploit/sploit-spec/pkg/env/container"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
	"github.com/ssst0n3/awesome_libs/awesome_error"
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
)

func (p *Namespace) Check() (bool, error) {
	if p.Checked {
		return p.Satisfied, nil
	}
	arbitrator, err := namespace.NewInoArbitrator()
	if err != nil {
		return false, err
	}
	namespaceLevels, _, err := namespace.CheckNamespaceLevel(arbitrator)
	if err != nil {
		return false, err
	}
	level, ok := namespaceLevels[container.NamespaceMapType2Name[p.Type]]
	if !ok {
		err = fmt.Errorf("unknown namespace type %s", p.Type)
		awesome_error.CheckErr(err)
		return false, err
	}
	p.Satisfied = level == container.NamespaceLevelHost
	p.Checked = true
	return p.Satisfied, nil
}
