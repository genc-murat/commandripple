package processes

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"
)

// KillAll terminates all processes with the given name
func KillAll(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: killall [name]")
	}
	processName := args[0]

	if runtime.GOOS == "windows" {
		return killAllWindows(processName)
	} else {
		return killAllUnix(processName)
	}
}

func killAllWindows(processName string) error {
	cmd := exec.Command("taskkill", "/F", "/IM", processName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func killAllUnix(processName string) error {
	// First, try using the 'killall' command if it exists
	killallCmd := exec.Command("killall", processName)
	if err := killallCmd.Run(); err == nil {
		return nil
	}

	// If 'killall' doesn't exist or fails, fall back to manual process killing
	processes, err := getProcessesByName(processName)
	if err != nil {
		return fmt.Errorf("error finding processes: %v", err)
	}

	for _, pid := range processes {
		process, err := os.FindProcess(pid)
		if err != nil {
			fmt.Printf("Warning: Could not find process %d: %v\n", pid, err)
			continue
		}

		err = process.Signal(syscall.SIGTERM)
		if err != nil {
			fmt.Printf("Warning: Could not send SIGTERM to process %d: %v\n", pid, err)
			err = process.Signal(syscall.SIGKILL)
			if err != nil {
				fmt.Printf("Error: Could not send SIGKILL to process %d: %v\n", pid, err)
			}
		}
	}

	return nil
}

func getProcessesByName(name string) ([]int, error) {
	cmd := exec.Command("pgrep", name)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error running pgrep: %v", err)
	}

	var pids []int
	for _, line := range strings.Split(string(output), "\n") {
		if line != "" {
			pid, err := strconv.Atoi(line)
			if err != nil {
				fmt.Printf("Warning: Could not parse PID '%s': %v\n", line, err)
				continue
			}
			pids = append(pids, pid)
		}
	}

	return pids, nil
}
