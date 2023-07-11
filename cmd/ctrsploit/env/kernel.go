package env

import (
	"github.com/ctrsploit/ctrsploit/env/kernel"
	"github.com/urfave/cli/v2"
)

var kernelCommand = &cli.Command{
	Name:    kernel.CommandName,
	Aliases: []string{"k"},
	Usage:   "collect kernel environment information",
	Action: func(context *cli.Context) (err error) {
		return
	},
}
