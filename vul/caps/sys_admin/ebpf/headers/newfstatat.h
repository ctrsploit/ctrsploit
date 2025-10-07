#include "vmlinux.h"

struct newfstatat_args {
    int dirfd;
    const char *filename;
    struct stat *stat_buf;
    int flags;
};

/*
* return:
*   - 0: success
*   - non-zero: failure (e.g., invalid args)
*/
static __inline int parse_newfstatat_args(struct bpf_raw_tracepoint_args *ctx, struct newfstatat_args *args) {
    // SYSCALL_DEFINE4(newfstatat, int, dfd, const char __user *, filename, struct stat __user *, statbuf, int, flag);
    struct pt_regs *regs = (struct pt_regs *)(ctx->args[0]);
    if (!regs) {
        return -1;
    }
    args->dirfd = (int)PT_REGS_PARM1_CORE(regs);
    args->filename = (const char*)PT_REGS_PARM2_CORE(regs);
    if (!args->filename) {
        return -2;
    }
    args->stat_buf = (struct stat*)PT_REGS_PARM3_CORE(regs);
    if (!args->stat_buf) {
        return -3;
    }
    args->flags = (int)PT_REGS_PARM4_CORE(regs);
    return 0;
}
