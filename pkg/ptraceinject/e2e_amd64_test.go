//go:build linux && amd64

package ptraceinject

import (
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// markerFile is created by the injected child's shellcode (creat syscall).
// Since the victim is started from within the test process, it shares the
// container's mount namespace, so the file is visible to the test.
const markerFile = "/tmp/pi"

// forkCreatExit is amd64 shellcode that follows the Inject convention:
//
//	fork(); if (child) { creat("/tmp/pi", 07); exit(0); } else { int3; }
//
// The parent traps on int3 so Inject can restore it; the child creates a
// marker file and exits, which the test detects.
var forkCreatExit = []byte{
	// fork()
	0x6a, 0x39, // push 57
	0x58,       // pop rax
	0x0f, 0x05, // syscall
	0x48, 0x85, 0xc0, // test rax, rax
	0x74, 0x01, // je +1 (jump over int3 into child code)
	0xcc, // int3 (parent traps here)

	// --- child: creat("/tmp/pi", 07) ---
	0x48, 0xb8, 0x2f, 0x74, 0x6d, 0x70, 0x2f, 0x70, 0x69, 0x00, // mov rax, "/tmp/pi\0"
	0x50,             // push rax
	0x48, 0x89, 0xe7, // mov rdi, rsp        ; rdi = path
	0x6a, 0x55, // push 85             ; SYS_creat (amd64)
	0x58,       // pop rax
	0x6a, 0x07, // push 7              ; mode = 07
	0x5e,       // pop rsi
	0x0f, 0x05, // syscall

	// exit(0)
	0x6a, 0x3c, // push 60             ; SYS_exit
	0x58,             // pop rax
	0x48, 0x31, 0xff, // xor rdi, rdi
	0x0f, 0x05, // syscall
}

// TestE2E_Inject verifies the full injection flow: a victim process is
// started, fork+creat+exit shellcode is injected, the forked child creates a
// marker file, and the victim continues running unaffected.
//
// This test requires CAP_SYS_PTRACE (and a permissive ptrace_scope). It is
// gated on TEST_ENV=exploitable, matching the e2e container setup used for
// vul/caps/sys_ptrace/pid_host.
func TestE2E_Inject(t *testing.T) {
	if os.Getenv("TEST_ENV") != "exploitable" {
		t.Skipf("Skipping e2e test; set TEST_ENV=exploitable (requires cap_sys_ptrace)")
	}

	// Clean up any stale marker from a previous run.
	os.Remove(markerFile)

	// Start a victim that spins in user space (no syscalls per iteration), so
	// its RIP is a user instruction and PtraceCont immediately executes the
	// injected shellcode. A blocking-syscall victim (e.g. sleep) would resume
	// the syscall instead and never reach the poked code.
	victim := exec.Command("sh", "-c", "while true; do :; done")
	if err := victim.Start(); err != nil {
		t.Fatalf("failed to start victim: %v", err)
	}
	victimPid := victim.Process.Pid
	t.Logf("victim pid: %d", victimPid)

	defer func() {
		victim.Process.Signal(syscall.SIGKILL)
		victim.Wait()
		os.Remove(markerFile)
	}()

	// Inject the fork+creat+exit shellcode.
	if err := Inject(victimPid, forkCreatExit); err != nil {
		t.Fatalf("Inject failed: %v", err)
	}

	// Wait for the marker file to appear (the forked child creates it).
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if _, err := os.Stat(markerFile); err == nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	_, err := os.Stat(markerFile)
	assert.NoError(t, err, "marker file %s not created by injected child", markerFile)

	// Verify the victim is still alive and running normally.
	err = victim.Process.Signal(syscall.Signal(0))
	assert.NoError(t, err, "victim process should still be alive after injection")
}
