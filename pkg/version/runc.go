package version

import "github.com/ctrsploit/ctrsploit/pkg/version/kernel"

type Releaser int

const (
	RuncReleaserUnknown = iota
	RuncReleaserGithubRelease
)

type Runc struct {
	BinaryUrl string
	Releaser  Releaser
	Static    bool
}

func NewRunc(version string, static bool) kernel.Version {
	return kernel.Version{}
}

func (v Runc) String() (version string) {
	return
}

var (
	Runc_1_0_0_rc93_GithubRelease_Static = Runc{
		BinaryUrl: "https://github.com/opencontainers/runc/releases/download/v1.0.0-rc93/runc.amd64",
		Releaser:  RuncReleaserGithubRelease,
		Static:    true,
	}
)
