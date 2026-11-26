package main

import (
	"fmt"
	"site-health-monitoring/agent/collectors"
)

func main() {
	fmt.Println("Agent started")
	fmt.Println()

	result1, err := collectors.GetOSInfo()
	if err != nil {
		fmt.Println("Info collection error:", err)
		return
	}

	fmt.Println("System Information:")
	fmt.Printf("OS: %s\n", result1.OS)
	fmt.Printf("Architecture: %s\n", result1.Architecture)
	fmt.Printf("Version: %s\n", result1.Version)
	fmt.Println()

	result2, err := collectors.GetMemoryInfo()
	if err != nil {
		fmt.Println("Memory collection error:", err)
		return
	}

	totalGB := float64(result2.Total) / (1024 * 1024 * 1024)
	usedGB := float64(result2.Used) / (1024 * 1024 * 1024)

	fmt.Println("Memory Information:")
	fmt.Printf("Total: %.2f GB\n", totalGB)
	fmt.Printf("Used: %.2f GB\n", usedGB)
	fmt.Printf("Usage: %.2f%%\n", result2.UsedPercent)

	result3, err := collectors.GetCPUInfo()
	if err == nil {
		fmt.Println()
		fmt.Println("CPU Information:")
		fmt.Printf("Usage: %.2f%%\n", result3.UsagePercent)
	}

	result4, err := collectors.GetLoggedInUsers()
	if err != nil {
		fmt.Println("Logged-in users collection error:", err)
	} else {
		fmt.Println()
		fmt.Println("Logged-in Users:")

		for _, user := range result4 {
			fmt.Printf(
				"Username: %s | Terminal: %s | Host: %s\n",
				user.Username,
				user.Terminal,
				user.Host,
			)
		}
	}

	result5, err := collectors.GetOpenPorts()
	if err != nil {
		fmt.Println("Open ports collection error:", err)
	} else {
		fmt.Println()
		fmt.Println("Open / Listening Ports:")

		for _, port := range result5 {
			fmt.Printf(
				"IP: %s | Port: %d | PID: %d\n",
				port.Ip,
				port.Port,
				port.PID,
			)
		}
	}
	result6, err := collectors.GetActiveConnections()
	if err != nil {
		fmt.Println("Active connections collection error:", err)
	} else {
		fmt.Println()
		fmt.Println("Active Network Connections:")

		for _, connection := range result6 {
			fmt.Printf(
				"Local: %s:%d | Remote: %s:%d | Status: %s | PID: %d\n",
				connection.LocalIP,
				connection.LocalPort,
				connection.RemoteIP,
				connection.RemotePort,
				connection.Status,
				connection.PID,
			)
		}
	}
	result7, err := collectors.GetInstalledApplications()
	if err != nil {
		fmt.Println("Installed applications collection error:", err)
	} else {
		fmt.Println()
		fmt.Println("Installed Applications:")

		for _, application := range result7 {
			fmt.Printf("Name: %s\n", application.Name)
		}
	}
	result8, err := collectors.GetFirewallInfo()
	if err != nil {
		fmt.Println("Firewall collection error:", err)
	} else {
		fmt.Println()
		fmt.Println("Firewall Information:")
		fmt.Printf("Status: %s\n", result8.Status)
	}
}
