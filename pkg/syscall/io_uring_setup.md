# io_uring_setup

## commits

* CommitWhitelistIoUring https://github.com/moby/moby/commit/f4d41f1dfa52caa8f12b070315e230e7eded5f4a
* CommitBlockIoUring https://github.com/moby/moby/commit/891241e7e74d4aae6de5f6125574eb994f25e169

## testcase

| kernel | docker        | runc                        | libseccomp | errno  | note                                    |
|--------|---------------|-----------------------------|------------|--------|-----------------------------------------|
| 4.4.0  | 19.03.14      | static-1.0.0-rc10           |            | EPERM  |
| 4.4.0  | 19.03.15      | static-1.0.0-rc10           |            | EPERM  |
| 4.4.0  | 20.10.0-beta1 | static-1.0.0-rc91           |            | EPERM  |
| 4.4.0  | 20.10.0-beta1 | static-1.0.0-rc92           |            | EPERM  |
| 4.4.0  | 20.10.3       | static-1.0.0-rc92           |            | EPERM  |
| 5.4.0  | 20.10.3       | static-1.0.0-rc92           |            | EPERM  |
| 5.4.0  | 20.10.4       | static-1.0.0-rc93           | 2.3.3      | ENOSYS |
| 5.4.0  | 20.10.5       | static-1.0.0-rc93 (12644e6) | 2.3.3      | ENOSYS |
| 5.4.0  | 20.10.6       | static-1.0.0-rc93 (12644e6) | 2.4.4      | EFAULT | https://github.com/moby/moby/pull/42147 |
| 5.4.0  | 20.10.7       | static-1.0.0-rc95 (b9ee9c6) | 2.4.4      | EFAULT |                                         |

| kernel | docker                            | runc                           | libseccomp | errno  | note |
|--------|-----------------------------------|--------------------------------|------------|--------|------|
| 5.4.0  | 19.03.14_dind_static_5eb3275      | 1.0.0-rc10_dind_static_dc9208a | 2.3.3      | EPERM  |
| 5.4.0  | 19.03.15_dind_static_99e3ed8      | 1.0.0-rc10_dind_static_dc9208a | 2.3.3      | EPERM  |
| 5.4.0  | 20.10.0-beta1_dind_static_9c15e82 | 1.0.0-rc92_dind_static_ff819c7 | 2.3.3      | EPERM  |
| 5.4.0  | 20.10.0-beta1_dind_static_9c15e82 | 1.0.0-rc92_runc-release_static | 2.4.1      | EPERM  |
| 5.4.0  | 20.10.4_dind_static_363e9a8       | 1.0.0-rc93_dind_static_12644e6 | 2.3.3      | ENOSYS |      |
| 5.4.0  | 20.10.4_dind_static_363e9a8       | 1.0.0-rc93_runc-release_static | 2.5.1      | EFAULT |      |

### kernel-4.4.0_docker-20.10.3_runc-1.0.0-rc92

| software | version    | errno  | note          |
|----------|------------|--------|---------------|
| kernel   | 4.4.0      | ENOSYS |
| docker   | 20.10.3    | EFAULT |
| runc     | 1.0.0-rc92 | EPERM  | ENOSYS->EPERM |
| final    | -          | EPERM  |

