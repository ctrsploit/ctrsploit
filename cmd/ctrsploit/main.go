package main

import (
	"github.com/ctrsploit/ctrsploit/cmd/ctrsploit/checksec"
	"github.com/ctrsploit/ctrsploit/cmd/ctrsploit/env"
	"github.com/ctrsploit/ctrsploit/cmd/ctrsploit/exploit"
	"github.com/ctrsploit/ctrsploit/cmd/ctrsploit/helper"
	"github.com/ctrsploit/ctrsploit/internal"
	"github.com/ctrsploit/ctrsploit/internal/log"
	"github.com/ctrsploit/ctrsploit/version"
	"github.com/docker/docker/pkg/reexec"
	"github.com/urfave/cli/v2"
	"os"
)

const usage = `A penetration toolkit for container environment

ctrsploit is a command line ... //TODO
`

func main() {
	if reexec.Init() {
		return
	}
	app := &cli.App{
		Name:  "ctrsploit",
		Usage: usage,
		Commands: []*cli.Command{
			env.Command,
			exploit.Command,
			checksec.Command,
			autoCommand,
			helper.Command,
			version.Command,
		},
	}
	internal.InstallGlobalFlags(app)
	err := app.Run(os.Args)
	if err != nil {
		log.Logger.Fatal(err)
	}
}
