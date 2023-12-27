package runc

import "github.com/ctrsploit/ctrsploit/pkg/version/libseccomp"

type Releaser int

const (
	Unknown Releaser = iota
	GithubRelease
	DockerhubDind
)

type Version struct {
	Url        string
	Releaser   Releaser
	Static     bool
	LibSeccomp libseccomp.Version
	Note       string
}

func (v Version) String() (version string) {
	return
}

type Versions map[string]Version
