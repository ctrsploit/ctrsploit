package vul

//goland:noinspection GoSnakeCaseUsage
import (
	"github.com/ctrsploit/ctrsploit/vul/caps"
	cve_2016_8867 "github.com/ctrsploit/ctrsploit/vul/cve-2016-8867"
	cve_2019_5736 "github.com/ctrsploit/ctrsploit/vul/cve-2019-5736"
	cve_2020_15257 "github.com/ctrsploit/ctrsploit/vul/cve-2020-15257"
	cve_2020_8558 "github.com/ctrsploit/ctrsploit/vul/cve-2020-8558"
	cve_2021_25741 "github.com/ctrsploit/ctrsploit/vul/cve-2021-25741"
	cve_2021_25748 "github.com/ctrsploit/ctrsploit/vul/cve-2021-25748"
	cve_2022_0492 "github.com/ctrsploit/ctrsploit/vul/cve-2022-0492"
	cve_2022_39253 "github.com/ctrsploit/ctrsploit/vul/cve-2022-39253"
	cve_2024_0132 "github.com/ctrsploit/ctrsploit/vul/cve-2024-0132"
	cve_2024_23650 "github.com/ctrsploit/ctrsploit/vul/cve-2024-23650"
	cve_2025_23266 "github.com/ctrsploit/ctrsploit/vul/cve-2025-23266"
	cve_2025_47290 "github.com/ctrsploit/ctrsploit/vul/cve-2025-47290"
	cve_2025_62725 "github.com/ctrsploit/ctrsploit/vul/cve-2025-62725"
	fork_bomb "github.com/ctrsploit/ctrsploit/vul/fork-bomb"
	"github.com/ctrsploit/ctrsploit/vul/naked"
	"github.com/ctrsploit/ctrsploit/vul/namespace"
	shared_socket "github.com/ctrsploit/ctrsploit/vul/shared-socket"
	"github.com/urfave/cli/v2"
)

var Command = &cli.Command{
	Name:    "vul",
	Aliases: []string{"v"},
	Usage:   "list vulnerabilities",
	Subcommands: []*cli.Command{
		cve_2016_8867.VulCmd,
		cve_2019_5736.VulCmd,
		cve_2020_8558.VulCmd,
		cve_2020_15257.VulCmd,
		cve_2021_25741.VulCmd,
		cve_2021_25748.VulCmd,
		cve_2022_0492.VulCmd,
		cve_2022_39253.VulCmd,
		cve_2024_0132.VulCmd,
		cve_2024_23650.VulCmd,
		cve_2025_23266.VulCmd,
		cve_2025_47290.VulCmd,
		cve_2025_62725.VulCmd,
		fork_bomb.VulCmd,
		naked.VulCmd,
		caps.VulCmd,
		namespace.VulCmd,
		shared_socket.VulCmd,
	},
}
