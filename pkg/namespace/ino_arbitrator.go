package namespace

import (
	"github.com/ctrsploit/ctrsploit/log"
	"github.com/ctrsploit/ctrsploit/prerequisite"
	"github.com/ctrsploit/ctrsploit/prerequisite/kernel"
	"github.com/pkg/errors"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"sort"
)

type InoArbitrator struct {
	InoList       []int
	MinIno        int
	MaxIno        int
	Prerequisites prerequisite.Prerequisites
}

func NewInoArbitrator() (arbitrator *InoArbitrator, err error) {
	arbitrator = &InoArbitrator{
		Prerequisites: prerequisite.Prerequisites{},
	}
	err = arbitrator.init()
	if err != nil {
		return
	}
	return
}

func (i *InoArbitrator) init() (err error) {
	proc := "/proc"
	inoList, err := ReadInodeNumberListUnderProc(proc)
	if err != nil {
		return
	}
	sort.Ints(inoList)
	for len(inoList) > 0 && inoList[0] < ProcDynamicFirst {
		inoList = inoList[1:]
	}
	i.InoList = inoList
	i.MinIno = inoList[0]
	i.MaxIno = inoList[len(inoList)-1]
	log.Logger.Debugf("min ino: %d, max ino: %d", i.MinIno, i.MaxIno)
	//log.Logger.Debugf("inoList: %+v", i.InoList)
	return
}

func (i *InoArbitrator) GuessNetworkNamespaceInitialIno() (initialIno int) {
	for index, ino := range i.InoList[1:] {
		if ino-i.InoList[index] >= 3 {
			initialIno = i.InoList[index] + 2
			log.Logger.Debugf("InoList[%d]=%d, InoList[%d]=%d", index+1, ino, index, i.InoList[index])
			log.Logger.Debugf("network ns initial ino maybe: %d", initialIno)
			return
		}
	}
	awesome_error.CheckWarning(errors.Errorf("cannot guess network ns initial inode number"))
	return
}

func (i *InoArbitrator) IsNetworkNamespaceInoBetweenTwoAdjacentMissingIno(ns Namespace) (is bool) {
	if len(i.InoList) == 0 {
		return
	}
	initialIno := i.GuessNetworkNamespaceInitialIno()
	if initialIno == 0 {
		return
	}
	if ns.InodeNumber < i.MinIno {
		is = true
		return
	}
	is = initialIno == ns.InodeNumber
	return
}

func (i *InoArbitrator) IsNetworkNamespaceInoBetweenProcInoList(ns Namespace) (is bool) {
	if ns.InodeNumber < ProcDynamicFirst {
		return
	}
	if ns.InodeNumber < i.MaxIno {
		is = true
		return
	}
	return
}

func (i *InoArbitrator) Arbitrate(ns Namespace) (namespaceLevel Level, err error) {
	var isHostNamespace, normal bool
	// linuxkit
	err = kernel.ReleasedByLinuxkit.Check()
	if err != nil {
		return
	}
	if kernel.ReleasedByLinuxkit.Satisfied {
		switch ns.Type {
		case TypeNetwork:
			isHostNamespace = LinuxKitNetNsInitIno == ns.InodeNumber
			break
		case TypeMount:
			isHostNamespace = LinuxKitMountNsInitIno == ns.InodeNumber
			break
		default:
			normal = true
		}
	} else {
		normal = true
	}
	if normal {
		switch ns.Type {
		case TypeNetwork:
			// network namespace:
			//	linuxkit: 0xF0000000
			//	normal:
			//		verify: first two hole => host
			// verify
			isHostNamespace = i.IsNetworkNamespaceInoBetweenTwoAdjacentMissingIno(ns)
		case TypeCGroup:
			// cgroups namespace:
			//     kernel not supports => host
			//     otherwise, compare with initIno
			err = kernel.SupportsCgroupNamespace.Check()
			if err != nil {
				return
			}
			log.Logger.Debugf("kernel supports cgroup ns: %t\n", kernel.SupportsCgroupNamespace.Satisfied)
			if kernel.SupportsCgroupNamespace.Satisfied {
				isHostNamespace = ns.InodeNumber == ns.InitInodeNumber
			} else {
				// kernel not supports <=> host ns
				isHostNamespace = true
			}
		case TypeTime:
			// time namespace:
			//     kernel not supports => host
			//     otherwise, compare with initIno
			err = kernel.SupportsTimeNamespace.Check()
			if err != nil {
				return
			}
			if kernel.SupportsTimeNamespace.Satisfied {
				isHostNamespace = ns.InodeNumber == ns.InitInodeNumber
			} else {
				// kernel not supports <=> host ns
				isHostNamespace = true
			}
		case TypeIPC, TypeUTS, TypeUser, TypePid, TypeMount:
			isHostNamespace = ns.InodeNumber == ns.InitInodeNumber
		default:
			log.Logger.Warningf("%s has an unknown namespace level: %+v\n", ns.Name, ns)
		}
	}
	if isHostNamespace {
		namespaceLevel = LevelHost
	} else {
		namespaceLevel = LevelChild
	}
	return
}

func (i *InoArbitrator) PrerequisitesSatisfied() (satisfied bool, err error) {
	// TODO: need check kernel >= v3.8?
	satisfied = true
	return
}
