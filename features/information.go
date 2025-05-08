package features

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// get information about given card
func Information(args []string) {
	if len(args) < 2 {
		fmt.Println("Please enter correct input for information")
		os.Exit(1)
	}

	brandsFile, issuersFile := "", ""
	stdinInput := false
	cardNumbers := []string{}

	for i := 0; i < len(args); i++ {
		if strings.HasPrefix(args[i], "--brands=") {
			brandsFile = strings.TrimPrefix(args[i], "--brands=")
		} else if strings.HasPrefix(args[i], "--issuers=") {
			issuersFile = strings.TrimPrefix(args[i], "--issuers=")
		} else if args[i] == "--stdin" {
			stdinInput = true
		} else {
			cardNumbers = append(cardNumbers, args[i])
		}
	}

	if brandsFile == "" || issuersFile == "" {
		fmt.Println("Missing flags --brands and --issuer")
		os.Exit(1)
	}

	brands, err1 := parseFile(brandsFile)
	issuers, err2 := parseFile(issuersFile)
	if err1 != nil || err2 != nil {
		fmt.Println("Error in reading files")
		os.Exit(1)
	}

	if stdinInput {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			field := strings.Fields(scanner.Text())
			cardNumbers = append(cardNumbers, field...)
		}
	}

	for _, card := range cardNumbers {
		fmt.Println(card)

		if len(card) < 13 || luhnAlgorithm(card) == false {
			fmt.Println("Correct: no")
			fmt.Println("Card Brand: -")
			fmt.Println("Card Issuer: -")
		} else {
			brand := matchData(card, brands)
			issuer := matchData(card, issuers)

			fmt.Println("Correct: yes")
			fmt.Println("Card Brand: " + brand)
			fmt.Println("Card Issuer: " + issuer)
		}
	}
}

// read file and parse to map[name-of-the-kind][]prefixes
func parseFile(fileName string) (map[string][]string, error) {
	res := make(map[string][]string)

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error in reading file")
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ":")
		if len(parts) != 2 {
			continue
		}
		name := parts[0]
		prefix := parts[1]
		res[name] = append(res[name], prefix)
	}

	return res, nil
}

// check if card number prefix matches given data
func matchData(number string, data map[string][]string) string {
	for name, prefixes := range data {
		for _, pref := range prefixes {
			if strings.HasPrefix(number, pref) {
				return name
			}
		}
	}
	return "-"
}
