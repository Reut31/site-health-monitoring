package collectors

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type SecurityUpdateInfo struct {
	Status         string
	AvailableCount int
	Scope          string
}

func GetSecurityUpdateInfo() (SecurityUpdateInfo, error) {
	updateInfo := SecurityUpdateInfo{
		Status: "Unavailable",
		Scope:  "Unknown",
	}

	ctx, cancel := context.WithTimeout(
		context.Background(), 120*time.Second,
	)
	defer cancel()

	switch runtime.GOOS {

	case "darwin":
		updateInfo.Scope = "GeneralUpdates"

		cmd := exec.CommandContext(
			ctx,
			"softwareupdate",
			"--list",
		)

		result, err := cmd.CombinedOutput()

		if err != nil {
			return updateInfo, err
		}

		info := strings.ToLower(string(result))

		for _, line := range strings.Split(info, "\n") {
			if strings.HasPrefix(
				strings.TrimSpace(line), "* label:",
			) {
				updateInfo.AvailableCount++
			}
		}

		if updateInfo.AvailableCount == 0 &&
			!strings.Contains(info, "no new software available") {

			updateInfo.Status = "Unknown"
			return updateInfo, nil
		}

	case "windows":
		updateInfo.Scope = "SecurityUpdates"

		script := `
$ErrorActionPreference = 'Stop'

$session = New-Object -ComObject Microsoft.Update.Session
$searcher = $session.CreateUpdateSearcher()
$searcher.Online = $true

$criteria = "IsInstalled=0 and Type='Software' and CategoryIDs contains '0FA1201D-4330-4FA8-8AE9-B877473B6441'"

$result = $searcher.Search($criteria)

if ($result.ResultCode -ne 2) {
    throw 'Windows Update search did not complete successfully'
}

$result.Updates.Count
`

		cmd := exec.CommandContext(
			ctx,
			"powershell.exe",
			"-NoProfile",
			"-NonInteractive",
			"-Command",
			script,
		)

		result, err := cmd.Output()

		if err != nil {
			return updateInfo, fmt.Errorf(
				"Windows Update: %w", err,
			)
		}

		count, err := strconv.Atoi(
			strings.TrimSpace(string(result)),
		)

		if err != nil || count < 0 {
			return updateInfo, fmt.Errorf(
				"invalid Windows Update count",
			)
		}

		updateInfo.AvailableCount = count

	case "linux":

		if path, err := exec.LookPath(
			"/usr/lib/update-notifier/apt-check",
		); err == nil {

			updateInfo.Scope = "SecurityUpdates"

			cmd := exec.CommandContext(ctx, path)
			cmd.Env = append(os.Environ(), "LC_ALL=C")

			result, err := cmd.CombinedOutput()

			if err != nil {
				return updateInfo, err
			}

			info := strings.TrimSpace(string(result))
			parts := strings.Split(info, ";")

			if len(parts) != 2 {
				updateInfo.Status = "Unknown"
				return updateInfo, nil
			}

			total, totalErr := strconv.Atoi(parts[0])
			security, securityErr := strconv.Atoi(parts[1])

			if totalErr != nil || securityErr != nil ||
				security < 0 || total < security {

				updateInfo.Status = "Unknown"
				return updateInfo, nil
			}

			updateInfo.AvailableCount = security

		} else if path, err := exec.LookPath("dnf5"); err == nil {

			return collectDNFSecurityUpdates(
				ctx, path, updateInfo,
			)

		} else if path, err := exec.LookPath("dnf"); err == nil {

			return collectDNFSecurityUpdates(
				ctx, path, updateInfo,
			)

		} else if path, err := exec.LookPath("apt"); err == nil {

			updateInfo.Scope = "GeneralUpdates"

			cmd := exec.CommandContext(
				ctx,
				path,
				"list",
				"--upgradable",
			)

			cmd.Env = append(os.Environ(), "LC_ALL=C")

			result, err := cmd.Output()

			if err != nil {
				return updateInfo, err
			}

			for _, line := range strings.Split(
				string(result), "\n",
			) {
				line = strings.TrimSpace(line)

				if line == "" ||
					line == "Listing..." ||
					line == "Listing... Done" {
					continue
				}

				if !strings.Contains(line, "[upgradable from:") {
					updateInfo.AvailableCount = 0
					updateInfo.Status = "Unknown"
					return updateInfo, nil
				}

				updateInfo.AvailableCount++
			}

		} else {
			return updateInfo, nil
		}

	default:
		return updateInfo, nil
	}

	if updateInfo.AvailableCount > 0 {
		updateInfo.Status = "UpdatesAvailable"
	} else {
		updateInfo.Status = "UpToDate"
	}

	return updateInfo, nil
}

func collectDNFSecurityUpdates(
	ctx context.Context,
	path string,
	updateInfo SecurityUpdateInfo,
) (SecurityUpdateInfo, error) {

	updateInfo.Scope = "SecurityUpdates"

	cmd := exec.CommandContext(
		ctx,
		path,
		"-q",
		"repoquery",
		"--upgrades",
		"--security",
		"--queryformat",
		`PACKAGE:%{name}.%{arch}\n`,
	)

	cmd.Env = append(os.Environ(), "LC_ALL=C")

	result, err := cmd.Output()

	if err != nil {
		return updateInfo, err
	}

	packages := make(map[string]bool)

	for _, line := range strings.Split(string(result), "\n") {
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		if !strings.HasPrefix(line, "PACKAGE:") ||
			line == "PACKAGE:" {

			updateInfo.Status = "Unknown"
			return updateInfo, nil
		}

		name := strings.TrimPrefix(line, "PACKAGE:")
		packages[name] = true
	}

	updateInfo.AvailableCount = len(packages)

	if updateInfo.AvailableCount > 0 {
		updateInfo.Status = "UpdatesAvailable"
	} else {
		updateInfo.Status = "UpToDate"
	}

	return updateInfo, nil
}
