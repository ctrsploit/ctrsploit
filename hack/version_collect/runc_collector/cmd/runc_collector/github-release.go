package main

import (
	"encoding/json"
	"fmt"
	"github.com/ctrsploit/ctrsploit/hack/version_collect/runc_collector/pkg"
	"github.com/ctrsploit/ctrsploit/pkg/version/libseccomp"
	"github.com/ctrsploit/ctrsploit/pkg/version/runc"
	"github.com/ctrsploit/ctrsploit/pkg/version/version"
	"github.com/google/uuid"
	"github.com/ssst0n3/awesome_libs"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"github.com/urfave/cli/v2"
	"io"
	"net/http"
	"os"
	"strings"
)

type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadUrl string `json:"browser_download_url"`
}

type Release struct {
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
	HtmlUrl string  `json:"html_url"`
}

const (
	template = `
"{.number}": Version{.lbrace}
	Url: "{.url}",
	Releaser: {.releaser},
	Static: {.static},
	LibSeccomp: libseccomp.Versions["{.seccomp}"].(libseccomp.Version),
{.rbrace},`
)

func ReadLibseccompVersion(runcDownloadUrl string) (ver libseccomp.Version, err error) {
	resp, err := http.DefaultClient.Get(runcDownloadUrl)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	path := fmt.Sprintf("/tmp/runc-%s", uuid.New())
	err = os.WriteFile(path, content, 0644)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	defer os.Remove(path)
	ver, err = pkg.ReadLibseccompVersion(path)
	if err != nil {
		return
	}
	return
}

var (
	GithubRelease = &cli.Command{
		Name:    "github-release",
		Aliases: []string{"g"},
		Action: func(context *cli.Context) (err error) {
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

			versionsMap := version.Map{}
			for _, release := range releases {
				seccomp := libseccomp.Version{}
				for _, asset := range release.Assets {
					if strings.Contains(asset.Name, "libseccomp") {
						// do not use libseccomp asset to assert libseccomp version
						// runc v1.0.0-rc92 actually use v2.4.1 instead of v2.4.3
						// seccompVersion := strings.TrimPrefix(asset.Name, "libseccomp-")
						// seccompVersion = strings.TrimSuffix(seccompVersion, ".tar.gz")
						// seccomp = *libseccomp.New(seccompVersion)
						// break
					}
					if asset.Name == "runc.amd64" {
						seccomp, _ = ReadLibseccompVersion(asset.BrowserDownloadUrl)
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
				versionsMap[v.Number.String()] = v
				//log.Logger.Infof("%s: %+v", release.TagName, v)
				//spew.Config.DisableMethods = true
				//spew.Dump(v)
				fmt.Println(awesome_libs.Format(template, awesome_libs.Dict{
					"number":   v.Number,
					"url":      v.Url,
					"releaser": runc.Releaser2Name[v.Releaser],
					"static":   v.Static,
					"seccomp":  v.LibSeccomp.Number,
					"lbrace":   "{",
					"rbrace":   "}",
				}))
			}
			return
		},
	}
)
