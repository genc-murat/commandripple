package commands

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"commandripple/internal/commands/processes"
	"commandripple/internal/commands/stat"
	"commandripple/internal/commands/which"

	"github.com/olekukonko/tablewriter"
)

var (
	history     []string
	aliases     = make(map[string]string)
	startTime   = time.Now()
	bgJobs      = make(map[int]*JobInfo) // Store background jobs
	bgJobsMutex sync.Mutex               // To handle concurrent access to bgJobs
	jobCounter  int                      // Unique identifier for jobs
)

type JobInfo struct {
	Cmd       *exec.Cmd
	StartTime time.Time
	Status    string
}

// IsBuiltinCommand checks if the command is a built-in command.
func IsBuiltinCommand(cmd string) bool {
	switch cmd {
	case "exit", "cd", "pwd", "echo", "clear", "mkdir", "mkdirp", "rmdir", "rm", "rmrf", "cp", "mv", "head", "tail", "grep", "find", "wc", "chmod", "chmodr", "env", "export", "history", "alias", "unalias", "date", "uptime", "kill", "ps", "whoami", "basename", "dirname", "sort", "uniq", "cut", "tee", "log", "calc", "truncate", "du", "df", "ln", "tr", "help", "ping", "ls", "lsc", "cal", "touch", "stat", "dfi", "which", "killall", "source", "jobs", "fg", "bg", "tree", "watch", "compress", "decompress":
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
		return processes.KillProcess(args)
	case "ps":
		return processes.ListProcesses()
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
	case "lsc":
		return LsColor(args)
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
		return stat.Stat(args)
	case "dfi":
		return DfInodes(args)
	case "which":
		return which.Which(args)
	case "killall":
		return processes.KillAll(args)
	case "source":
		return Source(args)
	case "jobs":
		return ListJobs()
	case "fg":
		return BringToForeground(args)
	case "bg":
		return SendToBackground(args)
	case "tree":
		return Tree(args)
	case "watch":
		return Watch(args)
	case "compress":
		return Compress(args)
	case "decompress":
		return Decompress(args)
	case "help":
		PrintHelp()
	default:
		return fmt.Errorf("unknown builtin command: %s", cmd)
	}
	return nil
}

// `jobs` command implementation
func ListJobs() error {
	bgJobsMutex.Lock()
	defer bgJobsMutex.Unlock()

	if len(bgJobs) == 0 {
		fmt.Println("No background jobs running.")
		return nil
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Command", "Status", "Started", "Runtime"})

	for id, jobInfo := range bgJobs {
		command := strings.Join(jobInfo.Cmd.Args, " ")
		status := jobInfo.Status
		started := jobInfo.StartTime.Format("2006-01-02 15:04:05")
		runtime := time.Since(jobInfo.StartTime).Round(time.Second).String()

		table.Append([]string{
			fmt.Sprintf("%d", id),
			command,
			status,
			started,
			runtime,
		})
	}

	table.Render()
	return nil
}

// `source` command implementation
func Source(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("'source' requires a filename")
	}
	file, err := os.Open(args[0])
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		commandLine := strings.TrimSpace(scanner.Text())
		if commandLine == "" {
			continue
		}
		if err := executeCommand(commandLine); err != nil {
			fmt.Fprintf(os.Stderr, "CommandRipple: %v\n", err)
		}
	}

	return scanner.Err()
}

