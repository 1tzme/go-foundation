package features

import (
	"fmt"
	"os"
	"bufio"
	"strings"
)

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

	brands, _ := parseFile(brandsFile)
	issuers, _ := parseFile(issuersFile)

	if stdinInput {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {

		}
	}

	for _, card := range cardNumbers {
		fmt.Println(card)

		
	}
}

// read file and parse to map[name-of-the-kind]number-prefix
func parseFile(fileName string) (map[string]string, error) {
	res := make(map[string]string)

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
		res[parts[1]] = parts[0]
	}

	return res, nil
}