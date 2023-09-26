# env/auto

```
$ docker run -ti -v $(pwd):/ctrsploit --cap-add cap_sys_admin ubuntu bash
root@08e2905722f1:/# cd /ctrsploit/bin/release/
root@08e2905722f1:/ctrsploit/bin/release# ./env_linux_amd64 --colorful auto
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

===========AppArmor===========
✘  Kernel Supported     # Kernel enabled apparmor module
✘  Container Enabled    # Current container enabled apparmor

===========SELinux===========
✘  Enabled      

===========Capability===========
[Capabilities (pid1)]   
0xa82425fb
✘  Equal to Docker's Default capability # 0xa82425fb
[Additional]    
["CAP_SYS_ADMIN"]

[Capabilities (current)]        
0xa82425fb
✘  Equal to Docker's Default capability # 0xa82425fb
[Additional]    
["CAP_SYS_ADMIN"]


===========CGroups===========
✘  v1   
✔  v2   

===========Overlay===========
✔  Enabled      
✔  Used 
The host path of container's rootfs     
/var/lib/docker/overlay2/4d3cb214adf736400170531bbe55bffa2c34611bfc73c159a08c467452563cf9/merged

===========DeviceMapper===========
✔  Enabled      
✘  Used 

===========Namespace Level===========
cgroup:                 child   
ipc:                    child   
mnt:                    child   
net:                    child   
pid:                    child   
pid_for_children:       child   
time:                   host    
time_for_children:      host    
user:                   host    
uts:                    child   

===========Seccomp===========
✔  Kernel Supported     
✔  Container Enabled    
Mode:                   filter
```