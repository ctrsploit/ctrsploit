# namespace

## namespace level check

These commands below are same:

```
$ ctrsploit env namespace
$ ctrsploit checksec namespace
$ checksec namespace
```

### Test Environment 2: MacBookPro

```
❯ uname -a    
Darwin MacBook-Pro.local 21.6.0 Darwin Kernel Version 21.6.0: Sat Jun 18 17:07:25 PDT 2022; root:xnu-8020.140.41~1/RELEASE_X86_64 x86_64

❯ docker --version
Docker version 20.10.21, build baeda1f

❯ docker run -ti --rm -v $(pwd):/ctrsploit --workdir /ctrsploit ubuntu 
root@d325c6a6d650:/ctrsploit# uname -a
Linux docker-desktop 5.15.49-linuxkit #1 SMP Tue Sep 13 07:51:46 UTC 2022 x86_64 x86_64 x86_64 GNU/Linux
```

**all namespaces level**

```
❯ docker run -ti --rm -v $(pwd):/ctrsploit --workdir /ctrsploit ubuntu
root@49a71bef0ee8:/ctrsploit# ./checksec_linux_amd64 namespace
========namespace level=======
cgroup:              child
ipc:                 child
mnt:                 child
net:                 child
pid:                 child
pid_for_children:    child
time:                host
time_for_children:   host
user:                host
uts:                 child
```

**cgroups**

```
❯ docker run --cgroupns=host -ti --rm -v $(pwd):/ctrsploit --workdir /ctrsploit ubuntu 
root@7b500ab9a66b:/ctrsploit# ./checksec_linux_amd64 namespace cgroup
cgroup: host

❯ docker run --cgroupns=private -ti --rm -v $(pwd):/ctrsploit --workdir /ctrsploit ubuntu
root@0181d6dd6069:/ctrsploit# ./checksec_linux_amd64 namespace cgroup
cgroup: child
```

**ipc**

```
❯ docker run --ipc=host -ti --rm -v $(pwd):/ctrsploit --workdir /ctrsploit ubuntu
root@20726663d0c3:/ctrsploit# ./checksec_linux_amd64 namespace ipc
ipc: host

❯ docker run --ipc=private -ti --rm -v $(pwd):/ctrsploit --workdir /ctrsploit ubuntu
root@0b55497ab2ea:/ctrsploit# ./checksec_linux_amd64 namespace ipc
ipc: child
```

**mnt**

current process's mnt ns must be private ns

```
❯ docker run  -ti --rm -v $(pwd):/ctrsploit --workdir /ctrsploit ubuntu 
root@b5cae05e76bb:/ctrsploit# ./checksec_linux_amd64 namespace mnt
mnt: child
```

**net**

```
❯ docker run --net=host -ti --rm -v $(pwd):/ctrsploit --workdir /ctrsploit ubuntu
root@docker-desktop:/ctrsploit# ./checksec_linux_amd64 namespace net
net: host

❯ docker run -ti --rm -v $(pwd):/ctrsploit --workdir /ctrsploit ubuntu 
root@753f679cde36:/ctrsploit# ./checksec_linux_amd64 namespace net
net: child
```

**pid/pid_for_children**

```
❯ docker run --pid=host -ti --rm -v $(pwd):/ctrsploit --workdir /ctrsploit ubuntu
root@e78bf8cce178:/ctrsploit# ./checksec_linux_amd64 namespace pid
pid: host
root@e78bf8cce178:/ctrsploit# ./checksec_linux_amd64 namespace pid_for_children
pid_for_children: host

❯ docker run -ti --rm -v $(pwd):/ctrsploit --workdir /ctrsploit ubuntu 
root@cae17c751027:/ctrsploit# ./checksec_linux_amd64 namespace pid
pid: child
root@cae17c751027:/ctrsploit# ./checksec_linux_amd64 namespace pid_for_children
pid_for_children: child
```

**time/time_for_children**

docker has not supported time ns

```
❯ docker run  -ti --rm -v $(pwd):/ctrsploit --workdir /ctrsploit ubuntu 
root@d8fd83310e65:/ctrsploit# ./checksec_linux_amd64 namespace time
time: host
root@d8fd83310e65:/ctrsploit# ./checksec_linux_amd64 namespace time_for_children
time_for_children: host
```

**user**

userns-remap seems not work for mac, see https://github.com/docker/for-mac/issues/3280

```
❯ docker run --userns=host -ti --rm -v $(pwd):/ctrsploit --workdir /ctrsploit ubuntu
root@55f606353ec4:/ctrsploit# ./checksec_linux_amd64 namespace user
user: host
```

**uts**

```
❯ docker run --uts=host -ti --rm -v $(pwd):/ctrsploit --workdir /ctrsploit ubuntu
root@docker-desktop:/ctrsploit# ./checksec_linux_amd64 namespace uts
uts: host

❯ docker run -ti --rm -v $(pwd):/ctrsploit --workdir /ctrsploit ubuntu 
root@56fafd2a375c:/ctrsploit# ./checksec_linux_amd64 namespace uts
uts: child
```

### Test Environment 2: ubuntu 20.04

```
# uname -a
Linux wanglei-svc-gz 5.4.0-66-generic #74-Ubuntu SMP Wed Jan 27 22:54:38 UTC 2021 x86_64 x86_64 x86_64 GNU/Linux

# docker --version
Docker version 20.10.12, build e91ed57
```

**all namespaces level**

