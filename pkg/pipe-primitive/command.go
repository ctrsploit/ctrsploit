package pipeprimitive

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

func Command(primitive Primitive, aliases []string, usage string) *cli.Command {
	return &cli.Command{
		Name:    primitive.GetExpName(),
		Aliases: aliases,
		Usage:   usage,
		Commands: []*cli.Command{
			{
				Name:    escalateName(primitive),
				Aliases: []string{"pe"},
				Usage:   fmt.Sprintf("permission escalate by using %s", primitive.GetExpName()),
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return Escalate(primitive)
				},
			},
			{
				Name:    escapeName(primitive),
				Aliases: []string{"e"},
				Usage:   fmt.Sprintf("escape by using %s", primitive.GetExpName()),
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return Escape(primitive)
				},
			},
			{
				Name:    imagePollutionName(primitive),
				Aliases: []string{"i"},
				Usage:   fmt.Sprintf("image pollution using %s", primitive.GetExpName()),
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "source", Aliases: []string{"s"}, Required: true,
						Usage: "the path of file with evil content"},
					&cli.StringFlag{Name: "destination", Aliases: []string{"d"}, Required: true,
						Usage: "the path of file you want to pollute"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return ImagePollution(primitive, cmd.String("source"), cmd.String("destination"))
				},
			},
		},
	}
}
