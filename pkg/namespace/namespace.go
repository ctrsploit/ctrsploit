package namespace

import "path/filepath"

type Namespace struct {
	Name            string
	Path            string
	Type            Type
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
			Type:            MapName2Type[name],
			InodeNumber:     ino,
			InitInodeNumber: InitInoMap[MapName2Type[name]],
		}
		namespaces = append(namespaces, namespace)
	}
	return
}
