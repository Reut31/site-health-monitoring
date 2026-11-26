package collectors

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/shirou/gopsutil/v4/mem"
)

type MemoryInfo struct {
	Total       uint64
	Used        uint64
	UsedPercent float64
}

func GetMemoryInfo() (MemoryInfo, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(), 5*time.Second,
	)
	defer cancel()

	vm, err := mem.VirtualMemoryWithContext(ctx)

	if err != nil {
		return MemoryInfo{}, err
	}

	if vm == nil {
		return MemoryInfo{}, fmt.Errorf(
			"no memory information received",
		)
	}

	if vm.Total == 0 {
		return MemoryInfo{}, fmt.Errorf(
			"total memory cannot be zero",
		)
	}

	if vm.Used > vm.Total {
		return MemoryInfo{}, fmt.Errorf(
			"used memory exceeds total memory",
		)
	}

	if math.IsNaN(vm.UsedPercent) ||
		math.IsInf(vm.UsedPercent, 0) ||
		vm.UsedPercent < 0 ||
		vm.UsedPercent > 100 {

		return MemoryInfo{}, fmt.Errorf(
			"invalid memory usage: %v",
			vm.UsedPercent,
		)
	}

	return MemoryInfo{
		Total:       vm.Total,
		Used:        vm.Used,
		UsedPercent: vm.UsedPercent,
	}, nil
}
