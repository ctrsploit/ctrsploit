package kernel

import (
	"fmt"
	"time"

	"github.com/ctrsploit/ctrsploit/pkg/kernel"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
)

type CompiledBefore struct {
	prerequisite.BasePrerequisite
	Expected time.Time
}

func (p *CompiledBefore) Check() (bool, error) {
	if p.Checked {
		return p.Satisfied, nil
	}
	actualCompiledDate, err := kernel.GetCompiledDate()
	if err != nil {
		return false, fmt.Errorf("failed to determine compiled date: %w", err)
	}
	p.Satisfied = actualCompiledDate.Before(p.Expected)
	p.Checked = true
	return p.Satisfied, nil
}
