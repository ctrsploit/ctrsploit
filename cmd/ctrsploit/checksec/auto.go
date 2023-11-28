package checksec

import (
	vul2 "github.com/ctrsploit/ctrsploit/vul"
	"github.com/ctrsploit/sploit-spec/pkg/vul"
	"github.com/urfave/cli/v2"
)

const (
	CommandNameAuto = "auto"
)

var Auto = &cli.Command{
	Name:    CommandNameAuto,
	Usage:   "auto check security",
	Aliases: []string{"a"},
	Action: func(context *cli.Context) (err error) {
		vulnerabilities := vul.Vulnerabilities{
			&vul2.SysadminCgroupV1,
			&vul2.NetworkNamespaceHostLevel,
		}
		err = vulnerabilities.Check()
		if err != nil {
			return
		}
		vulnerabilities.Output()
		return
	},
}
