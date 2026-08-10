package cpusec

import (
	"github.com/ctrsploit/ctrsploit/pkg/cpusec"
)

const (
	CommandName = "cpusec"
)

// Machine is the machine-readable result of the `env cpusec` subcommand. It is
// a json-tagged mirror of cpusec.Mitigations; keeping a local struct (rather
// than extending sploit-spec) avoids a cross-repo schema release. The layout
// matches what a future sploit-spec/pkg/env/container/cpusec.go would hold,
// so promoting it is a 1-file move.
type Machine struct {
	Arch         string   `json:"arch"`
	ConfigSource string   `json:"config_source"`
	CPUFlags     []string `json:"cpu_flags"`

	// x86-64
	SMEP    bool `json:"smep"`
	SMAP    bool `json:"smap"`
	KPTI    bool `json:"kpti"`
	IBT     bool `json:"ibt"`
	KCFI    bool `json:"kcfi"`
	FgKaslr bool `json:"fg_kaslr"`

	// arm64
	PAC bool `json:"pac"`
	BTI bool `json:"bti"`
	PAN bool `json:"pan"`
	MTE bool `json:"mte"`
}

// Cpusec gathers the active CPU/kernel mitigations and returns the
// machine-readable result.
func Cpusec() (machine Machine, err error) {
	m, err := cpusec.Detect()
	if err != nil {
		return
	}
	machine = Machine{
		Arch:         m.Arch,
		ConfigSource: m.ConfigSource,
		CPUFlags:     m.CPUFlags,
		SMEP:         m.SMEP,
		SMAP:         m.SMAP,
		KPTI:         m.KPTI,
		IBT:          m.IBT,
		KCFI:         m.KCFI,
		FgKaslr:      m.FgKaslr,
		PAC:          m.PAC,
		BTI:          m.BTI,
		PAN:          m.PAN,
		MTE:          m.MTE,
	}
	return
}