```
root@ubuntu:~/ctrsploit/bin/release# docker run  --privileged --name dind -d -v $(pwd):/ctrsploit docker:20.10.3-dind
root@ubuntu:~/ctrsploit/bin/release# docker exec -ti dind sh
/ # docker version
Client: Docker Engine - Community
 Version:           20.10.3
 API version:       1.41
 Go version:        go1.13.15
 Git commit:        48d30b5
 Built:             Fri Jan 29 14:28:23 2021
 OS/Arch:           linux/amd64
 Context:           default
 Experimental:      true

Server: Docker Engine - Community
 Engine:
  Version:          20.10.3
  API version:      1.41 (minimum version 1.12)
  Go version:       go1.13.15
  Git commit:       46229ca
  Built:            Fri Jan 29 14:31:57 2021
  OS/Arch:          linux/amd64
  Experimental:     false
 containerd:
  Version:          v1.4.3
  GitCommit:        269548fa27e0089a8b8278fc4fc781d7f65a939b
 runc:
  Version:          1.0.0-rc92
  GitCommit:        ff819c7e9184c13b7c2607fe6c30ae19403a7aff
 docker-init:
  Version:          0.19.0
  GitCommit:        de40ad0
/ # uname -a
Linux ae6772f01057 4.4.0-210-generic #242-Ubuntu SMP Fri Apr 16 09:57:56 UTC 2021 x86_64 Linux
/ # docker run -ti -v /ctrsploit:/ctrsploit busybox sh
/ # /ctrsploit/env_linux_amd64 --debug dv
DEBU[0000]/go/pkg/mod/github.com/ctrsploit/sploit-spec@v0.4.3/pkg/app/flags.go:60 github.com/ctrsploit/sploit-spec/pkg/app.InstallGlobalFlagDebug.func1() debug mode on                                
DEBU[0000]/root/ctrsploit/pkg/seccomp/syscall.go:58 github.com/ctrsploit/ctrsploit/pkg/seccomp.Syscall.State() syscall 425 errno: operation not permitted   
DEBU[0000]/root/ctrsploit/env/version/docker.go:27 github.com/ctrsploit/ctrsploit/env/version.Docker() io_uring_setup false                         
dockerd-version	# dockerd version range
dockerd is in [0.1.0, 20.10.0-beta1) U [25.0.0-beta.1, 99.99.99)
```

### kernel-4.4.0_docker-20.10.4_runc-1.0.0-rc93

| software | version    | errno  | note   |
|----------|------------|--------|--------|
| kernel   | 4.4.0      | ENOSYS |
| docker   | 20.10.4    | EFAULT |
| runc     | 1.0.0-rc93 | ENOSYS | static |
| final    | -          | ENOSYS |

```
root@ubuntu:~/ctrsploit/bin/release# docker run  --privileged --name dind -d -v $(pwd):/ctrsploit docker:20.10.4-dind
ebc921057269210237bdca941c0b06ac1fd98140187565c7f2f79beff9fc643b
root@ubuntu:~/ctrsploit/bin/release# docker exec -ti dind sh
/ # docker version
Client: Docker Engine - Community
 Version:           20.10.4
 API version:       1.41
 Go version:        go1.13.15
 Git commit:        d3cb89e
 Built:             Thu Feb 25 07:01:39 2021
 OS/Arch:           linux/amd64
 Context:           default
 Experimental:      true

Server: Docker Engine - Community
 Engine:
  Version:          20.10.4
  API version:      1.41 (minimum version 1.12)
  Go version:       go1.13.15
  Git commit:       363e9a8
  Built:            Thu Feb 25 07:05:55 2021
  OS/Arch:          linux/amd64
  Experimental:     false
 containerd:
  Version:          v1.4.3
  GitCommit:        269548fa27e0089a8b8278fc4fc781d7f65a939b
 runc:
  Version:          1.0.0-rc93
  GitCommit:        12644e614e25b05da6fd08a38ffa0cfe1903fdec
 docker-init:
  Version:          0.19.0
  GitCommit:        de40ad0
/ # uname -a
Linux ebc921057269 4.4.0-210-generic #242-Ubuntu SMP Fri Apr 16 09:57:56 UTC 2021 x86_64 Linux
/ # docker run -ti -v /ctrsploit:/ctrsploit busybox sh
Unable to find image 'busybox:latest' locally
latest: Pulling from library/busybox
3f4d90098f5b: Pull complete 
Digest: sha256:3fbc632167424a6d997e74f52b878d7cc478225cffac6bc977eedfe51c7f4e79
Status: Downloaded newer image for busybox:latest
/ # 
/ # /ctrsploit/env_linux_amd64 --debug dv
DEBU[0000]/go/pkg/mod/github.com/ctrsploit/sploit-spec@v0.4.3/pkg/app/flags.go:60 github.com/ctrsploit/sploit-spec/pkg/app.InstallGlobalFlagDebug.func1() debug mode on                                
DEBU[0000]/root/ctrsploit/pkg/seccomp/syscall.go:58 github.com/ctrsploit/ctrsploit/pkg/seccomp.Syscall.State() syscall 425 errno: function not implemented  
DEBU[0000]/root/ctrsploit/env/version/docker.go:27 github.com/ctrsploit/ctrsploit/env/version.Docker() io_uring_setup true                          
dockerd-version	# dockerd version range
dockerd is in [20.10.0-beta1, 25.0.0-beta.1)
```

