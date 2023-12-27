# guess version by syscall

## releases

* dockerhub/docker/dind
* runc/release

| docker(server)                    | runc                           | libseccomp | note | url                                                                                                                                                                                         |
|-----------------------------------|--------------------------------|------------|------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| 1.9.0-rc2_dind_static_60d36f7     | -                              |
| 19.03.14_dind_static_5eb3275      | 1.0.0-rc10_dind_static_dc9208a | 2.3.3      |
| 19.03.15_dind_static_99e3ed8      | 1.0.0-rc10_dind_static_dc9208a | 2.3.3      |
| 20.10.0-beta1_dind_static_9c15e82 | 1.0.0-rc92_dind_static_ff819c7 | 2.3.3      |
| 20.10.4_dind_static_363e9a8       | 1.0.0-rc93_dind_static_12644e6 | 2.3.3      |      | [docker:20.10.4-dind](https://hub.docker.com/layers/library/docker/20.10.4-dind/images/sha256-5a4bc48824d0a363b997211382ede81fa8a6a71ec92a454ae661cda668bd0415?context=explore)             |
| 20.10.23_dind_static_6051f14      | 1.1.4_dind_static_g5fd4c4d1    | 2.5.1      |      | [docker:20.10.23-dind](https://hub.docker.com/layers/library/docker/20.10.23-dind/images/sha256-88c01ecbf7def2c25f729b7a6052861b20c85913dfd9938f40bd2ab49c017951?context=explore)           |
| 20.10.24_dind_static_5d6db84      | 1.1.5_dind_static_gf19387a6    | 2.5.1      |      | [docker:20.10.24-dind](https://hub.docker.com/layers/library/docker/20.10.24-dind/images/sha256-11a8556b63283fb6edb98fc990166476bc2e14b33164039ab6132b27d84882d8?context=explore)           |
| 23.0.0-rc.4_dind_static_e92dd87   | 1.1.4_dind_static_g5fd4c4d     | 2.5.1      |      | [docker:23.0.0-rc.4-dind](https://hub.docker.com/layers/library/docker/23.0.0-rc.4-dind/images/sha256-140ae586d6981c3c768308ad4d7ffb85cc859551625e601b2d4caf830aaf1eb5?context=explore)     |
| 23.0.0_dind_static_d7573ab        | 1.1.4_dind_static_g5fd4c4d     | 2.5.1      |      | [docker:23.0.0-dind](https://hub.docker.com/layers/library/docker/23.0.0-dind/images/sha256-be7c1cb42809b910473a8a1b195736758fc8c10b395001b90968d5f31ad6a40b?context=explore)               |
| 23.0.1_dind_static_bc3805a        | 1.1.4_dind_static_g5fd4c4d     | 2.5.1      |      | [docker:23.0.1-dind](https://hub.docker.com/layers/library/docker/23.0.1-dind/images/sha256-021ad18301281a4367bc3f1e0ef998c1c54a7eba02d49b90892d01985904bb31?context=explore)               |
| 23.0.2_dind_static_219f21b        | 1.1.4_dind_static_g5fd4c4d     | 2.5.1      |      | [docker:23.0.2-dind](https://hub.docker.com/layers/library/docker/23.0.2-dind/images/sha256-3351a0a8ce728f938f8c15ebbacf533fd0054317544ad74348b5c651e0502a31?context=explore)               |
| 23.0.3_dind_static_59118bf        | 1.1.5_dind_static_gf19387a     | 2.5.1      |      | [docker:23.0.3-dind](https://hub.docker.com/layers/library/docker/23.0.3-dind/images/sha256-81b982fcfa4971a75286eb19c697c957bc29678ce0009e1b6a1f6e26ebb3ce6c?context=explore)               |
| 23.0.4_dind_static_cbce331        | 1.1.5_dind_static_gf19387a     | 2.5.1      |      | [docker:23.0.4-dind](https://hub.docker.com/layers/library/docker/23.0.4-dind/images/sha256-bbf1431e22956b26e333543422484882bf9fe5f2cfc236459c5945d5450c6e4f?context=explore)               |
| 23.0.5_dind_static_94d3ad6        | 1.1.5_dind_static_gf19387a     | 2.5.1      |      | [docker:23.0.5-dind](https://hub.docker.com/layers/library/docker/23.0.5-dind/images/sha256-280555387880b1fd76b45251d76122ea01ec00f4920373c46833564a0a7b24da?context=explore)               |
| 23.0.6_dind_static_9dbdbd4        | 1.1.7_dind_static_g860f061     | 2.5.1      |      | [docker:23.0.6-dind](https://hub.docker.com/layers/library/docker/23.0.6-dind/images/sha256-283f828961edf69aba8d784ccfdb9d2eacd7a718a67700302331c45b413eb9cb?context=explore)               |
| 24.0.0-beta.1_dind_static_348f836 | 1.1.5_dind_static_gf19387a     | 2.5.1      |      | [docker:24.0.0-beta.1-dind](https://hub.docker.com/layers/library/docker/24.0.0-beta.1-dind/images/sha256-9530f81960d06da36fa69cdf448e20e9df1392b5cbf334036f35d505d9314d14?context=explore) |
| 24.0.0-beta.2_dind_static_5b1282c | 1.1.5_dind_static_gf19387a     | 2.5.1      |      | [docker:24.0.0-beta.2-dind](https://hub.docker.com/layers/library/docker/24.0.0-beta.2-dind/images/sha256-de6af7ab610f6ae7c9c30123df9873895b9f969a7bd9b1c48578d961cb3a5b9d?context=explore) |
| 24.0.0-rc.1_dind_static_f117aef   | 1.1.6_dind_static_g0f48801     | 2.5.1      |      | [docker:24.0.0-rc.1-dind](https://hub.docker.com/layers/library/docker/24.0.0-rc.1-dind/images/sha256-61e10fd22216fcb93000fdf579b37be28224d24f905c85a83f94771dbaf42c63?context=explore)     |
| 24.0.0-rc.2_dind_static_8d9a40a   | 1.1.7_dind_static_g860f061     | 2.5.1      |      | [docker:24.0.0-rc.2-dind](https://hub.docker.com/layers/library/docker/24.0.0-rc.2-dind/images/sha256-81915f22c7da6419ff3be8f990cc533516030a4ee4748234070effce4335e954?context=explore)     |
| 24.0.0-rc.3_dind_static_807e415   | 1.1.7_dind_static_g860f061     | 2.5.1      |      | [docker:24.0.0-rc.3-dind](https://hub.docker.com/layers/library/docker/24.0.0-rc.3-dind/images/sha256-c09bc1b4d6881dc35cce5458c212e7f030d380777d15e74d726269d19a86e57a?context=explore)     |
| 24.0.0-rc.4_dind_static_a5b597e   | 1.1.7_dind_static_g860f061     | 2.5.1      |      | [docker:24.0.0-rc.4-dind](https://hub.docker.com/layers/library/docker/24.0.0-rc.4-dind/images/sha256-c286e59386b4cd75cd5a8f0e95785974a2c70fcb077db994f8e3367b97da8d85?context=explore)     |
| 24.0.0_dind_static_1331b8c        | 1.1.7_dind_static_g860f061     | 2.5.1      |      | [docker:24.0.0-dind](https://hub.docker.com/layers/library/docker/24.0.0-dind/images/sha256-09b5af3bc4554a5138807ab7bb6b9560de69ae28c8834687f114ac41a8d5d31f?context=explore)               |
| 24.0.1_dind_static_463850e        | 1.1.7_dind_static_g860f061     | 2.5.1      |      | [docker:24.0.1-dind](https://hub.docker.com/layers/library/docker/24.0.1-dind/images/sha256-22679c5c3a1967133f1e817c3ba187024375205d1ef084be7aecf6b0d4599c99?context=explore)               |
| 24.0.2_dind_static_659604f        | 1.1.7_dind_static_g860f061     | 2.5.1      |      | [docker:24.0.2-dind](https://hub.docker.com/layers/library/docker/24.0.2-dind/images/sha256-d981b86555d31fb68b974505d0815223c027362a5802ad3ac394a0dabb0da658?context=explore)               |
| 24.0.3_dind_static_1d9c8619       | 1.1.7_dind_static_g860f061     | 2.5.1      |      | [docker:24.0.3-dind](https://hub.docker.com/layers/library/docker/24.0.3-dind/images/sha256-cd573a0e55aa7ef9f2b767816baabd96f1b7e0edb4e3409ed14b2247a75e978d?context=explore)               |
| 24.0.4_dind_static_4ffc614        | 1.1.7-dind_static_g860f061     | 2.5.1      |      | [docker:24.0.4-dind](https://hub.docker.com/layers/library/docker/24.0.4-dind/images/sha256-adb52cf47859063409c1826d5c09b633d4d8e4f960f43e928b128caa7223c23e?context=explore)               |
| 24.0.5_dind_static_a61e2b4        | 1.1.8_dind_static_g82f18fe     | 2.5.1      |      | [docker:24.0.5-dind](https://hub.docker.com/layers/library/docker/24.0.5-dind/images/sha256-78e33a72a1df68d3b4341630a215626e79e45725e205f0355dc05df1b81ad68a?context=explore)               |
| 24.0.6_dind_static_1a79695        | 1.1.9_dind_static_gccaecfc     | 2.5.1      |      | [docker:24.0.6-dind](https://hub.docker.com/layers/library/docker/24.0.6-dind/images/sha256-74e78208fc18da48ddf8b569abe21563730845c312130bd0f0b059746a7e10f5?context=explore)               |
| 24.0.7_dind_static_311b9ff        | 1.1.9_dind_static_gccaecfc     | 2.5.1      |      | [docker:24.0.7-dind](https://hub.docker.com/layers/library/docker/24.0.7-dind/images/sha256-4c92bd9328191f76e8eec6592ceb2e248aa7406dfc9505870812cf8ebee9326a?context=explore)               |
| 25.0.0-beta.1_dind_static_6af7d6e | 1.1.10_dind_static_g18a0cb0    | 2.5.1      |      | [docker:25.0.0-beta.1-dind](https://hub.docker.com/layers/library/docker/25.0.0-beta.1-dind/images/sha256-0a71a25e703d5ec1c5a9460aa0bc80cac80910eb08d5abab3db5fb4490091ffb?context=explore) |

| runc                           | libseccomp | note | url                                                                                            |
|--------------------------------|------------|------|------------------------------------------------------------------------------------------------|
| 1.0.0-rc92_runc-release_static | 2.4.1      |      | [v1.0.0-rc92](https://github.com/opencontainers/runc/releases/download/v1.0.0-rc92/runc.amd64) |
| 1.0.0-rc93_runc-release_static | 2.5.1      |      | [v1.0.0-rc93](https://github.com/opencontainers/runc/releases/download/v1.0.0-rc93/runc.amd64) | 

## dynamic/static

dynamic

```
$ readelf -d /usr/bin/runc |grep Shared
 0x0000000000000001 (NEEDED)             Shared library: [libpthread.so.0]
 0x0000000000000001 (NEEDED)             Shared library: [libseccomp.so.2]
 0x0000000000000001 (NEEDED)             Shared library: [libc.so.6]
```

static

```
$ readelf -d runc
```

## libseccomp

get the libseccomp version of static runc binary

https://man7.org/linux/man-pages/man3/seccomp_version.3.html

```
$ wget https://hub.gitmirror.com/https://github.com/opencontainers/runc/releases/download/v1.0.0-rc92/runc.amd64
$ r2 -c "pf iii @ obj.library_version" runc.amd64
0x00b42090 = 2
0x00b42094 = 4
0x00b42098 = 1
```