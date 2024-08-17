package commands

import (
	"fmt"
	"os"
	"time"
)

func Stat(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: stat [file]")
	}

	filePath := args[0]
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	fmt.Printf("  File: %s\n", filePath)
	fmt.Printf("  Size: %d bytes\n", fileInfo.Size())
	fmt.Printf("  Mode: %s\n", fileInfo.Mode())
	fmt.Printf("  Modified: %s\n", fileInfo.ModTime().Format(time.RFC1123))

	printDetailedStats(filePath, fileInfo)

	return nil
}
