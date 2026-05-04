package pipeprimitive

import (
	"bytes"
	"fmt"
	"os"
)

func rootPasswdPasswordOffset(path string) (int, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("read passwd file %q: %w", path, err)
	}
	return passwdPasswordOffset(content, "root")
}

func passwdPasswordOffset(content []byte, username string) (int, error) {
	prefix := []byte(username + ":")
	offset := 0
	for len(content) > 0 {
		line := content
		if newline := bytes.IndexByte(content, '\n'); newline >= 0 {
			line = content[:newline]
		}

		if bytes.HasPrefix(line, prefix) {
			return offset + len(prefix) - 1, nil
		}

		if len(line) == len(content) {
			break
		}
		offset += len(line) + 1
		content = content[len(line)+1:]
	}

	return 0, fmt.Errorf("user %q not found", username)
}
