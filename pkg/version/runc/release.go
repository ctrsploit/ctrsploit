package runc

var (
	Versions = append(GithubReleaseVersions.Values(), DindVersions.Values()...)
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
