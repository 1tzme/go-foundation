package features

import (
	"fmt"
	"os"
	"strings"
	"math/rand"
	"time"
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

		pattern = args[1]
	} else {
		if len(args) != 1 {
			fmt.Println("Too many arguments.")
			os.Exit(1)
		}
		pattern = args[0]
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

	for i := 0; i < maxPossible; i++ {
		endNumber := fmt.Sprintf("%0*d", starsCount, i)
		cardNumber := preNumber + endNumber
		if luhnAlgorithm(cardNumber) {
			validCards = append(validCards, cardNumber)
		}
	}

	if len(validCards) == 0 {
		fmt.Println("No valid card numbers")
		os.Exit(1)
	}

	if pick {
		rand.Seed(time.Now().UnixNano())
		selectedNumber := validCards[rand.Intn(len(validCards))]
		fmt.Println(selectedNumber)
	} else {
		for i := 0; i < len(validCards); i++ {
			fmt.Println(validCards[i])
		}
	}

	os.Exit(0)
}