package version

import "fmt"

type Range []Version
type SoftwareRanges map[TypeSoftware]Range

func (v Range) String() (s string) {
	s = fmt.Sprintf("[%s, %s]", v[0], v[1])
	return
}
