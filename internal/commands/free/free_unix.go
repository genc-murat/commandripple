//go:build linux
// +build linux

package free

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func getMemoryInfo() (MemoryInfo, error) {
	cmd := exec.Command("free", "-b")
	output, err := cmd.Output()
	if err != nil {
		return MemoryInfo{}, err
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return MemoryInfo{}, fmt.Errorf("unexpected output format")
	}

	fields := strings.Fields(lines[1])
	if len(fields) < 7 {
		return MemoryInfo{}, fmt.Errorf("unexpected output format")
	}

	total, _ := strconv.ParseUint(fields[1], 10, 64)
	used, _ := strconv.ParseUint(fields[2], 10, 64)
	free, _ := strconv.ParseUint(fields[3], 10, 64)
	available, _ := strconv.ParseUint(fields[6], 10, 64)

	return MemoryInfo{
		Total:     total,
		Used:      used,
		Free:      free,
		Available: available,
	}, nil
}
