package collectors

import (
	"github.com/shirou/gopsutil/v4/net"
)

type ActiveConnection struct {
	LocalIP    string
	LocalPort  uint32
	RemoteIP   string
	RemotePort uint32
	Status     string
	PID        int32
}

func GetActiveConnections() ([]ActiveConnection, error) {
	var activeConnections []ActiveConnection
	var newConnection ActiveConnection
	result, err := net.Connections("tcp")
	if err != nil {
		return nil, err
	}
	for _, connection := range result {
		if connection.Status == "ESTABLISHED" {
			newConnection.LocalIP = connection.Laddr.IP
			newConnection.LocalPort = connection.Laddr.Port
			newConnection.RemoteIP = connection.Raddr.IP
			newConnection.RemotePort = connection.Raddr.Port
			newConnection.Status = connection.Status
			newConnection.PID = connection.Pid
			activeConnections = append(activeConnections, newConnection)
		}
	}
	return activeConnections, nil
}
