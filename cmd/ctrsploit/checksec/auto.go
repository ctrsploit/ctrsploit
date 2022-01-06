package checksec

import (
	"github.com/ctrsploit/ctrsploit/env/auto"
	"github.com/urfave/cli/v2"
)

var autoCommand = &cli.Command{
	Name:  "auto",
	Usage: "auto collect security information",
	Action: func(context *cli.Context) (err error) {
		err = auto.Auto()
		return
	},
}
