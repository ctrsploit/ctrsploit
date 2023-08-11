package version

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

var Command = &cli.Command{
	Name:    "version",
	Aliases: []string{},
	Usage:   "Show the sploit version information",
	Action: func(context *cli.Context) error {
		fmt.Println(DefaultVer())
		return nil
	},
}
