# CommandRipple - A Simple Shell in Go

CommandRipple is a simple, lightweight command-line shell written in Go. It supports various built-in commands, input/output redirection, piping, background jobs, and much more, making it a powerful tool for developers and system administrators.

## Features

- **Built-in Commands**: CommandRipple supports a wide range of built-in commands like `cd`, `ls`, `echo`, `mkdir`, `rm`, `cp`, `mv`, and many more.
- **Input/Output Redirection**: Redirect command input and output using `<`, `>`, and `>>` for reading from files and writing to or appending to files.
- **Piping**: Chain commands together using the `|` operator to pass the output of one command as input to another.
- **Background Jobs**: Run commands in the background using the `bg` command and manage them with `jobs`, `fg`, and `kill`.
- **Color-Coded Output**: Enhanced `ls` command with color-coded output for better readability.
- **Customizable**: Easily extendable with new commands and features.

## Installation

To install CommandRipple, make sure you have Go installed on your system. Then, clone the repository and build the project:

```bash
git clone https://github.com/genc-murat/commandripple.git
cd commandripple
go build -o commandripple
```

## Usage

To start CommandRipple, simply run the executable:

```bash
./commandripple
```

You will be presented with the `CommandRipple>` prompt, where you can enter commands just like in any other shell.

### Examples

- **Change Directory**:
  ```bash
  cd /path/to/directory
  ```

- **List Files**:
  ```bash
  ls
  ```

- **Run a Command in the Background**:
  ```bash
  bg sleep 30
  ```

- **Show Background Jobs**:
  ```bash
  jobs
  ```

- **Bring a Job to the Foreground**:
  ```bash
  fg 1
  ```

- **Piping and Redirection**:
  ```bash
  ls | grep 'filename' > result.txt
  ```

## Supported Commands

CommandRipple supports the following built-in commands:

- `cd [dir]` - Change the current directory
- `pwd` - Print the current working directory
- `echo [text]` - Echo the input text back to the user
- `clear` - Clear the terminal screen
- `mkdir [dir]` - Create a new directory
- `mkdirp [dir]` - Create directories and parent directories if needed
- `rmdir [dir]` - Remove an empty directory
- `rm [file]` - Remove a file
- `rmrf [dir]` - Recursively remove a directory and its contents
- `cp [src] [dest]` - Copy a file
- `mv [src] [dest]` - Move or rename a file or directory
- `touch [file]` - Create an empty file or update the timestamp
- `touch -t [timestamp] [file]` - Create or update a file with a specific timestamp
- `chmod [permissions] [file]` - Change file permissions
- `chmodr [permissions] [dir]` - Recursively change permissions of a directory
- `cat [file]` - Display the content of a file
- `head [file]` - Display the first few lines of a file
- `tail [file]` - Display the last few lines of a file
- `grep [pattern] [file]` - Search for a pattern in a file
- `find [dir] [name]` - Search for a file or directory by name
- `wc [file]` - Count lines, words, and characters in a file
- `env` - Print environment variables
- `export NAME=VALUE` - Set or modify environment variables
- `history` - Display command history
- `alias name=command` - Create an alias for a command
- `unalias name` - Remove an alias
- `date` - Display the current date and time
- `uptime` - Display how long the shell has been running
- `kill [PID]` - Terminate a process by PID
- `killall [name]` - Kill all processes by name
- `ps` - List currently running processes
- `whoami` - Display the current user's username
- `basename [path]` - Strip directory and suffix from filenames
- `dirname [path]` - Extract the directory path from a full path
- `sort [file]` - Sort lines of a text file
- `uniq [file]` - Remove duplicate lines from a file
- `cut [file] -d [delimiter] -f [field]` - Extract selected portions of each line
- `tee [file]` - Read from standard input and write to standard output and files
- `log [message]` - Append a log message to a log file
- `calc [expression]` - Evaluate a simple arithmetic expression
- `truncate [file] -s [size]` - Truncate or extend the size of a file
- `du [dir]` - Estimate file space usage of a directory
- `df` - Report file system disk space usage
- `dfi` - Report file system inode usage
- `ln [target] [link]` - Create a symbolic link between files
- `tr [set1] [set2]` - Translate or delete characters in a string
- `ping [hostname]` - Send ICMP ECHO_REQUEST to network hosts
- `which [command]` - Locate a command in the PATH
- `ls [dir]` - List directory contents with detailed file information
- `lsc [dir]` - List directory contents with detailed file information and color-coded output
- `stat [file]` - Display file or file system status
- `cal` - Display a calendar
- `source [file]` - Execute commands from a file
- `jobs` - List background jobs
- `fg [job]` - Bring a background job to the foreground
- `bg [job]` - Send a job to the background
- `help` - Show this help message

### Pipes

You can use the `|` character to pipe the output of one command to the input of another:

```bash
cat file.txt | grep 'search' | sort
```

## Contributing

Contributions are welcome! Feel free to submit issues, pull requests, or suggestions.

## License

CommandRipple is released under the MIT License.