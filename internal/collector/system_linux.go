//go:build linux

package collector

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

// LinuxSystemCollector contains the state needed for tracking CPU usage.
type linuxSystemState struct {
	mu          sync.Mutex
	lastUser    uint64
	lastNice    uint64
	lastSystem  uint64
	lastIdle    uint64
	lastIowait  uint64
	lastIrq     uint64
	lastSoftirq uint64
	lastSteal   uint64
	hasPrev     bool
}

var state linuxSystemState

// Collect collects memory, load averages, and CPU utilization on Linux.
func (s *SystemCollector) Collect(ctx context.Context) (map[string]any, error) {
	metrics := make(map[string]any)

	// Collect Load Average
	if load, err := readLoadAvg(); err == nil {
		metrics["load"] = load
	}

	// Collect Memory
	if mem, err := readMemInfo(); err == nil {
		metrics["memory"] = mem
	}

	// Collect CPU
	if cpu, err := readCPUUsage(); err == nil {
		metrics["cpu"] = cpu
	}

	return metrics, nil
}

func readLoadAvg() (map[string]float64, error) {
	data, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return nil, err
	}
	fields := strings.Fields(string(data))
	if len(fields) < 3 {
		return nil, fmt.Errorf("invalid loadavg format")
	}

	load1, err1 := strconv.ParseFloat(fields[0], 64)
	load5, err2 := strconv.ParseFloat(fields[1], 64)
	load15, err3 := strconv.ParseFloat(fields[2], 64)
	if err1 != nil || err2 != nil || err3 != nil {
		return nil, fmt.Errorf("failed to parse load averages")
	}

	return map[string]float64{
		"load1":  load1,
		"load5":  load5,
		"load15": load15,
	}, nil
}

func readMemInfo() (map[string]uint64, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	mem := make(map[string]uint64)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		key := strings.TrimSuffix(fields[0], ":")
		val, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			continue
		}
		// Convert kB to Bytes
		bytesVal := val * 1024

		switch key {
		case "MemTotal":
			mem["total_bytes"] = bytesVal
		case "MemFree":
			mem["free_bytes"] = bytesVal
		case "MemAvailable":
			mem["available_bytes"] = bytesVal
		}
	}
	return mem, nil
}

func readCPUUsage() (map[string]float64, error) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return nil, fmt.Errorf("empty /proc/stat")
	}
	line := scanner.Text()
	fields := strings.Fields(line)
	if len(fields) < 5 || fields[0] != "cpu" {
		return nil, fmt.Errorf("invalid cpu format in /proc/stat")
	}

	var user, nice, system, idle, iowait, irq, softirq, steal uint64
	var errs [8]error
	user, errs[0] = strconv.ParseUint(fields[1], 10, 64)
	nice, errs[1] = strconv.ParseUint(fields[2], 10, 64)
	system, errs[2] = strconv.ParseUint(fields[3], 10, 64)
	idle, errs[3] = strconv.ParseUint(fields[4], 10, 64)
	if len(fields) > 5 {
		iowait, errs[4] = strconv.ParseUint(fields[5], 10, 64)
	}
	if len(fields) > 6 {
		irq, errs[5] = strconv.ParseUint(fields[6], 10, 64)
	}
	if len(fields) > 7 {
		softirq, errs[6] = strconv.ParseUint(fields[7], 10, 64)
	}
	if len(fields) > 8 {
		steal, errs[7] = strconv.ParseUint(fields[8], 10, 64)
	}

	for _, e := range errs {
		if e != nil {
			return nil, fmt.Errorf("failed to parse cpu values: %w", e)
		}
	}

	state.mu.Lock()
	defer state.mu.Unlock()

	if !state.hasPrev {
		state.lastUser = user
		state.lastNice = nice
		state.lastSystem = system
		state.lastIdle = idle
		state.lastIowait = iowait
		state.lastIrq = irq
		state.lastSoftirq = softirq
		state.lastSteal = steal
		state.hasPrev = true
		return map[string]float64{"utilization_percent": 0.0}, nil
	}

	prevIdle := state.lastIdle + state.lastIowait
	idleNow := idle + iowait

	prevNonIdle := state.lastUser + state.lastNice + state.lastSystem + state.lastIrq + state.lastSoftirq + state.lastSteal
	nonIdleNow := user + nice + system + irq + softirq + steal

	prevTotal := prevIdle + prevNonIdle
	totalNow := idleNow + nonIdleNow

	totalDiff := totalNow - prevTotal
	idleDiff := idleNow - prevIdle

	state.lastUser = user
	state.lastNice = nice
	state.lastSystem = system
	state.lastIdle = idle
	state.lastIowait = iowait
	state.lastIrq = irq
	state.lastSoftirq = softirq
	state.lastSteal = steal

	if totalDiff == 0 {
		return map[string]float64{"utilization_percent": 0.0}, nil
	}

	utilization := float64(totalDiff-idleDiff) / float64(totalDiff) * 100.0
	return map[string]float64{"utilization_percent": utilization}, nil
}
