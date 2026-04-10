package services

import (
	"context"
	"encoding/binary"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
)

const dnsTimeout = 2 * time.Second

// CheckKubeDNS validates that K8s cluster DNS is available by probing well-known records.
func CheckKubeDNS(zone string) bool {
	r := net.DefaultResolver
	ctx, cancel := context.WithTimeout(context.Background(), dnsTimeout)
	defer cancel()

	// kubernetes.default.svc.<zone>
	ips, err := r.LookupIPAddr(ctx, "kubernetes.default.svc."+zone)
	if err == nil && len(ips) > 0 {
		log.Debugf("[services/dns] kubernetes.default.svc.%s resolved to %v", zone, ips)
		return true
	}

	// Try short form (in case zone is wrong but search domain is set)
	ctx2, cancel2 := context.WithTimeout(context.Background(), dnsTimeout)
	defer cancel2()
	ips, err = r.LookupIPAddr(ctx2, "kubernetes.default.svc")
	if err == nil && len(ips) > 0 {
		log.Debugf("[services/dns] kubernetes.default.svc resolved to %v", ips)
		return true
	}

	// dns-version TXT
	ctx3, cancel3 := context.WithTimeout(context.Background(), dnsTimeout)
	defer cancel3()
	txts, err := r.LookupTXT(ctx3, "dns-version."+zone)
	if err == nil && len(txts) > 0 {
		log.Debugf("[services/dns] dns-version.%s = %v", zone, txts)
		return true
	}

	return false
}

// ScanWildcard uses CoreDNS wildcard SRV queries to discover services.
// Queries any.any.svc.<zone> and any.any.any.svc.<zone>.
func ScanWildcard(zone string) []Service {
	searchDomains := []string{
		dns.Fqdn("any.any.svc." + zone),
		dns.Fqdn("any.any.any.svc." + zone),
	}

	r := net.DefaultResolver
	var allServices []Service

	for _, domain := range searchDomains {
		ctx, cancel := context.WithTimeout(context.Background(), dnsTimeout*2)
		_, srvs, err := r.LookupSRV(ctx, "", "", domain)
		cancel()
		if err != nil {
			log.Debugf("[services/dns] wildcard SRV query %s failed: %v", domain, err)
			continue
		}
		for _, srv := range srvs {
			svc := parseSRVTarget(srv, zone, "wildcard")
			if svc != nil {
				allServices = append(allServices, *svc)
			}
		}
	}

	return allServices
}

// ScanAXFR performs a DNS zone transfer to dump all records.
func ScanAXFR(zone string) []Service {
	dnsServer := "ns.dns." + zone + ":53"
	fqdnZone := dns.Fqdn(zone)

	t := new(dns.Transfer)
	m := new(dns.Msg)
	m.SetAxfr(fqdnZone)

	ch, err := t.In(m, dnsServer)
	if err != nil {
		log.Debugf("[services/dns] AXFR to %s failed: %v", dnsServer, err)
		return nil
	}

	seen := make(map[string]*Service)
	for envelope := range ch {
		if envelope.Error != nil {
			log.Debugf("[services/dns] AXFR envelope error: %v", envelope.Error)
			break
		}
		for _, rr := range envelope.RR {
			name := rr.Header().Name
			if !isServiceDomain(name, zone) {
				continue
			}
			svcName, ns := parseServiceDomain(name, zone)
			if svcName == "" {
				continue
			}

			key := ns + "/" + svcName
			svc, ok := seen[key]
			if !ok {
				svc = &Service{
					Name:      svcName,
					Namespace: ns,
					Sources:   []string{"axfr"},
				}
				seen[key] = svc
			}

			switch v := rr.(type) {
			case *dns.A:
				if svc.ClusterIP == "" {
					svc.ClusterIP = v.A.String()
				}
			case *dns.SRV:
				svc.Ports = append(svc.Ports, ServicePort{
					Port: v.Port,
				})
			}
			if svc.Domain == "" {
				svc.Domain = name
			}
		}
	}

	var result []Service
	for _, svc := range seen {
		result = append(result, *svc)
	}
	return result
}

