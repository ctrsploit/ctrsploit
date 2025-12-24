package env

import (
	"context"

	"github.com/ctrsploit/ctrsploit/env/kernel"
	"github.com/urfave/cli/v3"
)

var Kernel = &cli.Command{
	Name:    kernel.CommandName,
	Aliases: []string{"k"},
	Usage:   "collect kernel environment information",
	Action: func(ctx context.Context, cmd *cli.Command) (err error) {
		err = kernel.Print()
		if err != nil {
			return
		}
		return
	},
}
