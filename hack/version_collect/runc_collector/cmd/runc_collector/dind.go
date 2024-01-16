package main

import (
	"fmt"
	"github.com/ctrsploit/ctrsploit/hack/version_collect/runc_collector/pkg"
	"github.com/ctrsploit/ctrsploit/pkg/version/libseccomp"
	"github.com/ctrsploit/sploit-spec/pkg/log"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"os"
	"strings"
)

func ImageLibseccompVersion(image string) (ver libseccomp.Version, err error) {
	path, e := pkg.GetFileHostPath(image, "/usr/local/bin/runc")
	if e != nil {
		if os.IsNotExist(e) {
			log.Logger.Warn("/usr/local/bin/runc not exists")
		} else {
			err = e
			awesome_error.CheckErr(err)
		}
		return
	}
	ver, err = pkg.ReadLibseccompVersion(path)
	return
}

func Dind() (err error) {
	tags, err := pkg.ListTags("docker")
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	for _, tag := range tags {
		if !strings.Contains(tag, "dind") {
			continue
		}
		image := fmt.Sprintf("docker:%s", tag)
		fmt.Println(image)
		err = pkg.Pull(image)
		if err != nil {
			return
		}
		ver, _ := ImageLibseccompVersion(image)
		fmt.Printf("%+v\n", ver)
	}
	return
}
