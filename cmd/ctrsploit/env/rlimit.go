package env

import (
	"context"

	"github.com/ctrsploit/ctrsploit/env/kernel/rlimit"
	"github.com/urfave/cli/v3"
)

var Rlimit = &cli.Command{
	Name:  "rlimit",
	Usage: "get process resource limits",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		err := rlimit.Print()
		if err != nil {
			return err
		}
		return nil
	},
}
