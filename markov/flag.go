package markov

import (
	"fmt"
	"os"
	"strconv"
)

func ParseFlags() (int, int, string) {
	wordLimit := 100
	prefixLen := 2
	startPref := ""

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--help":
			usage()
		case "-w":
			i++
			if i >= len(args) {
				fmt.Println("Missing value for -w")
				os.Exit(1)
			}
			val, err := strconv.Atoi(args[i])
			if err != nil || val <= 0 || val > 10000 {
				fmt.Println("Number of words should be between 1 and 10000")
				os.Exit(1)
			}
			wordLimit = val
		case "-l":
			i++
			if i >= len(args) {
				fmt.Println("Missing value for -l")
				os.Exit(1)
			}
			val, err := strconv.Atoi(args[i])
			if err != nil || val < 1 || val > 5 {
				fmt.Println("Prefix len should be between 1 and 5")
				os.Exit(1)
			}
			prefixLen = val
		case "-p":
			i++
			if i >= len(args) {
				fmt.Println("Missing value for -p")
				os.Exit(1)
			}
			startPref = args[i]
		default:
			fmt.Println("Unknown flag: " + args[i])
			os.Exit(1)
		}
	}

	return wordLimit, prefixLen, startPref
}

func usage() {
	fmt.Println(`
Markov Chain text generator.

Usage:
  markovchain [-w <N>] [-p <S>] [-l <N>]
  markovchain --help

Options:
  --help  Show this screen.
  -w N    Number of maximum words
  -p S    Starting prefix
  -l N    Prefix length
	`)
	os.Exit(0)
}
