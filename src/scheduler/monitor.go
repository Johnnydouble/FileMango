package scheduler

import (
	"github.com/shirou/gopsutil/v3/process"
)

func determinePoolSize() {

}

//returns the cpu usage percentage for a given process
func getProcessCpu(pid int) float64 {
	p, _ := process.NewProcess(int32(pid))
	val, _ := p.CPUPercent()
	return val
}

//returns the memory usage percentage for a given process
//not as precise as getProcessCpu() because the library does not natively produce float64 here and it must be casted
func getProcessMem(pid int) float64 {
	p, _ := process.NewProcess(int32(pid))
	val, _ := p.MemoryPercent()
	return float64(val)
}
