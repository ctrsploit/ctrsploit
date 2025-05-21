package vul

import (
	cve_2020_15257 "github.com/ctrsploit/ctrsploit/vul/cve-2020-15257"
	"github.com/ctrsploit/sploit-spec/pkg/app"
	"github.com/urfave/cli/v2"
)

var Command = &cli.Command{
	Name:    "vul",
	Aliases: []string{"v"},
	Usage:   "check security inside a container",
	Subcommands: []*cli.Command{
		app.Vul2VulCmd(&cve_2020_15257.Vul, []string{"15257"}, nil, nil),
	},
}
