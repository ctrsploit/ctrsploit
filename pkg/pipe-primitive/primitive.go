package pipeprimitive

import "fmt"

// Primitive writes content into a file without opening it writable.
type Primitive interface {
	GetExpName() string
	MinOffset() int64
	Write(path string, offset int64, content []byte) error
}

func escapeName(primitive Primitive) string {
	return fmt.Sprintf("%s-escape", primitive.GetExpName())
}

func escalateName(primitive Primitive) string {
	return fmt.Sprintf("%s-permission-escalate", primitive.GetExpName())
}

func imagePollutionName(primitive Primitive) string {
	return fmt.Sprintf("%s-image-pollution", primitive.GetExpName())
}
