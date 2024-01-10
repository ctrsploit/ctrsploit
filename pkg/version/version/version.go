package version

import (
	"fmt"
	"strconv"
	"strings"
)

type Version interface {
	String() string
}

type Number struct {
	Major int
	Minor int
	Patch int
	Rc    int
	Beta  int
}

func New(number string) *Number {
	number = strings.TrimPrefix(number, "v")
	parts := strings.SplitN(number, "-", 2)
	part0 := parts[0]
	main := strings.Split(part0, ".")
	major, _ := strconv.Atoi(main[0])
	minor, _ := strconv.Atoi(main[1])
	patch, _ := strconv.Atoi(main[2])
	rc := -1
	beta := -1
	if len(parts) > 1 {
		ext := parts[1]
		if strings.Contains(ext, "beta") {
			ext = strings.TrimPrefix(ext, "beta")
			ext = strings.TrimPrefix(ext, ".")
			beta, _ = strconv.Atoi(ext)
		}
		if strings.Contains(ext, "rc") {
			ext = strings.TrimPrefix(ext, "rc")
			ext = strings.TrimPrefix(ext, ".")
			rc, _ = strconv.Atoi(ext)
		}
	}

	return &Number{
		Major: major,
		Minor: minor,
		Patch: patch,
		Rc:    rc,
		Beta:  beta,
	}
}

func (n Number) String() (s string) {
	s = fmt.Sprintf("v%d.%d.%d", n.Major, n.Minor, n.Patch)
	if n.Rc != -1 {
		s = fmt.Sprintf("%s-rc%d", s, n.Rc)
	}
	if n.Beta != -1 {
		s = fmt.Sprintf("%s-rc%d", s, n.Beta)
	}
	return
}
