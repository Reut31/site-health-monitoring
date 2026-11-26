package collectors

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type AntivirusInfo struct {
	Name   string
	Status string
}

func GetAntivirusInfo() (AntivirusInfo, error) {
	antivirus := AntivirusInfo{
		Name:   "Unknown",
		Status: "Unknown",
	}

	ctx, cancel := context.WithTimeout(
		context.Background(), 15*time.Second,
	)
	defer cancel()

	switch runtime.GOOS {

	case "darwin":
		cmd := exec.CommandContext(
			ctx,
			"pkgutil",
			"--pkgs",
		)

		result, err := cmd.Output()

		if err != nil {
			antivirus.Status = "Unavailable"
			return antivirus, err
		}

		info := strings.ToLower(string(result))

		if strings.Contains(info, "xprotect") {
			antivirus.Name = "XProtect"
		}

	case "windows":
		script := `
[Console]::OutputEncoding = [System.Text.UTF8Encoding]::new($false)

$defender = $null
$products = @()
$names = @()

try {
    $defender = Get-MpComputerStatus -ErrorAction Stop
} catch {
    $defender = $null
}

try {
    $products = @(Get-CimInstance ` + "`" + `
        -Namespace root/SecurityCenter2 ` + "`" + `
        -ClassName AntivirusProduct ` + "`" + `
        -ErrorAction Stop)

    $names = @($products | ForEach-Object {
        $_.displayName
    } | Where-Object { $_ } | Sort-Object -Unique)
} catch {
    $names = @()
}

if ($defender -and
    $defender.AntivirusEnabled -eq $true -and
    $defender.RealTimeProtectionEnabled -eq $true -and
    $defender.AMRunningMode -eq 'Normal') {

    $result = @{
        Name = 'Microsoft Defender'
        Status = 'Enabled'
    }

} elseif ($names.Count -gt 0) {

    $result = @{
        Name = ($names -join ', ')
        Status = 'Unknown'
    }

} elseif ($defender) {

    $status = 'Unknown'

    if ($defender.AntivirusEnabled -eq $false) {
        $status = 'Disabled'
    } elseif ($defender.AntivirusEnabled -eq $true -and
              $defender.RealTimeProtectionEnabled -eq $false) {
        $status = 'PartiallyEnabled'
    }

    $result = @{
        Name = 'Microsoft Defender'
        Status = $status
    }

} else {
    throw 'Unable to collect antivirus information'
}

ConvertTo-Json -InputObject $result -Compress
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
			antivirus.Status = "Unavailable"
			return antivirus, fmt.Errorf(
				"Windows antivirus collection: %w", err,
			)
		}

		if err := json.Unmarshal(result, &antivirus); err != nil {
			antivirus.Status = "Unavailable"
			return antivirus, fmt.Errorf(
				"Windows antivirus parsing: %w", err,
			)
		}

	case "linux":

		if path, err := exec.LookPath("mdatp"); err == nil {

			antivirus.Name = "Microsoft Defender for Endpoint"

			cmd := exec.CommandContext(
				ctx,
				path,
				"health",
				"--field",
				"real_time_protection_enabled",
			)

			result, err := cmd.Output()

			if err != nil {
				antivirus.Status = "Unavailable"
				return antivirus, err
			}

			info := strings.ToLower(
				strings.TrimSpace(string(result)),
			)

			if info == "true" {
				antivirus.Status = "Enabled"
			} else if info == "false" {
				antivirus.Status = "PartiallyEnabled"
			}

		} else if _, err := exec.LookPath("clamscan"); err == nil {

			antivirus.Name = "ClamAV"
			antivirus.Status = "Unknown"

		} else if _, err := exec.LookPath("clamd"); err == nil {

			antivirus.Name = "ClamAV"
			antivirus.Status = "Unknown"

		} else if _, err := exec.LookPath("clamonacc"); err == nil {

			antivirus.Name = "ClamAV"
			antivirus.Status = "Unknown"

		} else {
			antivirus.Status = "Unavailable"
		}

	default:
		antivirus.Status = "Unavailable"
	}

	return antivirus, nil
}
