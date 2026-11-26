package collectors

import (
	"runtime"

	"github.com/shirou/gopsutil/v4/host"
)

type OSInfo struct {
	OS           string
	Architecture string
	Version      string
}

func GetOSInfo() (OSInfo, error) {
	var info OSInfo
	info.OS = runtime.GOOS
	info.Architecture = runtime.GOARCH
	hostInfo, err := host.Info()
	if err != nil {
		return OSInfo{}, err
	}
	info.Version = hostInfo.PlatformVersion
	return info, nil
}
