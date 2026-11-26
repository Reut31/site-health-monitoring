package collectors

import (
	"context"
	"os/exec"
	"strings"
	"time"
)

type FailedLoginInfo struct {
	Count     int
	Status    string
	TimeRange string
}

func GetFailedLoginInfo() (FailedLoginInfo, error) {
	var loginInfo FailedLoginInfo
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	loginInfo.TimeRange = "24h"

	cmd := exec.CommandContext(
		ctx,
		"log",
		"show",
		"--last", "24h",
		"--info",
		"--style", "compact",
		"--predicate", `eventMessage CONTAINS[c] "Authentication failed"`,
	)
	result, err := cmd.Output()

	if err != nil {
		loginInfo.Status = "Unavailable"
		return loginInfo, err
	}
	info := string(result)
	loginInfo.Count = strings.Count(info, "Authentication failed for")
	loginInfo.Status = "Available"
	return loginInfo, nil
}
