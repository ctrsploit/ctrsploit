package checksec

//goland:noinspection GoSnakeCaseUsage
import (
	"context"

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
	cve_2021_25748 "github.com/ctrsploit/ctrsploit/vul/cve-2021-25748"
	cve_2021_3493 "github.com/ctrsploit/ctrsploit/vul/cve-2021-3493"
	cve_2022_0492 "github.com/ctrsploit/ctrsploit/vul/cve-2022-0492"
	cve_2022_0847 "github.com/ctrsploit/ctrsploit/vul/cve-2022-0847"
	cve_2022_39253 "github.com/ctrsploit/ctrsploit/vul/cve-2022-39253"
	cve_2024_0132 "github.com/ctrsploit/ctrsploit/vul/cve-2024-0132"
	cve_2024_23650 "github.com/ctrsploit/ctrsploit/vul/cve-2024-23650"
	cve_2025_23266 "github.com/ctrsploit/ctrsploit/vul/cve-2025-23266"
	cve_2025_47290 "github.com/ctrsploit/ctrsploit/vul/cve-2025-47290"
	cve_2025_62725 "github.com/ctrsploit/ctrsploit/vul/cve-2025-62725"
	cve_2026_23111 "github.com/ctrsploit/ctrsploit/vul/cve-2026-23111"
	cve_2026_31431 "github.com/ctrsploit/ctrsploit/vul/cve-2026-31431"
	cve_2026_43284 "github.com/ctrsploit/ctrsploit/vul/cve-2026-43284"
	cve_2026_43500 "github.com/ctrsploit/ctrsploit/vul/cve-2026-43500"
	cve_2026_46300 "github.com/ctrsploit/ctrsploit/vul/cve-2026-46300"
	fork_bomb "github.com/ctrsploit/ctrsploit/vul/fork-bomb"
	"github.com/ctrsploit/ctrsploit/vul/naked"
	"github.com/ctrsploit/ctrsploit/vul/namespace/net"
	"github.com/ctrsploit/ctrsploit/vul/namespace/pid"
	access_secrets "github.com/ctrsploit/ctrsploit/vul/sa-token/access-secrets"
	sa_token_policy "github.com/ctrsploit/ctrsploit/vul/sa-token/policy"
	docker_sock "github.com/ctrsploit/ctrsploit/vul/shared-socket/docker-sock"
	"github.com/ctrsploit/sploit-spec/pkg/vul"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"github.com/urfave/cli/v3"
)

const (
	CommandNameAuto = "auto"
)

var All = vul.Vulnerabilities{
	&cve_2016_8867.Vul,
	&cve_2019_5736.Vul,
	&cve_2020_8558.Vul,
	&cve_2020_15257.Vul,
	&cve_2021_25741.Vul,
	&cve_2021_25748.Vul,
	&cve_2021_3493.Vul,
	&cve_2022_0492.Vul,
	&cve_2022_0847.Vul,
	&cve_2022_39253.Vul,
	&cve_2024_0132.Vul,
	&cve_2024_23650.Vul,
	&cve_2025_23266.Vul,
	&cve_2025_47290.Vul,
	&cve_2025_62725.Vul,
	&cve_2026_23111.Vul,
	&cve_2026_31431.Vul,
	&cve_2026_43284.Vul,
	&cve_2026_43500.Vul,
	&cve_2026_46300.Vul,
	&fork_bomb.Vul,
	&shocker.Vul,
	&sys_admin.Vul,
	&bpf.Vul,
	&sys_ptrace.Vul,
	&pid_host.Vul,
	&naked.Vul,
	&net.Vul,
	&pid.Vul,
	&access_secrets.Vul,
	&sa_token_policy.Vul,
	&docker_sock.Vul,
}

var Auto = &cli.Command{
	Name:    CommandNameAuto,
	Usage:   "auto check security",
	Aliases: []string{"a"},
	Action: func(ctx context.Context, cmd *cli.Command) (err error) {
		err = All.Check(cmd)
		if err != nil {
			awesome_error.CheckWarning(err)
			err = nil
		}
		All.Output()
		return
	},
}
