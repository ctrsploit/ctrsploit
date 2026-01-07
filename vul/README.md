* :heavy_check_mark: : Fully Supported
* :o: : Partially Supported
* :bug: : Known Bug
* :x: : Not Supported
* `-` : Not Applicable

## module

### containerd

| vul | desc | check | exploit | test | doc | video | case |
|-----|------|-------|---------|------|-----|-------|------|
| [cve-2024-40635](cve-2024-40635) | bypass runAsNonRoot via integer overflow in User ID handling in containerd | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :x: | :x: |

### kernel

| vul | desc | check | exploit | test | doc | video | case |
|-----|------|-------|---------|------|-----|-------|------|
| [cve-2022-0492](cve-2022-0492) | escape via cgroup's release agent without CAP_SYS_ADMIN if kernel is vulnerable to CVE-2022-0492 | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :x: |

### runc

| vul | desc | check | exploit | test | doc | video | case |
|-----|------|-------|---------|------|-----|-------|------|
| [cve-2019-5736](cve-2019-5736) | escape by overwrite runc executable file via /proc/self/exe | :heavy_check_mark: | - | - | :heavy_check_mark: | - | :x: |
| └─[exec](cve-2019-5736/exec) | cve-2019-5736 exploit via runc exec process | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | - | :heavy_check_mark: | :x: |
| └─[image](cve-2019-5736/image) | cve-2019-5736 exploit via a malicious image | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | - | :heavy_check_mark: | :x: |

### nvidia-container-toolkit

| vul | desc | check | exploit | test | doc | video | case |
|-----|------|-------|---------|------|-----|-------|------|
| [cve-2024-0132](cve-2024-0132) | gpu container escape via nvidia-container-toolkit CVE-2024-0132 | :o: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :x: |

## other