// `fg` command implementation
func BringToForeground(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("'fg' requires a job ID")
	}

	jobID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid job ID: %s", args[0])
	}

	bgJobsMutex.Lock()
	jobInfo, exists := bgJobs[jobID]
	if !exists {
		bgJobsMutex.Unlock()
		return fmt.Errorf("no such job: %d", jobID)
	}

	if jobInfo.Status == "Completed" {
		bgJobsMutex.Unlock()
		return fmt.Errorf("job %d has already completed", jobID)
	}

	delete(bgJobs, jobID)
	bgJobsMutex.Unlock()

	// Bring the job to the foreground
	jobInfo.Cmd.Stdout = os.Stdout
	jobInfo.Cmd.Stderr = os.Stderr
	jobInfo.Cmd.Stdin = os.Stdin

	fmt.Printf("Bringing job %d to foreground: %s\n", jobID, strings.Join(jobInfo.Cmd.Args, " "))

	err = jobInfo.Cmd.Wait()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("job exited with status %d", exitErr.ExitCode())
		}
		return fmt.Errorf("error waiting for job: %v", err)
	}

	return nil
}

// Helper function to create a new background job
func StartBackgroundJob(cmd *exec.Cmd) (int, error) {
	bgJobsMutex.Lock()
	defer bgJobsMutex.Unlock()

	jobCounter++
	id := jobCounter

	bgJobs[id] = &JobInfo{
		Cmd:       cmd,
		StartTime: time.Now(),
		Status:    "Running",
	}

	go func() {
		cmd.Wait()
		bgJobsMutex.Lock()
		if job, exists := bgJobs[id]; exists {
			job.Status = "Completed"
		}
		bgJobsMutex.Unlock()
	}()

	return id, nil
}

// `bg` command implementation
func SendToBackground(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("'bg' requires a command")
	}

	cmdName := args[0]
	cmdArgs := args[1:]

	command := exec.Command(cmdName, cmdArgs...)

	// Start the command
	err := command.Start()
	if err != nil {
		return fmt.Errorf("failed to start command: %v", err)
	}

	// Add the job to bgJobs
	jobID, err := StartBackgroundJob(command)
	if err != nil {
		return fmt.Errorf("failed to add job to background: %v", err)
	}

	fmt.Printf("[%d] %s\n", jobID, strings.Join(command.Args, " "))
	return nil
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

// Ls lists directory contents with detailed file information
func Ls(args []string) error {
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		mode := file.Mode()
		size := file.Size()
		modTime := file.ModTime().Format(time.RFC822)
		name := file.Name()

		// Create a string that mimics the output of 'ls -l' on Unix
		var fileInfo string
		if runtime.GOOS == "windows" {
			// On Windows, we'll use a simplified format
			fileType := "f"
			if file.IsDir() {
				fileType = "d"
			}
			fileInfo = fmt.Sprintf("%s %10d %s %s", fileType, size, modTime, name)
		} else {
			// On Unix-like systems, we'll try to mimic 'ls -l' more closely
			perms := mode.String()
			owner := getOwner(file)
			group := getGroup(file)
			fileInfo = fmt.Sprintf("%s %s %s %8d %s %s", perms, owner, group, size, modTime, name)
		}

		fmt.Println(fileInfo)
	}

	return nil
}

// LsColor lists directory contents with colors (for file types) and detailed information
func LsColor(args []string) error {
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		mode := file.Mode()
		size := file.Size()
		modTime := file.ModTime().Format(time.RFC822)
		name := file.Name()

		// Determine color based on file type
		color := "\033[0m" // Default color (reset)
		if file.IsDir() {
			color = "\033[1;34m" // Blue for directories
		} else if mode&0111 != 0 {
			color = "\033[1;32m" // Green for executable files
		} else if strings.HasPrefix(strings.ToLower(filepath.Ext(name)), ".") {
			color = "\033[1;37m" // White for hidden files
		}

		// Create a string that mimics the output of 'ls -l' on Unix
		var fileInfo string
		if runtime.GOOS == "windows" {
			// On Windows, we'll use a simplified format
			fileType := "f"
			if file.IsDir() {
				fileType = "d"
			}
			fileInfo = fmt.Sprintf("%s %10d %s %s%s%s", fileType, size, modTime, color, name, "\033[0m")
		} else {
			// On Unix-like systems, we'll try to mimic 'ls -l' more closely
			perms := mode.String()
			owner := getOwner(file)
			group := getGroup(file)
			fileInfo = fmt.Sprintf("%s %s %s %8d %s %s%s%s", perms, owner, group, size, modTime, color, name, "\033[0m")
		}

		fmt.Println(fileInfo)
	}

	return nil
}

