package runc

import (
	"github.com/ctrsploit/ctrsploit/pkg/version/libseccomp"
	"github.com/ctrsploit/ctrsploit/pkg/version/version"
)

var (
	Versions              = append(GithubReleaseVersions.Values(), DindVersions.Values()...)
	GithubReleaseVersions = version.Map{
		"1.0.0-rc92-github_release": Version{
			Url:        "",
			Releaser:   GithubRelease,
			Static:     true,
			LibSeccomp: libseccomp.Versions["2.4.1"].(libseccomp.Version),
		},
		"1.0.0-rc93-github_release": Version{
			Url:        "https://github.com/opencontainers/runc/releases/download/v1.0.0-rc93/runc.amd64",
			Releaser:   GithubRelease,
			Static:     true,
			LibSeccomp: libseccomp.Versions["2.5.1"].(libseccomp.Version),
		},
	}
	DindVersions = version.Map{
		"1.0.0-rc92-dind": Version{
			Url:        "",
			Releaser:   DockerhubDind,
			Static:     true,
			LibSeccomp: libseccomp.Versions["2.3.3"].(libseccomp.Version),
			Note:       "",
		},
		"1.0.0-rc93-dind": Version{
			Url:        "",
			Releaser:   DockerhubDind,
			Static:     true,
			LibSeccomp: libseccomp.Versions["2.3.3"].(libseccomp.Version),
			Note: `
docker:20.10.4-dind
docker:20.10.5-dind`,
		},
		"1.0.0-rc93_docker:20.10.6-dind": Version{
			Url:        "",
			Releaser:   DockerhubDind,
			Static:     true,
			LibSeccomp: libseccomp.Versions["2.4.4"].(libseccomp.Version),
			Note:       `docker:20.10.6-dind`,
		},
		"1.1.4-dind": Version{
			Url:        "",
			Releaser:   DockerhubDind,
			Static:     true,
			LibSeccomp: libseccomp.Versions["2.5.1"].(libseccomp.Version),
			Note: `
docker:20.10.23-dind
docker:23.0.0-rc.4-dind
docker:23.0.0-dind
docker:23.0.1-dind
docker:23.0.2-dind`,
		},
		"1.1.5-dind": Version{
			Url:        "",
			Releaser:   DockerhubDind,
			Static:     true,
			LibSeccomp: libseccomp.Versions["2.5.1"].(libseccomp.Version),
			Note: `
docker:20.10.24-dind
docker:23.0.3-dind
docker:23.0.4-dind
docker:23.0.5-dind
docker:24.0.0-beta.1-dind
docker:24.0.0-beta.2-dind`,
		},
		"1.1.6-dind": Version{
			Url:        "",
			Releaser:   DockerhubDind,
			Static:     true,
			LibSeccomp: libseccomp.Versions["2.5.1"].(libseccomp.Version),
			Note:       `docker:24.0.0-rc.1-dind`,
		},
		"1.1.7-dind": Version{
			Url:        "",
			Releaser:   DockerhubDind,
			Static:     true,
			LibSeccomp: libseccomp.Versions["2.5.1"].(libseccomp.Version),
			Note: `
docker:23.0.6-dind
docker:24.0.0-rc.2-dind
docker:24.0.0-rc.3-dind
docker:24.0.0-rc.4-dind
docker:24.0.0-dind
docker:24.0.1-dind
docker:24.0.2-dind
docker:24.0.3-dind
docker:24.0.4-dind`,
		},
		"1.1.8-dind": Version{
			Url:        "",
			Releaser:   DockerhubDind,
			Static:     true,
			LibSeccomp: libseccomp.Versions["2.5.1"].(libseccomp.Version),
			Note:       `docker:24.0.5-dind`,
		},
		"1.1.9-dind": Version{
			Url:        "",
			Releaser:   DockerhubDind,
			Static:     true,
			LibSeccomp: libseccomp.Versions["2.5.1"].(libseccomp.Version),
			Note: `
docker:24.0.6-dind
docker:24.0.7-dind`,
		},
		"1.1.10-dind": Version{
			Url:        "",
			Releaser:   DockerhubDind,
			Static:     true,
			LibSeccomp: libseccomp.Versions["2.5.1"].(libseccomp.Version),
			Note:       `docker:25.0.0-beta.1`,
		},
	}
)

var (
	// StaticBeforeSupportEnosys <= 1.0.0-rc92
	StaticBeforeSupportEnosys = append(
		GithubReleaseVersions.Get([]string{
			"1.0.0-rc92-github_release",
		}),
		DindVersions.Get([]string{
			"1.0.0-rc92-dind",
		})...,
	)
)
