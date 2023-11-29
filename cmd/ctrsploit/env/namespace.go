package env

import (
	"github.com/ctrsploit/ctrsploit/env/namespace"
	namespace2 "github.com/ctrsploit/ctrsploit/pkg/namespace"
	"github.com/ctrsploit/sploit-spec/pkg/log"
	"github.com/urfave/cli/v2"
)

var Namespace = &cli.Command{
	Name:    namespace.CommandName,
	Aliases: []string{"n", "ns"},
	Usage:   "check namespace is host ns",
	Before: func(context *cli.Context) (err error) {
		var ns string
		if context.NArg() > 0 {
			ns = context.Args().First()
		}
		log.Logger.Debugf("namespace = %s\n", ns)
		if namespace2.CheckNamespaceValid(ns) {
			context.App.Metadata["namespace"] = ns
		}
		return
	},
	Action: func(context *cli.Context) (err error) {
		ns, ok := context.App.Metadata["namespace"]
		if !ok {
			ns = ""
		}
		err = namespace.CurrentNamespaceLevel(ns.(string))
		if err != nil {
			return
		}
		return
	},
}
