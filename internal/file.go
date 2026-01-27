package internal

import (
	"bytes"
	"fmt"
	"os"
	"strconv"

	"github.com/ssst0n3/awesome_libs/awesome_error"
)

func CheckPathExists(path string) (bool, error) {
	_, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("checkPathExists: %w", err)
	}
	return true, nil
}

func ReadIntFromFile(path string) (result int, err error) {
	content, err := os.ReadFile(path)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	content = bytes.TrimSpace(content)
	result, err = strconv.Atoi(string(content))
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	return
}

func ReplaceContent(path string, old, new []byte) (err error) {
	content, err := os.ReadFile(path)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	err = os.WriteFile(path, bytes.Replace(content, old, new, -1), 0)
	if err != nil {
		awesome_error.CheckWarning(err)
		return
	}
	return
}

func ReadUint64(file string) (uint64, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return 0, fmt.Errorf("ReadUint64(%s): failed to read file: %w", file, err)
	}
	content = bytes.TrimSpace(content)
	result, err := strconv.ParseUint(string(content), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("ReadUint64(%s): failed to parse content: %w", file, err)
	}
	return result, nil
}
