package main

import (
	"ctrsploit/cmd/ctrsploit/env"
	"ctrsploit/cmd/ctrsploit/exploit"
	"ctrsploit/log"
	"github.com/docker/docker/pkg/reexec"
	"github.com/sirupsen/logrus"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"github.com/ssst0n3/awesome_libs/awesome_error/exporter"
	log2 "github.com/ssst0n3/awesome_libs/log"
	"github.com/urfave/cli/v2"
	"io/ioutil"
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
			autoCommand,
			exploit.ExploitCommand,
			env.Command,
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "lang",
				Value: "english",
				Usage: "language for the greeting",
			},
			&cli.BoolFlag{
				Name:  "debug",
				Value: false,
				Usage: "Output information for helping debugging ctrsploit",
			},
		},
		Before: func(context *cli.Context) (err error) {
			awesome_error.Default = exporter.GetAwesomeError(log.Logger)
			debug := context.Bool("debug")
			if !debug {
				log2.Logger.SetOutput(ioutil.Discard)
			} else {
				log.Logger.Level = logrus.DebugLevel
				log.Logger.SetReportCaller(true)
				log.Logger.SetFormatter(&logrus.TextFormatter{
					ForceColors: true,
				})
				log2.Logger.Level = logrus.DebugLevel
				log2.Logger.Debug("debug mode on")
			}
			return
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Logger.Fatal(err)
	}
}
