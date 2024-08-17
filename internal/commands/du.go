package commands

import (
	"fmt"
	"io/fs"
	"path/filepath"
)

// Du estimates file space usage of a directory
func Du(args []string) error {
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	var totalSize int64
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing %s: %v\n", path, err)
			return nil // Continue walking
		}
		if !info.IsDir() {
			totalSize += info.Size()
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking directory: %v", err)
	}

	// Convert bytes to human-readable format
	sizeStr := formatSize(totalSize)
	fmt.Printf("%s\t%s\n", sizeStr, dir)

	return nil
}

func formatSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}
