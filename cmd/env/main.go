package main

import (
	"context"
	"os"

	"github.com/ctrsploit/ctrsploit/cmd/ctrsploit/env"
	"github.com/ctrsploit/sploit-spec/pkg/app"
	"github.com/ssst0n3/awesome_libs/awesome_error"
)

const (
	name = `ctrsploit/env`
)

func main() {
	sploit := app.Command2App(env.Command)
	sploit.Name = name
	app.InstallGlobalFlags(sploit)
	err := sploit.Run(context.Background(), os.Args)
	if err != nil {
		awesome_error.CheckFatal(err)
	}
}
