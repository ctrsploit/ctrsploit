package runc

import (
	"github.com/ctrsploit/ctrsploit/pkg/version/libseccomp"
	"github.com/ctrsploit/ctrsploit/pkg/version/version"
)

var (
	GithubReleaseVersions = version.Map{
		"1.0.0-rc92": Version{
			Url:        "",
			Releaser:   GithubRelease,
			Static:     true,
			LibSeccomp: libseccomp.Versions["2.4.1"].(libseccomp.Version),
		},
		"1.0.0-rc93": Version{
			Url:        "https://github.com/opencontainers/runc/releases/download/v1.0.0-rc93/runc.amd64",
			Releaser:   GithubRelease,
			Static:     true,
			LibSeccomp: libseccomp.Versions["2.5.1"].(libseccomp.Version),
		},
	}
)
