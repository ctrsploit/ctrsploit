package internal

import (
	"fmt"
	"github.com/ctrsploit/ctrsploit/vul"
	"github.com/urfave/cli/v2"
	"k8s.io/apimachinery/pkg/util/json"
)

func Vul2cmd(v vul.Vulnerability) *cli.Command {
	return &cli.Command{
		Name:  v.GetName(),
		Usage: v.GetDescription(),
		Action: func(context *cli.Context) (err error) {
			_, err = v.CheckSec()
			if err != nil {
				return
			}
			output, err := json.Marshal(v)
			if err != nil {
				return
			}
			fmt.Println(string(output))
			return
		},
	}
}
