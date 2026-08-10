package cpusec

import (
	"fmt"

	"github.com/ctrsploit/sploit-spec/pkg/printer"
	"github.com/ctrsploit/sploit-spec/pkg/result"
	"github.com/ctrsploit/sploit-spec/pkg/result/item"
)

// Result is the human-readable rendering of Machine. Only fields relevant to
// machine.Arch are populated; off-arch item.Bool fields are left zero-valued
// (empty Name, false Result) so the reflective printer skips them via
// item.Bool.IsEmpty().
type Result struct {
	Name         result.Title `json:"name"`
	Arch         item.Short   `json:"arch"`
	ConfigSource item.Short   `json:"config_source"`

	// x86-64
	SMEP    item.Bool `json:"smep"`
	SMAP    item.Bool `json:"smap"`
	KPTI    item.Bool `json:"kpti"`
	IBT     item.Bool `json:"ibt"`
	KCFI    item.Bool `json:"kcfi"`
	FgKaslr item.Bool `json:"fg_kaslr"`

	// arm64
	PAC item.Bool `json:"pac"`
	BTI item.Bool `json:"bti"`
	PAN item.Bool `json:"pan"`
	MTE item.Bool `json:"mte"`
}

// Human converts the machine-readable result into the human-readable form.
// Off-arch fields are left empty so the printer omits them.
func Human(machine Machine) (human Result) {
	human = Result{
		Name: result.Title{
			Name: "CPU/Kernel Security Mitigations",
		},
		Arch: item.Short{
			Name:   "Arch",
			Result: machine.Arch,
		},
		ConfigSource: item.Short{
			Name:   "Config Source",
			Result: machine.ConfigSource,
		},
	}

	switch machine.Arch {
	case "amd64":
		human.SMEP = item.Bool{Name: "SMEP", Description: "blocks kernel executing user pages (CET)", Result: machine.SMEP}
		human.SMAP = item.Bool{Name: "SMAP", Description: "blocks kernel read/write user pages", Result: machine.SMAP}
		human.KPTI = item.Bool{Name: "KPTI", Description: "user pagetables map no kernel text/data", Result: machine.KPTI}
		human.IBT = item.Bool{Name: "IBT", Description: "indirect jumps must land on ENDBR (CET-IBT)", Result: machine.IBT}
		human.KCFI = item.Bool{Name: "KCFI", Description: "Clang indirect-call target type check", Result: machine.KCFI}
		human.FgKaslr = item.Bool{Name: "FG-KASLR", Description: "per-function randomization beyond base KASLR", Result: machine.FgKaslr}
	case "arm64":
		human.PAC = item.Bool{Name: "PAC", Description: "pointer authentication signs code/data pointers", Result: machine.PAC}
		human.BTI = item.Bool{Name: "BTI", Description: "indirect branches must land on BTI instructions", Result: machine.BTI}
		human.KPTI = item.Bool{Name: "KPTI", Description: "user pagetables map no kernel text/data", Result: machine.KPTI}
		human.PAN = item.Bool{Name: "PAN", Description: "blocks kernel accessing user pages", Result: machine.PAN}
		human.MTE = item.Bool{Name: "MTE", Description: "tag-based memory access checks", Result: machine.MTE}
	}
	return
}

// Print gathers the mitigations and emits them via the global printer (text,
// colorful, or json per the global --json/--colorful flags).
func Print() (err error) {
	machine, err := Cpusec()
	if err != nil {
		return
	}
	u := result.Union{
		Machine: machine,
		Human:   Human(machine),
	}
	fmt.Println(printer.Printer.Print(u))
	return
}
