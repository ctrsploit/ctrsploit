package namespace

import "github.com/pkg/errors"

type Arbitrator interface {
	Arbitrate(namespace Namespace) (namespaceLevel Level, err error)
	PrerequisitesSatisfied() (satisfied bool, err error)
}

var (
	ErrPrerequisiteNotSatisfied = errors.Errorf("not support because of prerequisite not satisfied")
)

func Arbitrate(arbitrator Arbitrator, namespace Namespace) (namespaceLevel Level, err error) {
	satisfied, err := arbitrator.PrerequisitesSatisfied()
	if err != nil {
		return
	}
	if satisfied {
		namespaceLevel, err = arbitrator.Arbitrate(namespace)
	} else {
		err = ErrPrerequisiteNotSatisfied
	}
	return
}

func CheckNamespaceLevel(arbitrator Arbitrator) (result map[string]Level, names []string, err error) {
	result = map[string]Level{}
	satisfied, err := arbitrator.PrerequisitesSatisfied()
	if err != nil {
		return
	}
	namespaces, names, err := ParseNamespaces()
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
