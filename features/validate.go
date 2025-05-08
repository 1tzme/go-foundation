package features

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// check if creaditcard num is valid
func Validate(args []string) {
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
	for _, number := range numbers {
		if len(number) < 13 {
			fmt.Println("INCORRECT")
			os.Exit(1)
		}
		if luhnAlgorithm(number) {
			fmt.Println("OK")
		} else {
			fmt.Println("INCORRECT")
			os.Exit(1)
		}
	}
	os.Exit(0)
}
