# ctrsploit: A penetration toolkit for container environment

[中文文档](./README-ZH.md)

ctrsploit [kənˈteɪnər splɔɪt] , follows [sploit-spec](https://github.com/ctrsploit/sploit-spec) v0.4.3

## Why ctrsploit

see [here](https://github.com/ctrsploit/ctrsploit/discussions/11)

## Pre-Built Release

https://github.com/ctrsploit/ctrsploit/releases

```shell
$ wget -q https://github.com/ctrsploit/ctrsploit/releases/latest/download/ctrsploit_linux_amd64 -O /usr/bin/ctrsploit
$ chmod +x /usr/bin/ctrsploit
$ ctrsploit --help
```

## Build

### Build in Container

```bash
make binary
```

## Usage

### env

```shell
$ ctrsploit env     
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

### vul

```shell
$ ctrsploit vul    
NAME:
   ctrsploit vul - list vulnerabilities

USAGE:
   ctrsploit vul [command options]

COMMANDS:
   cve-2016-8867, 8867, amb                        Ambient Capabilities in the Linux kernel allow local users to gain privileges
   cve-2019-5736, 5736                             escape by overwrite runc executable file via /proc/self/exe
   cve-2020-15257, 15257                           abuse the containerd-shim's abstract unix socket in a container with host network namespace
   cve-2021-25741, 25741, kubelet-subpath-symlink  kubelet symlink exchange vulnerability allows mounting node filesystem inside a pod
   cve-2022-0492, 0492                             escape via cgroup's release agent without CAP_SYS_ADMIN if kernel is vulnerable to CVE-2022-0492
   cve-2022-39253, 39253                           read host file during docker build via git CVE-2022-39253
   cve-2024-0132, 0132                             gpu container escape via nvidia-container-toolkit CVE-2024-0132
   cve-2024-23650, 23650                           dos buildkit via oci exporter by sending a crafted request
   cve-2025-23266, 23266                           gpu container escape via nvidia-container-toolkit cve-2025-23266 by running a malicious container image
   cve-2025-47290, 47290                           modify host file via containerd cve-2025-47290 during pulling image
   naked                                           we call containers running without seccomp, AppArmor, or SELinux enabled 'naked containers', which leaves them highly vulnerable to kernel exploits and potential container escapes
   capability, caps                                abuse dangerous capabilities in container
   namespace, ns                                   host level namespaces break the isolations
   shared-socket, sock                             abuse runtime's api via shared socket
   help, h                                         Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```

* :heavy_check_mark: : Fully Supported
* :o: : Partially Supported
* :bug: : Known Bug
* :x: : Not Supported
* `-` : Not Applicable

| vul                                                         | desc                                                                                                                                                                                            | check              | exploit            |
|-------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------------------|--------------------|
| [cve-2016-8867](./vul/cve-2016-8867)                        | ambient capabilities allow local users to gain <br>privileges                                                                                                                                   | :heavy_check_mark: | :heavy_check_mark: |
| cve-2016-9962                                               |                                                                                                                                                                                                 | :x:                | :x:                |
| CVE-2017-1002101                                            |                                                                                                                                                                                                 | :x:                | :heavy_check_mark: |
| [cve-2019-5736](./vul/cve-2019-5736)                        | escape by overwrite runc executable file via <br>/proc/self/exe                                                                                                                                 | :heavy_check_mark: | -                  |
| └─[exec](./vul/cve-2019-5736/exec)                          | cve-2019-5736 exploit via runc exec process                                                                                                                                                     | :heavy_check_mark: | :heavy_check_mark: |
| └─[image](./vul/cve-2019-5736/image)                        | cve-2019-5736 exploit via a malicious image                                                                                                                                                     | :heavy_check_mark: | :heavy_check_mark: |
| CVE-2019-14271                                              |                                                                                                                                                                                                 | :x:                | :heavy_check_mark: |
| CVE-2019-16884                                              |                                                                                                                                                                                                 | :x:                | :heavy_check_mark: |
| CVE-2020-8555                                               |                                                                                                                                                                                                 | :x:                | :heavy_check_mark: |
| CVE-2020-8558                                               |                                                                                                                                                                                                 | :x:                | :heavy_check_mark: |
| CVE-2020-15157                                              |                                                                                                                                                                                                 | :x:                | :heavy_check_mark: |
| [cve-2020-15257](./vul/cve-2020-15257)                      | abuse the containerd-shim's abstract unix socket <br>in a container with host network namespace                                                                                                 | :heavy_check_mark: | :heavy_check_mark: |
| CVE-2021-3493                                               |                                                                                                                                                                                                 | :x:                | :heavy_check_mark: |
| CVE-2021-21285                                              |                                                                                                                                                                                                 | :x:                | :heavy_check_mark: |
| CVE-2021-22555                                              |                                                                                                                                                                                                 | :x:                | :heavy_check_mark: |
| CVE-2021-41091                                              |                                                                                                                                                                                                 | :x:                | :heavy_check_mark: |
| [cve-2021-25741](./vul/cve-2021-25741)                      | kubelet symlink exchange vulnerability allows <br>mounting node filesystem inside a pod                                                                                                         | :heavy_check_mark: | :heavy_check_mark: |
| [cve-2022-0492](./vul/cve-2022-0492)                        | escape via cgroup's release agent without <br>CAP_SYS_ADMIN if kernel is vulnerable to <br>CVE-2022-0492                                                                                        | :heavy_check_mark: | :heavy_check_mark: |
| CVE-2022-0847                                               |                                                                                                                                                                                                 | :x:                | :heavy_check_mark: |
| [cve-2022-39253](./vul/cve-2022-39253)                      | read host file during docker build via git <br>CVE-2022-39253                                                                                                                                   | :o:                | :heavy_check_mark: |
| CVE-2023-28642                                              |                                                                                                                                                                                                 | :x:                | :x:                |
| [cve-2024-0132](./vul/cve-2024-0132)                        | gpu container escape via nvidia-container-toolkit <br>CVE-2024-0132                                                                                                                             | :o:                | :heavy_check_mark: |
| CVE-2024-21626                                              |                                                                                                                                                                                                 | :x:                | :x:                |
| [cve-2024-23650](./vul/cve-2024-23650)                      | dos buildkit via oci exporter by sending a crafted <br>request                                                                                                                                  | :heavy_check_mark: | :heavy_check_mark: |
| [cve-2025-23266](./vul/cve-2025-23266)                      | gpu container escape via nvidia-container-toolkit <br>cve-2025-23266 by running a malicious container image                                                                                     | :o:                | :heavy_check_mark: |
| cve-2025-23267                                              |                                                                                                                                                                                                 | :x:                | :x:                |
| cve-2025-23359                                              |                                                                                                                                                                                                 | :x:                | :x:                |
| cve-2025-31133                                              |                                                                                                                                                                                                 | :x:                | :x:                |
| [cve-2025-47290](./vul/cve-2025-47290)                      | modify host file via containerd cve-2025-47290 <br>during pulling image                                                                                                                         | :heavy_check_mark: | :heavy_check_mark: |
| cve-2025-52565                                              |                                                                                                                                                                                                 | :x:                | :x:                |
| [caps](./vul/caps)                                          | abuse dangerous capabilities in container                                                                                                                                                       | -                  | -                  |
| └─[shocker](./vul/caps/shocker)                             | escape by CAP_DAC_READ_SEARCH, alias shocker, <br>found by Sebastian Krahmer (stealth) in 2014                                                                                                  | :heavy_check_mark: | :heavy_check_mark: |
| └─[sys_admin](./vul/caps/sys_admin)                         | abuse cap_sys_admin                                                                                                                                                                             | :heavy_check_mark: | -                  |
| &emsp;└─[release_agent](./vul/caps/sys_admin/release_agent) | escape by cap_sys_admin via cgroups v1 <br>release_agent                                                                                                                                        | :heavy_check_mark: | :heavy_check_mark: |
| &emsp;└─mount-device                                        |                                                                                                                                                                                                 | :x:                | :x:                |
| &emsp;└─mount-proc                                          |                                                                                                                                                                                                 | :x:                | :x:                |
| &emsp;└─device.allow                                        |                                                                                                                                                                                                 | :x:                | :x:                |
| &emsp;└─[ebpf](./vul/caps/sys_admin/ebpf)                   | escape by loading evil eBPF programs into the <br>kernel                                                                                                                                        | :heavy_check_mark: | -                  |
| &emsp;&emsp;└─[bash](./vul/caps/sys_admin/ebpf/bash)        | abuse eBPF to inject malicious commands into <br>bash processes running on host                                                                                                                 | :heavy_check_mark: | :heavy_check_mark: |
| &emsp;&emsp;└─[cron](./vul/caps/sys_admin/ebpf/cron)        | abuse eBPF to inject malicious job into host's <br>crontab                                                                                                                                      | :heavy_check_mark: | :heavy_check_mark: |
| &emsp;&emsp;└─[execve](./vul/caps/sys_admin/ebpf/execve)    | abuse eBPF to hijack execve syscall to run <br>arbitrary commands                                                                                                                               | :heavy_check_mark: | :heavy_check_mark: |
| &emsp;&emsp;└─[kubelet](./vul/caps/sys_admin/ebpf/kubelet)  | abuse eBPF to leak services account token from <br>kubelet                                                                                                                                      | :heavy_check_mark: | :heavy_check_mark: |
| &emsp;&emsp;└─sshd                                          |                                                                                                                                                                                                 | :x:                | :x:                |
| └─[bpf](./vul/caps/bpf)                                     | load evil bpf programs via cap_bpf                                                                                                                                                              | -                  | -                  |
| &emsp;└─[ebpf](./vul/caps/sys_admin/ebpf)                   | same as caps/sys_admin/ebpf                                                                                                                                                                     | :heavy_check_mark: | -                  |
| └─[sys_ptrace](./vul/caps/sys_ptrace)                       | abuse cap_sys_ptrace                                                                                                                                                                            | :heavy_check_mark: | -                  |
| &emsp;└─[pid_host](./vul/caps/sys_ptrace/pid_host)          | ptrace host processes in a container with <br>cap_sys_ptrace and host pid namespace                                                                                                             | :heavy_check_mark: | :heavy_check_mark: |
| └─sys_module                                                |                                                                                                                                                                                                 | :x:                | :x:                |
| └─net_admin                                                 |                                                                                                                                                                                                 | :x:                | :x:                |
| [naked](./vul/naked)                                        | we call containers running without seccomp, <br>AppArmor, or SELinux enabled 'naked containers', <br>which leaves them highly vulnerable to <br>kernel exploits and potential container escapes | :heavy_check_mark: | -                  |
| [namespace](./vul/namespace)                                | shared host namespaces break the isolations                                                                                                                                                     | -                  | -                  |
| └─[net](./vul/namespace/net)                                | shared host network namespace breaks the network <br>isolation                                                                                                                                  | :heavy_check_mark: | :x:                |
| &emsp;└─shijack                                             |                                                                                                                                                                                                 | :x:                | :x:                |
| &emsp;&emsp;└─basic                                         |                                                                                                                                                                                                 | :x:                | :x:                |
| &emsp;&emsp;└─ali                                           |                                                                                                                                                                                                 | :x:                | :x:                |
| &emsp;&emsp;└─hw                                            |                                                                                                                                                                                                 | :x:                | :x:                |
| &emsp;&emsp;└─gcp                                           |                                                                                                                                                                                                 | :x:                | :x:                |
| &emsp;&emsp;└─aws                                           |                                                                                                                                                                                                 | :x:                | :x:                |
| └─[pid](./vul/namespace/pid)                                | shared host pid namespace breaks the process <br>isolation                                                                                                                                      | -                  | -                  |
| &emsp;└─[proc_root](./vul/namespace/pid/proc_root)          | escape by abusing host pid ns via /proc/[pid]/root                                                                                                                                              | :heavy_check_mark: | :heavy_check_mark: |
| fs                                                          |                                                                                                                                                                                                 | :x:                | :x:                |
| └─proc-rw                                                   |                                                                                                                                                                                                 | :x:                | -                  |
| &emsp;└─core_pattern                                        |                                                                                                                                                                                                 | :x:                | :x:                |
| &emsp;└─binfmt                                              |                                                                                                                                                                                                 | :x:                | :x:                |
| └─sys-rw                                                    |                                                                                                                                                                                                 | :x:                | :x:                |
| └─lxcfs-rw                                                  |                                                                                                                                                                                                 | :x:                | :x:                |
| [shared-socket](./vul/shared-socket)                        | abuse runtime's api via shared socket                                                                                                                                                           | -                  | -                  |
| └─[docker.sock](./vul/shared-socket/docker-sock)            | escape by shared docker.sock via running a <br>privileged container                                                                                                                             | :heavy_check_mark: | :heavy_check_mark: |
| └─containerd.sock                                           |                                                                                                                                                                                                 | :x:                | :x:                |
| exposed-api                                                 |                                                                                                                                                                                                 | -                  | -                  |
| └─docker-2375                                               |                                                                                                                                                                                                 | :x:                | :x:                |
| lxcfs                                                       |                                                                                                                                                                                                 | :x:                | :x:                |

### exploit

```shell
$ ctrsploit exploit                                       
NAME:
   ctrsploit exploit - run a exploit

USAGE:
   ctrsploit exploit [command options]

COMMANDS:
   cve-2016-8867, 8867, amb                           Ambient Capabilities in the Linux kernel allow local users to gain privileges
   cve-2019-5736, 5736                                escape by overwrite runc executable file via /proc/self/exe
   cve-2020-15257, 15257                              abuse the containerd-shim's abstract unix socket in a container with host network namespace
   cve-2021-25741, 25741, kubelet-subpath-symlink     kubelet symlink exchange vulnerability allows mounting node filesystem inside a pod
   cve-2022-0492, 0492                                escape via cgroup's release agent without CAP_SYS_ADMIN if kernel is vulnerable to CVE-2022-0492
   cve-2022-39253, 39253                              read host file during docker build via git CVE-2022-39253
   cve-2024-0132, 0132                                gpu container escape via nvidia-container-toolkit CVE-2024-0132
   cve-2024-23650, 23650                              dos buildkit via oci exporter by sending a crafted request
   cve-2025-23266, 23266                              gpu container escape via nvidia-container-toolkit cve-2025-23266 by running a malicious container image
   cve-2025-47290, 47290                              modify host file via containerd cve-2025-47290 during pulling image
   shocker, cap_dac_read_search, open_by_handle_at    escape by CAP_DAC_READ_SEARCH, alias shocker, found by Sebastian Krahmer (stealth) in 2014
   cap_sys_admin, sys_admin                           abuse cap_sys_admin
   release_agent, ra                                  escape by cap_sys_admin via cgroups v1 release_agent
   ebpf                                               escape by loading evil eBPF programs into the kernel
   ebpf-bash, bash                                    abuse eBPF to inject malicious commands into bash processes running on host
   ebpf-execve, execve                                abuse eBPF to hijack execve syscall to run arbitrary commands
   ebpf-cron, cron                                    abuse eBPF to inject malicious job into host's crontab
   ebpf-kubelet, kubelet                              abuse eBPF to leak services account token from kubelet
   cap_bpf, bpf                                       load evil bpf programs via cap_bpf
   cap_sys_ptrace, sys_ptrace, ptrace                 abuse cap_sys_ptrace
   ptrace-pid-host, ptrace-pid                        ptrace host processes in a container with cap_sys_ptrace and host pid namespace
   host-pid, pid                                      shared host pid namespace breaks process isolation
   host-pid-proc-root, proc                           escape by abusing host pid ns via /proc/[pid]/root
   docker.sock, docker                                escape by shared docker.sock via running a privileged container
   CVE-2021-22555, 22555                              escape tech by using the CVE-2021-22555
   CVE-2020-8555, 8555                                k8s CVE-2020-8555 SSRF
   CVE-2017-1002101, subPath1, 1002101, 2017-1002101  CVE-2017-1002101
   dirty-pipe, dp, CVE-2022-0847, 0847                dirty-pipe
   crash, c                                           make container crash
   help, h                                            Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```

### checksec

```shell
$ ctrsploit checksec     
NAME:
   ctrsploit checksec - check security inside a container

USAGE:
   ctrsploit checksec [command options]

COMMANDS:
   auto, a                                          auto check security
   env, e                                           gather information
   cve-2016-8867, 8867, amb                         Ambient Capabilities in the Linux kernel allow local users to gain privileges
   cve-2019-5736, 5736                              
   cve-2020-15257, 15257                            Abuse the containerd-shim's abstract unix socket when running in a container with host network namespace.
   cve-2021-25741, 25741, kubelet-subpath-symlink   Kubernetes kubelet symlink exchange vulnerability allows mounting Node filesystem inside POD with read-write privileges
   cve-2022-0492, 0492                              Container escape using cgroup's release agent without CAP_SYS_ADMIN if kernel has CVE-2022-0492
   cve-2022-39253, 39253                            docker build host file read by git CVE-2022-39253
   cve-2024-0132, 0132                              nvidia-container-toolkit CVE-2024-0132 container escape. Affected versions: libnvidia-container >= 1.0.0, <= 1.16.1
   cve-2024-23650, 23650                            BuildKit OCI exporter DoS vulnerability by sending a crafted request.
   cve-2025-23266, 23266                            NVIDIA Container Toolkit allows an attacker to execute arbitrary code on the host by running a specially crafted container image.
   cve-2025-47290, 47290                            TOCTOU vulnerability in containerd that allows modification of the host filesystem during image pull.
   shocker, cap_dac_read_search, open_by_handle_at  Container escape with CAP_DAC_READ_SEARCH, alias shocker, found by Sebastian Krahmer (stealth) in 2014.
   cap_sys_admin, sys_admin                         Container can be escaped when has cap_sys_admin
   cap_bpf, bpf                                     Container can load evil bpf program when has cap_bpf, may cause container escape
   cap_sys_ptrace, sys_ptrace, ptrace               Container can be escaped when has cap_sys_ptrace
   ptrace-pid-host, ptrace-pid                      Container can be escaped when has cap_sys_ptrace and host pid namespace
   naked                                            We call containers running without seccomp, AppArmor, or SELinux enabled 'naked containers', which leaves them highly vulnerable to kernel exploits and potential container escapes
   host-net, net                                    The network namespace of the host is shared
   host-pid, pid                                    container can be escaped with host pid namespace
   docker.sock, docker                              escape by shared docker socket
   help, h                                          Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```

```shell
$ ctrsploit --colorful checksec auto
✘  cve-2025-47290       # TOCTOU vulnerability in containerd that allows modification of the host filesystem during image pull.
✔  naked        # We call containers running without seccomp, AppArmor, or SELinux enabled 'naked containers', which leaves them highly vulnerable to kernel exploits and potential container escapes
✘  cve-2019-5736        
✘  cve-2021-25741       # Kubernetes kubelet symlink exchange vulnerability allows mounting Node filesystem inside POD with read-write privileges
✘  cve-2022-0492        # Container escape using cgroup's release agent without CAP_SYS_ADMIN if kernel has CVE-2022-0492
✔  cap_sys_admin        # Container can be escaped when has cap_sys_admin
✔  cap_bpf      # Container can load evil bpf program when has cap_bpf, may cause container escape
✔  host-net     # The network namespace of the host is shared
✘  cve-2016-8867        # Ambient Capabilities in the Linux kernel allow local users to gain privileges
✘  cve-2022-39253       # docker build host file read by git CVE-2022-39253
✘  cve-2024-23650       # BuildKit OCI exporter DoS vulnerability by sending a crafted request.
✘  cve-2025-23266       # NVIDIA Container Toolkit allows an attacker to execute arbitrary code on the host by running a specially crafted container image.
✔  shocker      # Container escape with CAP_DAC_READ_SEARCH, alias shocker, found by Sebastian Krahmer (stealth) in 2014.
✔  cap_sys_ptrace       # Container can be escaped when has cap_sys_ptrace
✘  docker.sock  # escape by shared docker socket
✘  cve-2020-15257       # Abuse the containerd-shim's abstract unix socket when running in a container with host network namespace.
✘  cve-2024-0132        # nvidia-container-toolkit CVE-2024-0132 container escape. Affected versions: libnvidia-container >= 1.0.0, <= 1.16.1
✔  ptrace-pid-host      # Container can be escaped when has cap_sys_ptrace and host pid namespace
✔  host-pid     # container can be escaped with host pid namespac
```

### helper

// TODO

## Progress of Development

### vul

* :heavy_check_mark: : Fully Supported
* :o: : Partially Supported
* :bug: : Known Bug
* :x: : Not Supported
* `-` : Not Applicable

| vul                                                         | check              | exploit            | test                | doc                | video              | case               |
|-------------------------------------------------------------|--------------------|--------------------|---------------------|--------------------|--------------------|--------------------|
| [cve-2016-8867](./vul/cve-2016-8867)                        | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark:  | :heavy_check_mark: | :heavy_check_mark: | :x:                |
| cve-2016-9962                                               | :x:                | :x:                | :x:                 | :x:                | :x:                | :x:                |
| CVE-2017-1002101                                            | :x:                | :heavy_check_mark: | :x:                 | :x:                | :x:                | :x:                |
| [cve-2019-5736](./vul/cve-2019-5736)                        | :heavy_check_mark: | -                  | -                   | :heavy_check_mark: | -                  | :x:                |
| └─[exec](./vul/cve-2019-5736/exec)                          | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark:  | -                  | :heavy_check_mark: | :x:                |
| └─[image](./vul/cve-2019-5736/image)                        | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark:  | -                  | :heavy_check_mark: | :x:                |
| CVE-2019-14271                                              | :x:                | :heavy_check_mark: | :x:                 | :x:                | :x:                | :x:                |
| CVE-2019-16884                                              | :x:                | :heavy_check_mark: | :x:                 | :x:                | :x:                | :x:                |
| CVE-2020-8555                                               | :x:                | :heavy_check_mark: | :x:                 | :x:                | :x:                | :x:                |
| CVE-2020-8558                                               | :x:                | :heavy_check_mark: | :x:                 | :x:                | :x:                | :x:                |
| CVE-2020-15157                                              | :x:                | :heavy_check_mark: | :x:                 | :x:                | :x:                | :x:                |
| [cve-2020-15257](./vul/cve-2020-15257)                      | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark:  | :heavy_check_mark: | :heavy_check_mark: | :x:                |
| CVE-2021-3493                                               | :x:                | :heavy_check_mark: | :x:                 | :x:                | :x:                | :x:                |
| CVE-2021-21285                                              | :x:                | :heavy_check_mark: | :x:                 | :x:                | :x:                | :x:                |
| CVE-2021-22555                                              | :x:                | :heavy_check_mark: | :x:                 | :x:                | :x:                | :x:                |
| CVE-2021-41091                                              | :x:                | :heavy_check_mark: | :x:                 | :x:                | :x:                | :x:                |
| [cve-2021-25741](./vul/cve-2021-25741)                      | :heavy_check_mark: | :heavy_check_mark: | :x:                 | :heavy_check_mark: | :x:                | :x:                |
| [cve-2022-0492](./vul/cve-2022-0492)                        | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark:  | :heavy_check_mark: | :heavy_check_mark: | :x:                |
| CVE-2022-0847                                               | :x:                | :heavy_check_mark: | :x:                 | :x:                | :x:                | :x:                |
| [cve-2022-39253](./vul/cve-2022-39253)                      | :o:                | :heavy_check_mark: | :heavy_check_mark:  | :heavy_check_mark: | :heavy_check_mark: | :x:                |
| CVE-2023-28642                                              | :x:                | :x:                | :x:                 | :x:                | :x:                | :x:                |
| [cve-2024-0132](./vul/cve-2024-0132)                        | :o:                | :heavy_check_mark: | :heavy_check_mark:  | :heavy_check_mark: | :heavy_check_mark: | :x:                |
| CVE-2024-21626                                              | :x:                | :x:                | :x:                 | :x:                | :x:                | :x:                |
| [cve-2024-23650](./vul/cve-2024-23650)                      | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark:  | :heavy_check_mark: | :heavy_check_mark: | :x:                |
| [cve-2025-23266](./vul/cve-2025-23266)                      | :o:                | :heavy_check_mark: | :heavy_check_mark:  | :heavy_check_mark: | :heavy_check_mark: | :x:                |
| cve-2025-23267                                              | :x:                | :x:                | :x:                 | :x:                | :x:                | :x:                |
| cve-2025-23359                                              | :x:                | :x:                | :x:                 | :x:                | :x:                | :x:                |
| cve-2025-31133                                              | :x:                | :x:                | :x:                 | :x:                | :x:                | :x:                |
| [cve-2025-47290](./vul/cve-2025-47290)                      | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark:  | :heavy_check_mark: | :heavy_check_mark: | :x:                |
| cve-2025-52565                                              | :x:                | :x:                | :x:                 | :x:                | :x:                | :x:                |
| [caps](./vul/caps)                                          | -                  | -                  | -                   | -                  | -                  | -                  |
| └─[shocker](./vul/caps/shocker)                             | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark:️ | :heavy_check_mark: | :heavy_check_mark: | :x:                |
| └─[sys_admin](./vul/caps/sys_admin)                         | :heavy_check_mark: | -                  | -                   | -                  | -                  | -                  |
| &emsp;└─[release_agent](./vul/caps/sys_admin/release_agent) | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark:  | :heavy_check_mark: | :heavy_check_mark: | :x:                |
| &emsp;└─mount-device                                        | :x:                | :x:                | :x:                 | :x:                | :x:                | :x:                |
| &emsp;└─mount-proc                                          | :x:                | :x:                | :x:                 | :x:                | :x:                | :x:                |
| &emsp;└─device.allow                                        | :x:                | :x:                | :x:                 | :x:                | :x:                | :x:                |
| &emsp;└─[ebpf](./vul/caps/sys_admin/ebpf)                   | :heavy_check_mark: | -                  | -                   | -                  | -                  | -                  |
| &emsp;&emsp;└─[bash](./vul/caps/sys_admin/ebpf/bash)        | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark:  | :heavy_check_mark: | :heavy_check_mark: | :x:                |
| &emsp;&emsp;└─[cron](./vul/caps/sys_admin/ebpf/cron)        | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark:  | :heavy_check_mark: | :heavy_check_mark: | :x:                |
| &emsp;&emsp;└─[execve](./vul/caps/sys_admin/ebpf/execve)    | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark:  | :heavy_check_mark: | :heavy_check_mark: | :x:                |
| &emsp;&emsp;└─[kubelet](./vul/caps/sys_admin/ebpf/kubelet)  | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark:  | :heavy_check_mark: | :heavy_check_mark: | :x:                |
| &emsp;&emsp;└─sshd                                          | :x:                | :x:                | :x:                 | :x:                | :x:                | :x:                |
| └─[bpf](./vul/caps/bpf)                                     | -                  | -                  | -                   | -                  | -                  | -                  |
| &emsp;└─[ebpf](./vul/caps/sys_admin/ebpf)                   | :heavy_check_mark: | -                  | -                   | -                  | -                  | -                  |
| └─[sys_ptrace](./vul/caps/sys_ptrace)                       | :heavy_check_mark: | -                  | -                   | -                  | -                  | -                  |
| &emsp;└─[pid_host](./vul/caps/sys_ptrace/pid_host)          | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark:  | :heavy_check_mark: | :heavy_check_mark: | :x:                |
| └─sys_module                                                | :x:                | :x:                | :x:                 | :x:                | :x:                | :x:                |
| └─net_admin                                                 | :x:                | :x:                | :x:                 | :x:                | :x:                | :x:                |
| [naked](./vul/naked)                                        | :heavy_check_mark: | -                  | :heavy_check_mark:  | :x:                | :x:                | :x:                |
| [namespace](./vul/namespace)                                | -                  | -                  | -                   | -                  | -                  | -                  |
| └─[net](./vul/namespace/net)                                | :heavy_check_mark: | :x:                | :x:                 | :x:                | :x:                | :x:                |
| &emsp;└─shijack                                             | :x:                | :x:                | :x:                 | :x:                | :x:                | :x:                |
| &emsp;&emsp;└─basic                                         | :x:                | :x:                | :x:                 | :x:                | :x:                | :x:                |
| &emsp;&emsp;└─ali                                           | :x:                | :x:                | :x:                 | :x:                | :x:                | :x:                |
| &emsp;&emsp;└─hw                                            | :x:                | :x:                | :x:                 | :x:                | :x:                | :x:                |
| &emsp;&emsp;└─gcp                                           | :x:                | :x:                | :x:                 | :x:                | :x:                | :x:                |
| &emsp;&emsp;└─aws                                           | :x:                | :x:                | :x:                 | :x:                | :x:                | :x:                |
| └─[pid](./vul/namespace/pid)                                | -                  | -                  | -                   | -                  | -                  | -                  |
| &emsp;└─[proc_root](./vul/namespace/pid/proc_root)          | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark:  | :heavy_check_mark: | :heavy_check_mark: | :x:                |
| fs                                                          | :x:                | :x:                | :x:                 | :x:                | :x:                | :x:                |
| └─proc-rw                                                   | :x:                | -                  | -                   | :x:                | :x:                | :x:                |
| &emsp;└─core_pattern                                        | :x:                | :x:                | :x:                 | :x:                | :x:                | :x:                |
| &emsp;└─binfmt                                              | :x:                | :x:                | :x:                 | :x:                | :x:                | :x:                |
| └─sys-rw                                                    | :x:                | :x:                | :x:                 | :x:                | :x:                | :x:                |
| └─lxcfs-rw                                                  | :x:                | :x:                | :x:                 | :x:                | :x:                | :x:                |
| [shared-socket](./vul/shared-socket)                        | -                  | -                  | -                   | -                  | -                  | -                  |
| └─[docker.sock](./vul/shared-socket/docker-sock)            | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark:  | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: |
| └─containerd.sock                                           | :x:                | :x:                | :x:                 | :x:                | :x:                | :x:                |
| exposed-api                                                 | -                  | -                  | -                   | -                  | -                  | -                  |
| └─docker-2375                                               | :x:                | :x:                | :x:                 | :x:                | :x:                | :x:                |
| lxcfs                                                       | :x:                | :x:                | :x:                 | :x:                | :x:                | :x:                |
