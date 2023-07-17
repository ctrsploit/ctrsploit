package namespace

type Level int

const (
	TypeNamespaceLevelUnknown Level = iota
	TypeNamespaceLevelBoot
	TypeNamespaceLevelChild
	TypeNamespaceLevelNotSupported
	TypeNamespaceLevelHost = TypeNamespaceLevelBoot
)

var (
	LevelMap = map[Level]string{
		TypeNamespaceLevelChild:        "child",
		TypeNamespaceLevelHost:         "host",
		TypeNamespaceLevelNotSupported: "not supported( <=> host)",
		TypeNamespaceLevelUnknown:      "unknown",
	}
)

func (l Level) String() string {
	return LevelMap[l]
}
