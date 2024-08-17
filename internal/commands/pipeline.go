package commands

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

type Command struct {
	Name string
	Args []string
}

// ExecutePipeline executes a series of commands connected by pipes
func ExecutePipeline(commandsChain []Command) error {
	var lastOutput io.ReadCloser
	for i, cmd := range commandsChain {
		var err error
		if i == 0 {
			// First command: get the output directly
			lastOutput, err = executeCommandWithOutput(cmd, nil)
		} else {
			// Subsequent commands: get output from the previous command
			lastOutput, err = executeCommandWithOutput(cmd, lastOutput)
		}
		if err != nil {
			return err
		}
	}

	// Consume and display the final output
	if lastOutput != nil {
		_, err := io.Copy(os.Stdout, lastOutput)
		return err
	}
	return nil
}

// executeCommandWithOutput executes a command, possibly using input from a previous command
func executeCommandWithOutput(cmd Command, input io.Reader) (io.ReadCloser, error) {
	var command *exec.Cmd

	if IsBuiltinCommand(cmd.Name) {
		// Handle built-in commands
		if input != nil {
			return nil, fmt.Errorf("input redirection not supported for built-in commands: %s", cmd.Name)
		}
		err := ExecuteBuiltin(cmd.Name, cmd.Args)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}

	command = exec.Command(cmd.Name, cmd.Args...)

	if input != nil {
		command.Stdin = input
	}

	// Get a pipe for the command's standard output
	output, err := command.StdoutPipe()
	if err != nil {
		return nil, err
	}

	command.Stderr = os.Stderr

	// Start the command
	if err := command.Start(); err != nil {
		return nil, err
	}

	go func() {
		command.Wait()
	}()

	return output, nil
}
