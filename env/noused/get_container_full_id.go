package noused

import (
	"errors"
	"fmt"
	"github.com/ctrsploit/ctrsploit/internal"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"strings"
)

func GetContainerFullId() (ID string, err error) {
	cgroups, err := internal.ParseCgroup("/proc/self/cgroup")
	if err != nil {
		return
	}
	cgroupName := string(cgroups[0].Name)
	if strings.HasPrefix(cgroupName, "/docker/") {
		ID = cgroupName[len("/docker/"):]
	} else {
		err = errors.New(fmt.Sprintf("there's no /docker/ in cgroup: %v", cgroups[0]))
		awesome_error.CheckErr(err)
		return
	}
	return
}
