package main

import (
	"github.com/ctrsploit/ctrsploit/log"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	app := &cli.App{
		Name:     name,
		Usage:    usage,
		Commands: []*cli.Command{},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Logger.Fatal(err)
	}
}
