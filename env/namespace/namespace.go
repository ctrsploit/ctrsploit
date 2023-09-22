package namespace

import (
	"fmt"
	"github.com/ctrsploit/ctrsploit/internal/colorful"
	"github.com/ctrsploit/ctrsploit/internal/log"
	"github.com/ctrsploit/ctrsploit/pkg/namespace"
	"github.com/ctrsploit/ctrsploit/prerequisite/kernel"
	"github.com/ssst0n3/awesome_libs"
)

const CommandName = "namespace"

var (
	tplNamespace = `
===========namespace level===========
{.out}`
)

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
		out := ""
		for _, name := range names {
			level := namespaceLevels[name]
			out += outputNamespaceLevelColorfully(name, level, true)
		}
		info := awesome_libs.Format(tplNamespace, awesome_libs.Dict{
			"out": out,
		})
		print(info)
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
		fmt.Println()
	}
	return
}

func outputNamespaceLevelColorfully(name string, level namespace.Level, padding bool) (info string) {
	var out string
	if level == namespace.LevelHost {
		out = colorful.Danger(level.String())
	} else {
		out = colorful.Safe(level.String())
	}
	if padding {
		info = fmt.Sprintf("%-20s %s\n", name+":", out)
	} else {
		info = fmt.Sprintf("%s: %s\n", name, out)
	}
	return
}

func OutputNamespaceLevelColorfully(name string, level namespace.Level, padding bool) {
	print(outputNamespaceLevelColorfully(name, level, padding))
}
