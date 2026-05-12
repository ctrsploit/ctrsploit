package nonewprivs

import (
	"fmt"

	"github.com/ctrsploit/sploit-spec/pkg/printer"
	"github.com/ctrsploit/sploit-spec/pkg/result"
	"github.com/ctrsploit/sploit-spec/pkg/result/item"
)

type Result struct {
	Name        result.Title `json:"name"`
	Enabled     item.Bool    `json:"enabled"`
	SUIDBlocked item.Bool    `json:"suid_blocked"`
	StatusPath  item.Short   `json:"status_path"`
}

func Human(machine Info) Result {
	return Result{
		Name: result.Title{
			Name: "NoNewPrivs",
		},
		Enabled: item.Bool{
			Name:        "NoNewPrivs enabled",
			Description: "current process cannot gain new privileges across execve",
			Result:      machine.Enabled,
		},
		SUIDBlocked: item.Bool{
			Name:        "SUID privilege gain blocked",
			Description: "SUID and file capabilities cannot raise privileges when NoNewPrivs is enabled",
			Result:      machine.Enabled,
		},
		StatusPath: item.Short{
			Name:   "source",
			Result: machine.StatusPath,
		},
	}
}

func Print() error {
	machine, err := Current()
	if err != nil {
		return err
	}
	u := result.Union{
		Machine: machine,
		Human:   Human(machine),
	}
	fmt.Println(printer.Printer.Print(u))
	return nil
}
