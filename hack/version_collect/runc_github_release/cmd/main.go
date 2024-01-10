package main

import (
	"github.com/ctrsploit/ctrsploit/pkg/version/libseccomp"
	"github.com/ctrsploit/ctrsploit/pkg/version/runc"
	"github.com/ctrsploit/ctrsploit/pkg/version/version"
	"github.com/ctrsploit/sploit-spec/pkg/log"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"io"
	"k8s.io/apimachinery/pkg/util/json"
	"net/http"
	"strings"
)

type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadUrl string `json:"browserDownloadUrl"`
}

type Release struct {
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
	HtmlUrl string  `json:"html_url"`
}

func main() {
	resp, err := http.DefaultClient.Get("https://api.github.com/repos/opencontainers/runc/releases?per_page=1000")
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	var releases []Release
	err = json.Unmarshal(content, &releases)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	for _, release := range releases {
		seccomp := libseccomp.Version{}
		for _, asset := range release.Assets {
			if strings.Contains(asset.Name, "libseccomp") {
				seccompVersion := strings.TrimPrefix(asset.Name, "libseccomp-")
				seccompVersion = strings.TrimSuffix(asset.Name, ".tar.gz")
				seccomp = *libseccomp.New(seccompVersion)
				break
			}
		}
		v := runc.Version{
			Number:     *version.New(release.TagName),
			Url:        release.HtmlUrl,
			Releaser:   runc.GithubRelease,
			Static:     true,
			LibSeccomp: seccomp,
		}
		log.Logger.Infof("%s: %+v", release.TagName, v)
	}
}
