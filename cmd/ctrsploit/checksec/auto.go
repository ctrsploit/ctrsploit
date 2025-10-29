package checksec

//goland:noinspection GoSnakeCaseUsage
import (
	"github.com/ctrsploit/ctrsploit/vul/caps/bpf"
	"github.com/ctrsploit/ctrsploit/vul/caps/shocker"
	"github.com/ctrsploit/ctrsploit/vul/caps/sys_admin"
	"github.com/ctrsploit/ctrsploit/vul/caps/sys_ptrace"
	"github.com/ctrsploit/ctrsploit/vul/caps/sys_ptrace/pid_host"
	cve_2016_8867 "github.com/ctrsploit/ctrsploit/vul/cve-2016-8867"
	cve_2019_5736 "github.com/ctrsploit/ctrsploit/vul/cve-2019-5736"
	cve_2020_15257 "github.com/ctrsploit/ctrsploit/vul/cve-2020-15257"
	cve_2020_8558 "github.com/ctrsploit/ctrsploit/vul/cve-2020-8558"
	cve_2021_25741 "github.com/ctrsploit/ctrsploit/vul/cve-2021-25741"
	cve_2022_0492 "github.com/ctrsploit/ctrsploit/vul/cve-2022-0492"
	cve_2022_39253 "github.com/ctrsploit/ctrsploit/vul/cve-2022-39253"
	cve_2024_0132 "github.com/ctrsploit/ctrsploit/vul/cve-2024-0132"
	cve_2024_23650 "github.com/ctrsploit/ctrsploit/vul/cve-2024-23650"
	cve_2025_23266 "github.com/ctrsploit/ctrsploit/vul/cve-2025-23266"
	cve_2025_47290 "github.com/ctrsploit/ctrsploit/vul/cve-2025-47290"
	"github.com/ctrsploit/ctrsploit/vul/naked"
	"github.com/ctrsploit/ctrsploit/vul/namespace/net"
	"github.com/ctrsploit/ctrsploit/vul/namespace/pid"
	docker_sock "github.com/ctrsploit/ctrsploit/vul/shared-socket/docker-sock"
	"github.com/ctrsploit/sploit-spec/pkg/vul"
	"github.com/ssst0n3/awesome_libs/awesome_error"
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
			&cve_2016_8867.Vul,
			&cve_2019_5736.Vul,
			&cve_2020_8558.Vul,
			&cve_2020_15257.Vul,
			&cve_2021_25741.Vul,
			&cve_2022_39253.Vul,
			&cve_2022_0492.Vul,
			&cve_2024_0132.Vul,
			&cve_2024_23650.Vul,
			&cve_2025_23266.Vul,
			&cve_2025_47290.Vul,
			&shocker.Vul,
			&sys_admin.Vul,
			&bpf.Vul,
			&sys_ptrace.Vul,
			&pid_host.Vul,
			&naked.Vul,
			&net.Vul,
			&pid.Vul,
			&docker_sock.Vul,
		}
		err = vulnerabilities.Check(context)
		if err != nil {
			awesome_error.CheckWarning(err)
			err = nil
		}
		vulnerabilities.Output()
		return
	},
}
