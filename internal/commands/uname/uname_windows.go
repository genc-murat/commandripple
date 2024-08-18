//go:build windows
// +build windows

package uname

import (
	"os/exec"
	"strings"
)

func getSystemInfo(info SystemInfo) (SystemInfo, error) {
	out, err := exec.Command("ver").Output()
	if err == nil {
		info.Version = strings.TrimSpace(string(out))
	}

	out, err = exec.Command("wmic", "os", "get", "Version", "/Value").Output()
	if err == nil {
		info.Kernel = strings.TrimSpace(strings.TrimPrefix(string(out), "Version="))
	}

	return info, nil
}
