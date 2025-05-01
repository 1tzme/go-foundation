package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Missing command")
		os.Exit(0)
	}

	command := os.Args[1]
	args := os.Args[2:]

	if command == "validate" {
		validate(args)
	} else {
		fmt.Println("Unknown command: ", command)
	}
}

// check if creaditcard num is valid
func validate(args []string) {
	if len(args) == 0 {
		fmt.Println("Please enter creditcard number.")
	}

	numbers := []string{}
	if args[0] == "--stdin" {
		input := make([]byte, 0)
		n, error := os.Stdin.Read(input)
		if error != nil {
			fmt.Println("Can not read from stdin")
		}
		text := string(input[:n])
		numbers = strings.Fields(text)
	} else {
		numbers = args
	}

	for _, number := range numbers {
		if len(number) < 13 {
			fmt.Println("INCORRECT")
			continue
		}
		if luhnAlgorithm(number) {
			fmt.Println("OK")
		} else {
			fmt.Println("INCORRECT")
		}
	}
}

// luhn algorithm for creaditcard num checking
func luhnAlgorithm(number string) bool {
    sum := 0
    isSecond := false

    for i := len(number) - 1; i >= 0; i-- {
        digit := int(number[i] - '0')
        if isSecond {
            digit *= 2
        }
        sum += digit / 10
        sum += digit % 10

        isSecond = !isSecond
    }
	if sum%10 == 0 {
		return true
	}
	return false
}