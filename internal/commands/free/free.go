package free

import (
	"fmt"
)

type MemoryInfo struct {
	Total     uint64
	Used      uint64
	Free      uint64
	Available uint64
}

func Free(args []string) error {
	memInfo, err := getMemoryInfo()
	if err != nil {
		return fmt.Errorf("error getting memory info: %v", err)
	}

	printMemoryInfo(memInfo)
	return nil
}

func printMemoryInfo(info MemoryInfo) {
	fmt.Println("Memory Information:")
	fmt.Printf("Total:     %s\n", formatBytes(info.Total))
	fmt.Printf("Used:      %s\n", formatBytes(info.Used))
	fmt.Printf("Free:      %s\n", formatBytes(info.Free))
	fmt.Printf("Available: %s\n", formatBytes(info.Available))
}

func formatBytes(bytes uint64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}
