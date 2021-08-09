package env

import (
	"github.com/ctrsploit/ctrsploit/env"
	"github.com/urfave/cli/v2"
)

var seccompCommand = &cli.Command{
	Name:    env.CommandSeccompName,
	Aliases: []string{"s"},
	Usage:   "show the seccomp info",
	Action: func(context *cli.Context) (err error) {
		err = env.Seccomp()
		if err != nil {
			return
		}
		return
	},
}
