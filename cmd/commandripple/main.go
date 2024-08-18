package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"commandripple/internal/commands"

	"github.com/chzyer/readline"
)

var (
	history     []string
	historyFile = "/tmp/commandripple_history"
)

func main() {
	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for sig := range sigChan {
			if sig == syscall.SIGINT {
				fmt.Println("\nCTRL-C detected. Use 'exit' to quit the shell.")
			}
		}
	}()

	// Initialize readline
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          getPrompt(),
		HistoryFile:     historyFile,
		AutoComplete:    completer{},
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil { // io.EOF, readline.ErrInterrupt
			break
		}

		line = strings.TrimSpace(line)

		if line == "exit" {
			break
		}

		if line == "" {
			continue
		}

		// Add command to history
		readline.AddHistory(line)

		if err := executePipeline(line); err != nil {
			fmt.Fprintf(os.Stderr, "CommandRipple: %v\n", err)
		}

		// Update prompt after each command execution
		rl.SetPrompt(getPrompt())
	}
}

func getPrompt() string {
	pwd, err := os.Getwd()
	if err != nil {
		pwd = "unknown"
	}
	// Shorten the path if it's too long
	if len(pwd) > 30 {
		pwd = "..." + pwd[len(pwd)-27:]
	}
	return fmt.Sprintf("\033[1;34m%s\033[0m CommandRipple> ", pwd)
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

// completer implements readline.AutoCompleter interface
type completer struct{}

func (c completer) Do(line []rune, pos int) (newLine [][]rune, length int) {
	completions := []string{
		"exit", "cd", "pwd", "echo", "clear", "mkdir", "mkdirp",
		"rmdir", "rm", "rmrf", "cp", "mv", "head", "tail", "grep",
		"find", "wc", "chmod", "chmodr", "env", "export", "history",
		"alias", "unalias", "date", "uptime", "kill", "ps", "whoami",
		"basename", "dirname", "sort", "uniq", "cut", "tee", "log", "calc",
		"truncate", "du", "df", "ln", "tr", "help", "ping", "ls", "cal", "touch",
		"stat", "dfi", "which", "killall", "source", "jobs", "fg", "bg", "compress",
		"decompress", "tree", "watch",
	}

	lineStr := string(line[:pos])
	var matches []string

	for _, comp := range completions {
		if strings.HasPrefix(comp, lineStr) {
			matches = append(matches, comp)
		}
	}

	if len(matches) == 0 {
		return
	}

	if len(matches) == 1 {
		newLine = [][]rune{[]rune(matches[0][pos:])}
		length = len(matches[0]) - pos
		return
	}

	for _, match := range matches {
		newLine = append(newLine, []rune(match))
	}
	length = pos
	return
}
