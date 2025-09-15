package checksec

import (
	"github.com/ctrsploit/ctrsploit/vul/caps/sys_admin/release_agent"
	"github.com/ctrsploit/ctrsploit/vul/namespace/net"
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
			&release_agent.Vul,
			&net.Vul,
		}
		err = vulnerabilities.Check(context)
		if err != nil {
			return
		}
		vulnerabilities.Output()
		return
	},
}
