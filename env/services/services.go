package services

import (
	"fmt"
	"net"
	"os"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"
)

const CommandName = "services"

type ServicePort struct {
	Name     string `json:"name,omitempty"`
	Port     uint16 `json:"port"`
	Protocol string `json:"protocol,omitempty"`
}

func (p ServicePort) String() string {
	if p.Name != "" {
		return fmt.Sprintf("%d/%s(%s)", p.Port, p.Protocol, p.Name)
	}
	if p.Protocol != "" {
		return fmt.Sprintf("%d/%s", p.Port, p.Protocol)
	}
	return fmt.Sprintf("%d", p.Port)
}

type Service struct {
	Name      string        `json:"name"`
	Namespace string        `json:"namespace"`
	ClusterIP string        `json:"cluster_ip,omitempty"`
	Domain    string        `json:"domain,omitempty"`
	Ports     []ServicePort `json:"ports,omitempty"`
	PodIPs    []string      `json:"pod_ips,omitempty"`
	Sources   []string      `json:"sources"`
}

func (s Service) Key() string {
	return s.Namespace + "/" + s.Name
}

type Options struct {
	Zone       string
	CIDR       string
	Threads    int
	Methods    []string
	OutputFile string
}

func DefaultOptions() Options {
	cidr := ""
	if host := os.Getenv("KUBERNETES_SERVICE_HOST"); host != "" {
		cidr = host + "/16"
	}
	return Options{
		Zone:    "cluster.local",
		CIDR:    cidr,
		Threads: 16,
		Methods: []string{"all"},
	}
}

func methodEnabled(methods []string, name string) bool {
	for _, m := range methods {
		if m == "all" || m == name {
			return true
		}
	}
	return false
}

func Discover(opts Options) (result []Service, err error) {
	svcMap := make(map[string]*Service)

	if methodEnabled(opts.Methods, "env") {
		envServices := DiscoverFromEnv()
		for i := range envServices {
			mergeService(svcMap, &envServices[i])
		}
		log.Debugf("[services] env discovery found %d services", len(envServices))
	}

	if methodEnabled(opts.Methods, "wildcard") || methodEnabled(opts.Methods, "axfr") || methodEnabled(opts.Methods, "cidr") {
		if !CheckKubeDNS(opts.Zone) {
			log.Warnf("[services] K8s DNS not detected, skipping DNS-based discovery")
		} else {
			if methodEnabled(opts.Methods, "wildcard") {
				wildcardServices := ScanWildcard(opts.Zone)
				for i := range wildcardServices {
					mergeService(svcMap, &wildcardServices[i])
				}
				log.Debugf("[services] wildcard discovery found %d services", len(wildcardServices))
			}

			if methodEnabled(opts.Methods, "axfr") {
				axfrServices := ScanAXFR(opts.Zone)
				for i := range axfrServices {
					mergeService(svcMap, &axfrServices[i])
				}
				log.Debugf("[services] AXFR discovery found %d services", len(axfrServices))
			}

			if methodEnabled(opts.Methods, "cidr") {
				cidr := opts.CIDR
				if cidr == "" {
					log.Warnf("[services] no CIDR specified, skipping CIDR scan")
				} else {
					_, ipNet, parseErr := net.ParseCIDR(cidr)
					if parseErr != nil {
						log.Warnf("[services] invalid CIDR %q: %v", cidr, parseErr)
					} else {
						cidrServices := ScanCIDR(ipNet, opts.Zone, opts.Threads)
						for i := range cidrServices {
							mergeService(svcMap, &cidrServices[i])
						}
						log.Debugf("[services] CIDR scan found %d services", len(cidrServices))
					}
				}
			}
		}
	}

	for _, svc := range svcMap {
		svc.PodIPs = uniqueStrings(svc.PodIPs)
		svc.Sources = uniqueStrings(svc.Sources)
		svc.Ports = uniquePorts(svc.Ports)
		result = append(result, *svc)
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Namespace != result[j].Namespace {
			return result[i].Namespace < result[j].Namespace
		}
		return result[i].Name < result[j].Name
	})

	return result, nil
}

func mergeService(m map[string]*Service, incoming *Service) {
	key := incoming.Key()
	existing, ok := m[key]
	if !ok {
		copied := *incoming
		m[key] = &copied
		return
	}
	if existing.ClusterIP == "" && incoming.ClusterIP != "" {
		existing.ClusterIP = incoming.ClusterIP
	}
	if existing.Domain == "" && incoming.Domain != "" {
		existing.Domain = incoming.Domain
	}
	existing.Ports = append(existing.Ports, incoming.Ports...)
	existing.PodIPs = append(existing.PodIPs, incoming.PodIPs...)
	existing.Sources = append(existing.Sources, incoming.Sources...)
}

func uniqueStrings(s []string) []string {
	seen := make(map[string]struct{})
	var out []string
	for _, v := range s {
		if v == "" {
			continue
		}
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			out = append(out, v)
		}
	}
	return out
}

func uniquePorts(ports []ServicePort) []ServicePort {
	seen := make(map[string]struct{})
	var out []ServicePort
	for _, p := range ports {
		key := fmt.Sprintf("%d/%s", p.Port, strings.ToLower(p.Protocol))
		if _, ok := seen[key]; !ok {
			seen[key] = struct{}{}
			out = append(out, p)
		}
	}
	return out
}
