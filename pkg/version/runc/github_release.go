package runc

import (
	"github.com/ctrsploit/ctrsploit/pkg/version/libseccomp"
	"github.com/ctrsploit/ctrsploit/pkg/version/version"
	"time"
)

var (
	GithubReleaseVersionsLastTimeUpdate, _ = time.Parse(time.RFC3339, "2024-01-10T18:00:00Z08:00")
	GithubReleaseVersions                  = version.Map{
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

		"v1.1.11": Version{
			Url:        "https://github.com/opencontainers/runc/releases/tag/v1.1.11",
			Releaser:   GithubRelease,
			Static:     true,
			LibSeccomp: libseccomp.Versions["v2.5.4"].(libseccomp.Version),
		},

		"v1.1.10": Version{
			Url:        "https://github.com/opencontainers/runc/releases/tag/v1.1.10",
			Releaser:   GithubRelease,
			Static:     true,
			LibSeccomp: libseccomp.Versions["v2.5.4"].(libseccomp.Version),
		},

		"v1.1.9": Version{
			Url:        "https://github.com/opencontainers/runc/releases/tag/v1.1.9",
			Releaser:   GithubRelease,
			Static:     true,
			LibSeccomp: libseccomp.Versions["v2.5.4"].(libseccomp.Version),
		},

		"v1.1.8": Version{
			Url:        "https://github.com/opencontainers/runc/releases/tag/v1.1.8",
			Releaser:   GithubRelease,
			Static:     true,
			LibSeccomp: libseccomp.Versions["v2.5.4"].(libseccomp.Version),
		},

		"v1.1.7": Version{
			Url:        "https://github.com/opencontainers/runc/releases/tag/v1.1.7",
			Releaser:   GithubRelease,
			Static:     true,
			LibSeccomp: libseccomp.Versions["v2.5.4"].(libseccomp.Version),
		},

		"v1.1.6": Version{
			Url:        "https://github.com/opencontainers/runc/releases/tag/v1.1.6",
			Releaser:   GithubRelease,
			Static:     true,
			LibSeccomp: libseccomp.Versions["v2.5.4"].(libseccomp.Version),
		},

		"v1.1.5": Version{
			Url:        "https://github.com/opencontainers/runc/releases/tag/v1.1.5",
			Releaser:   GithubRelease,
			Static:     true,
			LibSeccomp: libseccomp.Versions["v2.5.4"].(libseccomp.Version),
		},

		"v1.1.4": Version{
			Url:        "https://github.com/opencontainers/runc/releases/tag/v1.1.4",
			Releaser:   GithubRelease,
			Static:     true,
			LibSeccomp: libseccomp.Versions["v2.5.4"].(libseccomp.Version),
		},

		"v1.1.3": Version{
			Url:        "https://github.com/opencontainers/runc/releases/tag/v1.1.3",
			Releaser:   GithubRelease,
			Static:     true,
			LibSeccomp: libseccomp.Versions["v2.5.4"].(libseccomp.Version),
		},

		"v1.1.2": Version{
			Url:        "https://github.com/opencontainers/runc/releases/tag/v1.1.2",
			Releaser:   GithubRelease,
			Static:     true,
			LibSeccomp: libseccomp.Versions["v2.5.4"].(libseccomp.Version),
		},

		"v1.1.1": Version{
			Url:        "https://github.com/opencontainers/runc/releases/tag/v1.1.1",
			Releaser:   GithubRelease,
			Static:     true,
			LibSeccomp: libseccomp.Versions["v2.5.4"].(libseccomp.Version),
		},

		"v1.1.0": Version{
			Url:        "https://github.com/opencontainers/runc/releases/tag/v1.1.0",
			Releaser:   GithubRelease,
			Static:     true,
			LibSeccomp: libseccomp.Versions["v2.5.4"].(libseccomp.Version),
		},

		"v1.1.0-rc1": Version{
			Url:        "https://github.com/opencontainers/runc/releases/tag/v1.1.0-rc.1",
			Releaser:   GithubRelease,
			Static:     true,
			LibSeccomp: libseccomp.Versions["v2.5.4"].(libseccomp.Version),
		},

		"v1.0.3": Version{
			Url:        "https://github.com/opencontainers/runc/releases/tag/v1.0.3",
			Releaser:   GithubRelease,
			Static:     true,
			LibSeccomp: libseccomp.Versions["v2.5.4"].(libseccomp.Version),
		},

		"v1.0.2": Version{
			Url:        "https://github.com/opencontainers/runc/releases/tag/v1.0.2",
			Releaser:   GithubRelease,
			Static:     true,
			LibSeccomp: libseccomp.Versions["v2.5.4"].(libseccomp.Version),
		},

		"v1.0.1": Version{
			Url:        "https://github.com/opencontainers/runc/releases/tag/v1.0.1",
			Releaser:   GithubRelease,
			Static:     true,
			LibSeccomp: libseccomp.Versions["v2.5.1"].(libseccomp.Version),
		},

		"v1.0.0": Version{
			Url:        "https://github.com/opencontainers/runc/releases/tag/v1.0.0",
			Releaser:   GithubRelease,
			Static:     true,
			LibSeccomp: libseccomp.Versions["v2.5.1"].(libseccomp.Version),
		},

		"v1.0.0-rc95": Version{
			Url:        "https://github.com/opencontainers/runc/releases/tag/v1.0.0-rc95",
			Releaser:   GithubRelease,
			Static:     true,
			LibSeccomp: libseccomp.Versions["v2.5.1"].(libseccomp.Version),
		},

		"v1.0.0-rc94": Version{
			Url:        "https://github.com/opencontainers/runc/releases/tag/v1.0.0-rc94",
			Releaser:   GithubRelease,
			Static:     true,
			LibSeccomp: libseccomp.Versions["v2.5.1"].(libseccomp.Version),
		},

		"v1.0.0-rc93": Version{
			Url:        "https://github.com/opencontainers/runc/releases/tag/v1.0.0-rc93",
			Releaser:   GithubRelease,
			Static:     true,
			LibSeccomp: libseccomp.Versions["v2.5.1"].(libseccomp.Version),
		},

		"v1.0.0-rc92": Version{
			Url:        "https://github.com/opencontainers/runc/releases/tag/v1.0.0-rc92",
			Releaser:   GithubRelease,
			Static:     true,
			LibSeccomp: libseccomp.Versions["v2.4.1"].(libseccomp.Version),
		},

		"v1.0.0-rc91": Version{
			Url:        "https://github.com/opencontainers/runc/releases/tag/v1.0.0-rc91",
			Releaser:   GithubRelease,
			Static:     true,
			LibSeccomp: libseccomp.Versions["v2.5.1"].(libseccomp.Version),
		},

		"v1.0.0-rc90": Version{
			Url:        "https://github.com/opencontainers/runc/releases/tag/v1.0.0-rc90",
			Releaser:   GithubRelease,
			Static:     true,
			LibSeccomp: libseccomp.Versions["v2.5.1"].(libseccomp.Version),
		},

		"v1.0.0-rc10": Version{
			Url:        "https://github.com/opencontainers/runc/releases/tag/v1.0.0-rc10",
			Releaser:   GithubRelease,
			Static:     true,
			LibSeccomp: libseccomp.Versions["v2.5.1"].(libseccomp.Version),
		},

		"v1.0.0-rc9": Version{
			Url:        "https://github.com/opencontainers/runc/releases/tag/v1.0.0-rc9",
			Releaser:   GithubRelease,
			Static:     true,
			LibSeccomp: libseccomp.Versions["v2.5.1"].(libseccomp.Version),
		},

		"v1.0.0-rc8": Version{
			Url:        "https://github.com/opencontainers/runc/releases/tag/v1.0.0-rc8",
			Releaser:   GithubRelease,
			Static:     true,
			LibSeccomp: libseccomp.Versions["v2.5.1"].(libseccomp.Version),
		},

		"v1.0.0-rc7": Version{
			Url:        "https://github.com/opencontainers/runc/releases/tag/v1.0.0-rc7",
			Releaser:   GithubRelease,
			Static:     true,
			LibSeccomp: libseccomp.Versions["v2.5.1"].(libseccomp.Version),
		},

		"v1.0.0-rc6": Version{
			Url:        "https://github.com/opencontainers/runc/releases/tag/v1.0.0-rc6",
			Releaser:   GithubRelease,
			Static:     true,
			LibSeccomp: libseccomp.Versions["v2.5.1"].(libseccomp.Version),
		},

		"v1.0.0-rc5": Version{
			Url:        "https://github.com/opencontainers/runc/releases/tag/v1.0.0-rc5",
			Releaser:   GithubRelease,
			Static:     true,
			LibSeccomp: libseccomp.Versions["v2.5.1"].(libseccomp.Version),
		},

		"v1.0.0-rc4": Version{
			Url:        "https://github.com/opencontainers/runc/releases/tag/v1.0.0-rc4",
			Releaser:   GithubRelease,
			Static:     true,
			LibSeccomp: libseccomp.Versions["v2.5.1"].(libseccomp.Version),
		},

		"v1.0.0-rc3": Version{
			Url:        "https://github.com/opencontainers/runc/releases/tag/v1.0.0-rc3",
			Releaser:   GithubRelease,
			Static:     true,
			LibSeccomp: libseccomp.Versions["v2.5.1"].(libseccomp.Version),
		},

		"v1.0.0-rc2": Version{
			Url:        "https://github.com/opencontainers/runc/releases/tag/v1.0.0-rc2",
			Releaser:   GithubRelease,
			Static:     true,
			LibSeccomp: libseccomp.Versions["v2.5.1"].(libseccomp.Version),
		},

		"v1.0.0-rc1": Version{
			Url:      "https://github.com/opencontainers/runc/releases/tag/v1.0.0-rc1",
			Releaser: GithubRelease,
			Static:   true,
		},

		"v0.1.1": Version{
			Url:      "https://github.com/opencontainers/runc/releases/tag/v0.1.1",
			Releaser: GithubRelease,
			Static:   true,
		},

		"v0.1.0": Version{
			Url:      "https://github.com/opencontainers/runc/releases/tag/v0.1.0",
			Releaser: GithubRelease,
			Static:   true,
		},

		"v0.0.9": Version{
			Url:      "https://github.com/opencontainers/runc/releases/tag/v0.0.9",
			Releaser: GithubRelease,
			Static:   true,
		},

		"v0.0.8": Version{
			Url:      "https://github.com/opencontainers/runc/releases/tag/v0.0.8",
			Releaser: GithubRelease,
			Static:   true,
		},

		"v0.0.7": Version{
			Url:      "https://github.com/opencontainers/runc/releases/tag/v0.0.7",
			Releaser: GithubRelease,
			Static:   true,
		},

		"v0.0.6": Version{
			Url:      "https://github.com/opencontainers/runc/releases/tag/v0.0.6",
			Releaser: GithubRelease,
			Static:   true,
		},
	}
)
