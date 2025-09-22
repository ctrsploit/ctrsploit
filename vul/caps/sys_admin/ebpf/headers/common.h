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

static __inline bool is_root() {
    u64 uid_gid = bpf_get_current_uid_gid();
    u32 uid = uid_gid;
    return uid == 0;
}

#ifndef NULL
#define NULL ((void *)0)
#endif

#ifndef PROC_DYNAMIC_FIRST
#define PROC_DYNAMIC_FIRST 0xF0000000
#endif

#ifndef MNT_NS_INODE_INIT
#define MNT_NS_INODE_INIT PROC_DYNAMIC_FIRST+1
#endif

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
