package rlimit

import (
	"github.com/ctrsploit/ctrsploit/pkg/rlimit"
	spec "github.com/ctrsploit/sploit-spec/pkg/env/container/kernel/rlimit"
)

func Rlimit() (spec.Rlimit, error) {
	return rlimit.GetAll()
}
