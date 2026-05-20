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

type RestartLoaderProvider interface {
	RestartLoader(payload []byte) ([]byte, error)
}

type RuncOverwritePayloadProvider interface {
	RuncOverwritePayload(cmd string) ([]byte, error)
}
