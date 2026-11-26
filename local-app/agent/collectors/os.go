package collectors

import (
	"context"
	"runtime"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v4/host"
)

type OSInfo struct {
	OS           string
	Architecture string
	Version      string
}

func GetOSInfo() (OSInfo, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(), 5*time.Second,
	)
	defer cancel()

	info := OSInfo{
		OS:           runtime.GOOS,
		Architecture: runtime.GOARCH,
		Version:      "Unknown",
	}

	_, _, version, err := host.PlatformInformationWithContext(ctx)

	if err != nil {
		return info, err
	}

	version = strings.TrimSpace(version)

	if version != "" {
		info.Version = version
	}

	return info, nil
}
