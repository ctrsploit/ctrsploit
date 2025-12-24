package env

import (
	"context"

	"github.com/ctrsploit/ctrsploit/env/capability"
	"github.com/urfave/cli/v3"
)

var Capability = &cli.Command{
	Name:    capability.CommandCapabilityName,
	Aliases: []string{"cap"},
	Usage:   "show the capability of pid 1 and current process",
	Action: func(ctx context.Context, cmd *cli.Command) (err error) {
		err = capability.Print()
		if err != nil {
			return
		}
		return
	},
}
