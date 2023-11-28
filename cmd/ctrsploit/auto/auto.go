package auto

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

const (
	CommandNameAuto = "auto"
)

var (
	Command = &cli.Command{
		Name:    CommandNameAuto,
		Usage:   "auto gathering information, detect vulnerabilities and run exploits",
		Aliases: []string{"a"},
		Action: func(context *cli.Context) (err error) {
			fmt.Println("TODO")
			return
		},
	}
)
