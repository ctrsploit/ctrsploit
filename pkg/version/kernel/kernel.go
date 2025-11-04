package kernel

import "github.com/ctrsploit/ctrsploit/pkg/version/version"

type VersionOld struct {
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

func New(version string) VersionOld {
	return VersionOld{
		Version: version,
	}
}

func News(versions []string) (result []VersionOld) {
	for _, v := range versions {
		result = append(result, New(v))
	}
	return
}

func (v VersionOld) String() (version string) {
	return v.Version
}
