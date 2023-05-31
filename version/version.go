package version

import (
	"github.com/ssst0n3/awesome_libs"
)

// Reference: https://github.com/moby/moby/blob/v24.0.1/dockerversion/version_lib.go
// Default build-time variable for library-import.
// These variables are overridden on build with build-time information.
var (
	ProductName           = "ctrsploit"
	DefaultProductLicense = "public"
	Version               = "library-import"
	GitCommit             = "library-import"
	BuildTime             = "library-import"
	PlatformName          = "library-import"
)

type Ver struct {
	ProductName string
	Version     string
	GitCommit   string
	License     string
	BuildTime   string
}

func DefaultVer() *Ver {
	return &Ver{
		ProductName: ProductName,
		Version:     Version,
		GitCommit:   GitCommit,
		License:     DefaultProductLicense,
		BuildTime:   BuildTime,
	}
}

func (v Ver) String() string {
	return awesome_libs.Format(
		"{.product_name} {.license} version {.version}, build {.git_commit} at {.build_time}",
		awesome_libs.Dict{
			"product_name": v.ProductName,
			"license":      v.License,
			"version":      v.Version,
			"git_commit":   v.GitCommit,
			"build_time":   v.BuildTime,
		})
}
