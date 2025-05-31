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

func (p *MyCSVParser) ReadLine(r io.Reader) (string, error) {
	p.buffer = p.buffer[:0]
	temp := [1]byte{}
	for {
		n, err := r.Read(temp[:])
		if n > 0 {
			ch := temp[0]
			if ch == '\n' {
				break
			}
			if ch == '\r' {
				continue
			}
			p.buffer = append(p.buffer, ch)
		}
		if err != nil {
			if err == io.EOF && len(p.buffer) > 0 {
				break
			}
			return "", err
		}
	}

	line := string(p.buffer)
	fields, err := parseCSVLine(line)
	if err != nil {
		p.fields = nil
		return "", err
	}

	p.fields = fields
	return line, nil
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

func parseCSVLine(line string) ([]string, error) {
	fields := []string{}
	field := make([]byte, 0, len(line))
	inQuotes := false

	for i := 0; i < len(line); i++ {
		switch line[i] {
		case '"':
			if inQuotes && i+1 < len(line) && line[i+1] == '"' {
				field = append(field, '"')
				i++
			} else {
				inQuotes = !inQuotes
			}
		case ',':
			if inQuotes {
				field = append(field, line[i])
			} else {
				fields = append(fields, string(field))
				field = field[:0]
			}
		default:
			field = append(field, line[i])
		}
	}

	if inQuotes {
		return nil, ErrQuote
	}

	fields = append(fields, string(field))
	return fields, nil
}
