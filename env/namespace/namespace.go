package namespace

import (
	"fmt"
	"github.com/ctrsploit/ctrsploit/internal"
	"github.com/ctrsploit/ctrsploit/internal/log"
	"github.com/ctrsploit/ctrsploit/pkg/namespace"
	"github.com/ctrsploit/ctrsploit/prerequisite/kernel"
	"github.com/ctrsploit/sploit-spec/pkg/colorful"
	"github.com/ctrsploit/sploit-spec/pkg/result"
	"github.com/ctrsploit/sploit-spec/pkg/result/item"
)

const CommandName = "namespace"

type Result struct {
	Name   result.Title
	Levels []item.Short `json:"levels"`
}

func (r Result) String() (s string) {
	s += internal.Print(r.Name)
	for _, l := range r.Levels {
		s += internal.Print(l)
	}
	return
}

func getNamespaceLevels() (namespaceLevels map[string]namespace.Level, names []string, err error) {
	arbitrator, err := namespace.NewInoArbitrator()
	if err != nil {
		return
	}
	namespaceLevels, names, err = namespace.CheckNamespaceLevel(arbitrator)
	if err != nil {
		return
	}
	return
}

func level2result(name string, level namespace.Level) item.Short {
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

func CurrentNamespaceLevel(ns string) (err error) {
	r := Result{
		Name: result.Title{
			Name: "Namespace Level",
		},
		Levels: []item.Short{},
	}
	namespaceLevels, names, err := getNamespaceLevels()
	if err != nil {
		return
	}
	if ns == "" {
		for _, name := range names {
			level := namespaceLevels[name]
			r.Levels = append(r.Levels, level2result(name, level))
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
		r.Levels = append(r.Levels, level2result(ns, level))
	}
	fmt.Println(r)
	return
}
