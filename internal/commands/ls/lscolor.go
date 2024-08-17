package ls

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// LsColor lists directory contents with colors (for file types) and detailed information
func LsColor(args []string) error {
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		mode := file.Mode()
		size := file.Size()
		modTime := file.ModTime().Format(time.RFC822)
		name := file.Name()

		// Determine color based on file type
		color := "\033[0m" // Default color (reset)
		if file.IsDir() {
			color = "\033[1;34m" // Blue for directories
		} else if mode&0111 != 0 {
			color = "\033[1;32m" // Green for executable files
		} else if strings.HasPrefix(strings.ToLower(filepath.Ext(name)), ".") {
			color = "\033[1;37m" // White for hidden files
		}

		// Create a string that mimics the output of 'ls -l' on Unix
		var fileInfo string
		if runtime.GOOS == "windows" {
			// On Windows, we'll use a simplified format
			fileType := "f"
			if file.IsDir() {
				fileType = "d"
			}
			fileInfo = fmt.Sprintf("%s %10d %s %s%s%s", fileType, size, modTime, color, name, "\033[0m")
		} else {
			// On Unix-like systems, we'll try to mimic 'ls -l' more closely
			perms := mode.String()
			owner := getOwner(file)
			group := getGroup(file)
			fileInfo = fmt.Sprintf("%s %s %s %8d %s %s%s%s", perms, owner, group, size, modTime, color, name, "\033[0m")
		}

		fmt.Println(fileInfo)
	}

	return nil
}
