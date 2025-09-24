# runtime type

## docker

1. /.dockerenv exists
2. rootfs's mountinfo contains docker
3. cgroup v1: /proc/self/cgroup contains docker
4. /etc/hosts's mountinfo contains docker
5. hostname matches regex pattern `^[0-9a-f]{12}$`

