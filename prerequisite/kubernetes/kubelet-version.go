package kubernetes

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ctrsploit/ctrsploit/pkg/version/kubelet"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
	"github.com/ssst0n3/awesome_libs/awesome_error"
)

type AllKubeletVersionConstraint struct {
	prerequisite.BasePrerequisite
	Constraint      string
	matchedNodes    []string
	notMatchedNodes []string
}

// Check checks if any kubelet version satisfies the constraint.
func (p *AllKubeletVersionConstraint) Check() (bool, error) {
	return p.CheckTemplate(func() (bool, error) {
		p.Satisfied = true
		cons, err := semver.NewConstraint(p.Constraint)
		if err != nil {
			err = fmt.Errorf("failed to parse constraint %s: %w", p.Constraint, err)
			// it's fatal if the constraint is invalid, it means the code is wrong
			awesome_error.CheckFatal(err)
		}
		versions, err := kubelet.Versions()
		if err != nil {
			p.Err = err
			return false, p.Err
		}
		if len(versions) == 0 {
			p.Err = fmt.Errorf("no kubelet versions found")
			return false, p.Err
		}
		for name, version := range versions {
			if cons.Check(version) {
				p.matchedNodes = append(p.matchedNodes, name)
			} else {
				p.notMatchedNodes = append(p.notMatchedNodes, name)
				p.Satisfied = false
			}
		}
		return p.Satisfied, p.Err
	})
}

func GetNodes(s prerequisite.Set) ([]string, error) {
	_, _ = s.Check()
	switch s.(type) {
	case *AllKubeletVersionConstraint:
		return s.(*AllKubeletVersionConstraint).matchedNodes, nil
	case *prerequisite.SetNot:
		if p, ok := s.(*prerequisite.SetNot).Set.(*AllKubeletVersionConstraint); ok {
			return p.notMatchedNodes, nil
		} else {
			return nil, fmt.Errorf("unsupported prerequisite type inside Not: %T", s.(*prerequisite.SetNot).Set)
		}
	default:
		return nil, fmt.Errorf("unsupported prerequisite type: %T", s)
	}
}

//goland:noinspection GoSnakeCaseUsage
var (
	ConstraintKubeletFixedCVE_2020_8558 = ">= 1.18.4 || >= 1.17.7, < 1.18.0 || >= 1.16.11, < 1.17.0"
	FixedCVE_2020_8558                  = AllKubeletVersionConstraint{
		BasePrerequisite: prerequisite.BasePrerequisite{
			Name:   "cve-2020-8558-fixed",
			Info:   "kubelet >= 1.18.4, 1.17.7, 1.16.11",
			ExeEnv: exeenv.InHost | exeenv.InContainer,
		},
		Constraint: ConstraintKubeletFixedCVE_2020_8558,
	}
	MaybeVulnerableToCVE_2020_8558 = prerequisite.Not(&FixedCVE_2020_8558)

	ConstraintKubeletFixedCVE_2021_25741 = ">= 1.22.2 || >= 1.21.5, < 1.22.0 || >= 1.20.11, < 1.21.0 || >= 1.19.15, < 1.20.0"
	FixedCVE_2021_25741                  = AllKubeletVersionConstraint{
		BasePrerequisite: prerequisite.BasePrerequisite{
			Name: "cve-2021-25741-fixed",
			// https://github.com/kubernetes/kubernetes/issues/104980
			Info:   "kubelet >= 1.22.2, 1.21.5, 1.20.11, 1.19.15",
			ExeEnv: exeenv.InHost | exeenv.InContainer,
		},
		Constraint: ConstraintKubeletFixedCVE_2021_25741,
	}
	MaybeVulnerableToCVE_2021_25741 = prerequisite.Not(&FixedCVE_2021_25741)
)
