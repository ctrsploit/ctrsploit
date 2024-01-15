package libseccomp

import (
	"time"
)

var (
	LastTimeUpdate, _ = time.Parse(time.RFC3339, "2023-12-07T14:19:00Z08:00")
	Versions          = NewMap([]string{
		"v2.3.3",
		"v2.4.1",
		"v2.4.4",
		"v2.5.1",
		"v2.5.4",
	})
)
