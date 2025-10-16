package kernel

import (
	"fmt"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
)

// once ensures the kernel compile date is read and parsed only once.
var once sync.Once

func getCompiledDate(content string) (compiledDate time.Time, err error) {
	var re *regexp.Regexp
	var dateStr string

	// "Mon Jan 2 15:04:05 MST 2006"
	re = regexp.MustCompile(`\w{3}\s+\w{3}\s+\d{1,2}\s+\d{2}:\d{2}:\d{2}\s+\w+\s+\d{4}`)
	dateStr = re.FindString(content)
	if dateStr != "" {
		compiledDate, err = time.Parse("Mon Jan 2 15:04:05 MST 2006", dateStr)
		return
	}

	// "Tue, 13 Feb 2024 15:21:03 +0000"
	re = regexp.MustCompile(`\w{3},\s+\d{1,2}\s+\w{3}\s+\d{4}\s+\d{2}:\d{2}:\d{2}\s+[+-]\d{4}`)
	dateStr = re.FindString(content)
	if dateStr != "" {
		compiledDate, err = time.Parse(time.RFC1123Z, dateStr)
		return
	}

	// ISO-like
	re = regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
	dateStr = re.FindString(content)
	if dateStr != "" {
		compiledDate, err = time.Parse("2006-01-02", dateStr)
		return
	}

	err = fmt.Errorf("could not find a recognizable date string in /proc/version (content: %s)", content)
	return
}

// GetKernelCompiledDate reads and parses the kernel compilation date from /proc/version.
// It uses sync.Once to perform this operation only once.
func GetKernelCompiledDate() (t time.Time, err error) {
	once.Do(func() {
		// Read the content of /proc/version
		content, err := os.ReadFile("/proc/version")
		if err != nil {
			err = fmt.Errorf("failed to read /proc/version: %w", err)
			return
		}
		t, err = getCompiledDate(string(content))
	})
	return
}

type CompiledBefore struct {
	prerequisite.BasePrerequisite
	Expected time.Time
}

func (p *CompiledBefore) Check() (bool, error) {
	if p.Checked {
		return p.Satisfied, nil
	}
	actualCompiledDate, err := GetKernelCompiledDate()
	if err != nil {
		return false, fmt.Errorf("failed to determine compiled date: %w", err)
	}
	p.Satisfied = actualCompiledDate.Before(p.Expected)
	p.Checked = true
	return p.Satisfied, nil
}
