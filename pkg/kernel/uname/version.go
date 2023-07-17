package uname

import (
	"strings"
)

func VersionEqual(v1, v2 string) bool {
	c1, v1, f1 := strings.Cut(v1, ".")
	c2, v2, f2 := strings.Cut(v2, ".")
	if c1 == "" {
		c1 = "0"
	}
	if c2 == "" {
		c2 = "0"
	}
	equal := c1 == c2
	if f1 || f2 {
		equal = equal && VersionEqual(v1, v2)
	}
	return equal
}
