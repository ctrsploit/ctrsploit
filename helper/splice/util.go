package splice

import (
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"io/ioutil"
	"strings"
)

func getRootPasswdOffset() (offset int, err error) {
	content, err := ioutil.ReadFile("/etc/passwd")
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	offset = strings.Index(string(content), "root:") + len("root:") - 1
	return
}
