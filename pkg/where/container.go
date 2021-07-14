package where

type Container interface {
	IsIn() (in bool, err error)
}
