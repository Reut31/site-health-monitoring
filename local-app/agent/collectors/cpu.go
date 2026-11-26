package collectors

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
)

type CPUInfo struct {
	UsagePercent float64
}

func GetCPUInfo() (CPUInfo, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(), 5*time.Second,
	)
	defer cancel()

	values, err := cpu.PercentWithContext(
		ctx,
		time.Second,
		false,
	)

	if err != nil {
		return CPUInfo{}, err
	}

	if len(values) != 1 {
		return CPUInfo{}, fmt.Errorf(
			"expected one CPU usage value, received %d",
			len(values),
		)
	}

	usage := values[0]

	if math.IsNaN(usage) || math.IsInf(usage, 0) ||
		usage < 0 || usage > 100 {

		return CPUInfo{}, fmt.Errorf(
			"invalid CPU usage: %v", usage,
		)
	}

	return CPUInfo{
		UsagePercent: usage,
	}, nil
}
