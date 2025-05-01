package main

import (
	"fmt"
	"os"
	"bufio"
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
		os.Exit(1)
	}

	numbers := []string{}
	if args[0] == "--stdin" {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			fields := strings.Fields(line)
			numbers = append(numbers, fields...)
		}

		error := scanner.Err()
		if error != nil {
			fmt.Println("Can not read from stdin")
			os.Exit(1)
		}
	} else {
		numbers = args
	}
	exitCode := 0
	for _, number := range numbers {
		if len(number) < 13 {
			fmt.Println("INCORRECT")
			exitCode = 1
			continue
		}
		if luhnAlgorithm(number) {
			fmt.Println("OK")
			exitCode = 0
		} else {
			fmt.Println("INCORRECT")
			exitCode = 1
		}
	}
	os.Exit(exitCode)
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