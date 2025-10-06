package env

import (
	"github.com/ctrsploit/ctrsploit/env/seccomp"
	"github.com/urfave/cli/v2"
)

var Seccomp = &cli.Command{
	Name:    seccomp.CommandName,
	Aliases: []string{"sc"},
	Usage:   "show the seccomp info",
	Action: func(context *cli.Context) (err error) {
		err = seccomp.Print()
		if err != nil {
			return
		}
		return
	},
}
