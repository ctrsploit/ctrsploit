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
