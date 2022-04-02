package splice

import (
	"github.com/docker/docker/pkg/reexec"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"os"
)

func ReexecRegister(splice Splice) {
	{
		reexec.Register(escalateExpName(splice), func() {
			awesome_error.CheckFatal(Escalate(splice))
		})
	}
	{
		reexec.Register(escapeExpName(splice), func() {
			awesome_error.CheckFatal(Escape(splice))
		})
	}
}

func InvokeEscalate(splice Splice) error {
	cmd := reexec.Command(escalateExpName(splice))
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func InvokeEscape(splice Splice) error {
	cmd := reexec.Command(escapeExpName(splice))
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
