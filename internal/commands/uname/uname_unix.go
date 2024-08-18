//go:build linux
// +build linux

package uname

import (
	"os/exec"
	"strings"
)

func getSystemInfo(info SystemInfo) (SystemInfo, error) {
	out, err := exec.Command("uname", "-r").Output()
	if err == nil {
		info.Kernel = strings.TrimSpace(string(out))
	}

	out, err = exec.Command("cat", "/etc/os-release").Output()
	if err == nil {
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "VERSION=") {
				info.Version = strings.Trim(strings.TrimPrefix(line, "VERSION="), "\"")
				break
			}
		}
	}

	return info, nil
}
