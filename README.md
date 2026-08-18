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
| └─[shocker](./vul/caps/shocker) | escape by CAP_DAC_READ_SEARCH | :heavy_check_mark: | :heavy_check_mark: |
| └─[sys_admin](./vul/caps/sys_admin) | abuse cap_sys_admin | :heavy_check_mark: | - |
| &emsp;└─[release_agent](./vul/caps/sys_admin/release_agent) | escape via cgroups v1 release_agent | :heavy_check_mark: | :heavy_check_mark: |
| &emsp;└─[ebpf](./vul/caps/sys_admin/ebpf) | escape by loading evil eBPF programs into the kernel | :heavy_check_mark: | - |
| &emsp;&emsp;└─[bash](./vul/caps/sys_admin/ebpf/bash) | inject malicious commands into host bash processes | :heavy_check_mark: | :heavy_check_mark: |
| &emsp;&emsp;└─[cron](./vul/caps/sys_admin/ebpf/cron) | inject malicious job into host's crontab | :heavy_check_mark: | :heavy_check_mark: |
| &emsp;&emsp;└─[execve](./vul/caps/sys_admin/ebpf/execve) | hijack execve syscall | :heavy_check_mark: | :heavy_check_mark: |
| &emsp;&emsp;└─[kubelet](./vul/caps/sys_admin/ebpf/kubelet) | leak service account token from kubelet | :heavy_check_mark: | :heavy_check_mark: |
| └─[bpf](./vul/caps/bpf) | load evil bpf programs via cap_bpf | - | - |
| └─[sys_ptrace](./vul/caps/sys_ptrace) | abuse cap_sys_ptrace | :heavy_check_mark: | - |
| &emsp;└─[pid_host](./vul/caps/sys_ptrace/pid_host) | ptrace host processes | :heavy_check_mark: | :heavy_check_mark: |
| └─[cve-2016-8867](./vul/cve-2016-8867) | ambient capabilities allow local users to gain privileges | :heavy_check_mark: | :heavy_check_mark: |
| [naked](./vul/naked) | containers running without seccomp, AppArmor, or SELinux | :heavy_check_mark: | - |
| [shared-socket](./vul/shared-socket) | abuse runtime's api via shared socket | - | - |
| └─[docker.sock](./vul/shared-socket/docker-sock) | escape by shared docker.sock | :heavy_check_mark: | :heavy_check_mark: |
| [kubeconfig](./vul/kubeconfig) | check kubeconfig related vulnerabilities | - | - |
| └─[user-exec](./vul/kubeconfig/user-exec) | loading an untrusted kubeconfig can execute arbitrary client-side commands via users[].user.exec | :heavy_check_mark: | - |
| [sa-token](./vul/sa-token) | check service account token related vulnerabilities | - | - |
| └─[secret](./vul/sa-token/access-secrets) | check if service account token can access Kubernetes Secrets | :heavy_check_mark: | - |
| └─[policy](./vul/sa-token/policy) | check if service account token has dangerous permissions | :heavy_check_mark: | - |
| [namespace](./vul/namespace) | shared host namespaces break the isolations | - | - |
| └─[net](./vul/namespace/net) | shared host network namespace | :heavy_check_mark: | :x: |
| └─[pid](./vul/namespace/pid) | shared host pid namespace | - | - |
| &emsp;└─[proc_root](./vul/namespace/pid/proc_root) | escape via /proc/[pid]/root | :heavy_check_mark: | :heavy_check_mark: |

### runc

| vul | desc | check | exploit |
|-----|------|-------|---------|
| [cve-2016-8867](./vul/cve-2016-8867) | ambient capabilities allow local users to gain privileges | :heavy_check_mark: | :heavy_check_mark: |
| [cve-2019-5736](./vul/cve-2019-5736) | escape by overwrite runc executable via /proc/self/exe | :heavy_check_mark: | - |
| └─[exec](./vul/cve-2019-5736/exec) | exploit via runc exec process | :heavy_check_mark: | :heavy_check_mark: |
| └─[image](./vul/cve-2019-5736/image) | exploit via malicious image | :heavy_check_mark: | :heavy_check_mark: |
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
| [cve-2020-15257](./vul/cve-2020-15257) | abuse containerd-shim abstract unix socket | :heavy_check_mark: | :heavy_check_mark: |
| [cve-2024-40635](./vul/cve-2024-40635) | bypass runAsNonRoot via integer overflow | :heavy_check_mark: | :heavy_check_mark: |
| [cve-2025-47290](./vul/cve-2025-47290) | modify host file during image pull | :heavy_check_mark: | :heavy_check_mark: |

