//go:build !windows
// +build !windows

package stat

import (
	"fmt"
	"os"
	"syscall"
)

func printDetailedStats(filePath string, fileInfo os.FileInfo) {
	stat, ok := fileInfo.Sys().(*syscall.Stat_t)
	if !ok {
		return
	}
	fmt.Printf("  Inode: %d\n", stat.Ino)
	fmt.Printf("  Links: %d\n", stat.Nlink)
	fmt.Printf("  UID: %d\n", stat.Uid)
	fmt.Printf("  GID: %d\n", stat.Gid)
}
