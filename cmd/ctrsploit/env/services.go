package env

import (
	"context"
	"strings"

	"github.com/ctrsploit/ctrsploit/env/services"
	"github.com/urfave/cli/v3"
)

var Services = &cli.Command{
	Name:    services.CommandName,
	Aliases: []string{"svc"},
	Usage:   "discover K8s cluster services and ports via env vars and DNS",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "zone",
			Aliases: []string{"z"},
			Value:   "cluster.local",
			Usage:   "K8s cluster DNS zone",
		},
		&cli.StringFlag{
			Name:    "cidr",
			Aliases: []string{"c"},
			Usage:   "service CIDR to scan (auto-detected from KUBERNETES_SERVICE_HOST if empty)",
		},
		&cli.IntFlag{
			Name:    "threads",
			Aliases: []string{"t"},
			Value:   16,
			Usage:   "number of threads for CIDR scanning",
		},
		&cli.StringFlag{
			Name:    "methods",
			Aliases: []string{"m"},
			Value:   "all",
			Usage:   "discovery methods: all,env,wildcard,axfr,cidr (comma-separated)",
		},
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "export results to file (NDJSON format)",
		},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		opts := services.DefaultOptions()
		if cmd.IsSet("zone") {
			opts.Zone = cmd.String("zone")
		}
		if cmd.IsSet("cidr") {
			opts.CIDR = cmd.String("cidr")
		}
		if cmd.IsSet("threads") {
			opts.Threads = int(cmd.Int("threads"))
		}
		if cmd.IsSet("methods") {
			opts.Methods = strings.Split(cmd.String("methods"), ",")
		}
		opts.OutputFile = cmd.String("output")
		return services.Print(opts)
	},
}
