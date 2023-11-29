package version

import (
	"fmt"
	"strconv"
	"strings"
)

type Docker struct {
	Major  int
	Minor  int
	Patch  int
	IsBeta bool
	Beta   int
}

func NewDocker(version string) Docker {
	parts := strings.SplitN(version, "-", 2)
	main := parts[0]
	part0 := strings.Split(main, ".")
	major, _ := strconv.Atoi(part0[0])
	minor, _ := strconv.Atoi(part0[1])
	patch, _ := strconv.Atoi(part0[2])

	var isBeta bool
	var beta int
	if len(parts) > 1 {
		ext := parts[1]
		if strings.Contains(ext, "beta.") {
			isBeta = true
			beta, _ = strconv.Atoi(ext[5:])
		}
	}

	return Docker{
		Major:  major,
		Minor:  minor,
		Patch:  patch,
		IsBeta: isBeta,
		Beta:   beta,
	}
}

func (v Docker) String() (version string) {
	version = fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	if v.IsBeta {
		version += fmt.Sprintf("-beta.%d", v.Beta)
	}
	return
}
