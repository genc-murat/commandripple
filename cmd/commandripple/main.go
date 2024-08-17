package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
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

	// To track the previous command's output pipe
	var previousOutput io.ReadCloser

	for i, cmd := range commandsList {
		trimmedCmd := strings.TrimSpace(cmd)
		if trimmedCmd == "" {
			continue
		}

		cmdName := strings.Fields(trimmedCmd)[0]
		cmdArgs := strings.Fields(trimmedCmd)[1:]

		// Create the command object
		command := exec.Command(cmdName, cmdArgs...)

		// Set stdin to the previous command's output
		if previousOutput != nil {
			command.Stdin = previousOutput
		}

		// If it's not the last command, create a pipe
		if i < numCommands-1 {
			var err error
			previousOutput, err = command.StdoutPipe()
			if err != nil {
				return fmt.Errorf("failed to create stdout pipe: %v", err)
			}
		} else {
			// For the last command, connect stdout to the terminal
			command.Stdout = os.Stdout
		}

		command.Stderr = os.Stderr

		// Start the command
		if err := command.Start(); err != nil {
			return fmt.Errorf("failed to start command: %v", err)
		}

		// Wait for the command to finish
		if err := command.Wait(); err != nil {
			return fmt.Errorf("command failed: %v", err)
		}

		// Close the previous output if necessary
		if previousOutput != nil {
			previousOutput.Close()
		}
	}

	return nil
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
