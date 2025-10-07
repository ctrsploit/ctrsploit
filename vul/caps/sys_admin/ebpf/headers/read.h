struct read_args {
    int fd;
    const char *buf;
    size_t count;
    int ret;
};

static __inline int parse_read_args(struct bpf_raw_tracepoint_args *ctx, struct read_args *args) {
    // SYSCALL_DEFINE3(read, unsigned int, fd, char __user *, buf, size_t, count)
    struct pt_regs *regs = (struct pt_regs *)(ctx->args[0]);
    if (!regs) {
        return -1;
    }
    args->fd = (int)PT_REGS_PARM1_CORE(regs);
    args->buf = (const char *)PT_REGS_PARM2_CORE(regs);
    if (!args->buf) {
        return -2;
    }
    args->count = (size_t)PT_REGS_PARM3_CORE(regs);
    args->ret = PT_REGS_RC_CORE(regs);
    return 0;
}
