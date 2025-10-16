package kernel

import (
	"github.com/ctrsploit/ctrsploit/env/kernel/sysctl"
	"github.com/ctrsploit/sploit-spec/pkg/env/container"
	"github.com/ctrsploit/sploit-spec/pkg/printer"
	"github.com/ctrsploit/sploit-spec/pkg/result"
	"github.com/ctrsploit/sploit-spec/pkg/result/item"
)

type Result struct {
	Name         result.Title
	CompiledDate item.Short `json:"compiled_date"`
	sysctl.Result
}

func Human(machine container.Kernel) (human Result) {
	human = Result{
		Name: result.Title{
			Name: "Kernel",
		},
		CompiledDate: item.Short{
			Name:        "Compiled Date",
			Description: "The date when the kernel was compiled.",
			Result:      machine.CompiledDate.String(),
		},
		Result: sysctl.Human(machine.Sysctl),
	}
	return
}

func Print() (err error) {
	k, err := Kernel()
	if err != nil {
		return
	}
	u := result.Union{
		Machine: k,
		Human:   Human(k),
	}
	println(printer.Printer.Print(u))
	return
}
