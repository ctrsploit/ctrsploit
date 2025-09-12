package namespace

import (
	"path/filepath"

	"github.com/ctrsploit/sploit-spec/pkg/env/container"
)

type Namespace struct {
	Name            string
	Path            string
	Type            container.NamespaceType
	InodeNumber     int
	InitInodeNumber int
}

func ParseNamespaces() (namespaces []Namespace, names []string, err error) {
	proc := "/proc/self/ns"
	namespaceInoMap, names, err := ListNamespaceDir(proc)
	if err != nil {
		return
	}
	for _, name := range names {
		ino := namespaceInoMap[name]
		namespace := Namespace{
			Name:            name,
			Path:            filepath.Join(proc, name),
			Type:            container.NamespaceMapName2Type[name],
			InodeNumber:     ino,
			InitInodeNumber: InitInoMap[container.NamespaceMapName2Type[name]],
		}
		namespaces = append(namespaces, namespace)
	}
	return
}
