//go:build !windows
// +build !windows

package commands

import (
	"os"
	"path/filepath"
)

func isExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir() && info.Mode()&0111 != 0
}

func checkCurrentDir(command string) string {
	path := filepath.Join(".", command)
	if isExecutable(path) {
		return path
	}
	return ""
}

func checkInDir(command, dir string) string {
	path := filepath.Join(dir, command)
	if isExecutable(path) {
		return path
	}
	return ""
}
