package namespace

import (
	"fmt"
	"github.com/ctrsploit/ctrsploit/log"
	"github.com/ctrsploit/ctrsploit/pkg/namespace"
	"github.com/ctrsploit/ctrsploit/prerequisite"
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
				err := prerequisite.KernelSupportsTimeNamespace.Check()
				if err != nil {
					break
				}
				if !prerequisite.KernelSupportsTimeNamespace.Satisfied {
					level = namespace.TypeNamespaceLevelHost
				}
			case namespace.NameCGroup:
				err := prerequisite.KernelSupportsCgroupNamespace.Check()
				if err != nil {
					break
				}
				if !prerequisite.KernelSupportsCgroupNamespace.Satisfied {
					level = namespace.TypeNamespaceLevelHost
				}
			}
		}
		log.Logger.Debugf("%s: %+v \n", ns, level)
		OutputNamespaceLevelColorfully(ns, level, false)
	}
	return
}

func CheckNamespaceValid(ns string) (valid bool) {
	_, valid = namespace.TypeMap[ns]
	return
}

func OutputNamespaceLevelColorfully(name string, level namespace.Level, padding bool) {
	var out string
	if level == namespace.TypeNamespaceLevelHost {
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
