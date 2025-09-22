package bash

import (
	"testing"

	"github.com/urfave/cli/v2"
)

func Test_vulnerability_Exploit(t *testing.T) {
	f := VulCmd.Subcommands[1].Action
	f(&cli.Context{})
}
