package version

type TypeState int

const (
	StateUnknown TypeState = iota
	StateValid
	StateDisable
	StateUnsupported
)

func (t TypeState) String() (s string) {
	switch t {
	case StateValid:
		s = "valid"
	case StateDisable:
		s = "disable"
	case StateUnsupported:
		s = "unsupported"
	default:
		s = "unknown"
	}
	return
}

type TypeSoftware int

const (
	SoftwareUnknown TypeSoftware = iota
	SoftwareDocker
	SoftwareRunc
	SoftwareKernel
)
