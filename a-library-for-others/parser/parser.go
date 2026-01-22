package parser

import (
	"io"
)

type CSVParser interface {
	ReadLine(r io.Reader) (string, error)
	GetField(n int) (string, error)
	GetNumberOfFields() int
}

type MyCSVParser struct {
	buffer []byte
	fields []string
}

func (p *MyCSVParser) GetField(n int) (string, error) {
	if n < 0 || n >= len(p.fields) {
		return "", ErrFieldCount
	}
	return p.fields[n], nil
}

func (p *MyCSVParser) GetNumberOfFields() int {
	return len(p.fields)
}
