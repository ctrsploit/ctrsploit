package crash

import (
	"context"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	defaultOOMChunkBytes = 16 << 20
	maxOOMBytesFallback  = 512 << 20
)

type OOM struct {
	PID           int
	Delay         time.Duration
	ProcRoot      string
	CgroupRoot    string
	MaxBytes      uint64
	ChunkBytes    int
	TargetOOMBias int
	SelfOOMBias   int
}

func (o OOM) Name() string {
	return MethodOOM
}

func (o OOM) Trigger(ctx context.Context) error {
	if err := wait(ctx, o.Delay); err != nil {
		return err
	}

	pid := normalizedPID(o.PID)
	if err := writeOOMScoreAdj(o.ProcRoot, pid, o.TargetOOMBias); err != nil {
		return err
	}
	if err := writeOOMScoreAdj(o.ProcRoot, os.Getpid(), o.SelfOOMBias); err != nil {
		return err
	}

	maxBytes, err := o.limit()
	if err != nil {
		return err
	}
	chunkBytes := o.ChunkBytes
	if chunkBytes <= 0 {
		chunkBytes = defaultOOMChunkBytes
	}

	var held [][]byte
	var allocated uint64
	for allocated < maxBytes {
		if err := ctx.Err(); err != nil {
			return err
		}
		size := uint64(chunkBytes)
		if remaining := maxBytes - allocated; remaining < size {
			size = remaining
		}
		if size > uint64(math.MaxInt) {
			return fmt.Errorf("oom chunk size %d exceeds int max", size)
		}
		chunk := make([]byte, int(size))
		touchPages(chunk)
		held = append(held, chunk)
		allocated += size
	}
	return fmt.Errorf("%w: exhausted configured OOM allocation limit %d without container exit", ErrUnsupported, maxBytes)
}

func (o OOM) limit() (uint64, error) {
	if o.MaxBytes > 0 {
		return o.MaxBytes, nil
	}
	path, err := cgroupControlPath(o.ProcRoot, o.CgroupRoot, normalizedPID(o.PID), "memory.max")
	if err != nil {
		return 0, err
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("read %s: %w", path, err)
	}
	value := strings.TrimSpace(string(content))
	if value == "" || value == "max" {
		return maxOOMBytesFallback, nil
	}
	limit, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse %s value %q: %w", path, value, err)
	}
	return limit + uint64(defaultOOMChunkBytes), nil
}

func writeOOMScoreAdj(procRoot string, pid int, score int) error {
	if procRoot == "" {
		procRoot = "/proc"
	}
	path := filepath.Join(procRoot, strconv.Itoa(pid), "oom_score_adj")
	if err := writeFile(path, []byte(strconv.Itoa(score)), 0); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

func touchPages(buf []byte) {
	const pageSize = 4096
	for i := 0; i < len(buf); i += pageSize {
		buf[i] = 1
	}
	if len(buf) > 0 {
		buf[len(buf)-1] = 1
	}
}
