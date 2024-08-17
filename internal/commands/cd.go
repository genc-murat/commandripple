package commands

import (
	"fmt"
	"os"
	"path/filepath"
)

func ChangeDirectory(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("'cd' requires an argument")
	}

	targetDir := args[0]

	switch targetDir {
	case "~":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %v", err)
		}
		targetDir = homeDir
	case "..":
		currentDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %v", err)
		}
		targetDir = filepath.Dir(currentDir)
	case "-":
		// TODO: Implement changing to previous directory
		return fmt.Errorf("changing to previous directory is not implemented yet")
	}

	err := os.Chdir(targetDir)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("directory does not exist: %s", targetDir)
		} else if os.IsPermission(err) {
			return fmt.Errorf("permission denied: %s", targetDir)
		}
		return fmt.Errorf("failed to change directory: %v", err)
	}

	newDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}

	absPath, err := filepath.Abs(newDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %v", err)
	}

	fmt.Printf("Changed to directory: %s%s%s\n", ColorBlue, absPath, ColorReset)

	// Optional: Update shell prompt or environment variable with new directory
	os.Setenv("PWD", absPath)

	return nil
}
