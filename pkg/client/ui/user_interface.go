package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

const (
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorReset  = "\033[0m"
)

/*
	func PromptUserName() string {
		scanner := bufio.NewScanner(os.Stdin)

		scanner.Scan()
		return strings.TrimSpace(scanner.Text())

}
*/
func PromptUserName() string {
	scanner := bufio.NewScanner(os.Stdin)

	if scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		if input != "" {
			return input // Return the non-empty username
		}
	}
	// If input is empty, prompt again
	fmt.Println("Username cannot be empty. Please try again.")
	return ""
}

func DisplayMessage(msg string) {
	fmt.Printf("\r%s\n> ", msg)
}

func DisplayUsers(users []interface{}) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, ColorYellow+"\n**ACTIVE USERS**\t"+ColorReset)
	fmt.Fprintln(w, ColorYellow+"--------\t"+ColorReset)
	for _, user := range users {
		fmt.Fprintf(w, "%s\t\n", user)
	}
	fmt.Fprintln(w, ColorYellow+"--------\t "+ColorReset)
	w.Flush()
	fmt.Print("> ")
}

func DisplayError(text string, exiting bool) {
	if !exiting {
		fmt.Printf("\r%s%s%s\n> ", ColorRed, text, ColorReset)
	}
	fmt.Printf("\r%s%s%s", ColorRed, text, ColorReset)
}

func DisplaySuccess(text string) {
	fmt.Printf("\r%s%s%s\n> ", ColorGreen, text, ColorReset)
}

func DisplayExiting(text string) {
	fmt.Printf("\r%s%s%s\n> ", ColorBlue, text, ColorReset)
}

func DisplayNeutral(text string) {
	fmt.Printf("\r%s%s%s\n> ", ColorYellow, text, ColorReset)
}
