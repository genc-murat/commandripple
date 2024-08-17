package commands

import (
	"fmt"
)

// Color codes
const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"
)

// Print in color
func PrintColor(color, text string) {
	fmt.Println(color + text + Reset)
}

func PrintColorInline(color, text string) {
	fmt.Print(color + text + Reset)
}
