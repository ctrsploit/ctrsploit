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
	LinuxKitNetNsInitIno   = ProcDynamicFirst
	LinuxKitMountNsInitIno = ProcDynamicFirst + 1
)

const (
	NameCGroup          = "cgroup"
	NameIpc             = "ipc"
	NameMnt             = "mnt"
	NameNet             = "net"
	NamePid             = "pid"
	NamePidForChildren  = "pid_for_children"
	NameUser            = "user"
	NameUts             = "uts"
	NameTime            = "time"
	NameTimeForChildren = "time_for_children"
)

var (
	TypeMap = map[string]int{
		NameCGroup: TypeNamespaceTypeCGroup,
		NameIpc:    TypeNamespaceTypeIPC,
		NameMnt:    TypeNamespaceTypeMount,
		NameNet:    TypeNamespaceTypeNetwork,
		NamePid:    TypeNamespaceTypePid,
		// TODO: not sure pid_for_children is same as pid?
		NamePidForChildren: TypeNamespaceTypePid,
		NameUser:           TypeNamespaceTypeUser,
		NameUts:            TypeNamespaceTypeUTS,
		NameTime:           TypeNamespaceTypeTime,
		// TODO: not sure time_for_children is same as time?
		NameTimeForChildren: TypeNamespaceTypeTime,
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
			Type:            TypeMap[name],
			InodeNumber:     ino,
			InitInodeNumber: InitInoMap[TypeMap[name]],
		}
		namespaces = append(namespaces, namespace)
	}
	return
}
