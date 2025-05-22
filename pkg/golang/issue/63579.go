package issue

import "runtime"

/*
I63579 showed a bug in Go's net package,
In the syscall file there is code that special-cases handling of unix socket names beginning with @,
by removing the trailing NUL character,
but it does not do the same for unix socket names starting with the NUL character.
https://github.com/golang/go/issues/63579

This function returns whether the issue exists in the current golang sdk,
returns true if the issue exists,
*/
func I63579() bool {
	// It fixed in go1.22.0, [3de6033](https://github.com/golang/go/commit/3de6033d0e8022dffee85bd9537f90b1a5ba5e30)
	return runtime.Version() < "go1.22.0"
}
