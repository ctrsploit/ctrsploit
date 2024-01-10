package runc

import (
	"github.com/ctrsploit/ctrsploit/pkg/version/libseccomp"
	"github.com/ctrsploit/ctrsploit/pkg/version/version"
)

type Releaser int

const (
	Unknown Releaser = iota
	GithubRelease
	DockerhubDind
)

type Version struct {
	Number     version.Number
	Url        string
	Releaser   Releaser
	Static     bool
	LibSeccomp libseccomp.Version
	Note       string
}

func (v Version) String() (version string) {
	return v.Number.String()
}
