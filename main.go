package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Missing command")
	}

	command := os.Args[1]
	args := os.Args[2:]

	if command == "validate" {
		validate(args)
	} else {
		fmt.Println("Unknown command: ", command)
	}
}


func validate(args []string) {
	if len(args) == 0 {
		fmt.Println("Please enter creditcard number.")
	}

	numbers := []string{}
	if args[0] == "--stdin" {
		d := make([]byte, 0)
		n, error := os.Stdin.Read(d)
		if error != nil {
			fmt.Println("Can not read from stdin")
		}

		input := string(d[:n])
		numbers = strings.Fields(input)
	} else {
		numbers = args
	}
}

func luhnAlgorithm(number string) bool {
	sum := 0
	doubleDigit := false

	for i := len(number); i >= 0; i-- {
		digit, error := strconv.Atoi(string(number[i]))
		if error != nil {
			return false
		}
		if doubleDigit {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		doubleDigit = !doubleDigit
	}

	if sum%10 == 0 {
		return true
	}
	return false
}