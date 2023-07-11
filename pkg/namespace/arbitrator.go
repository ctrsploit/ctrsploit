package namespace

import "github.com/pkg/errors"

type Arbitrator interface {
	Arbitrate(namespace Namespace) (isHostNamespace bool, err error)
	PrerequisitesSatisfied() (satisfied bool, err error)
}

var (
	ErrPrerequisiteNotSatisfied = errors.Errorf("not support because of prerequisite not satisfied")
)

func Arbitrate(arbitrator Arbitrator, namespace Namespace) (isHostNamespace bool, err error) {
	satisfied, err := arbitrator.PrerequisitesSatisfied()
	if err != nil {
		return
	}
	if satisfied {
		isHostNamespace, err = arbitrator.Arbitrate(namespace)
	} else {
		err = ErrPrerequisiteNotSatisfied
	}
	return
}

func CheckNamespaceIsHost(arbitrator Arbitrator) (result map[string]bool, err error) {
	result = map[string]bool{}
	satisfied, err := arbitrator.PrerequisitesSatisfied()
	if err != nil {
		return
	}
	namespaces, err := ParseNamespaces()
	if err != nil {
		return
	}
	if satisfied {
		for _, namespace := range namespaces {
			result[namespace.Name], err = arbitrator.Arbitrate(namespace)
			if err != nil {
				err = nil
				continue
			}
		}
	} else {
		err = ErrPrerequisiteNotSatisfied
	}
	return
}
