package env

import (
	"context"

	"github.com/ctrsploit/ctrsploit/env/cpusec"
	"github.com/urfave/cli/v3"
)

var Cpusec = &cli.Command{
	Name:    cpusec.CommandName,
	Aliases: []string{"cs"},
	Usage:   "show CPU/kernel security mitigations (SMEP/SMAP/KPTI/IBT/KCFI/FG-KASLR on x86; PAC/BTI/KPTI/PAN/MTE on arm64)",
	Action: func(ctx context.Context, cmd *cli.Command) (err error) {
		err = cpusec.Print()
		if err != nil {
			return
		}
		return
	},
}
