package copyfail

import (
	"fmt"
	"io"
	"os"
)

const ChunkSize = 4

type Write4Func func(fd int, offset int64, block [ChunkSize]byte) error

func Write(path string, offset int64, content []byte, write4 Write4Func) error {
	if offset < 0 {
		return fmt.Errorf("invalid negative offset %d", offset)
	}
	if len(content) == 0 {
		return nil
	}

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open target file %q: %w", path, err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("stat target file %q: %w", path, err)
	}
	size := info.Size()
	if offset > size {
		return fmt.Errorf("offset %d is outside target file %q with size %d", offset, path, size)
	}
	if int64(len(content)) > size-offset {
		return fmt.Errorf(
			"write length %d at offset %d would enlarge target file %q with size %d",
			len(content), offset, path, size,
		)
	}

	for i := 0; i < len(content); i += ChunkSize {
		currentOffset := offset + int64(i)
		end := i + ChunkSize
		if end > len(content) {
			end = len(content)
		}

		var block [ChunkSize]byte
		written := copy(block[:], content[i:end])
		if written < ChunkSize {
			if currentOffset+ChunkSize > size {
				return fmt.Errorf(
					"cannot preserve trailing bytes for target file %q at offset %d: 4-byte primitive would cross EOF",
					path, currentOffset,
				)
			}
			if _, err := file.ReadAt(block[written:], currentOffset+int64(written)); err != nil {
				if err == io.EOF {
					return fmt.Errorf(
						"cannot preserve trailing bytes for target file %q at offset %d: short read",
						path, currentOffset+int64(written),
					)
				}
				return fmt.Errorf(
					"read target file %q to preserve trailing bytes at offset %d: %w",
					path, currentOffset+int64(written), err,
				)
			}
		}

		if err := write4(int(file.Fd()), currentOffset, block); err != nil {
			return fmt.Errorf("copy-fail write %q at offset %d: %w", path, currentOffset, err)
		}
	}
	return nil
}
