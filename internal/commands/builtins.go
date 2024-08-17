package commands

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	history   []string
	aliases   = make(map[string]string)
	startTime = time.Now()
)

// IsBuiltinCommand checks if the command is a built-in command.
func IsBuiltinCommand(cmd string) bool {
	switch cmd {
	case "exit", "cd", "pwd", "echo", "clear", "mkdir", "mkdirp", "rmdir", "rm", "rmrf", "cp", "mv", "head", "tail", "grep", "find", "wc", "chmod", "chmodr", "env", "export", "history", "alias", "unalias", "date", "uptime", "kill", "ps", "whoami", "basename", "dirname", "sort", "uniq", "cut", "tee", "log", "calc", "truncate", "du", "df", "ln", "tr", "help", "ping", "ls", "cal", "touch", "stat", "dfi", "which", "killall":
		return true
	default:
		return false
	}
}

// ExecuteBuiltin executes the built-in commands.
func ExecuteBuiltin(cmd string, args []string) error {
	history = append(history, cmd+" "+strings.Join(args, " "))

	switch cmd {
	case "exit":
		fmt.Println("Exiting CommandRipple...")
		os.Exit(0)
	case "cd":
		return ChangeDirectory(args)
	case "pwd":
		return PrintWorkingDirectory()
	case "echo":
		return Echo(args)
	case "clear":
		return ClearScreen()
	case "mkdir":
		return MakeDirectory(args)
	case "rmdir":
		return RemoveDirectory(args)
	case "cat":
		return Cat(args)
	case "rm":
		return RemoveFile(args)
	case "cp":
		return CopyFile(args)
	case "mv":
		return MoveFile(args)
	case "head":
		return Head(args)
	case "tail":
		return Tail(args)
	case "grep":
		return Grep(args)
	case "find":
		return Find(args)
	case "wc":
		return WordCount(args)
	case "chmod":
		return Chmod(args)
	case "env":
		return PrintEnv()
	case "export":
		return ExportEnv(args)
	case "history":
		return ShowHistory()
	case "alias":
		return CreateAlias(args)
	case "unalias":
		return RemoveAlias(args)
	case "date":
		return ShowDate()
	case "uptime":
		return ShowUptime()
	case "kill":
		return KillProcess(args)
	case "ps":
		return ListProcesses()
	case "basename":
		return Basename(args)
	case "dirname":
		return Dirname(args)
	case "sort":
		return SortFile(args)
	case "uniq":
		return Uniq(args)
	case "cut":
		return Cut(args)
	case "tee":
		return Tee(args)
	case "log":
		return LogMessage(args)
	case "calc":
		return Calc(args)
	case "truncate":
		return Truncate(args)
	case "du":
		return Du(args)
	case "df":
		return Df(args)
	case "ln":
		return Ln(args)
	case "tr":
		return Tr(args)
	case "mkdirp":
		return MkdirP(args)
	case "rmrf":
		return RmRf(args)
	case "ping":
		return Ping(args)
	case "ls":
		return Ls(args)
	case "cal":
		return Cal(args)
	case "chmodr":
		return ChmodRecursive(args)
	case "touch":
		if len(args) > 0 && args[0] == "-t" {
			return TouchWithTimestamp(args)
		}
		return Touch(args)
	case "whoami":
		return Whoami(args)
	case "stat":
		return Stat(args)
	case "dfi":
		return DfInodes(args)
	case "which":
		return Which(args)
	case "killall":
		return KillAll(args)
	case "help":
		PrintHelp()
	default:
		return fmt.Errorf("unknown builtin command: %s", cmd)
	}
	return nil
}

// Locate a command in the PATH
func Which(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: which [command]")
	}
	cmd := exec.Command("which", args[0])
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Command History
func ShowHistory() error {
	for i, cmd := range history {
		fmt.Printf("%d %s\n", i+1, cmd)
	}
	return nil
}

// Alias Management
func CreateAlias(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("'alias' requires an argument in the format name=command")
	}
	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("'alias' argument must be in the format name=command")
		}
		name := parts[0]
		command := parts[1]
		aliases[name] = command
	}
	return nil
}

func RemoveAlias(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("'unalias' requires an argument")
	}
	for _, name := range args {
		delete(aliases, name)
	}
	return nil
}

// Date and Uptime
func ShowDate() error {
	fmt.Println(time.Now().Format(time.RFC1123))
	return nil
}

func ShowUptime() error {
	uptime := fmt.Sprintf("Uptime: %s", time.Since(startTime).String())
	PrintColor(Green, uptime)
	return nil
}

