package env

import (
	"github.com/ctrsploit/ctrsploit/env/kernel/rlimit"
	"github.com/urfave/cli/v2"
)

var Rlimit = &cli.Command{
	Name:  "rlimit",
	Usage: "get process resource limits",
	Action: func(c *cli.Context) error {
		err := rlimit.Print()
		if err != nil {
			return err
		}
		return nil
	},
}
