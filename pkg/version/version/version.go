package version

import "fmt"

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
