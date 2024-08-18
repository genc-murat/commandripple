package uname

import (
	"fmt"
	"os"
	"runtime"
)

type SystemInfo struct {
	OS           string
	Kernel       string
	Architecture string
	Hostname     string
	Version      string
}

func Uname(args []string) error {
	info := SystemInfo{
		OS:           runtime.GOOS,
		Architecture: runtime.GOARCH,
	}

	var err error
	info.Hostname, err = os.Hostname()
	if err != nil {
		return fmt.Errorf("error getting hostname: %v", err)
	}

	sysInfo, err := getSystemInfo(info)
	if err != nil {
		return fmt.Errorf("error getting system info: %v", err)
	}

	if len(args) > 0 && args[0] == "-a" {
		printDetailedSystemInfo(sysInfo)
	} else {
		printBasicSystemInfo(sysInfo)
	}
	return nil
}

func printBasicSystemInfo(info SystemInfo) {
	fmt.Printf("%s %s %s\n", info.OS, info.Kernel, info.Architecture)
}

func printDetailedSystemInfo(info SystemInfo) {
	fmt.Printf("Operating System: %s\n", info.OS)
	fmt.Printf("Kernel: %s\n", info.Kernel)
	fmt.Printf("Architecture: %s\n", info.Architecture)
	fmt.Printf("Hostname: %s\n", info.Hostname)
	fmt.Printf("Version: %s\n", info.Version)
}
