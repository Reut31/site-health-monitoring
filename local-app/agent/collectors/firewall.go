package collectors

import (
	"context"
	"errors"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type FirewallInfo struct {
	Status string
}

func GetFirewallInfo() (FirewallInfo, error) {
	var fireWall FirewallInfo

	ctx, cancel := context.WithTimeout(
		context.Background(), 5*time.Second,
	)
	defer cancel()

	switch runtime.GOOS {

	case "darwin":
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

	case "windows":
		cmd := exec.CommandContext(
			ctx,
			"powershell.exe",
			"-NoProfile",
			"-NonInteractive",
			"-Command",
			"Get-NetFirewallProfile -PolicyStore ActiveStore -ErrorAction Stop | Select-Object -ExpandProperty Enabled",
		)

		result, err := cmd.Output()

		if err != nil {
			fireWall.Status = "Unavailable"
			return fireWall, err
		}

		info := strings.ToLower(string(result))
		profiles := strings.Fields(info)

		enabledCount := 0
		disabledCount := 0

		for _, profile := range profiles {
			if profile == "true" {
				enabledCount++
			} else if profile == "false" {
				disabledCount++
			}
		}

		if len(profiles) != 3 {
			fireWall.Status = "Unknown"
		} else if enabledCount == 3 {
			fireWall.Status = "Enabled"
		} else if disabledCount == 3 {
			fireWall.Status = "Disabled"
		} else if enabledCount+disabledCount == 3 {
			fireWall.Status = "PartiallyEnabled"
		} else {
			fireWall.Status = "Unknown"
		}

	case "linux":
		fireWall.Status = "Unavailable"
		var collectionErr error

		// Check UFW
		if _, err := exec.LookPath("ufw"); err == nil {

			cmd := exec.CommandContext(
				ctx, "ufw", "status",
			)

			result, err := cmd.Output()

			if err != nil {
				collectionErr = err
			} else {
				info := strings.TrimSpace(string(result))

				if strings.HasPrefix(info, "Status: active") {
					fireWall.Status = "Enabled"
					return fireWall, nil
				}

				fireWall.Status = "Unknown"
			}
		}

		// Check firewalld
		if _, err := exec.LookPath("firewall-cmd"); err == nil {

			cmd := exec.CommandContext(
				ctx, "firewall-cmd", "--state",
			)

			result, err := cmd.Output()

			if err != nil {
				exitErr, ok := err.(*exec.ExitError)

				if ok && exitErr.ExitCode() == 252 {
					fireWall.Status = "Unknown"
				} else {
					collectionErr = errors.Join(collectionErr, err)
				}
			} else {
				info := strings.TrimSpace(string(result))

				if info == "running" {
					fireWall.Status = "Enabled"
					return fireWall, nil
				}

				fireWall.Status = "Unknown"
			}
		}

		if collectionErr != nil {
			fireWall.Status = "Unavailable"
			return fireWall, collectionErr
		}

	default:
		fireWall.Status = "Unavailable"
	}

	return fireWall, nil
}
