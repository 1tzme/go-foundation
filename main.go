package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Missing command")
		os.Exit(1)
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
	var numbers []string

}
