# runtime type

## docker

1. /.dockerenv exists
2. rootfs's mountinfo contains docker
3. cgroup v1: /proc/1/cgroup contains docker
4. /etc/hosts's mountinfo contains docker
5. apparmor: /proc/1/attr/current contains docker
6. ~~hostname matches regex pattern `^[0-9a-f]{12}$`~~: both docker, nerdctl match this behavior
7. /proc/net/unix contains docker.sock (prerequisite: --net=host)

## containerd
1. rootfs's mountinfo contains containerd
2. /etc/hosts's mountinfo contains nerdctl
3. /etc/hostname's mountinfo contains containerd/nerdctl
4. /proc/net/unix contains containerd.sock, no docker.sock (prerequisite: --net=host)