// ScanCIDR walks the service CIDR, performs PTR lookups to find service domains,
// then SRV lookups to discover ports. Multi-threaded.
func ScanCIDR(ipNet *net.IPNet, zone string, threads int) []Service {
	if threads < 1 {
		threads = 1
	}

	ips := cidrToIPs(ipNet)
	r := net.DefaultResolver

	type record struct {
		ip     net.IP
		domain string
	}

	// Phase 1: PTR scan
	ptrCh := make(chan record, 100)
	var wg sync.WaitGroup

	chunkSize := len(ips) / threads
	if chunkSize < 1 {
		chunkSize = 1
	}

	for i := 0; i < len(ips); i += chunkSize {
		end := i + chunkSize
		if end > len(ips) {
			end = len(ips)
		}
		wg.Add(1)
		go func(chunk []net.IP) {
			defer wg.Done()
			for _, ip := range chunk {
				ctx, cancel := context.WithTimeout(context.Background(), dnsTimeout)
				names, err := r.LookupAddr(ctx, ip.String())
				cancel()
				if err != nil || len(names) == 0 {
					continue
				}
				for _, name := range names {
					log.Debugf("[services/dns] PTR %v -> %s", ip, name)
					ptrCh <- record{ip: copyIP(ip), domain: name}
				}
			}
		}(ips[i:end])
	}
	go func() {
		wg.Wait()
		close(ptrCh)
	}()

	// Collect PTR results
	var ptrRecords []record
	for rec := range ptrCh {
		ptrRecords = append(ptrRecords, rec)
	}

	// Phase 2: Parse domains and SRV enrich
	svcMap := make(map[string]*Service)
	for _, rec := range ptrRecords {
		domain := rec.domain

		if isPodServiceDomain(domain, zone) {
			svcName, ns := parsePodServiceDomain(domain, zone)
			if svcName == "" {
				continue
			}
			key := ns + "/" + svcName
			svc, ok := svcMap[key]
			if !ok {
				svcDomain := svcName + "." + ns + ".svc." + dns.Fqdn(zone)
				svc = &Service{
					Name:      svcName,
					Namespace: ns,
					Domain:    svcDomain,
					Sources:   []string{"cidr"},
				}
				svcMap[key] = svc
			}
			podIP := podIPFromDomain(domain, zone)
			if podIP != "" {
				svc.PodIPs = append(svc.PodIPs, podIP)
			}
		} else if isServiceDomain(domain, zone) {
			svcName, ns := parseServiceDomain(domain, zone)
			if svcName == "" {
				continue
			}
			key := ns + "/" + svcName
			svc, ok := svcMap[key]
			if !ok {
				svc = &Service{
					Name:      svcName,
					Namespace: ns,
					Domain:    domain,
					Sources:   []string{"cidr"},
				}
				svcMap[key] = svc
			}
			if svc.ClusterIP == "" {
				svc.ClusterIP = rec.ip.String()
			}
		}
	}

	// SRV enrichment for each discovered service
	for _, svc := range svcMap {
		if svc.Domain == "" {
			continue
		}
		ctx, cancel := context.WithTimeout(context.Background(), dnsTimeout)
		_, srvs, err := r.LookupSRV(ctx, "", "", svc.Domain)
		cancel()
		if err != nil {
			continue
		}
		for _, srv := range srvs {
			svc.Ports = append(svc.Ports, ServicePort{
				Port: srv.Port,
			})
		}
	}

	var result []Service
	for _, svc := range svcMap {
		result = append(result, *svc)
	}
	return result
}

// --- Domain parsing helpers ---
// K8s DNS spec: https://github.com/kubernetes/dns/blob/master/docs/specification.md
// Service: <svc>.<ns>.svc.<zone>
// Pod service: <ip-dashed>.<svc>.<ns>.svc.<zone>

func isServiceDomain(domain, zone string) bool {
	fqdn := dns.Fqdn(domain)
	fqdnZone := dns.Fqdn(zone)
	parts := strings.Split(fqdn, ".")
	zoneParts := strings.Split(fqdnZone, ".")
	if len(parts) < len(zoneParts)+3 {
		return false
	}
	// Check for "svc" at the right position
	reversed := reverseSlice(parts)
	if len(reversed) > len(zoneParts) && reversed[len(zoneParts)] == "svc" {
		return true
	}
	return false
}

