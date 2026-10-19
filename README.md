# ctrsploit: A penetration toolkit for container environment

ctrsploit [kənˈteɪnər splɔɪt], follows [sploit-spec](https://github.com/ctrsploit/sploit-spec)

![](./docs/images/logo-white-256.png)

## Quick Start

```bash
# Download
wget -q https://github.com/ctrsploit/ctrsploit/releases/latest/download/ctrsploit_linux_amd64 -O /usr/bin/ctrsploit
chmod +x /usr/bin/ctrsploit

# Check what's exploitable
ctrsploit checksec auto

# Run an exploit
ctrsploit exploit cve-2022-0847
```

## Vul — Vulnerability Modules

Modules are grouped by component or configuration type. Each module may
support `checksec` (detection) and/or `exploit`.

| Legend | |
|--------|---|
| :heavy_check_mark: | Supported |
| :o: | Partial |
| :bug: | Known bug |
| :x: | Not supported |
| `-` | Not applicable |

### config

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
| [kubeconfig](./vul/kubeconfig) | check kubeconfig related vulnerabilities | - | - |
| └─[user-exec](./vul/kubeconfig/user-exec) | loading an untrusted kubeconfig can execute arbitrary client-side commands via users[].user.exec | :heavy_check_mark: | - |
| [sa-token](./vul/sa-token) | check service account token related vulnerabilities | - | - |
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

### runc

| vul | desc | check | exploit |
|-----|------|-------|---------|
| [cve-2016-8867](./vul/cve-2016-8867) | ambient capabilities allow local users to gain privileges | :heavy_check_mark: | :heavy_check_mark: |
| [cve-2019-5736](./vul/cve-2019-5736) | escape by overwrite runc executable file via /proc/self/exe | :heavy_check_mark: | - |
| └─[exec](./vul/cve-2019-5736/exec) | cve-2019-5736 exploit via runc exec process | :heavy_check_mark: | :heavy_check_mark: |
| └─[image](./vul/cve-2019-5736/image) | cve-2019-5736 exploit via a malicious image | :heavy_check_mark: | :heavy_check_mark: |
| [cve-2021-30465](./vul/cve-2021-30465) | escape via runc mount destination TOCTOU (symlink exchange) | :heavy_check_mark: | - |
| └─[docker](./vul/cve-2021-30465/docker) | cve-2021-30465 exploit against docker | :heavy_check_mark: | :heavy_check_mark: |
| └─[kubernetes](./vul/cve-2021-30465/kubernetes) | cve-2021-30465 exploit against kubernetes (kubelet emptyDir mounts) | :heavy_check_mark: | :heavy_check_mark: |
| └─[podman](./vul/cve-2021-30465/podman) | cve-2021-30465 exploit against podman (tmpcopyup) | :heavy_check_mark: | :heavy_check_mark: |
| cve-2019-16884 |  | :x: | :heavy_check_mark: |
| cve-2023-28642 |  | :x: | :x: |
| cve-2024-21626 |  | :x: | :x: |
| cve-2025-31133 |  | :x: | :x: |
| cve-2025-52565 |  | :x: | :x: |
| cve-2025-52881 |  | :x: | :x: |

### containerd

| vul | desc | check | exploit |
|-----|------|-------|---------|
| cve-2020-15157 |  | :x: | :heavy_check_mark: |
| [cve-2020-15257](./vul/cve-2020-15257) | abuse the containerd-shim's abstract unix socket in a container with host network namespace | :heavy_check_mark: | :heavy_check_mark: |
| [cve-2024-40635](./vul/cve-2024-40635) | bypass runAsNonRoot via integer overflow in User ID handling in containerd | :heavy_check_mark: | :heavy_check_mark: |
| [cve-2025-47290](./vul/cve-2025-47290) | modify host file via containerd cve-2025-47290 during pulling image | :heavy_check_mark: | :heavy_check_mark: |

### docker

| vul | desc | check | exploit |
|-----|------|-------|---------|
| [docker.sock](./vul/shared-socket/docker-sock) | escape by shared docker.sock via running a privileged container | :heavy_check_mark: | :heavy_check_mark: |
| cve-2016-9962 |  | :x: | :x: |
| cve-2019-14271 |  | :x: | :heavy_check_mark: |
| cve-2021-41091 |  | :x: | :heavy_check_mark: |
| cve-2021-21285 |  | :x: | :heavy_check_mark: |

### kernel

| vul | desc | check | exploit |
|-----|------|-------|---------|
| cve-2021-22555 |  | :x: | :heavy_check_mark: |
| [cve-2021-3493](./vul/cve-2021-3493) | local privilege escalation in Ubuntu OverlayFS that may lead to container escape when the kernel attack surface is exposed | :x: | :x: |
| [cve-2022-0492](./vul/cve-2022-0492) | escape via cgroup's release agent without CAP_SYS_ADMIN if kernel is vulnerable to CVE-2022-0492 | :heavy_check_mark: | :heavy_check_mark: |
| [cve-2022-0847](./vul/cve-2022-0847) | local privilege escalation and container escape in Linux kernel Dirty Pipe | :o: | :heavy_check_mark: |
| [cve-2026-31431](./vul/cve-2026-31431) | local privilege escalation and container escape in Linux kernel AF_ALG AEAD interface | :o: | :heavy_check_mark: |
| [cve-2026-43284](./vul/cve-2026-43284) | local privilege escalation and container escape in Linux kernel xfrm ESP Dirty Frag path | :o: | :heavy_check_mark: |
| [cve-2026-43500](./vul/cve-2026-43500) | local privilege escalation in Linux kernel RxRPC/rxkad Dirty Frag path | :o: | :heavy_check_mark: |
| [cve-2026-46300](./vul/cve-2026-46300) | local privilege escalation and container escape in Linux kernel xfrm ESP-in-TCP Fragnesia path | :o: | :heavy_check_mark: |
| [cve-2026-23111](./vul/cve-2026-23111) | local privilege escalation via nf_tables UAF (inverted check in nft_map_catchall_activate) | :heavy_check_mark: | :heavy_check_mark: |

### kubernetes

| vul | desc | check | exploit |
|-----|------|-------|---------|
| cve-2017-1002101 |  | :x: | :heavy_check_mark: |
| cve-2020-8555 |  | :x: | :heavy_check_mark: |
| [cve-2020-8558](./vul/cve-2020-8558) | access services bound to 127.0.0.1 from adjacent hosts | :heavy_check_mark: | :heavy_check_mark: |
| [cve-2021-25741](./vul/cve-2021-25741) | kubelet symlink exchange vulnerability allows mounting node filesystem inside a pod | :heavy_check_mark: | :heavy_check_mark: |

### nvidia-container-toolkit

| vul | desc | check | exploit |
|-----|------|-------|---------|
| [cve-2024-0132](./vul/cve-2024-0132) | gpu container escape via nvidia-container-toolkit CVE-2024-0132 | :o: | :heavy_check_mark: |
| [cve-2025-23266](./vul/cve-2025-23266) | gpu container escape via nvidia-container-toolkit cve-2025-23266 by running a malicious container image | :o: | :heavy_check_mark: |
| cve-2025-23267 |  | :x: | :x: |
| cve-2025-23359 |  | :x: | :x: |

### docker-compose

| vul | desc | check | exploit |
|-----|------|-------|---------|
| [cve-2025-62725](./vul/cve-2025-62725) | path traversal in docker compose oci artifacts allows arbitrary file write via malicious registry | :heavy_check_mark: | :heavy_check_mark: |

### buildkit

| vul | desc | check | exploit |
|-----|------|-------|---------|
| [cve-2024-23650](./vul/cve-2024-23650) | dos buildkit via oci exporter by sending a crafted request | :heavy_check_mark: | :heavy_check_mark: |

### ingress-nginx

| vul | desc | check | exploit |
|-----|------|-------|---------|
| [cve-2021-25748](./vul/cve-2021-25748) | ingress-nginx path validation bypass vulnerability allows credential leakage through newline injection | :heavy_check_mark: | :heavy_check_mark: |

### git

| vul | desc | check | exploit |
|-----|------|-------|---------|
| [cve-2022-39253](./vul/cve-2022-39253) | read host file during docker build via git CVE-2022-39253 | :o: | :heavy_check_mark: |

## env — Environment Gathering

```bash
ctrsploit env auto       # auto-detect container environment
ctrsploit env where      # check if in container and what type
ctrsploit env capability # show capabilities
ctrsploit env seccomp    # show seccomp status
ctrsploit env services   # discover cluster services without API access
ctrsploit env cpusec     # show CPU/kernel security mitigations (SMEP/SMAP/KPTI/IBT/KCFI/FG-KASLR on x86; PAC/BTI/KPTI/PAN/MTE on arm64)
```

For full `env` subcommands and flags, run `ctrsploit env --help`.

## tool — Convenience Tools

Small convenience tools used during penetration testing. No subcommands yet.

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md).
