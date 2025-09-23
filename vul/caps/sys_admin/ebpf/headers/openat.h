struct openat_args {
    int dirfd;
    const char *pathname;
    int flags;
};


/*
* return:
*   - 0: success
*   - non-zero: failure (e.g., invalid args)
*/
static __inline int parse_openat_args(struct bpf_raw_tracepoint_args *ctx, struct openat_args *args) {
    // int openat(int dirfd, const char *pathname, int flags, ...);
    struct pt_regs *regs = (struct pt_regs *)(ctx->args[0]);
    if (!regs) {
        return -1;
    }
    args->dirfd = (int)PT_REGS_PARM1_CORE(regs);
    args->pathname = (const char*)PT_REGS_PARM2_CORE(regs);
    if (!args->pathname) {
        return -2;
    }
    args->flags = (int)PT_REGS_PARM3_CORE(regs);
    return 0;
}