### kernel-4.4.0_docker-24.0.7_runc-1.0.0-rc93

| software | version | errno  | note |
|----------|---------|--------|------|
| kernel   | 4.4.0   | ENOSYS |
| docker   | 24.0.7  | EFAULT |
| runc     | 1.1.9   | ENOSYS |      |
| final    | -       | ENOSYS |

```
root@ubuntu:~/ctrsploit/bin/release# docker run  --privileged --name dind -d -v $(pwd):/ctrsploit docker:24.0.7-dind
20a5836a9bf574cc6ba0486352c7dc019dda11b06763f4e746f07e9d45bb09b5
root@ubuntu:~/ctrsploit/bin/release# docker exec -ti dind sh
/ # docker version
Client:
 Version:           24.0.7
 API version:       1.43
 Go version:        go1.20.10
 Git commit:        afdd53b
 Built:             Thu Oct 26 09:04:00 2023
 OS/Arch:           linux/amd64
 Context:           default

Server: Docker Engine - Community
 Engine:
  Version:          24.0.7
  API version:      1.43 (minimum version 1.12)
  Go version:       go1.20.10
  Git commit:       311b9ff
  Built:            Thu Oct 26 09:05:28 2023
  OS/Arch:          linux/amd64
  Experimental:     false
 containerd:
  Version:          v1.7.6
  GitCommit:        091922f03c2762540fd057fba91260237ff86acb
 runc:
  Version:          1.1.9
  GitCommit:        v1.1.9-0-gccaecfc
 docker-init:
  Version:          0.19.0
  GitCommit:        de40ad0
/ # uname -a
Linux 20a5836a9bf5 4.4.0-210-generic #242-Ubuntu SMP Fri Apr 16 09:57:56 UTC 2021 x86_64 Linux
/ # docker run -ti -v /ctrsploit:/ctrsploit busybox sh
Unable to find image 'busybox:latest' locally
latest: Pulling from library/busybox
3f4d90098f5b: Pull complete 
Digest: sha256:3fbc632167424a6d997e74f52b878d7cc478225cffac6bc977eedfe51c7f4e79
Status: Downloaded newer image for busybox:latest
/ # /ctrsploit/env_linux_amd64 --debug dv
DEBU[0000]/go/pkg/mod/github.com/ctrsploit/sploit-spec@v0.4.3/pkg/app/flags.go:60 github.com/ctrsploit/sploit-spec/pkg/app.InstallGlobalFlagDebug.func1() debug mode on                                
DEBU[0000]/root/ctrsploit/pkg/seccomp/syscall.go:58 github.com/ctrsploit/ctrsploit/pkg/seccomp.Syscall.State() syscall 425 errno: function not implemented  
DEBU[0000]/root/ctrsploit/env/version/docker.go:27 github.com/ctrsploit/ctrsploit/env/version.Docker() io_uring_setup true                          
dockerd-version	# dockerd version range
dockerd is in [20.10.0-beta1, 25.0.0-beta.1)
```