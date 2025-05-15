package main

import (
	"fmt"
	"os"
	"flag"
	"markov-chain/markov"
)

func usage() {
	fmt.Println(`Markov Chain text generator.
	
Usage:
	markovchain [-w <N>] [-p <S>] [-l <N>]
	markovchain --help

Options:
	--help  Show this screen.
	-w N    Number of maximum words
	-p S    Starting prefix
	-l N    Prefix length`)
	os.Exit(0)
}

func main() {
	numWords := flag.Int("w", 100, "maximum number of words to print")
	prefixLen := flag.Int("p", 2, "prefix length in words")
	flag.Parse()
	c := markov.NewChain(*prefixLen)
	c.Build(os.Stdin)
	text := c.Generate(*numWords)
	fmt.Println(text)
}