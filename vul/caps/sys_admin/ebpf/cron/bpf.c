//go:build ignore
#include "vmlinux.h"
#include <bpf/bpf_tracing.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_core_read.h>
#include "common.h"
#include "cron.h"
#include "openat.h"
#include "newfstatat.h"

char LICENSE[] SEC("license") = "GPL";

struct event {
    u32 pid;
};

struct {
    __uint(type, BPF_MAP_TYPE_RINGBUF);
    __uint(max_entries, 1 << 24);
    __type(value, struct event);
} events SEC(".maps");

static __inline int handle_exit_read(struct bpf_raw_tracepoint_args *ctx) {
    // read(5, "# /etc/crontab: system-wide cron"..., 4096) = 1136
    return 0;
}

#define CRONTAB_FILE "/etc/crontab"

static __inline int handle_exit_openat(struct bpf_raw_tracepoint_args *ctx) {
    // openat(AT_FDCWD, "/etc/crontab", O_RDONLY) = 5
    struct openat_args args = {};
    if (parse_openat_args(ctx, &args) != 0) {
        return -1;
    }
    bpf_printk("pathname=%s", args.pathname);
//    if (const_memcmp(args.pathname, CRONTAB_FILE, sizeof(CRONTAB_FILE)) != 0) {
//        return -1;
//    }
//    bpf_printk("/etc/crontab");
    // send event to user space
    struct event *e;
    e = bpf_ringbuf_reserve(&events, sizeof(*e), 0);
    if (!e) {
        return 1;
    }
    e->pid = bpf_get_current_pid_tgid() >> 32;
    bpf_ringbuf_submit(e, 0);
    return 0;
}

static __inline int handle_exit_newfstatat(struct bpf_raw_tracepoint_args *ctx) {
    // newfstatat(AT_FDCWD, "/etc/crontab", {st_mode=S_IFREG|0644, st_size=1136, ...}, 0) = 0
    struct newfstatat_args args = {};
    if (parse_newfstatat_args(ctx, &args) != 0) {
        return -1;
    }
    bpf_printk("[handle_exit_newfstatat] filename: %s", args.filename);
//    bpf_printk("st_mtime: %d", args.statbuf->st_mtime);
    return 0;
}

SEC("raw_tracepoint/sys_exit")
int raw_tracepoint(struct bpf_raw_tracepoint_args *ctx) {
    // filter comm==cron
    if (!is_cron()) {
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
    case 262:
        handle_exit_newfstatat(ctx);
    default:
        return 0;
    }
    return 0;
}