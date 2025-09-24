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
