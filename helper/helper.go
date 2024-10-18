package helper

import (
	"fmt"
	"os"
	"runtime"
)

// ANSI color codes
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Bold   = "\033[1m"
)

var enableColors = true

func init() {
	if runtime.GOOS == "windows" {
		enableColors = true
	}
}

func colorize(text string, colorCode string) string {
	if enableColors {
		return fmt.Sprintf("%s%s%s", colorCode, text, Reset)
	}
	return text
}

func PrintError(errMsg string) {
	fmt.Fprintln(os.Stderr, colorize(errMsg, Red))
	os.Exit(1)
}

func PrintOutput(data any) {
	switch v := data.(type) {
	case string:
		fmt.Fprintln(os.Stdout, colorize(v, Green))
	default:
		fmt.Fprintln(os.Stdout, colorize(fmt.Sprint(v), Green))
	}
}

func PrintDeleted(info string) {
	fmt.Fprintln(os.Stdout, colorize(info, Red))
}

func PrintInfo(info string) {
	fmt.Fprintln(os.Stdout, colorize(info, Yellow))
}

func PrintSuccess(msg string) {
	fmt.Fprintln(os.Stdout, colorize(msg, Blue+Bold))
}
