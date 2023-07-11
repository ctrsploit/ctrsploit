package namespace

import (
	"github.com/ctrsploit/ctrsploit/log"
	"github.com/ctrsploit/ctrsploit/pkg/namespace"
)

const CommandName = "namespace"

func CheckCurrentNamespaceIsHost() (err error) {
	arbitrator, err := namespace.NewInoArbitrator()
	if err != nil {
		return
	}
	isHost, err := namespace.CheckNamespaceIsHost(arbitrator)
	if err != nil {
		return
	}
	for name, isHost := range isHost {
		log.Logger.Infof("%s: %t\n", name, isHost)
	}
	return
}
