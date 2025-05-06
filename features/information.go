package features

import (
	"fmt"
	"os"
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

}