package runc

import (
	"fmt"
	"github.com/ctrsploit/ctrsploit/pkg/version/libseccomp"
	"github.com/ctrsploit/ctrsploit/pkg/version/version"
)

type Releaser int

const (
	Unknown Releaser = iota
	GithubRelease
	DockerhubDind
)

var (
	Releaser2Name = map[Releaser]string{
		GithubRelease: "GithubRelease",
		DockerhubDind: "DockerhubDind",
	}
)

func (r Releaser) String() string {
	switch r {
	case GithubRelease:
		return "github-release"
	case DockerhubDind:
		return "dind"
	default:
		return "unknown"
	}
}

type Version struct {
	Number     version.Number
	Url        string
	Releaser   Releaser
	Static     bool
	LibSeccomp libseccomp.Version
	Note       string
}

func (v Version) String() (version string) {
	version = fmt.Sprintf("%s-%s", v.Number, v.Releaser)
	if v.Static {
		version = fmt.Sprintf("%s-static", version)
	}
	return
}
