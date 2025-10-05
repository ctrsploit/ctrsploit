package runtime

type Type int

const (
	TypeDocker Type = 1 << iota
	TypeContainerd
	TypeCtr
	TypeNerdCtl
	TypePodman
	TypeCrio
	TypeCri
	TypeK8s
)

func (t Type) String() string {
	switch t {
	case TypeDocker:
		return "docker"
	case TypeContainerd:
		return "containerd"
	case TypeCtr:
		return "ctr"
	case TypeNerdCtl:
		return "nerdctl"
	case TypePodman:
		return "podman"
	case TypeCrio:
		return "crio"
	case TypeCri:
		return "cri"
	case TypeK8s:
		return "k8s"
	default:
		return "Unknown"
	}
}

func GetType() (t Type) {
	if is, _ := Docker().Is(); is {
		t |= TypeDocker
	}
	if is, _ := Containerd().Is(); is {
		t |= TypeContainerd
	}
	if is, _ := Nerdctl().Is(); is {
		t |= TypeNerdCtl
	}
	return
}
