---

tags: sploit
author:
    - Hpd0ger
version: v0.1.1
changelog:
    - v0.1.1: add reproduce steps
    - v0.1.0: init

---

# Service Account Token Dangerous Permissions

## 1. Vulnerability Overview

This module checks if the current Pod's ServiceAccount token has dangerous RBAC permissions that could lead to privilege escalation, remote code execution (RCE), or sensitive data access.

Unlike the `access-secrets` module which requires cluster admin privileges to scan the entire cluster, this module only needs the current Pod's ServiceAccount token and uses the Kubernetes `SelfSubjectAccessReview` API to check permissions.

## 2. Dangerous Permissions

The module checks for the following dangerous permissions (defined in YAML configuration):

### Critical Level

| Resource | Subresource | Verbs | Risk |
|----------|-------------|-------|------|
| `nodes` | `proxy` | `get` | RCE via WebSocket bypass ([Reference](https://grahamhelton.com/blog/nodes-proxy-rce)) |
| `pods` | `exec` | `create` | Command execution in other Pods |
| `pods` | `attach` | `create` | Attach to processes in other Pods |
| `*` | - | `*` | Full cluster control |
| `clusterroles` | - | `create`, `update`, `patch` | Privilege escalation |
| `clusterrolebindings` | - | `create`, `update`, `patch` | Bind high-privilege roles |

### High Level

| Resource | Subresource | Verbs | Risk |
|----------|-------------|-------|------|
| `pods` | - | `create` | Create privileged containers |
| `secrets` | - | `get`, `list` | Read sensitive credentials |
| `serviceaccounts` | `token` | `create` | Create tokens for any SA |
| `configmaps` | - | `get`, `list` | Read sensitive configuration |
| `roles` | - | `create`, `update`, `patch` | Namespace-level privilege escalation |
| `rolebindings` | - | `create`, `update`, `patch` | Bind roles in namespace |
| `*` | - | `create` | Create any resource |

### Medium Level

| Resource | Subresource | Verbs | Risk |
|----------|-------------|-------|------|
| `pods` | `log` | `get` | Read container logs |
| `pods` | `portforward` | `create` | Network access to Pods |
| `nodes` | - | `get`, `list` | Read node information |
| `persistentvolumes` | - | `create` | Access host storage |

## 3. Prerequisites

* Running inside a Kubernetes Pod with a mounted ServiceAccount token
* Or access to a kubeconfig file

## 4. Reproduce

### 4.1 Reproduce Environment

```shell
$ git clone https://github.com/ssst0n3/docker_archive.git
$ cd docker_archive/kubernetes/v1.34.0-calico
$ docker compose -f docker-compose.yml -f docker-compose.kvm.yml up -d
```

### 4.2 Reproduce Steps

```shell
$ ./ssh
root@kubernetes-1-34-0:~# wget -q https://github.com/ctrsploit/ctrsploit/releases/latest/download/ctrsploit_linux_amd64 -O /usr/bin/ctrsploit
root@kubernetes-1-34-0:~# chmod +x /usr/bin/ctrsploit
root@kubernetes-1-34-0:~# ctrsploit checksec sa-token-policy
INFO[0001]                                              
INFO[0001] === ServiceAccount Token Dangerous Permissions === 
ERRO[0001]                                              
ERRO[0001] [CRITICAL] 6 critical permissions found:     
INFO[0001]   - nodes/proxy [get] Cluster-Wide           
INFO[0001]     Risk: nodes/proxy GET allows RCE in any Pod via WebSocket bypass 
INFO[0001]     Ref: https://grahamhelton.com/blog/nodes-proxy-rce 
INFO[0001]   - pods/exec [create] Cluster-Wide          
INFO[0001]     Risk: pods/exec allows command execution in other Pods 
INFO[0001]   - pods/attach [create] Cluster-Wide        
INFO[0001]     Risk: pods/attach allows attaching to processes in other Pods 
INFO[0001]   - * [*] Cluster-Wide                       
INFO[0001]     Risk: Wildcard permission grants full cluster control 
INFO[0001]   - clusterroles [create] Cluster-Wide       
INFO[0001]     Risk: Can create/modify ClusterRoles to escalate privileges 
INFO[0001]   - clusterrolebindings [create] Cluster-Wide 
INFO[0001]     Risk: Can bind high-privilege roles to any ServiceAccount 
WARN[0001]                                              
WARN[0001] [HIGH] 7 high-risk permissions found:        
INFO[0001]   - pods [create] Cluster-Wide               
INFO[0001]     Risk: Can create privileged containers for container escape 
INFO[0001]   - secrets [get] Cluster-Wide               
INFO[0001]     Risk: Can read sensitive credentials and tokens 
INFO[0001]   - serviceaccounts/token [create] Cluster-Wide 
INFO[0001]     Risk: Can create tokens for any ServiceAccount 
INFO[0001]   - configmaps [get] Cluster-Wide            
INFO[0001]     Risk: ConfigMaps may contain sensitive configuration data 
INFO[0001]   - roles [create] Cluster-Wide              
INFO[0001]     Risk: Can create/modify Roles to escalate privileges in namespace 
INFO[0001]   - rolebindings [create] Cluster-Wide       
INFO[0001]     Risk: Can bind roles to ServiceAccounts in namespace 
INFO[0001]   - * [create] Cluster-Wide                  
INFO[0001]     Risk: Can create any resource type       
INFO[0001]                                              
INFO[0001] [MEDIUM] 4 medium-risk permissions found:    
INFO[0001]   - pods/log [get] Cluster-Wide              
INFO[0001]     Risk: Can read container logs which may contain sensitive data 
INFO[0001]   - pods/portforward [create] Cluster-Wide   
INFO[0001]     Risk: Can forward ports to access Pod network services 
INFO[0001]   - nodes [get] Cluster-Wide                 
INFO[0001]     Risk: Can read node information and metadata 
INFO[0001]   - persistentvolumes [create] Cluster-Wide  
INFO[0001]     Risk: Can create PVs potentially accessing host storage 
INFO[0001]                                              
INFO[0001] === Summary ===                              
INFO[0001] Critical: 6, High: 7, Medium: 4              


[Y]  sa-token-policy	# Check if service account token has dangerous permissions
```

## 5. Advanced Usage

### 5.1 Options

| Flag | Alias | Description | Default |
|------|-------|-------------|---------|
| `--namespace` | `-n` | Namespace to check | cluster-wide |
| `--level` | `-l` | Minimum level to report: critical, high, medium | medium |
| `--config` | `-c` | Path to custom dangerous permissions YAML file | - |

### 5.2 Custom Permissions Configuration

You can create a custom YAML file to define additional dangerous permissions:

```yaml
# custom_permissions.yaml
permissions:
  - resource: my-custom-resource
    group: example.com
    verbs: [create, delete]
    level: high
    description: "Custom resource with dangerous permissions"
    reference: "https://example.com/docs"
```

Use with `--config` flag:

```shell
ctrsploit checksec sa-token-policy --config custom_permissions.yaml
```

## 6. References

* [Kubernetes RBAC Good Practices](https://kubernetes.io/docs/concepts/security/rbac-good-practices/)
* [nodes/proxy RCE via WebSocket](https://grahamhelton.com/blog/nodes-proxy-rce)
* [Privilege Escalation in Kubernetes RBAC](https://kubernetes.io/docs/reference/access-authn-authz/authorization/#privilege-escalation-prevention)
