package commands

import (
	"os"
	"os/exec"
)

// ExecuteExternal executes external commands using cmd.exe /c on Windows.
func ExecuteExternal(cmdName string, args []string) error {
	cmd := exec.Command("cmd", "/c", cmdName)
	cmd.Args = append(cmd.Args, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
