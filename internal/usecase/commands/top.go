package commands

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/process"
	"math"
	"os"
	"proxy-server-with-tg-admin/internal/helper"
	"runtime"
	"strings"
	"time"
)

type top struct {
	ip   string
	port uint
}

func (c *top) Id() string {
	return "top"
}

func (c *top) Arguments() []string {
	return []string{}
}

func (c *top) Run(args ...string) (string, error) {
	return c.top() + c.selfStatus() + c.runtime(), nil
}

func (c *top) top() string {
	pid := os.Getpid()
	if pid > math.MaxInt32 || pid < math.MinInt32 {
		return fmt.Sprintf("%d pid is out of range", pid)
	}

	proc, err := process.NewProcess(int32(pid))
	if err != nil {
		return ""
	}

	var total float64

	samples := 5

	for range samples {
		cpu, err := proc.CPUPercent()
		if err != nil {
			break
		}

		total += cpu

		time.Sleep(100 * time.Millisecond)
	}

	averageCPU := total / float64(samples)
	memInfo, _ := proc.MemoryInfo()

	return fmt.Sprintf("%s:%d\n", c.ip, c.port) +
		fmt.Sprintf("CPU: %.2f%%\n", averageCPU) +
		fmt.Sprintf("RSS: %s\n", helper.BytesFormat(memInfo.RSS)) +
		fmt.Sprintf("VMS: %s\n", helper.BytesFormat(memInfo.VMS))
}

func (c *top) selfStatus() string {
	data, _ := os.ReadFile("/proc/self/status")
	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "Threads:") {
			return line + "\n"
		}
	}

	return ""
}

func (c *top) runtime() string {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return fmt.Sprintf("NumGoroutine = %v\n", runtime.NumGoroutine()) +
		fmt.Sprintf("Sys = %s\n", helper.BytesFormat(m.Sys)) +
		fmt.Sprintf("Alloc = %s\n", helper.BytesFormat(m.Alloc)) +
		fmt.Sprintf("TotalAlloc = %s\n", helper.BytesFormat(m.TotalAlloc)) +
		fmt.Sprintf("HeapInuse = %s\n", helper.BytesFormat(m.HeapInuse)) +
		fmt.Sprintf("HeapAlloc  = %s\n", helper.BytesFormat(m.HeapAlloc)) +
		fmt.Sprintf("NumGC = %v\n", m.NumGC)
}
