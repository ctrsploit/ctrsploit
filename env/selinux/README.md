---
maintainer:
    - ssst0n3
---
# SELinux

```
root@b1f9f8da70c3:~/app# ./bin/release/env_linux_amd64 --colorful selinux
===========SELinux===========
✘  Enabled	

root@b1f9f8da70c3:~/app# ./bin/release/env_linux_amd64 selinux
===========SELinux===========
[N]  Enabled	

root@b1f9f8da70c3:~/app# ./bin/release/env_linux_amd64 --json selinux
{"kernel_supported":false,"container_enabled":false,"mode":"disabled","mount_point":""}
```