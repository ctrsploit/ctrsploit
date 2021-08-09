package env

import (
	"github.com/ctrsploit/ctrsploit/env/capability"
	"github.com/urfave/cli/v2"
)

var capabilityCommand = &cli.Command{
	Name:    capability.CommandCapabilityName,
	Aliases: []string{"cap"},
	Usage:   "show the capability of pid 1 and current process",
	Action: func(context *cli.Context) (err error) {
		err = capability.Capability()
		if err != nil {
			return
		}
		return
	},
}
