package pipe_primitive

import (
	"fmt"
)

type Primitive interface {
	GetExpName() string
	MinOffset() int64
	Write(filepath string, offset int64, content []byte) (err error)
}

func escapeExpName(primitive Primitive) string {
	return fmt.Sprintf("%s-escape", primitive.GetExpName())
}

func escalateExpName(primitive Primitive) string {
	return fmt.Sprintf("%s-permission-escalate", primitive.GetExpName())
}

func imagePollutionExpName(primitive Primitive) string {
	return fmt.Sprintf("%s-image-pollution", primitive.GetExpName())
}
