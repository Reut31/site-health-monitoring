package collectors

import (
	"context"
	"os/exec"
	"strings"
	"time"
)

type SecurityUpdateInfo struct {
	Status         string
	AvailableCount int
}

func GetSecurityUpdateInfo() (SecurityUpdateInfo, error) {
	var updateInfo SecurityUpdateInfo
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	cmd := exec.CommandContext(
		ctx,
		"softwareupdate",
		"--list",
	)
	result, err := cmd.Output()

	if err != nil {
		updateInfo.Status = "Unavailable"
		return updateInfo, err
	}
	info := string(result)
	updateInfo.AvailableCount = strings.Count(info, "* Label:")
	if strings.Contains(strings.ToLower(info), "software update found the following") {
		updateInfo.Status = "UpdatesAvailable"
	} else if strings.Contains(strings.ToLower(info), "no new software available") {
		updateInfo.Status = "UpToDate"
	} else {
		updateInfo.Status = "Unknown"
	}
	return updateInfo, nil
}
