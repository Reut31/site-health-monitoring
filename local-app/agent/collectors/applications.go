package collectors

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strings"
	"time"
)

type InstalledApplication struct {
	Name string
}

func GetInstalledApplications() ([]InstalledApplication, error) {
	names := make([]string, 0)

	ctx, cancel := context.WithTimeout(
		context.Background(), 15*time.Second,
	)
	defer cancel()

	switch runtime.GOOS {

	case "darwin":
		result, err := os.ReadDir("/Applications")

		if err != nil {
			return nil, err
		}

		for _, application := range result {
			if application.IsDir() &&
				strings.HasSuffix(application.Name(), ".app") {
				names = append(names, application.Name())
			}
		}

	case "windows":
		script := `
$ErrorActionPreference = 'Stop'
[Console]::OutputEncoding = [System.Text.UTF8Encoding]::new($false)

$paths = @(
    'HKLM:\Software\Microsoft\Windows\CurrentVersion\Uninstall',
    'HKLM:\Software\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall',
    'HKCU:\Software\Microsoft\Windows\CurrentVersion\Uninstall'
)

$names = @(foreach ($path in $paths) {
    if (Test-Path -LiteralPath $path) {
        Get-ChildItem -LiteralPath $path | ForEach-Object {
            $app = Get-ItemProperty -LiteralPath $_.PSPath

            if ($app.DisplayName) {
                [string]$app.DisplayName
            }
        }
    }
})

ConvertTo-Json -InputObject $names -Compress
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
			return nil, fmt.Errorf(
				"Windows application collection: %w", err,
			)
		}

		if err := json.Unmarshal(result, &names); err != nil {
			return nil, fmt.Errorf(
				"Windows application parsing: %w", err,
			)
		}

	case "linux":

		if path, err := exec.LookPath("dpkg-query"); err == nil {

			cmd := exec.CommandContext(
				ctx,
				path,
				"-W",
				`-f=${db:Status-Status}\t${binary:Package}\n`,
			)

			result, err := cmd.Output()

			if err != nil {
				return nil, fmt.Errorf("dpkg-query: %w", err)
			}

			for _, line := range strings.Split(string(result), "\n") {
				fields := strings.SplitN(line, "\t", 2)

				if len(fields) == 2 && fields[0] == "installed" {
					names = append(names, fields[1])
				}
			}

		} else if path, err := exec.LookPath("rpm"); err == nil {

			cmd := exec.CommandContext(
				ctx,
				path,
				"-qa",
				"--queryformat",
				`%{NAME}\n`,
			)

			result, err := cmd.Output()

			if err != nil {
				return nil, fmt.Errorf("rpm: %w", err)
			}

			names = strings.Split(string(result), "\n")

		} else {
			return nil, fmt.Errorf(
				"no supported Linux package manager found",
			)
		}

	default:
		return nil, fmt.Errorf(
			"installed application collection unsupported on %s",
			runtime.GOOS,
		)
	}

	// Remove duplicates and empty names
	applications := make([]InstalledApplication, 0)
	seen := make(map[string]bool)

	for _, name := range names {
		name = strings.TrimSpace(name)

		if name == "" || seen[name] {
			continue
		}

		seen[name] = true

		applications = append(
			applications,
			InstalledApplication{Name: name},
		)
	}

	// Sort applications by name
	sort.Slice(applications, func(i, j int) bool {
		return applications[i].Name < applications[j].Name
	})

	return applications, nil
}
