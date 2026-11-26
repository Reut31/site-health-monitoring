package collectors

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type FailedLoginInfo struct {
	Count     int
	Status    string
	TimeRange string
}

func GetFailedLoginInfo() (FailedLoginInfo, error) {
	loginInfo := FailedLoginInfo{
		Status:    "Unavailable",
		TimeRange: "24h",
	}

	ctx, cancel := context.WithTimeout(
		context.Background(), 30*time.Second,
	)
	defer cancel()

	switch runtime.GOOS {

	case "darwin":
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
			return loginInfo, err
		}

		loginInfo.Count = strings.Count(
			strings.ToLower(string(result)),
			"authentication failed for",
		)

	case "windows":
		script := `
$ErrorActionPreference = 'Stop'

try {
    $count = (Get-WinEvent -FilterHashtable @{
        LogName = 'Security'
        Id = 4625
        StartTime = (Get-Date).AddHours(-24)
    } -ErrorAction Stop | Measure-Object).Count

    $count
} catch {
    if ($_.FullyQualifiedErrorId -like 'NoMatchingEventsFound,*') {
        0
    } else {
        throw
    }
}
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
			return loginInfo, fmt.Errorf(
				"Windows Security log: %w", err,
			)
		}

		count, err := strconv.Atoi(
			strings.TrimSpace(string(result)),
		)

		if err != nil || count < 0 {
			return loginInfo, fmt.Errorf(
				"invalid Windows failed logon count",
			)
		}

		loginInfo.Count = count

	case "linux":
		path, err := exec.LookPath("journalctl")

		if err != nil {
			return loginInfo, fmt.Errorf(
				"journalctl unavailable: %w", err,
			)
		}

		// Check journal accessibility
		probe, err := readLoginJournal(
			ctx, path, "--lines=1",
		)

		if err != nil {
			return loginInfo, err
		}

		if len(bytes.TrimSpace(probe)) == 0 ||
			strings.TrimSpace(string(probe)) == "-- No entries --" {

			return loginInfo, fmt.Errorf(
				"no readable system journal entries in the last 24h",
			)
		}

		// Collect authentication-related events
		result, err := readLoginJournal(
			ctx,
			path,
			"SYSLOG_FACILITY=4",
			"SYSLOG_FACILITY=10",
		)

		if err != nil {
			return loginInfo, err
		}

		count, err := countLinuxAuthFailures(result)

		if err != nil {
			return loginInfo, err
		}

		loginInfo.Count = count

	default:
		return loginInfo, nil
	}

	loginInfo.Status = "Available"

	return loginInfo, nil
}

func readLoginJournal(
	ctx context.Context,
	path string,
	extra ...string,
) ([]byte, error) {

	args := []string{
		"--system",
		"--since", "24 hours ago",
		"--no-pager",
		"--output=json",
	}

	args = append(args, extra...)

	cmd := exec.CommandContext(ctx, path, args...)
	cmd.Env = append(os.Environ(), "LC_ALL=C")

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	result, err := cmd.Output()

	if err != nil {
		return nil, fmt.Errorf(
			"journalctl: %w", err,
		)
	}

	if strings.TrimSpace(stderr.String()) != "" {
		return nil, fmt.Errorf(
			"journalctl reported warnings; log coverage cannot be verified",
		)
	}

	return result, nil
}

func countLinuxAuthFailures(result []byte) (int, error) {
	count := 0

	for _, line := range bytes.Split(result, []byte("\n")) {
		line = bytes.TrimSpace(line)

		if len(line) == 0 || string(line) == "-- No entries --" {
			continue
		}

		var entry struct {
			Message string `json:"MESSAGE"`
		}

		if err := json.Unmarshal(line, &entry); err != nil {
			return 0, fmt.Errorf(
				"invalid journal entry: %w", err,
			)
		}

		message := strings.ToLower(entry.Message)

		if strings.Contains(message, "failed password for ") ||
			strings.Contains(message, "authentication failure") ||
			strings.Contains(message, "authentication failed") {

			count++
		}
	}

	return count, nil
}
