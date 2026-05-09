package crash

import (
	"context"

	pkgcrash "github.com/ctrsploit/ctrsploit/pkg/crash"
)

type Crash interface {
	Crash() (err error)
}

func MakeContainerCrash(cs ...Crash) (err error) {
	if len(cs) == 0 {
		triggers, err := pkgcrash.NewTriggers(nil, pkgcrash.Options{})
		if err != nil {
			return err
		}
		return pkgcrash.TriggerFirst(context.Background(), triggers...)
	}
	for _, c := range cs {
		err = c.Crash()
		if err != nil {
			return
		}
	}
	return
}
