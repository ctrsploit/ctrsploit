package services

import (
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type svcInfo struct {
	host  string
	ports map[uint16]ServicePort
}

// DiscoverFromEnv discovers K8s services by parsing environment variables.
// Kubernetes injects env vars like:
//
//	REDIS_SERVICE_HOST=10.0.0.11
//	REDIS_SERVICE_PORT=6379
//	REDIS_PORT_6379_TCP_PROTO=tcp
//	REDIS_PORT_6379_TCP_PORT=6379
//	REDIS_PORT_6379_TCP_ADDR=10.0.0.11
func DiscoverFromEnv() []Service {
	services := make(map[string]*svcInfo)

	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key, val := parts[0], parts[1]

		if strings.HasSuffix(key, "_SERVICE_HOST") {
			name := strings.TrimSuffix(key, "_SERVICE_HOST")
			name = envNameToServiceName(name)
			if _, ok := services[name]; !ok {
				services[name] = &svcInfo{ports: make(map[uint16]ServicePort)}
			}
			services[name].host = val
			continue
		}

		if strings.HasSuffix(key, "_SERVICE_PORT") {
			name := strings.TrimSuffix(key, "_SERVICE_PORT")
			name = envNameToServiceName(name)
			port, err := strconv.ParseUint(val, 10, 16)
			if err != nil {
				continue
			}
			if _, ok := services[name]; !ok {
				services[name] = &svcInfo{ports: make(map[uint16]ServicePort)}
			}
			services[name].ports[uint16(port)] = ServicePort{Port: uint16(port)}
			continue
		}

		// Parse _PORT_<port>_TCP / _PORT_<port>_UDP patterns for named ports
		if idx := strings.Index(key, "_PORT_"); idx > 0 {
			rest := key[idx+6:] // after _PORT_
			portProto := strings.SplitN(rest, "_", 2)
			if len(portProto) != 2 {
				continue
			}
			portStr, protoAndSuffix := portProto[0], portProto[1]
			port, err := strconv.ParseUint(portStr, 10, 16)
			if err != nil {
				continue
			}

			proto := ""
			suffix := protoAndSuffix
			for _, p := range []string{"TCP", "UDP"} {
				if strings.HasPrefix(protoAndSuffix, p) {
					proto = strings.ToLower(p)
					suffix = strings.TrimPrefix(protoAndSuffix, p)
					break
				}
			}
			if proto == "" {
				continue
			}

			name := envNameToServiceName(key[:idx])

			if _, ok := services[name]; !ok {
				services[name] = &svcInfo{ports: make(map[uint16]ServicePort)}
			}

			sp := services[name].ports[uint16(port)]
			sp.Port = uint16(port)
			sp.Protocol = proto

			switch suffix {
			case "_ADDR":
				if services[name].host == "" {
					services[name].host = val
				}
			case "_PORT":
				// redundant port info
			case "_PROTO":
				sp.Protocol = val
			case "":
				// _PORT_<port>_TCP=tcp://host:port (the compound URL)
			}

			services[name].ports[uint16(port)] = sp
		}
	}

	var result []Service
	for name, info := range services {
		if name == "kubernetes" {
			// Skip the default kubernetes API server service unless it has non-standard ports
			if len(info.ports) <= 1 {
				hasOnly443 := true
				for p := range info.ports {
					if p != 443 {
						hasOnly443 = false
					}
				}
				if hasOnly443 || len(info.ports) == 0 {
					svc := serviceFromEnvInfo(name, info)
					result = append(result, svc)
					continue
				}
			}
		}
		svc := serviceFromEnvInfo(name, info)
		result = append(result, svc)
	}

	log.Debugf("[services/env] discovered %d services from environment variables", len(result))
	return result
}

func serviceFromEnvInfo(name string, info *svcInfo) Service {
	svc := Service{
		Name:      name,
		Namespace: guessNamespaceFromEnv(name),
		ClusterIP: info.host,
		Sources:   []string{"env"},
	}
	for _, p := range info.ports {
		svc.Ports = append(svc.Ports, p)
	}
	return svc
}

// envNameToServiceName converts K8s env var prefix to a service name.
// K8s uppercases and replaces hyphens with underscores:
// my-service -> MY_SERVICE -> my-service
func envNameToServiceName(envPrefix string) string {
	return strings.ToLower(strings.ReplaceAll(envPrefix, "_", "-"))
}

// guessNamespaceFromEnv tries to determine the namespace.
// Env vars don't carry namespace info, so we default to the pod's own namespace.
func guessNamespaceFromEnv(name string) string {
	if ns := os.Getenv("POD_NAMESPACE"); ns != "" {
		return ns
	}
	// Try reading namespace from the service account mount
	data, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err == nil && len(data) > 0 {
		return strings.TrimSpace(string(data))
	}
	return "default"
}
