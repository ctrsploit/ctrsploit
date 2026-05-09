package pipeprimitive

import (
	"bytes"
	"fmt"
	"os"
)

func ImagePollution(primitive Primitive, source, destination string) error {
	payload, err := os.ReadFile(source)
	if err != nil {
		return fmt.Errorf("read image pollution source %q: %w", source, err)
	}

	if err := WriteImage(primitive, destination, payload); err != nil {
		return fmt.Errorf("pollute image destination %q with source %q: %w", destination, source, err)
	}
	return nil
}

func WriteImage(primitive Primitive, path string, payload []byte) error {
	minOffset := primitive.MinOffset()
	if minOffset < 0 {
		return fmt.Errorf("%s has invalid min offset %d", primitive.GetExpName(), minOffset)
	}
	if int64(len(payload)) < minOffset {
		return fmt.Errorf("payload length %d is smaller than min offset %d", len(payload), minOffset)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read target image %q: %w", path, err)
	}
	if int64(len(content)) < minOffset {
		return fmt.Errorf("target image %q length %d is smaller than min offset %d", path, len(content), minOffset)
	}

	prefixLen := int(minOffset)
	if !bytes.Equal(content[:prefixLen], payload[:prefixLen]) {
		return fmt.Errorf(
			"target image %q prefix mismatch: first %d bytes are % x, payload wants % x",
			path, prefixLen, content[:prefixLen], payload[:prefixLen],
		)
	}

	if prefixLen == len(payload) {
		return nil
	}
	writePayload := payload[prefixLen:]
	if err := primitive.Write(path, minOffset, writePayload); err != nil {
		return fmt.Errorf(
			"%s write target image %q at offset %d with %d bytes: %w",
			primitive.GetExpName(), path, minOffset, len(writePayload), err,
		)
	}
	return nil
}
