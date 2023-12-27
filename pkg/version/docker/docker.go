package docker

import (
	"fmt"
	"github.com/ctrsploit/ctrsploit/pkg/version/version"
	"strconv"
	"strings"
)

type Version struct {
	Major  int
	Minor  int
	Patch  int
	IsBeta bool
	Beta   int
}

func NewMap(versions []string) (m version.Map) {
	m = version.Map{}
	for _, v := range versions {
		m[v] = New(v)
	}
	return
}

func New(version string) Version {
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
		if strings.Contains(ext, "beta") {
			isBeta = true
			ext = strings.TrimPrefix(ext, "beta")
			ext = strings.TrimPrefix(ext, ".")
			beta, _ = strconv.Atoi(ext)
		}
	}

	return Version{
		Major:  major,
		Minor:  minor,
		Patch:  patch,
		IsBeta: isBeta,
		Beta:   beta,
	}
}

func (v Version) String() (version string) {
	version = fmt.Sprintf("%d.%s.%d", v.Major, v.minor(), v.Patch)
	if v.IsBeta {
		version += v.beta()
	}
	return
}

/*
v0.1.0
...
v1.13.1-rc2
v17.03.0-ce
...
v22.06.0-beta.0
v23.0.0
...
*/
func (v Version) minor() (version string) {
	if v.Minor == 0 {
		return "0"
	}
	if v.Major < 17 {
		return fmt.Sprintf("%d", v.Minor)
	}
	version = fmt.Sprintf("%02d", v.Minor)
	return
}

/*
beta:

$ git --no-pager tag | grep beta
v18.09.0-beta3
...
v20.10.0-beta1
v22.06.0-beta.0
...
*/
func (v Version) beta() (version string) {
	if v.IsBeta {
		point := ""
		if v.Major >= 22 {
			point = "."
		}
		version += fmt.Sprintf("-beta%s%d", point, v.Beta)
	}
	return
}
