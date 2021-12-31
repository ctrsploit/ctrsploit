# env/cgroups

show the cgroups version, and list subsystems

查看cgroup版本, 及cgroup子系统

```
root@ctr:/# ./ctrsploit env c
INFO[0000] ===========Cgroups=========
is cgroupv1: ✘
is cgroupv2: ✔ 
```

```
root@ctr:/# ./ctrsploit env c
INFO[0000] ===========Cgroups=========
is cgroupv1: ✔
is cgroupv2: ✘ 
INFO[0000] ------sub systems-------
["blkio" "cpu,cpuacct" "cpuset" "devices" "files" "freezer" "hugetlb" "memory" "net_cls,net_prio" "perf_event" "pids" "rdma" "systemd"]
--------top level subsystem----------
["rdma"]
```