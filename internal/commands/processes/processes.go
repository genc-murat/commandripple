package processes

import (
	"fmt"
	"os/exec"
	"runtime"
)

func ListProcesses() error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("tasklist")
	case "darwin":
		cmd = exec.Command("ps", "-eo", "pid,ppid,%cpu,%mem,command")
	default:
		cmd = exec.Command("ps", "aux")
	}

	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("error listing processes: %v", err)
	}

	formattedOutput, err := FormatProcessList(string(output))
	if err != nil {
		return fmt.Errorf("error formatting process list: %v", err)
	}

	fmt.Println(formattedOutput)
	return nil
}
