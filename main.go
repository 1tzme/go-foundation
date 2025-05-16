package main

import (
	"fmt"
	"markov-chain/markov"
	"os"
	"strings"
)

func main() {
	wordLimit, prefixLen, startPref := markov.ParseFlags()

	stat, _ := os.Stdin.Stat()
	if (stat.Mode()&os.ModeCharDevice) != 0 && stat.Size() == 0 {
		fmt.Fprintln(os.Stderr, "Error: no input text")
		os.Exit(1)
	}

	chain := markov.NewChain(prefixLen)
	chain.Build(os.Stdin)

	start := make(markov.Prefix, prefixLen)
	if startPref != "" {
		words := strings.Fields(startPref)
		if len(words) != prefixLen {
			fmt.Println("Starting prefix len mismatch")
			os.Exit(1)
		}
		copy(start, words)

		if !chain.HasPrefix(start) {
			fmt.Println("Prefix not found in the text")
			os.Exit(1)
		}
	}

	fmt.Println(chain.Generate(wordLimit, start))
}
