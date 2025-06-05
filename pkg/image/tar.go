package image

import (
	"archive/tar"
	"bytes"
	"github.com/ssst0n3/awesome_libs/awesome_error"
)

type Tar struct {
	buf    bytes.Buffer
	writer *tar.Writer
}

func NewTar() *Tar {
	t := &Tar{}
	t.writer = tar.NewWriter(&t.buf)
	return t
}

func (t *Tar) Build() (bytes []byte, err error) {
	err = t.writer.Close()
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	bytes = t.buf.Bytes()
	return
}

func (t *Tar) Dir(name string) (err error) {
	err = t.writer.WriteHeader(&tar.Header{
		Typeflag: tar.TypeDir,
		Name:     name,
		Mode:     0o755,
	})
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	return
}

func (t *Tar) Reg(name string, content []byte) (err error) {
	err = t.writer.WriteHeader(&tar.Header{
		Typeflag: tar.TypeReg,
		Name:     name,
		Mode:     0o755,
		Size:     int64(len(content)),
	})
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	if len(content) > 0 {
		if _, err = t.writer.Write(content); err != nil {
			awesome_error.CheckErr(err)
			return
		}
	}
	return
}

func (t *Tar) Symlink(name string, linkname string) (err error) {
	err = t.writer.WriteHeader(&tar.Header{
		Typeflag: tar.TypeSymlink,
		Name:     name,
		Linkname: linkname,
		Mode:     0o777,
	})
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	return
}
