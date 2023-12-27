package kernel

import "github.com/ctrsploit/ctrsploit/pkg/version/version"

type Version struct {
	Version string
	// TODO: major, minor, ...
	// TODO: rc
}

func NewMap(versions []string) (m version.Map) {
	m = version.Map{}
	for _, v := range versions {
		m[v] = New(v)
	}
	return
}

func New(version string) Version {
	return Version{
		Version: version,
	}
}

func News(versions []string) (result []Version) {
	for _, v := range versions {
		result = append(result, New(v))
	}
	return
}

func (v Version) String() (version string) {
	return v.Version
}
