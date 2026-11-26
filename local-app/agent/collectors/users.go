package collectors

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v4/host"
)

type LoggedInUser struct {
	Username string
	Terminal string
	Host     string
}

func GetLoggedInUsers() ([]LoggedInUser, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(), 10*time.Second,
	)
	defer cancel()

	users := make([]LoggedInUser, 0)

	result, err := host.UsersWithContext(ctx)

	if err != nil {
		if runtime.GOOS == "windows" {
			return getWindowsLoggedInUsers(ctx)
		}

		return nil, err
	}

	seen := make(map[LoggedInUser]bool)

	for _, user := range result {
		newUser := LoggedInUser{
			Username: strings.TrimSpace(user.User),
			Terminal: strings.TrimSpace(user.Terminal),
			Host:     strings.TrimSpace(user.Host),
		}

		if newUser.Username == "" {
			continue
		}

		if seen[newUser] {
			continue
		}

		seen[newUser] = true
		users = append(users, newUser)
	}

	return users, nil
}

func getWindowsLoggedInUsers(
	ctx context.Context,
) ([]LoggedInUser, error) {

	users := make([]LoggedInUser, 0)

	script := `
$ErrorActionPreference = 'Stop'
[Console]::OutputEncoding = [System.Text.UTF8Encoding]::new($false)

$user = (Get-CimInstance -ClassName Win32_ComputerSystem).UserName

ConvertTo-Json -InputObject ([string]$user) -Compress
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
			"Windows user collection: %w", err,
		)
	}

	var username string

	if err := json.Unmarshal(result, &username); err != nil {
		return nil, fmt.Errorf(
			"Windows user parsing: %w", err,
		)
	}

	username = strings.TrimSpace(username)

	if username != "" {
		users = append(users, LoggedInUser{
			Username: username,
			Terminal: "console",
			Host:     "",
		})
	}

	return users, nil
}
