# apparmor
## usage

```
root@b93597d64108:/# ./env_linux_amd64 --colorful apparmor
===========AppArmor===========
✔  Kernel Supported     # Kernel enabled apparmor module
✔  Container Enabled    # Current container enabled apparmor
Profile:        docker-default (enforce)        
Mode:           enforce 

root@b93597d64108:/# ./env_linux_amd64 apparmor
===========AppArmor===========
[Y]  Kernel Supported   # Kernel enabled apparmor module
[Y]  Container Enabled  # Current container enabled apparmor
Profile:        docker-default (enforce)        
Mode:           enforce 
```
