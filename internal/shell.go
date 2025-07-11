package internal

import (
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"io"
	"os"
	"os/exec"
	"syscall"
)

func InvokeRootShellBySu() {
	shell := exec.Command("su", "-", "root")
	shell.Stdout = os.Stdout
	shell.Stdin = os.Stdin
	shell.Stderr = os.Stderr
	shell.Run()
}

func InvokeRootShellBySuid(i io.Reader, o, e io.Writer) (err error) {
	// Incorrect call:
	//     os.Chmod("/bin/dash", 06755)
	//
	// Why it fails:
	// os.Chmod applies
	//     mode & os.ModePerm          // keeps only 0o777
	// so the set-uid, set-gid and sticky bits (0o4000/0o2000/0o1000)
	// are silently discarded, turning 06755 into 0755.
	//
	// Correct approaches:
	// 1. Use syscall.Chmod / unix.Chmod, which preserve all bits, or
	// 2. Call os.Chmod for the regular permission bits, then add the
	//    special bits with a syscall.
	// e.g.
	// syscall.Chmod("/bin/dash", 06755)
	// os.Chmod("/bin/dash", os.ModePerm|os.ModeSetuid|os.ModeSetgid)
	err = syscall.Chmod("/bin/sh", 04755)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	// prevent the text-busy error caused by overlayfs's copy-up
	syscall.Sync()
	shell := exec.Command("/bin/sh", "-p")
	shell.Stdin = i
	shell.Stdout = o
	shell.Stderr = e
	err = shell.Run()
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	return
}
