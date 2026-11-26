package collectors

import (
	"context"
	"os/exec"
	"strings"
	"time"
)

type FirewallInfo struct {
	Status string
}

func GetFirewallInfo() (FirewallInfo, error) {
	var fireWall FirewallInfo
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(
		ctx,
		"/usr/libexec/ApplicationFirewall/socketfilterfw",
		"--getglobalstate",
	)
	result, err := cmd.Output()

	if err != nil {
		fireWall.Status = "Unavailable"
		return fireWall, err
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
