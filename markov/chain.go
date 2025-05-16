package markov

import (
	"bufio"
	"io"
	"math/rand"
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
}

func NewChain(prefixLen int) *Chain {
	return &Chain{
		chain:     make(map[string][]string),
		prefixLen: prefixLen,
	}
}

func (c *Chain) HasPrefix(p Prefix) bool {
	_, ok := c.chain[p.String()]
	return ok
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
	}
}

func (c *Chain) Generate(n int, start Prefix) string {
	words := append([]string{}, start...)
	p := make(Prefix, c.prefixLen)
	copy(p, start)

	for i := 0; i < n; i++ {
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
