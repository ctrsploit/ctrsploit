package libseccomp

import "github.com/ctrsploit/ctrsploit/pkg/version/version"

// Version https://man7.org/linux/man-pages/man3/seccomp_version.3.html
type Version struct {
	version.Number
}

func NewMap(versions []string) (m version.Map) {
	m = version.Map{}
	for _, v := range versions {
		m[v] = *New(v)
	}
	return
}

func New(v string) *Version {
	return &Version{
		Number: *version.New(v),
	}
}

func (v Version) String() (version string) {
	return
}
