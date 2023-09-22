# env/auto

```
$ docker run -ti -v $(pwd):/ctrsploit --cap-add cap_sys_admin ubuntu bash
root@08e2905722f1:/# cd /ctrsploit/bin/release/
root@08e2905722f1:/ctrsploit/bin/release# ./env_linux_amd64 --colorful auto

===========Container===========
✔  Is in Container

===========Docker===========
✔  .dockerenv exists
✘  rootfs contains 'docker'     
✘  cgroups contains 'docker'
✔  the mount source of /etc/hosts contains 'docker'     
✔  hostname match regex ^[0-9a-f]12$
---
✔  => Is in docker

===========k8s===========
✘  /var/run/secrets/kubernetes.io exists
✘  hostname match k8s pattern
✘  the mount source of /etc/hosts contains 'pods'
✘  contains 'kubepods'
---
✘  => is in k8s

===========Apparmor===========
✘  Kernel Supported
✘  Container Enabled

===========SELinux===========
✘  Enabled

===========Capability(process 1)===========
[Capabilities]
0xa82425fb != 0xa80425fb(docker's default caps)
[Additional Capabilities]
["CAP_SYS_ADMIN"]

===========Capability(process current)===========
[Capabilities]
0xa82425fb != 0xa80425fb(docker's default caps)
[Additional Capabilities]

["CAP_SYS_ADMIN"]

===========Cgroups===========
✘  cgroups v1
✔  cgroups v2


===========Overlay===========
✘  Enabled

===========DeviceMapper===========
✔  Enabled
✘  Used

===========namespace level===========
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

===========Seccomp===========
✔  Kernel Supported
✔  Container Enabled
---
mode:   filter
```