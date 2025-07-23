package utils

import (
	"bufio"
	"os"
	"strings"

	"github.com/fatih/color"
)

var (
	MenuTitle  = color.New(color.FgCyan, color.Bold, color.Underline)
	MenuItem   = color.New(color.FgWhite)
	MenuNumber = color.New(color.FgYellow, color.Bold)
	MenuChoice = color.New(color.FgGreen, color.Bold)

	ResultTitle   = color.New(color.FgBlue, color.Bold)
	ResultItem    = color.New(color.FgWhite)
	ResultValue   = color.New(color.FgGreen, color.Bold)
	ResultError   = color.New(color.FgRed, color.Bold)
	ResultSuccess = color.New(color.FgGreen, color.Bold)

	StatusInfo     = color.New(color.FgCyan)
	StatusWarn     = color.New(color.FgYellow)
	StatusError    = color.New(color.FgRed)
	StatusContinue = color.New(color.FgHiBlack)
)

func PrintMenuTitle(title string) {
	MenuTitle.Printf("\n%s\n", title)
}

func PrintMenuItem(number int, label string) {
	MenuNumber.Printf("%d. ", number)
	MenuItem.Printf("%s\n", label)
}

func PrintMenuChoice() {
	MenuChoice.Printf("Choice: ")
}

func PrintResult(title string, values []string) {
	ResultTitle.Printf("%s: ", title)
	if len(values) == 0 {
		ResultValue.Printf("No results\n")
	} else if len(values) == 1 {
		ResultValue.Printf("%s\n", values[0])
	} else {
		ResultValue.Printf("\n%s\n", strings.Join(values, "\n"))
	}
}

func PrintError(message string) {
	ResultError.Printf("Error: %s\n", message)
}

func PrintSuccess(message string) {
	ResultSuccess.Printf("✓ %s\n", message)
}

func PrintInfo(message string) {
	StatusInfo.Printf("ℹ %s\n", message)
}

func PrintWarning(message string) {
	StatusWarn.Printf("⚠ %s\n", message)
}

func PrintPressEnterToContinue() {
	StatusContinue.Printf("\nPress Enter to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
