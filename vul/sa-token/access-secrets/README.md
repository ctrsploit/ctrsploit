---
maintainer:
    - ssst0n3
---
# Service Account Token Access to Secrets

## 1. Vulnerability Overview

This vulnerability check identifies Pods in a Kubernetes cluster that have access to Kubernetes Secrets through their ServiceAccount tokens. If an attacker gains access to a Pod, they can extract the ServiceAccount token from `/var/run/secrets/kubernetes.io/serviceaccount/token` and use it to access secrets that the ServiceAccount has permissions to read via the Kubernetes API.

## 2. Exploit Scenario

Insecure RBAC configuration

## 3. Prerequisites

* ServiceAccounts with RBAC permissions to access Secrets (via RoleBindings or ClusterRoleBindings)
* Pods using these ServiceAccounts
* Permissions include verbs: `get`, `list`, `watch`, or `*` on the `secrets` resource

**Note**: This module is only used for detection purposes. It assumes you have the highest permissions in the cluster to read RBAC configurations, ServiceAccounts, and Pods across all namespaces.

## 4. Vulnerability Existence Check

```shell
ctrsploit vul sa-token-access-secrets checksec
```

Or use the alias:

```shell
ctrsploit vul secret checksec
```

## 5. Reproduce

![](./video.svg)

### 5.1 Reproduce Environment

```shell
$ git clone https://github.com/ssst0n3/docker_archive.git
$ cd docker_archive/kubernetes/v1.34.0-calico
$ docker compose -f docker-compose.yml -f docker-compose.kvm.yml up -d
```

### 5.2 Reproduce Steps

```shell
$ ./ssh
root@kubernetes-1-34-0:~# wget -q https://github.com/ctrsploit/ctrsploit/releases/latest/download/ctrsploit_linux_amd64 -O /usr/bin/ctrsploit
root@kubernetes-1-34-0:~# chmod +x /usr/bin/ctrsploit
root@kubernetes-1-34-0:~# ctrsploit vul secret c
INFO[0001]                                              
INFO[0001] === Cluster-wide Pods with Secret Access === 
INFO[0001] Found 1 pod(s) with secret access permissions across the cluster 
INFO[0001]                                              
INFO[0001] === Pods by Namespace ===                    
INFO[0001]                                              
INFO[0001] Namespace: tigera-operator                   
INFO[0001]   Pods: 1                                    
INFO[0001]   Service Accounts: 1                        
INFO[0001]   RBAC Permissions:                          
INFO[0001]     - [ClusterRoleBinding] tigera-operator (Scope: Cluster-Wide, Verbs: list,get,watch) 
INFO[0001]                                              
INFO[0001] === Summary ===                              
INFO[0001] Total pods with secret access: 1             
INFO[0001] Namespaces affected: 1                       
INFO[0001] Unique RBAC permissions: 1                   
[Y]  sa-token-access-secrets	# Check if service account token can access Kubernetes Secrets
```
