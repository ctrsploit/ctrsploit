---
maintainer:
    - ssst0n3
---
# env/auto

```
root@b1f9f8da70c3:~/app# ./bin/release/env_linux_amd64 --colorful auto
===========Container===========
✔  Is in Container	
===========Docker===========
✔  dockerenv	# .dockerenv exists
✔  rootfs	# rootfs contains 'docker'
✘  cgroups	# cgroups contains 'docker'
✔  hosts	# the mount source of /etc/hosts contains 'docker'
✔  hostname	# hostname match regex ^[0-9a-f]{12}$
✔  Is in docker	
===========K8S===========
✘  secret	# secret path /var/run/secrets/kubernetes.io exists
✘  hostname	# hostname match k8s pattern
✘  hosts	# the mount source of /etc/hosts contains 'pods'
✘  cgroups	# cgroups contains 'kubepods'
✘  is in k8s	
===========AppArmor===========
✔  Kernel Supported	# Kernel enabled apparmor module
✔  Container Enabled	# Current container enabled apparmor
Profile:		docker-default (enforce)	
Mode:			enforce	
===========SELinux===========
✘  Enabled	
Mode:			disabled	
Mount point:			
===========Capability===========
[pid1]
capabilities:		0xa80425fb	
✘  Not Equal to Docker's Default Capability (0xa80425fb)	# 0xa80425fb
[current]
capabilities:		0xa80425fb	
✘  Not Equal to Docker's Default Capability (0xa80425fb)	# 0xa80425fb
===========CGroups===========
✘  v1	
✔  v2	
===========Overlay===========
✔  Enabled	
✔  Used	
The number of graph driver mounted	# equal to the number of containers
5
The host path of container's rootfs	
/var/lib/docker/overlay2/254080c7a7418fb2db4d5d27e03bc86f3fe4ae52b5c15b36ad0196adfc887f9e/merged
===========DeviceMapper===========
✔  Enabled	
✘  Used	
The number of graph driver mounted	# equal to the number of containers
0
The host path of container's rootfs	

===========Namespace Level===========
cgroup:			child	
ipc:			child	
mnt:			child	
net:			child	
pid:			child	
pid_for_children:	child	
time:			host	
time_for_children:	host	
user:			host	
uts:			child	
===========Seccomp===========
✔  Kernel Supported	
✔  Container Enabled	
Mode:			filter	

root@b1f9f8da70c3:~/app# ./bin/release/env_linux_amd64 auto
===========Container===========
[Y]  Is in Container	
===========Docker===========
[Y]  dockerenv	# .dockerenv exists
[Y]  rootfs	# rootfs contains 'docker'
[N]  cgroups	# cgroups contains 'docker'
[Y]  hosts	# the mount source of /etc/hosts contains 'docker'
[Y]  hostname	# hostname match regex ^[0-9a-f]{12}$
[Y]  Is in docker	
===========K8S===========
[N]  secret	# secret path /var/run/secrets/kubernetes.io exists
[N]  hostname	# hostname match k8s pattern
[N]  hosts	# the mount source of /etc/hosts contains 'pods'
[N]  cgroups	# cgroups contains 'kubepods'
[N]  is in k8s	
===========AppArmor===========
[Y]  Kernel Supported	# Kernel enabled apparmor module
[Y]  Container Enabled	# Current container enabled apparmor
Profile:		docker-default (enforce)	
Mode:			enforce	
===========SELinux===========
[N]  Enabled	
Mode:			disabled	
Mount point:			
===========Capability===========
[pid1]
capabilities:		0xa80425fb	
[N]  Not Equal to Docker's Default Capability (0xa80425fb)	# 0xa80425fb
[current]
capabilities:		0xa80425fb	
[N]  Not Equal to Docker's Default Capability (0xa80425fb)	# 0xa80425fb
===========CGroups===========
[N]  v1	
[Y]  v2	
===========DeviceMapper===========
[Y]  Enabled	
[N]  Used	
The number of graph driver mounted	# equal to the number of containers
0
The host path of container's rootfs	

===========Overlay===========
[Y]  Enabled	
[Y]  Used	
The number of graph driver mounted	# equal to the number of containers
5
The host path of container's rootfs	
/var/lib/docker/overlay2/254080c7a7418fb2db4d5d27e03bc86f3fe4ae52b5c15b36ad0196adfc887f9e/merged
===========Namespace Level===========
cgroup:			child	
ipc:			child	
mnt:			child	
net:			child	
pid:			child	
pid_for_children:	child	
time:			host	
time_for_children:	host	
user:			host	
uts:			child	
===========Seccomp===========
[Y]  Kernel Supported	
[Y]  Container Enabled	
Mode:			filter	

root@b1f9f8da70c3:~/app# ./bin/release/env_linux_amd64 --json auto |jq
{
  "where": {
    "container": {
      "in": true,
      "rules": {}
    },
    "k8s": {
      "in": false,
      "rules": {
        "cgroups": false,
        "hostname": false,
        "hosts": false,
        "secret": false
      }
    },
    "containerd": {
      "in": false,
      "rules": null
    },
    "docker": {
      "in": true,
      "rules": {
        "cgroups": false,
        "dockerenv": true,
        "hostname": true,
        "hosts": true,
        "rootfs": true
      }
    }
  },
  "kernel_version": "",
  "credential": {
    "uid": 0,
    "gid": 0
  },
  "capability": {
    "pid1": 2818844155,
    "self": 2818844155
  },
  "apparmor": {
    "kernel_supported": true,
    "container_enabled": true,
    "profile": "docker-default (enforce)",
    "mode": "enforce"
  },
  "selinux": {
    "kernel_supported": false,
    "container_enabled": false,
    "mode": "disabled",
    "mount_point": ""
  },
  "seccomp": {
    "kernel_supported": true,
    "container_enabled": true,
    "mode": "filter"
  },
  "namespace": {
    "levels": {
      "cgroup": 2,
      "ipc": 2,
      "mnt": 2,
      "net": 2,
      "pid": 2,
      "pid_for_children": 2,
      "time": 1,
      "time_for_children": 1,
      "user": 1,
      "uts": 2
    }
  },
  "cgroups": {
    "version": 2,
    "sub": [],
    "top": []
  }, 
  "overlay": {
    "loaded": true,
    "used": true,
    "refcnt": 5,
    "hostPath": "/var/lib/docker/overlay2/254080c7a7418fb2db4d5d27e03bc86f3fe4ae52b5c15b36ad0196adfc887f9e/merged"
  },
  "device_mapper": {
    "loaded": true,
    "used": false,
    "refcnt": 0,
    "hostPath": ""
  },
  "version": "",
  "apiserver": "",
  "pod_env": null,
  "pod_mount": null,
  "dns_service": null,
  "apiserver_access": 0,
  "kublet_access": 0,
  "hostpath": null,
  "credentials": null,
  "plugins": null,
  "config": {
    "ConfigMap": null,
    "Secret": null,
    "NameSpaces": null,
    "Pods": null
  },
  "custom_resources": null,
  "runtime_version": {},
  "ctr_cnt": {}
}
```