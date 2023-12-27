package version

type TypeState int

const (
	StateUnknown TypeState = iota
	StateValid
	StateDisable
	StateUnsupported
)

type TypeSoftware int

const (
	SoftwareUnknown TypeSoftware = iota
	SoftwareDocker
	SoftwareRunc
	SoftwareKernel
)
