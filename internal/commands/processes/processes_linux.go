//go:build linux
// +build linux

package processes

import (
	"os/exec"
)

func getProcessListCommand() *exec.Cmd {
	return exec.Command("ps", "aux")
}
