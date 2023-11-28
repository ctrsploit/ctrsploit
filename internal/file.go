package internal

import (
	"bytes"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"os"
	"strconv"
)

func CheckPathExists(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
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
