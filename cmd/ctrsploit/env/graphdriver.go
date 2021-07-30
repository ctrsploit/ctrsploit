package env

import (
	"ctrsploit/env/graphdriver"
	"github.com/urfave/cli/v2"
)

var graphdriverCommand = &cli.Command{
	Name:    graphdriver.CommandGraphdriverName,
	Aliases: []string{"g"},
	Usage:   "detect graphdriver type and extend information",
	Action: func(context *cli.Context) (err error) {
		err = graphdriver.Overlay()
		if err != nil {
			return
		}
		err = graphdriver.DeviceMapper()
		if err != nil {
			return
		}
		return
	},
}
