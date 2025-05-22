package cgroup

import (
	"bufio"
	"fmt"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"os"
	"regexp"
)

func GetContainerID() (id string, err error) {
	file, err := os.Open("/proc/self/cgroup")
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	re := regexp.MustCompile(`[0-9a-f]{64}`)
	for scanner.Scan() {
		line := scanner.Text()
		if match := re.FindString(line); match != "" {
			id = match
			return
		}
	}

	if err = scanner.Err(); err != nil {
		awesome_error.CheckErr(err)
		return
	}
	err = fmt.Errorf("cannot find cgroup container")
	awesome_error.CheckWarning(err)
	return
}
