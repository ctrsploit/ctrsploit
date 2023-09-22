package internal

import (
	"github.com/ctrsploit/ctrsploit/internal/colorful"
	"github.com/ctrsploit/ctrsploit/internal/log"
	"github.com/ctrsploit/ctrsploit/version"
	"github.com/sirupsen/logrus"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"github.com/ssst0n3/awesome_libs/awesome_error/exporter"
	log2 "github.com/ssst0n3/awesome_libs/log"
	"github.com/urfave/cli/v2"
	"io"
)

var (
	debugFlag = &cli.BoolFlag{
		Name:  "debug",
		Value: false,
		Usage: "Output information for helping debugging ctrsploit",
	}
	experimentalFlag = &cli.BoolFlag{
		Name:  "experimental",
		Value: false,
		Usage: "enable experimental feature",
	}
	colorfulFlag = &cli.BoolFlag{
		Name:  "colorful",
		Value: false,
		Usage: "output colorfully",
	}
)

func InstallGlobalFlagDebug(app *cli.App) {
	app.Flags = append(app.Flags, debugFlag)
	before := app.Before
	app.Before = func(context *cli.Context) (err error) {
		if before != nil {
			err = before(context)
			if err != nil {
				return
			}
		}
		debug := context.Bool("debug")
		awesome_error.Default = exporter.GetAwesomeError(log.Logger, debug)
		if !debug {
			log2.Logger.SetOutput(io.Discard)
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
	}
}

func InstallGlobalFlagExperimentalFlag(app *cli.App) {
	app.Flags = append(app.Flags, experimentalFlag)
}

func InstallGlobalFlagColorfulFlag(app *cli.App) {
	app.Flags = append(app.Flags, colorfulFlag)
	before := app.Before
	app.Before = func(ctx *cli.Context) (err error) {
		if before != nil {
			err = before(ctx)
			if err != nil {
				return
			}
		}
		flag := ctx.Bool("colorful")
		if flag {
			colorful.O = colorful.Colorful{}
		}
		return
	}
}

func InstallGlobalFlags(app *cli.App) {
	InstallGlobalFlagDebug(app)
	InstallGlobalFlagExperimentalFlag(app)
	InstallGlobalFlagColorfulFlag(app)
}

func Command2App(command *cli.Command, installGlobalFlags bool) (app *cli.App) {
	return &cli.App{
		Name:     command.Name,
		Usage:    command.Usage,
		Action:   command.Action,
		Commands: append(command.Subcommands, version.Command),
		Flags:    command.Flags,
		Before:   command.Before,
	}
}
