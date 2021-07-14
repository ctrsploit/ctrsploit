package env

import (
	"ctrsploit/env/capability"
	"github.com/urfave/cli/v2"
)

var capabilityCommand = &cli.Command{
	Name:    capability.CommandName,
	Aliases: []string{"cap"},
	Usage:   "show the capability of pid 1",
	Action: func(context *cli.Context) (err error) {
		err = capability.GetPid1Capability()
		if err != nil {
			return
		}
		return
	},
}
