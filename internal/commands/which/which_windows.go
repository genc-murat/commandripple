//go:build windows
// +build windows

package which

import (
	"os"
	"path/filepath"
	"strings"
)

func isExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func checkCurrentDir(command string) string {
	return findWindowsExecutable(command, ".")
}

func checkInDir(command, dir string) string {
	return findWindowsExecutable(command, dir)
}

func findWindowsExecutable(command, dir string) string {
	exts := strings.Split(os.Getenv("PATHEXT"), ";")
	if len(exts) == 0 {
		exts = []string{".COM", ".EXE", ".BAT", ".CMD"}
	}

	// First, check if the command already has an extension
	if filepath.Ext(command) != "" {
		path := filepath.Join(dir, command)
		if isExecutable(path) {
			return path
		}
	}

	// If not, try appending each possible extension
	for _, ext := range exts {
		path := filepath.Join(dir, command+ext)
		if isExecutable(path) {
			return path
		}
	}

	return ""
}
