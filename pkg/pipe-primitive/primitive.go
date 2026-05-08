package pipeprimitive

// Primitive writes content into a file without opening it writable.
type Primitive interface {
	GetExpName() string
	MinOffset() int64
	Write(path string, offset int64, content []byte) error
}

type EscapeImageWriterProvider interface {
	EscapeImageWriter() []byte
}

type EscapeImageExtraFileProvider interface {
	EscapeImageExtraFiles() map[string][]byte
}

func escapeName(primitive Primitive) string {
	return "escape"
}

func escalateName(primitive Primitive) string {
	return "privilege-escalate"
}

func imagePollutionName(primitive Primitive) string {
	return "image-pollution"
}
