---
maintainer:
    - ssst0n3
---
# seccomp
## usage
```
root@b1f9f8da70c3:~/app# ./bin/release/env_linux_amd64 --colorful seccomp
===========Seccomp===========
✔  Kernel Supported	
✔  Container Enabled	
Mode:			filter	

root@b1f9f8da70c3:~/app# ./bin/release/env_linux_amd64 seccomp
===========Seccomp===========
[Y]  Kernel Supported	
[Y]  Container Enabled	
Mode:			filter	

root@b1f9f8da70c3:~/app# ./bin/release/env_linux_amd64 --json seccomp
{"kernel_supported":true,"container_enabled":true,"mode":"filter"}
```