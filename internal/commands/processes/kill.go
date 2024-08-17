package processes

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"syscall"
)

// KillProcess terminates a process by its PID
func KillProcess(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("'kill' requires a PID")
	}
	pid, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid PID: %s", args[0])
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("process not found: %v", err)
	}

	var killErr error
	if runtime.GOOS == "windows" {
		killErr = process.Kill()
	} else {
		// On Unix-like systems, first try a graceful termination
		killErr = process.Signal(syscall.SIGTERM)
		if killErr == nil {
			// Wait a bit to see if the process terminates
			_, err := process.Wait()
			if err != nil {
				// If waiting fails, force kill the process
				killErr = process.Signal(syscall.SIGKILL)
			}
		}
	}

	if killErr != nil {
		return fmt.Errorf("failed to kill process: %v", killErr)
	}

	fmt.Printf("Process with PID %d has been terminated\n", pid)
	return nil
}
