struct ioctl_args {
    int fd;
    unsigned long request;
    unsigned long arg;
};

/*
* return:
*   - 0: success
*   - non-zero: failure (e.g., invalid args)
*/
static __inline int parse_ioctl_args(struct bpf_raw_tracepoint_args *ctx, struct ioctl_args *args) {
    // int ioctl(int fd, unsigned long op, ...)
    struct pt_regs *regs = (struct pt_regs *)(ctx->args[0]);
    if (!regs) {
        return -1;
    }
    args->fd = (int)PT_REGS_PARM1_CORE(regs);
    args->request = (unsigned long)PT_REGS_PARM2_CORE(regs);
    args->arg = (unsigned long)PT_REGS_PARM3_CORE(regs);
    return 0;
}
