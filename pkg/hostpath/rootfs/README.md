# infer host path of container's rootfs

## docker

### overlay

rootfs=[upperdir]/../merged

### devicemapper

rootfs=/var/lib/docker/devicemapper/mnt/[dm]/

## containerd

### overlay

ctr:
upperdir=/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/17/fs/
rootfs=/run/containerd/io.containerd.runtime.v2.task/default/[container-id]/rootfs/

nerdctl:
upperdir=/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/15/fs/
rootfs=/run/containerd/io.containerd.runtime.v2.task/default/[container-id]/rootfs/

k8s+containerd:
rootfs=/run/containerd/io.containerd.runtime.v2.task/k8s.io/[container-id]/rootfs/ssst0n3-eee
