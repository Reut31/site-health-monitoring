package collectors

import (
	"os/exec"
	"strings"
)

type FirewallInfo struct {
	Status string
}

func GetFirewallInfo() (FirewallInfo, error) {
	var fireWall FirewallInfo

	cmd := exec.Command(
		"/usr/libexec/ApplicationFirewall/socketfilterfw",
		"--getglobalstate",
	)
	result, err := cmd.Output()

	if err != nil {
		return FirewallInfo{}, err
	}
	info := string(result)
	if strings.Contains(info, "State = 1") {
		fireWall.Status = "Enabled"
	} else if strings.Contains(info, "State = 0") {
		fireWall.Status = "Disabled"
	} else {
		fireWall.Status = "Unknown"
	}
	return fireWall, nil
}
