package module

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ctrsploit/ctrsploit/pkg/kernel/uname"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
)

// Loaded checks whether a kernel module is loaded by reading /proc/modules.
type Loaded struct {
	prerequisite.BasePrerequisite
	// Name of the module to check, e.g. "nf_tables".
	Module string
}

type Available struct {
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

// NewAvailable checks whether a module is loaded, built in, or available for
// kernel autoloading. It does not load the module.
func NewAvailable(module string) *Available {
	return &Available{
		BasePrerequisite: prerequisite.BasePrerequisite{
			Name: fmt.Sprintf("%s available", module),
			Info: fmt.Sprintf("%s kernel module is loaded, built in, or available for autoloading.", module),
		},
		Module: module,
	}
}

func (c *Loaded) Check() (bool, error) {
	return c.CheckTemplate(func() {
		loaded, err := loaded(c.Module)
		if err != nil {
			c.Err = c.WrapErr(err)
			return
		}
		c.Satisfied = loaded
	})
}

func (c *Available) Check() (bool, error) {
	return c.CheckTemplate(func() {
		available, err := available(c.Module)
		if err != nil {
			c.Err = c.WrapErr(err)
			return
		}
		c.Satisfied = available
	})
}

func available(module string) (bool, error) {
	loaded, err := loaded(module)
	if err != nil {
		return false, err
	}
	if loaded {
		return true, nil
	}

	if exists(filepath.Join("/sys/module", module)) {
		return true, nil
	}

	release, err := uname.Release()
	if err != nil {
		return false, fmt.Errorf("getting kernel release: %w", err)
	}
	modulesDir := filepath.Join("/lib/modules", release)
	for _, name := range moduleIndexFiles(modulesDir) {
		found, err := fileContainsModule(name, module)
		if err != nil {
			return false, err
		}
		if found {
			return true, nil
		}
	}
	return false, nil
}

func loaded(module string) (bool, error) {
	file, err := os.Open("/proc/modules")
	if err != nil {
		return false, fmt.Errorf("reading /proc/modules: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) > 0 && fields[0] == module {
			return true, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return false, fmt.Errorf("scanning /proc/modules: %w", err)
	}
	return false, nil
}

func moduleIndexFiles(modulesDir string) []string {
	return []string{
		filepath.Join(modulesDir, "modules.alias"),
		filepath.Join(modulesDir, "modules.dep"),
		filepath.Join(modulesDir, "modules.builtin"),
		filepath.Join(modulesDir, "modules.builtin.modinfo"),
	}
}

func fileContainsModule(path string, module string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("reading %s: %w", path, err)
	}
	defer file.Close()

	normalizedModule := strings.ReplaceAll(module, "-", "_")
	pathModule := strings.ReplaceAll(module, "_", "-")
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, module+".ko") ||
			strings.Contains(line, pathModule+".ko") ||
			strings.Contains(line, "name="+normalizedModule) {
			return true, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return false, fmt.Errorf("scanning %s: %w", path, err)
	}
	return false, nil
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
