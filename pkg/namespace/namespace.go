package namespace

import "path/filepath"

type Namespace struct {
	Name            string
	Path            string
	Type            int
	InodeNumber     int
	InitInodeNumber int
}

const (
	TypeNamespaceLevelUnknown = iota
	TypeNamespaceLevelBoot
	TypeNamespaceLevelChild
	TypeNamespaceLevelHost = TypeNamespaceLevelBoot
)

const (
	TypeNamespaceTypeUnknown = iota
	TypeNamespaceTypeIPC
	TypeNamespaceTypeUTS
	TypeNamespaceTypeUser
	TypeNamespaceTypePid
	TypeNamespaceTypeCGroup
	TypeNamespaceTypeTime
	TypeNamespaceTypeMount
	TypeNamespaceTypeNetwork
)

const (
	// ProcDynamicFirst https://git.kernel.org/pub/scm/linux/kernel/git/stable/linux.git/tree/fs/proc/generic.c?h=v5.11.11#n201
	ProcDynamicFirst = 0xF0000000
	// ProcIpcInitIno https://git.kernel.org/pub/scm/linux/kernel/git/stable/linux.git/tree/include/linux/proc_ns.h?h=v5.6#n43
	ProcIpcInitIno    = 0xEFFFFFFF
	ProcUtsInitIno    = 0xEFFFFFFE
	ProcUserInitIno   = 0xEFFFFFFD
	ProcPidInitIno    = 0xEFFFFFFC
	ProcCGroupInitIno = 0xEFFFFFFB
	ProcTimeInitIno   = 0xEFFFFFFA
	// ProcMountInitIno mount init ns create once procfs mounted
	ProcMountInitIno = ProcDynamicFirst
	// LinuxKitNetNsInitIno not sure, just use in practice
	LinuxKitNetNsInitIno = ProcDynamicFirst
)

var (
	TypeMap = map[string]int{
		"cgroup":            TypeNamespaceTypeCGroup,
		"ipc":               TypeNamespaceTypeIPC,
		"mnt":               TypeNamespaceTypeMount,
		"net":               TypeNamespaceTypeNetwork,
		"pid":               TypeNamespaceTypePid,
		"pid_for_children":  TypeNamespaceTypePid,
		"user":              TypeNamespaceTypeUser,
		"uts":               TypeNamespaceTypeUTS,
		"time":              TypeNamespaceTypeTime,
		"time_for_children": TypeNamespaceTypeTime,
	}
	InitInoMap = map[int]int{
		TypeNamespaceTypeCGroup:  ProcCGroupInitIno,
		TypeNamespaceTypeIPC:     ProcIpcInitIno,
		TypeNamespaceTypeMount:   ProcMountInitIno,
		TypeNamespaceTypeNetwork: -1,
		TypeNamespaceTypePid:     ProcPidInitIno,
		TypeNamespaceTypeUser:    ProcUserInitIno,
		TypeNamespaceTypeUTS:     ProcUtsInitIno,
		TypeNamespaceTypeTime:    ProcTimeInitIno,
	}
)

func ParseNamespaces() (namespaces []Namespace, err error) {
	proc := "/proc/self/ns"
	namespaceInoMap, err := ListNamespaceDir(proc)
	if err != nil {
		return
	}
	for name, ino := range namespaceInoMap {
		namespace := Namespace{
			Name:            name,
			Path:            filepath.Join(proc, name),
			Type:            TypeMap[name],
			InodeNumber:     ino,
			InitInodeNumber: InitInoMap[TypeMap[name]],
		}
		namespaces = append(namespaces, namespace)
	}
	return
}
