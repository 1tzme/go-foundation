package markov

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
)

type Prefix []string

func (p Prefix) String() string {
	return strings.Join(p, " ")
}

func (p Prefix) Shift(word string) {
	copy(p, p[1:])
	p[len(p)-1] = word
}

type Chain struct {
	chain     map[string][]string
	prefixLen int
	size      int
}

func NewChain(prefixLen int) *Chain {
	return &Chain{
		chain:     make(map[string][]string),
		prefixLen: prefixLen,
	}
}

func (c *Chain) HasPrefix(p Prefix) bool {
	val, ok := c.chain[p.String()]
	return ok && len(val) > 0
}

func (c *Chain) Size() int {
	return c.size
}

func (c *Chain) Raw() map[string][]string {
	return c.chain
}

func (c *Chain) Build(r io.Reader) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)

	p := make(Prefix, c.prefixLen)
	for scanner.Scan() {
		word := scanner.Text()
		key := p.String()
		c.chain[key] = append(c.chain[key], word)
		p.Shift(word)
		c.size++
	}

	key := p.String()
	_, exists := c.chain[key]
	if !exists {
		c.chain[key] = []string{}
	}
	err := scanner.Err()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Scanner error: ", err)
	}
}

func (c *Chain) Generate(n int, start Prefix) string {
	words := append([]string{}, start...)
	p := make(Prefix, c.prefixLen)
	copy(p, start)

	for i := len(start); i < n; i++ {
		choices := c.chain[p.String()]
		if len(choices) == 0 {
			break
		}

		next := choices[rand.Intn(len(choices))]
		words = append(words, next)
		p.Shift(next)
	}
	return strings.Join(words, " ")
}
