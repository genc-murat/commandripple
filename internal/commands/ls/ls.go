package ls

import (
	"fmt"
	"io/ioutil"
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

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		mode := file.Mode()
		size := file.Size()
		modTime := file.ModTime().Format(time.RFC822)
		name := file.Name()

		// Create a string that mimics the output of 'ls -l' on Unix
		var fileInfo string
		if runtime.GOOS == "windows" {
			// On Windows, we'll use a simplified format
			fileType := "f"
			if file.IsDir() {
				fileType = "d"
			}
			fileInfo = fmt.Sprintf("%s %10d %s %s", fileType, size, modTime, name)
		} else {
			// On Unix-like systems, we'll try to mimic 'ls -l' more closely
			perms := mode.String()
			owner := getOwner(file)
			group := getGroup(file)
			fileInfo = fmt.Sprintf("%s %s %s %8d %s %s", perms, owner, group, size, modTime, name)
		}

		fmt.Println(fileInfo)
	}

	return nil
}

func getOwner(file os.FileInfo) string {
	if runtime.GOOS == "windows" {
		return "owner"
	}
	return "owner" // Replace with actual owner retrieval for Unix systems
}

func getGroup(file os.FileInfo) string {
	if runtime.GOOS == "windows" {
		return "group"
	}
	return "group" // Replace with actual group retrieval for Unix systems
}
