package release_agent

import (
	"flag"
	"os"

	"github.com/moby/sys/reexec"
	"github.com/ssst0n3/awesome_libs/awesome_error"
)

const (
	CmdName = "release_agent"
)

func init() {
	reexec.Register(CmdName, func() {
		flagSet := flag.NewFlagSet(CmdName, flag.ContinueOnError)
		var path string
		flagSet.StringVar(&path, "path", "", "")
		awesome_error.CheckFatal(flagSet.Parse(os.Args[1:]))
		err := ReleaseAgent(path)
		awesome_error.CheckErr(err)
	})
}
