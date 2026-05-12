package env

import (
	"context"

	"github.com/ctrsploit/ctrsploit/env/suid"
	"github.com/urfave/cli/v3"
)

var SUID = &cli.Command{
	Name:    suid.CommandName,
	Aliases: []string{"setuid"},
	Usage:   "find and list SUID files",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "all",
			Aliases: []string{"a"},
			Usage:   "scan from / instead of common executable directories",
		},
		&cli.StringFlag{
			Name:    "path",
			Aliases: []string{"p"},
			Usage:   "comma-separated paths to scan (default: common executable directories)",
		},
		&cli.StringFlag{
			Name:    "skip",
			Aliases: []string{"s"},
			Usage:   "additional comma-separated directories to skip",
		},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		opts := suid.DefaultOptions()
		if cmd.Bool("all") {
			opts.Paths = []string{"/"}
		}
		if cmd.IsSet("path") {
			opts.Paths = suid.ParsePaths(cmd.String("path"))
		}
		if cmd.IsSet("skip") {
			opts.SkipDirs = append(opts.SkipDirs, suid.ParsePaths(cmd.String("skip"))...)
		}
		return suid.Print(opts)
	},
}
