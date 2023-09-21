package where

type Where interface {
	IsIn() (in bool, err error)
}
