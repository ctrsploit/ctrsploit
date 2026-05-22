---
maintainer:
    - Hpd0ger
---
# env/services

Discover Kubernetes services visible from the current pod. The module combines
Kubernetes-injected environment variables with DNS-based discovery to list
service names, namespaces, cluster IPs, ports, pod IPs when available, and the
source method that found each record.

## Usage

```shell
$ ctrsploit env services
$ ctrsploit env services --methods env
$ ctrsploit env services --methods wildcard,axfr,cidr --zone cluster.local --cidr 10.96.0.0/12 --threads 32
$ ctrsploit env services --output /tmp/services.json
```

## Discovery Methods

- `env`: parse service variables such as `*_SERVICE_HOST`, `*_SERVICE_PORT`,
  and `*_PORT_<port>_TCP_*` from the current process environment.
- `wildcard`: query Kubernetes DNS wildcard SRV records under
  `svc.<zone>`.
- `axfr`: attempt a DNS zone transfer for the configured cluster zone.
- `cidr`: scan a service CIDR with PTR and SRV lookups to discover services.
- `all`: run all supported methods.

DNS-based methods first verify that Kubernetes DNS is reachable. If the cluster
DNS check fails, only environment-variable discovery can return results.

## Options

- `--zone, -z`: Kubernetes DNS zone. Defaults to `cluster.local`.
- `--cidr, -c`: service CIDR for DNS reverse scanning. When omitted, ctrsploit
  derives a `/16` from `KUBERNETES_SERVICE_HOST` if that variable exists.
- `--threads, -t`: worker count for CIDR scanning. Defaults to `16`.
- `--methods, -m`: comma-separated discovery methods. Defaults to `all`.
- `--output, -o`: write discovered services as newline-delimited JSON.

## Output

Human output is printed through the standard ctrsploit printer. Machine output
contains records like:

```json
{
  "name": "kubernetes",
  "namespace": "default",
  "cluster_ip": "10.96.0.1",
  "ports": [
    {
      "port": 443,
      "protocol": "tcp"
    }
  ],
  "sources": [
    "env"
  ]
}
```
