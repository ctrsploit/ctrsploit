package runtime

type Runtime interface {
	Is() (bool, error)
}
