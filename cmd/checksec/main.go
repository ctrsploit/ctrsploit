package main

import (
	"github.com/ctrsploit/ctrsploit/cmd/ctrsploit/checksec"
	"github.com/ctrsploit/ctrsploit/internal"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"os"
)

const (
	name = `ctrsploit/checksec`
)

func main() {
	app := internal.Command2App(checksec.Command, true)
	app.Name = name
	internal.InstallGlobalFlags(app)
	err := app.Run(os.Args)
	if err != nil {
		awesome_error.CheckFatal(err)
	}
}
