package main

import (
	"fmt"
	"markov-chain/markov"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	wordLimit, prefixLen, startPref := markov.ParseFlags()

	stat, _ := os.Stdin.Stat()
	if (stat.Mode()&os.ModeCharDevice) != 0 && stat.Size() == 0 {
		fmt.Fprintln(os.Stderr, "Error: no input text")
		os.Exit(1)
	}

	chain := markov.NewChain(prefixLen)
	chain.Build(os.Stdin)

	if chain.Size() == 0 {
		fmt.Println("Not enough length to build chain")
		os.Exit(1)
	}
	if wordLimit > chain.Size() {
		fmt.Println("Word limit exceeds imput size")
		os.Exit(1)
	}

	start := make(markov.Prefix, prefixLen)
	if startPref != "" {
		words := strings.Fields(startPref)
		if len(words) != prefixLen {
			fmt.Println("Starting prefix len mismatch")
			os.Exit(1)
		}
		copy(start, words)

		raw := chain.Raw()
		values, exists := raw[start.String()]
		if !exists {
			fmt.Println("Prefix not found in the text")
			os.Exit(1)
		}
		if values == nil || len(values) == 0 {
			fmt.Println("Prefix can not be continued")
			os.Exit(1)
		}
		if len(start) > wordLimit {
			fmt.Println("-w should be equal or more than prefix length")
			os.Exit(1)
		}
	}

	fmt.Println(chain.Generate(wordLimit, start))
}
