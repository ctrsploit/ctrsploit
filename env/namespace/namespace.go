package namespace

import (
	"fmt"
	"github.com/ctrsploit/ctrsploit/log"
	"github.com/ctrsploit/ctrsploit/pkg/namespace"
	"github.com/ctrsploit/ctrsploit/prerequisite/kernel"
	"github.com/ctrsploit/ctrsploit/util"
)

const CommandName = "namespace"

func CheckCurrentNamespaceLevel(ns string) (err error) {
	arbitrator, err := namespace.NewInoArbitrator()
	if err != nil {
		return
	}
	namespaceLevels, names, err := namespace.CheckNamespaceLevel(arbitrator)
	if err != nil {
		return
	}
	if ns == "" {
		fmt.Printf("========namespace level=======\n")
		for _, name := range names {
			level := namespaceLevels[name]
			OutputNamespaceLevelColorfully(name, level, true)
		}
	} else {
		level, ok := namespaceLevels[ns]
		if !ok {
			// maybe kernel not support
			switch ns {
			case namespace.NameTime, namespace.NameTimeForChildren:
				err := kernel.SupportsTimeNamespace.Check()
				if err != nil {
					break
				}
				if !kernel.SupportsTimeNamespace.Satisfied {
					level = namespace.LevelHost
				}
			case namespace.NameCGroup:
				err := kernel.SupportsCgroupNamespace.Check()
				if err != nil {
					break
				}
				if !kernel.SupportsCgroupNamespace.Satisfied {
					level = namespace.LevelHost
				}
			}
		}
		log.Logger.Debugf("%s: %+v \n", ns, level)
		OutputNamespaceLevelColorfully(ns, level, false)
	}
	fmt.Println()
	return
}

func OutputNamespaceLevelColorfully(name string, level namespace.Level, padding bool) {
	var out string
	if level == namespace.LevelHost {
		out = util.Danger(level.String())
	} else {
		out = util.Success(level.String())
	}
	if padding {
		fmt.Printf("%-20s %s\n", name+":", out)
	} else {
		fmt.Printf("%s: %s\n", name, out)
	}
}
