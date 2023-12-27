package runc

import "github.com/ctrsploit/ctrsploit/pkg/version/libseccomp"

var (
	GithubReleaseVersions = Versions{
		"1.0.0-rc92": {
			Url:        "",
			Releaser:   GithubRelease,
			Static:     true,
			LibSeccomp: libseccomp.Versions["2.4.1"],
		},
		"1.0.0-rc93": {
			Url:        "https://github.com/opencontainers/runc/releases/download/v1.0.0-rc93/runc.amd64",
			Releaser:   GithubRelease,
			Static:     true,
			LibSeccomp: libseccomp.Versions["2.5.1"],
		},
	}
	DindVersions = Versions{
		"1.0.0-rc92": {
			Url:        "",
			Releaser:   DockerhubDind,
			Static:     true,
			LibSeccomp: libseccomp.Versions["2.3.3"],
			Note:       "",
		},
		"1.0.0-rc93": {
			Url:        "",
			Releaser:   DockerhubDind,
			Static:     true,
			LibSeccomp: libseccomp.Versions["2.3.3"],
			Note: `
docker:20.10.4-dind
docker:20.10.5-dind`,
		},
		"1.0.0-rc93_docker:20.10.6-dind": {
			Url:        "",
			Releaser:   DockerhubDind,
			Static:     true,
			LibSeccomp: libseccomp.Versions["2.4.4"],
			Note:       `docker:20.10.6-dind`,
		},
		"1.1.4": {
			Url:        "",
			Releaser:   DockerhubDind,
			Static:     true,
			LibSeccomp: libseccomp.Versions["2.5.1"],
			Note: `
docker:20.10.23-dind
docker:23.0.0-rc.4-dind
docker:23.0.0-dind
docker:23.0.1-dind
docker:23.0.2-dind`,
		},
		"1.1.5": {
			Url:        "",
			Releaser:   DockerhubDind,
			Static:     true,
			LibSeccomp: libseccomp.Versions["2.5.1"],
			Note: `
docker:20.10.24-dind
docker:23.0.3-dind
docker:23.0.4-dind
docker:23.0.5-dind
docker:24.0.0-beta.1-dind
docker:24.0.0-beta.2-dind`,
		},
		"1.1.6": {
			Url:        "",
			Releaser:   DockerhubDind,
			Static:     true,
			LibSeccomp: libseccomp.Versions["2.5.1"],
			Note:       `docker:24.0.0-rc.1-dind`,
		},
		"1.1.7": {
			Url:        "",
			Releaser:   DockerhubDind,
			Static:     true,
			LibSeccomp: libseccomp.Versions["2.5.1"],
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
		"1.1.8": {
			Url:        "",
			Releaser:   DockerhubDind,
			Static:     true,
			LibSeccomp: libseccomp.Versions["2.5.1"],
			Note:       `docker:24.0.5-dind`,
		},
		"1.1.9": {
			Url:        "",
			Releaser:   DockerhubDind,
			Static:     true,
			LibSeccomp: libseccomp.Versions["2.5.1"],
			Note: `
docker:24.0.6-dind
docker:24.0.7-dind`,
		},
		"1.1.10": {
			Url:        "",
			Releaser:   DockerhubDind,
			Static:     true,
			LibSeccomp: libseccomp.Versions["2.5.1"],
			Note:       `docker:25.0.0-beta.1`,
		},
	}
)
