package collectors

import (
	"github.com/shirou/gopsutil/v4/cpu"
	"time"
)

type CPUInfo struct {
	UsagePercent float64
}
func GetCPUInfo() (CPUInfo, error){
	var cpuinf CPUInfo
	values, err := cpu.Percent(time.Second, false)
	if err != nil {
		return CPUInfo{}, err
	}
	cpuinf.UsagePercent = values[0]
	return  cpuinf, nil
}