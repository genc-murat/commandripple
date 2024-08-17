//go:build windows
// +build windows

package processes

import (
	"os/exec"
)

func getProcessListCommand() *exec.Cmd {
	return exec.Command("wmic", "process", "get", "ProcessId,ParentProcessId,UserModeTime,KernelModeTime,WorkingSetSize,CommandLine")
}
