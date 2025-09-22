package namespace

import (
	"github.com/ctrsploit/sploit-spec/pkg/env/container"
)

// Init Inode Number

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
	ProcMountInitIno = ProcDynamicFirst + 1
	// LinuxKitNetNsInitIno not sure, just use in practice
	LinuxKitNetNsInitIno   = ProcDynamicFirst
	LinuxKitMountNsInitIno = ProcDynamicFirst + 1
)

// Map

var (
	InitInoMap = map[container.NamespaceType]int{
		container.NamespaceTypeCGroup:  ProcCGroupInitIno,
		container.NamespaceTypeIPC:     ProcIpcInitIno,
		container.NamespaceTypeMount:   ProcMountInitIno,
		container.NamespaceTypeNetwork: -1,
		container.NamespaceTypePid:     ProcPidInitIno,
		container.NamespaceTypeUser:    ProcUserInitIno,
		container.NamespaceTypeUTS:     ProcUtsInitIno,
		container.NamespaceTypeTime:    ProcTimeInitIno,
	}
)
