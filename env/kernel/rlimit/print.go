package rlimit

import (
	"fmt"

	"github.com/ctrsploit/ctrsploit/pkg/rlimit"
	spec "github.com/ctrsploit/sploit-spec/pkg/env/container/kernel/rlimit"
	"github.com/ctrsploit/sploit-spec/pkg/printer"
	"github.com/ctrsploit/sploit-spec/pkg/result"
	"github.com/ctrsploit/sploit-spec/pkg/result/item"
)

type Result struct {
	Name       result.Title
	Core       item.Short
	Cpu        item.Short
	Data       item.Short
	Fsize      item.Short
	Locks      item.Short
	Msgqueue   item.Short
	Nice       item.Short
	Rtprio     item.Short
	Rttime     item.Short
	Sigpending item.Short
	Stack      item.Short
	As         item.Short
	Memlock    item.Short
	Nofile     item.Short
	Nproc      item.Short
	Rss        item.Short
}

func Human(machine spec.Rlimit) Result {
	return Result{
		Name: result.Title{
			Name: "Rlimit",
		},
		Core: item.Short{
			Name:        "core",
			Description: "max core file size",
			Result:      rlimit.Resource(machine.Core).String(),
		},
		Cpu: item.Short{
			Name:        "cpu",
			Description: "CPU time in seconds",
			Result:      rlimit.Resource(machine.Cpu).String(),
		},
		Data: item.Short{
			Name:        "data",
			Description: "max data size",
			Result:      rlimit.Resource(machine.Data).String(),
		},
		Fsize: item.Short{
			Name:        "fsize",
			Description: "max file size",
			Result:      rlimit.Resource(machine.Fsize).String(),
		},
		Locks: item.Short{
			Name:        "locks",
			Description: "max number of file locks",
			Result:      rlimit.Resource(machine.Locks).String(),
		},
		Msgqueue: item.Short{
			Name:        "msgqueue",
			Description: "max bytes in POSIX message queues",
			Result:      rlimit.Resource(machine.Msgqueue).String(),
		},
		Nice: item.Short{
			Name:        "nice",
			Description: "max nice priority",
			Result:      rlimit.Resource(machine.Nice).String(),
		},
		Rtprio: item.Short{
			Name:        "rtprio",
			Description: "max real-time priority",
			Result:      rlimit.Resource(machine.Rtprio).String(),
		},
		Rttime: item.Short{
			Name:        "rttime",
			Description: "timeout for real-time tasks",
			Result:      rlimit.Resource(machine.Rttime).String(),
		},
		Sigpending: item.Short{
			Name:        "sigpending",
			Description: "max number of pending signals",
			Result:      rlimit.Resource(machine.Sigpending).String(),
		},
		Stack: item.Short{
			Name:        "stack",
			Description: "max stack size",
			Result:      rlimit.Resource(machine.Stack).String(),
		},
		As: item.Short{
			Name:        "as",
			Description: "address space (virtual memory)",
			Result:      rlimit.Resource(machine.As).String(),
		},
		Memlock: item.Short{
			Name:        "memlock",
			Description: "max locked-in-memory address space",
			Result:      rlimit.Resource(machine.Memlock).String(),
		},
		Nofile: item.Short{
			Name:        "nofile",
			Description: "max number of open files",
			Result:      rlimit.Resource(machine.Nofile).String(),
		},
		Nproc: item.Short{
			Name:        "nproc",
			Description: "max number of processes",
			Result:      rlimit.Resource(machine.Nproc).String(),
		},
		Rss: item.Short{
			Name:        "rss",
			Description: "max resident set size",
			Result:      rlimit.Resource(machine.Rss).String(),
		},
	}
}

func Print() (err error) {
	r, err := Rlimit()
	if err != nil {
		return err
	}
	u := result.Union{
		Machine: r,
		Human:   Human(r),
	}
	fmt.Println(printer.Printer.Print(u))
	return
}
