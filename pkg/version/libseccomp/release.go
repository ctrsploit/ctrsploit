package libseccomp

import "time"

var (
	LastTimeUpdate, _ = time.Parse(time.RFC3339, "2023-12-07T14:19:00Z08:00")
	Versions          = NewMap([]string{
		"2.3.3",
		"2.4.1",
		"2.5.1",
	})
)
