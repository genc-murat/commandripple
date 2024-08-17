package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"commandripple/internal/commands"
)

func main() {
	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for sig := range sigChan {
			if sig == syscall.SIGINT {
				fmt.Println("\nCTRL-C detected. Use 'exit' to quit the shell.")
				fmt.Print("CommandRipple> ")
			}
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("CommandRipple> ")

		if !scanner.Scan() { // Reads input from the user
			break
		}

		commandLine := strings.TrimSpace(scanner.Text())

		if commandLine == "" {
			continue
		}

		if err := executePipeline(commandLine); err != nil {
			fmt.Fprintf(os.Stderr, "CommandRipple: %v\n", err)
		}
	}
}

func executePipeline(commandLine string) error {
	commandsList := strings.Split(commandLine, "|")
	numCommands := len(commandsList)

	if numCommands == 1 {
		// If there's no pipe, execute the command normally
		return executeCommand(commandLine)
	}

	var commandsChain []commands.Command

	// Create a list of commands to execute
	for _, cmd := range commandsList {
		trimmedCmd := strings.TrimSpace(cmd)
		if trimmedCmd == "" {
			continue
		}

		cmdName := strings.Fields(trimmedCmd)[0]
		cmdArgs := strings.Fields(trimmedCmd)[1:]

		commandsChain = append(commandsChain, commands.Command{
			Name: cmdName,
			Args: cmdArgs,
		})
	}

	// Execute the command pipeline
	return commands.ExecutePipeline(commandsChain)
}

func executeCommand(commandLine string) error {
	parts := strings.Fields(commandLine)
	cmdName := parts[0]
	cmdArgs := parts[1:]

	if commands.IsBuiltinCommand(cmdName) {
		return commands.ExecuteBuiltin(cmdName, cmdArgs)
	} else {
		return commands.ExecuteExternal(cmdName, cmdArgs)
	}
}
