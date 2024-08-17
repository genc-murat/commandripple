package commands

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// File operation command implementations

func Cat(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("'cat' requires an argument")
	}
	for _, file := range args {
		if err := printFileContent(file); err != nil {
			return err
		}
	}
	return nil
}

func Touch(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("'touch' requires an argument")
	}
	for _, file := range args {
		if err := touchFile(file); err != nil {
			return err
		}
	}
	return nil
}

func RemoveFile(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("'rm' requires an argument")
	}
	for _, file := range args {
		if err := os.Remove(file); err != nil {
			return err
		}
	}
	return nil
}

func CopyFile(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("'cp' requires two arguments")
	}
	src := args[0]
	dest := args[1]

	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

func MoveFile(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("'mv' requires two arguments")
	}
	return os.Rename(args[0], args[1])
}

// New file operation commands

func Head(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("'head' requires an argument")
	}
	file := args[0]
	lines := 10 // Default number of lines to display
	if len(args) > 1 {
		fmt.Sscanf(args[1], "%d", &lines)
	}
	return printHead(file, lines)
}

func Tail(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("'tail' requires an argument")
	}
	file := args[0]
	lines := 10 // Default number of lines to display
	if len(args) > 1 {
		fmt.Sscanf(args[1], "%d", &lines)
	}
	return printTail(file, lines)
}

func Grep(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("'grep' requires a pattern and a file")
	}
	pattern := args[0]
	file := args[1]
	return grepPattern(file, pattern)
}

func Find(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("'find' requires a directory and a name")
	}
	directory := args[0]
	name := args[1]
	return findFile(directory, name)
}

func WordCount(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("'wc' requires a file")
	}
	file := args[0]
	return printWordCount(file)
}

func Chmod(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("'chmod' requires permissions and a file")
	}
	// Windows doesn't fully support Unix-style permissions,
	// so this is a placeholder to demonstrate the structure.
	fmt.Println("chmod is not fully supported on Windows.")
	return nil
}

// Helper functions for the new commands

func printHead(filename string, lines int) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for i := 0; i < lines && scanner.Scan(); i++ {
		fmt.Println(scanner.Text())
	}
	return scanner.Err()
}

func printTail(filename string, lines int) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var buffer []string
	for scanner.Scan() {
		buffer = append(buffer, scanner.Text())
		if len(buffer) > lines {
			buffer = buffer[1:]
		}
	}
	for _, line := range buffer {
		fmt.Println(line)
	}
	return scanner.Err()
}

func grepPattern(filename, pattern string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, pattern) {
			fmt.Println(line)
		}
	}
	return scanner.Err()
}

func findFile(directory, name string) error {
	return filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.Contains(info.Name(), name) {
			fmt.Println(path)
		}
		return nil
	})
}

func printWordCount(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	var words, lines, characters int
	for scanner.Scan() {
		words++
		characters += len(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	file.Seek(0, 0)
	scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		lines++
	}

	fmt.Printf(" %d %d %d %s\n", lines, words, characters, filename)
	return nil
}

func printFileContent(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(os.Stdout, file)
	return err
}

func touchFile(filename string) error {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	return file.Close()
}

func Calc(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("'calc' requires an expression")
	}
	expression := strings.Join(args, " ")
	result, err := eval(expression)
	if err != nil {
		return fmt.Errorf("failed to evaluate expression: %v", err)
	}
	fmt.Println(result)
	return nil
}

func Basename(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("'basename' requires a path")
	}
	path := args[0]
	fmt.Println(filepath.Base(path))
	return nil
}

func Dirname(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("'dirname' requires a path")
	}
	path := args[0]
	fmt.Println(filepath.Dir(path))
	return nil
}

func SortFile(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("'sort' requires a file")
	}
	file := args[0]

	lines, err := readLines(file)
	if err != nil {
		return err
	}

	sort.Strings(lines)

	for _, line := range lines {
		fmt.Println(line)
	}
	return nil
}

func Uniq(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("'uniq' requires a file")
	}
	file := args[0]

	lines, err := readLines(file)
	if err != nil {
		return err
	}

	uniqLines := uniq(lines)

	for _, line := range uniqLines {
		fmt.Println(line)
	}
	return nil
}

func Cut(args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("'cut' requires a file, delimiter, and field number")
	}
	file := args[0]
	delimiter := args[1]
	field := 1
	fmt.Sscanf(args[2], "%d", &field)

	lines, err := readLines(file)
	if err != nil {
		return err
	}

	for _, line := range lines {
		fields := strings.Split(line, delimiter)
		if len(fields) >= field {
			fmt.Println(fields[field-1])
		}
	}
	return nil
}

func Tee(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("'tee' requires at least one file")
	}

	files := args
	writers := []io.Writer{os.Stdout}

	for _, file := range files {
		f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		defer f.Close()
		writers = append(writers, f)
	}

	multiWriter := io.MultiWriter(writers...)
	_, err := io.Copy(multiWriter, os.Stdin)
	return err
}

func LogMessage(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("'log' requires a message")
	}
	message := strings.Join(args, " ")
	logFile := "commandripple.log"
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(time.Now().Format(time.RFC1123) + " - " + message + "\n")
	if err != nil {
		return err
	}
	fmt.Println("Log entry added.")
	return nil
}

// Utility functions

func readLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func uniq(lines []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range lines {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func eval(expression string) (float64, error) {
	return strconv.ParseFloat(expression, 64)
}
