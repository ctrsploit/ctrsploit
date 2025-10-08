struct close_args {
    int fd;
};

static __inline int parse_close_args(struct bpf_raw_tracepoint_args *ctx, struct close_args *args) {
    // SYSCALL_DEFINE1(close, unsigned int, fd)
    struct pt_regs *regs = (struct pt_regs *)(ctx->args[0]);
    if (!regs) {
        return -1;
    }
    args->fd = (int)PT_REGS_PARM1_CORE(regs);
    return 0;
}