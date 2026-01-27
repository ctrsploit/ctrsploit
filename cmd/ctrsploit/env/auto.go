package env

import (
	"context"

	"github.com/ctrsploit/ctrsploit/env/auto"
	"github.com/urfave/cli/v3"
)

const (
	CommandNameAuto = "auto"
)

var (
	Auto = &cli.Command{
		Name:  CommandNameAuto,
		Usage: "auto",
		Action: func(ctx context.Context, cmd *cli.Command) (err error) {
			return auto.Print()
		},
	}
)
