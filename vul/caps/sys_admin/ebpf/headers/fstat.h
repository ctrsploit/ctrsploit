#include "vmlinux.h"

struct fstat_args {
    int fd;
    struct stat *stat_buf;
};

static __inline int parse_fstat_args(struct bpf_raw_tracepoint_args *ctx, struct fstat_args *args) {
    // SYSCALL_DEFINE2(newfstat, unsigned int, fd, struct stat __user *, statbuf)
    struct pt_regs *regs = (struct pt_regs *)(ctx->args[0]);
    if (!regs) {
        return -1;
    }
    args->fd = (int)PT_REGS_PARM1_CORE(regs);
    args->stat_buf = (struct stat*)PT_REGS_PARM2_CORE(regs);
    if (!args->stat_buf) {
        return -2;
    }
    return 0;
}