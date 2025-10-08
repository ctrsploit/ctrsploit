#include "common.h"
#define TASK_COMM_LEN 16

static __inline bool is_kubelet() {
    char comm[TASK_COMM_LEN] = {0};
    bpf_get_current_comm(&comm, sizeof(comm));
    return const_memcmp(comm, "kubelet", 7);
}