// Process Management
func KillProcess(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("'kill' requires a PID")
	}
	pid, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid PID: %s", args[0])
	}
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return process.Kill()
}

func ListProcesses() error {
	cmd := exec.Command("tasklist") // On Unix-like systems, use "ps"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Built-in command implementations
func ChangeDirectory(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("'cd' requires an argument")
	}
	return os.Chdir(args[0])
}

func PrintWorkingDirectory() error {
	if dir, err := os.Getwd(); err != nil {
		return err
	} else {
		fmt.Println(dir)
		return nil
	}
}

func Echo(args []string) error {
	PrintColor(Cyan, strings.Join(args, " "))
	return nil
}

func ClearScreen() error {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func MakeDirectory(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("'mkdir' requires an argument")
	}
	return os.Mkdir(args[0], os.ModePerm)
}

func RemoveDirectory(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("'rmdir' requires an argument")
	}
	return os.Remove(args[0])
}

func PrintEnv() error {
	for _, env := range os.Environ() {
		fmt.Println(env)
	}
	return nil
}

func ExportEnv(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("'export' requires an argument in the format NAME=VALUE")
	}
	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("'export' argument must be in the format NAME=VALUE")
		}
		name := parts[0]
		value := parts[1]
		if err := os.Setenv(name, value); err != nil {
			return fmt.Errorf("failed to set environment variable: %v", err)
		}
	}
	return nil
}

// Truncate or extend the size of a file
func Truncate(args []string) error {
	if len(args) < 3 || args[1] != "-s" {
		return fmt.Errorf("usage: truncate [file] -s [size]")
	}
	size, err := strconv.ParseInt(args[2], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid size: %s", args[2])
	}
	return os.Truncate(args[0], size)
}

// Estimate file space usage of a directory
func Du(args []string) error {
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}
	cmd := exec.Command("du", "-sh", dir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Report file system disk space usage
func Df(args []string) error {
	cmd := exec.Command("df", "-h")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Create a symbolic link between files
func Ln(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("usage: ln [target] [link]")
	}
	return os.Symlink(args[0], args[1])
}

// Translate or delete characters in a string
func Tr(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("usage: tr [set1] [set2]")
	}

	set1 := args[0]
	set2 := args[1]

	// Ensure both sets have the same length
	if len(set1) != len(set2) {
		return fmt.Errorf("set1 and set2 must have the same length")
	}

	// Create a slice to hold pairs of strings for the replacer
	replacements := make([]string, 0, len(set1)*2)
	for i := 0; i < len(set1); i++ {
		replacements = append(replacements, string(set1[i]), string(set2[i]))
	}

	// Create a new Replacer with the pairs
	replacer := strings.NewReplacer(replacements...)

	input := bufio.NewReader(os.Stdin)
	output := os.Stdout

	for {
		line, err := input.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Apply the replacement to the line
		transformedLine := replacer.Replace(line)
		output.WriteString(transformedLine)
	}

	return nil
}

// Recursively remove a directory and its contents
func RmRf(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: rm -rf [directory]")
	}
	for _, dir := range args {
		err := os.RemoveAll(dir)
		if err != nil {
			return err
		}
	}
	return nil
}

