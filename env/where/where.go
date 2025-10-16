package where

import (
	"github.com/ctrsploit/ctrsploit/pkg/runtime"
	"github.com/ctrsploit/ctrsploit/pkg/where"
	"github.com/ctrsploit/sploit-spec/pkg/env/container"
	"github.com/ssst0n3/awesome_libs/awesome_error"
)

const CommandName = "where"

func runtime2Where(r *runtime.Runtime) container.Type {
	in, _ := r.Is()
	t := container.Type{
		In:    in,
		Rules: map[string]bool{},
	}
	for pre := range r.Prerequisites.Range() {
		satisfied, err := pre.Check()
		if err != nil && !pre.GetChecked() {
			awesome_error.CheckWarning(err)
		}
		t.Rules[pre.GetName()] = satisfied
	}
	return t
}

func Container() (t container.Type, err error) {
	c := where.Container{}
	in, err := c.IsIn()
	if err != nil {
		return
	}
	t = container.Type{
		In:    in,
		Rules: map[string]bool{},
	}
	return
}

func Where() (machine container.Where, err error) {
	c, err := Container()
	if err != nil {
		return
	}
	machine = container.Where{
		Container:  c,
		K8s:        runtime2Where(runtime.K8s()),
		Containerd: runtime2Where(runtime.Containerd()),
		Docker:     runtime2Where(runtime.Docker()),
		Nerdctl:    runtime2Where(runtime.Nerdctl()),
	}
	return
}
