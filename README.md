# ctrsploit: A penetration toolkit for container environment

[中文文档](./README-ZH.md)

ctrsploit [kənˈteɪnər splɔɪt] , follows [sploit-spec](https://github.com/ctrsploit/sploit-spec) v0.4.3

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
wget -O ctrsploit https://github.com/ctrsploit/ctrsploit/releases/latest/download/ctrsploit_linux_amd64 && chmod +x ctrsploit
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

### gather information

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

where

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

### run a exploit

```
root@2aa13a052102:/# ./ctrsploit exploit
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
./checksec_linux_amd64 auto
[N]  cap_sys_admin      # Container can be escaped when has cap_sys_admin and use cgroups v1
[N]  host_net_ns        # The network namespace of the host is shared
...
```

## Details

### env

| command                               | alias | description                                                              |
|---------------------------------------|-------|--------------------------------------------------------------------------|
| [auto](./env/auto)                    |       | auto gather environment information                                      |
| [where](./env/where)                  | w     | detect whether you are in the container, and which type of the container |
| [storage-driver](./env/storagedriver) | sd    | detect storage driver type and extend information                        |
| [cgroups](./env/cgroups)              | c     | gather cgroup information                                                |
| [capability](./env/capability)        | cap   | show the capability of pid 1 and current process                         |
| [seccomp](./env/seccomp)              | sc    | show the seccomp info                                                    |
| [apparmor](./env/apparmor)            | a     | show the apparmor info                                                   |
| [namespace](./env/namespace)          | n, ns | check namespace is host ns                                               |

### exploit

| exploit                                                                                  | alias               | description                                                                                                                                                                         |
|------------------------------------------------------------------------------------------|---------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [cgroupv1-release_agent-unknown_rootfs](./exploit/cgroupv1-release_agent-unknown_rootfs) | ra3                 | escape tech by using the notify_on_release of cgroup v1 without known rootfs                                                                                                        |
| [cve-2021-22555_ubuntu18.04](./exploit/CVE-2021-22555_ubuntu18.04)                       | 22555               | escape tech by using the CVE-2021-22555 (ubuntu18.04)                                                                                                                               |
| [cve-2025-23266](./vul/cve-2025-23266)                                                   | 23266               | nvidia-container-toolkit CVE-2025-23266 GPU container escape                                                                                                                        | 
| [cve-2025-47290](./vul/cve-2025-47290)                                                   | 47290               | containerd cve-2025-47290 host filesystem access during image unpack                                                                                                                | 
| [cve-2024-0132](./vul/cve-2024-0132)                                                     | 0132                | nvidia-container-toolkit CVE-2024-0132 GPU container escape                                                                                                                         | 
| [cve-2024-23650](./vul/cve-2024-23650)                                                   | 23650               | buildkitd cve-2024-23650 panic when incorrect parameters sent from frontend                                                                                                         | 
| [cve-2022-39253](./vul/cve-2022-39253)                                                   | 39253               | docker build host file read by git CVE-2022-39253                                                                                                                                   | 
| [cve-2020-15257](./vul/cve-2020-15257)                                                   | 15257               | containerd cve-2020-15257 host network container escape                                                                                                                             | 
| [cve-2016-8867](./vul/cve-2016-8867)                                                     | 8867                | runc cve-2016-8867 container normal user privilege escalation by ambient capabilities                                                                                               | 
| [release_agent](vul/caps/sys_admin/release_agent)                                        | ra                  | escape tech by using the notify_on_release of cgroup v1                                                                                                                             |
| [shocker](vul/caps/shocker)                                                              | cap_dac_read_search | Container escape with CAP_DAC_READ_SEARCH, alias shocker, found by Sebastian Krahmer (stealth) in 2014.                                                                             |
| [naked](./vul/naked)                                                                     |                     | We call containers running without seccomp, AppArmor, or SELinux enabled 'naked containers', which leaves them highly vulnerable to kernel exploits and potential container escapes |
| [ptrace-pid-host](./vul/caps/sys_ptrace/pid_host)                                        | pid                 | Container can be escaped when has cap_sys_ptrace and host pid namespace                                                                                                             |
| [host-pid-proc-root](./vul/namespace/pid/proc_root)                                      | proc                | Host level pid ns escape via /proc/[pid]/root, accessing host or other container's filesystem                                                                                       |
| [docker.sock](./vul/shared-socket/docker-sock)                                           | docker              | Container with shared docker.sock can be escaped via creating privileged containers                                                                                                 |
| [ebpf-bash](./vul/caps/sys_admin/ebpf/bash)                                              | bash                | Ebpf escape by bash                                                                                                                                                                 |
| [ebpf-execve](./vul/caps/sys_admin/ebpf/execve)                                          | execve              | Ebpf escape by execve                                                                                                                                                               |
| [ebpf-cron](./vul/caps/sys_admin/ebpf/cron)                                              | cron                | Ebpf escape by cron                                                                                                                                                                 |
| [ebpf-kubelet](./vul/caps/sys_admin/ebpf/kubelet)                                        | cron                | Ebpf escape by leaking k8s service account token via kubelet                                                                                                                        |

### vul

* :heavy_check_mark: : Fully Supported
* :o: : Partially Supported
* :bug: : Known Bug
* :x: : Not Supported
* :calendar: : Planned
* :heavy_minus_sign: : Not Applicable

| vul                                                                     | check              | exploit            | test                | doc                | video              |
|-------------------------------------------------------------------------|--------------------|--------------------|---------------------|--------------------|--------------------|
| [cve-2016-8867](./vul/cve-2016-8867)                                    | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark:  | :heavy_check_mark: | :heavy_check_mark: |
| CVE-2017-1002101                                                        | :x:                | :heavy_check_mark: | :x:                 | :x:                | :x:                |
| CVE-2019-5736                                                           | :x:                | :heavy_check_mark: | :x:                 | :x:                | :x:                |
| CVE-2019-14271                                                          | :x:                | :heavy_check_mark: | :x:                 | :x:                | :x:                |
| CVE-2019-16884                                                          | :x:                | :heavy_check_mark: | :x:                 | :x:                | :x:                |
| CVE-2020-8555                                                           | :x:                | :heavy_check_mark: | :x:                 | :x:                | :x:                |
| CVE-2020-8558                                                           | :x:                | :heavy_check_mark: | :x:                 | :x:                | :x:                |
| CVE-2020-15157                                                          | :x:                | :heavy_check_mark: | :x:                 | :x:                | :x:                |
| [cve-2020-15257](./vul/cve-2020-15257)                                  | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark:  | :heavy_check_mark: | :heavy_check_mark: |
| CVE-2021-21285                                                          | :x:                | :heavy_check_mark: | :x:                 | :x:                | :x:                |
| CVE-2021-22555                                                          | :x:                | :heavy_check_mark: | :x:                 | :x:                | :x:                |
| CVE-2021-41091                                                          | :x:                | :heavy_check_mark: | :x:                 | :x:                | :x:                |
| [cve-2021-25741](./vul/cve-2021-25741)                                  | :heavy_check_mark: | :heavy_check_mark: | :x:                 | :heavy_check_mark: | :x:                |
| [cve-2022-39253](./vul/cve-2022-39253)                                  | :o:                | :heavy_check_mark: | :heavy_check_mark:  | :heavy_check_mark: | :heavy_check_mark: |
| [cve-2022-0492](./vul/cve-2022-0492)                                    | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark:  | :heavy_check_mark: | :heavy_check_mark: |
| [cve-2024-0132](./vul/cve-2024-0132)                                    | :o:                | :heavy_check_mark: | :heavy_check_mark:  | :heavy_check_mark: | :heavy_check_mark: |
| [cve-2024-23650](./vul/cve-2024-23650)                                  | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark:  | :heavy_check_mark: | :heavy_check_mark: |
| [cve-2025-23266](./vul/cve-2025-23266)                                  | :o:                | :heavy_check_mark: | :heavy_check_mark:  | :heavy_check_mark: | :heavy_check_mark: |
| [cve-2025-47290](./vul/cve-2025-47290)                                  | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark:  | :heavy_check_mark: | :heavy_check_mark: |
| caps                                                                    | :heavy_minus_sign: | :heavy_minus_sign: | :heavy_minus_sign:  | :heavy_minus_sign: | :heavy_minus_sign: |
| └─[shocker](./vul/caps/shocker)                                         | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark:️ | :heavy_check_mark: | :heavy_check_mark: |
| └─[sys_admin](./vul/caps/sys_admin)                                     | :heavy_check_mark: | :heavy_minus_sign: | :heavy_minus_sign:  | :heavy_minus_sign: | :heavy_minus_sign: |
| &emsp;└─[ebpf](./vul/caps/sys_admin/ebpf)                               | :heavy_check_mark: | :heavy_minus_sign: | :heavy_minus_sign:  | :heavy_minus_sign: | :heavy_minus_sign: |
| &emsp;&emsp;└─[ebpf-bash](./vul/caps/sys_admin/ebpf/bash)               | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark:  | :heavy_check_mark: | :heavy_check_mark: |
| &emsp;&emsp;└─[ebpf-cron](./vul/caps/sys_admin/ebpf/cron)               | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark:  | :heavy_check_mark: | :heavy_check_mark: |
| &emsp;&emsp;└─[ebpf-execve](./vul/caps/sys_admin/ebpf/execve)           | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark:  | :heavy_check_mark: | :heavy_check_mark: |
| &emsp;&emsp;└─[ebpf-kubelet](./vul/caps/sys_admin/ebpf/kubelet)         | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark:  | :heavy_check_mark: | :heavy_check_mark: |
| &emsp;&emsp;└─ebpf-sshd                                                 | :calendar:         | :calendar:         | :calendar:          | :calendar:         | :calendar:         |
| └─[bpf](./vul/caps/bpf)                                                 | :heavy_minus_sign: | :heavy_minus_sign: | :heavy_minus_sign:  | :heavy_minus_sign: | :heavy_minus_sign: |
| &emsp;└─[release_agent](./vul/caps/sys_admin/release_agent)             | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark:  | :heavy_check_mark: | :heavy_check_mark: |
| &emsp;└─[ebpf](./vul/caps/sys_admin/ebpf) (same as caps/sys_admin/ebpf) | :heavy_check_mark: | :heavy_minus_sign: | :heavy_minus_sign:  | :heavy_minus_sign: | :heavy_minus_sign: |
| └─[sys_ptrace](./vul/caps/sys_ptrace)                                   | :heavy_check_mark: | :heavy_minus_sign: | :heavy_minus_sign:  | :heavy_minus_sign: | :heavy_minus_sign: |
| &emsp;└─[pid_host](./vul/caps/sys_ptrace/pid_host)                      | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark:  | :heavy_check_mark: | :heavy_check_mark: |
| └─sys_module                                                            | :calendar:         | :calendar:         | :calendar:          | :calendar:         | :calendar:         |
| └─net_admin                                                             | :calendar:         | :calendar:         | :calendar:          | :calendar:         | :calendar:         |
| [naked](./vul/naked)                                                    | :heavy_check_mark: | :heavy_minus_sign: | :heavy_check_mark:  | :x:                | :x:                |
| namespace                                                               | :heavy_minus_sign: | :heavy_minus_sign: | :heavy_minus_sign:  | :heavy_minus_sign: | :heavy_minus_sign: |
| └─[net](./vul/namespace/net)                                            | :heavy_minus_sign: | :x:                | :x:                 | :x:                | :x:                |
| └─[pid](./vul/namespace/pid)                                            | :heavy_minus_sign: | :heavy_minus_sign: | :heavy_minus_sign:  | :heavy_minus_sign: | :heavy_minus_sign: |
| &emsp;└─[proc_root](./vul/namespace/pid/proc_root)                      | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark:  | :heavy_check_mark: | :heavy_check_mark: |
| shared-socket                                                           | :heavy_minus_sign: | :heavy_minus_sign: | :heavy_minus_sign:  | :heavy_minus_sign: | :heavy_minus_sign: |
| └─[docker.sock](./vul/shared-socket/docker-sock)                        | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark:  | :heavy_check_mark: | :heavy_check_mark: |
| └─containerd.sock                                                       | :calendar:         | :calendar:         | :calendar:          | :calendar:         | :calendar:         |                                  
| crash                                                                   | :x:                | :heavy_check_mark: | :x:                 | :x:                | :x:                |
| dirty-pipe                                                              | :x:                | :heavy_check_mark: | :x:                 | :x:                | :x:                |

### helper

| helper                                  | alias                    | description                    |
|-----------------------------------------|--------------------------|--------------------------------|
| [cve-2021-3493](./helper/cve-2021-3493) | ubuntu-overlayfs-pe,3493 | Ubuntu OverlayFS Local Privesc |

### checksec

Just execute `ctrsploit checksec` or standalone binary file `checksec`.
