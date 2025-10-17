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
* `-` : Not Applicable

| vul                                                             | desc                                                                                                                                                                                            |
|-----------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [cve-2016-8867](./vul/cve-2016-8867)                            | ambient capabilities allow local users to gain <br>privileges                                                                                                                                   |
| cve-2016-9962                                                   |                                                                                                                                                                                                 |
| CVE-2017-1002101                                                |                                                                                                                                                                                                 |
| [cve-2019-5736](./vul/cve-2019-5736)                            | escape by overwrite runc executable file via <br>/proc/self/exe                                                                                                                                 |
| └─[exec](./vul/cve-2019-5736/exec)                              | cve-2019-5736 exploit via runc exec process                                                                                                                                                     |
| └─[image](./vul/cve-2019-5736/image)                            | cve-2019-5736 exploit via a malicious image                                                                                                                                                     |
| CVE-2019-14271                                                  |                                                                                                                                                                                                 |
| CVE-2019-16884                                                  |                                                                                                                                                                                                 |
| CVE-2020-8555                                                   |                                                                                                                                                                                                 |
| CVE-2020-8558                                                   |                                                                                                                                                                                                 |
| CVE-2020-15157                                                  |                                                                                                                                                                                                 |
| [cve-2020-15257](./vul/cve-2020-15257)                          | abuse the containerd-shim's abstract unix socket <br>in a container with host network namespace                                                                                                 |
| CVE-2021-3493                                                   |                                                                                                                                                                                                 |
| CVE-2021-21285                                                  |                                                                                                                                                                                                 |
| CVE-2021-22555                                                  |                                                                                                                                                                                                 |
| CVE-2021-41091                                                  |                                                                                                                                                                                                 |
| [cve-2021-25741](./vul/cve-2021-25741)                          | kubelet symlink exchange vulnerability allows <br>mounting node filesystem inside a pod                                                                                                         |
| [cve-2022-0492](./vul/cve-2022-0492)                            | escape via cgroup's release agent without <br>CAP_SYS_ADMIN if kernel is vulnerable to <br>CVE-2022-0492                                                                                        |
| CVE-2022-0847                                                   |                                                                                                                                                                                                 |
| [cve-2022-39253](./vul/cve-2022-39253)                          | read host file during docker build via git <br>CVE-2022-39253                                                                                                                                   |
| CVE-2023-28642                                                  |                                                                                                                                                                                                 |
| [cve-2024-0132](./vul/cve-2024-0132)                            | gpu container escape via nvidia-container-toolkit <br>CVE-2024-0132                                                                                                                             |
| CVE-2024-21626                                                  |                                                                                                                                                                                                 |
| [cve-2024-23650](./vul/cve-2024-23650)                          | dos buildkit via oci exporter by sending a crafted <br>request                                                                                                                                  |
| [cve-2025-23266](./vul/cve-2025-23266)                          | gpu container escape via nvidia-container-toolkit <br>cve-2025-23266 by running a malicious container image                                                                                     |
| cve-2025-23267                                                  |                                                                                                                                                                                                 |
| cve-2025-23359                                                  |                                                                                                                                                                                                 |
| cve-2025-31133                                                  |                                                                                                                                                                                                 |
| [cve-2025-47290](./vul/cve-2025-47290)                          | modify host file via containerd cve-2025-47290 <br>during pulling image                                                                                                                         |
| cve-2025-52565                                                  |                                                                                                                                                                                                 |
| caps                                                            | abuse dangerous capabilities in container                                                                                                                                                       |
| └─[shocker](./vul/caps/shocker)                                 | escape by CAP_DAC_READ_SEARCH, alias shocker, <br>found by Sebastian Krahmer (stealth) in 2014                                                                                                  |
| └─[sys_admin](./vul/caps/sys_admin)                             | abuse cap_sys_admin                                                                                                                                                                             |
| &emsp;└─[release_agent](./vul/caps/sys_admin/release_agent)     | escape by cap_sys_admin via cgroups v1 <br>release_agent                                                                                                                                        |
| &emsp;└─mount-device                                            |                                                                                                                                                                                                 |
| &emsp;└─mount-proc                                              |                                                                                                                                                                                                 |
| &emsp;└─device.allow                                            |                                                                                                                                                                                                 |
| &emsp;└─[ebpf](./vul/caps/sys_admin/ebpf)                       | escape by loading evil eBPF programs into the <br>kernel                                                                                                                                        |
| &emsp;&emsp;└─[ebpf-bash](./vul/caps/sys_admin/ebpf/bash)       | abuse eBPF to inject malicious commands into <br>bash processes running on host                                                                                                                 |
| &emsp;&emsp;└─[ebpf-cron](./vul/caps/sys_admin/ebpf/cron)       | abuse eBPF to inject malicious job into host's <br>crontab                                                                                                                                      |
| &emsp;&emsp;└─[ebpf-execve](./vul/caps/sys_admin/ebpf/execve)   | abuse eBPF to hijack execve syscall to run <br>arbitrary commands                                                                                                                               |
| &emsp;&emsp;└─[ebpf-kubelet](./vul/caps/sys_admin/ebpf/kubelet) | abuse eBPF to leak services account token from <br>kubelet                                                                                                                                      |
| &emsp;&emsp;└─ebpf-sshd                                         |                                                                                                                                                                                                 |
| └─[bpf](./vul/caps/bpf)                                         | load evil bpf programs via cap_bpf                                                                                                                                                              |
| &emsp;└─[ebpf](./vul/caps/sys_admin/ebpf)                       | same as caps/sys_admin/ebpf                                                                                                                                                                     |
| └─[sys_ptrace](./vul/caps/sys_ptrace)                           | abuse cap_sys_ptrace                                                                                                                                                                            |
| &emsp;└─[pid_host](./vul/caps/sys_ptrace/pid_host)              | ptrace host processes in a container with <br>cap_sys_ptrace and host pid namespace                                                                                                             |
| └─sys_module                                                    |                                                                                                                                                                                                 |
| └─net_admin                                                     |                                                                                                                                                                                                 |
| [naked](./vul/naked)                                            | we call containers running without seccomp, <br>AppArmor, or SELinux enabled 'naked containers', <br>which leaves them highly vulnerable to <br>kernel exploits and potential container escapes |
| namespace                                                       | shared host namespaces break the isolations                                                                                                                                                     |
| └─[net](./vul/namespace/net)                                    | shared host network namespace breaks the network <br>isolation                                                                                                                                  |
| &emsp;└─shijack                                                 |                                                                                                                                                                                                 |
| &emsp;&emsp;└─basic                                             |                                                                                                                                                                                                 |
| &emsp;&emsp;└─ali                                               |                                                                                                                                                                                                 |
| &emsp;&emsp;└─hw                                                |                                                                                                                                                                                                 |
| &emsp;&emsp;└─gcp                                               |                                                                                                                                                                                                 |
| &emsp;&emsp;└─aws                                               |                                                                                                                                                                                                 |
| └─[pid](./vul/namespace/pid)                                    | shared host pid namespace breaks the process <br>isolation                                                                                                                                      |
| &emsp;└─[proc_root](./vul/namespace/pid/proc_root)              | escape by abusing host pid ns via /proc/[pid]/root                                                                                                                                              |
| fs                                                              |                                                                                                                                                                                                 |
| └─proc-rw                                                       |                                                                                                                                                                                                 |
| &emsp;└─core_pattern                                            |                                                                                                                                                                                                 |
| &emsp;└─binfmt                                                  |                                                                                                                                                                                                 |
| └─sys-rw                                                        |                                                                                                                                                                                                 |
| └─lxcfs-rw                                                      |                                                                                                                                                                                                 |
| shared-socket                                                   | abuse runtime's api via shared socket                                                                                                                                                           |
| └─[docker.sock](./vul/shared-socket/docker-sock)                | escape by shared docker.sock via running a <br>privileged container                                                                                                                             |
| └─containerd.sock                                               |                                                                                                                                                                                                 |
| exposed-api                                                     |                                                                                                                                                                                                 |
| └─docker-2375                                                   |                                                                                                                                                                                                 |
| lxcfs                                                           |                                                                                                                                                                                                 |

### helper

| helper                                  | alias                    | description                    |
|-----------------------------------------|--------------------------|--------------------------------|
| [cve-2021-3493](./helper/cve-2021-3493) | ubuntu-overlayfs-pe,3493 | Ubuntu OverlayFS Local Privesc |
| crash                                   |

### checksec

Just execute `ctrsploit checksec` or standalone binary file `checksec`.
