package lsm

import (
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"io/ioutil"
	"strings"
)

const (
	PathAttrCurrent = "/proc/self/attr/current"
)

func Current() (current string, err error) {
	content, err := ioutil.ReadFile(PathAttrCurrent)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	current = strings.TrimSpace(string(content))
	return
}

func IsEnabled() bool {
	current, err := Current()
	if err != nil {
		return false
	}
	if len(current) > 0 {
		if current != "unconfined" { // TODO: not sure for selinux
			return true
		}
	}
	return false
}
