package env

import (
	"github.com/ctrsploit/ctrsploit/env/cgroups"
	"github.com/urfave/cli/v2"
)

var Cgroups = &cli.Command{
	Name:    cgroups.CommandCgroupsName,
	Aliases: []string{"c"},
	Usage:   "gather cgroup information",
	Action: func(context *cli.Context) (err error) {
		err = cgroups.Cgroups()
		if err != nil {
			return
		}
		return
	},
}
