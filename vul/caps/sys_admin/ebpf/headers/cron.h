#define TASK_COMM_LEN 16

static __inline bool is_cron() {
    char comm[TASK_COMM_LEN] = {0};
    bpf_get_current_comm(&comm, sizeof(comm));
    // 4 bytes is allow to use memcpy, the compiler will expand it
    return __builtin_memcmp(comm, "cron", 4) == 0;
//    if (comm[0] != 'c' || comm[1] != 'r' || comm[2] != 'o' || comm[3] != 'n' || comm[4] != '\0') {
//        return false;
//    }
//    return true
}
