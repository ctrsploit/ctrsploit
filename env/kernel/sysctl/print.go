package sysctl

import (
	"fmt"

	"github.com/ctrsploit/sploit-spec/pkg/env/container"
	"github.com/ctrsploit/sploit-spec/pkg/printer"
	"github.com/ctrsploit/sploit-spec/pkg/result"
	"github.com/ctrsploit/sploit-spec/pkg/result/item"
)

type Result struct {
	Name                    result.Title
	RouteLocalNet           item.Bool  `json:"route_localnet"`
	MaxUserNamespaces       item.Short `json:"max_user_namespaces"`
	UnprivilegedUsernsClone item.Bool  `json:"unprivileged_userns_clone"`
}

func Human(machine container.Sysctl) (human Result) {
	human = Result{
		Name: result.Title{
			Name: "Sysctl",
		},
		RouteLocalNet: item.Bool{
			Name:        "net.ipv4.conf.all.route_localnet",
			Description: "Indicates if the kernel parameter net.ipv4.conf.all.route_localnet is enabled.",
			Result:      machine.RouteLocalNet,
		},
		MaxUserNamespaces: item.Short{
			Name:        "user.max_user_namespaces",
			Description: "",
			Result:      fmt.Sprintf("%d", machine.MaxUserNamespaces),
		},
		UnprivilegedUsernsClone: item.Bool{
			Name:        "kernel.unprivileged_userns_clone",
			Description: "allow unprivileged process create user namespace",
			Result:      machine.UnprivilegedUsernsClone,
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
