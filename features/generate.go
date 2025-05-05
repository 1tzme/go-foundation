package features

import (
	"fmt"
	"os"
	"strings"
)

func Generate(args []string) {
	if len(args) == 0 {
		fmt.Println("Please enter creditcard number.")
		os.Exit(1)
	}

	pick := false
	pattern := ""
	
	if args[0] == "--pick" {
		pick = true

		if len(args) < 2 {
			fmt.Println("Please enter creditcard number.")
			os.Exit(1)
		}
		if len(args) > 2 {
			fmt.Println("Too many arguments")
			os.Exit(1)
		}

		pattern = args[0]
	} else {
		if len(args) != 1 {
			fmt.Println("Too many arguments.")
			os.Exit(1)
		}
	}

	starsCount := strings.Count(pattern, "*")
	if starsCount == 0 || starsCount > 4 {
		fmt.Println("The number of * should be between 1 and 4")
		os.Exit(1)
	}

	preNumber := pattern[:len(pattern)-starsCount]
	validCards := []string{}

	maxPossible := 1
	for i := 0; i < starsCount; i++ {
		maxPossible *= 10
	}
}