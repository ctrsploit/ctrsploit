package runtime

type Type int

const (
	TypeUnknown Type = iota
	TypeDocker
	TypeContainerd
	TypePodman
	TypeCrio
)

func (t Type) String() string {
	switch t {
	case TypeDocker:
		return "Overlay"
	case TypeContainerd:
		return "DeviceMapper"
	case TypePodman:
		return "Podman"
	default:
		return "Unknown"
	}
}
