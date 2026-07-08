package kernelprivesc

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

// recordingPrimitive records every WriteKmem/MemsetKmem call and resolves
// symbols from a fixed table. Mirrors pipeprimitive's recordingPrimitive.
type recordingPrimitive struct {
	expName string
	syms    map[string]uint64
	writes  []kmemWrite
	memsets []kmemMemset
}

type kmemWrite struct {
	kaddr   uint64
	content []byte
}
type kmemMemset struct {
	kaddr uint64
	fill  byte
	n    uint64
}

func (r *recordingPrimitive) ExpName() string { return r.expName }
func (r *recordingPrimitive) WriteKmem(kaddr uint64, content []byte) error {
	c := make([]byte, len(content))
	copy(c, content)
	r.writes = append(r.writes, kmemWrite{kaddr, c})
	return nil
}
func (r *recordingPrimitive) MemsetKmem(kaddr uint64, fill byte, n uint64) error {
	r.memsets = append(r.memsets, kmemMemset{kaddr, fill, n})
	return nil
}
func (r *recordingPrimitive) Kaddr(symbol string) (uint64, error) {
	a, ok := r.syms[symbol]
	if !ok {
		return 0, errors.New("unknown symbol: " + symbol)
	}
	return a, nil
}

func TestSelectMethodKnown(t *testing.T) {
	m, err := SelectMethod("modprobe_path")
	if err != nil {
		t.Fatalf("select modprobe_path: %v", err)
	}
	if m.Name() != "modprobe_path" {
		t.Errorf("got name %q", m.Name())
	}
	if _, err := SelectMethod("bogus"); err == nil {
		t.Error("expected error for bogus method")
	}
}

func TestModprobePathPrepareWritesScriptPathToKaddr(t *testing.T) {
	p := &recordingPrimitive{
		expName: "test",
		syms: map[string]uint64{
			"modprobe_path": 0x1000,
			"selinux_state": 0x2000,
		},
	}
	m := ModprobePathMethod{Prefix: "test"}
	a, err := m.Prepare(p)
	if err != nil {
		t.Fatalf("Prepare: %v", err)
	}
	if len(p.writes) != 1 {
		t.Fatalf("writes = %d, want 1", len(p.writes))
	}
	if p.writes[0].kaddr != 0x1000 {
		t.Errorf("write kaddr = %#x, want 0x1000", p.writes[0].kaddr)
	}
	if !bytes.HasSuffix(p.writes[0].content, []byte{0x00}) {
		t.Error("modprobe_path write should be NUL-terminated")
	}
	if !strings.HasPrefix(string(p.writes[0].content), a.ScriptPath) {
		t.Errorf("write content should start with script path %q", a.ScriptPath)
	}
	if len(p.memsets) != 1 || p.memsets[0].kaddr != 0x2000 {
		t.Errorf("expected one selinux memset at 0x2000, got %+v", p.memsets)
	}
}

func TestModprobePathPrepareSkipsSelinuxWhenSymbolMissing(t *testing.T) {
	p := &recordingPrimitive{
		expName: "test",
		syms:    map[string]uint64{"modprobe_path": 0x1000},
	}
	m := ModprobePathMethod{Prefix: "test"}
	if _, err := m.Prepare(p); err != nil {
		t.Fatalf("Prepare: %v", err)
	}
	if len(p.memsets) != 0 {
		t.Errorf("expected no selinux memset when symbol missing, got %+v", p.memsets)
	}
}

func TestCorePatternPrepareWritesPipeDirective(t *testing.T) {
	p := &recordingPrimitive{
		expName: "test",
		syms: map[string]uint64{
			"core_pattern":  0x3000,
			"selinux_state": 0x2000,
		},
	}
	m := CorePatternMethod{Prefix: "test"}
	if _, err := m.Prepare(p); err != nil {
		t.Fatalf("Prepare: %v", err)
	}
	if len(p.writes) != 1 {
		t.Fatalf("writes = %d, want 1", len(p.writes))
	}
	if p.writes[0].kaddr != 0x3000 {
		t.Errorf("write kaddr = %#x, want 0x3000", p.writes[0].kaddr)
	}
	if !bytes.HasPrefix(p.writes[0].content, []byte("|")) {
		t.Error("core_pattern write should be a pipe directive (leading |)")
	}
	if !bytes.Contains(p.writes[0].content, []byte("%P")) {
		t.Error("core_pattern write should contain %P")
	}
}

func TestMethodAvailableSkipsUnknownSymbol(t *testing.T) {
	p := &recordingPrimitive{expName: "test", syms: map[string]uint64{}}
	mm := ModprobePathMethod{}
	if ok, err := mm.Available(p); err != nil || ok {
		t.Errorf("modprobe Available with no symbol: ok=%v err=%v, want false/nil", ok, err)
	}
	cm := CorePatternMethod{}
	if ok, err := cm.Available(p); err != nil || ok {
		t.Errorf("core_pattern Available with no symbol: ok=%v err=%v, want false/nil", ok, err)
	}
}

func TestShellArtifactScriptContainsSetuidDrop(t *testing.T) {
	a := NewShellArtifact("test")
	s := a.RootShellScript()
	for _, want := range []string{"/bin/bash", "chmod 4755", a.RootBashPath, a.MarkerPath} {
		if !strings.Contains(s, want) {
			t.Errorf("script missing %q:\n%s", want, s)
		}
	}
}

func TestSelectMethodUnknownReturnsError(t *testing.T) {
	p := &recordingPrimitive{expName: "test", syms: map[string]uint64{}}
	_ = p
	if _, err := SelectMethod("bogus"); err == nil {
		t.Error("expected error for bogus method")
	}
}
