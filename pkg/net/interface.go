package net

import (
	"fmt"
	"net"

	"github.com/ssst0n3/awesome_libs/awesome_error"
)

// InterfaceByIP iterates through all network interfaces to find the one that owns this source IP.
func InterfaceByIP(target net.IP) (*net.Interface, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("could not list network interfaces: %w", err)
	}
	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			// it's unlikely but just skip this interface
			awesome_error.CheckWarning(err)
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip != nil && ip.Equal(target) {
				return &iface, nil
			}
		}
	}
	return nil, fmt.Errorf("could not find interface for source IP %s", target)
}
