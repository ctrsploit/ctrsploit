package bash

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
	"github.com/ctrsploit/sploit-spec/pkg/log"
	"github.com/davecgh/go-spew/spew"
)

// $BPF_CFLAGS are set by the Makefile
//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cflags $BPF_CFLAGS bpf ./bpf.c -- -I../headers

func Load() (err error) {
	// 1. gracefully handle shutdown on SIGINT and SIGTERM
	stopper := make(chan os.Signal, 1)
	signal.Notify(stopper, os.Interrupt, syscall.SIGTERM)
	// 2. allow the current process to lock memory for eBPF resources
	if err := rlimit.RemoveMemlock(); err != nil {
		return fmt.Errorf("removing memlock: %w", err)
	}
	// 3. load pre-compiled programs and maps into the kernel
	objs, tp, err := SetupBpf()
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
	// 4. start processing events
	return processEvents(objs.Events, stopper)
}

//goland:noinspection GoExportedFuncWithUnexportedType
func SetupBpf() (*bpfObjects, link.Link, error) {
	// 1. load pre-compiled programs and maps into the kernel
	objs := bpfObjects{}
	if err := loadBpfObjects(&objs, nil); err != nil {
		return nil, nil, fmt.Errorf("loading objects: %w", err)
	}
	// 2. attach the program to the raw tracepoint (sys_enter)
	tp, err := link.AttachRawTracepoint(link.RawTracepointOptions{
		Name:    "sys_exit",
		Program: objs.RawTracepoint,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("opening raw tracepoint: %w", err)
	}
	return &objs, tp, nil
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
				return
			}
			log.Logger.Errorf("reading from reader: %s", err)
			continue
		}
		// TODO: handle BigEndian
		if err := binary.Read(bytes.NewBuffer(record.RawSample), binary.LittleEndian, &event); err != nil {
			log.Logger.Errorf("parsing ringbuf event: %s", err)
			continue
		}
		spew.Dump(event)
	}
}
