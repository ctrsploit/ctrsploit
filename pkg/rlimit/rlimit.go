package rlimit

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/ctrsploit/sploit-spec/pkg/env/container/kernel/rlimit"
	"golang.org/x/sys/unix"
)

type Resource unix.Rlimit

func New(name string) (Resource, error) {
	t, ok := rlimit.MapNameToType[name]
	if !ok {
		return Resource{}, fmt.Errorf("rlimit %s not found", name)
	}
	r := unix.Rlimit{}
	if err := unix.Getrlimit(t, &r); err != nil {
		return Resource{}, fmt.Errorf("failed to get rlimit for %s: %w", name, err)
	}
	return Resource(r), nil
}

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

	for name, _ := range rlimit.MapNameToType {
		resource, err := New(name)
		if err != nil {
			return r, err
		}
		field := v.FieldByNameFunc(func(f string) bool {
			return strings.EqualFold(f, name)
		})
		if field.IsValid() && field.CanSet() {
			field.Set(reflect.ValueOf(unix.Rlimit(resource)))
		}
	}

	return r, nil
}
