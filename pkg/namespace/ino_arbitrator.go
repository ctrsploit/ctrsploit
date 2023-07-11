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
	if ns.InodeNumber < ProcDynamicFirst || ns.InodeNumber > ProcTimeInitIno {
		return
	}
	if ns.InodeNumber < i.MaxIno {
		is = true
		return
	}
	return
}

func (i *InoArbitrator) Arbitrate(ns Namespace) (isHostNamespace bool, err error) {
	switch ns.Type {
	case TypeNamespaceTypeNetwork:
		err = prerequisite.KernelReleasedByLinuxkit.Check()
		if err != nil {
			return
		}
		if prerequisite.KernelReleasedByLinuxkit.Satisfied {
			isHostNamespace = LinuxKitNetNsInitIno == ns.InodeNumber
			return
		}
		isHostNamespace = i.IsNetworkNamespaceInoBetweenProcInoList(ns)
		if err != nil {
			return
		}
		if !isHostNamespace {
			return
		}
		isHostNamespaceGuess := i.IsNetworkNamespaceInoBetweenTwoAdjacentMissingIno(ns)
		if !isHostNamespaceGuess {
			log.Logger.Warningf("ino BetweenProcInoList but not BetweenTwoAdjacentMissingIno")
		}
		return
	case TypeNamespaceTypeCGroup:
		err = prerequisite.KernelSupportsCgroupNamespace.Check()
		if err != nil {
			return
		}
		if !prerequisite.KernelSupportsCgroupNamespace.Satisfied {
			isHostNamespace = true
			return
		}
		isHostNamespace = ns.InodeNumber == ns.InitInodeNumber
		return
	case TypeNamespaceTypeTime:
		err = prerequisite.KernelSupportsTimeNamespace.Check()
		if err != nil {
			return
		}
		if !prerequisite.KernelSupportsTimeNamespace.Satisfied {
			isHostNamespace = true
			return
		}
		isHostNamespace = ns.InodeNumber == ns.InitInodeNumber
		return
	case TypeNamespaceTypeIPC, TypeNamespaceTypeUTS, TypeNamespaceTypeUser, TypeNamespaceTypePid:
		isHostNamespace = ns.InodeNumber == ns.InitInodeNumber
		return
	default:
		// TODO: raise unknown type error
	}

	return
}

func (i *InoArbitrator) PrerequisitesSatisfied() (satisfied bool, err error) {
	// TODO: need check kernel >= v3.8?
	satisfied = true
	return
}
