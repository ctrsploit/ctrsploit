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

## 4. Usage

### Basic Usage

```shell
# Check for dangerous permissions
ctrsploit checksec sa-token-policy

# Or use alias
ctrsploit checksec policy
ctrsploit checksec dp
```

### Command Options

```shell
# Check specific namespace only
ctrsploit checksec sa-token-policy --namespace kube-system

# Only report critical level permissions
ctrsploit checksec sa-token-policy --level critical

# Use custom permissions YAML file
ctrsploit checksec sa-token-policy --config /path/to/custom.yaml
```

### Options

| Flag | Alias | Description | Default |
|------|-------|-------------|---------|
| `--namespace` | `-n` | Namespace to check | cluster-wide |
| `--level` | `-l` | Minimum level to report: critical, high, medium | medium |
| `--config` | `-c` | Path to custom dangerous permissions YAML file | - |

## 5. Output Example

```
INFO === ServiceAccount Token Dangerous Permissions ===

ERRO [CRITICAL] 2 critical permissions found:
INFO   - nodes/proxy [get] Cluster-Wide
INFO     Risk: nodes/proxy GET allows RCE in any Pod via WebSocket bypass
INFO     Ref: https://grahamhelton.com/blog/nodes-proxy-rce
INFO   - pods/exec [create] Namespace: default
INFO     Risk: pods/exec allows command execution in other Pods

WARN [HIGH] 1 high-risk permission found:
INFO   - secrets [get] Namespace: kube-system
INFO     Risk: Can read sensitive credentials and tokens

INFO === Summary ===
INFO Critical: 2, High: 1, Medium: 0
[Y]  sa-token-policy    # Check if service account token has dangerous permissions
```

## 6. Custom Permissions Configuration

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

## 7. References

* [Kubernetes RBAC Good Practices](https://kubernetes.io/docs/concepts/security/rbac-good-practices/)
* [nodes/proxy RCE via WebSocket](https://grahamhelton.com/blog/nodes-proxy-rce)
* [Privilege Escalation in Kubernetes RBAC](https://kubernetes.io/docs/reference/access-authn-authz/authorization/#privilege-escalation-prevention)
