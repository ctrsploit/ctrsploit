package checksec

import (
	"github.com/ctrsploit/ctrsploit/cmd/ctrsploit/env"
	"github.com/ctrsploit/ctrsploit/vul"
	"github.com/urfave/cli/v2"
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
		env.NamespaceCommand,
	},
	Action: func(context *cli.Context) (err error) {
		//err = auto.Auto()
		//if err != nil {
		//	return
		//}
		vulnerabilities := vul.Vulnerabilities{
			&vul.SysadminCgroupV1,
			&vul.NetworkNamespaceHostLevel,
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
