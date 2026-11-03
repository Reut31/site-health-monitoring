package collectors

import (
	"runtime"
)

type OSInfo struct {
	OS string
  	Architecture string
}

func GetOSInfo() (OSInfo) {
	var info OSInfo
	info.OS = runtime.GOOS
	info.Architecture = runtime.GOARCH
	return info
}