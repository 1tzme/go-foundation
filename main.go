package main

import (
	"fmt"
	"os"
	"creditcard/features"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Missing command")
		os.Exit(0)
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "validate":
		features.Validate(args)
	default:
		fmt.Println("Unknown comman: ", command)
	}
}