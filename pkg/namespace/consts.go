package namespace

import "github.com/ctrsploit/ctrsploit/util"

// Name

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
	ProcMountInitIno = ProcDynamicFirst
	// LinuxKitNetNsInitIno not sure, just use in practice
	LinuxKitNetNsInitIno   = ProcDynamicFirst
	LinuxKitMountNsInitIno = ProcDynamicFirst + 1
)

// Level

type Level int

const (
	LevelUnknown Level = iota
	LevelBoot
	LevelChild
	LevelNotSupported
	LevelHost = LevelBoot
)

var (
	LevelMap = map[Level]string{
		LevelChild:        "child",
		LevelHost:         "host",
		LevelNotSupported: "not supported( <=> host)",
		LevelUnknown:      "unknown",
	}
)

func (l Level) String() string {
	return LevelMap[l]
}

// Type

type Type int

const (
	TypeUnknown Type = iota
	TypeIPC
	TypeUTS
	TypeUser
	TypePid
	TypeCGroup
	TypeTime
	TypeMount
	TypeNetwork
)

// Map

var (
	MapName2Type = map[string]Type{
		NameCGroup: TypeCGroup,
		NameIpc:    TypeIPC,
		NameMnt:    TypeMount,
		NameNet:    TypeNetwork,
		NamePid:    TypePid,
		// TODO: not sure pid_for_children is same as pid?
		NamePidForChildren: TypePid,
		NameUser:           TypeUser,
		NameUts:            TypeUTS,
		NameTime:           TypeTime,
		// TODO: not sure time_for_children is same as time?
		NameTimeForChildren: TypeTime,
	}
	MapType2Name = util.ReverseMap(MapName2Type).(map[Type]string)
	InitInoMap   = map[Type]int{
		TypeCGroup:  ProcCGroupInitIno,
		TypeIPC:     ProcIpcInitIno,
		TypeMount:   ProcMountInitIno,
		TypeNetwork: -1,
		TypePid:     ProcPidInitIno,
		TypeUser:    ProcUserInitIno,
		TypeUTS:     ProcUtsInitIno,
		TypeTime:    ProcTimeInitIno,
	}
)
