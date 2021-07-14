package env

import (
	"ctrsploit/env"
	"github.com/urfave/cli/v2"
)

var capabilityCommand = &cli.Command{
	Name:    env.CommandCapabilityName,
	Aliases: []string{"cap"},
	Usage:   "show the capability of pid 1 and current process",
	Action: func(context *cli.Context) (err error) {
		err = env.Capability()
		if err != nil {
			return
		}
		return
	},
}
