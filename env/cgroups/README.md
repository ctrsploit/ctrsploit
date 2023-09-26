# env/cgroups

show the cgroups version, and list subsystems

查看cgroup版本, 及cgroup子系统

```
root@b93597d64108:/# ./env_linux_amd64 cgroups 
===========CGroups===========
[Y]  v1 
[N]  v2 
sub systems     
["rdma" "freezer" "devices" "perf_event" "pids" "hugetlb" "cpu" "cpuacct" "net_cls" "net_prio" "blkio" "cpuset" "memory"]
top level subsystems    
[]
root@b93597d64108:/# ./env_linux_amd64 --colorful cgroups
===========CGroups===========
✔  v1   
✘  v2   
sub systems     
["hugetlb" "devices" "freezer" "rdma" "cpuacct" "cpuset" "pids" "net_prio" "net_cls" "blkio" "cpu" "memory" "perf_event"]
top level subsystems    
[]
```