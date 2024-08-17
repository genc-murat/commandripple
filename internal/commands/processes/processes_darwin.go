//go:build darwin
// +build darwin

package processes

import (
	"os/exec"
)

func getProcessListCommand() *exec.Cmd {
	return exec.Command("ps", "-eo", "pid,ppid,%cpu,%mem,command")
}
