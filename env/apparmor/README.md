---
maintainer:
    - ssst0n3
---
# apparmor
## usage

```
root@b1f9f8da70c3:~/app# ./bin/release/env_linux_amd64 --colorful apparmor
===========AppArmor===========
✔  Kernel Supported	# Kernel enabled apparmor module
✔  Container Enabled	# Current container enabled apparmor
Profile:		docker-default (enforce)	
Mode:			enforce	

root@b1f9f8da70c3:~/app# ./bin/release/env_linux_amd64 apparmor
===========AppArmor===========
[Y]  Kernel Supported	# Kernel enabled apparmor module
[Y]  Container Enabled	# Current container enabled apparmor
Profile:		docker-default (enforce)	
Mode:			enforce	

root@b1f9f8da70c3:~/app# ./bin/release/env_linux_amd64 --json apparmor
{"kernel_supported":true,"container_enabled":true,"profile":"docker-default (enforce)","mode":"enforce"}
```
