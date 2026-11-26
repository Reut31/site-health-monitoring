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

type ActiveConnection struct {
	LocalIP    string
	LocalPort  uint32
	RemoteIP   string
	RemotePort uint32
	Status     string
	PID        int32
	Protocol   string
	Family     string
}

func GetActiveConnections() ([]ActiveConnection, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(), 10*time.Second,
	)
	defer cancel()

	activeConnections := make([]ActiveConnection, 0)
	seen := make(map[ActiveConnection]bool)

	result, err := net.ConnectionsWithContext(ctx, "tcp")

	if err != nil {
		return nil, fmt.Errorf(
			"active connection collection: %w", err,
		)
	}

	for _, connection := range result {

		// Collect established TCP connections only
		if connection.Status != "ESTABLISHED" {
			continue
		}

		// Identify IP address family
		var family string

		switch connection.Family {
		case syscall.AF_INET:
			family = "IPv4"

		case syscall.AF_INET6:
			family = "IPv6"

		default:
			continue
		}

		// Validate local and remote ports
		if connection.Laddr.Port == 0 ||
			connection.Laddr.Port > 65535 ||
			connection.Raddr.Port == 0 ||
			connection.Raddr.Port > 65535 {

			continue
		}

		newConnection := ActiveConnection{
			LocalIP:    strings.TrimSpace(connection.Laddr.IP),
			LocalPort:  connection.Laddr.Port,
			RemoteIP:   strings.TrimSpace(connection.Raddr.IP),
			RemotePort: connection.Raddr.Port,
			Status:     connection.Status,
			PID:        connection.Pid,
			Protocol:   "TCP",
			Family:     family,
		}

		// Validate IP addresses are present
		if newConnection.LocalIP == "" ||
			newConnection.RemoteIP == "" {

			continue
		}

		// Remove duplicate records
		if seen[newConnection] {
			continue
		}

		seen[newConnection] = true

		activeConnections = append(
			activeConnections,
			newConnection,
		)
	}

	// Sort active connections
	sort.Slice(activeConnections, func(i, j int) bool {
		a := activeConnections[i]
		b := activeConnections[j]

		if a.PID != b.PID {
			return a.PID < b.PID
		}

		if a.Family != b.Family {
			return a.Family < b.Family
		}

		if a.LocalIP != b.LocalIP {
			return a.LocalIP < b.LocalIP
		}

		if a.LocalPort != b.LocalPort {
			return a.LocalPort < b.LocalPort
		}

		if a.RemoteIP != b.RemoteIP {
			return a.RemoteIP < b.RemoteIP
		}

		return a.RemotePort < b.RemotePort
	})

	return activeConnections, nil
}