func isPodServiceDomain(domain, zone string) bool {
	fqdn := dns.Fqdn(domain)
	fqdnZone := dns.Fqdn(zone)
	trimmed := strings.TrimSuffix(fqdn, fqdnZone)
	parts := strings.Split(strings.TrimSuffix(trimmed, "."), ".")
	// pod service format: <ip-dashed>.<svc>.<ns>.svc
	if len(parts) >= 4 && parts[len(parts)-1] == "svc" {
		return true
	}
	return false
}

func parseServiceDomain(domain, zone string) (svcName, namespace string) {
	fqdn := dns.Fqdn(domain)
	fqdnZone := dns.Fqdn(zone)
	trimmed := strings.TrimSuffix(fqdn, fqdnZone)
	trimmed = strings.TrimSuffix(trimmed, ".")
	parts := strings.Split(trimmed, ".")
	// Expected: <svc>.<ns>.svc or <port>.<proto>.<svc>.<ns>.svc
	if len(parts) < 3 {
		return "", ""
	}
	// Find "svc" marker
	for i, p := range parts {
		if p == "svc" && i >= 2 {
			return parts[i-2], parts[i-1]
		}
	}
	return "", ""
}

func parsePodServiceDomain(domain, zone string) (svcName, namespace string) {
	fqdn := dns.Fqdn(domain)
	fqdnZone := dns.Fqdn(zone)
	trimmed := strings.TrimSuffix(fqdn, fqdnZone)
	trimmed = strings.TrimSuffix(trimmed, ".")
	parts := strings.Split(trimmed, ".")
	// <ip-dashed>.<svc>.<ns>.svc
	for i, p := range parts {
		if p == "svc" && i >= 3 {
			return parts[i-2], parts[i-1]
		}
	}
	return "", ""
}

func podIPFromDomain(domain, zone string) string {
	fqdn := dns.Fqdn(domain)
	fqdnZone := dns.Fqdn(zone)
	trimmed := strings.TrimSuffix(fqdn, fqdnZone)
	trimmed = strings.TrimSuffix(trimmed, ".")
	parts := strings.Split(trimmed, ".")
	if len(parts) < 4 {
		return ""
	}
	return strings.ReplaceAll(parts[0], "-", ".")
}

func parseSRVTarget(srv *net.SRV, zone, source string) *Service {
	target := srv.Target
	if isServiceDomain(target, zone) {
		svcName, ns := parseServiceDomain(target, zone)
		if svcName == "" {
			return nil
		}
		return &Service{
			Name:      svcName,
			Namespace: ns,
			Domain:    target,
			Ports: []ServicePort{
				{Port: srv.Port},
			},
			Sources: []string{source},
		}
	}
	if isPodServiceDomain(target, zone) {
		svcName, ns := parsePodServiceDomain(target, zone)
		if svcName == "" {
			return nil
		}
		podIP := podIPFromDomain(target, zone)
		svc := &Service{
			Name:      svcName,
			Namespace: ns,
			Ports: []ServicePort{
				{Port: srv.Port},
			},
			Sources: []string{source},
		}
		if podIP != "" {
			svc.PodIPs = []string{podIP}
		}
		return svc
	}
	return nil
}

// --- Network helpers ---

func cidrToIPs(ipNet *net.IPNet) []net.IP {
	ip4 := ipNet.IP.To4()
	if ip4 == nil {
		return nil
	}
	mask := binary.BigEndian.Uint32(ipNet.Mask)
	start := binary.BigEndian.Uint32(ip4)
	finish := (start & mask) | (mask ^ 0xffffffff)

	var ips []net.IP
	for i := start; i <= finish; i++ {
		ip := make(net.IP, 4)
		binary.BigEndian.PutUint32(ip, i)
		ips = append(ips, ip)
	}
	return ips
}

func copyIP(ip net.IP) net.IP {
	dup := make(net.IP, len(ip))
	copy(dup, ip)
	return dup
}

func reverseSlice(s []string) []string {
	out := make([]string, len(s))
	for i, v := range s {
		out[len(s)-1-i] = v
	}
	return out
}

