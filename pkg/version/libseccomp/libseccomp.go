package libseccomp

import "github.com/ctrsploit/ctrsploit/pkg/version/version"

// Version https://man7.org/linux/man-pages/man3/seccomp_version.3.html
type Version struct {
	Major int
	Minor int
	Micro int
}

func NewMap(versions []string) (m version.Map) {
	m = version.Map{}
	for _, v := range versions {
		m[v] = New(v)
	}
	return
}

func New(version string) Version {
	return Version{}
}

func (v Version) String() (version string) {
	return
}
