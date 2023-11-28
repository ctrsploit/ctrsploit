package env

import (
	"github.com/ctrsploit/ctrsploit/env/graphdriver"
	"github.com/urfave/cli/v2"
)

var graphdriverCommand = &cli.Command{
	Name:    graphdriver.CommandName,
	Aliases: []string{"g"},
	Usage:   "detect graphdriver type and extend information",
	Action: func(context *cli.Context) (err error) {
		err = graphdriver.GraphDrivers()
		if err != nil {
			return
		}
		return
	},
}
