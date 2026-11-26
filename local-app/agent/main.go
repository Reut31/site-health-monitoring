package main

import (
	"fmt"
	"net"
	"strconv"

	"site-health-monitoring/agent/collectors"
)

func main() {
	fmt.Println("Agent started")

	// Operating System
	osInfo, err := collectors.GetOSInfo()

	fmt.Println()
	fmt.Println("System Information:")

	fmt.Printf("OS: %s\n", osInfo.OS)
	fmt.Printf("Architecture: %s\n", osInfo.Architecture)
	fmt.Printf("Version: %s\n", osInfo.Version)

	if err != nil {
		fmt.Println("Status: Unavailable")
		fmt.Println("Collection error:", err)
	}

	// Memory
	memoryInfo, err := collectors.GetMemoryInfo()

	fmt.Println()
	fmt.Println("Memory Information:")

	if err != nil {
		fmt.Println("Status: Unavailable")
		fmt.Println("Collection error:", err)
	} else {
		totalGiB := float64(memoryInfo.Total) / (1024 * 1024 * 1024)
		usedGiB := float64(memoryInfo.Used) / (1024 * 1024 * 1024)

		fmt.Printf("Total: %.2f GiB\n", totalGiB)
		fmt.Printf("Used: %.2f GiB\n", usedGiB)
		fmt.Printf("Usage: %.2f%%\n", memoryInfo.UsedPercent)
	}

	// CPU
	cpuInfo, err := collectors.GetCPUInfo()

	fmt.Println()
	fmt.Println("CPU Information:")

	if err != nil {
		fmt.Println("Status: Unavailable")
		fmt.Println("Collection error:", err)
	} else {
		fmt.Printf("Usage: %.2f%%\n", cpuInfo.UsagePercent)
	}

	// Logged-in Users
	users, err := collectors.GetLoggedInUsers()

	fmt.Println()
	fmt.Println("Logged-in Users:")

	if err != nil {
		fmt.Println("Status: Unavailable")
		fmt.Println("Collection error:", err)
	} else {
		fmt.Printf("Total Records: %d\n", len(users))

		if len(users) == 0 {
			fmt.Println("No logged-in users detected")
		}

		for _, user := range users {
			fmt.Printf(
				"Username: %s | Terminal: %s | Host: %s\n",
				user.Username,
				user.Terminal,
				user.Host,
			)
		}
	}

	// Open Ports
	ports, err := collectors.GetOpenPorts()

	fmt.Println()
	fmt.Println("Open / Listening Ports:")

	if err != nil {
		fmt.Println("Status: Unavailable")
		fmt.Println("Collection error:", err)
	} else {
		fmt.Printf("Total Records: %d\n", len(ports))

		if len(ports) == 0 {
			fmt.Println("No open or listening ports detected")
		}

		for _, port := range ports {
			fmt.Printf(
				"IP: %s | Port: %d | PID: %d | Protocol: %s | Family: %s\n",
				port.Ip,
				port.Port,
				port.PID,
				port.Protocol,
				port.Family,
			)
		}
	}

	// Active Network Connections
	connections, err := collectors.GetActiveConnections()

	fmt.Println()
	fmt.Println("Active Network Connections:")

	if err != nil {
		fmt.Println("Status: Unavailable")
		fmt.Println("Collection error:", err)
	} else {
		fmt.Printf("Total Records: %d\n", len(connections))

		if len(connections) == 0 {
			fmt.Println("No established TCP connections detected")
		}

		for _, connection := range connections {
			localAddress := net.JoinHostPort(
				connection.LocalIP,
				strconv.FormatUint(uint64(connection.LocalPort), 10),
			)

			remoteAddress := net.JoinHostPort(
				connection.RemoteIP,
				strconv.FormatUint(uint64(connection.RemotePort), 10),
			)

			fmt.Printf(
				"Local: %s | Remote: %s | Status: %s | PID: %d | Protocol: %s | Family: %s\n",
				localAddress,
				remoteAddress,
				connection.Status,
				connection.PID,
				connection.Protocol,
				connection.Family,
			)
		}
	}

	// Installed Applications
	applications, err := collectors.GetInstalledApplications()

	fmt.Println()
	fmt.Println("Installed Applications:")

	if err != nil {
		fmt.Println("Status: Unavailable")
		fmt.Println("Collection error:", err)
	} else {
		fmt.Printf("Total Applications: %d\n", len(applications))

		if len(applications) == 0 {
			fmt.Println("No installed applications detected")
		}

		for _, application := range applications {
			fmt.Printf("Name: %s\n", application.Name)
		}
	}

	// Firewall
	firewall, err := collectors.GetFirewallInfo()

	fmt.Println()
	fmt.Println("Firewall Information:")

	if err != nil {
		fmt.Println("Status: Unavailable")
		fmt.Println("Collection error:", err)
	} else {
		fmt.Printf("Status: %s\n", firewall.Status)
	}

	// Antivirus
	antivirus, err := collectors.GetAntivirusInfo()

	fmt.Println()
	fmt.Println("Antivirus Information:")
	fmt.Printf("Name: %s\n", antivirus.Name)

	if err != nil {
		fmt.Println("Status: Unavailable")
		fmt.Println("Collection error:", err)
	} else {
		fmt.Printf("Status: %s\n", antivirus.Status)
	}

	// Security Updates
	updates, err := collectors.GetSecurityUpdateInfo()

	fmt.Println()
	fmt.Println("Security Update Information:")
	fmt.Printf("Scope: %s\n", updates.Scope)

	if err != nil {
		fmt.Println("Status: Unavailable")
		fmt.Println("Available Updates: Unknown")
		fmt.Println("Collection error:", err)
	} else {
		fmt.Printf("Status: %s\n", updates.Status)

		if updates.Status == "UpdatesAvailable" ||
			updates.Status == "UpToDate" {

			fmt.Printf("Available Updates: %d\n", updates.AvailableCount)
		} else {
			fmt.Println("Available Updates: Unknown")
		}
	}

	// Failed Login Attempts
	failedLogins, err := collectors.GetFailedLoginInfo()

	fmt.Println()
	fmt.Println("Failed Login Information:")
	fmt.Printf("Time Range: %s\n", failedLogins.TimeRange)

	if err != nil {
		fmt.Println("Status: Unavailable")
		fmt.Println("Failed Authentication Events: Unknown")
		fmt.Println("Collection error:", err)
	} else {
		fmt.Printf("Status: %s\n", failedLogins.Status)

		if failedLogins.Status == "Available" {
			fmt.Printf(
				"Failed Authentication Events: %d\n",
				failedLogins.Count,
			)
		} else {
			fmt.Println("Failed Authentication Events: Unknown")
		}
	}

	fmt.Println()
	fmt.Println("Collection completed")
}
