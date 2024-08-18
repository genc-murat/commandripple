//go:build darwin
// +build darwin

package free

import (
	"os/exec"
	"strconv"
	"strings"
)

func getMemoryInfo() (MemoryInfo, error) {
	cmd := exec.Command("vm_stat")
	output, err := cmd.Output()
	if err != nil {
		return MemoryInfo{}, err
	}

	lines := strings.Split(string(output), "\n")
	pageSize := uint64(4096) // Default page size for macOS

	var free, active, inactive, wired uint64

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		value, err := strconv.ParseUint(strings.TrimRight(fields[1], "."), 10, 64)
		if err != nil {
			continue
		}

		switch fields[0] {
		case "Pages free:":
			free = value * pageSize
		case "Pages active:":
			active = value * pageSize
		case "Pages inactive:":
			inactive = value * pageSize
		case "Pages wired down:":
			wired = value * pageSize
		}
	}

	total := free + active + inactive + wired
	used := active + inactive + wired
	available := free + inactive

	return MemoryInfo{
		Total:     total,
		Used:      used,
		Free:      free,
		Available: available,
	}, nil
}
