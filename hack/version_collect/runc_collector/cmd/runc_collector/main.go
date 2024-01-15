package main

import (
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	app := &cli.App{
		Name: "runc_collector",
		Commands: []*cli.Command{
			GithubRelease,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
