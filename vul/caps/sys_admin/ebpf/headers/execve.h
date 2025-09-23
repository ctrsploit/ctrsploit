//#include "vmlinux.h"

struct execve_args {
    char *pathname;
    char *argv;
    char *envp;
};

/*
* return:
*   - 0: success
*   - non-zero: failure (e.g., invalid args)
*/
static __inline int parse_execve_args(struct bpf_raw_tracepoint_args *ctx, struct execve_args *args) {
    // int execve(const char *pathname, char *const _Nullable argv[], char *const _Nullable envp[]);
    struct pt_regs *regs = (struct pt_regs *)(ctx->args[0]);
    if (!regs) {
        return -1;
    }
    args->pathname = (char*)PT_REGS_PARM1_CORE(regs);
    if (!args->pathname) {
        return -2;
    }
    return 0;
}
