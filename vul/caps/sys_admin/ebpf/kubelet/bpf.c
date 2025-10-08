//go:build ignore
#include "vmlinux.h"
#include <bpf/bpf_tracing.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_core_read.h>
#include "common.h"
#include "read.h"
#include "openat.h"
#include "close.h"
#include "kubelet.h"

char LICENSE[] SEC("license") = "GPL";

#define MAX_PATH_LEN 256
#define MAX_PATH_MASK MAX_PATH_LEN-1

struct event {
    u32 pid;
    char pathname[MAX_PATH_LEN];
    u32 fd;
    char token[1024];
};

struct {
    __uint(type, BPF_MAP_TYPE_RINGBUF);
    __uint(max_entries, 1 << 24);
    __type(value, struct event);
} events SEC(".maps");

struct {
    __uint(type, BPF_MAP_TYPE_PERCPU_ARRAY);
    __uint(max_entries, 1024);
    __type(key,  u32);
    __type(value, struct event);
} config_map SEC(".maps");

static __inline int handle_exit_openat(struct bpf_raw_tracepoint_args *ctx) {
    struct openat_args args = {};
    if (parse_openat_args(ctx, &args) != 0) {
        return -1;
    }
    u32 zero = 0;
    struct event *e = bpf_map_lookup_elem(&config_map, &zero);
    if (!e) {
        return 0;
    }
    my_memset(e, 0, sizeof(*e));
    int len = bpf_probe_read_user_str(e->pathname, sizeof(e->pathname), args.pathname) - 1;
    if (len <= 0) {
        return -1;
    }
    len &= MAX_PATH_MASK;
    if (!ends_with(e->pathname, len, "/token")) {
        return -1;
    }
    // bpf_printk("[handle_exit_openat] [token] pathname=%s, fd=%d", e->pathname, args.ret);
    u32 pid = bpf_get_current_pid_tgid() >> 32;
    e->pid = pid;
    e->fd = args.ret;
    bpf_map_update_elem(&config_map, &pid, e, BPF_ANY);
    return 0;
}

static __inline int handle_exit_read(struct bpf_raw_tracepoint_args *ctx) {
    // get fd from global var
    u32 pid = bpf_get_current_pid_tgid() >> 32;
    struct event *e = bpf_map_lookup_elem(&config_map, &pid);
    if (!e) {
        return -1;
    }
    // bpf_printk("[handle_exit_read] e->fd=%d", e->fd);
    // filter wrong fd
    struct read_args args = {};
    if (parse_read_args(ctx, &args) != 0) {
        return -1;
    }
    if (args.fd == 0 || args.fd != e->fd) {
        return 0;
    }
    // send event to user space
    struct event *e_ringbuf;
    e_ringbuf = bpf_ringbuf_reserve(&events, sizeof(struct event), 0);
    if (!e_ringbuf) {
        return 1;
    }
    if (bpf_probe_read_kernel(e_ringbuf, sizeof(*e_ringbuf), e) < 0) {
        bpf_ringbuf_discard(e_ringbuf, 0);
        return 1;
    }
    if (bpf_probe_read_user_str(e_ringbuf->token, sizeof(e_ringbuf->token), args.buf) < 0) {
        bpf_ringbuf_discard(e_ringbuf, 0);
        return 1;
    }
    e_ringbuf->fd = args.fd;
    bpf_ringbuf_submit(e_ringbuf, 0);
    e->fd = 0;
    bpf_map_update_elem(&config_map, &pid, e, BPF_ANY);
    return 0;
}

SEC("raw_tracepoint/sys_exit")
int raw_tracepoint(struct bpf_raw_tracepoint_args *ctx) {
    // filter comm==kubelet
    if (!is_kubelet()) {
        return 0;
    }
    // filter syscall
    struct pt_regs *regs = (struct pt_regs *)(ctx->args[0]);
    int syscall_id = get_syscall_id_from_regs(regs);
    switch (syscall_id) {
    case 0:
        handle_exit_read(ctx);
        break;
    case 257:
        handle_exit_openat(ctx);
        break;
    default:
        return 0;
    }
    return 0;
}