| vul | desc | check | exploit | test | doc | video | case |
|-----|------|-------|---------|------|-----|-------|------|
| [cve-2016-8867](cve-2016-8867) | ambient capabilities allow local users to gain privileges | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :x: |
| cve-2016-9962 |  | :x: | :x: | :x: | :x: | :x: | :x: |
| CVE-2017-1002101 |  | :x: | :heavy_check_mark: | :x: | :x: | :x: | :x: |
| CVE-2019-14271 |  | :x: | :heavy_check_mark: | :x: | :x: | :x: | :x: |
| CVE-2019-16884 |  | :x: | :heavy_check_mark: | :x: | :x: | :x: | :x: |
| CVE-2020-8555 |  | :x: | :heavy_check_mark: | :x: | :x: | :x: | :x: |
| [cve-2020-8558](cve-2020-8558) | access services bound to 127.0.0.1 from adjacent hosts | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :x: |
| CVE-2020-15157 |  | :x: | :heavy_check_mark: | :x: | :x: | :x: | :x: |
| [cve-2020-15257](cve-2020-15257) | abuse the containerd-shim's abstract unix socket in a container with host network namespace | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :x: |
| CVE-2021-3493 |  | :x: | :heavy_check_mark: | :x: | :x: | :x: | :x: |
| CVE-2021-21285 |  | :x: | :heavy_check_mark: | :x: | :x: | :x: | :x: |
| CVE-2021-22555 |  | :x: | :heavy_check_mark: | :x: | :x: | :x: | :x: |
| CVE-2021-41091 |  | :x: | :heavy_check_mark: | :x: | :x: | :x: | :x: |
| [cve-2021-25741](cve-2021-25741) | kubelet symlink exchange vulnerability allows mounting node filesystem inside a pod | :heavy_check_mark: | :heavy_check_mark: | :x: | :heavy_check_mark: | :x: | :x: |
| CVE-2022-0847 |  | :x: | :heavy_check_mark: | :x: | :x: | :x: | :x: |
| [cve-2022-39253](cve-2022-39253) | read host file during docker build via git CVE-2022-39253 | :o: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :x: |
| CVE-2023-28642 |  | :x: | :x: | :x: | :x: | :x: | :x: |
| CVE-2024-21626 |  | :x: | :x: | :x: | :x: | :x: | :x: |
| [cve-2024-23650](cve-2024-23650) | dos buildkit via oci exporter by sending a crafted request | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :x: |
| [cve-2025-23266](cve-2025-23266) | gpu container escape via nvidia-container-toolkit cve-2025-23266 by running a malicious container image | :o: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :x: |
| cve-2025-23267 |  | :x: | :x: | :x: | :x: | :x: | :x: |
| cve-2025-23359 |  | :x: | :x: | :x: | :x: | :x: | :x: |
| cve-2025-31133 |  | :x: | :x: | :x: | :x: | :x: | :x: |
| [cve-2025-47290](cve-2025-47290) | modify host file via containerd cve-2025-47290 during pulling image | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :x: |
| cve-2025-52565 |  | :x: | :x: | :x: | :x: | :x: | :x: |
| [cve-2025-62725](cve-2025-62725) | path traversal in docker compose oci artifacts allows arbitrary file write via malicious registry | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :x: |
| [caps](caps) | abuse dangerous capabilities in container | - | - | - | - | - | - |
| └─[shocker](caps/shocker) | escape by CAP_DAC_READ_SEARCH, alias shocker, found by Sebastian Krahmer (stealth) in 2014 | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark:️ | :heavy_check_mark: | :heavy_check_mark: | :x: |
| └─[sys_admin](caps/sys_admin) | abuse cap_sys_admin | :heavy_check_mark: | - | - | - | - | - |
| &emsp;└─[release_agent](caps/sys_admin/release_agent) | escape by cap_sys_admin via cgroups v1 release_agent | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :x: |
| &emsp;└─mount-device |  | :x: | :x: | :x: | :x: | :x: | :x: |
| &emsp;└─mount-proc |  | :x: | :x: | :x: | :x: | :x: | :x: |
| &emsp;└─device.allow |  | :x: | :x: | :x: | :x: | :x: | :x: |
| &emsp;└─[ebpf](caps/sys_admin/ebpf) | escape by loading evil eBPF programs into the kernel | :heavy_check_mark: | - | - | - | - | - |
| &emsp;&emsp;└─[bash](caps/sys_admin/ebpf/bash) | abuse eBPF to inject malicious commands into bash processes running on host | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :x: |
| &emsp;&emsp;└─[cron](caps/sys_admin/ebpf/cron) | abuse eBPF to inject malicious job into host's crontab | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :x: |
| &emsp;&emsp;└─[execve](caps/sys_admin/ebpf/execve) | abuse eBPF to hijack execve syscall to run arbitrary commands | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :x: |
| &emsp;&emsp;└─[kubelet](caps/sys_admin/ebpf/kubelet) | abuse eBPF to leak services account token from kubelet | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :x: |
| &emsp;&emsp;└─sshd |  | :x: | :x: | :x: | :x: | :x: | :x: |
| └─[bpf](caps/bpf) | load evil bpf programs via cap_bpf | - | - | - | - | - | - |
| &emsp;└─[ebpf](caps/sys_admin/ebpf) | same as caps/sys_admin/ebpf | :heavy_check_mark: | - | - | - | - | - |
| └─[sys_ptrace](caps/sys_ptrace) | abuse cap_sys_ptrace | :heavy_check_mark: | - | - | - | - | - |
| &emsp;└─[pid_host](caps/sys_ptrace/pid_host) | ptrace host processes in a container with cap_sys_ptrace and host pid namespace | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :x: |
| └─sys_module |  | :x: | :x: | :x: | :x: | :x: | :x: |
| └─net_admin |  | :x: | :x: | :x: | :x: | :x: | :x: |
| [naked](naked) | we call containers running without seccomp, AppArmor, or SELinux enabled 'naked containers', which leaves them highly vulnerable to kernel exploits and potential container escapes | :heavy_check_mark: | - | :heavy_check_mark: | :x: | :x: | :x: |
| [namespace](namespace) | shared host namespaces break the isolations | - | - | - | - | - | - |
| └─[net](namespace/net) | shared host network namespace breaks the network isolation | :heavy_check_mark: | :x: | :x: | :x: | :x: | :x: |
| &emsp;└─shijack |  | :x: | :x: | :x: | :x: | :x: | :x: |
| &emsp;&emsp;└─basic |  | :x: | :x: | :x: | :x: | :x: | :x: |
| &emsp;&emsp;└─ali |  | :x: | :x: | :x: | :x: | :x: | :x: |
| &emsp;&emsp;└─hw |  | :x: | :x: | :x: | :x: | :x: | :x: |
| &emsp;&emsp;└─gcp |  | :x: | :x: | :x: | :x: | :x: | :x: |
| &emsp;&emsp;└─aws |  | :x: | :x: | :x: | :x: | :x: | :x: |
| └─[pid](namespace/pid) | shared host pid namespace breaks the process isolation | - | - | - | - | - | - |
| &emsp;└─[proc_root](namespace/pid/proc_root) | escape by abusing host pid ns via /proc/[pid]/root | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :x: |
| fs |  | :x: | :x: | :x: | :x: | :x: | :x: |
| └─proc-rw |  | :x: | - | - | :x: | :x: | :x: |
| &emsp;└─core_pattern |  | :x: | :x: | :x: | :x: | :x: | :x: |
| &emsp;└─binfmt |  | :x: | :x: | :x: | :x: | :x: | :x: |
| └─sys-rw |  | :x: | :x: | :x: | :x: | :x: | :x: |
| └─lxcfs-rw |  | :x: | :x: | :x: | :x: | :x: | :x: |
| [shared-socket](shared-socket) | abuse runtime's api via shared socket | - | - | - | - | - | - |
| └─[docker.sock](shared-socket/docker-sock) | escape by shared docker.sock via running a privileged container | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: |
| └─containerd.sock |  | :x: | :x: | :x: | :x: | :x: | :x: |
| exposed-api |  | - | - | - | - | - | - |
| └─docker-2375 |  | :x: | :x: | :x: | :x: | :x: | :x: |
| lxcfs |  | :x: | :x: | :x: | :x: | :x: | :x: |
| [fork-bomb](fork-bomb) |  | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :x: |
| [sa-token-access-secrets](sa-token/access-secrets) | check if service account token can access Kubernetes Secrets | :heavy_check_mark: | - | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :x: |
