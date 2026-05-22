---
maintainer:
    - ssst0n3
---
# env/cgroups

show the cgroups version, and list subsystems

查看cgroup版本, 及cgroup子系统

```
root@76459c21a665:/# ./env_linux_amd64 --colorful cgroups
===========CGroups===========
✔  v1	
✘  v2	
sub systems	
["blkio" "perf_event" "memory" "rdma" "devices" "net_cls" "cpuset" "cpu" "cpuacct" "hugetlb" "pids" "freezer" "net_prio"]
top level subsystems	
[]

root@76459c21a665:/# ./env_linux_amd64 cgroups
===========CGroups===========
[Y]  v1	
[N]  v2	
sub systems	
["hugetlb" "blkio" "perf_event" "net_cls" "devices" "cpu" "net_prio" "memory" "freezer" "cpuset" "cpuacct" "pids" "rdma"]
top level subsystems	
[]

root@76459c21a665:/# ./env_linux_amd64 --json cgroups
{"version":1,"sub":["memory","hugetlb","perf_event","net_cls","pids","freezer","rdma","devices","blkio","net_prio","cpuset","cpu","cpuacct"],"top":[]}
```