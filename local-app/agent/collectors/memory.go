package collectors

import (
	"github.com/shirou/gopsutil/v4/mem"
)

type MemoryInfo struct {
	Total uint64
  	Used uint64
	UsedPercent float64
}

func GetMemoryInfo() (MemoryInfo, error){
	var memory MemoryInfo
	vm, err := mem.VirtualMemory()
	if( err != nil){
		return MemoryInfo{}, err
	}
	memory.Total = vm.Total
	memory.Used = vm.Used
	memory.UsedPercent = vm.UsedPercent
	return  memory , nil
}