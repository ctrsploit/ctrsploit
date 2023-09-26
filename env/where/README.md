# env/where

To see whether you are in the container, and in which type container

查看当前是否在容器内，在何容器内：

```
root@e7d5f53eb35b:~/ctrsploit/bin/release# ./env_linux_amd64 --colorful where
===========Container===========
✔  Is in Container      
===========Docker===========
✔  dockerenv    # .dockerenv exists
✔  rootfs       # rootfs contains 'docker'
✘  cgroups      # cgroups contains 'docker'
✔  hosts        # the mount source of /etc/hosts contains 'docker'
✔  hostname     # hostname match regex ^[0-9a-f]{12}$
✔  Is in docker 
===========K8S===========
✘  secret       # secret path /var/run/secrets/kubernetes.io exists
✘  hostname     # hostname match k8s pattern
✘  hosts        # the mount source of /etc/hosts contains 'pods'
✘  cgroups      # cgroups contains 'kubepods'
✘  is in k8s 
```