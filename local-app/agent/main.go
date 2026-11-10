package main

import (
	"fmt"
	"site-health-monitoring/agent/collectors"
)

func main() {
	fmt.Println("Agent started")
	fmt.Println()

	result1 := collectors.GetOSInfo()

	fmt.Println("System Information:")
	fmt.Printf("OS: %s\n", result1.OS)
	fmt.Printf("Architecture: %s\n", result1.Architecture)
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
}