package env

type Container interface {

	Rootfs() (string, error)
}

type Docker struct {
}
