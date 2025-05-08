package features

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// card issuance by specific brand and issuer
func Issue(args []string) {
	if len(args) != 4 {
		fmt.Println("Number of arguments should be 4")
		os.Exit(1)
	}

	brandsFile, issuersFile, brand, issuer := "", "", "", ""

	for i := 0; i < 4; i++ {
		if strings.HasPrefix(args[i], "--brands=") {
			brandsFile = strings.TrimPrefix(args[i], "--brands=")
		} else if strings.HasPrefix(args[i], "--issuers=") {
			issuersFile = strings.TrimPrefix(args[i], "--issuers=")
		} else if strings.HasPrefix(args[i], "--brand=") {
			brand = strings.TrimPrefix(args[i], "--brand=")
		} else if strings.HasPrefix(args[i], "--issuer=") {
			issuer = strings.TrimPrefix(args[i], "--issuer=")
		} else {
			fmt.Println("Invalid command: " + args[i])
			os.Exit(1)
		}
	}

	if brandsFile == "" || issuersFile == "" || brand == "" || issuer == "" {
		fmt.Println("Missing flag")
		os.Exit(1)
	}

	brandsMap, err1 := parseFile(brandsFile)
	issuersMap, err2 := parseFile(issuersFile)
	if err1 != nil || err2 != nil {
		fmt.Println("Error in reading files")
		os.Exit(1)
	}

	brandsPrefs, brandsSuccess := brandsMap[brand]
	issuersPrefs, issuersSuccess := issuersMap[issuer]
	if !brandsSuccess || len(brandsPrefs) == 0 {
		fmt.Println("Error in reading brands")
		os.Exit(1)
	}
	if !issuersSuccess || len(issuersPrefs) == 0 {
		fmt.Println("Error in reading issuers")
		os.Exit(1)
	}

	rand.Seed(time.Now().UnixNano())
	brandPref := brandsPrefs[rand.Intn(len(brandsPrefs))]
	issuerPref := issuersPrefs[rand.Intn(len(issuersPrefs))]

	cardLen := 0
	switch brand {
	case "VISA":
		cardLen = 16
	case "MASTERCARD":
		cardLen = 16
	case "AMEX":
		cardLen = 15
	default:
		fmt.Println("Unsupported card length")
		os.Exit(1)
	}

	base := issuerPref + brandPref
	remainingDigits := cardLen - len(base) - 1

	for {
		card := base
		for i := 0; i < remainingDigits; i++ {
			card += strconv.Itoa(rand.Intn(10))
		}

		for i := 0; i < 10; i++ {
			full := card + strconv.Itoa(i)
			if luhnAlgorithm(full) {
				fmt.Println(full)
				os.Exit(0)
			}
		}
	}
}