### docker

| vul | desc | check | exploit |
|-----|------|-------|---------|
| [docker.sock](./vul/shared-socket/docker-sock) | escape by shared docker.sock | :heavy_check_mark: | :heavy_check_mark: |
| cve-2016-9962 |  | :x: | :x: |
| cve-2019-14271 |  | :x: | :heavy_check_mark: |
| cve-2021-41091 |  | :x: | :heavy_check_mark: |
| cve-2021-21285 |  | :x: | :heavy_check_mark: |

### kernel

| vul | desc | check | exploit |
|-----|------|-------|---------|
| cve-2021-22555 |  | :x: | :heavy_check_mark: |
| [cve-2021-3493](./vul/cve-2021-3493) | local privilege escalation in Ubuntu OverlayFS | :x: | :x: |
| [cve-2022-0492](./vul/cve-2022-0492) | escape via cgroup release_agent without CAP_SYS_ADMIN | :heavy_check_mark: | :heavy_check_mark: |
| [cve-2022-0847](./vul/cve-2022-0847) | local privilege escalation and container escape via Dirty Pipe | :o: | :heavy_check_mark: |
| [cve-2026-31431](./vul/cve-2026-31431) | local privilege escalation and container escape via AF_ALG AEAD | :o: | :heavy_check_mark: |
| [cve-2026-43284](./vul/cve-2026-43284) | local privilege escalation and container escape via xfrm ESP Dirty Frag | :o: | :heavy_check_mark: |
| [cve-2026-43500](./vul/cve-2026-43500) | local privilege escalation and container escape via RxRPC/rxkad Dirty Frag | :o: | :heavy_check_mark: |
| [cve-2026-46300](./vul/cve-2026-46300) | local privilege escalation and container escape via xfrm ESP-in-TCP Fragnesia | :o: | :heavy_check_mark: |
| [cve-2026-23111](./vul/cve-2026-23111) | local privilege escalation via nf_tables use-after-free (inverted check in nft_map_catchall_activate) | :heavy_check_mark: | :heavy_check_mark: |

### kubernetes

| vul | desc | check | exploit |
|-----|------|-------|---------|
| cve-2017-1002101 |  | :x: | :heavy_check_mark: |
| cve-2020-8555 |  | :x: | :heavy_check_mark: |
| [cve-2020-8558](./vul/cve-2020-8558) | access 127.0.0.1 services from adjacent hosts | :heavy_check_mark: | :heavy_check_mark: |
| [cve-2021-25741](./vul/cve-2021-25741) | kubelet symlink exchange | :heavy_check_mark: | :heavy_check_mark: |

### nvidia-container-toolkit

| vul | desc | check | exploit |
|-----|------|-------|---------|
| [cve-2024-0132](./vul/cve-2024-0132) | gpu container escape via nvidia-container-toolkit CVE-2024-0132 | :o: | :heavy_check_mark: |
| [cve-2025-23266](./vul/cve-2025-23266) | gpu container escape via nvidia-container-toolkit CVE-2025-23266 by running a malicious container image | :o: | :heavy_check_mark: |
| cve-2025-23267 |  | :x: | :x: |
| cve-2025-23359 |  | :x: | :x: |

### docker-compose

| vul | desc | check | exploit |
|-----|------|-------|---------|
| [cve-2025-62725](./vul/cve-2025-62725) | path traversal in Docker Compose OCI artifacts allows arbitrary file write via malicious registry | :heavy_check_mark: | :heavy_check_mark: |

### buildkit

| vul | desc | check | exploit |
|-----|------|-------|---------|
| [cve-2024-23650](./vul/cve-2024-23650) | dos buildkit via OCI exporter by sending a crafted request | :heavy_check_mark: | :heavy_check_mark: |

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