func getOwner(file os.FileInfo) string {
	if runtime.GOOS == "windows" {
		return "owner"
	}
	return "owner" // Replace with actual owner retrieval for Unix systems
}

func getGroup(file os.FileInfo) string {
	if runtime.GOOS == "windows" {
		return "group"
	}
	return "group" // Replace with actual group retrieval for Unix systems
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

// Command execution helper
func executeCommand(commandLine string) error {
	parts := strings.Fields(commandLine)
	cmdName := parts[0]
	cmdArgs := parts[1:]

	if IsBuiltinCommand(cmdName) {
		return ExecuteBuiltin(cmdName, cmdArgs)
	} else {
		return ExecuteExternal(cmdName, cmdArgs)
	}
}

// Help function
func PrintHelp() {
	PrintColor(Cyan, "CommandRipple - A simple shell implemented in Go")
	PrintColor(White, "Built-in commands:")
	PrintColor(Green, "  cd [dir]          ")
	fmt.Println("Change the current directory")
	PrintColor(Green, "  pwd               ")
	fmt.Println("Print the current working directory")
	PrintColor(Green, "  echo [text]       ")
	fmt.Println("Echo the input text back to the user")
	PrintColor(Green, "  clear             ")
	fmt.Println("Clear the terminal screen")
	PrintColor(Green, "  mkdir [dir]       ")
	fmt.Println("Create a new directory")
	PrintColor(Green, "  mkdirp [dir]      ")
	fmt.Println("Create directories and parent directories if needed")
	PrintColor(Green, "  rmdir [dir]       ")
	fmt.Println("Remove an empty directory")
	PrintColor(Green, "  rm [file]         ")
	fmt.Println("Remove a file")
	PrintColor(Green, "  rmrf [dir]        ")
	fmt.Println("Recursively remove a directory and its contents")
	PrintColor(Green, "  cp [src] [dest]   ")
	fmt.Println("Copy a file")
	PrintColor(Green, "  mv [src] [dest]   ")
	fmt.Println("Move or rename a file or directory")
	PrintColor(Green, "  touch [file]      ")
	fmt.Println("Create an empty file or update timestamp")
	PrintColor(Green, "  touch -t [timestamp] [file] ")
	fmt.Println("Create or update a file with a specific timestamp")
	PrintColor(Green, "  chmod [permissions] [file] ")
	fmt.Println("Change file permissions")
	PrintColor(Green, "  chmodr [permissions] [dir] ")
	fmt.Println("Recursively change permissions of a directory")
	PrintColor(Green, "  cat [file]        ")
	fmt.Println("Display the content of a file")
	PrintColor(Green, "  head [file]       ")
	fmt.Println("Display the first few lines of a file")
	PrintColor(Green, "  tail [file]       ")
	fmt.Println("Display the last few lines of a file")
	PrintColor(Green, "  grep [pattern] [file] ")
	fmt.Println("Search for a pattern in a file")
	PrintColor(Green, "  find [dir] [name] ")
	fmt.Println("Search for a file or directory by name")
	PrintColor(Green, "  wc [file]         ")
	fmt.Println("Count lines, words, and characters in a file")
	PrintColor(Green, "  env               ")
	fmt.Println("Print environment variables")
	PrintColor(Green, "  export NAME=VALUE ")
	fmt.Println("Set or modify environment variables")
	PrintColor(Green, "  history           ")
	fmt.Println("Display command history")
	PrintColor(Green, "  alias name=command ")
	fmt.Println("Create an alias for a command")
	PrintColor(Green, "  unalias name      ")
	fmt.Println("Remove an alias")
	PrintColor(Green, "  date              ")
	fmt.Println("Display the current date and time")
	PrintColor(Green, "  uptime            ")
	fmt.Println("Display how long the shell has been running")
	PrintColor(Green, "  kill [PID]        ")
	fmt.Println("Terminate a process by PID")
	PrintColor(Green, "  killall [name]    ")
	fmt.Println("Kill all processes by name")
	PrintColor(Green, "  ps                ")
	fmt.Println("List currently running processes")
	PrintColor(Green, "  whoami            ")
	fmt.Println("Display the current user's username")
	PrintColor(Green, "  basename [path]   ")
	fmt.Println("Strip directory and suffix from filenames")
	PrintColor(Green, "  dirname [path]    ")
	fmt.Println("Extract the directory path from a full path")
	PrintColor(Green, "  sort [file]       ")
	fmt.Println("Sort lines of a text file")
	PrintColor(Green, "  uniq [file]       ")
	fmt.Println("Remove duplicate lines from a file")
	PrintColor(Green, "  cut [file] -d [delimiter] -f [field] ")
	fmt.Println("Extract selected portions of each line")
	PrintColor(Green, "  tee [file]        ")
	fmt.Println("Read from standard input and write to standard output and files")
	PrintColor(Green, "  log [message]     ")
	fmt.Println("Append a log message to a log file")
	PrintColor(Green, "  calc [expression] ")
	fmt.Println("Evaluate a simple arithmetic expression")
	PrintColor(Green, "  truncate [file] -s [size] ")
	fmt.Println("Truncate or extend the size of a file")
	PrintColor(Green, "  du [dir]          ")
	fmt.Println("Estimate file space usage of a directory")
	PrintColor(Green, "  df                ")
	fmt.Println("Report file system disk space usage")
	PrintColor(Green, "  dfi               ")
	fmt.Println("Report file system inode usage")
	PrintColor(Green, "  ln [target] [link] ")
	fmt.Println("Create a symbolic link between files")
	PrintColor(Green, "  tr [set1] [set2]  ")
	fmt.Println("Translate or delete characters in a string")
	PrintColor(Green, "  ping [hostname]   ")
	fmt.Println("Send ICMP ECHO_REQUEST to network hosts")
	PrintColor(Green, "  which [command]   ")
	fmt.Println("Locate a command in the PATH")
	PrintColor(Green, "  ls [dir]          ")
	fmt.Println("List directory contents with detailed file information")
	PrintColor(Green, "  lsc [dir]          ")
	fmt.Println("List directory contents with detailed file information. color-coded output")
	PrintColor(Green, "  stat [file]       ")
	fmt.Println("Display file or file system status")
	PrintColor(Green, "  cal               ")
	fmt.Println("Display a calendar")
	PrintColor(Green, "  source [file]     ")
	fmt.Println("Execute commands from a file")
	PrintColor(Green, "  jobs              ")
	fmt.Println("List background jobs")
	PrintColor(Green, "  fg [job]          ")
	fmt.Println("Bring a background job to the foreground")
	PrintColor(Green, "  bg [job]          ")
	fmt.Println("Send a job to the background")
	PrintColor(Green, "  tree [directory] [-a|--all]")
	fmt.Println("Visualize directory structure as a colorful tree")
	fmt.Println("    -a, --all    Show hidden files and directories")
	fmt.Println("    If no directory is specified, the current directory is used")
	fmt.Println("    Colors indicate directory levels and file types")
	PrintColor(Green, "  watch              ")
	fmt.Println("Runs a specified command periodically and displays its output")
	PrintColor(Green, "  compress [file]   ")
	fmt.Println("Compress a file using gzip")
	PrintColor(Green, "  decompress [file] ")
	fmt.Println("Decompress a file using gzip")
	PrintColor(Green, "  help              ")
	fmt.Println("Show this help message")
	PrintColor(White, "\nPipes:")
	PrintColor(Green, "  Use the '|' character to pipe the output of one command to the input of another.")
	fmt.Println("Example: cat file.txt | grep 'search' | sort")
}
