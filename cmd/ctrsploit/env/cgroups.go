package env

import (
	"context"

	"github.com/ctrsploit/ctrsploit/env/cgroups"
	"github.com/urfave/cli/v3"
)

var Cgroups = &cli.Command{
	Name:    cgroups.CommandCgroupsName,
	Aliases: []string{"c"},
	Usage:   "gather cgroup information",
	Action: func(ctx context.Context, cmd *cli.Command) (err error) {
		err = cgroups.Print()
		if err != nil {
			return
		}
		return
	},
}
