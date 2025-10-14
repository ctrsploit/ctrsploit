package sysctl

import (
	"fmt"

	"github.com/ctrsploit/sploit-spec/pkg/env/container"
	"github.com/ctrsploit/sploit-spec/pkg/printer"
	"github.com/ctrsploit/sploit-spec/pkg/result"
	"github.com/ctrsploit/sploit-spec/pkg/result/item"
)

type Result struct {
	Name          result.Title
	RouteLocalNet item.Bool `json:"route_localnet"`
}

func Human(machine container.Sysctl) (human Result) {
	human = Result{
		Name: result.Title{
			Name: "Sysctl",
		},
		RouteLocalNet: item.Bool{
			Name:        "net.ipv4.conf.all.route_localnet",
			Description: "Indicates if the kernel parameter net.ipv4.conf.all.route_localnet is enabled.",
			Result:      machine.ProcSysNetIpv4ConfAllRouteLocalNet,
		},
	}
	return
}

func Print() (err error) {
	s, err := Sysctl()
	if err != nil {
		return
	}
	u := result.Union{
		Machine: s,
		Human:   Human(s),
	}
	fmt.Println(printer.Printer.Print(u))
	return
}
