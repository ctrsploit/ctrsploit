package noused

type Container interface {

	Rootfs() (string, error)
}

type Docker struct {
}
