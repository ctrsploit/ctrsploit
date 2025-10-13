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

// compiledDate stores the parsed kernel compilation time.
var compiledDate time.Time

// compileDateErr stores any error that occurred during reading or parsing.
var compileDateErr error

// GetKernelCompiledDate reads and parses the kernel compilation date from /proc/version.
// It uses sync.Once to perform this operation only once.
func GetKernelCompiledDate() (time.Time, error) {
	once.Do(func() {
		// Read the content of /proc/version
		content, err := os.ReadFile("/proc/version")
		if err != nil {
			compileDateErr = fmt.Errorf("failed to read /proc/version: %w", err)
			return
		}

		// The date in /proc/version typically looks like:
		// "... #1 SMP PREEMPT_DYNAMIC Tue, 13 Feb 2024 15:21:03 +0000 ..."
		// or "... #83-Ubuntu SMP Mon Jun 5 14:18:34 UTC 2023"
		// We use a regex to find a common pattern for the date and time.
		// This regex is designed to capture various common formats.
		// Example: "Mon Jan _2 15:04:05 MST 2006"
		re := regexp.MustCompile(`\w{3}\s+\w{3}\s+\d{1,2}\s+\d{2}:\d{2}:\d{2}\s+\w+\s+\d{4}`)
		dateStr := re.FindString(string(content))
		if dateStr == "" {
			// Try another common format, e.g., RFC1123Z used in some distros
			// Example: "Tue, 13 Feb 2024 15:21:03 +0000"
			re = regexp.MustCompile(`\w{3},\s+\d{2}\s+\w{3}\s+\d{4}\s+\d{2}:\d{2}:\d{2}\s+[+-]\d{4}`)
			dateStr = re.FindString(string(content))
			if dateStr == "" {
				compileDateErr = fmt.Errorf("could not find a recognizable date string in /proc/version")
				return
			}
			// Parse using RFC1123Z layout
			compiledDate, compileDateErr = time.Parse(time.RFC1123Z, dateStr)
			return
		}

		// Define the layout for the first matched format.
		// Go's reference time is: Mon Jan 2 15:04:05 MST 2006
		const layout = "Mon Jan _2 15:04:05 MST 2006"
		compiledDate, compileDateErr = time.Parse(layout, dateStr)
	})
	return compiledDate, compileDateErr
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
