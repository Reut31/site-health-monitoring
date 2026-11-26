package collectors

import (
	"context"
	"os/exec"
	"strings"
	"time"
)

type AntivirusInfo struct {
	Name   string
	Status string
}

func GetAntivirusInfo() (AntivirusInfo, error) {
	var antivirus AntivirusInfo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(
		ctx,
		"pkgutil",
		"--pkgs",
	)
	result, err := cmd.Output()

	if err != nil {
		antivirus.Name = "Unknown"
		antivirus.Status = "Unavailable"
		return antivirus, err
	}
	info := string(result)
	if strings.Contains(strings.ToLower(info), "xprotect") {
		antivirus.Name = "XProtect"
		antivirus.Status = "Unknown"
	} else {
		antivirus.Name = "Unknown"
		antivirus.Status = "Unknown"
	}
	return antivirus, nil
}
