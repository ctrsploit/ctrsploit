// #include "vmlinux.h"
// #include <bpf/bpf_tracing.h>

#ifndef NULL
#define NULL ((void *)0)
#endif

#ifndef PROC_DYNAMIC_FIRST
#define PROC_DYNAMIC_FIRST 0xF0000000
#endif

#ifndef MNT_NS_INODE_INIT
#define MNT_NS_INODE_INIT PROC_DYNAMIC_FIRST+1
#endif

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

static __inline bool is_root() {
    u64 uid_gid = bpf_get_current_uid_gid();
    u32 uid = uid_gid;
    return uid == 0;
}

static __inline bool is_host() {
    struct task_struct *task = (struct task_struct *) bpf_get_current_task();
    if (task == NULL)
        return false;
    struct nsproxy *namespaceproxy = BPF_CORE_READ(task, nsproxy);
    if (!namespaceproxy) {
        return false;
    }
    u32 mount_ns_id = BPF_CORE_READ(namespaceproxy, mnt_ns, ns.inum);
    return mount_ns_id == MNT_NS_INODE_INIT;
}

static __inline bool const_memcmp(const char *s1, const char *s2, int n) {
    for (int i = 0; i < n; i++) {
        if (s1[i] != s2[i]) {
            return false;
        }
    }
    return true;
}

static __inline bool starts_with(const char *cmdline, int cmd_len, const char *target, int target_size) {
    // clang will do "Loop Unrolling" here for small size,
    // but for sizeof(BIN_BASH) it will call real memcmp,
    // which will cause:
    // symbol "memcmp": unsatisfied program reference
    // so, use custom const_memcmp here.
    if (cmd_len >= target_size && const_memcmp(cmdline, target, target_size)) {
        return true;
    }
    return false;
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
