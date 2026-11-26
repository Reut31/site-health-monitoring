package collectors

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/v4/net"
)

type OpenPort struct {
	Port     uint32
	Ip       string
	PID      int32
	Protocol string
	Family   string
}

func GetOpenPorts() ([]OpenPort, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(), 10*time.Second,
	)
	defer cancel()

	ports := make([]OpenPort, 0)
	seen := make(map[OpenPort]bool)

	result, err := net.ConnectionsWithContext(ctx, "inet")

	if err != nil {
		return nil, fmt.Errorf(
			"network port collection: %w", err,
		)
	}

	for _, port := range result {
		var protocol string
		var family string

		// Identify protocol
		switch port.Type {
		case syscall.SOCK_STREAM:
			if port.Status != "LISTEN" {
				continue
			}
			protocol = "TCP"

		case syscall.SOCK_DGRAM:
			protocol = "UDP"

		default:
			continue
		}

		// Identify IP address family
		switch port.Family {
		case syscall.AF_INET:
			family = "IPv4"

		case syscall.AF_INET6:
			family = "IPv6"

		default:
			continue
		}

		// Validate port number
		if port.Laddr.Port == 0 || port.Laddr.Port > 65535 {
			continue
		}

		newPort := OpenPort{
			Port:     port.Laddr.Port,
			Ip:       strings.TrimSpace(port.Laddr.IP),
			PID:      port.Pid,
			Protocol: protocol,
			Family:   family,
		}

		// Remove duplicate records
		if seen[newPort] {
			continue
		}

		seen[newPort] = true
		ports = append(ports, newPort)
	}

	// Sort ports
	sort.Slice(ports, func(i, j int) bool {
		if ports[i].Port != ports[j].Port {
			return ports[i].Port < ports[j].Port
		}

		if ports[i].Protocol != ports[j].Protocol {
			return ports[i].Protocol < ports[j].Protocol
		}

		if ports[i].Family != ports[j].Family {
			return ports[i].Family < ports[j].Family
		}

		if ports[i].Ip != ports[j].Ip {
			return ports[i].Ip < ports[j].Ip
		}

		return ports[i].PID < ports[j].PID
	})

	return ports, nil
}
