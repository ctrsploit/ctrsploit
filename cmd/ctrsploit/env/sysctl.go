package env

import (
	"context"

	"github.com/ctrsploit/ctrsploit/env/kernel/sysctl"
	"github.com/urfave/cli/v3"
)

var Sysctl = &cli.Command{
	Name:    sysctl.CommandName,
	Aliases: []string{},
	Usage:   "display sysctl information",
	Action: func(ctx context.Context, cmd *cli.Command) (err error) {
		err = sysctl.Print()
		if err != nil {
			return
		}
		return
	},
}
