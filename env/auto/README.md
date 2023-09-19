# env/auto

```
$ docker run -ti -v $(pwd):/ctrsploit --cap-add cap_sys_admin ubuntu bash
root@08e2905722f1:/# cd /ctrsploit/bin/release/
root@08e2905722f1:/ctrsploit/bin/release# ./env_linux_amd64 auto
===========Docker=========
.dockerenv exists: ✔
rootfs contains 'docker': ✘
cgroup contains 'docker': ✘
the mount source of /etc/hosts contains 'docker': ✔
hostname match regex ^[0-9a-f]{12}$: ✔
=> is in docker: ✔

===========k8s=========
/var/run/secrets/kubernetes.io exists: ✘
hostname match k8s pattern: ✘
the mount source of /etc/hosts contains 'pods': ✘
cgroup contains 'kubepods': ✘
=> is in k8s: ✘

===========Apparmor=========
Kernel Supported: ✘
Container Enabled: ✘

===========SELinux=========
Enabled: ✘
mode: disabled
SELinux filesystem mount point: 

===========Capability=========
pid 1
[caps]
0xa82425fb != default(0xa80425fb)
[Additional Capabilities]
["CAP_SYS_ADMIN"]

current process
[caps]
0xa82425fb != default(0xa80425fb)
[Additional Capabilities]
["CAP_SYS_ADMIN"]

===========Cgroups=========
is cgroupv1: ✘
is cgroupv2: ✔

===========Overlay=========
Overlay enabled: ✘

===========DeviceMapper=========
DeviceMapper enabled: ✔
DeviceMapper used: ✘

========namespace level=======
cgroup:              child
ipc:                 child
mnt:                 child
net:                 child
pid:                 child
pid_for_children:    child
time:                host
time_for_children:   host
user:                host
uts:                 child

===========Seccomp=========
kernel supported: ✔
seccomp enabled in current container: ✔
seccomp mode: filter


```