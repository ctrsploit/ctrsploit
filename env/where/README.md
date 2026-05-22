---
maintainer:
    - ssst0n3
---
# env/where

To see whether you are in the container, and in which type container

查看当前是否在容器内，在何容器内：

```
root@b1f9f8da70c3:~/app# ./bin/release/env_linux_amd64 where
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

root@b1f9f8da70c3:~/app# ./bin/release/env_linux_amd64 --colorful where
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
===========Container===========
✔  Is in Container	

root@b1f9f8da70c3:~/app# ./bin/release/env_linux_amd64 --json where
{"container":{"in":true,"rules":{}},"k8s":{"in":false,"rules":{"cgroups":false,"hostname":false,"hosts":false,"secret":false}},"containerd":{"in":false,"rules":null},"docker":{"in":true,"rules":{"cgroups":false,"dockerenv":true,"hostname":true,"hosts":true,"rootfs":true}}}
```