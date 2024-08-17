package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func Diff(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("usage: diff [file1] [file2]")
	}

	file1, err := os.ReadFile(args[0])
	if err != nil {
		return fmt.Errorf("error reading file1: %v", err)
	}

	file2, err := os.ReadFile(args[1])
	if err != nil {
		return fmt.Errorf("error reading file2: %v", err)
	}

	dmp := diffmatchpatch.New()

	diffs := dmp.DiffMain(string(file1), string(file2), false)

	fmt.Printf("Differences between %s and %s:\n\n", args[0], args[1])

	lineNum1, lineNum2 := 1, 1
	for _, diff := range diffs {
		switch diff.Type {
		case diffmatchpatch.DiffEqual:
			lines := strings.Split(diff.Text, "\n")
			lineNum1 += len(lines) - 1
			lineNum2 += len(lines) - 1
		case diffmatchpatch.DiffInsert:
			lines := strings.Split(diff.Text, "\n")
			for _, line := range lines {
				if line != "" {
					fmt.Printf("\033[32m+ %d: %s\033[0m\n", lineNum2, line)
					lineNum2++
				}
			}
		case diffmatchpatch.DiffDelete:
			lines := strings.Split(diff.Text, "\n")
			for _, line := range lines {
				if line != "" {
					fmt.Printf("\033[31m- %d: %s\033[0m\n", lineNum1, line)
					lineNum1++
				}
			}
		}
	}

	return nil
}
