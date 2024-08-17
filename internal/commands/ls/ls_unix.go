//go:build !windows
// +build !windows

package ls

import (
	"os"
	"strconv"
	"syscall"
)

func getOwner(file os.FileInfo) string {
	stat := file.Sys().(*syscall.Stat_t)
	uid := stat.Uid
	return strconv.Itoa(int(uid)) // Returns the owner's user ID as a string
}

func getGroup(file os.FileInfo) string {
	stat := file.Sys().(*syscall.Stat_t)
	gid := stat.Gid
	return strconv.Itoa(int(gid)) // Returns the group's ID as a string
}
