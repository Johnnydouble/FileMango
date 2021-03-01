package scheduler

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/process"
	"strconv"
	"time"
)

func determinePoolSize(currentSize int) int {
	const threshold = 25.0
	const poolSizeMax = 10
	const poolSizeMin = 1

	//propose size
	proposedSize := 0
	if getInternalCpu() > threshold {
		proposedSize = currentSize - 1
	} else {
		proposedSize = currentSize + 1
	}

	//check reasonability
	if proposedSize > poolSizeMax {
		return poolSizeMax
	}
	if proposedSize < poolSizeMin {
		return poolSizeMin
	}
	return proposedSize
}

func getInternalCpu() float64 {
	total := 0.0
	for _, job := range pool.jobs {
		job := job                                 //todo: might not be necessary
		respChan := make(chan message)             //make a channel to wait for a response on
		job.input <- constructCpuGetMsg(&respChan) //construct a message
		select {
		case <-time.After(50 * time.Millisecond): //timeout if worker dies before it can respond, this should be rare
		case resp := <-respChan:
			individual, _ := strconv.ParseFloat(resp.Output.Pairs[0].Value, 64)
			total = total + individual
		}
	}
	return total
}

func constructCpuGetMsg(respChan *chan message) message {
	return message{
		Header: header{respChan},
		Input: input{
			Type: queryUsage,
			Data: "cpu",
		},
		Output: output{},
	}
}

//get total cpu usage as a percentage
func getTotalCpuUsage() float64 {
	p, _ := cpu.Percent(10*time.Second, false)
	val := p[0]
	return val
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
