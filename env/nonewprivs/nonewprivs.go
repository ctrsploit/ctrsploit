package nonewprivs

import (
	"github.com/ctrsploit/ctrsploit/pkg/proc/status"
)

const (
	CommandName        = "no-new-privs"
	procSelfStatusPath = "/proc/self/status"
)

type Info struct {
	StatusPath string `json:"status_path"`
	Enabled    bool   `json:"enabled"`
}

func Current() (Info, error) {
	return FromStatusFile(procSelfStatusPath)
}

func FromStatusFile(path string) (Info, error) {
	procStatus, err := status.ParseStatusFile(path)
	if err != nil {
		return Info{}, err
	}
	return Info{
		StatusPath: path,
		Enabled:    procStatus.NoNewPrivs,
	}, nil
}
