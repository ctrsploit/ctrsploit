package env

import (
	"ctrsploit/env"
	"github.com/urfave/cli/v2"
)

var cgroupsCommand = &cli.Command{
	Name:    env.CommandCgroupsName,
	Aliases: []string{"c"},
	Usage:   "gather cgroup information",
	Action: func(context *cli.Context) (err error) {
		err = env.Version()
		if err != nil {
			return
		}
		return
	},
}
