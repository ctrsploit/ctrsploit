package storagedriver

import (
	"errors"
	"fmt"
	"os"

	"github.com/ctrsploit/ctrsploit/pkg/module"
	"github.com/ctrsploit/ctrsploit/pkg/mountinfo"
	"github.com/ctrsploit/sploit-spec/pkg/env/container/storagedriver"
)

type Overlay struct {
}

func NewOverlay() *Overlay {
	return &Overlay{}
}

func (o *Overlay) Type() storagedriver.Type {
	return storagedriver.TypeOverlay
}

func (o *Overlay) Enabled() (bool, error) {
	return module.Loaded("overlay")
}

func (o *Overlay) Number() (int, error) {
	number, err := module.RefCount("overlay")
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		} else {
			return 0, fmt.Errorf("error getting overlay number: %w", err)
		}
	}
	return number, nil
}

func (o *Overlay) Used() (bool, error) {
	var errs []error
	info, err := mountinfo.RootMount()
	if err != nil {
		errs = append(errs, fmt.Errorf("error getting root mount info: %w", err))
	} else if mountinfo.IsOverlay(info) {
		return true, nil
	}
	number, err := o.Number()
	if err != nil {
		errs = append(errs, fmt.Errorf("error getting overlay number: %w", err))
	} else if number > 0 {
		return true, nil
	}
	return false, errors.Join(errs...)
}
