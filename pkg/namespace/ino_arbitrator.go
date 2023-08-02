package namespace

import (
	"github.com/ctrsploit/ctrsploit/log"
	"github.com/ctrsploit/ctrsploit/prerequisite"
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
	return
}

func (i *InoArbitrator) GuessNetworkNamespaceInitialIno() (initialIno int) {
	for index, ino := range i.InoList[1:] {
		if ino-i.InoList[index] >= 3 {
			initialIno = i.InoList[index] + 2
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
	err = prerequisite.KernelReleasedByLinuxkit.Check()
	if err != nil {
		return
	}
	if prerequisite.KernelReleasedByLinuxkit.Satisfied {
		switch ns.Type {
		case TypeNamespaceTypeNetwork:
			isHostNamespace = LinuxKitNetNsInitIno == ns.InodeNumber
			break
		case TypeNamespaceTypeMount:
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
		case TypeNamespaceTypeNetwork:
			// network namespace:
			//	linuxkit: 0xF0000000
			//	normal:
			//		verify: first two hole => host
			// verify
			isHostNamespace = i.IsNetworkNamespaceInoBetweenTwoAdjacentMissingIno(ns)
		case TypeNamespaceTypeCGroup:
			// cgroups namespace:
			//     kernel not supports => host
			//     otherwise, compare with initIno
			err = prerequisite.KernelSupportsCgroupNamespace.Check()
			if err != nil {
				return
			}
			log.Logger.Debugf("kernel supports cgroup ns: %t\n", prerequisite.KernelSupportsCgroupNamespace.Satisfied)
			if prerequisite.KernelSupportsCgroupNamespace.Satisfied {
				isHostNamespace = ns.InodeNumber == ns.InitInodeNumber
			} else {
				// kernel not supports <=> host ns
				isHostNamespace = true
			}
		case TypeNamespaceTypeTime:
			// time namespace:
			//     kernel not supports => host
			//     otherwise, compare with initIno
			err = prerequisite.KernelSupportsTimeNamespace.Check()
			if err != nil {
				return
			}
			if prerequisite.KernelSupportsTimeNamespace.Satisfied {
				isHostNamespace = ns.InodeNumber == ns.InitInodeNumber
			} else {
				// kernel not supports <=> host ns
				isHostNamespace = true
			}
		case TypeNamespaceTypeIPC, TypeNamespaceTypeUTS, TypeNamespaceTypeUser, TypeNamespaceTypePid, TypeNamespaceTypeMount:
			isHostNamespace = ns.InodeNumber == ns.InitInodeNumber
		default:
			log.Logger.Warningf("%s has an unknown namespace level: %+v\n", ns.Name, ns)
		}
	}
	if isHostNamespace {
		namespaceLevel = TypeNamespaceLevelHost
	} else {
		namespaceLevel = TypeNamespaceLevelChild
	}
	return
}

func (i *InoArbitrator) PrerequisitesSatisfied() (satisfied bool, err error) {
	// TODO: need check kernel >= v3.8?
	satisfied = true
	return
}
