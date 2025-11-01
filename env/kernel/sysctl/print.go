package sysctl

import (
	"fmt"

	spec "github.com/ctrsploit/sploit-spec/pkg/env/container/kernel/sysctl"
	"github.com/ctrsploit/sploit-spec/pkg/printer"
	"github.com/ctrsploit/sploit-spec/pkg/result"
	"github.com/ctrsploit/sploit-spec/pkg/result/item"
)

type Result struct {
	Name                    result.Title
	RouteLocalNet           item.Bool  `json:"route_localnet"`
	MaxUserNamespaces       item.Short `json:"max_user_namespaces"`
	UnprivilegedUsernsClone item.Bool  `json:"unprivileged_userns_clone"`
	PidMax                  item.Short `json:"pid_max"`
}

func Human(machine spec.Sysctl) (human Result) {
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
			Description: "Specifies the maximum number of user namespaces that may exist on the system.",
			Result:      fmt.Sprintf("%d", machine.MaxUserNamespaces),
		},
		UnprivilegedUsernsClone: item.Bool{
			Name:        "kernel.unprivileged_userns_clone",
			Description: "Allow unprivileged process create user namespace",
			Result:      machine.UnprivilegedUsernsClone,
		},
		PidMax: item.Short{
			Name:        "kernel.pid_max",
			Description: "Sets the maximum PID value, controlling how many processes the system can run concurrently.",
			Result:      fmt.Sprintf("%d", machine.PidMax),
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
