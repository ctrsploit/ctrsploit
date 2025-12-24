package pipe_primitive

import (
	"context"
	"fmt"

	"github.com/ssst0n3/awesome_libs/log"
	"github.com/urfave/cli/v3"
)

func Command(primitive Primitive, aliases []string, usage string) (cmd *cli.Command) {
	return &cli.Command{
		Name:    primitive.GetExpName(),
		Aliases: aliases,
		Usage:   usage,
		Commands: []*cli.Command{
			{
				Name:    escalateExpName(primitive),
				Aliases: []string{"pe"},
				Usage:   fmt.Sprintf("permission escalate by using %s", primitive.GetExpName()),
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return InvokeEscalate(primitive)
				},
			},
			{
				Name:    escapeExpName(primitive),
				Aliases: []string{"e"},
				Usage:   fmt.Sprintf("escape by using %s", primitive.GetExpName()),
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return InvokeEscape(primitive)
				},
			},
			{
				Name:    imagePollutionExpName(primitive),
				Aliases: []string{"i"},
				Usage:   fmt.Sprintf("image pollusion using %s", primitive.GetExpName()),
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "source", Aliases: []string{"s"}, Required: true,
						Usage: "the path of file with evil content"},
					&cli.StringFlag{Name: "destination", Aliases: []string{"d"}, Required: true,
						Usage: "the path of file you want to pollution"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					source := cmd.String("source")
					dest := cmd.String("destination")
					log.Logger.Info(fmt.Sprintf("Overwrite %s with %s", source, dest))
					return InvokeImagePollution(primitive, source, dest)
				},
			},
		},
	}
}
