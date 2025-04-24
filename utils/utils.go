package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ReadInput reads user input from the terminal
func ReadInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// ValidateInput checks if the input string is not empty
func ValidateInput(input string) bool {
	return len(input) > 0
}

// PrintMessage prints a message to the terminal
func PrintMessage(message string) {
	fmt.Println(message)
}

// ClearScreen clears the terminal screen
func ClearScreen() {
	fmt.Print("\033[H\033[2J")
}
