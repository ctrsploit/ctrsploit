package util

import (
	"bytes"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"io/ioutil"
)

type Cgroup struct {
	Id       []byte
	Resource []byte
	Name     []byte
}

func ParseCgroup(filepath string) (cgroups []Cgroup, err error) {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	for _, line := range bytes.Split(content, []byte{'\n'}) {
		item := bytes.Split(line, []byte{':'})
		if len(item) != 3 {
			continue
		}
		cgroups = append(cgroups, Cgroup{
			Id:       item[0],
			Resource: item[1],
			Name:     item[2],
		})
	}
	return
}
