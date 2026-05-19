//go:build ignore

package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	asmPath = "material/suid-shell-amd64.s"
	binPath = "material/suid-shell-amd64.bin"

	loadBase   = 0x400000
	textOffset = 0x78
)

func main() {
	dir, err := os.MkdirTemp("", "ctrsploit-suid-payload-*")
	must(err)
	defer os.RemoveAll(dir)

	objPath := filepath.Join(dir, "suid-shell-amd64.o")
	textPath := filepath.Join(dir, "suid-shell-amd64.text")
	run("as", "--64", "-o", objPath, asmPath)
	run("objcopy", "-O", "binary", "-j", ".text", objPath, textPath)

	text, err := os.ReadFile(textPath)
	must(err)
	payload := buildMinimalELF(text)

	must(os.WriteFile(binPath, payload, 0o755))
}

func buildMinimalELF(text []byte) []byte {
	var out bytes.Buffer
	out.Grow(textOffset + len(text))

	out.Write([]byte{
		0x7f, 'E', 'L', 'F',
		2, 1, 1, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
	})
	write16(&out, 2)    // ET_EXEC
	write16(&out, 0x3e) // EM_X86_64
	write32(&out, 1)    // EV_CURRENT
	write64(&out, loadBase+textOffset)
	write64(&out, 64) // e_phoff
	write64(&out, 0)  // e_shoff
	write32(&out, 0)  // e_flags
	write16(&out, 64) // e_ehsize
	write16(&out, 56) // e_phentsize
	write16(&out, 1)  // e_phnum
	write16(&out, 0)  // e_shentsize
	write16(&out, 0)  // e_shnum
	write16(&out, 0)  // e_shstrndx

	write32(&out, 1) // PT_LOAD
	write32(&out, 7) // PF_R | PF_W | PF_X
	write64(&out, textOffset)
	write64(&out, loadBase+textOffset)
	write64(&out, loadBase+textOffset)
	write64(&out, len(text))
	write64(&out, len(text))
	write64(&out, 8)

	if out.Len() != textOffset {
		panic(fmt.Sprintf("ELF header length = %d, want %d", out.Len(), textOffset))
	}
	out.Write(text)
	return out.Bytes()
}

func run(name string, args ...string) {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s %v failed: %v\n%s", name, args, err, out)
		os.Exit(1)
	}
}

func write16(out *bytes.Buffer, v int) {
	var b [2]byte
	binary.LittleEndian.PutUint16(b[:], uint16(v))
	out.Write(b[:])
}

func write32(out *bytes.Buffer, v int) {
	var b [4]byte
	binary.LittleEndian.PutUint32(b[:], uint32(v))
	out.Write(b[:])
}

func write64(out *bytes.Buffer, v int) {
	var b [8]byte
	binary.LittleEndian.PutUint64(b[:], uint64(v))
	out.Write(b[:])
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
