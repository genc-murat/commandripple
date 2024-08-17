package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func findExecutable(command string) (string, error) {
	// If the command contains a path separator, check if it's an absolute path
	if strings.ContainsRune(command, os.PathSeparator) {
		path, err := filepath.Abs(command)
		if err == nil {
			if isExecutable(path) {
				return path, nil
			}
		}
		return "", fmt.Errorf("command not found: %s", command)
	}

	// Check in the current directory first
	path := checkCurrentDir(command)
	if path != "" {
		return path, nil
	}

	// Check in PATH
	pathEnv := os.Getenv("PATH")
	pathDirs := filepath.SplitList(pathEnv)

	for _, dir := range pathDirs {
		path := checkInDir(command, dir)
		if path != "" {
			return path, nil
		}
	}

	return "", fmt.Errorf("command not found: %s", command)
}
