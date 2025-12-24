package env

import (
	"context"

	"github.com/ctrsploit/ctrsploit/env/namespace"
	namespace2 "github.com/ctrsploit/ctrsploit/pkg/namespace"
	"github.com/ctrsploit/sploit-spec/pkg/log"
	"github.com/urfave/cli/v3"
)

var Namespace = &cli.Command{
	Name:    namespace.CommandName,
	Aliases: []string{"n", "ns"},
	Usage:   "check namespace is host ns",
	Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
		var ns string
		if cmd.NArg() > 0 {
			ns = cmd.Args().First()
		}
		log.Logger.Debugf("namespace = %s\n", ns)
		if namespace2.CheckNamespaceValid(ns) {
			ctx = context.WithValue(ctx, "namespace", ns)
		}
		return ctx, nil
	},
	Action: func(ctx context.Context, cmd *cli.Command) (err error) {
		ns := ""
		if val := ctx.Value("namespace"); val != nil {
			ns = val.(string)
		}
		err = namespace.Print(ns)
		if err != nil {
			return
		}
		return
	},
}
