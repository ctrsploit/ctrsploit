package file

import (
	"fmt"
	"os"
	"strings"

	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
)

type Contains struct {
	prerequisite.BasePrerequisite
	Target   string
	Expected string
}

func (p *Contains) Check() (bool, error) {
	return p.CheckTemplate(func() (bool, error) {
		content, err := os.ReadFile(p.Target)
		if err != nil {
			return false, fmt.Errorf("failed to check [%s], caused by reading file %s: %w", p.GetName(), p.Target, err)
		}
		p.Satisfied = strings.Contains(string(content), p.Expected)
		return p.Satisfied, p.Err
	})
}

var (
	// HostsContainsNerdctlMarker
	// https://github.com/containerd/nerdctl/blob/v2.1.6/pkg/dnsutil/hostsstore/hosts.go#L42
	// https://github.com/containerd/nerdctl/blob/v2.1.6/pkg/dnsutil/hostsstore/hostsstore.go#L331
	HostsContainsNerdctlMarker = Contains{
		BasePrerequisite: prerequisite.BasePrerequisite{
			Name:   "/etc/hosts",
			Info:   "/etc/hosts contains '<nerdctl>'",
			ExeEnv: exeenv.InContainer,
		},
		Target:   "/etc/hosts",
		Expected: "<nerdctl>",
	}
)
