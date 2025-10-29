package iptables

import (
	"fmt"

	"github.com/coreos/go-iptables/iptables"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
)

type RuleExists struct {
	prerequisite.BasePrerequisite
	rule []string
}

func (p *RuleExists) Check() (bool, error) {
	return p.CheckTemplate(func() (bool, error) {
		ipt, err := iptables.New()
		if err != nil {
			p.Err = fmt.Errorf("failed to check [%s], caused by init iptables: %w", p.GetName(), err)
			return p.Satisfied, p.Err
		}

		ok, err := ipt.Exists("filter", "KUBE-FIREWALL", p.rule...)
		if err != nil {
			p.Err = fmt.Errorf("failed to check [%s], caused by checking rule '%s': %w", p.GetName(), p.rule, err)
			return p.Satisfied, p.Err
		}
		p.Satisfied = ok
		return p.Satisfied, p.Err
	})
}

var (
	KubeletPR91569 = RuleExists{
		BasePrerequisite: prerequisite.BasePrerequisite{
			Name:   "iptables drop all non-localnet packets",
			Info:   "Kubelet PR#91569 iptables rule exists in KUBE-FIREWALL chain",
			ExeEnv: exeenv.InHost,
		},
		rule: []string{
			"!", "-s", "127.0.0.0/8",
			"-d", "127.0.0.0/8",
			"-m", "comment", "--comment", "block incoming localnet connections",
			"-m", "conntrack", "!", "--ctstate", "RELATED,ESTABLISHED,DNAT",
			"-j", "DROP",
		},
	}
	KubeletPR91569Missing = prerequisite.Not(&KubeletPR91569)
)
