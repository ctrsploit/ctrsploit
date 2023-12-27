package libseccomp

// Version https://man7.org/linux/man-pages/man3/seccomp_version.3.html
type Version struct {
	Major int
	Minor int
	Micro int
}

func New(version string) Version {
	return Version{}
}

type Map map[string]Version
type Slice []Version

var (
	Versions = Map{
		"2.3.3": New("2.3.3"),
		"2.4.1": New("2.4.1"),
		"2.5.1": New("2.5.1"),
	}
	Order = Slice{
		Versions["2.3.3"],
		Versions["2.4.1"],
		Versions["2.5.1"],
	}
)
