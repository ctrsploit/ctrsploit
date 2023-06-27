package checksec

import (
	"github.com/ctrsploit/ctrsploit/cmd/ctrsploit/env"
	"github.com/ctrsploit/ctrsploit/log"
	"github.com/ctrsploit/ctrsploit/vul"
	"github.com/sirupsen/logrus"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"github.com/ssst0n3/awesome_libs/awesome_error/exporter"
	log2 "github.com/ssst0n3/awesome_libs/log"
	"github.com/urfave/cli/v2"
	"io/ioutil"
)

var Command = &cli.Command{
	Name:    "checksec",
	Aliases: []string{"c"},
	Usage:   "check security inside a container",
	Subcommands: []*cli.Command{
		autoCommand,
		env.WhereCommand,
		env.SeccompCommand,
		env.ApparmorCommand,
		env.CgroupsCommand,
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
		debug := context.Bool("debug")
		awesome_error.Default = exporter.GetAwesomeError(log.Logger, debug)
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
	Action: func(context *cli.Context) (err error) {
		//err = auto.Auto()
		//if err != nil {
		//	return
		//}
		vulnerabilities := vul.Vulnerabilities{
			vul.SysadminCgroupV1,
		}
		for _, v := range vulnerabilities {
			_, err = v.CheckSec()
			if err != nil {
				return
			}
		}
		err = vulnerabilities.Output()
		if err != nil {
			return
		}
		return
	},
}
