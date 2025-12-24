package env

import (
	"context"

	"github.com/ctrsploit/ctrsploit/env/mountinfo"
	"github.com/urfave/cli/v3"
)

var Mountinfo = &cli.Command{
	Name:    "mountinfo",
	Aliases: []string{"m"},
	Usage:   "list mount points",
	Action: func(ctx context.Context, cmd *cli.Command) (err error) {
		err = mountinfo.Print()
		if err != nil {
			return
		}
		return
	},
}
