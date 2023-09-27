package checksec

import (
	"github.com/ctrsploit/ctrsploit/vul"
	"github.com/urfave/cli/v2"
)

var Auto = &cli.Command{
	Name:  "auto",
	Usage: "auto check security",
	Action: func(context *cli.Context) (err error) {
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
