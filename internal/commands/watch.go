package commands

import (
	"fmt"
	"strings"
	"time"
)

func Watch(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: watch [interval] [command]")
	}

	interval, err := time.ParseDuration(args[0])
	if err != nil {
		return fmt.Errorf("invalid interval: %v", err)
	}

	command := strings.Join(args[1:], " ")

	for {
		fmt.Printf("\033[2J") // Clear screen
		fmt.Printf("\033[H")  // Move cursor to top-left
		fmt.Printf("Every %v: %s\n\n", interval, command)

		err := executeCommand(command)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}

		time.Sleep(interval)
	}
}
