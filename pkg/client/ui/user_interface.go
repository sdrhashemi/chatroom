package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

func PromptUserName() string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter your username (must be unique): ")
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())

}

func DisplayMessage(msg string) {
	fmt.Printf("\r%s\n> ", msg)
}

func DisplayUsers(users []interface{}) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "**USERNAMES**\t")
	fmt.Fprintln(w, "--------\t")
	for _, user := range users {
		fmt.Fprintf(w, "%s\t\n", user)
	}
	fmt.Fprintln(w, "--------\t")
	w.Flush()
	fmt.Print("> ")
}
