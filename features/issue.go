package features

import (
	"fmt"
	"os"
	"strings"
)

func Issue(args []string) {
	if len(args) != 4 {
		fmt.Println("Number of arguments should be 4")
		os.Exit(1)
	}

	brandsFile, issuersFile, brand, issuer := "", "", "", ""

	for i := 0; i < 4; i ++ {
		if strings.HasPrefix(args[i], "--brands=") {
			brandsFile = strings.TrimPrefix(args[i], "--brands=")
		} else if strings.HasPrefix(args[i], "--issuers=") {
			issuersFile = strings.TrimPrefix(args[i], "--issuers=")
		} else if strings.HasPrefix(args[i], "--brand=") {
			brand = strings.TrimPrefix(args[i], "--brand=")
		} else if strings.HasPrefix(args[i], "--issuer=") {
			issuer = strings.TrimPrefix(args[i], "--issuer=")
		}
	}

	if brandsFile == "" || issuersFile == "" || brand == "" || issuer == "" {
		fmt.Println("Missing flag")
		os.Exit(1)
	}

	
}