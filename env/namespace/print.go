package namespace

import (
	"fmt"

	"github.com/ctrsploit/ctrsploit/prerequisite/kernel"
	"github.com/ctrsploit/sploit-spec/pkg/colorful"
	"github.com/ctrsploit/sploit-spec/pkg/env/container"
	"github.com/ctrsploit/sploit-spec/pkg/log"
	"github.com/ctrsploit/sploit-spec/pkg/printer"
	"github.com/ctrsploit/sploit-spec/pkg/result"
	"github.com/ctrsploit/sploit-spec/pkg/result/item"
)

type Result struct {
	Name   result.Title `json:"name"`
	Levels []item.Short `json:"levels"`
}

func level2result(name string, level container.NamespaceLevel) item.Short {
	levelResult := level.String()
	if levelResult == "host" {
		levelResult = colorful.O.Danger(levelResult)
	}
	return item.Short{
		Name:        name,
		Description: "",
		Result:      levelResult,
	}
}

func Human(namespace container.Namespace, ns string) (human Result) {
	human = Result{
		Name: result.Title{
			Name: "Namespace Level",
		},
		Levels: []item.Short{},
	}
	if ns == "" {
		for _, name := range namespace.Names {
			level := namespace.Levels[name]
			human.Levels = append(human.Levels, level2result(name, level))
		}
	} else {
		level, ok := namespace.Levels[ns]
		if !ok {
			// maybe kernel not support
			switch ns {
			case container.NamespaceNameTime, container.NamespaceNameTimeForChildren:
				_, err := kernel.SupportsTimeNamespace.Check()
				if err != nil {
					break
				}
				if !kernel.SupportsTimeNamespace.Satisfied {
					level = container.NamespaceLevelHost
				}
			case container.NamespaceNameCGroup:
				_, err := kernel.SupportsCgroupNamespace.Check()
				if err != nil {
					break
				}
				if !kernel.SupportsCgroupNamespace.Satisfied {
					level = container.NamespaceLevelHost
				}
			}
		}
		log.Logger.Debugf("%s: %+v \n", ns, level)
		human.Levels = append(human.Levels, level2result(ns, level))
	}
	return
}

func Print(ns string) (err error) {
	machine, err := Namespace()
	if err != nil {
		return
	}
	u := result.Union{
		Machine: machine,
		Human:   Human(machine, ns),
	}
	fmt.Println(printer.Printer.Print(u))
	return
}
