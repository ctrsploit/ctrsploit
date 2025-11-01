package rlimit

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/ctrsploit/sploit-spec/pkg/env/container/kernel/rlimit"
	"golang.org/x/sys/unix"
)

type Resource unix.Rlimit

func (r Resource) String() string {
	format := func(v uint64) string {
		if v == unix.RLIM_INFINITY {
			return "unlimited"
		}
		return fmt.Sprintf("%d", v)
	}

	return fmt.Sprintf("soft:%s, hard:%s", format(r.Cur), format(r.Max))
}

func GetAll() (rlimit.Rlimit, error) {
	r := rlimit.Rlimit{}
	v := reflect.ValueOf(&r).Elem()

	for name, typ := range rlimit.MapNameToType {
		var resource unix.Rlimit
		if err := unix.Getrlimit(typ, &resource); err != nil {
			return r, fmt.Errorf("failed to get rlimit for %s: %w", name, err)
		}
		field := v.FieldByNameFunc(func(f string) bool {
			return strings.EqualFold(f, name)
		})
		if field.IsValid() && field.CanSet() {
			field.Set(reflect.ValueOf(resource))
		}
	}

	return r, nil
}
