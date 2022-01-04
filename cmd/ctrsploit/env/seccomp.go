package env

import (
	"github.com/ctrsploit/ctrsploit/env/seccomp"
	"github.com/urfave/cli/v2"
)

var SeccompCommand = &cli.Command{
	Name:    seccomp.CommandSeccompName,
	Aliases: []string{"s"},
	Usage:   "show the seccomp info",
	Action: func(context *cli.Context) (err error) {
		err = seccomp.Seccomp()
		if err != nil {
			return
		}
		return
	},
}
