package kubelet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
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
	objs = &bpfObjects{}
)

type Event struct {
	Path  string
	Token string
}

func Load(c chan Event) error {
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
	// 5. start processing events
	return processEvents(objs.Events, stopper, c)
}

func SetupBpf() (link.Link, error) {
	// 1. load pre-compiled programs and maps into the kernel
	if err := loadBpfObjects(objs, nil); err != nil {
		return nil, fmt.Errorf("loading objects: %w", err)
	}
	// 2. attach the program to the raw tracepoint (sys_enter)
	tp, err := link.AttachRawTracepoint(link.RawTracepointOptions{
		Name:    "sys_exit",
		Program: objs.RawTracepoint,
	})
	if err != nil {
		return nil, fmt.Errorf("opening raw tracepoint: %w", err)
	}
	return tp, nil
}

func processEvents(events *ebpf.Map, stopper chan os.Signal, notifier chan Event) error {
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
		pathname := util.Int8ToStr(event.Pathname[:])
		token := util.Int8ToStr(event.Token[:])
		log.Logger.Infof("pid: %d, fd=%d, pathname: %s\ntoken: %s", event.Pid, event.Fd, pathname, token)
		if notifier != nil {
			notifier <- Event{
				Path:  pathname,
				Token: token,
			}
		}
	}
}
