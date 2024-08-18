//go:build windows
// +build windows

package free

import (
	"os/exec"
	"strconv"
	"strings"
)

func getMemoryInfo() (MemoryInfo, error) {
	cmd := exec.Command("wmic", "OS", "get", "FreePhysicalMemory,TotalVisibleMemorySize", "/Value")
	output, err := cmd.Output()
	if err != nil {
		return MemoryInfo{}, err
	}

	lines := strings.Split(string(output), "\n")
	var total, free uint64

	for _, line := range lines {
		if strings.HasPrefix(line, "FreePhysicalMemory=") {
			value := strings.TrimPrefix(line, "FreePhysicalMemory=")
			free, _ = strconv.ParseUint(strings.TrimSpace(value), 10, 64)
			free *= 1024 // Convert from KB to bytes
		} else if strings.HasPrefix(line, "TotalVisibleMemorySize=") {
			value := strings.TrimPrefix(line, "TotalVisibleMemorySize=")
			total, _ = strconv.ParseUint(strings.TrimSpace(value), 10, 64)
			total *= 1024 // Convert from KB to bytes
		}
	}

	used := total - free
	available := free // On Windows, free memory is considered available

	return MemoryInfo{
		Total:     total,
		Used:      used,
		Free:      free,
		Available: available,
	}, nil
}
