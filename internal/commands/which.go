package commands

import (
	"fmt"
)

// Which locates a command in the PATH
func Which(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: which [command]")
	}

	command := args[0]
	path, err := findExecutable(command)
	if err != nil {
		return fmt.Errorf("command not found: %s", command)
	}

	fmt.Println(path)
	return nil
}
