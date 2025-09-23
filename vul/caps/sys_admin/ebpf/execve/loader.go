package execve

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/ringbuf"
	"github.com/cilium/ebpf/rlimit"
	"github.com/ctrsploit/ctrsploit/pkg/util"
	"github.com/ctrsploit/sploit-spec/pkg/log"
)

// $BPF_CFLAGS are set by the Makefile
//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cflags $BPF_CFLAGS bpf ./bpf.c -- -I../headers

var (
	cfg     = bpfConfig{}
	objs    = &bpfObjects{}
	hostPid = make(chan uint32, 1)
)

func Load(path string, relative bool) (err error) {
	// 1. gracefully handle shutdown on SIGINT and SIGTERM
	stopper := make(chan os.Signal, 1)
	signal.Notify(stopper, os.Interrupt, syscall.SIGTERM)
	// 2. allow the current process to lock memory for eBPF resources
	if err := rlimit.RemoveMemlock(); err != nil {
		return fmt.Errorf("removing memlock: %w", err)
	}
	// 3. load pre-compiled programs and maps into the kernel
	tp, err := SetupBpf()
	if err != nil {
		return fmt.Errorf("loading BPF objects: %w", err)
	}
	defer func() {
		if err := objs.Close(); err != nil {
			log.Logger.Errorf("closing objects: %v", err)
		}
	}()
	defer func() {
		if err := tp.Close(); err != nil {
			log.Logger.Errorf("closing raw tracepoint: %v", err)
		}
	}()
	// 4. setup config_map
	go SetupConfig(path, relative)
	// 5. trigger ebpf to get container root path by ioctl
	if _, _, err := syscall.Syscall(syscall.SYS_IOCTL, os.Stdin.Fd(), uintptr(cfg.CallerId), 0); err != 0 {
		// ignore ENOTTY, it does not matter
		if !errors.Is(err, syscall.ENOTTY) {
			return fmt.Errorf("triggering ioctl failed: %w", err)
		}
	}
	// 6. start processing events
	return processEvents(objs.Events, stopper)
}

func SetupConfig(path string, relative bool) (err error) {
	if relative {
		pid := <-hostPid
		path = fmt.Sprintf("/proc/%d/root%s", pid, path)
	}
	path = fmt.Sprintf("%s\x00", path)
	cfg.LenCommand = uint32(len(path))
	copy(cfg.Command[:], util.StrToInt8(path))
	cfg.CallerId = rand.Uint32()
	key := int32(0)
	if err := objs.ConfigMap.Update(&key, &cfg, ebpf.UpdateAny); err != nil {
		return fmt.Errorf("updating config_map failed: %w", err)
	}
	log.Logger.Infof("set up command as: %q", util.Int8ToStr(cfg.Command[:len(path)]))
	return
}

//goland:noinspection GoExportedFuncWithUnexportedType
func SetupBpf() (link.Link, error) {
	// 1. load pre-compiled programs and maps into the kernel
	if err := loadBpfObjects(objs, nil); err != nil {
		return nil, fmt.Errorf("loading objects: %w", err)
	}
	// 2. attach the program to the raw tracepoint (sys_enter)
	tp, err := link.AttachRawTracepoint(link.RawTracepointOptions{
		Name:    "sys_enter",
		Program: objs.RawTracepoint,
	})
	if err != nil {
		return nil, fmt.Errorf("opening raw tracepoint: %w", err)
	}
	return tp, nil
}

func processEvents(events *ebpf.Map, stopper chan os.Signal) (err error) {
	rd, err := ringbuf.NewReader(events)
	if err != nil {
		return fmt.Errorf("opening ringbuf reader: %s", err)
	}
	defer func() {
		if err := rd.Close(); err != nil {
			log.Logger.Errorf("closing ringbuf reader: %v", err)
		}
	}()
	go func() {
		<-stopper
		if err := rd.Close(); err != nil {
			log.Logger.Fatalf("closing ringbuf reader: %s", err)
		}
	}()
	log.Logger.Info("Waiting for events..")
	var event bpfEvent
	for {
		record, err := rd.Read()
		if err != nil {
			if errors.Is(err, ringbuf.ErrClosed) {
				log.Logger.Info("Received signal, exiting..")
				return err
			}
			log.Logger.Errorf("reading from reader: %s", err)
			continue
		}
		// TODO: handle BigEndian
		if err := binary.Read(bytes.NewBuffer(record.RawSample), binary.LittleEndian, &event); err != nil {
			log.Logger.Errorf("parsing ringbuf event: %s", err)
			continue
		}
		if event.Loader {
			log.Logger.Infof("Host pid: %d", event.Pid)
			hostPid <- event.Pid
		} else {
			pathname := util.Int8ToStr(event.Pathname[:event.LenPathname])
			log.Logger.Infof("pid: %d, pathname: %s, injected: %t", event.Pid, pathname, event.Injected)
		}
	}
}
