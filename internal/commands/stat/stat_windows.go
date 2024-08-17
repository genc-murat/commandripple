//go:build windows
// +build windows

package stat

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func printDetailedStats(filePath string, fileInfo os.FileInfo) {
	cmd := exec.Command("powershell", "-Command",
		fmt.Sprintf("(Get-Item '%s').CreationTime, (Get-Item '%s').LastAccessTime", filePath, filePath))
	output, err := cmd.Output()
	if err == nil {
		times := strings.Split(strings.TrimSpace(string(output)), "\n")
		if len(times) == 2 {
			fmt.Printf("  Created: %s\n", strings.TrimSpace(times[0]))
			fmt.Printf("  Accessed: %s\n", strings.TrimSpace(times[1]))
		}
	}
}
