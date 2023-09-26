package where

import (
	"bytes"
	"github.com/ctrsploit/ctrsploit/pkg/mountinfo"
	"github.com/ctrsploit/ctrsploit/util"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

const (
	PathDirSecrets = "/var/run/secrets/kubernetes.io"
	// PatternHostname a DNS-1123 subdomain must consist of lower case alphanumeric characters, '-' or '.',
	// and must start and end with an alphanumeric character (e.g. 'example.com',
	// regex used for validation is '[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*')
	PatternHostname = "^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*-[0-9a-f]{10}-[0-9a-z]{5}$"
)

type K8s struct {
	DirSecretsExists             bool
	HostnameMatchPattern         bool
	HostsMountSourceContainsPods bool
	CgroupContainsKubepods       bool
}

func (k *K8s) CheckDirSecretsExists() {
	k.DirSecretsExists = util.CheckPathExists(PathDirSecrets)
}

func (k *K8s) CheckHostsMountSourceContainsPods() (err error) {
	mount, err := mountinfo.HostsMount()
	if err != nil {
		return
	}
	k.HostsMountSourceContainsPods = strings.Contains(mount.Root, "pods")
	return
}

func (k *K8s) CheckCgroup() (err error) {
	content, err := ioutil.ReadFile("/proc/self/cgroup")
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	k.CgroupContainsKubepods = bytes.Contains(content, []byte("kubepods"))
	return
}

func (k *K8s) CheckHostnameMatchPattern() (err error) {
	hostname, err := os.Hostname()
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	k.HostnameMatchPattern, err = regexp.MatchString(PatternHostname, hostname)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	return
}

func (k *K8s) IsIn() (in bool, err error) {
	err = k.CheckHostnameMatchPattern()
	if err != nil {
		return
	}
	if k.HostnameMatchPattern {
		in = true
	}

	k.CheckDirSecretsExists()
	if k.DirSecretsExists {
		in = true
	}

	// don't care this error in production mode
	if k.CheckHostsMountSourceContainsPods() == nil {
		if k.HostsMountSourceContainsPods {
			in = true
		}
	}

	err = k.CheckCgroup()
	if err != nil {
		return
	}
	if k.CgroupContainsKubepods {
		in = true
	}

	return
}
