package processes

import (
	"fmt"
	"runtime"
	"strings"
)

func FormatProcessList(output string) (string, error) {
	lines := strings.Split(output, "\n")
	if len(lines) < 2 {
		return "", fmt.Errorf("unexpected output format")
	}

	header := "PID\tPPID\tCPU%\tMEM%\tCOMMAND"
	formattedLines := []string{header}

	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}

		formattedLine := formatProcessLine(fields)
		if formattedLine != "" {
			formattedLines = append(formattedLines, formattedLine)
		}
	}

	return strings.Join(formattedLines, "\n"), nil
}

func formatProcessLine(fields []string) string {
	var pid, ppid, cpu, mem, command string

	switch runtime.GOOS {
	case "windows":
		pid = fields[1]
		mem = fields[4]
		command = strings.Join(fields[0:1], " ")
		ppid = "N/A"
		cpu = "N/A"
	case "darwin":
		pid = fields[0]
		ppid = fields[1]
		cpu = fields[2]
		mem = fields[3]
		command = strings.Join(fields[4:], " ")
	default:
		pid = fields[1]
		ppid = "N/A"
		cpu = fields[2]
		mem = fields[3]
		command = strings.Join(fields[10:], " ")
	}

	return fmt.Sprintf("%s\t%s\t%s\t%s\t%s", pid, ppid, cpu, mem, command)
}
