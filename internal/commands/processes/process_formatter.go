package processes

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
)

func FormatProcessList(output string) (string, error) {
	lines := strings.Split(output, "\n")
	if len(lines) < 2 {
		return "", fmt.Errorf("unexpected output format")
	}

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetHeader([]string{"PID", "PPID", "CPU%", "MEM%", "COMMAND"})
	table.SetBorder(false)
	table.SetColumnSeparator("")
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.FgHiCyanColor},
		tablewriter.Colors{tablewriter.FgHiCyanColor},
		tablewriter.Colors{tablewriter.FgHiCyanColor},
		tablewriter.Colors{tablewriter.FgHiCyanColor},
		tablewriter.Colors{tablewriter.FgHiCyanColor},
	)

	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}

		formattedFields := formatProcessLine(fields)
		if len(formattedFields) == 5 {
			table.Rich(formattedFields, []tablewriter.Colors{
				{tablewriter.FgYellowColor},
				{tablewriter.FgGreenColor},
				{tablewriter.FgRedColor},
				{tablewriter.FgMagentaColor},
				{tablewriter.FgWhiteColor},
			})
		}
	}

	table.Render()
	return tableString.String(), nil
}

func formatProcessLine(fields []string) []string {
	var pid, ppid, cpu, mem, command string

	switch runtime.GOOS {
	case "windows":
		if len(fields) >= 6 {
			pid = fields[1]
			ppid = fields[2]
			cpu = calculateCPUUsage(fields[3], fields[4])
			mem = formatMemory(fields[5])
			command = strings.Join(fields[6:], " ")
		}
	default:
		if len(fields) >= 5 {
			pid = fields[0]
			ppid = fields[1]
			cpu = fields[2]
			mem = fields[3]
			command = strings.Join(fields[4:], " ")
		}
	}

	if len(command) > 30 {
		command = command[:27] + "..."
	}

	return []string{pid, ppid, cpu, mem, command}
}

func calculateCPUUsage(userModeTime, kernelModeTime string) string {
	totalTime := parseTime(userModeTime) + parseTime(kernelModeTime)
	return fmt.Sprintf("%.2f%%", float64(totalTime)/100)
}

func parseTime(timeStr string) int {
	timeVal, err := strconv.Atoi(timeStr)
	if err != nil {
		return 0
	}
	return timeVal
}

func formatMemory(memStr string) string {
	memVal, err := strconv.Atoi(memStr)
	if err != nil {
		return memStr
	}
	return fmt.Sprintf("%.2f MB", float64(memVal)/1024/1024)
}
