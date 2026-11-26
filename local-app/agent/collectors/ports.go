package collectors

import (
	"github.com/shirou/gopsutil/v4/net"
)

type OpenPort struct {
	Port uint32
	Ip   string
	PID  int32
}

func GetOpenPorts() ([]OpenPort, error) {
	var ports []OpenPort
	var newPort OpenPort
	result, err := net.Connections("tcp")
	if err != nil {
		return nil, err
	}
	for _, port := range result {
		if port.Status == "LISTEN" {
			newPort.Port = port.Laddr.Port
			newPort.Ip = port.Laddr.IP
			newPort.PID = port.Pid
			ports = append(ports, newPort)
		}
	}
	return ports, nil
}