// Create a directory and its parent directories if they do not exist
func MkdirP(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: mkdir -p [directory]")
	}
	for _, dir := range args {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

// Send ICMP ECHO_REQUEST to network hosts
func Ping(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: ping [hostname]")
	}
	cmd := exec.Command("ping", args[0])
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// List directory contents with detailed file information
func Ls(args []string) error {
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}
	cmd := exec.Command("ls", "-l", dir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Display a calendar
func Cal(args []string) error {
	cmd := exec.Command("cal")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Create or update a file with a specific timestamp
func TouchWithTimestamp(args []string) error {
	if len(args) < 3 || args[0] != "-t" {
		return fmt.Errorf("usage: touch -t [timestamp] [file]")
	}
	timestamp := args[1]
	file := args[2]

	parsedTime, err := time.Parse("200601021504.05", timestamp)
	if err != nil {
		return fmt.Errorf("invalid timestamp format: %v", err)
	}

	err = os.Chtimes(file, parsedTime, parsedTime)
	if err != nil {
		return err
	}
	return nil
}

// Recursively change permissions of a directory
func ChmodRecursive(args []string) error {
	if len(args) != 3 || args[0] != "-R" {
		return fmt.Errorf("usage: chmod -R [permissions] [directory]")
	}
	mode, err := strconv.ParseUint(args[1], 8, 32)
	if err != nil {
		return fmt.Errorf("invalid permissions: %s", args[1])
	}
	dir := args[2]
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		return os.Chmod(path, os.FileMode(mode))
	})
}

// Display the current username
func Whoami(args []string) error {
	user := os.Getenv("USERNAME") // On Unix-like systems, use "USER"
	fmt.Println(user)
	return nil
}

// Report file system inode usage
func DfInodes(args []string) error {
	cmd := exec.Command("df", "-i")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Display file or file system status
func Stat(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: stat [file]")
	}
	cmd := exec.Command("stat", args[0])
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Kill all processes by name
func KillAll(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: killall [name]")
	}
	cmd := exec.Command("killall", args[0])
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Help function
func PrintHelp() {
	PrintColor(Cyan, "CommandRipple - A simple shell implemented in Go")
	PrintColor(White, "Built-in commands:")
	fmt.Println("  cd [dir]          Change the current directory")
	fmt.Println("  pwd               Print the current working directory")
	fmt.Println("  echo [text]       Echo the input text back to the user")
	fmt.Println("  clear             Clear the terminal screen")
	fmt.Println("  mkdir [dir]       Create a new directory")
	fmt.Println("  mkdirp [dir]      Create directories and parent directories if needed")
	fmt.Println("  rmdir [dir]       Remove an empty directory")
	fmt.Println("  rm [file]         Remove a file")
	fmt.Println("  rmrf [dir]        Recursively remove a directory and its contents")
	fmt.Println("  cp [src] [dest]   Copy a file")
	fmt.Println("  mv [src] [dest]   Move or rename a file or directory")
	fmt.Println("  touch [file]      Create an empty file or update timestamp")
	fmt.Println("  touch -t [timestamp] [file] Create or update a file with a specific timestamp")
	fmt.Println("  chmod [permissions] [file] Change file permissions")
	fmt.Println("  chmodr [permissions] [dir] Recursively change permissions of a directory")
	fmt.Println("  cat [file]        Display the content of a file")
	fmt.Println("  head [file]       Display the first few lines of a file")
	fmt.Println("  tail [file]       Display the last few lines of a file")
	fmt.Println("  grep [pattern] [file] Search for a pattern in a file")
	fmt.Println("  find [dir] [name] Search for a file or directory by name")
	fmt.Println("  wc [file]         Count lines, words, and characters in a file")
	fmt.Println("  env               Print environment variables")
	fmt.Println("  export NAME=VALUE Set or modify environment variables")
	fmt.Println("  history           Display command history")
	fmt.Println("  alias name=command Create an alias for a command")
	fmt.Println("  unalias name      Remove an alias")
	fmt.Println("  date              Display the current date and time")
	fmt.Println("  uptime            Display how long the shell has been running")
	fmt.Println("  kill [PID]        Terminate a process by PID")
	fmt.Println("  killall [name]    Kill all processes by name")
	fmt.Println("  ps                List currently running processes")
	fmt.Println("  whoami            Display the current user's username")
	fmt.Println("  basename [path]   Strip directory and suffix from filenames")
	fmt.Println("  dirname [path]    Extract the directory path from a full path")
	fmt.Println("  sort [file]       Sort lines of a text file")
	fmt.Println("  uniq [file]       Remove duplicate lines from a file")
	fmt.Println("  cut [file] -d [delimiter] -f [field] Extract selected portions of each line")
	fmt.Println("  tee [file]        Read from standard input and write to standard output and files")
	fmt.Println("  log [message]     Append a log message to a log file")
	fmt.Println("  calc [expression] Evaluate a simple arithmetic expression")
	fmt.Println("  truncate [file] -s [size] Truncate or extend the size of a file")
	fmt.Println("  du [dir]          Estimate file space usage of a directory")
	fmt.Println("  df                Report file system disk space usage")
	fmt.Println("  dfi               Report file system inode usage")
	fmt.Println("  ln [target] [link] Create a symbolic link between files")
	fmt.Println("  tr [set1] [set2]  Translate or delete characters in a string")
	fmt.Println("  ping [hostname]   Send ICMP ECHO_REQUEST to network hosts")
	fmt.Println("  which [command]   Locate a command in the PATH")
	fmt.Println("  ls [dir]          List directory contents with detailed file information")
	fmt.Println("  stat [file]       Display file or file system status")
	fmt.Println("  cal               Display a calendar")
	fmt.Println("  help              Show this help message")
	fmt.Println("\nPipes:")
	fmt.Println("  Use the '|' character to pipe the output of one command to the input of another.")
	fmt.Println("  Example: cat file.txt | grep 'search' | sort")
}
