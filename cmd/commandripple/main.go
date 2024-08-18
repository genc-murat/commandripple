package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"commandripple/internal/commands"

	"github.com/chzyer/readline"
)

var (
	history     []string
	historyFile = filepath.Join(os.TempDir(), "commandripple_history")
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
		AutoComplete:    newCompleter(),
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

func newCompleter() *completer {
	return &completer{}
}

func (c *completer) Do(line []rune, pos int) (newLine [][]rune, length int) {
	lineStr := string(line[:pos])
	parts := strings.Fields(lineStr)

	if len(parts) == 0 {
		return c.completeCommands(lineStr)
	}

	cmdName := parts[0]
	if len(parts) == 1 {
		return c.completeCommands(lineStr)
	}

	// Check if we're completing a path
	if strings.Contains(parts[len(parts)-1], string(os.PathSeparator)) || strings.Contains(parts[len(parts)-1], ".") {
		return c.completeFilesDirs(lineStr)
	}

	// Command-specific argument completion
	switch cmdName {
	case "cd", "ls", "rm", "mv", "cp":
		return c.completeFilesDirs(lineStr)
	case "chmod":
		if len(parts) == 2 {
			return c.completeChmodArgs(lineStr)
		}
		return c.completeFilesDirs(lineStr)
	// Add more command-specific completions here
	default:
		return c.completeFilesDirs(lineStr)
	}
}

func (c *completer) completeCommands(lineStr string) (newLine [][]rune, length int) {
	commands := []string{
		"exit", "cd", "pwd", "echo", "clear", "mkdir", "mkdirp",
		"rmdir", "rm", "rmrf", "cp", "mv", "head", "tail", "grep",
		"find", "wc", "chmod", "chmodr", "env", "export", "history",
		"alias", "unalias", "date", "uptime", "kill", "ps", "whoami",
		"basename", "dirname", "sort", "uniq", "cut", "tee", "log", "calc",
		"truncate", "du", "df", "ln", "tr", "help", "ping", "ls", "cal", "touch",
		"stat", "dfi", "which", "killall", "source", "jobs", "fg", "bg", "compress",
		"decompress", "tree", "watch", "free", "uname", "remote_execute", "file_transfer",
	}

	return c.filterCompletions(lineStr, commands)
}

func (c *completer) completeFilesDirs(lineStr string) (newLine [][]rune, length int) {
	parts := strings.Fields(lineStr)
	lastPart := parts[len(parts)-1]
	dir := filepath.Dir(lastPart)
	prefix := filepath.Base(lastPart)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, 0
	}

	var matches []string

	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, prefix) {
			if entry.IsDir() {
				name += string(os.PathSeparator)
			}
			matches = append(matches, filepath.Join(dir, name))
		}
	}

	return c.filterCompletions(lastPart, matches)
}

func (c *completer) completeChmodArgs(lineStr string) (newLine [][]rune, length int) {
	chmodArgs := []string{"644", "755", "777", "600", "400"}
	return c.filterCompletions(lineStr, chmodArgs)
}

func (c *completer) filterCompletions(lineStr string, completions []string) (newLine [][]rune, length int) {
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
		newLine = [][]rune{[]rune(matches[0][len(lineStr):])}
		length = len(matches[0])
		return
	}

	for _, match := range matches {
		newLine = append(newLine, []rune(match))
	}
	length = len(lineStr)
	return
}
