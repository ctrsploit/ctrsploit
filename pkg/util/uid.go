package util

import (
	"fmt"
	"os"
	"runtime"
	"syscall"

	"github.com/ctrsploit/sploit-spec/pkg/log"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"golang.org/x/sys/unix"
)

func WithFsuidFsgid(fsuid, fsgid int, f func() error) (err error) {
	err = unix.Setfsuid(fsuid)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	err = unix.Setfsgid(fsgid)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	err = f()
	return
}

// WithUid execute f() with uid
func WithUid(uid int, f func() error) (err error) {
	if os.Getuid() == uid {
		return f()
	}
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	origEUID := syscall.Geteuid()
	if err := syscall.Setresuid(uid, uid, origEUID); err != nil {
		return fmt.Errorf("failed to switch user: %w", err)
	}
	log.Logger.Infof("fsuid=%d, fgid=%d", syscall.Setfsuid(-1), unix.Setfsgid(-1))
	defer syscall.Setresuid(-1, origEUID, -1)
	// Execute the function
	err = f()

	return err
}

// WithUidAndCaps execute f() with uid and preserved original CapEff
func WithUidAndCaps(uid int, f func() error) (err error) {
	if os.Getuid() == uid {
		return f()
	}
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	origUID := syscall.Getuid()

	// Get original capabilities
	hdr := unix.CapUserHeader{Version: unix.LINUX_CAPABILITY_VERSION_3}
	data := [2]unix.CapUserData{}
	if err := unix.Capget(&hdr, &data[0]); err != nil {
		return fmt.Errorf("failed to get original capabilities: %w", err)
	}
	// Switch user using setresuid, keeping real UID as root
	if err := syscall.Setresuid(uid, uid, uid); err != nil {
		return fmt.Errorf("failed to switch user: %w", err)
	}
	defer syscall.Setresuid(origUID, origUID, origUID)

	// Preserve original capabilities
	if err := unix.Capset(&hdr, &data[0]); err != nil {
		return fmt.Errorf("failed to set capabilities: %w", err)
	}

	// Execute the function
	err = f()

	return err
}
