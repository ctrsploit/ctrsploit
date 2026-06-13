package module

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
)

// Loaded checks whether a kernel module is loaded by reading /proc/modules.
type Loaded struct {
	prerequisite.BasePrerequisite
	// Name of the module to check, e.g. "nf_tables".
	Module string
}

// New checks for a loaded kernel module.
func New(module string) *Loaded {
	return &Loaded{
		BasePrerequisite: prerequisite.BasePrerequisite{
			Name: fmt.Sprintf("%s loaded", module),
			Info: fmt.Sprintf("%s kernel module is loaded.", module),
		},
		Module: module,
	}
}

func (c *Loaded) Check() (bool, error) {
	return c.CheckTemplate(func() {
		file, err := os.Open("/proc/modules")
		if err != nil {
			c.Err = c.WrapErr(fmt.Errorf("reading /proc/modules: %w", err))
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			fields := strings.Fields(scanner.Text())
			if len(fields) > 0 && fields[0] == c.Module {
				c.Satisfied = true
				return
			}
		}
		if err := scanner.Err(); err != nil {
			c.Err = c.WrapErr(fmt.Errorf("scanning /proc/modules: %w", err))
		}
	})
}
