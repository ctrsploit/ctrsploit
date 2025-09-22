//go:build ignore
#include "vmlinux.h"
#include <bpf/bpf_core_read.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>

struct event {
    u32 uid;
};

struct {
    __uint(type, BPF_MAP_TYPE_RINGBUF);
    __uint(max_entries, 1 << 24);
    __type(value, struct event);
} events SEC(".maps");

SEC("raw_tracepoint/sys_exit")
int raw_tracepoint(struct bpf_raw_tracepoint_args *ctx) {
    return 0;
}