package parser

import (
	"io"
)

func (p *MyCSVParser) ReadLine(r io.Reader) (string, error) {
	p.buffer = p.buffer[:0]
	temp := [1]byte{}
	isCR := false

	for {
		n, err := r.Read(temp[:])
		if n > 0 {
			ch := temp[0]
			if ch == '\n' && isCR {
				isCR = false
				break
			}
			if ch == '\n' {
				break
			}
			if ch == '\r' {
				isCR = true
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

	if len(p.buffer) == 0 {
		return "", nil
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

func parseCSVLine(line string) ([]string, error) {
	if len(line) == 0 {
		return []string{""}, nil
	}

	fields := []string{}
	field := []byte{}
	inQuotes := false
	quotedField := false

	for i := 0; i < len(line); i++ {
		switch line[i] {
		case '"':
			if !inQuotes {
				inQuotes = true
				quotedField = true
			} else {
				if i+1 < len(line) && line[i+1] == '"' {
					field = append(field, '"')
					i++
				} else {
					inQuotes = false
				}
			}
		case ',':
			if inQuotes {
				field = append(field, line[i])
			} else {
				fields = append(fields, string(trimSpace(field)))
				field = field[:0]
				quotedField = false
			}
		default:
			if !inQuotes && len(field) == 0 && line[i] == ' ' {
				continue
			}
			field = append(field, line[i])
		}
	}

	if inQuotes {
		return nil, ErrQuote
	}
	if quotedField {
		if len(field) > 0 && field[0] == '"' {
			field = field[1:]
		}
		if len(field) > 0 && field[len(field)-1] == '"' {
			field = field[:len(field)-1]
		}
	}

	fields = append(fields, string(trimSpace(field)))
	return fields, nil
}

func trimSpace(s []byte) []byte {
	start := 0
	end := len(s)
	for start < end && s[start] == ' ' {
		start++
	}
	for end > start && s[end-1] == ' ' {
		end--
	}
	return s[start:end]
}
