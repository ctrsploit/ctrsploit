package pipe_primitive

import (
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"os"
	"strings"
)

func getRootPasswdOffset() (offset int, err error) {
	content, err := os.ReadFile("/etc/passwd")
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	offset = strings.Index(string(content), "root:") + len("root:") - 1
	return
}
