package ptraceinject

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestMain implements a "victim mode" for e2e tests: when the env var
// PTRACEINJECT_VICTIM is set, the test binary becomes a pure user-space spin
// loop (for {}) with GOMAXPROCS=1 and GC disabled — no syscalls, no goroutine
// preemption — so RIP is always a user instruction when a tracer attaches.
// This makes ptrace code-injection deterministic. The e2e tests re-exec the
// test binary with this env var to obtain a reliable victim.
//
// A readiness file is written before the spin loop so the test can wait until
// the victim has fully entered the loop (past all Go runtime init) before
// attaching.
func TestMain(m *testing.M) {
	if os.Getenv("PTRACEINJECT_VICTIM") == "1" {
		runtime.GOMAXPROCS(1)
		debug.SetGCPercent(-1)
		pid := os.Getpid()
		readyFile := fmt.Sprintf("/tmp/ptraceinject_victim_ready_%d", pid)
		os.WriteFile(readyFile, []byte("ready"), 0644)
		for {
		}
	}
	os.Exit(m.Run())
}

// TestInject_NonExistentPid verifies that Inject fails cleanly when the target
// pid does not exist. This does not require cap_sys_ptrace and can run anywhere.
func TestInject_NonExistentPid(t *testing.T) {
	// 999999 is almost certainly not a valid pid.
	err := Inject(999999, []byte{0xcc})
	assert.Error(t, err)
}
