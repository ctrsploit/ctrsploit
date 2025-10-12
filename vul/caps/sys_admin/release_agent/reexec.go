package release_agent

import (
	"flag"
	"os"

	"github.com/ctrsploit/sploit-spec/pkg/log"
	"github.com/moby/sys/reexec"
	"github.com/ssst0n3/awesome_libs/awesome_error"
)

const (
	CmdName = "release_agent"
)

func init() {
	reexec.Register(CmdName, func() {
		flagSet := flag.NewFlagSet(CmdName, flag.ContinueOnError)
		var cmd string
		flagSet.StringVar(&cmd, "cmd", "", "")
		awesome_error.CheckFatal(flagSet.Parse(os.Args[1:]))
		result, err := Exploit(cmd)
		awesome_error.CheckErr(err)
		log.Logger.Infof("result:\n%s", result)
	})
}
