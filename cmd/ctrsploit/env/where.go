package env

import (
	"context"

	"github.com/ctrsploit/ctrsploit/env/where"
	"github.com/urfave/cli/v3"
)

var Where = &cli.Command{
	Name:    where.CommandName,
	Aliases: []string{"w"},
	Usage:   "detect whether you are in the container, and which type of the container",
	Action: func(ctx context.Context, cmd *cli.Command) (err error) {
		err = where.Print()
		if err != nil {
			return
		}
		return
	},
}
