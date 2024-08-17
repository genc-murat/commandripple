package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	ColorReset   = "\033[0m"
	ColorBlue    = "\033[34m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorCyan    = "\033[36m"
	ColorRed     = "\033[31m"
	ColorMagenta = "\033[35m"
)

var LineColors = []string{ColorBlue, ColorGreen, ColorYellow, ColorCyan, ColorRed, ColorMagenta}

type TreeStats struct {
	Directories int
	Files       int
}

type TreeOptions struct {
	ShowHidden bool
}

func Tree(args []string) error {
	options := TreeOptions{}
	root := "."

	// Parse arguments
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-a", "--all":
			options.ShowHidden = true
		default:
			root = args[i]
		}
	}

	fmt.Printf("Starting tree from root: %s\n", root)

	stats := &TreeStats{}
	err := printTree(root, "", stats, options, 0)
	if err != nil {
		return fmt.Errorf("error in printTree: %v", err)
	}
	fmt.Printf("\n%d directories, %d files\n", stats.Directories, stats.Files)
	return nil
}

func printTree(path string, prefix string, stats *TreeStats, options TreeOptions, depth int) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("error getting file info for %s: %v", path, err)
	}

	if !options.ShowHidden && isHidden(fileInfo.Name()) {
		return nil
	}

	lineColor := LineColors[depth%len(LineColors)]
	coloredPrefix := colorizePrefix(prefix, lineColor)

	if fileInfo.IsDir() {
		fmt.Printf("%s%s%s%s\n", coloredPrefix, ColorBlue, fileInfo.Name(), ColorReset)
		stats.Directories++
	} else {
		fmt.Printf("%s%s%s%s\n", coloredPrefix, getFileColor(fileInfo), fileInfo.Name(), ColorReset)
		stats.Files++
	}

	if !fileInfo.IsDir() {
		return nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("error reading directory %s: %v", path, err)
	}

	for i, entry := range entries {
		if !options.ShowHidden && isHidden(entry.Name()) {
			continue
		}

		newPrefix := prefix + "├── "
		if i == len(entries)-1 {
			newPrefix = prefix + "└── "
		}
		err := printTree(filepath.Join(path, entry.Name()), newPrefix, stats, options, depth+1)
		if err != nil {
			fmt.Printf("Error processing %s: %v\n", entry.Name(), err)
		}
	}

	return nil
}

func colorizePrefix(prefix string, color string) string {
	parts := strings.Split(prefix, "──")
	if len(parts) > 1 {
		coloredParts := make([]string, len(parts))
		for i, part := range parts {
			if i == len(parts)-1 {
				coloredParts[i] = color + part + ColorReset
			} else {
				coloredParts[i] = color + part + "──" + ColorReset
			}
		}
		return strings.Join(coloredParts, "")
	}
	return prefix
}

func getFileColor(fileInfo os.FileInfo) string {
	switch {
	case fileInfo.Mode()&os.ModeSymlink != 0:
		return ColorCyan
	case fileInfo.Mode()&0111 != 0:
		return ColorGreen
	case isHidden(fileInfo.Name()):
		return ColorYellow
	default:
		return ColorReset
	}
}

func isHidden(name string) bool {
	return strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_")
}
