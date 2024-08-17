//go:build windows
// +build windows

package processes

import (
	"os/exec"
)

func getProcessListCommand() *exec.Cmd {
	return exec.Command("tasklist")
}
