package env

import (
	"context"

	"github.com/ctrsploit/ctrsploit/env/seccomp"
	"github.com/urfave/cli/v3"
)

var Seccomp = &cli.Command{
	Name:    seccomp.CommandName,
	Aliases: []string{"sc"},
	Usage:   "show the seccomp info",
	Action: func(ctx context.Context, cmd *cli.Command) (err error) {
		err = seccomp.Print()
		if err != nil {
			return
		}
		return
	},
}
