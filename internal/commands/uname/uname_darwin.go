//go:build darwin
// +build darwin

package uname

import (
	"os/exec"
	"strings"
)

func getSystemInfo(info SystemInfo) (SystemInfo, error) {
	out, err := exec.Command("uname", "-v").Output()
	if err == nil {
		info.Kernel = strings.TrimSpace(string(out))
	}

	out, err = exec.Command("sw_vers", "-productVersion").Output()
	if err == nil {
		info.Version = strings.TrimSpace(string(out))
	}

	return info, nil
}
