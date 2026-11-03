package main

import (
	"fmt"
	"site-health-monitoring/agent/collectors"
)

func main() {
	const Agent = "Agent started"
	fmt.Println(Agent)
	result := collectors.GetOSInfo()
	fmt.Println(result.OS)
	fmt.Println(result.Architecture)
}