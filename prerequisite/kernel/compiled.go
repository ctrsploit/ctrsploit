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
	return p.CheckTemplate(func() (bool, error) {
		actualCompiledDate, err := kernel.GetCompiledDate()
		if err != nil {
			p.Err = fmt.Errorf("failed to check [%s], caused by unable to determine compiled date: %w", p.GetName(), err)
			return p.Satisfied, p.Err
		}
		p.Satisfied = actualCompiledDate.Before(p.Expected)
		return p.Satisfied, p.Err
	})
}
