package env

import (
	"ctrsploit/env"
	"github.com/urfave/cli/v2"
)

var graphdriverCommand = &cli.Command{
	Name:    env.CommandGraphdriverName,
	Aliases: []string{"g"},
	Usage:   "detect graphdriver type and extend information",
	Action: func(context *cli.Context) (err error) {
		err = env.Overlay()
		if err != nil {
			return
		}
		err = env.DeviceMapper()
		if err != nil {
			return
		}
		return
	},
}
