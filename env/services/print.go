package services

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ctrsploit/sploit-spec/pkg/printer"
	"github.com/ctrsploit/sploit-spec/pkg/result"
	"github.com/ctrsploit/sploit-spec/pkg/result/item"
)

type ServiceItem struct {
	SubTitle result.SubTitle `json:"-"`
	Name     item.Short     `json:"name"`
	IP       item.Short     `json:"cluster_ip"`
	Ports    item.Short     `json:"ports"`
	PodIPs   item.Long      `json:"pod_ips"`
	Sources  item.Short     `json:"sources"`
}

type ServicesResult struct {
	Name     result.Title  `json:"name"`
	Count    item.Short    `json:"count"`
	Services []ServiceItem `json:"services"`
}

func Human(services []Service) ServicesResult {
	h := ServicesResult{
		Name: result.Title{
			Name: "K8s Services",
		},
		Count: item.Short{
			Name:   "Total Services Discovered",
			Result: fmt.Sprintf("%d", len(services)),
		},
	}

	for _, svc := range services {
		var portStrs []string
		for _, p := range svc.Ports {
			portStrs = append(portStrs, p.String())
		}

		si := ServiceItem{
			SubTitle: result.SubTitle{
				Name: svc.Key(),
			},
			Name: item.Short{
				Name:   "service",
				Result: svc.Name,
			},
			IP: item.Short{
				Name:   "cluster-ip",
				Result: svc.ClusterIP,
			},
			Ports: item.Short{
				Name:   "ports",
				Result: strings.Join(portStrs, ", "),
			},
			Sources: item.Short{
				Name:   "sources",
				Result: strings.Join(svc.Sources, ", "),
			},
		}

		if len(svc.PodIPs) > 0 {
			si.PodIPs = item.Long{
				Name:   "pod-ips",
				Result: strings.Join(svc.PodIPs, ", "),
			}
		}

		h.Services = append(h.Services, si)
	}

	return h
}

func Print(opts Options) error {
	services, err := Discover(opts)
	if err != nil {
		return err
	}

	u := result.Union{
		Machine: services,
		Human:   Human(services),
	}
	fmt.Println(printer.Printer.Print(u))

	if opts.OutputFile != "" {
		return exportJSON(services, opts.OutputFile)
	}
	return nil
}

func exportJSON(services []Service, path string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("open output file: %w", err)
	}
	defer f.Close()

	for _, svc := range services {
		data, err := json.Marshal(svc)
		if err != nil {
			continue
		}
		fmt.Fprintf(f, "%s\n", data)
	}
	return nil
}
