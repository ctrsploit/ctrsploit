package noused

import (
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"os/user"
)

func Whoami() (uid string, err error) {
	i, err := user.Current()
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	uid = i.Uid
	return
}

func AmIRoot() (result bool, err error) {
	i, err := Whoami()
	if err != nil {
		return
	}
	result = i == "0"
	return
}
