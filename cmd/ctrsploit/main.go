package main

import (
	"github.com/ctrsploit/ctrsploit/cmd/ctrsploit/auto"
	"github.com/ctrsploit/ctrsploit/cmd/ctrsploit/checksec"
	"github.com/ctrsploit/ctrsploit/cmd/ctrsploit/env"
	"github.com/ctrsploit/ctrsploit/cmd/ctrsploit/exploit"
	"github.com/ctrsploit/ctrsploit/cmd/ctrsploit/helper"
	"github.com/ctrsploit/sploit-spec/pkg/app"
	"github.com/ctrsploit/sploit-spec/pkg/log"
	"github.com/ctrsploit/sploit-spec/pkg/version"
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
	sploit := &cli.App{
		Name:  "ctrsploit",
		Usage: usage,
		Commands: []*cli.Command{
			auto.Command,
			env.Command,
			exploit.Command,
			checksec.Command,
			helper.Command,
			version.Command,
		},
	}
	app.InstallGlobalFlags(sploit)
	err := sploit.Run(os.Args)
	if err != nil {
		log.Logger.Fatal(err)
	}
}
