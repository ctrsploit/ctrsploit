# ctrsploit: A penetration toolkit for container environment

ctrsploit [kənˈteɪnər splɔɪt] , follows [sploit-spec](https://github.com/ctrsploit/sploit-spec)

![](./docs/images/logo-white-256.png)

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
   ctrsploit env [command [command options]]

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
   sysctl              display sysctl information
   rlimit              get process resource limits
   namespace, n, ns    check namespace is host ns
   no-new-privs, nnp   show NoNewPrivs status for the current process
   suid, setuid        find and list SUID files
   docker-version, dv  guess dockerd version range
   services, svc       discover K8s cluster services and ports via env vars and DNS
   upload, up          upload <servicename> <filename> <obs> [host]

OPTIONS:
   --help, -h  show help
```

#### env no-new-privs

Check whether the current process has `NoNewPrivs` enabled. When enabled, SUID
programs and file capabilities cannot grant new privileges across `execve`.

```shell
$ ctrsploit env no-new-privs
$ ctrsploit env nnp
```

#### env suid

List SUID files visible to the current process. By default, the command scans
common executable directories; use `--all` to scan from `/`, or `--path` for
explicit paths.

```shell
$ ctrsploit env suid
$ ctrsploit env suid --all
$ ctrsploit env suid --path /bin,/usr/bin
```

#### env services

Discover K8s cluster services and their ports without API server access, using multiple layered techniques:

| Method | Technique | Speed | Scope |
|--------|-----------|-------|-------|
| `env` | Parse K8s-injected env vars (`*_SERVICE_HOST`, `*_PORT_*_TCP`, …) | Instant | Same namespace only |
| `wildcard` | CoreDNS wildcard SRV (`any.any.svc.<zone>`) | Fast | All namespaces |
| `axfr` | DNS zone transfer from `ns.dns.<zone>:53` | Fast | All namespaces (if enabled) |
| `cidr` | PTR scan service CIDR + SRV enrichment | Slow (65536 IPs for /16) | All namespaces |

```shell
# Discover all services (default: all methods)
$ ctrsploit env services

# Quick scan - env vars only (zero network, milliseconds)
$ ctrsploit env services --methods env

# Skip slow CIDR scan
$ ctrsploit env services --methods env,wildcard,axfr

# Custom CIDR with more threads
$ ctrsploit env services --cidr 10.96.0.0/24 --threads 64

# Export results as NDJSON
$ ctrsploit env services --output /tmp/services.json

# Custom DNS zone
$ ctrsploit env services --zone mycluster.local
```

### vul

```shell
$ ctrsploit vul
NAME:
   ctrsploit vul - list vulnerabilities

USAGE:
   ctrsploit vul [command [command options]]

COMMANDS:
   cve-2016-8867, 8867, amb                        Ambient Capabilities in the Linux kernel allow local users to gain privileges
   cve-2019-5736, 5736                             escape by overwrite runc executable file via /proc/self/exe
   cve-2020-8558, 8558                             access services bound to 127.0.0.1 from adjacent hosts
   cve-2020-15257, 15257                           abuse the containerd-shim's abstract unix socket in a container with host network namespace
   cve-2021-25741, 25741, kubelet-subpath-symlink  kubelet symlink exchange vulnerability allows mounting node filesystem inside a pod
   cve-2021-25748, 25748, ingress-nginx-path-leak  ingress-nginx path validation bypass vulnerability allows credential leakage through newline injection
   cve-2021-3493, 3493, ubuntu-overlayfs-pe, CVE-2021-3493  local privilege escalation in Ubuntu OverlayFS that may lead to container escape when the kernel attack surface is exposed
   cve-2022-0492, 0492                             escape via cgroup's release agent without CAP_SYS_ADMIN if kernel is vulnerable to CVE-2022-0492
   cve-2022-0847, 0847, dirty-pipe, dirtypipe, dp, CVE-2022-0847  local privilege escalation and container escape in Linux kernel Dirty Pipe
   cve-2022-39253, 39253                           read host file during docker build via git CVE-2022-39253
   cve-2024-0132, 0132                             gpu container escape via nvidia-container-toolkit CVE-2024-0132
   cve-2024-23650, 23650                           dos buildkit via oci exporter by sending a crafted request
   cve-2024-40635, 40635                           bypass runAsNonRoot via integer overflow in User ID handling in containerd
   cve-2025-23266, 23266                           gpu container escape via nvidia-container-toolkit cve-2025-23266 by running a malicious container image
   cve-2025-47290, 47290                           modify host file via containerd cve-2025-47290 during pulling image
   cve-2025-62725, 62725                           path traversal in Docker Compose OCI artifacts allows arbitrary file write via malicious registry
    cve-2026-31431, 31431, copy-fail                local privilege escalation and container escape in Linux kernel AF_ALG AEAD interface
    cve-2026-43284, 43284, dirty-frag, dirtyfrag    local privilege escalation and container escape in Linux kernel xfrm ESP Dirty Frag path
    cve-2026-43500, 43500, dirty-frag-rxrpc, dirtyfrag-rxrpc  local privilege escalation in Linux kernel RxRPC/rxkad Dirty Frag path
    fork-bomb                                       
   naked                                           we call containers running without seccomp, AppArmor, or SELinux enabled 'naked containers', which leaves them highly vulnerable to kernel exploits and potential container escapes
   capability, caps                                abuse dangerous capabilities in container
   namespace, ns                                   host level namespaces break the isolations
   service-account-token, sa-token, token          check service account token related vulnerabilities
   shared-socket, sock                             abuse runtime's api via shared socket

OPTIONS:
   --help, -h  show help
```

### module

Group vulnerabilities by component or configuration type, and use the
tables below to see their check/exploit support status.

```shell
$ ctrsploit module
NAME:
   ctrsploit module - group vulnerabilities by component or config type

USAGE:
   ctrsploit module [component|config] [name]

DESCRIPTION:
   Classify and operate vulnerabilities by logical module
   such as kernel, runc, containerd, or config (e.g. capability).

COMMANDS:
   config, cfg                            insecure configuration and misconfiguration issues
   runc, r                                runc related vulnerabilities
   containerd, c                          containerd related vulnerabilities
   docker, d                              docker related vulnerabilities
   nvidia-container-toolkit, nvidia, nct  nvidia-container-toolkit related vulnerabilities
   docker-compose, compose                docker-compose related vulnerabilities
   buildkit, bk                           buildkit related vulnerabilities
   kubernetes, k8s                        kubernetes related vulnerabilities
   ingress-nginx, ingress                 ingress-nginx related vulnerabilities
   git, g                                 git related vulnerabilities
   kernel, k                              kernel related vulnerabilities

OPTIONS:
   --help, -h  show help
```

* :heavy_check_mark: : Fully Supported
* :o: : Partially Supported
* :bug: : Known Bug
* :x: : Not Supported
* `-` : Not Applicable

#### config

| vul | desc | check | exploit |
|-----|------|-------|---------|
| [fork-bomb](./vul/fork-bomb) | fork bomb causes denial of service when resource limits or cgroup configs are unsafe | :heavy_check_mark: | :heavy_check_mark: |
| [caps](./vul/caps) | abuse dangerous capabilities in container | - | - |
| └─[shocker](./vul/caps/shocker) | escape by CAP_DAC_READ_SEARCH, alias shocker, found by Sebastian Krahmer (stealth) in 2014 | :heavy_check_mark: | :heavy_check_mark: |
| └─[sys_admin](./vul/caps/sys_admin) | abuse cap_sys_admin | :heavy_check_mark: | - |
| &emsp;└─[release_agent](./vul/caps/sys_admin/release_agent) | escape by cap_sys_admin via cgroups v1 release_agent | :heavy_check_mark: | :heavy_check_mark: |
| &emsp;└─mount-device |  | :x: | :x: |
| &emsp;└─mount-proc |  | :x: | :x: |
| &emsp;└─device.allow |  | :x: | :x: |
| &emsp;└─[ebpf](./vul/caps/sys_admin/ebpf) | escape by loading evil eBPF programs into the kernel | :heavy_check_mark: | - |
| &emsp;&emsp;└─[bash](./vul/caps/sys_admin/ebpf/bash) | abuse eBPF to inject malicious commands into bash processes running on host | :heavy_check_mark: | :heavy_check_mark: |
| &emsp;&emsp;└─[cron](./vul/caps/sys_admin/ebpf/cron) | abuse eBPF to inject malicious job into host's crontab | :heavy_check_mark: | :heavy_check_mark: |
| &emsp;&emsp;└─[execve](./vul/caps/sys_admin/ebpf/execve) | abuse eBPF to hijack execve syscall to run arbitrary commands | :heavy_check_mark: | :heavy_check_mark: |
| &emsp;&emsp;└─[kubelet](./vul/caps/sys_admin/ebpf/kubelet) | abuse eBPF to leak services account token from kubelet | :heavy_check_mark: | :heavy_check_mark: |
| &emsp;&emsp;└─sshd |  | :x: | :x: |
| └─[bpf](./vul/caps/bpf) | load evil bpf programs via cap_bpf | - | - |
| &emsp;└─[ebpf](./vul/caps/sys_admin/ebpf) | same as caps/sys_admin/ebpf | :heavy_check_mark: | - |
| └─[sys_ptrace](./vul/caps/sys_ptrace) | abuse cap_sys_ptrace | :heavy_check_mark: | - |
| &emsp;└─[pid_host](./vul/caps/sys_ptrace/pid_host) | ptrace host processes in a container with cap_sys_ptrace and host pid namespace | :heavy_check_mark: | :heavy_check_mark: |
| └─sys_module |  | :x: | :x: |
| └─net_admin |  | :x: | :x: |
| └─[cve-2016-8867](./vul/cve-2016-8867) | ambient capabilities allow local users to gain privileges | :heavy_check_mark: | :heavy_check_mark: |
| [naked](./vul/naked) | we call containers running without seccomp, AppArmor, or SELinux enabled 'naked containers', which leaves them highly vulnerable to kernel exploits and potential container escapes | :heavy_check_mark: | - |
| [shared-socket](./vul/shared-socket) | abuse runtime's api via shared socket | - | - |
| └─[docker.sock](./vul/shared-socket/docker-sock) | escape by shared docker.sock via running a privileged container | :heavy_check_mark: | :heavy_check_mark: |
| └─containerd.sock |  | :x: | :x: |
| [sa-token](./vul/sa-token) |  | - | - |
| └─[secret](./vul/sa-token/access-secrets) | check if service account token can access Kubernetes Secrets | :heavy_check_mark: | - |
| └─[policy](./vul/sa-token/policy) | check if service account token has dangerous permissions | :heavy_check_mark: | - |
| [namespace](./vul/namespace) | shared host namespaces break the isolations | - | - |
| └─[net](./vul/namespace/net) | shared host network namespace breaks the network isolation | :heavy_check_mark: | :x: |
| &emsp;└─shijack |  | :x: | :x: |
| &emsp;&emsp;└─basic |  | :x: | :x: |
| &emsp;&emsp;└─ali |  | :x: | :x: |
| &emsp;&emsp;└─hw |  | :x: | :x: |
| &emsp;&emsp;└─gcp |  | :x: | :x: |
| &emsp;&emsp;└─aws |  | :x: | :x: |
| └─[pid](./vul/namespace/pid) | shared host pid namespace breaks the process isolation | - | - |
| &emsp;└─[proc_root](./vul/namespace/pid/proc_root) | escape by abusing host pid ns via /proc/[pid]/root | :heavy_check_mark: | :heavy_check_mark: |
| fs |  | :x: | :x: |
| └─proc-rw |  | :x: | - |
| &emsp;└─core_pattern |  | :x: | :x: |
| &emsp;└─binfmt |  | :x: | :x: |
| └─sys-rw |  | :x: | :x: |
| └─lxcfs-rw |  | :x: | :x: |
| exposed-api |  | - | - |
| └─docker-2375 |  | :x: | :x: |
| lxcfs |  | :x: | :x: |

#### runc

| vul | desc | check | exploit |
|-----|------|-------|---------|
| [cve-2016-8867](./vul/cve-2016-8867) | ambient capabilities allow local users to gain privileges | :heavy_check_mark: | :heavy_check_mark: |
| [cve-2019-5736](./vul/cve-2019-5736) | escape by overwrite runc executable file via /proc/self/exe | :heavy_check_mark: | - |
| └─[exec](./vul/cve-2019-5736/exec) | cve-2019-5736 exploit via runc exec process | :heavy_check_mark: | :heavy_check_mark: |
| └─[image](./vul/cve-2019-5736/image) | cve-2019-5736 exploit via a malicious image | :heavy_check_mark: | :heavy_check_mark: |
| cve-2019-16884 |  | :x: | :heavy_check_mark: |
| cve-2023-28642 |  | :x: | :x: |
| cve-2024-21626 |  | :x: | :x: |
| cve-2025-31133 |  | :x: | :x: |
| cve-2025-52565 |  | :x: | :x: |
| cve-2025-52881 |  | :x: | :x: |

#### containerd

| vul | desc | check | exploit |
|-----|------|-------|---------|
| cve-2020-15157 |  | :x: | :heavy_check_mark: |
| [cve-2020-15257](./vul/cve-2020-15257) | abuse the containerd-shim's abstract unix socket in a container with host network namespace | :heavy_check_mark: | :heavy_check_mark: |
| [cve-2024-40635](./vul/cve-2024-40635) | bypass runAsNonRoot via integer overflow in User ID handling in containerd | :heavy_check_mark: | :heavy_check_mark: |
| [cve-2025-47290](./vul/cve-2025-47290) | modify host file via containerd cve-2025-47290 during pulling image | :heavy_check_mark: | :heavy_check_mark: |

#### docker

| vul | desc | check | exploit |
|-----|------|-------|---------|
| [docker.sock](./vul/shared-socket/docker-sock) | escape by shared docker.sock via running a privileged container | :heavy_check_mark: | :heavy_check_mark: |
| cve-2016-9962 |  | :x: | :x: |
| cve-2019-14271 |  | :x: | :heavy_check_mark: |
| cve-2021-41091 |  | :x: | :heavy_check_mark: |
| cve-2021-21285 |  | :x: | :heavy_check_mark: |

#### nvidia-container-toolkit

| vul | desc | check | exploit |
|-----|------|-------|---------|
| [cve-2024-0132](./vul/cve-2024-0132) | gpu container escape via nvidia-container-toolkit CVE-2024-0132 | :o: | :heavy_check_mark: |
| [cve-2025-23266](./vul/cve-2025-23266) | gpu container escape via nvidia-container-toolkit cve-2025-23266 by running a malicious container image | :o: | :heavy_check_mark: |
| cve-2025-23267 |  | :x: | :x: |
| cve-2025-23359 |  | :x: | :x: |

#### docker-compose

| vul | desc | check | exploit |
|-----|------|-------|---------|
| [cve-2025-62725](./vul/cve-2025-62725) | path traversal in docker compose oci artifacts allows arbitrary file write via malicious registry | :heavy_check_mark: | :heavy_check_mark: |

#### buildkit

| vul | desc | check | exploit |
|-----|------|-------|---------|
| [cve-2024-23650](./vul/cve-2024-23650) | dos buildkit via oci exporter by sending a crafted request | :heavy_check_mark: | :heavy_check_mark: |

#### kubernetes

| vul | desc | check | exploit |
|-----|------|-------|---------|
| cve-2017-1002101 |  | :x: | :heavy_check_mark: |
| cve-2020-8555 |  | :x: | :heavy_check_mark: |
| [cve-2020-8558](./vul/cve-2020-8558) | access services bound to 127.0.0.1 from adjacent hosts | :heavy_check_mark: | :heavy_check_mark: |
| [cve-2021-25741](./vul/cve-2021-25741) | kubelet symlink exchange vulnerability allows mounting node filesystem inside a pod | :heavy_check_mark: | :heavy_check_mark: |

#### ingress-nginx

| vul | desc | check | exploit |
|-----|------|-------|---------|
| [cve-2021-25748](./vul/cve-2021-25748) | ingress-nginx path validation bypass vulnerability allows credential leakage through newline injection | :heavy_check_mark: | :heavy_check_mark: |

#### git

| vul | desc | check | exploit |
|-----|------|-------|---------|
| [cve-2022-39253](./vul/cve-2022-39253) | read host file during docker build via git CVE-2022-39253 | :o: | :heavy_check_mark: |

#### kernel

| vul | desc | check | exploit |
|-----|------|-------|---------|
| cve-2021-22555 |  | :x: | :heavy_check_mark: |
| [cve-2021-3493](./vul/cve-2021-3493) | local privilege escalation in Ubuntu OverlayFS that may lead to container escape when the kernel attack surface is exposed | :x: | :x: |
| [cve-2022-0492](./vul/cve-2022-0492) | escape via cgroup's release agent without CAP_SYS_ADMIN if kernel is vulnerable to CVE-2022-0492 | :heavy_check_mark: | :heavy_check_mark: |
| [cve-2022-0847](./vul/cve-2022-0847) | local privilege escalation and container escape in Linux kernel Dirty Pipe | :o: | :heavy_check_mark: |
| [cve-2026-31431](./vul/cve-2026-31431) | local privilege escalation and container escape in Linux kernel AF_ALG AEAD interface | :o: | :heavy_check_mark: |
| [cve-2026-43284](./vul/cve-2026-43284) | local privilege escalation and container escape in Linux kernel xfrm ESP Dirty Frag path | :o: | :heavy_check_mark: |
| [cve-2026-43500](./vul/cve-2026-43500) | local privilege escalation in Linux kernel RxRPC/rxkad Dirty Frag path | :o: | :heavy_check_mark: |

### exploit

```shell
$ ctrsploit exploit
NAME:
   ctrsploit exploit - run a exploit

USAGE:
   ctrsploit exploit [command [command options]]

COMMANDS:
   cve-2016-8867, 8867, amb                           Ambient Capabilities in the Linux kernel allow local users to gain privileges
   cve-2019-5736, 5736                                escape by overwrite runc executable file via /proc/self/exe
   cve-2020-8558, 8558                                access services bound to 127.0.0.1 from adjacent hosts
   cve-2020-15257, 15257                              abuse the containerd-shim's abstract unix socket in a container with host network namespace
   cve-2021-25741, 25741, kubelet-subpath-symlink     kubelet symlink exchange vulnerability allows mounting node filesystem inside a pod
   cve-2021-25748, 25748, ingress-nginx-path-leak     ingress-nginx path validation bypass vulnerability allows credential leakage through newline injection
   cve-2022-0492, 0492                                escape via cgroup's release agent without CAP_SYS_ADMIN if kernel is vulnerable to CVE-2022-0492
   cve-2022-0847, 0847, dirty-pipe, dirtypipe, dp, CVE-2022-0847  local privilege escalation and container escape in Linux kernel Dirty Pipe
   cve-2022-39253, 39253                              read host file during docker build via git CVE-2022-39253
   cve-2024-0132, 0132                                gpu container escape via nvidia-container-toolkit CVE-2024-0132
   cve-2024-23650, 23650                              dos buildkit via oci exporter by sending a crafted request
   cve-2024-40635, 40635                              bypass runAsNonRoot via integer overflow in User ID handling in containerd
   cve-2025-23266, 23266                              gpu container escape via nvidia-container-toolkit cve-2025-23266 by running a malicious container image
   cve-2025-47290, 47290                              modify host file via containerd cve-2025-47290 during pulling image
   cve-2026-31431, 31431, copy-fail                   local privilege escalation and container escape in Linux kernel AF_ALG AEAD interface
   cve-2026-43284, 43284, dirty-frag, dirtyfrag       local privilege escalation and container escape in Linux kernel xfrm ESP Dirty Frag path
   cve-2026-43500, 43500, dirty-frag-rxrpc, dirtyfrag-rxrpc  local privilege escalation in Linux kernel RxRPC/rxkad Dirty Frag path
   fork-bomb                                          
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
   crash, c                                           make container crash

OPTIONS:
   --help, -h  show help
```

### checksec

```shell
$ ctrsploit checksec
NAME:
   ctrsploit checksec - check security inside a container

USAGE:
   ctrsploit checksec [command [command options]]

COMMANDS:
   auto, a                                             auto check security
   env, e                                              gather information
   cve-2016-8867, 8867, amb                            Ambient Capabilities in the Linux kernel allow local users to gain privileges
   cve-2019-5736, 5736                                 escape by overwrite runc executable file via /proc/self/exe
   cve-2020-8558, 8558                                 access services bound to 127.0.0.1 from adjacent hosts
   cve-2020-15257, 15257                               abuse the containerd-shim's abstract unix socket in a container with host network namespace
   cve-2021-25741, 25741, kubelet-subpath-symlink      kubelet symlink exchange vulnerability allows mounting node filesystem inside a pod
   cve-2021-25748, 25748, ingress-nginx-path-leak      ingress-nginx path validation bypass vulnerability allows credential leakage through newline injection
   cve-2021-3493, 3493, ubuntu-overlayfs-pe, CVE-2021-3493  local privilege escalation in Ubuntu OverlayFS that may lead to container escape when the kernel attack surface is exposed
   cve-2022-0492, 0492                                 escape via cgroup's release agent without CAP_SYS_ADMIN if kernel is vulnerable to CVE-2022-0492
   cve-2022-0847, 0847, dirty-pipe, dirtypipe, dp, CVE-2022-0847  local privilege escalation and container escape in Linux kernel Dirty Pipe
   cve-2022-39253, 39253                               read host file during docker build via git CVE-2022-39253
   cve-2024-0132, 0132                                 gpu container escape via nvidia-container-toolkit CVE-2024-0132
   cve-2024-23650, 23650                               dos buildkit via oci exporter by sending a crafted request
   cve-2024-40635, 40635                               bypass runAsNonRoot via integer overflow in User ID handling in containerd
   cve-2025-23266, 23266                               gpu container escape via nvidia-container-toolkit cve-2025-23266 by running a malicious container image
   cve-2025-47290, 47290                               modify host file via containerd cve-2025-47290 during pulling image
   cve-2025-62725, 62725                               path traversal in Docker Compose OCI artifacts allows arbitrary file write via malicious registry
    cve-2026-31431, 31431, copy-fail                    local privilege escalation and container escape in Linux kernel AF_ALG AEAD interface
    cve-2026-43284, 43284, dirty-frag, dirtyfrag        local privilege escalation and container escape in Linux kernel xfrm ESP Dirty Frag path
    cve-2026-43500, 43500, dirty-frag-rxrpc, dirtyfrag-rxrpc  local privilege escalation in Linux kernel RxRPC/rxkad Dirty Frag path
    fork-bomb                                           
    shocker, cap_dac_read_search, open_by_handle_at     escape by CAP_DAC_READ_SEARCH, alias shocker, found by Sebastian Krahmer (stealth) in 2014
   cap_sys_admin, sys_admin                            abuse cap_sys_admin
   cap_bpf, bpf                                        load evil bpf programs via cap_bpf
   cap_sys_ptrace, sys_ptrace, ptrace                  abuse cap_sys_ptrace
   ptrace-pid-host, ptrace-pid                         ptrace host processes in a container with cap_sys_ptrace and host pid namespace
   naked                                               we call containers running without seccomp, AppArmor, or SELinux enabled 'naked containers', which leaves them highly vulnerable to kernel exploits and potential container escapes
   host-net, net                                       shared host network namespace breaks the network isolation
   host-pid, pid                                       shared host pid namespace breaks process isolation
   sa-token-access-secrets, secret                     Check if service account token can access Kubernetes Secrets
   sa-token-policy, policy, dangerous-permissions, dp  Check if service account token has dangerous permissions
   docker.sock, docker                                 escape by shared docker.sock via running a privileged container

OPTIONS:
   --help, -h  show help
```

## Progress of Development

* [vul](./vul/README.md)
