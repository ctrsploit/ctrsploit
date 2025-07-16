package checksec

import (
	"github.com/ctrsploit/ctrsploit/cmd/ctrsploit/env"
	"github.com/ctrsploit/ctrsploit/vul"
	cve_2016_8867 "github.com/ctrsploit/ctrsploit/vul/cve-2016-8867"
	cve_2020_15257 "github.com/ctrsploit/ctrsploit/vul/cve-2020-15257"
	cve_2025_47290 "github.com/ctrsploit/ctrsploit/vul/cve-2025-47290"
	"github.com/ctrsploit/ctrsploit/vul/shocker"
	"github.com/ctrsploit/ctrsploit/vul/sys_admin"
	"github.com/ctrsploit/ctrsploit/vul/sys_admin/release_agent"
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
		app.Vul2ChecksecCmd(&vul.NetworkNamespaceHostLevel, []string{"host"}, nil),
		cve_2025_47290.CheckSecCmd,
		cve_2020_15257.CheckSecCmd,
		sys_admin.CheckSecCmd,
		release_agent.CheckSecCmd,
		cve_2016_8867.CheckSecCmd,
		shocker.CheckSecCmd,
	},
}
