package main

import (
	"fmt"
	"site-health-monitoring/agent/collectors"
)

func main() {
	fmt.Println("Agent started")

	// Operating System
	result1, err := collectors.GetOSInfo()

	fmt.Println()
	fmt.Println("System Information:")

	if err != nil {
		fmt.Println("Status: Unavailable")
		fmt.Println("Collection error:", err)
	} else {
		fmt.Printf("OS: %s\n", result1.OS)
		fmt.Printf("Architecture: %s\n", result1.Architecture)
		fmt.Printf("Version: %s\n", result1.Version)
	}

	// Memory
	result2, err := collectors.GetMemoryInfo()

	fmt.Println()
	fmt.Println("Memory Information:")

	if err != nil {
		fmt.Println("Status: Unavailable")
		fmt.Println("Collection error:", err)
	} else {
		totalGB := float64(result2.Total) / (1024 * 1024 * 1024)
		usedGB := float64(result2.Used) / (1024 * 1024 * 1024)

		fmt.Printf("Total: %.2f GB\n", totalGB)
		fmt.Printf("Used: %.2f GB\n", usedGB)
		fmt.Printf("Usage: %.2f%%\n", result2.UsedPercent)
	}

	// CPU
	result3, err := collectors.GetCPUInfo()

	fmt.Println()
	fmt.Println("CPU Information:")

	if err != nil {
		fmt.Println("Status: Unavailable")
		fmt.Println("Collection error:", err)
	} else {
		fmt.Printf("Usage: %.2f%%\n", result3.UsagePercent)
	}

	// Logged-in Users
	result4, err := collectors.GetLoggedInUsers()

	fmt.Println()
	fmt.Println("Logged-in Users:")

	if err != nil {
		fmt.Println("Status: Unavailable")
		fmt.Println("Collection error:", err)
	} else {
		for _, user := range result4 {
			fmt.Printf(
				"Username: %s | Terminal: %s | Host: %s\n",
				user.Username,
				user.Terminal,
				user.Host,
			)
		}
	}

	// Open Ports
	result5, err := collectors.GetOpenPorts()

	fmt.Println()
	fmt.Println("Open / Listening Ports:")

	if err != nil {
		fmt.Println("Status: Unavailable")
		fmt.Println("Collection error:", err)
	} else {
		for _, port := range result5 {
			fmt.Printf(
				"IP: %s | Port: %d | PID: %d\n",
				port.Ip,
				port.Port,
				port.PID,
			)
		}
	}

	// Active Network Connections
	result6, err := collectors.GetActiveConnections()

	fmt.Println()
	fmt.Println("Active Network Connections:")

	if err != nil {
		fmt.Println("Status: Unavailable")
		fmt.Println("Collection error:", err)
	} else {
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

	// Installed Applications
	result7, err := collectors.GetInstalledApplications()

	fmt.Println()
	fmt.Println("Installed Applications:")

	if err != nil {
		fmt.Println("Status: Unavailable")
		fmt.Println("Collection error:", err)
	} else {
		for _, application := range result7 {
			fmt.Printf("Name: %s\n", application.Name)
		}
	}

	// Firewall
	result8, err := collectors.GetFirewallInfo()

	fmt.Println()
	fmt.Println("Firewall Information:")

	if err != nil {
		fmt.Println("Status: Unavailable")
		fmt.Println("Collection error:", err)
	} else {
		fmt.Printf("Status: %s\n", result8.Status)
	}

	// Antivirus
	result9, err := collectors.GetAntivirusInfo()

	fmt.Println()
	fmt.Println("Antivirus Information:")

	if err != nil {
		fmt.Println("Name: Unknown")
		fmt.Println("Status: Unavailable")
		fmt.Println("Collection error:", err)
	} else {
		fmt.Printf("Name: %s\n", result9.Name)
		fmt.Printf("Status: %s\n", result9.Status)
	}

	// Security Updates
	result10, err := collectors.GetSecurityUpdateInfo()

	fmt.Println()
	fmt.Println("Security Update Information:")

	if err != nil {
		fmt.Println("Status: Unavailable")
		fmt.Println("Collection error:", err)
	} else {
		fmt.Printf("Status: %s\n", result10.Status)

		if result10.Status == "UpdatesAvailable" ||
			result10.Status == "UpToDate" {
			fmt.Printf("Available Updates: %d\n", result10.AvailableCount)
		} else {
			fmt.Println("Available Updates: Unknown")
		}
	}

	// Failed Login Attempts
	result11, err := collectors.GetFailedLoginInfo()

	fmt.Println()
	fmt.Println("Failed Login Information:")
	fmt.Printf("Time Range: %s\n", result11.TimeRange)

	if err != nil {
		fmt.Println("Status: Unavailable")
		fmt.Println("Collection error:", err)
	} else {
		fmt.Printf("Status: %s\n", result11.Status)

		if result11.Status == "Available" {
			fmt.Printf("Failed Authentication Events: %d\n", result11.Count)
		} else {
			fmt.Println("Failed Authentication Events: Unknown")
		}
	}

	fmt.Println()
	fmt.Println("Collection completed")
}
