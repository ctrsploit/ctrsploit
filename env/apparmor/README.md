# apparmor
## usage

```
root@b93597d64108:/# ./env_linux_amd64 --colorful apparmor
===========AppArmor===========
✔  Kernel Supported     # Kernel enabled apparmor module
✔  Container Enabled    # Current container enabled apparmor
Profile:                docker-default (enforce)        
Mode:                   enforce

root@b93597d64108:/# ./env_linux_amd64 apparmor
===========AppArmor===========
[Y]  Kernel Supported   # Kernel enabled apparmor module
[Y]  Container Enabled  # Current container enabled apparmor
Profile:                docker-default (enforce)        
Mode:                   enforce

root@c551710496f3:/# ./env_linux_amd64 --json a
{"name":{"name":"AppArmor"},"kernel":{"name":"Kernel Supported","description":"Kernel enabled apparmor module","result":true},"container":{"name":"Container Enabled","description":"Current container enabled apparmor","result":true},"profile":{"name":"Profile","description":"","result":"docker-default (enforce)"},"mode":{"name":"Mode","description":"","result":"enforce"}} 
```
