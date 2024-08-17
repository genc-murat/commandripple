//go:build windows
// +build windows

package ls

import "os"

// getOwner returns the owner of the file for Unix-like systems.
func getOwner(file os.FileInfo) string {
	return "owner" // Placeholder for Windows, requires specific Windows API calls to get the owner.
}

// getGroup returns the group of the file for Unix-like systems.
func getGroup(file os.FileInfo) string {
	return "group" // Placeholder for Windows, requires specific Windows API calls to get the group.
}
