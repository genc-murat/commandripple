//go:build !windows
// +build !windows

package clear

import (
	"os"
	"os/exec"
)

func clear() error {
	cmd := exec.Command("clear")

	cmd.Stdout = os.Stdout
	return cmd.Run()
}
