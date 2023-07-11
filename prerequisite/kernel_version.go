package prerequisite

type KernelVersion struct {
	ExpectedMinVersion string
	ExpectedMaxVersion string
	BasePrerequisite
}

var (
	KernelSupportsCgroupNamespace = KernelVersion{
		ExpectedMinVersion: "4.6",
		ExpectedMaxVersion: "",
		BasePrerequisite: BasePrerequisite{
			Name: "Kernel Supports Cgroup namespace",
			Info: "kernel version >= v4.6",
		},
	}
	KernelSupportsTimeNamespace = KernelVersion{
		ExpectedMinVersion: "5.6",
		ExpectedMaxVersion: "",
		BasePrerequisite: BasePrerequisite{
			Name: "Kernel Supports Time namespace",
			Info: "kernel version >= v5.6",
		},
	}
)

func (p *KernelVersion) Check() (err error) {
	// TODO: get kernel version
	version := "xxx"
	p.Satisfied = p.ExpectedMinVersion <= version && p.ExpectedMaxVersion >= version
	return
}
