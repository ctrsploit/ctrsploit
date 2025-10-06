package runtime

import "strings"

type Type int

const (
	TypeDocker Type = 1 << iota
	TypeContainerd
	TypeCtr
	TypeNerdCtl
	TypePodman
	TypeCrio
	TypeK8s
)

var Type2String = map[Type]string{
	TypeDocker:     "docker",
	TypeContainerd: "containerd",
	TypeCtr:        "ctr",
	TypeNerdCtl:    "nerdctl",
	TypePodman:     "podman",
	TypeCrio:       "cri-o",
	TypeK8s:        "k8s",
}

var OrderedTypes = []Type{
	TypeDocker,
	TypeContainerd,
	TypeCtr,
	TypeNerdCtl,
	TypePodman,
	TypeCrio,
	TypeK8s,
}

func (t Type) String() string {
	var runtimes []string
	for _, typ := range OrderedTypes {
		if t&typ != 0 {
			runtimes = append(runtimes, Type2String[typ])
		}
	}
	if len(runtimes) == 0 {
		return "unknown"
	} else {
		return strings.Join(runtimes, "|")
	}
}

func GetType() (t Type) {
	if is, _ := K8s().Is(); is {
		t |= TypeK8s
	}
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
