---
maintainer:
    - ssst0n3
    - Axsl666
---
# env/capability

```
$ docker run -ti -v $(pwd):/ctrsploit -w /ctrsploit/bin/release --cap-add CAP_SYS_ADMIN ubuntu bash
root@d382ed5d4f04:/ctrsploit/bin/release# ./env_linux_amd64 --colorful capability
===========Capability===========
[pid1]
capabilities:		0xa82425fb	
✔  Not Equal to Docker's Default Capability (0xa80425fb)	# 0xa82425fb
[Additional]	
["CAP_SYS_ADMIN"]
[current]
capabilities:		0xa82425fb	
✔  Not Equal to Docker's Default Capability (0xa80425fb)	# 0xa82425fb
[Additional]	
["CAP_SYS_ADMIN"]

root@d382ed5d4f04:/ctrsploit/bin/release# ./env_linux_amd64 capability
===========Capability===========
[pid1]
capabilities:		0xa82425fb	
[Y]  Not Equal to Docker's Default Capability (0xa80425fb)	# 0xa82425fb
[Additional]	
["CAP_SYS_ADMIN"]
[current]
capabilities:		0xa82425fb	
[Y]  Not Equal to Docker's Default Capability (0xa80425fb)	# 0xa82425fb
[Additional]	
["CAP_SYS_ADMIN"]

root@d382ed5d4f04:/ctrsploit/bin/release# ./env_linux_amd64 --json capability
{"pid1":2820941307,"self":2820941307}
```