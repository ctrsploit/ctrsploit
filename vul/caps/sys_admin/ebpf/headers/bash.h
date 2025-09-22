// #include "vmlinux.h"
// #include <bpf/bpf_tracing.h>

#define BASH "bash"
#define BIN_BASH "/bin/bash"
#define USR_BIN_BASH "/usr/bin/bash"

static inline int get_syscall_id_from_regs(struct pt_regs *regs)
{
    int id = -1; // Unsupported architecture
#if defined(bpf_target_x86)
    id = BPF_CORE_READ(regs, orig_ax);
#elif defined(bpf_target_arm64)
    id = BPF_CORE_READ(regs, syscallno);
#endif
    return id;
}

// return the length of cmdline, or 0 on error
static __inline int get_cmdline(char *buf, int size) {
    struct task_struct *task = (struct task_struct *) bpf_get_current_task();
    if (!task) {
        return 0;
    }
    long unsigned int args_start = BPF_CORE_READ(task, mm, arg_start);
    long unsigned int args_end = BPF_CORE_READ(task, mm, arg_end);
    if (args_end <= args_start) {
        return 0;
    }
    int len = (args_end - args_start) & 0x7F;
    if (len >= size) {
        len = size - 1;
    }
    if (bpf_probe_read_user(buf, len, (const void *)args_start)) {
        return 0;
    }
    buf[len] = '\0';
    return len;
}

static __inline bool cmdline_starts_with_bash(char *cmdline, int len) {
    return starts_with(cmdline, len, BASH, sizeof(BASH)-1) ||
           starts_with(cmdline, len, BIN_BASH, sizeof(BIN_BASH)-1) ||
           starts_with(cmdline, len, USR_BIN_BASH, sizeof(USR_BIN_BASH)-1);
}

static __inline bool is_bash() {
    char cmdline[0x7F+1] = {0};
    int len = get_cmdline(cmdline, sizeof(cmdline));
    if (len <= 0) {
        return false;
    }
    return cmdline_starts_with_bash(cmdline, len);
}
