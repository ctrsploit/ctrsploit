package main

import (
	"context"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{
		Name: "runc_collector",
		Commands: []*cli.Command{
			GithubRelease,
		},
	}
	err := app.Run(context.Background(), os.Args)
	if err != nil {
		panic(err)
	}
}