```
# docker run -ti --rm -v $(pwd):/ctrsploit --workdir /ctrsploit ubuntu
root@027304b6d035:/ctrsploit# ./checksec_linux_amd64 namespace
========namespace level=======
cgroup:              host
ipc:                 child
mnt:                 child
net:                 child
pid:                 child
pid_for_children:    child
user:                host
uts:                 child
```

**cgroups**

```
# docker run --cgroupns=host -ti --rm -v $(pwd):/ctrsploit --workdir /ctrsploit ubuntu
root@386b17079fb9:/ctrsploit# ./checksec_linux_amd64 namespace cgroup
cgroup: host

# docker run --cgroupns=private -ti --rm -v $(pwd):/ctrsploit --workdir /ctrsploit ubuntu
root@02079076bbe3:/ctrsploit# ./checksec_linux_amd64 namespace cgroup
cgroup: child
```

**ipc**

```
# docker run --ipc=host -ti --rm -v $(pwd):/ctrsploit --workdir /ctrsploit ubuntu
root@d4b5cc1a6f2a:/ctrsploit# ./checksec_linux_amd64 namespace ipc
ipc: host

# docker run --ipc=private -ti --rm -v $(pwd):/ctrsploit --workdir /ctrsploit ubuntu
root@1299bb43ca83:/ctrsploit# ./checksec_linux_amd64 namespace ipc
ipc: child
```

**mnt**

current process's mnt ns must be private ns

```
# docker run  -ti --rm -v $(pwd):/ctrsploit --workdir /ctrsploit ubuntu 
root@028f8b29263a:/ctrsploit# ./checksec_linux_amd64 namespace mnt
mnt: child
```

**net**

```
# docker run --net=host -ti --rm -v $(pwd):/ctrsploit --workdir /ctrsploit ubuntu
root@host:/ctrsploit# ./checksec_linux_amd64 namespace net
net: host

# docker run -ti --rm -v $(pwd):/ctrsploit --workdir /ctrsploit ubuntu 
root@0b6f2aec03f8:/ctrsploit# ./checksec_linux_amd64 namespace net
net: child
```

**pid/pid_for_children**

```
# docker run --pid=host -ti --rm -v $(pwd):/ctrsploit --workdir /ctrsploit ubuntu
root@0672356d0686:/ctrsploit# ./checksec_linux_amd64 namespace pid
pid: host
root@0672356d0686:/ctrsploit# ./checksec_linux_amd64 namespace pid_for_children
pid_for_children: host

# docker run -ti --rm -v $(pwd):/ctrsploit --workdir /ctrsploit ubuntu
root@1aa7f3550c4f:/ctrsploit# ./checksec_linux_amd64 namespace pid
pid: child
root@1aa7f3550c4f:/ctrsploit# ./checksec_linux_amd64 namespace pid_for_children
pid_for_children: child
```

**time/time_for_children**

* this version of kernel has not supported time ns
* this version of docker has not supported time ns

```
# docker run  -ti --rm -v $(pwd):/ctrsploit --workdir /ctrsploit ubuntu
root@4cbd7d848437:/ctrsploit# ./checksec_linux_amd64 namespace time
time: host
root@4cbd7d848437:/ctrsploit# ./checksec_linux_amd64 namespace time_for_children
time_for_children: host
```

**user**

```
❯ docker run --userns=host -ti --rm -v $(pwd):/ctrsploit --workdir /ctrsploit ubuntu
root@11edf641ed54:/ctrsploit# ./checksec_linux_amd64 namespace user
user: host

# mkdir -p /etc/docker
# cat << EOF > /etc/docker/daemon.json
{
    "userns-remap": "default"
}
EOF
# service docker restart
# docker run -tid ubuntu
0a1930c3d975f6846f24d32698b1bb520c4be7f72fe17f4f350a9a927e6115ce
# docker run -tid --name userns ubuntu
3e29e0755c65d26eaaadc49ed54d1c1d0b6b62bffac7fbbda8842b617b6a9c99
# docker cp checksec_linux_amd64 userns:/usr/bin/
# docker exec -ti userns bash
root@3e29e0755c65:/# checksec_linux_amd64 namespace user
user: child
```

**uts**

```
# docker run --uts=host -ti --rm -v $(pwd):/ctrsploit --workdir /ctrsploit ubuntu
root@wanglei-svc-gz:/ctrsploit# ./checksec_linux_amd64 namespace uts
uts: host

# docker run -ti --rm -v $(pwd):/ctrsploit --workdir /ctrsploit ubuntu
root@075b5ad696af:/ctrsploit# ./checksec_linux_amd64 namespace uts
uts: child
```


### json format

```
root@4fe779fc104c:~/ctrsploit/bin/release# ./env_linux_amd64 --json namespace | jq
{
  "name": {
    "name": "Namespace Level"
  },
  "levels": [
    {
      "name": "cgroup",
      "description": "",
      "result": "child"
    },
    {
      "name": "ipc",
      "description": "",
      "result": "child"
    },
    {
      "name": "mnt",
      "description": "",
      "result": "child"
    },
    {
      "name": "net",
      "description": "",
      "result": "child"
    },
    {
      "name": "pid",
      "description": "",
      "result": "child"
    },
    {
      "name": "pid_for_children",
      "description": "",
      "result": "child"
    },
    {
      "name": "time",
      "description": "",
      "result": "host"
    },
    {
      "name": "time_for_children",
      "description": "",
      "result": "host"
    },
    {
      "name": "user",
      "description": "",
      "result": "host"
    },
    {
      "name": "uts",
      "description": "",
      "result": "child"
    }
  ]
}
```