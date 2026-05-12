package env

import (
	"context"

	"github.com/ctrsploit/ctrsploit/env/nonewprivs"
	"github.com/urfave/cli/v3"
)

var NoNewPrivs = &cli.Command{
	Name: nonewprivs.CommandName,
	Aliases: []string{
		"nnp",
		"no-new-privilege",
		"no-new-privileges",
	},
	Usage: "show NoNewPrivs status for the current process",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		return nonewprivs.Print()
	},
}
