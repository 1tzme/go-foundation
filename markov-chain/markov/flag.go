package markov

import (
	"flag"
	"fmt"
	"os"
)

func ParseFlags() (int, int, string) {
	wordLimit := flag.Int("w", 100, "Num of words to generate (1-10000)")
	prefixLen := flag.Int("l", 2, "Prefix length (1-5)")
	startPref := flag.String("p", "", "Starting prefix")
	help := flag.Bool("help", false, "Show usage")

	flag.Usage = usage
	flag.Parse()

	if *help {
		flag.Usage()
	}
	if *wordLimit <= 0 || *wordLimit > 10000 {
		fmt.Println("Number of words should be between 1 and 10000")
		os.Exit(1)
	}
	if *prefixLen < 1 || *prefixLen > 5 {
		fmt.Println("Prefix len should be between 1 and 5")
		os.Exit(1)
	}
	return *wordLimit, *prefixLen, *startPref
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
