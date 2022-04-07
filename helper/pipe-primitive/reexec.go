package pipe_primitive

import (
	"github.com/docker/docker/pkg/reexec"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"os"
)

func ReexecRegister(primitive Primitive) {
	{
		reexec.Register(escalateExpName(primitive), func() {
			awesome_error.CheckFatal(Escalate(primitive))
		})
	}
	{
		reexec.Register(escapeExpName(primitive), func() {
			awesome_error.CheckFatal(Escape(primitive))
		})
	}
	{
		reexec.Register(imagePollutionExpName(primitive), func() {
			awesome_error.CheckFatal(ImagePollution(primitive))
		})
	}
}

func InvokeEscalate(primitive Primitive) error {
	cmd := reexec.Command(escalateExpName(primitive))
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func InvokeEscape(primitive Primitive) error {
	cmd := reexec.Command(escapeExpName(primitive))
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func InvokeImagePollution(primitive Primitive, source, dest string) error {
	cmd := reexec.Command(imagePollutionExpName(primitive), "--source", source, "--destination", dest)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
