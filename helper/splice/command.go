package splice

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

func Command(splice Splice, aliases []string, usage string) (cmd *cli.Command) {
	return &cli.Command{
		Name:    splice.GetExpName(),
		Aliases: aliases,
		Usage:   usage,
		Subcommands: []*cli.Command{
			{
				Name:    escalateExpName(splice),
				Aliases: []string{"pe"},
				Usage:   fmt.Sprintf("permission escalate by using %s", splice.GetExpName()),
				Action: func(context *cli.Context) error {
					return InvokeEscalate(splice)
				},
			},
			{
				Name:    escapeExpName(splice),
				Aliases: []string{"e"},
				Usage:   fmt.Sprintf("escape by using %s", splice.GetExpName()),
				Action: func(context *cli.Context) error {
					return InvokeEscape(splice)
				},
			},
		},
	}
}
