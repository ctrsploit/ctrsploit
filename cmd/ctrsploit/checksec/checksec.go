package checksec

import (
	"github.com/ctrsploit/ctrsploit/cmd/ctrsploit/env"
	"github.com/ctrsploit/ctrsploit/vul/caps/shocker"
	"github.com/ctrsploit/ctrsploit/vul/caps/sys_admin"
	"github.com/ctrsploit/ctrsploit/vul/caps/sys_admin/ebpf"
	"github.com/ctrsploit/ctrsploit/vul/caps/sys_admin/release_agent"
	"github.com/ctrsploit/ctrsploit/vul/caps/sys_ptrace"
	"github.com/ctrsploit/ctrsploit/vul/caps/sys_ptrace/pid_host"
	cve_2016_8867 "github.com/ctrsploit/ctrsploit/vul/cve-2016-8867"
	cve_2020_15257 "github.com/ctrsploit/ctrsploit/vul/cve-2020-15257"
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
	"github.com/ctrsploit/sploit-spec/pkg/app"
	"github.com/urfave/cli/v2"
)

var Command = &cli.Command{
	Name:    "checksec",
	Aliases: []string{"c"},
	Usage:   "check security inside a container",
	Subcommands: []*cli.Command{
		Auto,
		env.Command,
		app.Vul2ChecksecCmd(&net.Vul, []string{"host"}, nil),
		cve_2025_23266.CheckSecCmd,
		cve_2025_47290.CheckSecCmd,
		cve_2024_0132.CheckSecCmd,
		cve_2024_23650.CheckSecCmd,
		cve_2022_39253.CheckSecCmd,
		cve_2022_0492.CheckSecCmd,
		cve_2021_25741.CheckSecCmd,
		cve_2020_15257.CheckSecCmd,
		cve_2016_8867.CheckSecCmd,
		naked.CheckSecCmd,
		sys_admin.CheckSecCmd,
		ebpf.CheckSecCmd,
		release_agent.CheckSecCmd,
		shocker.CheckSecCmd,
		pid_host.CheckSecCmd,
		sys_ptrace.CheckSecCmd,
		pid.CheckSecCmd,
		docker_sock.CheckSecCmd,
	},
}
