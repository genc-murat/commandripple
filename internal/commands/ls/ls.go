package ls

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

// Ls lists directory contents with detailed file information
func Ls(args []string) error {
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return err
		}

		mode := info.Mode()
		size := info.Size()
		modTime := info.ModTime().Format(time.RFC822)
		name := info.Name()

		// Create a string that mimics the output of 'ls -l' on Unix
		var fileInfo string
		if runtime.GOOS == "windows" {
			// On Windows, we'll use a simplified format
			fileType := "f"
			if entry.IsDir() {
				fileType = "d"
			}
			fileInfo = fmt.Sprintf("%s %10d %s %s", fileType, size, modTime, name)
		} else {
			// On Unix-like systems, we'll try to mimic 'ls -l' more closely
			perms := mode.String()
			owner := getOwner(info)
			group := getGroup(info)
			fileInfo = fmt.Sprintf("%s %s %s %8d %s %s", perms, owner, group, size, modTime, name)
		}

		fmt.Println(fileInfo)
	}

	return nil
}
